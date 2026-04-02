package service

import (
	"context"
	"errors"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type directMessageService struct {
	dmRepo     domain.DirectMessageRepository
	friendRepo domain.FriendshipRepository
	notifier   domain.NotificationService
}

func NewDirectMessageService(
	dmr domain.DirectMessageRepository,
	fr domain.FriendshipRepository,
	notifier domain.NotificationService,
) domain.DirectMessageService {
	return &directMessageService{
		dmRepo:     dmr,
		friendRepo: fr,
		notifier:   notifier,
	}
}

func (s *directMessageService) SendDM(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*domain.DirectMessage, error) {
	// Verify they are friends before allowing DM
	friendship, err := s.friendRepo.GetBetweenUsers(ctx, senderID, receiverID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFriends
		}
		return nil, err
	}
	if friendship.Status != domain.FriendshipStatusAccepted {
		return nil, domain.ErrNotFriends
	}

	msg := &domain.DirectMessage{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	if err := s.dmRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// Deliver in real time via WebSocket
	if s.notifier != nil {
		s.notifier.NotifyUser(receiverID, domain.WSTypeDMMessage, map[string]interface{}{
			"message_id": msg.ID,
			"sender_id":  senderID,
			"content":    content,
			"created_at": msg.CreatedAt,
		})
	}

	return msg, nil
}

func (s *directMessageService) GetConversation(ctx context.Context, userID, friendID uuid.UUID, page, pageSize int) ([]domain.DirectMessage, int64, error) {
	// Verify friendship exists
	friendship, err := s.friendRepo.GetBetweenUsers(ctx, userID, friendID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrNotFriends
		}
		return nil, 0, err
	}
	if friendship.Status != domain.FriendshipStatusAccepted {
		return nil, 0, domain.ErrNotFriends
	}

	offset := (page - 1) * pageSize

	// Auto-mark messages from friend as read when viewing the conversation
	_ = s.dmRepo.MarkAsRead(ctx, friendID, userID)

	return s.dmRepo.GetConversation(ctx, userID, friendID, offset, pageSize)
}
