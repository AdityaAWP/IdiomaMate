package postgres

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type vocabularyRepository struct {
	db *gorm.DB
}

func NewVocabularyRepository(db *gorm.DB) domain.VocabularyRepository {
	return &vocabularyRepository{db: db}
}

func (r *vocabularyRepository) Create(ctx context.Context, vocab *domain.Vocabulary) error {
	return r.db.WithContext(ctx).Create(vocab).Error
}

// GetByUserID returns a paginated list of the user's saved vocabulary words.
func (r *vocabularyRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]domain.Vocabulary, int64, error) {
	var words []domain.Vocabulary
	var count int64

	query := r.db.WithContext(ctx).
		Model(&domain.Vocabulary{}).
		Where("user_id = ?", userID)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&words).Error; err != nil {
		return nil, 0, err
	}

	return words, count, nil
}

// Delete removes a vocabulary word. Only the owner can delete their own words.
func (r *vocabularyRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Vocabulary{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
