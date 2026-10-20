package notifications

import (
	"context"
	"fmt"
	"log"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TypeDailyReview es el único tipo de aviso generado hoy (D-21: string + const).
const TypeDailyReview = "daily_review"

// Pusher entrega las notificaciones al móvil. `push` lo implementa; en web-only
// puede ser nil (SPEC-push, D-28).
type Pusher interface {
	Push(ctx context.Context, items []models.Notification) error
}

type Service interface {
	List(userID uuid.UUID, params ListNotificationsParams) ([]NotificationResponse, int64, error)
	UnreadCount(userID uuid.UUID) (int64, error)
	MarkRead(id, userID uuid.UUID) (*NotificationResponse, error)
	MarkAllRead(userID uuid.UUID) (int64, error)
	RunDailyJob(ctx context.Context) error
}

type notificationService struct {
	store  Store
	pusher Pusher
}

func NewService(store Store, pusher Pusher) Service {
	return &notificationService{store: store, pusher: pusher}
}

func (s *notificationService) List(userID uuid.UUID, params ListNotificationsParams) ([]NotificationResponse, int64, error) {
	items, total, err := s.store.FindPage(userID, params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las notificaciones", "")
	}

	result := make([]NotificationResponse, len(items))
	for i, item := range items {
		result[i] = ToNotificationResponse(item)
	}
	return result, total, nil
}

func (s *notificationService) UnreadCount(userID uuid.UUID) (int64, error) {
	count, err := s.store.CountUnread(userID)
	if err != nil {
		return 0, apierr.ErrInternal("Error al contar las notificaciones no leídas", "")
	}
	return count, nil
}

func (s *notificationService) MarkRead(id, userID uuid.UUID) (*NotificationResponse, error) {
	notification, err := s.store.FindByIDAndUser(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Notificacion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la notificacion", "")
	}

	// Idempotente: si ya estaba leída no se vuelve a escribir.
	if !notification.IsRead {
		if err := s.store.MarkRead(id, userID); err != nil {
			return nil, apierr.ErrInternal("Error al marcar la notificacion como leida", "")
		}
		notification.IsRead = true
	}

	resp := ToNotificationResponse(*notification)
	return &resp, nil
}

func (s *notificationService) MarkAllRead(userID uuid.UUID) (int64, error) {
	updated, err := s.store.MarkAllRead(userID)
	if err != nil {
		return 0, apierr.ErrInternal("Error al marcar las notificaciones como leidas", "")
	}
	return updated, nil
}

func (s *notificationService) RunDailyJob(ctx context.Context) error {
	due, err := s.store.DueReviewCounts()
	if err != nil {
		return apierr.ErrInternal("Error al calcular los repasos pendientes", "")
	}
	if len(due) == 0 {
		return nil
	}

	items := make([]models.Notification, 0, len(due))
	for _, d := range due {
		items = append(items, models.Notification{
			UserID: d.UserID,
			Type:   TypeDailyReview,
			Title:  "Repaso de hoy",
			Body:   fmt.Sprintf("Tienes %d preguntas para repasar", d.N),
		})
	}

	inserted, err := s.store.InsertBatch(items)
	if err != nil {
		return apierr.ErrInternal("Error al crear las notificaciones", "")
	}

	log.Printf("✅ Job diario de notificaciones: %d nuevas (%d usuarios con repaso pendiente)", len(inserted), len(due))

	// Best-effort: el aviso in-app ya existe; un fallo de push no rompe el job.
	if s.pusher != nil && len(inserted) > 0 {
		if pErr := s.pusher.Push(ctx, inserted); pErr != nil {
			log.Printf("⚠️ No se pudo entregar el push de notificaciones: %v", pErr)
		}
	}

	return nil
}
