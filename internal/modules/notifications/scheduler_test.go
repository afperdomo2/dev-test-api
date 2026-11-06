package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextRun(t *testing.T) {
	loc := time.FixedZone("Test/UTC-5", -5*60*60)
	day := func(y int, m time.Month, d, h, min int) time.Time {
		return time.Date(y, m, d, h, min, 0, 0, loc)
	}

	tests := []struct {
		name string
		now  time.Time
		want time.Time
	}{
		{"antes de las 08:00 → hoy", day(2026, time.October, 20, 7, 30), day(2026, time.October, 20, 8, 0)},
		{"justo a las 08:00 → mañana", day(2026, time.October, 20, 8, 0), day(2026, time.October, 21, 8, 0)},
		{"después de las 08:00 → mañana", day(2026, time.October, 20, 9, 15), day(2026, time.October, 21, 8, 0)},
		{"cambio de mes", day(2026, time.October, 31, 23, 0), day(2026, time.November, 1, 8, 0)},
		{"cambio de año", day(2026, time.December, 31, 12, 0), day(2027, time.January, 1, 8, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, nextRun(tt.now).Equal(tt.want), "got %s, want %s", nextRun(tt.now), tt.want)
		})
	}
}

func TestSchedulerStartReturnsOnCancel(t *testing.T) {
	svc := NewService(newStoreMock(), nil)
	scheduler := NewScheduler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		scheduler.Start(ctx)
	}()

	cancel()

	// 10 s, no 2 s: en un hook pre-push con Go + web + mobile en paralelo, el
	// scheduler puede tardar más de 2 s en ser despertado por el runtime. El test
	// valida que `Start` retorna, no cuánto tarda en arrancar.
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Start no retornó al cancelar el context")
	}
}
