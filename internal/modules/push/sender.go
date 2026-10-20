package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
)

const (
	expoBatchSize    = 100
	expoBatchTimeout = 10 * time.Second
)

type expoMessage struct {
	To    string            `json:"to"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

type expoTicketDetails struct {
	Error string `json:"error"`
}

type expoTicket struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Details expoTicketDetails `json:"details"`
}

type expoResponse struct {
	Data []expoTicket `json:"data"`
}

// ExpoSender entrega mensajes a la Expo Push API en lotes de 100, con timeout
// por lote. Si falla, el error se propaga para que el job lo loguee y siga: el
// aviso in-app ya está persistido (best-effort).
type ExpoSender struct {
	url         string
	accessToken string
	client      *http.Client
	store       Store
}

func NewExpoSender(cfg config.PushConfig, store Store) *ExpoSender {
	return &ExpoSender{
		url:         cfg.ExpoURL,
		accessToken: cfg.AccessToken,
		client:      &http.Client{Timeout: 15 * time.Second},
		store:       store,
	}
}

func (s *ExpoSender) Send(ctx context.Context, messages []expoMessage) error {
	for start := 0; start < len(messages); start += expoBatchSize {
		end := start + expoBatchSize
		if end > len(messages) {
			end = len(messages)
		}

		if err := s.sendBatch(ctx, messages[start:end]); err != nil {
			log.Printf("❌ No se pudo enviar el lote de push (%d mensajes): %v", end-start, err)
			return err
		}
	}
	return nil
}

func (s *ExpoSender) sendBatch(ctx context.Context, batch []expoMessage) error {
	batchCtx, cancel := context.WithTimeout(ctx, expoBatchTimeout)
	defer cancel()

	payload, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("serializar mensajes: %w", err)
	}

	req, err := http.NewRequestWithContext(batchCtx, http.MethodPost, s.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("crear petición a Expo: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if s.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.accessToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("llamar a Expo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("leer respuesta de Expo: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Expo respondió %d", resp.StatusCode)
	}

	var parsed expoResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("decodificar respuesta de Expo: %w", err)
	}

	s.handleTickets(batch, parsed.Data)
	log.Printf("✅ Push enviado: %d mensajes", len(batch))
	return nil
}

// handleTickets borra los tokens que Expo marca como DeviceNotRegistered. Los
// demás errores de ticket se loguean y se descartan.
func (s *ExpoSender) handleTickets(batch []expoMessage, tickets []expoTicket) {
	for i, ticket := range tickets {
		if i >= len(batch) {
			break
		}

		if ticket.Details.Error == "DeviceNotRegistered" {
			token := batch[i].To
			if err := s.store.DeleteByToken(token); err != nil {
				log.Printf("❌ No se pudo borrar el token push inválido: %v", err)
				continue
			}
			log.Printf("🗑️ Token push inválido eliminado: %s", token)
			continue
		}

		if ticket.Status == "error" {
			log.Printf("❌ Expo rechazó un push: %s", ticket.Message)
		}
	}
}
