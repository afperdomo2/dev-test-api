package push

import (
	"context"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
)

// Service registra/desregistra dispositivos. No conoce Expo: eso vive en el
// sender. También implementa notifications.Pusher (D-28).
type Service interface {
	Register(userID uuid.UUID, token, platform string) error
	Unregister(userID uuid.UUID, token string) (bool, error)
	Push(ctx context.Context, items []models.Notification) error
}

// sender entrega mensajes ya construidos (ExpoSender en producción).
type sender interface {
	Send(ctx context.Context, messages []expoMessage) error
}

type pushService struct {
	store  Store
	sender sender
}

func NewService(store Store, sender sender) Service {
	return &pushService{store: store, sender: sender}
}

func (s *pushService) Register(userID uuid.UUID, token, platform string) error {
	if err := s.store.Upsert(userID, token, platform); err != nil {
		return apierr.ErrInternal("Error al registrar el dispositivo", "")
	}
	return nil
}

func (s *pushService) Unregister(userID uuid.UUID, token string) (bool, error) {
	removed, err := s.store.DeleteByTokenAndUser(userID, token)
	if err != nil {
		return false, apierr.ErrInternal("Error al eliminar el dispositivo", "")
	}
	return removed, nil
}

func (s *pushService) Push(ctx context.Context, items []models.Notification) error {
	if len(items) == 0 || s.sender == nil {
		return nil
	}

	messages := make([]expoMessage, 0, len(items))
	for _, item := range items {
		tokens, err := s.store.FindByUser(item.UserID)
		if err != nil {
			return apierr.ErrInternal("Error al obtener los dispositivos del usuario", "")
		}
		for _, token := range tokens {
			messages = append(messages, expoMessage{
				To:    token.Token,
				Title: item.Title,
				Body:  item.Body,
				Data: map[string]string{
					"notificationId": item.ID.String(),
					"type":           item.Type,
				},
			})
		}
	}

	if len(messages) == 0 {
		return nil
	}
	return s.sender.Send(ctx, messages)
}
