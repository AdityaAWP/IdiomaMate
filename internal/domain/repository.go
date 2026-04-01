package domain

import (
	"context"

	"github.com/google/uuid"
)

// =============================================================================
// REPOSITORY INTERFACES
// These define the contracts for data access. Each has a concrete implementation
// in internal/repository/postgres/ (and internal/repository/valkey/ for matchmaking).
// =============================================================================

// --- User Repository ---

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdateReputationScore(ctx context.Context, userID uuid.UUID, newScore float64) error
	SetShadowBanned(ctx context.Context, userID uuid.UUID, banned bool) error
}

// --- Matchmaking Repository ---
// THESIS CRITICAL: This interface MUST have two implementations:
//   1. internal/repository/postgres/matchmaking_repository.go (PostgresMatchmakingRepo)
//   2. internal/repository/valkey/matchmaking_repository.go   (ValkeyMatchmakingRepo)
// Both are benchmarked under heavy concurrency for the thesis.

type MatchmakingRepository interface {
	// Enqueue places a user into the matchmaking queue for their language+proficiency.
	Enqueue(ctx context.Context, req MatchRequest) error

	// Dequeue atomically pops the next waiting user from the queue.
	// Returns nil if no one is waiting (caller should then enqueue themselves).
	Dequeue(ctx context.Context, targetLanguage, proficiencyLevel string) (*MatchRequest, error)

	// Remove cancels a pending matchmaking request.
	Remove(ctx context.Context, userID uuid.UUID, targetLanguage, proficiencyLevel string) error
}

// --- Room Repository ---

type RoomRepository interface {
	Create(ctx context.Context, room *Room) error
	GetByID(ctx context.Context, id uuid.UUID) (*Room, error)
	Update(ctx context.Context, room *Room) error
	Close(ctx context.Context, roomID uuid.UUID) error
	ListActiveLobbies(ctx context.Context, targetLanguage, proficiencyLevel string, offset, limit int) ([]Room, int64, error)

	AddParticipant(ctx context.Context, participant *RoomParticipant) error
	RemoveParticipant(ctx context.Context, roomID, userID uuid.UUID) error
	GetParticipants(ctx context.Context, roomID uuid.UUID) ([]RoomParticipant, error)
	CountParticipants(ctx context.Context, roomID uuid.UUID) (int, error)
	IsUserInRoom(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
}

// --- Message Repository ---

type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	GetByRoomID(ctx context.Context, roomID uuid.UUID, offset, limit int) ([]Message, int64, error)
}

// --- Friendship Repository ---

type FriendshipRepository interface {
	Create(ctx context.Context, friendship *Friendship) error
	GetByID(ctx context.Context, id uuid.UUID) (*Friendship, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status FriendshipStatus) error
	GetBetweenUsers(ctx context.Context, userID1, userID2 uuid.UUID) (*Friendship, error)
	ListFriends(ctx context.Context, userID uuid.UUID) ([]Friendship, error)
	ListPendingRequests(ctx context.Context, userID uuid.UUID) ([]Friendship, error)
}

// --- Direct Message Repository ---

type DirectMessageRepository interface {
	Create(ctx context.Context, message *DirectMessage) error
	GetConversation(ctx context.Context, userID1, userID2 uuid.UUID, offset, limit int) ([]DirectMessage, int64, error)
	MarkAsRead(ctx context.Context, senderID, receiverID uuid.UUID) error
}

// --- Vocabulary Repository ---

type VocabularyRepository interface {
	Create(ctx context.Context, vocab *Vocabulary) error
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Vocabulary, int64, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

// --- Report Repository ---

type ReportRepository interface {
	Create(ctx context.Context, report *Report) error
	CountByReportedID(ctx context.Context, reportedID uuid.UUID) (int64, error)
}

// --- Rating Repository ---

type RatingRepository interface {
	Create(ctx context.Context, rating *Rating) error
	GetAverageByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
	ExistsByRaterAndRoom(ctx context.Context, raterID, ratedID, roomID uuid.UUID) (bool, error)
}

// --- Discussion Topic Repository ---

type TopicRepository interface {
	GetRandom(ctx context.Context, targetLanguage, proficiencyLevel string) (*DiscussionTopic, error)
}
