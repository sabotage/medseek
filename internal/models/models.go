package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatSession represents a doctor chat session
type ChatSession struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	DoctorID  string     `json:"doctor_id,omitempty"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Status    string     `json:"status"` // active, closed, archived
}

// Message represents a message in a chat
type Message struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // user, assistant
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// DeepSeekRequest represents a request to DeepSeek API
type DeepSeekRequest struct {
	Model    string        `json:"model"`
	Messages []DeepSeekMsg `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

// DeepSeekMsg represents a message in DeepSeek format
type DeepSeekMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekResponse represents the response from DeepSeek API
type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string `json:"type"` // message, status, error
	Content   string `json:"content"`
	UserID    string `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// DoctorProfile represents a doctor's profile
type DoctorProfile struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Specialty      string   `json:"specialty"`
	LicenseNo      string   `json:"license_no"`
	Bio            string   `json:"bio"`
	Qualifications []string `json:"qualifications"`
	Available      bool     `json:"available"`
}
