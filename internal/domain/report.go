package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReportReason categorizes why a user is being reported.
type ReportReason string

const (
	ReportReasonAbuse        ReportReason = "ABUSE"
	ReportReasonSpam         ReportReason = "SPAM"
	ReportReasonInappropriate ReportReason = "INAPPROPRIATE"
	ReportReasonOther        ReportReason = "OTHER"
)

// Report represents a user report for safety/moderation.
type Report struct {
	ID         uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReporterID uuid.UUID    `json:"reporter_id" gorm:"type:uuid;not null;index"`
	ReportedID uuid.UUID    `json:"reported_id" gorm:"type:uuid;not null;index"`
	RoomID     uuid.UUID    `json:"room_id" gorm:"type:uuid;not null;index"`
	Reason     ReportReason `json:"reason" gorm:"type:varchar(20);not null"`
	ReasonText string       `json:"reason_text" gorm:"type:text"`
	IsReviewed bool         `json:"is_reviewed" gorm:"default:false"`
	CreatedAt  time.Time    `json:"created_at"`

	// Relations
	Reporter *User `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	Reported *User `json:"reported,omitempty" gorm:"foreignKey:ReportedID"`
	Room     *Room `json:"-" gorm:"foreignKey:RoomID"`
}

// Rating represents a post-session peer rating.
type Rating struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RaterID   uuid.UUID `json:"rater_id" gorm:"type:uuid;not null;index"`
	RatedID   uuid.UUID `json:"rated_id" gorm:"type:uuid;not null;index"`
	RoomID    uuid.UUID `json:"room_id" gorm:"type:uuid;not null;index"`
	Score     int       `json:"score" gorm:"not null"` // 1-5 star rating
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Rater *User `json:"-" gorm:"foreignKey:RaterID"`
	Rated *User `json:"-" gorm:"foreignKey:RatedID"`
	Room  *Room `json:"-" gorm:"foreignKey:RoomID"`
}

// --- Request / Response DTOs ---

type CreateReportRequest struct {
	ReportedID uuid.UUID    `json:"reported_id" binding:"required"`
	RoomID     uuid.UUID    `json:"room_id" binding:"required"`
	Reason     ReportReason `json:"reason" binding:"required,oneof=ABUSE SPAM INAPPROPRIATE OTHER"`
	ReasonText string       `json:"reason_text" binding:"max=1000"`
}

type CreateRatingRequest struct {
	RatedID uuid.UUID `json:"rated_id" binding:"required"`
	RoomID  uuid.UUID `json:"room_id" binding:"required"`
	Score   int       `json:"score" binding:"required,min=1,max=5"`
	Comment string    `json:"comment" binding:"max=500"`
}
