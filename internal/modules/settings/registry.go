package settings

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/pkg/apierr"
)

const (
	KeyAIMaxPerJob         = "ai.max_per_job"
	KeyAIMaxPerDay         = "ai.max_per_day"
	KeyAIDefaultQuantity   = "ai.default_quantity"
	KeyAIModel             = "ai.model"
	KeyAITemperature       = "ai.temperature"
	KeyAIRequestTimeout    = "ai.request_timeout"
	KeyAIWorkerConcurrency = "ai.worker_concurrency"
)

type valueType string

const (
	typeString valueType = "string"
	typeInt    valueType = "int"
	typeFloat  valueType = "float"
)

type definition struct {
	Type        valueType
	Default     string
	Min         float64
	Max         float64
	Description string
}

// registry es la fuente de verdad de qué ajustes existen, de qué tipo son y
// qué rango aceptan. Los defaults de ai.model y ai.request_timeout provienen
// del entorno (.env) porque no tiene sentido persistir el valor por defecto
// de infraestructura en la base de datos.
func buildRegistry(cfg config.AIConfig) map[string]definition {
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}
	timeout := strconv.Itoa(cfg.RequestTimeout)
	if cfg.RequestTimeout <= 0 {
		timeout = "120"
	}

	return map[string]definition{
		KeyAIMaxPerJob: {
			Type:        typeInt,
			Default:     "10",
			Min:         1,
			Max:         10,
			Description: "Máximo de preguntas por solicitud de generación con IA",
		},
		KeyAIMaxPerDay: {
			Type:        typeInt,
			Default:     "200",
			Min:         1,
			Max:         5000,
			Description: "Máximo de preguntas reservadas por día para generación con IA (tope global)",
		},
		KeyAIDefaultQuantity: {
			Type:        typeInt,
			Default:     "5",
			Min:         1,
			Max:         10,
			Description: "Cantidad de preguntas preseleccionada por defecto al encolar un trabajo",
		},
		KeyAIModel: {
			Type:        typeString,
			Default:     model,
			Description: "Modelo de IA usado para generar preguntas",
		},
		KeyAITemperature: {
			Type:        typeFloat,
			Default:     "0.8",
			Min:         0,
			Max:         2,
			Description: "Temperatura del modelo de IA",
		},
		KeyAIRequestTimeout: {
			Type:        typeInt,
			Default:     timeout,
			Min:         5,
			Max:         600,
			Description: "Timeout de la llamada a la API de IA, en segundos",
		},
		KeyAIWorkerConcurrency: {
			Type:        typeInt,
			Default:     "2",
			Min:         1,
			Max:         8,
			Description: "Número de workers que procesan trabajos de generación en paralelo",
		},
	}
}

// validate normaliza y valida el valor de un ajuste contra su definición.
func (d definition) validate(key, raw string) (string, error) {
	if d.Type == typeString {
		if raw == "" {
			return "", apierr.ErrValidation(fmt.Sprintf("%s no puede estar vacío", key), "")
		}
		return raw, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return "", apierr.ErrValidation(fmt.Sprintf("%s debe ser un número", key), "")
	}
	if value < d.Min || value > d.Max {
		return "", apierr.ErrValidation(
			fmt.Sprintf("%s debe estar entre %s y %s", key, formatBound(d.Min), formatBound(d.Max)),
			"",
		)
	}
	if d.Type == typeInt {
		if value != float64(int(value)) {
			return "", apierr.ErrValidation(fmt.Sprintf("%s debe ser un número entero", key), "")
		}
		return strconv.Itoa(int(value)), nil
	}
	return strconv.FormatFloat(value, 'f', -1, 64), nil
}

func formatBound(v float64) string {
	if v == float64(int(v)) {
		return strconv.Itoa(int(v))
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func knownKeys(registry map[string]definition) []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
