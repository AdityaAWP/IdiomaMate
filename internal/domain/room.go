package domain

import (
	"time"

	"github.com/google/uuid"
)

// RoomType defines whether a room is a 1-on-1 match or a group lobby.
type RoomType string

const (
	RoomTypeOneOnOne RoomType = "1-ON-1"
	RoomTypeLobby    RoomType = "LOBBY"
)

// RoomStatus tracks the lifecycle state of a room.
type RoomStatus string

const (
	RoomStatusWaiting RoomStatus = "WAITING" // Lobby waiting for players
	RoomStatusActive  RoomStatus = "ACTIVE"
	RoomStatusClosed  RoomStatus = "CLOSED"
)

// JoinStatus tracks a user's request to join a lobby.
type JoinStatus string

const (
	JoinStatusPending  JoinStatus = "PENDING"
	JoinStatusAccepted JoinStatus = "ACCEPTED"
	JoinStatusDeclined JoinStatus = "DECLINED"
)

// ParticipantRole differentiates the group master from regular members.
type ParticipantRole string

const (
	ParticipantRoleMember      ParticipantRole = "MEMBER"
	ParticipantRoleGroupMaster ParticipantRole = "GROUP_MASTER"
)

// Room represents a language practice session (1-on-1 or lobby).
type Room struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type             RoomType       `json:"type" gorm:"type:varchar(10);not null"`
	Status           RoomStatus     `json:"status" gorm:"type:varchar(10);not null;default:'WAITING'"`
	TargetLanguage   string         `json:"target_language" gorm:"not null"`
	ProficiencyLevel string         `json:"proficiency_level" gorm:"not null"`
	MaxParticipants  int            `json:"max_participants" gorm:"not null;default:2"`
	Title            string         `json:"title"`                                  // Optional lobby title
	AgoraChannelName string         `json:"agora_channel_name" gorm:"uniqueIndex"`  // Agora RTC channel
	Participants     []RoomParticipant `json:"participants,omitempty" gorm:"foreignKey:RoomID"`
	CreatedAt        time.Time      `json:"created_at"`
	ClosedAt         *time.Time     `json:"closed_at,omitempty"`
}

// RoomParticipant is the junction table linking users to rooms.
type RoomParticipant struct {
	ID       uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID   uuid.UUID       `json:"room_id" gorm:"type:uuid;not null;index"`
	UserID   uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	Role     ParticipantRole `json:"role" gorm:"type:varchar(15);not null;default:'MEMBER'"`
	Status   JoinStatus      `json:"status" gorm:"type:varchar(15);not null;default:'PENDING'"`
	JoinedAt time.Time       `json:"joined_at"`

	// Relations (for GORM preloading)
	Room *Room `json:"-" gorm:"foreignKey:RoomID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// --- Request / Response DTOs ---

type CreateLobbyRequest struct {
	Title            string `json:"title" binding:"required,min=3,max=100"`
	TargetLanguage   string `json:"target_language" binding:"required"`
	ProficiencyLevel string `json:"proficiency_level" binding:"required,oneof=beginner intermediate advanced"`
}

type RespondJoinRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Accept bool      `json:"accept"`
}

type KickUserRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type LobbyListResponse struct {
	Rooms      []LobbyItem `json:"rooms"`
	TotalCount int64       `json:"total_count"`
}

type LobbyItem struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	TargetLanguage   string    `json:"target_language"`
	ProficiencyLevel string    `json:"proficiency_level"`
	CurrentCount     int       `json:"current_count"`
	MaxParticipants  int       `json:"max_participants"`
	MasterUsername   string    `json:"master_username"`
	CreatedAt        time.Time `json:"created_at"`
}
