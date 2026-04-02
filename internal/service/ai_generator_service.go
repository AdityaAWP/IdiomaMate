package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type aiGeneratorService struct {
	apiKey string
}

func NewAIGeneratorService(apiKey string) domain.AIGeneratorService {
	return &aiGeneratorService{apiKey: apiKey}
}

func (s *aiGeneratorService) GenerateTOD(ctx context.Context, targetLanguage, proficiencyLevel string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("gemini api key not configured")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3-flash-preview")

	// Creativity is fine for TODs
	temp := float32(0.9)
	model.Temperature = &temp

	prompt := fmt.Sprintf(
		"Generate a fun, thought-provoking 'Truth or Dare' (TOD) question for two people practicing a language together. "+
			"The target language they are practicing is %s. "+
			"Their proficiency level is %s. "+
			"Respond with ONLY the specific Truth or Dare text in the target language without any generic intro words. Keep it engaging, respectful, and appropriately matched to their proficiency level.",
		targetLanguage, proficiencyLevel,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if text, ok := part.(genai.Text); ok {
			return strings.TrimSpace(string(text)), nil
		}
	}

	return "", fmt.Errorf("failed to parse gemini response")
}
