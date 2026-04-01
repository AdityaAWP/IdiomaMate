package domain

import (
	"time"

	"github.com/google/uuid"
)

// MatchmakingQueue is the GORM model for the Postgres-based matchmaking queue.
// Used ONLY by the PostgresMatchmakingRepository implementation.
// The Valkey implementation uses native list structures instead.
type MatchmakingQueue struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	TargetLanguage   string    `gorm:"type:varchar(50);not null;index:idx_queue_lookup"`
	ProficiencyLevel string    `gorm:"type:varchar(20);not null;index:idx_queue_lookup"`
	CreatedAt        time.Time
}
