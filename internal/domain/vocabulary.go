package domain

import (
	"time"

	"github.com/google/uuid"
)

// Vocabulary represents a flashcard word saved by a user from chat.
type Vocabulary struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	TargetWord      string    `json:"target_word" gorm:"type:varchar(255);not null"`
	Translation     string    `json:"translation" gorm:"type:varchar(255)"`
	ContextSentence string    `json:"context_sentence" gorm:"type:text"`
	SourceRoomID    *uuid.UUID `json:"source_room_id,omitempty" gorm:"type:uuid"`
	CreatedAt       time.Time `json:"created_at"`

	// Relations
	User *User `json:"-" gorm:"foreignKey:UserID"`
	Room *Room `json:"-" gorm:"foreignKey:SourceRoomID"`
}

// --- Request / Response DTOs ---

type SaveVocabularyRequest struct {
	TargetWord      string     `json:"target_word" binding:"required,min=1,max=255"`
	Translation     string     `json:"translation" binding:"max=255"`
	ContextSentence string     `json:"context_sentence" binding:"max=1000"`
	SourceRoomID    *uuid.UUID `json:"source_room_id"`
}

type VocabularyListResponse struct {
	Words      []Vocabulary `json:"words"`
	TotalCount int64        `json:"total_count"`
}
