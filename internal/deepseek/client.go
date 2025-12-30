package deepseek

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"medseek/internal/models"
)

const (
	DeepSeekAPIEndpoint = "https://api.deepseek.com/chat/completions"
	DeepSeekModel       = "deepseek-chat"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new DeepSeek API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// ChatCompletion sends a chat request to DeepSeek and returns the response
func (c *Client) ChatCompletion(messages []models.DeepSeekMsg) (string, error) {
	req := models.DeepSeekRequest{
		Model:    DeepSeekModel,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", DeepSeekAPIEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var deepseekResp models.DeepSeekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return deepseekResp.Choices[0].Message.Content, nil
}

// ChatCompletionStream sends a streaming chat request to DeepSeek
func (c *Client) ChatCompletionStream(messages []models.DeepSeekMsg, callback func(string) error) error {
	req := models.DeepSeekRequest{
		Model:    DeepSeekModel,
		Messages: messages,
		Stream:   true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", DeepSeekAPIEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
				if err := callback(streamResp.Choices[0].Delta.Content); err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}

// GetDoctorConsultationPrompt returns a system prompt for doctor consultation
func GetDoctorConsultationPrompt() string {
	return `You are a helpful medical assistant AI designed to provide online doctor consultation services.
Your role is to:
1. Listen to patient symptoms and health concerns
2. Provide general medical information and guidance
3. Ask clarifying questions about symptoms
4. Suggest when to seek in-person medical care
5. Maintain patient confidentiality and privacy
6. Provide evidence-based medical information

Important: 
- Always remind patients that you're an AI assistant and not a substitute for professional medical advice
- For serious or emergency symptoms, always recommend immediate professional medical care
- Be empathetic and professional in your responses
- Ask about relevant medical history when appropriate

Start by greeting the patient and asking about their health concern.`
}
