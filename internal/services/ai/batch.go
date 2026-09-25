package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	// maxExistingInPrompt acota cuántas preguntas previas se envían al prompt.
	// Más que esto infla el prompt sin evitar repeticiones de forma efectiva.
	maxExistingInPrompt = 60
)

// GeneratedQuestion es una pregunta devuelta por la IA, ya validada.
type GeneratedQuestion struct {
	Type           string            `json:"type"`
	Content        string            `json:"content"`
	Explanation    string            `json:"explanation"`
	Difficulty     string            `json:"difficulty"`
	Options        []GeneratedOption `json:"options,omitempty"`
	StarterCode    string            `json:"starterCode,omitempty"`
	ExpectedOutput string            `json:"expectedOutput,omitempty"`
	Language       string            `json:"language,omitempty"`
}

type GeneratedOption struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect"`
}

// BatchRequest describe una generación en lote. Count nunca supera el máximo
// por trabajo (ai.max_per_job) y se resuelve en una sola llamada a la IA.
type BatchRequest struct {
	TopicNames      []string
	Difficulty      string
	Count           int
	ExistingContent []string
}

// GenerateQuestions pide Count preguntas en una sola llamada y devuelve las
// que el modelo produjo válidas. Los elementos inválidos se descartan con un
// warning en vez de tumbar todo el trabajo: el llamador recibe lo que hay sido
// guardable y decide el estado final del job.
func (g *Generator) GenerateQuestions(ctx context.Context, req BatchRequest) ([]*GeneratedQuestion, error) {
	if !g.IsConfigured() {
		return nil, fmt.Errorf("la IA no está configurada (faltan AI_API_URL o AI_API_KEY)")
	}
	if req.Count <= 0 {
		return nil, fmt.Errorf("la cantidad de preguntas debe ser mayor que cero")
	}
	if len(req.TopicNames) == 0 {
		return nil, fmt.Errorf("se requiere al menos un tema para generar preguntas")
	}

	userPrompt := buildBatchUserPrompt(req)
	log.Printf("🤖 Generando %d preguntas en lote sobre: %s (dificultad: %s)",
		req.Count, strings.Join(req.TopicNames, ", "), req.Difficulty)

	responseText, err := g.client.Chat(ctx, batchSystemPrompt, userPrompt, g.runtimeSettings())
	if err != nil {
		return nil, err
	}

	questions, err := parseAIBatchResponse(responseText, req.Difficulty)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta de la IA: %w", err)
	}

	log.Printf("✅ Lote de %d preguntas válido: %d generadas de %d solicitadas",
		req.Count, len(questions), req.Count)
	return questions, nil
}

// parseAIBatchResponse parsea el array JSON devuelto por el modelo, descarta
// los elementos inválidos y normaliza la dificultad con la solicitada.
func parseAIBatchResponse(text, requestedDifficulty string) ([]*GeneratedQuestion, error) {
	var raw []GeneratedQuestion
	if err := json.Unmarshal([]byte(stripCodeFence(text)), &raw); err != nil {
		return nil, fmt.Errorf("JSON inválido: %w\nRespuesta: %s", err, truncate(text, 300))
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("la IA devolvió un array vacío")
	}

	fallbackDifficulty := requestedDifficulty
	if fallbackDifficulty == "" || fallbackDifficulty == "random" {
		fallbackDifficulty = "intermediate"
	}

	valid := make([]*GeneratedQuestion, 0, len(raw))
	for i := range raw {
		q := raw[i]
		if err := validateGeneratedQuestion(&q); err != nil {
			log.Printf("⚠️ Pregunta %d del lote descartada: %v", i+1, err)
			continue
		}
		if q.Difficulty == "" || q.Difficulty == "random" {
			q.Difficulty = fallbackDifficulty
		}
		valid = append(valid, &q)
	}

	if len(valid) == 0 {
		return nil, fmt.Errorf("ninguna de las %d preguntas devueltas superó la validación", len(raw))
	}

	return valid, nil
}

func buildBatchUserPrompt(req BatchRequest) string {
	var sb strings.Builder

	promptDifficulty := req.Difficulty
	if promptDifficulty == "random" {
		promptDifficulty = "variada (elige aleatoriamente entre beginner, intermediate y advanced)"
	}

	sb.WriteString(fmt.Sprintf("Genera EXACTAMENTE %d preguntas single_choice de dificultad \"%s\" sobre: %s.\n",
		req.Count, promptDifficulty, strings.Join(req.TopicNames, ", ")))
	sb.WriteString("Recuerda: solo tipo single_choice, 4 opciones, exactamente 1 correcta, TODO en español. Si alguna opción contiene código, usa el formato con triple backtick y lenguaje.\n")
	sb.WriteString("Cada pregunta debe tratar un concepto DISTINTO: no repitas ni reformules el mismo concepto.\n")

	if len(req.ExistingContent) > 0 {
		existing := req.ExistingContent
		if len(existing) > maxExistingInPrompt {
			existing = existing[len(existing)-maxExistingInPrompt:]
		}
		sb.WriteString("\nYa existen las siguientes preguntas para estos temas. NO las repitas:\n")
		for _, c := range existing {
			sb.WriteString(fmt.Sprintf("- %s\n", truncate(c, 120)))
		}
	}

	return sb.String()
}
