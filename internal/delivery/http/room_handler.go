package http

import (
	"net/http"
	"strconv"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomHandler struct {
	roomService    domain.RoomService
	messageService domain.MessageService
}

func NewRoomHandler(rs domain.RoomService, ms domain.MessageService) *RoomHandler {
	return &RoomHandler{
		roomService:    rs,
		messageService: ms,
	}
}

func (h *RoomHandler) CreateLobby(c *gin.Context) {
	userID := getUserID(c)

	var req domain.CreateLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	room, err := h.roomService.CreateLobby(c.Request.Context(), userID, req)
	if err != nil {
		utils.HandleError(c, err, "failed to create lobby")
		return
	}

	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) ListLobbies(c *gin.Context) {
	lang := c.Query("target_language")
	level := c.Query("proficiency_level")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}

	resp, err := h.roomService.ListLobbies(c.Request.Context(), lang, level, page, pageSize)
	if err != nil {
		utils.HandleError(c, err, "failed to list lobbies")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid room id"})
		return
	}

	room, err := h.roomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		utils.HandleError(c, err, "failed to get room")
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) RequestJoin(c *gin.Context) {
	userID := getUserID(c)
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid room id"})
		return
	}

	if err := h.roomService.RequestJoin(c.Request.Context(), roomID, userID); err != nil {
		utils.HandleError(c, err, "failed to request join room")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "join request sent to master", "room_id": roomID, "status": "pending"})
}

func (h *RoomHandler) RespondJoinRequest(c *gin.Context) {
	masterID := getUserID(c)
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid room id"})
		return
	}

	var req domain.RespondJoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	if err := h.roomService.RespondJoinRequest(c.Request.Context(), roomID, masterID, req.UserID, req.Accept); err != nil {
		utils.HandleError(c, err, "failed to respond to join request")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "responded successfully"})
}

func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	userID := getUserID(c)
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid room id"})
		return
	}

	if err := h.roomService.LeaveRoom(c.Request.Context(), roomID, userID); err != nil {
		utils.HandleError(c, err, "failed to leave room")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left room successfully"})
}

func (h *RoomHandler) KickUser(c *gin.Context) {
	masterID := getUserID(c)
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid room id"})
		return
	}

	var req domain.KickUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	if err := h.roomService.KickUser(c.Request.Context(), roomID, masterID, req.UserID); err != nil {
		utils.HandleError(c, err, "failed to kick user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user kicked successfully"})
}

func (h *RoomHandler) GetChatHistory(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid room id"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if page < 1 {
		page = 1
	}

	history, err := h.messageService.GetChatHistory(c.Request.Context(), roomID, page, pageSize)
	if err != nil {
		utils.HandleError(c, err, "failed to get chat history")
		return
	}

	c.JSON(http.StatusOK, history)
}
