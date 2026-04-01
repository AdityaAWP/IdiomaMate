package domain

import (
	"time"

	"github.com/google/uuid"
)

// FriendshipStatus tracks friend request lifecycle.
type FriendshipStatus string

const (
	FriendshipStatusPending  FriendshipStatus = "PENDING"
	FriendshipStatusAccepted FriendshipStatus = "ACCEPTED"
	FriendshipStatusDeclined FriendshipStatus = "DECLINED"
)

// Friendship represents a directional friend request between two users.
// UserID1 is always the sender, UserID2 is the receiver.
type Friendship struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID1   uuid.UUID        `json:"user_id_1" gorm:"type:uuid;not null;index"`   // Requester
	UserID2   uuid.UUID        `json:"user_id_2" gorm:"type:uuid;not null;index"`   // Receiver
	Status    FriendshipStatus `json:"status" gorm:"type:varchar(10);not null;default:'PENDING'"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`

	// Relations
	User1 *User `json:"user_1,omitempty" gorm:"foreignKey:UserID1"`
	User2 *User `json:"user_2,omitempty" gorm:"foreignKey:UserID2"`
}

// DirectMessage represents a persistent DM between friends.
type DirectMessage struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SenderID   uuid.UUID `json:"sender_id" gorm:"type:uuid;not null;index"`
	ReceiverID uuid.UUID `json:"receiver_id" gorm:"type:uuid;not null;index"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`

	Sender   *User `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Receiver *User `json:"-" gorm:"foreignKey:ReceiverID"`
}

// --- Request / Response DTOs ---

type FriendRequestAction struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type FriendRequestResponse struct {
	Action string `json:"action" binding:"required,oneof=accept decline"` // "accept" or "decline"
}

type SendDMRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}
