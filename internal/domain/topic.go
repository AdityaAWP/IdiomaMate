package domain

import (
	"github.com/google/uuid"
)

// DiscussionTopic stores conversation starters for the "Suggest Topic" feature.
type DiscussionTopic struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TargetLanguage   string    `json:"target_language" gorm:"type:varchar(50);not null;index"`
	ProficiencyLevel string    `json:"proficiency_level" gorm:"type:varchar(20);not null;index"`
	TopicText        string    `json:"topic_text" gorm:"type:text;not null"`
}
