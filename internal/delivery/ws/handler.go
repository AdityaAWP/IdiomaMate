package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all initially
	},
}

// ServeWS handles WebSocket requests from the peer and registers them to the hub.
func ServeWS(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	rawID, exists := c.Get("userID")
	if !exists {
		conn.Close()
		return
	}
	userID := rawID.(uuid.UUID)

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
