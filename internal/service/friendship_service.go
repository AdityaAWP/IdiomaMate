package service

import (
	"context"
	"errors"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type friendshipService struct {
	friendRepo domain.FriendshipRepository
	userRepo   domain.UserRepository
	notifier   domain.NotificationService
}

func NewFriendshipService(
	fr domain.FriendshipRepository,
	ur domain.UserRepository,
	notifier domain.NotificationService,
) domain.FriendshipService {
	return &friendshipService{
		friendRepo: fr,
		userRepo:   ur,
		notifier:   notifier,
	}
}

func (s *friendshipService) SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) (*domain.Friendship, error) {
	// Cannot friend yourself
	if senderID == receiverID {
		return nil, domain.ErrCannotFriendSelf
	}

	// Verify receiver exists
	_, err := s.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Check for existing friendship in any direction
	existing, err := s.friendRepo.GetBetweenUsers(ctx, senderID, receiverID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existing != nil {
		switch existing.Status {
		case domain.FriendshipStatusAccepted:
			return nil, domain.ErrAlreadyFriends
		case domain.FriendshipStatusPending:
			return nil, domain.ErrFriendRequestExists
		case domain.FriendshipStatusDeclined:
			// Allow re-sending after a decline — update the old row
			existing.UserID1 = senderID
			existing.UserID2 = receiverID
			existing.Status = domain.FriendshipStatusPending
			if err := s.friendRepo.UpdateStatus(ctx, existing.ID, domain.FriendshipStatusPending); err != nil {
				return nil, err
			}
			// Notify the receiver in real time
			if s.notifier != nil {
				s.notifier.NotifyUser(receiverID, domain.WSTypeFriendRequest, map[string]interface{}{
					"friendship_id": existing.ID,
					"from_user_id":  senderID,
					"message":       "You have a new friend request!",
				})
			}
			return existing, nil
		}
	}

	// Create new friendship
	friendship := &domain.Friendship{
		ID:      uuid.New(),
		UserID1: senderID,
		UserID2: receiverID,
		Status:  domain.FriendshipStatusPending,
	}

	if err := s.friendRepo.Create(ctx, friendship); err != nil {
		return nil, err
	}

	// Notify receiver via WebSocket
	if s.notifier != nil {
		s.notifier.NotifyUser(receiverID, domain.WSTypeFriendRequest, map[string]interface{}{
			"friendship_id": friendship.ID,
			"from_user_id":  senderID,
			"message":       "You have a new friend request!",
		})
	}

	return friendship, nil
}

func (s *friendshipService) RespondToRequest(ctx context.Context, friendshipID, userID uuid.UUID, accept bool) error {
	friendship, err := s.friendRepo.GetByID(ctx, friendshipID)
	if err != nil {
		return domain.ErrFriendshipNotFound
	}

	// Only the receiver (UserID2) can respond
	if friendship.UserID2 != userID {
		return domain.ErrForbidden
	}

	if friendship.Status != domain.FriendshipStatusPending {
		return domain.ErrBadRequest
	}

	var newStatus domain.FriendshipStatus
	if accept {
		newStatus = domain.FriendshipStatusAccepted
	} else {
		newStatus = domain.FriendshipStatusDeclined
	}

	if err := s.friendRepo.UpdateStatus(ctx, friendshipID, newStatus); err != nil {
		return err
	}

	// Notify the original sender about the response
	if accept && s.notifier != nil {
		s.notifier.NotifyUser(friendship.UserID1, domain.WSTypeFriendAccepted, map[string]interface{}{
			"friendship_id": friendshipID,
			"friend_id":     userID,
			"message":       "Your friend request was accepted!",
		})
	}

	return nil
}

func (s *friendshipService) ListFriends(ctx context.Context, userID uuid.UUID) ([]domain.Friendship, error) {
	return s.friendRepo.ListFriends(ctx, userID)
}

func (s *friendshipService) ListPendingRequests(ctx context.Context, userID uuid.UUID) ([]domain.Friendship, error) {
	return s.friendRepo.ListPendingRequests(ctx, userID)
}
