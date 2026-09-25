package settings

import (
	"strconv"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
)

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
}

func NewService(store Store, aiCfg config.AIConfig) Service {
	return &service{store: store, registry: buildRegistry(aiCfg)}
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
func (s *service) getOne(key string) (string, error) {
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
