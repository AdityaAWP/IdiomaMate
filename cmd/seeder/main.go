package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AdityaAWP/IdiomaMate/internal/config"
	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/auth"
	"github.com/AdityaAWP/IdiomaMate/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func main() {
	// Load config from root
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	log.Println("Seeding 500 dummy users for k6 load testing...")

	// Open CSV to write tokens
	file, err := os.Create("tests/benchmark/users.csv")
	if err != nil {
		log.Fatalf("failed to create csv: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"token"})

	const totalUsers = 500
	const batchSize = 50

	// Clear out any previous benchmark users so our deterministic inserts work cleanly
	log.Println("Cleaning up old benchmark users...")
	db.Exec("DELETE FROM room_participants WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'benchuser%')")
	db.Exec("DELETE FROM matchmaking_queues")
	db.Exec("DELETE FROM users WHERE email LIKE 'benchuser%'")

	var users []domain.User

	for i := 1; i <= totalUsers; i++ {
		// Use deterministic UUID so re-running the seeder generates the EXACT SAME id
		userID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("benchuser%d", i)))
		user := domain.User{
			ID:               userID,
			Email:            fmt.Sprintf("benchuser%d@test.com", i),
			Username:         fmt.Sprintf("benchuser%d", i),
			TargetLanguage:   "english",
			ProficiencyLevel: "intermediate",
			NativeLanguage:   "indonesian",
			CreatedAt:        time.Now(),
		}
		users = append(users, user)

		// Generate real JWT token
		token, _ := auth.GenerateToken(user.ID, cfg.JWT.Secret, time.Duration(cfg.JWT.Expiration)*time.Hour)
		writer.Write([]string{token})

		// Insert in batches for performance
		if i%batchSize == 0 {
			// On conflict do nothing in case this script is run multiple times
			// GORM requires clauses for on-conflict ignore, but we just ignore errors nicely
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
				log.Printf("Ignored existing batch or insertion err: %v", err)
			}
			users = []domain.User{} // clear batch
			log.Printf("Processed %d users...", i)
		}
	}

	log.Println("Successfully generated `users.csv` for k6! You can safely ignore 'existing block' insert warnings if you run this multiple times.")
}
