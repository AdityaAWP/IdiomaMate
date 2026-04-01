package domain

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
