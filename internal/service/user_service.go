package service

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(ur domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: ur,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetPublicProfile(ctx context.Context, userID uuid.UUID) (*domain.UserPublicProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &domain.UserPublicProfile{
		ID:               user.ID,
		Username:         user.Username,
		AvatarURL:        user.AvatarURL,
		NativeLanguage:   user.NativeLanguage,
		TargetLanguage:   user.TargetLanguage,
		ProficiencyLevel: user.ProficiencyLevel,
		ReputationScore:  user.ReputationScore,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Only update fields that were provided (non-nil pointers)
	if req.NativeLanguage != nil {
		user.NativeLanguage = *req.NativeLanguage
	}
	if req.TargetLanguage != nil {
		if !domain.IsValidTargetLanguage(*req.TargetLanguage) {
			return nil, domain.ErrInvalidLanguage
		}
		user.TargetLanguage = *req.TargetLanguage
	}
	if req.ProficiencyLevel != nil {
		user.ProficiencyLevel = *req.ProficiencyLevel
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
