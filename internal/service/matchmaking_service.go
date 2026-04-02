package service

import (
	"context"
	"fmt"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
)

type matchmakingService struct {
	matchRepo domain.MatchmakingRepository
	userRepo  domain.UserRepository
	roomRepo  domain.RoomRepository
}

func NewMatchmakingService(
	matchRepo domain.MatchmakingRepository,
	userRepo domain.UserRepository,
	roomRepo domain.RoomRepository,
) domain.MatchmakingService {
	return &matchmakingService{
		matchRepo: matchRepo,
		userRepo:  userRepo,
		roomRepo:  roomRepo,
	}
}

// FindMatch attempts to pair the user with a waiting partner.
// If a partner is found → creates a room and returns MatchResult.
// If no partner is found → enqueues the user and returns nil (WebSocket will notify later).
func (s *matchmakingService) FindMatch(ctx context.Context, userID uuid.UUID, questions []string) (*domain.MatchResult, error) {
	// 1. Get the searching user's profile
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.TargetLanguage == "" || user.ProficiencyLevel == "" {
		return nil, domain.ErrProfileIncomplete
	}

	if !domain.IsValidTargetLanguage(user.TargetLanguage) {
		return nil, domain.ErrInvalidLanguage
	}

	// 2. Try to dequeue a waiting partner from the same queue
	partnerReq, err := s.matchRepo.Dequeue(ctx, user.TargetLanguage, user.ProficiencyLevel)
	if err != nil {
		return nil, err
	}

	// 3. No partner waiting — enqueue ourselves and wait
	if partnerReq == nil {
		enqueueReq := domain.MatchRequest{
			UserID:           userID,
			TargetLanguage:   user.TargetLanguage,
			ProficiencyLevel: user.ProficiencyLevel,
			Questions:        questions,
		}
		if err := s.matchRepo.Enqueue(ctx, enqueueReq); err != nil {
			return nil, err
		}
		return nil, nil // Caller (Hub) knows nil means "waiting"
	}

	// 4. Partner found! Don't match with ourselves (safety check)
	if partnerReq.UserID == userID {
		// Re-enqueue the partner and return nil
		_ = s.matchRepo.Enqueue(ctx, *partnerReq)
		return nil, nil
	}

	// 5. Get partner's user data
	partner, err := s.userRepo.GetByID(ctx, partnerReq.UserID)
	if err != nil {
		// Partner disappeared — re-try by recursing
		return s.FindMatch(ctx, userID, questions)
	}

	// 6. Create the room
	channelName := fmt.Sprintf("room_%s", uuid.New().String())
	room := &domain.Room{
		ID:               uuid.New(),
		Type:             domain.RoomTypeOneOnOne,
		Status:           domain.RoomStatusActive,
		TargetLanguage:   user.TargetLanguage,
		ProficiencyLevel: user.ProficiencyLevel,
		MaxParticipants:  2,
		AgoraChannelName: channelName,
	}
	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	// 7. Add both users as participants
	for _, uid := range []uuid.UUID{userID, partner.ID} {
		participant := &domain.RoomParticipant{
			ID:     uuid.New(),
			RoomID: room.ID,
			UserID: uid,
			Role:   domain.ParticipantRoleMember,
		}
		if err := s.roomRepo.AddParticipant(ctx, participant); err != nil {
			return nil, err
		}
	}

	// 8. Return the result (from searcher's perspective)
	return &domain.MatchResult{
		RoomID:           room.ID,
		AgoraChannelName: channelName,
		PartnerID:        partner.ID,
		PartnerUsername:  partner.Username,
		PartnerQuestions: partnerReq.Questions,
		MyUsername:       user.Username,
		MyQuestions:      questions,
	}, nil
}

// CancelMatch removes the user from the matchmaking queue.
// Called when user explicitly cancels or when WebSocket disconnects.
func (s *matchmakingService) CancelMatch(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	return s.matchRepo.Remove(ctx, userID, user.TargetLanguage, user.ProficiencyLevel)
}
