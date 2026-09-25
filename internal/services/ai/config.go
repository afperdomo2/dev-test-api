package ai

import "time"

// Settings son los ajustes de IA resueltos en runtime (tabla `settings`).
// Se leen en cada llamada para que un cambio hecho por el administrador surta
// efecto sin reiniciar el servidor.
//
// Temperature es puntero porque 0 es un valor válido (generación determinista)
// y no puede usarse como "sin configurar".
type Settings struct {
	Model          string
	Temperature    *float64
	RequestTimeout time.Duration
}

// SettingsProvider resuelve los ajustes de IA. Lo implementa el módulo
// internal/modules/settings.
type SettingsProvider interface {
	AISettings() (Settings, error)
}
