package http

import (
	"net/http"
	"strconv"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FriendshipHandler struct {
	friendService domain.FriendshipService
	dmService     domain.DirectMessageService
}

func NewFriendshipHandler(fs domain.FriendshipService, dms domain.DirectMessageService) *FriendshipHandler {
	return &FriendshipHandler{
		friendService: fs,
		dmService:     dms,
	}
}

// ==============================
// Friend Request Endpoints
// ==============================

// SendFriendRequest sends a friend request to another user.
// POST /api/v1/friends/request
func (h *FriendshipHandler) SendFriendRequest(c *gin.Context) {
	senderID := getUserID(c)

	var req domain.FriendRequestAction
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	friendship, err := h.friendService.SendFriendRequest(c.Request.Context(), senderID, req.UserID)
	if err != nil {
		utils.HandleError(c, err, "failed to send friend request")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "friend request sent",
		"friendship": friendship,
	})
}

// RespondToFriendRequest accepts or declines a pending friend request.
// POST /api/v1/friends/:id/respond
func (h *FriendshipHandler) RespondToFriendRequest(c *gin.Context) {
	userID := getUserID(c)

	friendshipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid friendship id"})
		return
	}

	var req domain.FriendRequestResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	accept := req.Action == "accept"
	if err := h.friendService.RespondToRequest(c.Request.Context(), friendshipID, userID, accept); err != nil {
		utils.HandleError(c, err, "failed to respond to friend request")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "responded successfully", "action": req.Action})
}

// ListFriends returns the authenticated user's accepted friends.
// The response shows only the OTHER user's profile, not the authenticated user.
// GET /api/v1/friends
func (h *FriendshipHandler) ListFriends(c *gin.Context) {
	userID := getUserID(c)

	friendships, err := h.friendService.ListFriends(c.Request.Context(), userID)
	if err != nil {
		utils.HandleError(c, err, "failed to list friends")
		return
	}

	items := make([]domain.FriendItem, 0, len(friendships))
	for _, f := range friendships {
		// Pick the OTHER user — if I'm user_id_1, the friend is user_2 and vice versa
		var friend *domain.User
		if f.UserID1 == userID {
			friend = f.User2
		} else {
			friend = f.User1
		}

		if friend == nil {
			continue
		}

		items = append(items, domain.FriendItem{
			FriendshipID: f.ID,
			Friend: domain.UserPublicProfile{
				ID:               friend.ID,
				Username:         friend.Username,
				AvatarURL:        friend.AvatarURL,
				NativeLanguage:   friend.NativeLanguage,
				TargetLanguage:   friend.TargetLanguage,
				ProficiencyLevel: friend.ProficiencyLevel,
				ReputationScore:  friend.ReputationScore,
			},
			Status: f.Status,
			Since:  f.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"friends": items})
}

// ListPendingRequests returns pending friend requests sent TO the authenticated user.
// GET /api/v1/friends/pending
func (h *FriendshipHandler) ListPendingRequests(c *gin.Context) {
	userID := getUserID(c)

	requests, err := h.friendService.ListPendingRequests(c.Request.Context(), userID)
	if err != nil {
		utils.HandleError(c, err, "failed to list pending requests")
		return
	}

	c.JSON(http.StatusOK, gin.H{"pending_requests": requests})
}

// ==============================
// Direct Message Endpoints
// ==============================

// SendDM sends a direct message to a friend.
// POST /api/v1/dm/:user_id
func (h *FriendshipHandler) SendDM(c *gin.Context) {
	senderID := getUserID(c)

	receiverID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid user id"})
		return
	}

	var req domain.SendDMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	msg, err := h.dmService.SendDM(c.Request.Context(), senderID, receiverID, req.Content)
	if err != nil {
		utils.HandleError(c, err, "failed to send direct message")
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetConversation retrieves the DM thread between the authenticated user and a friend.
// GET /api/v1/dm/:user_id
func (h *FriendshipHandler) GetConversation(c *gin.Context) {
	userID := getUserID(c)

	friendID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if page < 1 {
		page = 1
	}

	messages, total, err := h.dmService.GetConversation(c.Request.Context(), userID, friendID, page, pageSize)
	if err != nil {
		utils.HandleError(c, err, "failed to get conversation")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":    messages,
		"total_count": total,
	})
}
