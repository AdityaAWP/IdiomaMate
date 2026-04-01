package postgres

import (
	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"gorm.io/gorm"
)

// Migrate auto-migrates all domain models to PostgreSQL.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Room{},
		&domain.RoomParticipant{},
		&domain.Message{},
		&domain.Friendship{},
		&domain.DirectMessage{},
		&domain.Vocabulary{},
		&domain.Report{},
		&domain.Rating{},
		&domain.DiscussionTopic{},
	)
}
