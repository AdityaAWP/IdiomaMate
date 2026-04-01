package http

import (
	"errors"
	"net/http"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	authService    domain.AuthService
	googleClientID string
}

func NewAuthHandler(as domain.AuthService, googleClientID string) *AuthHandler {
	return &AuthHandler{
		authService:    as,
		googleClientID: googleClientID,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, domain.ErrUsernameAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req domain.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// Verify the Google ID token
	payload, err := idtoken.Validate(c.Request.Context(), req.IDToken, h.googleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid google token"})
		return
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	res, err := h.authService.GoogleLogin(c.Request.Context(), googleID, email, name, picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process google login"})
		return
	}

	c.JSON(http.StatusOK, res)
}
