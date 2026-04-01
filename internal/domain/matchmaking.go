package domain

import (
	"github.com/google/uuid"
)

// MatchRequest represents a user's intent to find a 1-on-1 practice partner.
type MatchRequest struct {
	UserID           uuid.UUID `json:"user_id"`
	TargetLanguage   string    `json:"target_language"`
	ProficiencyLevel string    `json:"proficiency_level"`
}

// MatchResult is returned when two users are successfully paired.
type MatchResult struct {
	RoomID           uuid.UUID `json:"room_id"`
	AgoraChannelName string    `json:"agora_channel_name"`
	AgoraToken       string    `json:"agora_token,omitempty"`
	PartnerID        uuid.UUID `json:"partner_id"`
	PartnerUsername   string    `json:"partner_username"`
}
