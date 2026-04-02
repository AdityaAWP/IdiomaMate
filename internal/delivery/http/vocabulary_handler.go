package http

import (
	"net/http"
	"strconv"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/AdityaAWP/IdiomaMate/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VocabularyHandler struct {
	vocabService domain.VocabularyService
	topicService domain.TopicService
}

func NewVocabularyHandler(vs domain.VocabularyService, ts domain.TopicService) *VocabularyHandler {
	return &VocabularyHandler{
		vocabService: vs,
		topicService: ts,
	}
}

// SaveWord saves a new vocabulary flashcard for the authenticated user.
// POST /api/v1/vocabulary
func (h *VocabularyHandler) SaveWord(c *gin.Context) {
	userID := getUserID(c)

	var req domain.SaveVocabularyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error: utils.FormatValidationError(err),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	vocab, err := h.vocabService.SaveWord(c.Request.Context(), userID, req)
	if err != nil {
		utils.HandleError(c, err, "failed to save vocabulary word")
		return
	}

	c.JSON(http.StatusCreated, vocab)
}

// ListWords returns the authenticated user's saved vocabulary with pagination.
// GET /api/v1/vocabulary
func (h *VocabularyHandler) ListWords(c *gin.Context) {
	userID := getUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}

	result, err := h.vocabService.ListWords(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		utils.HandleError(c, err, "failed to list vocabulary")
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteWord removes a vocabulary word by ID (owner-only).
// DELETE /api/v1/vocabulary/:id
func (h *VocabularyHandler) DeleteWord(c *gin.Context) {
	userID := getUserID(c)

	wordID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "invalid word id"})
		return
	}

	if err := h.vocabService.DeleteWord(c.Request.Context(), wordID, userID); err != nil {
		utils.HandleError(c, err, "failed to delete vocabulary word")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "word deleted successfully"})
}

// GetRandomTopic returns a random discussion topic for the given language/proficiency.
// GET /api/v1/topics/random
func (h *VocabularyHandler) GetRandomTopic(c *gin.Context) {
	lang := c.Query("target_language")
	level := c.Query("proficiency_level")

	topic, err := h.topicService.GetRandomTopic(c.Request.Context(), lang, level)
	if err != nil {
		utils.HandleError(c, err, "failed to get random topic")
		return
	}

	c.JSON(http.StatusOK, topic)
}
