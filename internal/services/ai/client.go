package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

type aiClient struct {
	apiURL     string
	apiKey     string
	model      string
	httpClient *http.Client
}

func newAIClient(cfg config.AIConfig) *aiClient {
	return &aiClient{
		apiURL: cfg.APIURL,
		apiKey: cfg.APIKey,
		model:  cfg.Model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *aiClient) IsConfigured() bool {
	return c.apiURL != "" && c.apiKey != ""
}

func (c *aiClient) Chat(systemPrompt, userPrompt string) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.8,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error al serializar request: %w", err)
	}

	url := buildChatURL(c.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al llamar a la IA: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ AI API error (%d) llamando %s: %s", resp.StatusCode, url, string(body))
		return "", fmt.Errorf("AI API error (%d)", resp.StatusCode)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("error al parsear respuesta: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("la IA no devolvió contenido")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func buildChatURL(apiURL string) string {
	base := strings.TrimRight(apiURL, "/")
	if strings.HasSuffix(base, "/chat/completions") {
		return base
	}
	return base + "/chat/completions"
}
