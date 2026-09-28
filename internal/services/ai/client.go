package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/google/uuid"
)

const defaultTemperature = 0.8

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
	timeout    time.Duration
	sessionID  string
	httpClient *http.Client
}

func newAIClient(cfg config.AIConfig) *aiClient {
	timeout := time.Duration(cfg.RequestTimeout) * time.Second
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	return &aiClient{
		apiURL:    cfg.APIURL,
		apiKey:    cfg.APIKey,
		model:     cfg.Model,
		timeout:   timeout,
		sessionID: uuid.NewString(),
		httpClient: &http.Client{
			// Tope duro del entorno; el timeout efectivo se aplica por request.
			Timeout: timeout + 30*time.Second,
		},
	}
}

func (c *aiClient) IsConfigured() bool {
	return c.apiURL != "" && c.apiKey != ""
}

// envSettings son los ajustes por defecto tomados de .env.
func (c *aiClient) envSettings() Settings {
	temperature := defaultTemperature
	return Settings{
		Model:          c.model,
		Temperature:    &temperature,
		RequestTimeout: c.timeout,
	}
}

// resolve combina los ajustes de runtime sobre los del entorno.
func (c *aiClient) resolve(s Settings) (string, float64, time.Duration) {
	model, timeout := c.model, c.timeout
	temperature := defaultTemperature

	if s.Model != "" {
		model = s.Model
	}
	if s.Temperature != nil {
		temperature = *s.Temperature
	}
	if s.RequestTimeout > 0 {
		timeout = s.RequestTimeout
	}
	return model, temperature, timeout
}

func (c *aiClient) Chat(ctx context.Context, systemPrompt, userPrompt string, s Settings) (string, error) {
	model, temperature, timeout := c.resolve(s)

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error al serializar request: %w", err)
	}

	url := buildChatURL(c.apiURL)
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dev-test-api/1.0")
	req.Header.Set("x-opencode-session", c.sessionID)

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
