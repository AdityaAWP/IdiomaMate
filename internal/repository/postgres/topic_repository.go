package postgres

import (
	"context"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"gorm.io/gorm"
)

type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) domain.TopicRepository {
	return &topicRepository{db: db}
}

// GetRandom returns a random discussion topic matching the given language and proficiency.
// Uses PostgreSQL's ORDER BY RANDOM() for simplicity.
func (r *topicRepository) GetRandom(ctx context.Context, targetLanguage, proficiencyLevel string) (*domain.DiscussionTopic, error) {
	var topic domain.DiscussionTopic

	query := r.db.WithContext(ctx).Model(&domain.DiscussionTopic{})

	if targetLanguage != "" {
		query = query.Where("target_language = ?", targetLanguage)
	}
	if proficiencyLevel != "" {
		query = query.Where("proficiency_level = ?", proficiencyLevel)
	}

	err := query.Order("RANDOM()").First(&topic).Error
	if err != nil {
		return nil, err
	}

	return &topic, nil
}
