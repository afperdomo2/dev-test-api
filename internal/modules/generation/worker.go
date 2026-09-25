package generation

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"gorm.io/gorm"
)

// pollInterval es con qué frecuencia los workers buscan trabajo pendiente. Es
// polling simple y suficiente: la encolación no necesita notificar a los
// workers porque no hay cola en memoria.
const pollInterval = 3 * time.Second

// Worker es el pool que ejecuta los trabajos de generación. Cada worker reclama
// el pending más antiguo y lo procesa en una única llamada batch a la IA.
//
// Un trabajo = una llamada. El paralelismo se aplica entre trabajos, nunca
// entre preguntas del mismo trabajo.
type Worker struct {
	service  Service
	interval time.Duration
	// concurrencySource relee ai.worker_concurrency para poder reescalar en
	// caliente. Se consulta entre trabajos, así que un cambio de ajuste surte
	// efecto sin reiniciar y sin matar ningún trabajo a medio generar.
	concurrencySource func() int

	mu          sync.Mutex
	concurrency int
	handles     []*handle
	poolCtx     context.Context

	// stop cierra el pool por completo (shutdown del servidor).
	stop     context.CancelFunc
	finished sync.WaitGroup
}

// handle es un worker individual del pool. quit se cierra cuando el admin baja
// la concurrencia y ese worker sobra.
type handle struct {
	id   int
	quit chan struct{}
}

func NewWorker(service Service, concurrency int, source func() int) *Worker {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Worker{
		service:           service,
		interval:          pollInterval,
		concurrency:       concurrency,
		concurrencySource: source,
	}
}

// Concurrency devuelve el tamaño actual del pool.
func (w *Worker) Concurrency() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.concurrency
}

// reconcile alinea el pool con el ajuste vigente. Un error de lectura se ignora
// a propósito: el tamaño actual sigue siendo válido.
func (w *Worker) reconcile() {
	if w.concurrencySource == nil {
		return
	}
	target := w.concurrencySource()
	if target < 1 {
		return
	}
	if target != w.Concurrency() {
		w.SetConcurrency(target)
	}
}

// SetConcurrency reescala el pool en caliente. Bajar la concurrencia no mata
// trabajos a medio generar: los workers sobrantes salen entre un trabajo y el
// siguiente.
func (w *Worker) SetConcurrency(target int) {
	if target < 1 {
		target = 1
	}

	w.mu.Lock()
	previous := w.concurrency
	w.concurrency = target
	handles := w.handles
	running := len(handles)
	w.mu.Unlock()

	switch {
	case target > running:
		log.Printf("⚙️ Pool de generación: %d → %d workers", previous, target)
		w.mu.Lock()
		for i := 0; i < target-running; i++ {
			w.spawnLocked(w.poolCtx, running+i+1)
		}
		w.mu.Unlock()
	case target < running:
		log.Printf("⚙️ Pool de generación: %d → %d workers", previous, target)
		w.mu.Lock()
		// Cierra los handles sobrantes (los de índice más alto).
		for _, h := range handles[target:] {
			close(h.quit)
		}
		w.handles = handles[:target]
		w.mu.Unlock()
	}
}

// Start recupera los trabajos interrumpidos por un reinicio previo y lanza el
// pool. Bloquea hasta que se para el servidor, para que el shutdown sea
// ordenado.
func (w *Worker) Start(ctx context.Context) {
	if affected, err := w.service.RecoverRunning(); err != nil {
		log.Printf("❌ No se pudieron recuperar los trabajos en curso: %v", err)
	} else if affected > 0 {
		log.Printf("🛢️ %d trabajo(s) reencolados tras un reinicio", affected)
	}

	poolCtx, cancel := context.WithCancel(ctx)
	w.mu.Lock()
	w.stop = cancel
	w.poolCtx = poolCtx
	for i := 1; i <= w.concurrency; i++ {
		w.spawnLocked(poolCtx, i)
	}
	size := len(w.handles)
	w.mu.Unlock()

	log.Printf("🚀 Pool de generación con IA iniciado (%d workers)", size)

	w.finished.Wait()
	log.Println("✅ Pool de generación con IA detenido")
}

// spawnLocked añade un worker al pool. Requiere el lock tomado.
func (w *Worker) spawnLocked(ctx context.Context, id int) {
	h := &handle{id: id, quit: make(chan struct{})}
	w.handles = append(w.handles, h)
	w.finished.Add(1)
	go w.run(ctx, h)
}

func (w *Worker) run(ctx context.Context, h *handle) {
	defer w.finished.Done()

	log.Printf("🤖 Worker de generación %d listo", h.id)
	for {
		// Un worker sobredimensionado sale entre trabajos, nunca a medio generar.
		if closed(h.quit) || ctx.Err() != nil {
			log.Printf("🤖 Worker de generación %d detenido", h.id)
			return
		}

		w.reconcile()

		job, err := w.service.NextPendingJob()
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			if !w.wait(ctx) {
				return
			}
			continue
		case err != nil:
			log.Printf("❌ Worker %d: error al buscar trabajo pendiente: %v", h.id, err)
			if !w.wait(ctx) {
				return
			}
			continue
		}

		w.process(ctx, h.id, job)
	}
}

func (w *Worker) process(ctx context.Context, id int, job *models.GenerationJob) {
	log.Printf("▶️ Worker %d: procesando trabajo %s (%d preguntas, dificultad %s)",
		id, job.ID, job.Quantity, job.Difficulty)

	if err := w.service.ProcessJob(ctx, *job); err != nil {
		// Si el error viene del cierre del servidor el trabajo no se marca
		// fallido: se queda running y lo recupera RecoverRunning al arrancar.
		if ctx.Err() != nil {
			log.Printf("⚠️ Worker %d: trabajo %s interrumpido por cierre del servidor", id, job.ID)
			return
		}
		log.Printf("❌ Worker %d: falló el trabajo %s: %v", id, job.ID, err)
		if failErr := w.service.FailJob(job.ID, truncateMessage(err.Error(), 500)); failErr != nil {
			log.Printf("❌ Worker %d: no se pudo marcar el trabajo %s como fallido: %v", id, job.ID, failErr)
		}
	}
}

// wait espera el intervalo de polling. Devuelve false si el pool se está parando.
func (w *Worker) wait(ctx context.Context) bool {
	timer := time.NewTimer(w.interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func closed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func truncateMessage(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
