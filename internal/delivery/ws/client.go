package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// Client represents a single WebSocket connection for an authenticated user.
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID uuid.UUID
	Send   chan []byte

	// Matchmaking state — stored in memory for O(1) cleanup on disconnect.
	TargetLanguage   string
	ProficiencyLevel string
	IsMatchmaking    bool
	Questions        []string
}

// ReadPump reads messages from the WebSocket connection and routes them to the Hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error for user %s: %v", c.UserID, err)
			}
			break
		}

		var wsMsg domain.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			c.SendMessage(domain.WSMessage{
				Type:    domain.WSTypeError,
				Payload: "invalid message format",
			})
			continue
		}

		c.Hub.HandleMessage(c, wsMsg)
	}
}

// WritePump sends messages from the Send channel to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage serializes and sends a WSMessage to this client.
func (c *Client) SendMessage(msg domain.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message for user %s: %v", c.UserID, err)
		return
	}

	select {
	case c.Send <- data:
	default:
		log.Printf("Send buffer full for user %s, dropping message", c.UserID)
	}
}
