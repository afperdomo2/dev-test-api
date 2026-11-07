package settings

import (
	"strconv"
	"sync"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
)

// cacheTTL es cuánto se reutiliza un ajuste leído de la DB. El pool de
// generación consulta ai.worker_concurrency cada 3 s y por worker; sin caché
// eso son decenas de SELECT por minuto para un valor que casi nunca cambia.
// Update() invalida la caché, así que un cambio del admin se aplica al
// instante en este proceso y el TTL solo cubre la relectura "por si acaso".
const cacheTTL = 60 * time.Second

type Service interface {
	List() ([]SettingResponse, error)
	Update(items []UpdateSettingRequest, updatedBy uuid.UUID) ([]SettingResponse, error)
	AISettings() (ai.Settings, error)
	MaxPerJob() (int, error)
	MaxPerDay() (int, error)
	DefaultQuantity() (int, error)
	WorkerConcurrency() (int, error)
}

type service struct {
	store    Store
	registry map[string]definition

	mu    sync.RWMutex
	cache map[string]cachedValue
}

// cachedValue es un ajuste ya resuelto (fila de la DB o default del registro)
// con su fecha de caducidad.
type cachedValue struct {
	value     string
	expiresAt time.Time
}

func NewService(store Store, aiCfg config.AIConfig) Service {
	registry := buildRegistry(aiCfg)
	return &service{
		store:    store,
		registry: registry,
		cache:    make(map[string]cachedValue, len(registry)),
	}
}

// resolve devuelve el valor efectivo de un ajuste: el de la base de datos si
// existe y es válido, o el default del registro.
func (s *service) resolve(key string, rows map[string]models.Setting) (string, error) {
	def, ok := s.registry[key]
	if !ok {
		return "", apierr.ErrValidation("Ajuste desconocido: "+key, "")
	}
	if row, exists := rows[key]; exists && row.Value != "" {
		if value, err := def.validate(key, row.Value); err == nil {
			return value, nil
		}
	}
	return def.Default, nil
}

func (s *service) loadResolved() (map[string]string, error) {
	rows, err := s.store.FindByKeys(knownKeys(s.registry))
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]models.Setting, len(rows))
	for _, r := range rows {
		byKey[r.Key] = r
	}

	resolved := make(map[string]string, len(s.registry))
	for _, key := range knownKeys(s.registry) {
		value, rErr := s.resolve(key, byKey)
		if rErr != nil {
			return nil, rErr
		}
		resolved[key] = value
	}
	return resolved, nil
}

func (s *service) List() ([]SettingResponse, error) {
	rows, err := s.store.FindAll()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la configuración", "")
	}
	byKey := make(map[string]models.Setting, len(rows))
	for _, r := range rows {
		byKey[r.Key] = r
	}

	keys := knownKeys(s.registry)
	out := make([]SettingResponse, 0, len(keys))
	for _, key := range keys {
		def := s.registry[key]
		value, rErr := s.resolve(key, byKey)
		if rErr != nil {
			return nil, rErr
		}

		item := SettingResponse{
			Key:         key,
			Value:       value,
			Type:        string(def.Type),
			Description: def.Description,
		}
		if def.Type != typeString {
			min, max := def.Min, def.Max
			item.Min, item.Max = &min, &max
		}
		if row, exists := byKey[key]; exists {
			updatedAt := row.UpdatedAt
			item.UpdatedAt = &updatedAt
			item.UpdatedBy = row.UpdatedBy
		}
		out = append(out, item)
	}

	return out, nil
}

func (s *service) Update(items []UpdateSettingRequest, updatedBy uuid.UUID) ([]SettingResponse, error) {
	if len(items) == 0 {
		return nil, apierr.ErrValidation("No se recibieron ajustes para actualizar", "")
	}

	// Normaliza y valida cada par clave/valor contra el registro.
	normalized := make(map[string]string, len(items))
	for _, item := range items {
		def, ok := s.registry[item.Key]
		if !ok {
			return nil, apierr.ErrValidation("Ajuste desconocido: "+item.Key, "")
		}
		value, err := def.validate(item.Key, item.Value)
		if err != nil {
			return nil, err
		}
		normalized[item.Key] = value
	}

	// La cantidad por defecto no puede superar el máximo por trabajo.
	effective, err := s.loadResolved()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la configuración actual", "")
	}
	for key, value := range normalized {
		effective[key] = value
	}
	defaultQuantity, _ := strconv.Atoi(effective[KeyAIDefaultQuantity])
	maxPerJob, _ := strconv.Atoi(effective[KeyAIMaxPerJob])
	if defaultQuantity > maxPerJob {
		return nil, apierr.ErrValidation(
			"ai.default_quantity no puede superar ai.max_per_job", "",
		)
	}

	// Cualquier escritura invalida la caché aunque falle a mitad: las filas
	// anteriores ya se persistieron.
	defer s.invalidateCache()

	for key, value := range normalized {
		row := &models.Setting{
			Key:         key,
			Value:       value,
			Description: s.registry[key].Description,
			UpdatedBy:   &updatedBy,
		}
		if err := s.store.Upsert(row); err != nil {
			return nil, apierr.ErrInternal("Error al guardar la configuración", "")
		}
	}

	return s.List()
}

// getOne resuelve un único ajuste (fila si existe, default del registro si no).
// Va a la DB como mucho una vez por clave cada cacheTTL.
func (s *service) getOne(key string) (string, error) {
	if value, ok := s.cached(key); ok {
		return value, nil
	}

	value, err := s.loadOne(key)
	if err != nil {
		// Un error no se cachea: el siguiente intento vuelve a consultar.
		return "", err
	}

	s.remember(key, value)
	return value, nil
}

// loadOne es la lectura real: fila persistida o default del registro.
func (s *service) loadOne(key string) (string, error) {
	rows, err := s.store.FindByKeys([]string{key})
	if err != nil {
		return "", err
	}
	byKey := make(map[string]models.Setting, len(rows))
	for _, r := range rows {
		byKey[r.Key] = r
	}
	return s.resolve(key, byKey)
}

func (s *service) cached(key string) (string, bool) {
	s.mu.RLock()
	entry, ok := s.cache[key]
	s.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.value, true
}

func (s *service) remember(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// El mapa puede ser nil si el service se construyó a mano (tests).
	if s.cache == nil {
		s.cache = make(map[string]cachedValue, len(s.registry))
	}
	s.cache[key] = cachedValue{value: value, expiresAt: time.Now().Add(cacheTTL)}
}

// invalidateCache se llama tras cualquier escritura: lo cacheado deja de
// reflejar la DB.
func (s *service) invalidateCache() {
	s.mu.Lock()
	clear(s.cache)
	s.mu.Unlock()
}

func (s *service) GetInt(key string) (int, error) {
	value, err := s.getOne(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (s *service) GetFloat(key string) (float64, error) {
	value, err := s.getOne(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(value, 64)
}

func (s *service) GetString(key string) (string, error) {
	return s.getOne(key)
}

func (s *service) AISettings() (ai.Settings, error) {
	temperature, err := s.GetFloat(KeyAITemperature)
	if err != nil {
		return ai.Settings{}, err
	}
	model, err := s.GetString(KeyAIModel)
	if err != nil {
		return ai.Settings{}, err
	}
	timeout, err := s.GetInt(KeyAIRequestTimeout)
	if err != nil {
		return ai.Settings{}, err
	}

	return ai.Settings{
		Model:          model,
		Temperature:    &temperature,
		RequestTimeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (s *service) MaxPerJob() (int, error)         { return s.GetInt(KeyAIMaxPerJob) }
func (s *service) MaxPerDay() (int, error)         { return s.GetInt(KeyAIMaxPerDay) }
func (s *service) DefaultQuantity() (int, error)   { return s.GetInt(KeyAIDefaultQuantity) }
func (s *service) WorkerConcurrency() (int, error) { return s.GetInt(KeyAIWorkerConcurrency) }
