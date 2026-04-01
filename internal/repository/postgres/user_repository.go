package postgres

import (
	"context"
	"strings"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				return domain.ErrEmailAlreadyExists
			}
			if strings.Contains(err.Error(), "username") {
				return domain.ErrUsernameAlreadyExists
			}
			return domain.ErrEmailAlreadyExists // fallback
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) UpdateReputationScore(ctx context.Context, userID uuid.UUID, newScore float64) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("reputation_score", newScore).Error
}

func (r *userRepository) SetShadowBanned(ctx context.Context, userID uuid.UUID, banned bool) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("is_shadow_banned", banned).Error
}
