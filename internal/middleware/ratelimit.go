package middleware

import (
	"sync"
	"time"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Limit define el rate limit de un middleware: peticiones por minuto y ráfaga
// inicial. PerMinute=0 deshabilita el limiter.
type Limit struct {
	PerMinute int
	Burst     int
}

const (
	maxRateLimitVisitors = 10_000
	rateLimitVisitorTTL  = 10 * time.Minute
	rateLimitRetryAfter  = "60"
)

type rateLimitVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimitStore struct {
	mu       sync.Mutex
	visitors map[string]*rateLimitVisitor
}

func newRateLimitStore() *rateLimitStore {
	return &rateLimitStore{visitors: make(map[string]*rateLimitVisitor)}
}

// limiterFor devuelve (y refresca) el limiter de una clave. El barrido es
// perezoso: solo cuando el mapa crece demasiado, para no lanzar goroutines.
func (s *rateLimitStore) limiterFor(key string, perMinute, burst int) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if len(s.visitors) > maxRateLimitVisitors {
		s.sweepLocked(now)
	}

	if v, ok := s.visitors[key]; ok {
		v.lastSeen = now
		return v.limiter
	}

	if burst <= 0 {
		burst = perMinute
	}
	v := &rateLimitVisitor{
		limiter:  rate.NewLimiter(rate.Limit(float64(perMinute)/60), burst),
		lastSeen: now,
	}
	s.visitors[key] = v
	return v.limiter
}

func (s *rateLimitStore) sweepLocked(now time.Time) {
	for k, v := range s.visitors {
		if now.Sub(v.lastSeen) > rateLimitVisitorTTL {
			delete(s.visitors, k)
		}
	}
}

// RateLimit limita por clave (por defecto la IP del cliente). Al exceder
// devuelve 429 con Retry-After. Es un limiter en memoria: válido para una
// instancia (ver D-39).
func RateLimit(limit Limit, key func(*gin.Context) string) gin.HandlerFunc {
	store := newRateLimitStore()
	if key == nil {
		key = func(c *gin.Context) string { return c.ClientIP() }
	}

	return func(c *gin.Context) {
		if limit.PerMinute <= 0 {
			c.Next()
			return
		}

		if !store.limiterFor(key(c), limit.PerMinute, limit.Burst).Allow() {
			c.Header("Retry-After", rateLimitRetryAfter)
			response.Problem(c, apierr.ErrTooManyRequests(
				"Demasiados intentos. Espera un minuto e inténtalo de nuevo.",
				c.Request.URL.Path,
			))
			return
		}

		c.Next()
	}
}
