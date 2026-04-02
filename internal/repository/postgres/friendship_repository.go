package postgres

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type friendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) domain.FriendshipRepository {
	return &friendshipRepository{db: db}
}

func (r *friendshipRepository) Create(ctx context.Context, friendship *domain.Friendship) error {
	return r.db.WithContext(ctx).Create(friendship).Error
}

func (r *friendshipRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Friendship, error) {
	var f domain.Friendship
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("id = ?", id).
		First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *friendshipRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.FriendshipStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Friendship{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetBetweenUsers checks for any existing friendship in either direction.
func (r *friendshipRepository) GetBetweenUsers(ctx context.Context, userID1, userID2 uuid.UUID) (*domain.Friendship, error) {
	var f domain.Friendship
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)",
			userID1, userID2, userID2, userID1).
		First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// ListFriends returns all ACCEPTED friendships for a user (in either direction).
func (r *friendshipRepository) ListFriends(ctx context.Context, userID uuid.UUID) ([]domain.Friendship, error) {
	var friendships []domain.Friendship
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("(user_id1 = ? OR user_id2 = ?) AND status = ?", userID, userID, domain.FriendshipStatusAccepted).
		Order("updated_at DESC").
		Find(&friendships).Error
	return friendships, err
}

// ListPendingRequests returns pending requests WHERE the user is the RECEIVER (user_id_2).
func (r *friendshipRepository) ListPendingRequests(ctx context.Context, userID uuid.UUID) ([]domain.Friendship, error) {
	var friendships []domain.Friendship
	err := r.db.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		Where("user_id2 = ? AND status = ?", userID, domain.FriendshipStatusPending).
		Order("created_at DESC").
		Find(&friendships).Error
	return friendships, err
}
