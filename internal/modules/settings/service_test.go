package settings

import (
	"testing"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeStore struct {
	rows     []models.Setting
	findErr  error
	upserted []models.Setting
	upsertFn func(*models.Setting) error
}

func (s *fakeStore) FindAll() ([]models.Setting, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	return s.rows, nil
}

func (s *fakeStore) FindByKeys(keys []string) ([]models.Setting, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	wanted := make(map[string]bool, len(keys))
	for _, k := range keys {
		wanted[k] = true
	}
	var out []models.Setting
	for _, r := range s.rows {
		if wanted[r.Key] {
			out = append(out, r)
		}
	}
	return out, nil
}

func (s *fakeStore) Upsert(setting *models.Setting) error {
	if s.upsertFn != nil {
		return s.upsertFn(setting)
	}
	s.upserted = append(s.upserted, *setting)
	for i := range s.rows {
		if s.rows[i].Key == setting.Key {
			s.rows[i].Value = setting.Value
			return nil
		}
	}
	s.rows = append(s.rows, *setting)
	return nil
}

func newTestService(rows []models.Setting) (*service, *fakeStore) {
	store := &fakeStore{rows: rows}
	return &service{
		store:    store,
		registry: buildRegistry(config.AIConfig{Model: "gpt-4o-mini", RequestTimeout: 120}),
	}, store
}

func TestListReturnsEveryRegisteredKey(t *testing.T) {
	svc, _ := newTestService(nil)

	items, err := svc.List()
	require.NoError(t, err)
	assert.Len(t, items, 7, "deben listarse las 7 claves registradas aunque la tabla esté vacía")
	assert.Equal(t, "10", settingValue(t, items, KeyAIMaxPerJob))
	assert.Equal(t, "5", settingValue(t, items, KeyAIDefaultQuantity))
	assert.Equal(t, "2", settingValue(t, items, KeyAIWorkerConcurrency))
}

func TestListPrefersPersistedValueOverDefault(t *testing.T) {
	svc, _ := newTestService([]models.Setting{{Key: KeyAIMaxPerJob, Value: "4"}})

	items, err := svc.List()
	require.NoError(t, err)
	assert.Equal(t, "4", settingValue(t, items, KeyAIMaxPerJob))
}

func TestListFallsBackToDefaultWhenPersistedValueIsInvalid(t *testing.T) {
	// Una fila corrupta no debe romper la generación: gana el default.
	svc, _ := newTestService([]models.Setting{{Key: KeyAIMaxPerJob, Value: "999"}})

	items, err := svc.List()
	require.NoError(t, err)
	assert.Equal(t, "10", settingValue(t, items, KeyAIMaxPerJob))
}

func TestUpdateValidatesTypeAndRange(t *testing.T) {
	svc, _ := newTestService(nil)

	t.Run("rechaza un valor no numérico", func(t *testing.T) {
		_, err := svc.Update([]UpdateSettingRequest{{Key: KeyAIMaxPerJob, Value: "muchos"}}, uuid.New())
		requireAPIError(t, err, apierr.ErrValidation)
	})

	t.Run("rechaza un valor fuera de rango", func(t *testing.T) {
		_, err := svc.Update([]UpdateSettingRequest{{Key: KeyAIMaxPerJob, Value: "50"}}, uuid.New())
		requireAPIError(t, err, apierr.ErrValidation)
	})

	t.Run("rechaza un float en un campo entero", func(t *testing.T) {
		_, err := svc.Update([]UpdateSettingRequest{{Key: KeyAIMaxPerJob, Value: "2.5"}}, uuid.New())
		requireAPIError(t, err, apierr.ErrValidation)
	})

	t.Run("rechaza una clave desconocida", func(t *testing.T) {
		_, err := svc.Update([]UpdateSettingRequest{{Key: "ia.secreto", Value: "x"}}, uuid.New())
		requireAPIError(t, err, apierr.ErrValidation)
	})
}

func TestUpdateAcceptsZeroTemperature(t *testing.T) {
	// 0 es un valor válido: no puede tratarse como "sin configurar".
	svc, store := newTestService(nil)

	_, err := svc.Update([]UpdateSettingRequest{{Key: KeyAITemperature, Value: "0"}}, uuid.New())
	require.NoError(t, err)
	require.Len(t, store.upserted, 1)
	assert.Equal(t, "0", store.upserted[0].Value)
}

func TestUpdateRejectsDefaultQuantityAboveMaxPerJob(t *testing.T) {
	svc, _ := newTestService(nil)

	// max_per_job baja a 2 y default_quantity queda en 5: inválido.
	_, err := svc.Update([]UpdateSettingRequest{
		{Key: KeyAIMaxPerJob, Value: "2"},
		{Key: KeyAIDefaultQuantity, Value: "5"},
	}, uuid.New())
	requireAPIError(t, err, apierr.ErrValidation)
}

func TestUpdateRecordsWhoChangedTheSetting(t *testing.T) {
	adminID := uuid.New()
	svc, store := newTestService(nil)

	_, err := svc.Update([]UpdateSettingRequest{{Key: KeyAIMaxPerJob, Value: "7"}}, adminID)
	require.NoError(t, err)
	require.Len(t, store.upserted, 1)
	require.NotNil(t, store.upserted[0].UpdatedBy)
	assert.Equal(t, adminID, *store.upserted[0].UpdatedBy)
}

func TestAISettingsMapsModelTemperatureAndTimeout(t *testing.T) {
	svc, _ := newTestService([]models.Setting{
		{Key: KeyAIModel, Value: "gpt-4o"},
		{Key: KeyAITemperature, Value: "0.3"},
		{Key: KeyAIRequestTimeout, Value: "45"},
	})

	got, err := svc.AISettings()
	require.NoError(t, err)
	assert.Equal(t, "gpt-4o", got.Model)
	require.NotNil(t, got.Temperature)
	assert.InDelta(t, 0.3, *got.Temperature, 0.0001)
	assert.Equal(t, 45.0, got.RequestTimeout.Seconds())
}

func TestAISettingsUsesEnvFallbackWhenRowsMissing(t *testing.T) {
	store := &fakeStore{}
	svc := &service{
		store:    store,
		registry: buildRegistry(config.AIConfig{Model: "modelo-env", RequestTimeout: 90}),
	}

	got, err := svc.AISettings()
	require.NoError(t, err)
	assert.Equal(t, "modelo-env", got.Model, "sin fila, gana el valor de .env")
	assert.Equal(t, 90.0, got.RequestTimeout.Seconds())
}

func TestMaxPerDayIsNotCappedToThePerJobLimit(t *testing.T) {
	svc, _ := newTestService(nil)

	maxPerDay, err := svc.MaxPerDay()
	require.NoError(t, err)
	assert.Equal(t, 200, maxPerDay)
}

func settingValue(t *testing.T, items []SettingResponse, key string) string {
	t.Helper()
	for _, item := range items {
		if item.Key == key {
			return item.Value
		}
	}
	t.Fatalf("no se encontró el ajuste %s", key)
	return ""
}

func requireAPIError(t *testing.T, err error, check func(string, string) *apierr.APIError) {
	t.Helper()
	require.Error(t, err)
	apiErr, ok := err.(*apierr.APIError)
	require.True(t, ok, "se esperaba un *apierr.APIError, se obtuvo %T", err)
	expected := check("", "")
	assert.Equal(t, expected.Status, apiErr.Status)
}
