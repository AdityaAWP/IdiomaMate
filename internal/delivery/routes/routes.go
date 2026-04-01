package routes

import (
	"net/http"

	httpHandler "github.com/AdityaAWP/IdiomaMate/internal/delivery/http"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/middleware"
	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler *httpHandler.AuthHandler
	UserHandler *httpHandler.UserHandler
	UserRepo    domain.UserRepository
	JWTSecret   string
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
		// Future feature routes (matchmaking, rooms, friends, etc.) go here
	}
}
