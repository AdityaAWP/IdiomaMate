package routes

import (
	"net/http"

	httpHandler "github.com/AdityaAWP/IdiomaMate/internal/delivery/http"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/middleware"
	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler       *httpHandler.AuthHandler
	UserHandler       *httpHandler.UserHandler
	RoomHandler       *httpHandler.RoomHandler
	FriendshipHandler *httpHandler.FriendshipHandler
	VocabularyHandler *httpHandler.VocabularyHandler
	UserRepo          domain.UserRepository
	JWTSecret         string
}

func SetupRoutes(router *gin.Engine, deps *Dependencies) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// --- Public Routes (no auth required) ---
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/login", deps.AuthHandler.Login)
		auth.POST("/google", deps.AuthHandler.GoogleLogin)
	}

	// --- Protected: Profile routes (auth required, but no profile-complete check) ---
	// Users MUST be able to GET and PUT their profile even if it's incomplete.
	profile := router.Group("/api/v1")
	profile.Use(middleware.JWTAuth(deps.JWTSecret))
	{
		profile.GET("/profile", deps.UserHandler.GetProfile)
		profile.PUT("/profile", deps.UserHandler.UpdateProfile)
	}

	// --- Protected: Feature routes (auth + complete profile required) ---
	// All feature endpoints go here. Users with empty language fields get 403 PROFILE_INCOMPLETE.
	features := router.Group("/api/v1")
	features.Use(middleware.JWTAuth(deps.JWTSecret))
	features.Use(middleware.ProfileComplete(deps.UserRepo))
	{
		features.GET("/users/:id", deps.UserHandler.GetPublicProfile)
		
		// Room & Lobby Routes
		rooms := features.Group("/rooms")
		{
			rooms.POST("", deps.RoomHandler.CreateLobby)
			rooms.GET("", deps.RoomHandler.ListLobbies)
			rooms.GET("/:id", deps.RoomHandler.GetRoom)
			rooms.GET("/:id/token", deps.RoomHandler.GetAgoraToken)
			rooms.GET("/:id/messages", deps.RoomHandler.GetChatHistory)
			rooms.POST("/:id/request-join", deps.RoomHandler.RequestJoin)
			rooms.POST("/:id/respond", deps.RoomHandler.RespondJoinRequest)
			rooms.POST("/:id/leave", deps.RoomHandler.LeaveRoom)
			rooms.POST("/:id/kick", deps.RoomHandler.KickUser)
		}

		// Friendship Routes
		friends := features.Group("/friends")
		{
			friends.POST("/request", deps.FriendshipHandler.SendFriendRequest)
			friends.GET("", deps.FriendshipHandler.ListFriends)
			friends.GET("/pending", deps.FriendshipHandler.ListPendingRequests)
			friends.POST("/:id/respond", deps.FriendshipHandler.RespondToFriendRequest)
		}

		// Direct Message Routes
		dm := features.Group("/dm")
		{
			dm.POST("/:user_id", deps.FriendshipHandler.SendDM)
			dm.GET("/:user_id", deps.FriendshipHandler.GetConversation)
		}

		// Vocabulary Routes
		vocab := features.Group("/vocabulary")
		{
			vocab.POST("", deps.VocabularyHandler.SaveWord)
			vocab.GET("", deps.VocabularyHandler.ListWords)
			vocab.DELETE("/:id", deps.VocabularyHandler.DeleteWord)
		}

		// Topic Generator
		features.GET("/topics/random", deps.VocabularyHandler.GetRandomTopic)
	}
}
