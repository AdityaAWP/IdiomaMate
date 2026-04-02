package domain

import (
	"time"

	"github.com/google/uuid"
)

// WebSocket message types for real-time communication.
type WSMessageType string

const (
	// --- Matchmaking ---
	WSTypeMatchSearch     WSMessageType = "MATCH_SEARCH"
	WSTypeMatchFound      WSMessageType = "MATCH_FOUND"
	WSTypeMatchCancelled  WSMessageType = "MATCH_CANCELLED"
	WSTypeMatchError      WSMessageType = "MATCH_ERROR"

	// --- Lobby / Room ---
	WSTypeJoinRequest     WSMessageType = "JOIN_REQUEST"
	WSTypeJoinApproved    WSMessageType = "JOIN_APPROVED"
	WSTypeJoinRejected    WSMessageType = "JOIN_REJECTED"
	WSTypeUserJoined      WSMessageType = "USER_JOINED"
	WSTypeUserLeft        WSMessageType = "USER_LEFT"
	WSTypeUserKicked      WSMessageType = "USER_KICKED"
	WSTypeRoomClosed      WSMessageType = "ROOM_CLOSED"
	WSTypeSuggestTopic    WSMessageType = "SUGGEST_TOPIC"
	WSTypeGenerateTOD     WSMessageType = "GENERATE_TOD"
	WSTypeTODResult       WSMessageType = "TOD_RESULT"

	// --- Chat ---
	WSTypeChatMessage     WSMessageType = "CHAT_MESSAGE"

	// --- Friends ---
	WSTypeFriendRequest   WSMessageType = "FRIEND_REQUEST"
	WSTypeFriendAccepted  WSMessageType = "FRIEND_ACCEPTED"
	WSTypeDMMessage       WSMessageType = "DM_MESSAGE"

	// --- System ---
	WSTypePing            WSMessageType = "PING"
	WSTypePong            WSMessageType = "PONG"
	WSTypeError           WSMessageType = "ERROR"
)

// WSMessage is the envelope for all WebSocket communication.
type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload,omitempty"`
}

// --- Specific Payloads ---

type WSMatchSearchPayload struct {
	Questions []string `json:"questions,omitempty"`
}

type WSTODPayload struct {
	RoomID uuid.UUID `json:"room_id"`
}

type WSChatMessagePayload struct {
	RoomID  uuid.UUID `json:"room_id"`
	Content string    `json:"content"`
}

type WSChatBroadcast struct {
	MessageID uuid.UUID `json:"message_id"`
	RoomID    uuid.UUID `json:"room_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
