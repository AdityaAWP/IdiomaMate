package service

import (
	"context"
	"fmt"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
)

type roomService struct {
	roomRepo domain.RoomRepository
	notifier domain.NotificationService
}

func NewRoomService(rr domain.RoomRepository, notifier domain.NotificationService) domain.RoomService {
	return &roomService{
		roomRepo: rr,
		notifier: notifier,
	}
}

func (s *roomService) CreateLobby(ctx context.Context, masterID uuid.UUID, req domain.CreateLobbyRequest) (*domain.Room, error) {
	channelName := fmt.Sprintf("lobby_%s", uuid.New().String())

	room := &domain.Room{
		ID:               uuid.New(),
		Type:             domain.RoomTypeLobby,
		Status:           domain.RoomStatusWaiting,
		TargetLanguage:   req.TargetLanguage,
		ProficiencyLevel: req.ProficiencyLevel,
		MaxParticipants:  5,
		Title:            req.Title,
		AgoraChannelName: channelName,
	}

	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	// Add creator as master with ACCEPTED status
	participant := &domain.RoomParticipant{
		ID:     uuid.New(),
		RoomID: room.ID,
		UserID: masterID,
		Role:   domain.ParticipantRoleGroupMaster,
		Status: domain.JoinStatusAccepted,
	}
	if err := s.roomRepo.AddParticipant(ctx, participant); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) GetRoom(ctx context.Context, roomID uuid.UUID) (*domain.Room, error) {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, domain.ErrRoomNotFound
	}
	return room, nil
}

func (s *roomService) ListLobbies(ctx context.Context, targetLanguage, proficiencyLevel string, page, pageSize int) (*domain.LobbyListResponse, error) {
	offset := (page - 1) * pageSize
	rooms, count, err := s.roomRepo.ListActiveLobbies(ctx, targetLanguage, proficiencyLevel, offset, pageSize)
	if err != nil {
		return nil, err
	}

	var items []domain.LobbyItem
	for _, r := range rooms {
		masterUsername := "Unknown"
		for _, p := range r.Participants {
			if p.Role == domain.ParticipantRoleGroupMaster && p.User != nil {
				masterUsername = p.User.Username
				break
			}
		}

		items = append(items, domain.LobbyItem{
			ID:               r.ID,
			Title:            r.Title,
			TargetLanguage:   r.TargetLanguage,
			ProficiencyLevel: r.ProficiencyLevel,
			CurrentCount:     len(r.Participants),
			MaxParticipants:  r.MaxParticipants,
			MasterUsername:   masterUsername,
			CreatedAt:        r.CreatedAt,
		})
	}

	return &domain.LobbyListResponse{
		Rooms:      items,
		TotalCount: count,
	}, nil
}

func (s *roomService) RequestJoin(ctx context.Context, roomID, userID uuid.UUID) error {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return domain.ErrRoomNotFound
	}

	if room.Status == domain.RoomStatusClosed {
		return domain.ErrRoomClosed
	}

	inRoom, err := s.roomRepo.IsUserInRoom(ctx, roomID, userID)
	if err == nil && inRoom {
		return domain.ErrAlreadyInRoom
	}

	count, err := s.roomRepo.CountParticipants(ctx, roomID)
	if err != nil {
		return err
	}
	if count >= room.MaxParticipants {
		return domain.ErrRoomFull
	}

	participant := &domain.RoomParticipant{
		ID:     uuid.New(),
		RoomID: roomID,
		UserID: userID,
		Role:   domain.ParticipantRoleMember,
		Status: domain.JoinStatusPending,
	}

	if err := s.roomRepo.AddParticipant(ctx, participant); err != nil {
		return err
	}

	// Notify the master
	master, err := s.roomRepo.GetMaster(ctx, roomID)
	if err == nil && master != nil && s.notifier != nil {
		s.notifier.NotifyUser(master.UserID, domain.WSTypeJoinRequest, map[string]interface{}{
			"room_id": roomID,
			"user_id": userID,
			"message": "Someone wants to join your room",
		})
	}

	return nil
}

func (s *roomService) RespondJoinRequest(ctx context.Context, roomID, masterID, targetUserID uuid.UUID, accept bool) error {
	master, err := s.roomRepo.GetMaster(ctx, roomID)
	if err != nil || master.UserID != masterID {
		return domain.ErrNotRoomMaster
	}

	participant, err := s.roomRepo.GetParticipant(ctx, roomID, targetUserID)
	if err != nil || participant.Status != domain.JoinStatusPending {
		return fmt.Errorf("no pending request found for this user")
	}

	var status domain.JoinStatus
	var wsType domain.WSMessageType
	var message string

	if accept {
		status = domain.JoinStatusAccepted
		wsType = domain.WSTypeJoinApproved
		message = "Your request to join the room has been approved."
	} else {
		status = domain.JoinStatusDeclined
		wsType = domain.WSTypeJoinRejected
		message = "Your request to join the room was declined."
	}

	if err := s.roomRepo.UpdateParticipantStatus(ctx, roomID, targetUserID, status); err != nil {
		return err
	}

	if s.notifier != nil {
		s.notifier.NotifyUser(targetUserID, wsType, map[string]interface{}{
			"room_id": roomID,
			"message": message,
		})
	}

	return nil
}

func (s *roomService) LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return domain.ErrRoomNotFound
	}

	// Check if user is master
	isMaster := false
	for _, p := range room.Participants {
		if p.UserID == userID && p.Role == domain.ParticipantRoleGroupMaster {
			isMaster = true
			break
		}
	}

	// If master leaves, the room is closed
	if isMaster {
		return s.roomRepo.Close(ctx, roomID)
	}

	return s.roomRepo.RemoveParticipant(ctx, roomID, userID)
}

func (s *roomService) CloseRoom(ctx context.Context, roomID, masterID uuid.UUID) error {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return domain.ErrRoomNotFound
	}

	isMaster := false
	for _, p := range room.Participants {
		if p.UserID == masterID && p.Role == domain.ParticipantRoleGroupMaster {
			isMaster = true
			break
		}
	}

	if !isMaster {
		return domain.ErrNotRoomMaster
	}

	return s.roomRepo.Close(ctx, roomID)
}

func (s *roomService) KickUser(ctx context.Context, roomID, masterID, targetUserID uuid.UUID) error {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return domain.ErrRoomNotFound
	}

	isMaster := false
	for _, p := range room.Participants {
		if p.UserID == masterID && p.Role == domain.ParticipantRoleGroupMaster {
			isMaster = true
			break
		}
	}

	if !isMaster {
		return domain.ErrNotRoomMaster
	}

	if err := s.roomRepo.RemoveParticipant(ctx, roomID, targetUserID); err != nil {
		return err
	}

	if s.notifier != nil {
		s.notifier.NotifyUser(targetUserID, domain.WSTypeUserKicked, map[string]interface{}{
			"room_id": roomID,
			"message": "You have been kicked from the room by the master.",
		})
	}

	return nil
}
