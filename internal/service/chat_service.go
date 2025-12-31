package service

import (
	"fmt"
	"sync"

	"medseek/internal/deepseek"
	"medseek/internal/models"
)

type ChatService struct {
	deepseekClient *deepseek.Client
	sessions       map[string]*models.ChatSession
	messages       map[string][]*models.Message
	specialties    map[string]string // session_id -> specialty
	mu             sync.RWMutex
}

// NewChatService creates a new chat service
func NewChatService(deepseekAPIKey string) *ChatService {
	return &ChatService{
		deepseekClient: deepseek.NewClient(deepseekAPIKey),
		sessions:       make(map[string]*models.ChatSession),
		messages:       make(map[string][]*models.Message),
		specialties:    make(map[string]string),
	}
}

// CreateSession creates a new chat session with specialty
func (cs *ChatService) CreateSession(sessionID, userID, specialty string) *models.ChatSession {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	session := &models.ChatSession{
		ID:     sessionID,
		UserID: userID,
		Status: "active",
	}

	cs.sessions[sessionID] = session
	cs.messages[sessionID] = make([]*models.Message, 0)
	// Default to obstetrics if specialty not specified
	if specialty == "" {
		specialty = "obstetrics"
	}
	cs.specialties[sessionID] = specialty

	return session
}

// AddMessage adds a message to a session
func (cs *ChatService) AddMessage(sessionID, userID, role, content string) *models.Message {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	msg := &models.Message{
		ID:        fmt.Sprintf("%s-%d", sessionID, len(cs.messages[sessionID])),
		SessionID: sessionID,
		UserID:    userID,
		Role:      role,
		Content:   content,
	}

	if _, ok := cs.messages[sessionID]; !ok {
		cs.messages[sessionID] = make([]*models.Message, 0)
	}

	cs.messages[sessionID] = append(cs.messages[sessionID], msg)
	return msg
}

// GetSessionMessages returns all messages for a session
func (cs *ChatService) GetSessionMessages(sessionID string) []*models.Message {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if msgs, ok := cs.messages[sessionID]; ok {
		return msgs
	}
	return []*models.Message{}
}

// ProcessMessage sends a message to DeepSeek and returns the response
func (cs *ChatService) ProcessMessage(sessionID string, userMessage string) (string, error) {
	cs.mu.RLock()
	sessionMsgs := cs.messages[sessionID]
	specialty := cs.specialties[sessionID]
	cs.mu.RUnlock()

	// Build DeepSeek messages with appropriate system prompt based on specialty
	messages := []models.DeepSeekMsg{
		{
			Role:    "system",
			Content: deepseek.GetDoctorConsultationPrompt(specialty),
		},
	}

	// Add conversation history
	for _, msg := range sessionMsgs {
		messages = append(messages, models.DeepSeekMsg{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add user message
	messages = append(messages, models.DeepSeekMsg{
		Role:    "user",
		Content: userMessage,
	})

	// Get response from DeepSeek
	response, err := cs.deepseekClient.ChatCompletion(messages)
	if err != nil {
		return "", fmt.Errorf("failed to get deepseek response: %w", err)
	}

	return response, nil
}

// CloseSession closes a chat session
func (cs *ChatService) CloseSession(sessionID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if session, ok := cs.sessions[sessionID]; ok {
		session.Status = "closed"
		return nil
	}

	return fmt.Errorf("session not found: %s", sessionID)
}

// GetSession returns a session by ID
func (cs *ChatService) GetSession(sessionID string) *models.ChatSession {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cs.sessions[sessionID]
}
