package server

import (
	"fmt"
	"log"
	"time"

	"github.com/AdityaAWP/IdiomaMate/internal/config"
	httpHandler "github.com/AdityaAWP/IdiomaMate/internal/delivery/http"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/routes"
	"github.com/AdityaAWP/IdiomaMate/internal/repository/postgres"
	"github.com/AdityaAWP/IdiomaMate/internal/repository/valkey"
	"github.com/AdityaAWP/IdiomaMate/internal/service"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/ws"
	"github.com/AdityaAWP/IdiomaMate/internal/delivery/middleware"
	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	router := gin.Default()

	// Global CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust this in production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		cfg:    cfg,
		router: router,
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
	roomRepo := postgres.NewRoomRepository(db)

	var matchRepo domain.MatchmakingRepository
	
	valkeyClient, err := database.ConnectValkey(s.cfg.Valkey)
	if err != nil {
		log.Printf("Warning: Valkey connection failed, using Postgres fallback for Matchmaking Repo: %v", err)
		matchRepo = postgres.NewMatchmakingRepository(db)
	} else {
		matchRepo = valkey.NewMatchmakingRepository(valkeyClient)
	}

	// 2. Initialize Services (Bottom-Up)
	jwtExpiration := time.Duration(s.cfg.JWT.Expiration) * time.Hour
	authService := service.NewAuthService(userRepo, s.cfg.JWT.Secret, jwtExpiration)
	userService := service.NewUserService(userRepo)
	matchService := service.NewMatchmakingService(matchRepo, userRepo, roomRepo)

	messageRepo := postgres.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo, roomRepo)

	friendshipRepo := postgres.NewFriendshipRepository(db)
	dmRepo := postgres.NewDirectMessageRepository(db)
	vocabRepo := postgres.NewVocabularyRepository(db)
	topicRepo := postgres.NewTopicRepository(db)

	// Add External Services
	agoraService := service.NewAgoraService(s.cfg.Agora.AppID, s.cfg.Agora.AppCertificate)
	aiService := service.NewAIGeneratorService(s.cfg.Google.GeminiAPIKey)

	// Dependency loop resolution: Hub requires matchService/messageSvc, RoomService requires Hub.
	hub := ws.NewHub(matchService, messageService, roomRepo, aiService)
	go hub.Run()

	roomService := service.NewRoomService(roomRepo, hub)
	friendshipService := service.NewFriendshipService(friendshipRepo, userRepo, hub)
	dmService := service.NewDirectMessageService(dmRepo, friendshipRepo, hub)
	vocabService := service.NewVocabularyService(vocabRepo)
	topicService := service.NewTopicService(topicRepo)

	// 3. Initialize Handlers
	authHandler := httpHandler.NewAuthHandler(authService, s.cfg.Google.ClientID)
	userHandler := httpHandler.NewUserHandler(userService)
	roomHandler := httpHandler.NewRoomHandler(roomService, messageService, agoraService)
	friendshipHandler := httpHandler.NewFriendshipHandler(friendshipService, dmService)
	vocabHandler := httpHandler.NewVocabularyHandler(vocabService, topicService)

	// Handle WebSocket Route explicitly here since it requires the hub
	s.router.GET("/api/v1/ws", middleware.JWTAuth(s.cfg.JWT.Secret), middleware.ProfileComplete(userRepo), func(c *gin.Context) {
		ws.ServeWS(hub, c)
	})

	// 4. Setup Routes
	deps := &routes.Dependencies{
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		RoomHandler:       roomHandler,
		FriendshipHandler: friendshipHandler,
		VocabularyHandler: vocabHandler,
		UserRepo:          userRepo,
		JWTSecret:         s.cfg.JWT.Secret,
	}
	routes.SetupRoutes(s.router, deps)

	// 5. Start Server
	portStr := fmt.Sprintf(":%d", s.cfg.App.Port)
	log.Printf("Starting server on port %s...", portStr)
	return s.router.Run(portStr)
}
