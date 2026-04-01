package server

import (
	"fmt"
	"log"
	"time"

	"github.com/AdityaAWP/IdiomaMate/internal/config"
	httpHandler "github.com/AdityaAWP/IdiomaMate/internal/delivery/http"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/routes"
	"github.com/AdityaAWP/IdiomaMate/internal/repository/postgres"
	"github.com/AdityaAWP/IdiomaMate/internal/service"
	"github.com/AdityaAWP/IdiomaMate/pkg/database"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		router: gin.Default(),
	}
}

func (s *Server) Run() error {
	db, err := database.ConnectPostgres(s.cfg.Database)
	if err != nil {
		return err
	}

	log.Println("Running database migrations...")
	if err := postgres.Migrate(db); err != nil {
		return err
	}

	// 1. Initialize Repositories
	userRepo := postgres.NewUserRepository(db)

	// 2. Initialize Services
	jwtExpiration := time.Duration(s.cfg.JWT.Expiration) * time.Hour
	authService := service.NewAuthService(userRepo, s.cfg.JWT.Secret, jwtExpiration)

	// 3. Initialize Handlers
	authHandler := httpHandler.NewAuthHandler(authService, s.cfg.Google.ClientID)

	// 4. Setup Routes
	deps := &routes.Dependencies{
		AuthHandler: authHandler,
		JWTSecret:   s.cfg.JWT.Secret,
	}
	routes.SetupRoutes(s.router, deps)

	// 5. Start Server
	portStr := fmt.Sprintf(":%d", s.cfg.App.Port)
	log.Printf("Starting server on port %s...", portStr)
	return s.router.Run(portStr)
}
