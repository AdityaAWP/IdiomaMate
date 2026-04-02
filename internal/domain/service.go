package domain

import (
	"context"

	"github.com/google/uuid"
)

// =============================================================================
// SERVICE INTERFACES
// These define the business logic contracts consumed by the delivery layer.
// Concrete implementations live in internal/service/.
// =============================================================================

type NotificationService interface {
	NotifyUser(userID uuid.UUID, wsType WSMessageType, payload interface{})
}

// --- Auth Service ---

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	// GoogleLogin handles OAuth2 callback; creates user if first login.
	GoogleLogin(ctx context.Context, googleID, email, name, avatarURL string) (*AuthResponse, error)
}

// --- User Service ---

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
	GetPublicProfile(ctx context.Context, userID uuid.UUID) (*UserPublicProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*User, error)
}

// --- Matchmaking Service ---

type MatchmakingService interface {
	// FindMatch attempts to pair the user with a waiting partner.
	// If a partner is found, returns MatchResult immediately.
	// If no partner is found, enqueues the user and returns nil (WebSocket will notify later).
	FindMatch(ctx context.Context, userID uuid.UUID, questions []string) (*MatchResult, error)
	CancelMatch(ctx context.Context, userID uuid.UUID) error
}

// --- Room / Lobby Service ---

type RoomService interface {
	CreateLobby(ctx context.Context, masterID uuid.UUID, req CreateLobbyRequest) (*Room, error)
	GetRoom(ctx context.Context, roomID uuid.UUID) (*Room, error)
	ListLobbies(ctx context.Context, targetLanguage, proficiencyLevel string, page, pageSize int) (*LobbyListResponse, error)
	RequestJoin(ctx context.Context, roomID, userID uuid.UUID) error
	RespondJoinRequest(ctx context.Context, roomID, masterID, targetUserID uuid.UUID, accept bool) error
	LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error
	CloseRoom(ctx context.Context, roomID, masterID uuid.UUID) error
	KickUser(ctx context.Context, roomID, masterID, targetUserID uuid.UUID) error
}

// --- Message Service ---

type MessageService interface {
	SaveMessage(ctx context.Context, roomID, senderID uuid.UUID, content string) (*Message, error)
	GetChatHistory(ctx context.Context, roomID uuid.UUID, page, pageSize int) (*ChatHistoryResponse, error)
}

// --- Friendship Service ---

type FriendshipService interface {
	SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) (*Friendship, error)
	RespondToRequest(ctx context.Context, friendshipID, userID uuid.UUID, accept bool) error
	ListFriends(ctx context.Context, userID uuid.UUID) ([]Friendship, error)
	ListPendingRequests(ctx context.Context, userID uuid.UUID) ([]Friendship, error)
}

// --- Direct Message Service ---

type DirectMessageService interface {
	SendDM(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*DirectMessage, error)
	GetConversation(ctx context.Context, userID, friendID uuid.UUID, page, pageSize int) ([]DirectMessage, int64, error)
}

// --- Vocabulary Service ---

type VocabularyService interface {
	SaveWord(ctx context.Context, userID uuid.UUID, req SaveVocabularyRequest) (*Vocabulary, error)
	ListWords(ctx context.Context, userID uuid.UUID, page, pageSize int) (*VocabularyListResponse, error)
	DeleteWord(ctx context.Context, wordID, userID uuid.UUID) error
}

// --- Report & Rating Service ---

type SafetyService interface {
	ReportUser(ctx context.Context, reporterID uuid.UUID, req CreateReportRequest) error
	RateUser(ctx context.Context, raterID uuid.UUID, req CreateRatingRequest) error
}

// --- Topic Service ---

type TopicService interface {
	GetRandomTopic(ctx context.Context, targetLanguage, proficiencyLevel string) (*DiscussionTopic, error)
}

// --- Agora Token Service ---

type AgoraService interface {
	GenerateRTCToken(channelName string, userID uuid.UUID) (string, error)
}

// --- AI Generator Service ---

type AIGeneratorService interface {
	GenerateTOD(ctx context.Context, targetLanguage, proficiencyLevel string) (string, error)
}
