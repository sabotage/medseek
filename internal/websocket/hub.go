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
	sessions   map[string]map[*Client]bool // session_id -> clients in that session
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
		sessions:   make(map[string]map[*Client]bool),
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
			// Add client to session-specific map
			if h.sessions[client.SessionID] == nil {
				h.sessions[client.SessionID] = make(map[*Client]bool)
			}
			h.sessions[client.SessionID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s in session: %s", client.ID, client.SessionID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				// Remove from session-specific map
				if sessionClients, ok := h.sessions[client.SessionID]; ok {
					delete(sessionClients, client)
					if len(sessionClients) == 0 {
						delete(h.sessions, client.SessionID)
					}
				}
				close(client.send)
				h.mu.Unlock()
				log.Printf("Client unregistered: %s from session: %s", client.ID, client.SessionID)
			} else {
				h.mu.Unlock()
			}
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

		// Send user message to clients in this session only
		msgBytes, _ := json.Marshal(wsMsg)
		c.hub.broadcastToSession(c.SessionID, msgBytes)

		// Process message and get response
		response, err := c.hub.chatSvc.ProcessMessage(c.SessionID, wsMsg.Content)
		if err != nil {
			log.Printf("Failed to process message: %v", err)
			errMsg := models.WebSocketMessage{
				Type:    "error",
				Content: fmt.Sprintf("处理消息时出错: %v", err),
			}
			errBytes, _ := json.Marshal(errMsg)
			c.hub.broadcastToSession(c.SessionID, errBytes)
			continue
		}

		// Add assistant message to service
		c.hub.chatSvc.AddMessage(c.SessionID, "assistant", "assistant", response)

		// Send assistant response to this session only
		respMsg := models.WebSocketMessage{
			Type:    "message",
			Content: response,
			UserID:  "assistant",
		}
		respBytes, _ := json.Marshal(respMsg)
		c.hub.broadcastToSession(c.SessionID, respBytes)
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

// broadcastToSession sends a message to all clients in a specific session
func (h *Hub) broadcastToSession(sessionID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if sessionClients, ok := h.sessions[sessionID]; ok {
		for client := range sessionClients {
			select {
			case client.send <- message:
			default:
				// Client's send channel is full, skip
				log.Printf("Warning: Could not send message to client %s in session %s", client.ID, sessionID)
			}
		}
	}
}
