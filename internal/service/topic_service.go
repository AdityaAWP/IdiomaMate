package service

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
)

type topicService struct {
	topicRepo domain.TopicRepository
}

func NewTopicService(tr domain.TopicRepository) domain.TopicService {
	return &topicService{topicRepo: tr}
}

func (s *topicService) GetRandomTopic(ctx context.Context, targetLanguage, proficiencyLevel string) (*domain.DiscussionTopic, error) {
	return s.topicRepo.GetRandom(ctx, targetLanguage, proficiencyLevel)
}
