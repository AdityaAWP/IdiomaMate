package postgres

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID, offset, limit int) ([]domain.Message, int64, error) {
	var messages []domain.Message
	var count int64

	query := r.db.WithContext(ctx).
		Model(&domain.Message{}).
		Where("room_id = ?", roomID)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Sender"). // Preload user info like username and avatar
		Order("created_at ASC"). // Return chronological order
		Offset(offset).Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, count, nil
}
