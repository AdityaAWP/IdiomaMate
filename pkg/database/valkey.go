package database

import (
	"context"
	"fmt"
	"log"

	"github.com/AdityaAWP/IdiomaMate/internal/config"
	"github.com/valkey-io/valkey-go"
)

// ConnectValkey creates a connection to the Valkey (Redis-fork) server.
func ConnectValkey(cfg config.ValkeyConfig) (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Password:    cfg.Password,
		SelectDB:    cfg.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init valkey client: %w", err)
	}

	// Verify connection
	err = client.Do(context.Background(), client.B().Ping().Build()).Error()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	log.Println("Valkey Connection Success")
	return client, nil
}
