package middleware

import (
	"net/http"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProfileComplete ensures the authenticated user has filled in their
// native_language, target_language, and proficiency_level before
// accessing any feature endpoints. Profile routes themselves are exempt.
func ProfileComplete(userRepo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawID, _ := c.Get("userID")
		userID := rawID.(uuid.UUID)

		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
				Error: "user not found",
				Code:  "USER_NOT_FOUND",
			})
			return
		}

		if user.NativeLanguage == "" || user.TargetLanguage == "" || user.ProficiencyLevel == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse{
				Error: domain.ErrProfileIncomplete.Error(),
				Code:  "PROFILE_INCOMPLETE",
			})
			return
		}

		c.Next()
	}
}
