package http

import (
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
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		utils.HandleError(c, err, "failed to register user")
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		utils.HandleError(c, err, "failed to login")
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req domain.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	payload, err := idtoken.Validate(c.Request.Context(), req.IDToken, h.googleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Error: "invalid google token",
			Code:  "INVALID_GOOGLE_TOKEN",
		})
		return
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	res, err := h.authService.GoogleLogin(c.Request.Context(), googleID, email, name, picture)
	if err != nil {
		utils.HandleError(c, err, "failed to process google login")
		return
	}

	c.JSON(http.StatusOK, res)
}
