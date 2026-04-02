package postgres

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type directMessageRepository struct {
	db *gorm.DB
}

func NewDirectMessageRepository(db *gorm.DB) domain.DirectMessageRepository {
	return &directMessageRepository{db: db}
}

func (r *directMessageRepository) Create(ctx context.Context, message *domain.DirectMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetConversation fetches the DM thread between two users, ordered chronologically.
func (r *directMessageRepository) GetConversation(ctx context.Context, userID1, userID2 uuid.UUID, offset, limit int) ([]domain.DirectMessage, int64, error) {
	var messages []domain.DirectMessage
	var count int64

	query := r.db.WithContext(ctx).
		Model(&domain.DirectMessage{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Sender").
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, count, nil
}

// MarkAsRead bulk-marks all unread messages FROM senderID TO receiverID as read.
func (r *directMessageRepository) MarkAsRead(ctx context.Context, senderID, receiverID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.DirectMessage{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", senderID, receiverID, false).
		Update("is_read", true).Error
}
