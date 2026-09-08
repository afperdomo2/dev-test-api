package opencode

import (
	"context"
	"strings"

	opencodeSvc "github.com/felipe/dev-test-api/internal/services/opencode"
	"github.com/felipe/dev-test-api/pkg/apierr"
)

type usageClient interface {
	IsConfigured() bool
	GetUsage(ctx context.Context) (*opencodeSvc.UsageResponse, error)
}

type Service interface {
	GetUsage(ctx context.Context) (*opencodeSvc.UsageResponse, error)
}

type service struct {
	client usageClient
}

func NewService(client usageClient) Service {
	return &service{client: client}
}

func (s *service) GetUsage(ctx context.Context) (*opencodeSvc.UsageResponse, error) {
	if !s.client.IsConfigured() {
		return nil, apierr.NewAPIError(503, "Service Unavailable", "OpenCode no está configurado en el servidor (falta AI_API_KEY)", "")
	}

	res, err := s.client.GetUsage(ctx)
	if err != nil {
		msg := err.Error()
		// Mapea errores conocidos a códigos HTTP útiles.
		switch {
		case strings.Contains(msg, "(401)"):
			return nil, apierr.NewAPIError(502, "Bad Gateway", "OpenCode API key inválida o expirada", "")
		case strings.Contains(msg, "(403)"):
			return nil, apierr.NewAPIError(502, "Bad Gateway", "Sin suscripción Go activa en esta API key", "")
		default:
			return nil, apierr.ErrInternal("No se pudo obtener el consumo de OpenCode: "+msg, "")
		}
	}

	return res, nil
}
