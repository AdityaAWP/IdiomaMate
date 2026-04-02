package service

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
)

type vocabularyService struct {
	vocabRepo domain.VocabularyRepository
}

func NewVocabularyService(vr domain.VocabularyRepository) domain.VocabularyService {
	return &vocabularyService{vocabRepo: vr}
}

func (s *vocabularyService) SaveWord(ctx context.Context, userID uuid.UUID, req domain.SaveVocabularyRequest) (*domain.Vocabulary, error) {
	vocab := &domain.Vocabulary{
		ID:              uuid.New(),
		UserID:          userID,
		TargetWord:      req.TargetWord,
		Translation:     req.Translation,
		ContextSentence: req.ContextSentence,
		SourceRoomID:    req.SourceRoomID,
	}

	if err := s.vocabRepo.Create(ctx, vocab); err != nil {
		return nil, err
	}

	return vocab, nil
}

func (s *vocabularyService) ListWords(ctx context.Context, userID uuid.UUID, page, pageSize int) (*domain.VocabularyListResponse, error) {
	offset := (page - 1) * pageSize

	words, count, err := s.vocabRepo.GetByUserID(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	return &domain.VocabularyListResponse{
		Words:      words,
		TotalCount: count,
	}, nil
}

func (s *vocabularyService) DeleteWord(ctx context.Context, wordID, userID uuid.UUID) error {
	return s.vocabRepo.Delete(ctx, wordID, userID)
}
