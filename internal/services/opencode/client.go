package opencode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
)

type usageWindowRaw struct {
	Status   string  `json:"status"`
	Percent  float64 `json:"percent"`
	ResetsAt string  `json:"resetsAt"`
}

type usageResponseRaw struct {
	Usage struct {
		Rolling usageWindowRaw `json:"rolling"`
		Weekly  usageWindowRaw `json:"weekly"`
		Monthly usageWindowRaw `json:"monthly"`
	} `json:"usage"`
}

// UsageWindow representa una ventana de consumo de OpenCode Go.
type UsageWindow struct {
	Status   string  `json:"status"`
	Percent  float64 `json:"percent"`
	ResetsAt string  `json:"resetsAt"`
}

// UsageResponse representa el consumo rolling/semanal/mensual.
type UsageResponse struct {
	Rolling UsageWindow `json:"rolling"`
	Weekly  UsageWindow `json:"weekly"`
	Monthly UsageWindow `json:"monthly"`
}

// Client consulta el endpoint oficial de OpenCode Go.
type Client struct {
	usageURL   string
	apiKey     string
	httpClient *http.Client
}

func NewClient(cfg config.AIConfig, ocCfg config.OpenCodeConfig) *Client {
	return &Client{
		usageURL: ocCfg.UsageURL,
		apiKey:   cfg.APIKey,
		httpClient: &http.Client{
			Timeout: time.Duration(ocCfg.UsageTimeout) * time.Second,
		},
	}
}

func (c *Client) IsConfigured() bool {
	return c.usageURL != "" && c.apiKey != ""
}

// GetUsage llama a GET /zen/go/v1/usage con Bearer API key.
func (c *Client) GetUsage(ctx context.Context) (*UsageResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("OpenCode no está configurado: falta AI_API_KEY u OPENCODE_USAGE_URL")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.usageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al consultar OpenCode: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta de OpenCode: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ OpenCode usage error (%d) en %s: %s", resp.StatusCode, c.usageURL, string(body))
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("OpenCode API key inválida o expirada (401)")
		case http.StatusForbidden:
			return nil, fmt.Errorf("sin suscripción Go activa en esta API key (403)")
		default:
			return nil, fmt.Errorf("OpenCode API error (%d)", resp.StatusCode)
		}
	}

	var raw usageResponseRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta de OpenCode: %w", err)
	}

	return &UsageResponse{
		Rolling: UsageWindow(raw.Usage.Rolling),
		Weekly:  UsageWindow(raw.Usage.Weekly),
		Monthly: UsageWindow(raw.Usage.Monthly),
	}, nil
}
