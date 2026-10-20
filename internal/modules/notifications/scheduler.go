package notifications

import (
	"context"
	"log"
	"time"
)

// dailyHour es la hora local del servidor a la que corre el fan-out diario.
const dailyHour = 8

type Scheduler struct {
	service Service
}

func NewScheduler(service Service) *Scheduler {
	return &Scheduler{service: service}
}

// Start bloquea hasta que ctx se cancele (mismo patrón que el worker de
// generación: goroutine de larga vida + context cancelable, sin errgroup).
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("🚀 Scheduler de notificaciones iniciado (%02d:00 local)", dailyHour)

	for {
		timer := time.NewTimer(time.Until(nextRun(time.Now())))
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("🛑 Scheduler de notificaciones detenido")
			return
		case <-timer.C:
			if err := s.service.RunDailyJob(ctx); err != nil {
				log.Printf("❌ Job diario de notificaciones: %v", err)
			}
		}
	}
}

// nextRun devuelve hoy a las dailyHour:00 locales, o mañana si ya pasó.
func nextRun(now time.Time) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), dailyHour, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}
