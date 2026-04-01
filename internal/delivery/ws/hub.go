package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
)

// Hub maintains the set of active WebSocket clients and routes messages.
type Hub struct {
	// Registered clients mapped by UserID for O(1) lookup.
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex

	Register   chan *Client
	Unregister chan *Client

	matchService domain.MatchmakingService
}

func NewHub(matchService domain.MatchmakingService) *Hub {
	return &Hub{
		clients:      make(map[uuid.UUID]*Client),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		matchService: matchService,
	}
}

// Run starts the hub event loop. Must be launched as a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

			// Clean up matchmaking queue on disconnect
			if client.IsMatchmaking {
				_ = h.matchService.CancelMatch(context.Background(), client.UserID)
				log.Printf("Cleaned up matchmaking for disconnected user: %s", client.UserID)
			}

			log.Printf("Client disconnected: %s", client.UserID)
		}
	}
}

// GetClient returns the client for a given userID, or nil if not connected.
func (h *Hub) GetClient(userID uuid.UUID) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[userID]
}

// HandleMessage routes incoming WebSocket messages to the appropriate handler.
func (h *Hub) HandleMessage(client *Client, msg domain.WSMessage) {
	switch msg.Type {
	case domain.WSTypeMatchSearch:
		h.handleMatchSearch(client)
	case domain.WSTypeMatchCancelled:
		h.handleMatchCancel(client)
	case domain.WSTypePing:
		client.SendMessage(domain.WSMessage{Type: domain.WSTypePong})
	default:
		client.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeError,
			Payload: "unknown message type",
		})
	}
}

// handleMatchSearch attempts to find a match for the client.
func (h *Hub) handleMatchSearch(client *Client) {
	if client.IsMatchmaking {
		client.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeError,
			Payload: "already searching for a match",
		})
		return
	}

	result, err := h.matchService.FindMatch(context.Background(), client.UserID)
	if err != nil {
		client.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeMatchError,
			Payload: err.Error(),
		})
		return
	}

	if result == nil {
		// No partner yet — user has been enqueued. Mark client state.
		client.IsMatchmaking = true
		client.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeMatchSearch,
			Payload: "searching for a match...",
		})
		return
	}

	// Match found! Notify the searching user.
	client.IsMatchmaking = false
	client.SendMessage(domain.WSMessage{
		Type:    domain.WSTypeMatchFound,
		Payload: result,
	})

	// Notify the waiting partner via their WebSocket connection.
	partner := h.GetClient(result.PartnerID)
	if partner != nil {
		partner.IsMatchmaking = false

		// Build the complementary result for the partner (from their perspective).
		partnerResult := domain.MatchResult{
			RoomID:           result.RoomID,
			AgoraChannelName: result.AgoraChannelName,
			PartnerID:        client.UserID,
			PartnerUsername:   result.MyUsername, // searcher's username is the partner's "partner"
		}

		partner.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeMatchFound,
			Payload: partnerResult,
		})
	}
}

// handleMatchCancel removes the user from the matchmaking queue.
func (h *Hub) handleMatchCancel(client *Client) {
	if !client.IsMatchmaking {
		client.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeError,
			Payload: "not currently searching",
		})
		return
	}

	err := h.matchService.CancelMatch(context.Background(), client.UserID)
	if err != nil {
		client.SendMessage(domain.WSMessage{
			Type:    domain.WSTypeMatchError,
			Payload: "failed to cancel search",
		})
		return
	}

	client.IsMatchmaking = false
	client.TargetLanguage = ""
	client.ProficiencyLevel = ""

	client.SendMessage(domain.WSMessage{
		Type:    domain.WSTypeMatchCancelled,
		Payload: "search cancelled",
	})
}

// marshalPayload is a helper to decode the Payload from a WSMessage into a typed struct.
func marshalPayload(payload interface{}, target interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
