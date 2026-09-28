package ai

import (
	"fmt"
	"log"
	"strings"

	"github.com/felipe/dev-test-api/internal/config"
	"gorm.io/gorm"
)

type Generator struct {
	db       *gorm.DB
	client   *aiClient
	settings SettingsProvider
}

func NewGenerator(db *gorm.DB, cfg config.AIConfig, settings SettingsProvider) *Generator {
	return &Generator{
		db:       db,
		client:   newAIClient(cfg),
		settings: settings,
	}
}

func (g *Generator) IsConfigured() bool {
	return g.client.IsConfigured()
}

// runtimeSettings resuelve los ajustes de IA en cada llamada. Si el módulo de
// settings no está disponible o falla, se usan los defaults del entorno.
func (g *Generator) runtimeSettings() Settings {
	if g.settings == nil {
		return g.client.envSettings()
	}
	resolved, err := g.settings.AISettings()
	if err != nil {
		log.Printf("⚠️ No se pudieron leer los ajustes de IA, se usan los de .env: %v", err)
		return g.client.envSettings()
	}
	return resolved
}

// stripCodeFence quita la envoltura ```json ... ``` que añaden algunos modelos.
func stripCodeFence(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}
	return text
}

var validGeneratedTypes = map[string]bool{
	"single_choice":   true,
	"multiple_choice": true,
	"code_completion": true,
}

// validateGeneratedQuestion normaliza y valida una pregunta devuelta por la IA.
func validateGeneratedQuestion(q *GeneratedQuestion) error {
	if q.Type == "" {
		return fmt.Errorf("tipo de pregunta vacío")
	}
	if q.Content == "" {
		return fmt.Errorf("contenido de pregunta vacío")
	}
	if !validGeneratedTypes[q.Type] {
		return fmt.Errorf("tipo de pregunta no válido: %s", q.Type)
	}

	if q.Type == "code_completion" {
		if q.Language == "" {
			q.Language = "plaintext"
		}
		return nil
	}

	if len(q.Options) < 2 {
		return fmt.Errorf("se requieren al menos 2 opciones para preguntas de tipo choice")
	}

	correct := 0
	for i := range q.Options {
		if q.Options[i].Content == "" {
			return fmt.Errorf("opción %d sin contenido", i+1)
		}
		if q.Options[i].IsCorrect {
			correct++
		}
	}
	if correct != 1 {
		return fmt.Errorf("se requiere exactamente 1 opción correcta (se encontraron %d)", correct)
	}

	return nil
}
