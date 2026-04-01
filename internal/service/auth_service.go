package service

import (
	"context"
	"errors"
	"time"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	userRepo      domain.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewAuthService(ur domain.UserRepository, secret string, exp time.Duration) domain.AuthService {
	return &authService{
		userRepo:      ur,
		jwtSecret:     secret,
		jwtExpiration: exp,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedBytes),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.IsShadowBanned {
		return nil, domain.ErrShadowBanned
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) GoogleLogin(ctx context.Context, googleID, email, name, avatarURL string) (*domain.AuthResponse, error) {
	// Try to find existing user by Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// User doesn't exist yet — create them
		user = &domain.User{
			ID:        uuid.New(),
			Username:  name,
			Email:     email,
			GoogleID:  &googleID,
			AvatarURL: avatarURL,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Token: token, User: *user}, nil
}
