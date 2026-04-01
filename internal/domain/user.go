package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a registered user in the system.
type User struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username         string         `json:"username" gorm:"uniqueIndex;not null"`
	Email            string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string         `json:"-"`
	GoogleID         *string        `json:"-" gorm:"uniqueIndex"`
	AvatarURL        string         `json:"avatar_url"`
	NativeLanguage   string         `json:"native_language"`
	TargetLanguage   string         `json:"target_language"`
	ProficiencyLevel string         `json:"proficiency_level"` // beginner, intermediate, advanced
	ReputationScore  float64        `json:"reputation_score" gorm:"default:5.0"`
	IsShadowBanned   bool           `json:"-" gorm:"default:false"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// --- Request / Response DTOs ---

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type UpdateProfileRequest struct {
	NativeLanguage   *string `json:"native_language"`
	TargetLanguage   *string `json:"target_language"`
	ProficiencyLevel *string `json:"proficiency_level" binding:"omitempty,oneof=beginner intermediate advanced"`
	AvatarURL        *string `json:"avatar_url"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserPublicProfile struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	AvatarURL        string    `json:"avatar_url"`
	NativeLanguage   string    `json:"native_language"`
	TargetLanguage   string    `json:"target_language"`
	ProficiencyLevel string    `json:"proficiency_level"`
	ReputationScore  float64   `json:"reputation_score"`
}
