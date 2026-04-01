package service

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
)

type messageService struct {
	messageRepo domain.MessageRepository
	roomRepo    domain.RoomRepository
}

func NewMessageService(mr domain.MessageRepository, rr domain.RoomRepository) domain.MessageService {
	return &messageService{
		messageRepo: mr,
		roomRepo:    rr,
	}
}

func (s *messageService) SaveMessage(ctx context.Context, roomID, senderID uuid.UUID, content string) (*domain.Message, error) {
	// Optional: verify sender is in room
	inRoom, err := s.roomRepo.IsUserInRoom(ctx, roomID, senderID)
	if err != nil {
		return nil, err
	}
	if !inRoom {
		return nil, domain.ErrAlreadyInRoom // Or create ErrNotInRoom
	}

	msg := &domain.Message{
		ID:       uuid.New(),
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}

	if err := s.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *messageService) GetChatHistory(ctx context.Context, roomID uuid.UUID, page, pageSize int) (*domain.ChatHistoryResponse, error) {
	offset := (page - 1) * pageSize

	messages, count, err := s.messageRepo.GetByRoomID(ctx, roomID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	return &domain.ChatHistoryResponse{
		Messages:   messages,
		TotalCount: count,
	}, nil
}
