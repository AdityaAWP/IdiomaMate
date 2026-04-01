package utils

import (
	"errors"
	"net/http"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/gin-gonic/gin"
)

// ErrorResponse is the standard JSON error envelope for all API responses.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// Known domain error to HTTP status + code mappings.
var domainErrorMap = map[error]struct {
	Status int
	Code   string
}{
	// Auth & User
	domain.ErrUserNotFound:          {http.StatusNotFound, "USER_NOT_FOUND"},
	domain.ErrEmailAlreadyExists:    {http.StatusConflict, "EMAIL_EXISTS"},
	domain.ErrUsernameAlreadyExists: {http.StatusConflict, "USERNAME_EXISTS"},
	domain.ErrInvalidCredentials:    {http.StatusUnauthorized, "INVALID_CREDENTIALS"},
	domain.ErrUnauthorized:          {http.StatusUnauthorized, "UNAUTHORIZED"},
	domain.ErrShadowBanned:          {http.StatusForbidden, "ACCOUNT_RESTRICTED"},

	// Room & Matchmaking
	domain.ErrRoomNotFound:     {http.StatusNotFound, "ROOM_NOT_FOUND"},
	domain.ErrRoomFull:         {http.StatusConflict, "ROOM_FULL"},
	domain.ErrRoomClosed:       {http.StatusGone, "ROOM_CLOSED"},
	domain.ErrAlreadyInRoom:    {http.StatusConflict, "ALREADY_IN_ROOM"},
	domain.ErrNotRoomMaster:    {http.StatusForbidden, "NOT_ROOM_MASTER"},
	domain.ErrNoMatchAvailable: {http.StatusNotFound, "NO_MATCH"},

	// Friendship
	domain.ErrFriendshipNotFound:  {http.StatusNotFound, "FRIENDSHIP_NOT_FOUND"},
	domain.ErrAlreadyFriends:      {http.StatusConflict, "ALREADY_FRIENDS"},
	domain.ErrFriendRequestExists: {http.StatusConflict, "FRIEND_REQUEST_EXISTS"},
	domain.ErrCannotFriendSelf:    {http.StatusBadRequest, "CANNOT_FRIEND_SELF"},
	domain.ErrNotFriends:          {http.StatusBadRequest, "NOT_FRIENDS"},

	// General
	domain.ErrNotFound:        {http.StatusNotFound, "NOT_FOUND"},
	domain.ErrBadRequest:      {http.StatusBadRequest, "BAD_REQUEST"},
	domain.ErrForbidden:       {http.StatusForbidden, "FORBIDDEN"},
	domain.ErrDuplicateRating: {http.StatusConflict, "DUPLICATE_RATING"},
}

// HandleError writes a consistent JSON error response based on domain errors.
// For unknown errors, it returns a generic 500 with a safe message.
func HandleError(c *gin.Context, err error, fallbackMessage string) {
	for domainErr, info := range domainErrorMap {
		if errors.Is(err, domainErr) {
			c.JSON(info.Status, ErrorResponse{
				Error: domainErr.Error(),
				Code:  info.Code,
			})
			return
		}
	}

	// Unknown/internal error — never leak implementation details
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: fallbackMessage,
		Code:  "INTERNAL_ERROR",
	})
}
