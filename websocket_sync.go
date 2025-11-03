// websocket_sync.go - Example implementation of multi-device sync using WebSockets
package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/gorilla/websocket"
)


// WebSocket message types
type MessageType string

const (
	PlayMessage     MessageType = "play"
	PauseMessage    MessageType = "pause"
	SeekMessage     MessageType = "seek"
	VolumeMessage   MessageType = "volume"
	PlaylistMessage MessageType = "playlist"
)

// WebSocket message structure
type WebSocketMessage struct {
	Type      MessageType `json:"type"`
	UserID    string      `json:"user_id,omitempty"`
	SongID    string      `json:"song_id,omitempty"`
	Position  float64     `json:"position,omitempty"`
	Volume    float64     `json:"volume,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Client represents a WebSocket connection
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
	UserID string
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Process the message
		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Add timestamp if not present
		if msg.Timestamp == 0 {
			msg.Timestamp = time.Now().Unix()
		}

		// Broadcast to all clients
		broadcastMsg, _ := json.Marshal(msg)
		hub.Broadcast <- broadcastMsg
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// WebSocket upgrade handler
func websocketHandler(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// ServeWebSocket handles WebSocket connections
func ServeWebSocket(hub *Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		client := &Client{
			Socket: c,
			Send:   make(chan []byte, 256),
			UserID: c.Headers("User-ID"), // Get user ID from headers
		}

		hub.Register <- client

		go client.writePump()
		client.readPump(hub)
	})
}

// SendPlaySync sends a play synchronization message
func SendPlaySync(hub *Hub, userID, songID string, position float64) {
	message := WebSocketMessage{
		Type:      PlayMessage,
		UserID:    userID,
		SongID:    songID,
		Position:  position,
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(message)
	hub.Broadcast <- data
}

// SendPauseSync sends a pause synchronization message
func SendPauseSync(hub *Hub, userID string, position float64) {
	message := WebSocketMessage{
		Type:      PauseMessage,
		UserID:    userID,
		Position:  position,
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(message)
	hub.Broadcast <- data
}

// Example usage in a Fiber app
func main() {
	app := fiber.New()

	// Create hub
	hub := NewHub()
	go hub.Run()

	// Regular routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Music Player WebSocket Sync Server")
	})

	// WebSocket route
	app.Use("/ws", websocketHandler)
	app.Get("/ws/:user_id", ServeWebSocket(hub))

	// Example API endpoint that sends sync messages
	app.Post("/api/sync/play", func(c *fiber.Ctx) error {
		var req struct {
			UserID   string  `json:"user_id"`
			SongID   string  `json:"song_id"`
			Position float64 `json:"position"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		// Send play sync message
		SendPlaySync(hub, req.UserID, req.SongID, req.Position)

		return c.JSON(fiber.Map{"message": "Play sync sent"})
	})

	log.Fatal(app.Listen(":8080"))
}
