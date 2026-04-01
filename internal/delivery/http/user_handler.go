package http

import (
	"net/http"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(us domain.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// GetProfile returns the authenticated user's own profile.
// GET /api/v1/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := getUserID(c)

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.HandleError(c, err, "failed to get profile")
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetPublicProfile returns another user's public-safe profile.
// GET /api/v1/users/:id
func (h *UserHandler) GetPublicProfile(c *gin.Context) {
	idParam := c.Param("id")
	targetID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: "invalid user id",
			Code:  "INVALID_ID",
		})
		return
	}

	profile, err := h.userService.GetPublicProfile(c.Request.Context(), targetID)
	if err != nil {
		utils.HandleError(c, err, "failed to get user profile")
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates the authenticated user's profile.
// PUT /api/v1/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := getUserID(c)

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		utils.HandleError(c, err, "failed to update profile")
		return
	}

	c.JSON(http.StatusOK, user)
}

// getUserID extracts the authenticated user's UUID from the Gin context.
// The JWT middleware sets this value after token validation.
func getUserID(c *gin.Context) uuid.UUID {
	id, _ := c.Get("userID")
	return id.(uuid.UUID)
}
