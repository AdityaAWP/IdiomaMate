package valkey

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

// matchmakingRepository implements domain.MatchmakingRepository using Valkey.
// Queue key pattern: "queue:{targetLanguage}:{proficiencyLevel}"
// Time complexity: O(1) for all operations (RPUSH, LPOP, LREM).
type matchmakingRepository struct {
	client valkey.Client
}

func NewMatchmakingRepository(client valkey.Client) domain.MatchmakingRepository {
	return &matchmakingRepository{client: client}
}

// queueKey builds the Valkey list key for a given language+level pair.
func queueKey(targetLanguage, proficiencyLevel string) string {
	return fmt.Sprintf("queue:%s:%s", targetLanguage, proficiencyLevel)
}

// Enqueue pushes a user's ID to the right of the queue list.
// O(1) via RPUSH.
func (r *matchmakingRepository) Enqueue(ctx context.Context, req domain.MatchRequest) error {
	key := queueKey(req.TargetLanguage, req.ProficiencyLevel)
	dataKey := "queue:data:" + req.UserID.String()
	qBytes, _ := json.Marshal(req.Questions)
	
	// Store questions in a temporary key (1 hour expiry should be safe)
	_ = r.client.Do(ctx, r.client.B().Set().Key(dataKey).Value(string(qBytes)).ExSeconds(3600).Build()).Error()

	return r.client.Do(ctx, r.client.B().Rpush().Key(key).Element(req.UserID.String()).Build()).Error()
}

// Dequeue atomically pops the leftmost (longest-waiting) user from the queue.
// O(1) via LPOP. Returns nil if the queue is empty.
func (r *matchmakingRepository) Dequeue(ctx context.Context, targetLanguage, proficiencyLevel string) (*domain.MatchRequest, error) {
	key := queueKey(targetLanguage, proficiencyLevel)

	result, err := r.client.Do(ctx, r.client.B().Lpop().Key(key).Build()).ToString()
	if valkey.IsValkeyNil(err) {
		return nil, nil // Queue is empty — no match available
	}
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(result)
	if err != nil {
		return nil, fmt.Errorf("corrupt queue entry: %w", err)
	}

	dataKey := "queue:data:" + userID.String()
	qResult, _ := r.client.Do(ctx, r.client.B().Get().Key(dataKey).Build()).AsBytes()
	var questions []string
	if len(qResult) > 0 {
		_ = json.Unmarshal(qResult, &questions)
		_ = r.client.Do(ctx, r.client.B().Del().Key(dataKey).Build())
	}

	return &domain.MatchRequest{
		UserID:           userID,
		TargetLanguage:   targetLanguage,
		ProficiencyLevel: proficiencyLevel,
		Questions:        questions,
	}, nil
}

// Remove deletes a user from the queue when they cancel or disconnect.
// O(N) via LREM, but N is typically very small for a single queue partition.
func (r *matchmakingRepository) Remove(ctx context.Context, userID uuid.UUID, targetLanguage, proficiencyLevel string) error {
	key := queueKey(targetLanguage, proficiencyLevel)
	dataKey := "queue:data:" + userID.String()
	_ = r.client.Do(ctx, r.client.B().Del().Key(dataKey).Build())
	return r.client.Do(ctx, r.client.B().Lrem().Key(key).Count(0).Element(userID.String()).Build()).Error()
}
