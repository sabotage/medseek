package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"medseek/internal/models"
	"medseek/internal/service"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	chatSvc    *service.ChatService
}

type Client struct {
	ID        string
	SessionID string
	conn      *websocket.Conn
	send      chan []byte
	hub       *Hub
}

// NewHub creates a new WebSocket hub
func NewHub(chatSvc *service.ChatService) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		chatSvc:    chatSvc,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.mu.Unlock()
				log.Printf("Client unregistered: %s", client.ID)
			} else {
				h.mu.Unlock()
			}

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// HandleConnection handles a new WebSocket connection
func (h *Hub) HandleConnection(conn *websocket.Conn, clientID, sessionID string) {
	client := &Client{
		ID:        clientID,
		SessionID: sessionID,
		conn:      conn,
		send:      make(chan []byte, 256),
		hub:       h,
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

// readPump reads messages from WebSocket
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Time{}) // No deadline for reading

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		var wsMsg models.WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		wsMsg.UserID = c.ID
		wsMsg.SessionID = c.SessionID

		// Add user message to service
		c.hub.chatSvc.AddMessage(c.SessionID, c.ID, "user", wsMsg.Content)

		// Send message to all clients in session
		c.hub.broadcast <- message

		// Process message and get response
		response, err := c.hub.chatSvc.ProcessMessage(c.SessionID, wsMsg.Content)
		if err != nil {
			log.Printf("Failed to process message: %v", err)
			errMsg := models.WebSocketMessage{
				Type:    "error",
				Content: fmt.Sprintf("Error processing message: %v", err),
			}
			errBytes, _ := json.Marshal(errMsg)
			c.hub.broadcast <- errBytes
			continue
		}

		// Add assistant message to service
		c.hub.chatSvc.AddMessage(c.SessionID, "assistant", "assistant", response)

		// Send assistant response
		respMsg := models.WebSocketMessage{
			Type:    "message",
			Content: response,
			UserID:  "assistant",
		}
		respBytes, _ := json.Marshal(respMsg)
		c.hub.broadcast <- respBytes
	}
}

// writePump writes messages to WebSocket
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Time{})
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}
