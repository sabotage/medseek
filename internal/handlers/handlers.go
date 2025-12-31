package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"medseek/internal/service"
	wshub "medseek/internal/websocket"
)

type Handler struct {
	chatSvc *service.ChatService
	hub     *wshub.Hub
}

// NewHandler creates a new handler
func NewHandler(chatSvc *service.ChatService, hub *wshub.Hub) *Handler {
	return &Handler{
		chatSvc: chatSvc,
		hub:     hub,
	}
}

// CreateSessionRequest represents the request to create a new session
type CreateSessionRequest struct {
	UserID    string `json:"user_id"`
	Specialty string `json:"specialty"` // "obstetrics" or "pediatrics"
}

// CreateSessionResponse represents the response when creating a session
type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
}

// CreateSession creates a new chat session
func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()
	if req.Specialty == "" {
		req.Specialty = "obstetrics" // default specialty
	}
	session := h.chatSvc.CreateSession(sessionID, req.UserID, req.Specialty)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateSessionResponse{
		SessionID: session.ID,
		Status:    session.Status,
	})
}

// WebSocketUpgrade upgrades the connection to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // For development, allow all origins
	},
	// iOS Safari WebSocket support
	HandshakeTimeout: 45 * time.Second,
}

// WebSocket handles WebSocket connections
func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	userID := r.URL.Query().Get("user_id")

	if sessionID == "" || userID == "" {
		http.Error(w, "Missing session_id or user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	h.hub.HandleConnection(conn, userID, sessionID)
}

// GetSessionMessages returns messages for a session
func (h *Handler) GetSessionMessages(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	if sessionID == "" {
		http.Error(w, "Missing session_id", http.StatusBadRequest)
		return
	}

	messages := h.chatSvc.GetSessionMessages(sessionID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// CloseSession closes a chat session
func (h *Handler) CloseSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")

	if sessionID == "" {
		http.Error(w, "Missing session_id", http.StatusBadRequest)
		return
	}

	err := h.chatSvc.CloseSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "closed",
	})
}

// Health check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
