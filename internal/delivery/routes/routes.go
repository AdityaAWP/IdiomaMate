package routes

import (
	"net/http"

	httpHandler "github.com/AdityaAWP/IdiomaMate/internal/delivery/http"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/middleware"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler *httpHandler.AuthHandler
	JWTSecret   string
}

func SetupRoutes(router *gin.Engine, deps *Dependencies) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/login", deps.AuthHandler.Login)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth(deps.JWTSecret))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
			})
		})
	}
}
