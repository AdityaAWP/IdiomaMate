package postgres

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// matchmakingRepository implements domain.MatchmakingRepository using PostgreSQL.
// Uses SELECT FOR UPDATE SKIP LOCKED to safely dequeue under concurrency.
// This is the "control" implementation for the thesis benchmark comparison.
type matchmakingRepository struct {
	db *gorm.DB
}

func NewMatchmakingRepository(db *gorm.DB) domain.MatchmakingRepository {
	return &matchmakingRepository{db: db}
}

// Enqueue inserts a user into the matchmaking_queues table.
func (r *matchmakingRepository) Enqueue(ctx context.Context, req domain.MatchRequest) error {
	entry := domain.MatchmakingQueue{
		ID:               uuid.New(),
		UserID:           req.UserID,
		TargetLanguage:   req.TargetLanguage,
		ProficiencyLevel: req.ProficiencyLevel,
	}
	return r.db.WithContext(ctx).Create(&entry).Error
}

// Dequeue atomically finds and removes the oldest waiting user from the queue.
// Uses SELECT ... FOR UPDATE SKIP LOCKED to prevent double-booking under concurrency.
// Returns nil if the queue is empty.
func (r *matchmakingRepository) Dequeue(ctx context.Context, targetLanguage, proficiencyLevel string) (*domain.MatchRequest, error) {
	var entry domain.MatchmakingQueue

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// SELECT the oldest entry, locking it and skipping already-locked rows
		result := tx.
			Where("target_language = ? AND proficiency_level = ?", targetLanguage, proficiencyLevel).
			Order("created_at ASC").
			Clauses(clause.Locking{
				Strength: "UPDATE",
				Options:  "SKIP LOCKED",
			}).
			First(&entry)

		if result.Error != nil {
			return result.Error
		}

		// DELETE the matched entry from the queue
		return tx.Delete(&entry).Error
	})

	if err == gorm.ErrRecordNotFound {
		return nil, nil // Queue is empty
	}
	if err != nil {
		return nil, err
	}

	return &domain.MatchRequest{
		UserID:           entry.UserID,
		TargetLanguage:   entry.TargetLanguage,
		ProficiencyLevel: entry.ProficiencyLevel,
	}, nil
}

// Remove deletes a user's pending entry from the queue.
func (r *matchmakingRepository) Remove(ctx context.Context, userID uuid.UUID, targetLanguage, proficiencyLevel string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND target_language = ? AND proficiency_level = ?", userID, targetLanguage, proficiencyLevel).
		Delete(&domain.MatchmakingQueue{}).Error
}
