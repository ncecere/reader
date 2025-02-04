package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ncecere/reader-go/internal/common/config"
)

// Service handles AI-related operations
type Service struct {
	config *config.Config
	client *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the OpenAI chat completion request
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents the OpenAI chat completion response
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewService creates a new AI service
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
		client: &http.Client{},
	}
}

// Summarize generates a summary of the provided text using the configured AI model
func (s *Service) Summarize(ctx context.Context, text string) (string, error) {
	if !s.config.AI.Enabled {
		return "", fmt.Errorf("AI summarization is not enabled")
	}

	messages := []Message{
		{Role: "system", Content: s.config.AI.Prompt},
		{Role: "user", Content: text},
	}

	reqBody := ChatRequest{
		Model:    s.config.AI.Model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	endpoint := fmt.Sprintf("%s/chat/completions", s.config.AI.APIEndpoint)
	fmt.Printf("Debug - API Key: %s\n", s.config.AI.APIKey) // Debug log
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.AI.APIKey))

	// Log request details for debugging
	fmt.Printf("Making request to: %s\nHeaders: %+v\nBody: %s\n", endpoint, req.Header, jsonData)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for error cases
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create new reader from body for JSON decoding
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no summary generated")
	}

	return chatResp.Choices[0].Message.Content, nil
}
