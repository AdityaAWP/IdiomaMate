package domain

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a single chat message within a room.
// Messages are first buffered through Valkey Streams, then persisted here.
type Message struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID    uuid.UUID `json:"room_id" gorm:"type:uuid;not null;index"`
	SenderID  uuid.UUID `json:"sender_id" gorm:"type:uuid;not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Room   *Room `json:"-" gorm:"foreignKey:RoomID"`
	Sender *User `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}

// --- Request / Response DTOs ---

type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}

type ChatHistoryResponse struct {
	Messages   []Message `json:"messages"`
	TotalCount int64     `json:"total_count"`
}
