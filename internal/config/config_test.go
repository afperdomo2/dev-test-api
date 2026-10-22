package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// loadWith aísla el entorno: fija las variables obligatorias y deja el resto en
// blanco (equivale a no configuradas) salvo las del mapa `env`.
func loadWith(t *testing.T, env map[string]string) *Config {
	t.Helper()

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("GIN_MODE", "debug")
	t.Setenv("PORT", "8080")
	t.Setenv("MAX_BODY_BYTES", "")
	t.Setenv("TRUSTED_PROXIES", "")
	t.Setenv("HSTS_ENABLED", "")
	t.Setenv("SWAGGER_ENABLED", "")

	for k, v := range env {
		t.Setenv(k, v)
	}

	return Load()
}

func TestLoadHTTPHardeningDefaults(t *testing.T) {
	cfg := loadWith(t, nil)

	assert.Equal(t, int64(1<<20), cfg.MaxBodyBytes, "1 MiB por defecto")
	assert.Empty(t, cfg.TrustedProxies, "sin proxies configurados")
	assert.True(t, cfg.HSTSEnabled, "HSTS activo por defecto")
	assert.True(t, cfg.SwaggerEnabled, "Swagger habilitado en debug")
}

func TestLoadSwaggerGate(t *testing.T) {
	assert.False(t, loadWith(t, map[string]string{"GIN_MODE": "release"}).SwaggerEnabled,
		"release sin SWAGGER_ENABLED debe apagar Swagger")
	assert.True(t, loadWith(t, map[string]string{
		"GIN_MODE":        "release",
		"SWAGGER_ENABLED": "true",
	}).SwaggerEnabled, "la variable explícita gana en release")
	assert.False(t, loadWith(t, map[string]string{"SWAGGER_ENABLED": "false"}).SwaggerEnabled)
}

func TestLoadTrustedProxies(t *testing.T) {
	cfg := loadWith(t, map[string]string{
		"TRUSTED_PROXIES": " 10.0.0.1, ,172.16.0.0/12 ,",
	})

	assert.Equal(t, []string{"10.0.0.1", "172.16.0.0/12"}, cfg.TrustedProxies,
		"se recortan espacios y se descartan entradas vacías")
}

func TestLoadMaxBodyBytes(t *testing.T) {
	assert.Equal(t, int64(2048), loadWith(t, map[string]string{"MAX_BODY_BYTES": "2048"}).MaxBodyBytes)
	assert.Equal(t, int64(1<<20), loadWith(t, map[string]string{"MAX_BODY_BYTES": "no-es-numero"}).MaxBodyBytes,
		"valor inválido cae al default")
	assert.Equal(t, int64(1<<20), loadWith(t, map[string]string{"MAX_BODY_BYTES": "0"}).MaxBodyBytes,
		"0 no es un límite válido: cae al default")
}

func TestLoadHSTSFlag(t *testing.T) {
	assert.False(t, loadWith(t, map[string]string{"HSTS_ENABLED": "false"}).HSTSEnabled)
	assert.True(t, loadWith(t, map[string]string{"HSTS_ENABLED": "no-es-bool"}).HSTSEnabled,
		"valor inválido cae al default")
}
