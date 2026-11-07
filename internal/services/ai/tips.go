package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

const (
	// maxTipContentLength acota el largo de un tip. Los tips son consejos de
	// 1-3 oraciones; más que esto es señal de que el modelo se desvió.
	maxTipContentLength = 600
	// maxExistingTipsInPrompt acota cuántos tips previos se envían al prompt.
	maxExistingTipsInPrompt = 60
)

// GeneratedTip es un tip devuelto por la IA, ya validado.
type GeneratedTip struct {
	Content string `json:"content"`
	// Topic es el slug del tema al que pertenece el tip. Debe ser uno de los
	// slugs pedidos; el store descarta los que no matcheen.
	Topic string `json:"topic"`
}

// TipBatchRequest describe una generación de tips en lote. Count nunca supera
// el máximo por trabajo (ai.max_per_job) y se resuelve en una sola llamada.
type TipBatchRequest struct {
	TopicNames []string
	TopicSlugs []string
	Count      int
	// ExistingContent son tips ya guardados en los temas, para evitar
	// duplicados.
	ExistingContent []string
}

// GenerateTips pide Count tips en una sola llamada. Los elementos inválidos se
// descartan con un warning en vez de tumbar todo el trabajo.
func (g *Generator) GenerateTips(ctx context.Context, req TipBatchRequest) ([]*GeneratedTip, error) {
	if !g.IsConfigured() {
		return nil, fmt.Errorf("la IA no está configurada (faltan AI_API_URL o AI_API_KEY)")
	}
	if req.Count <= 0 {
		return nil, fmt.Errorf("la cantidad de tips debe ser mayor que cero")
	}
	if len(req.TopicNames) == 0 || len(req.TopicNames) != len(req.TopicSlugs) {
		return nil, fmt.Errorf("se requiere al menos un tema con su slug para generar tips")
	}

	userPrompt := buildTipsUserPrompt(req)
	log.Printf("🤖 Generando %d tips en lote sobre: %s", req.Count, strings.Join(req.TopicNames, ", "))

	responseText, err := g.client.Chat(ctx, tipsSystemPrompt, userPrompt, g.runtimeSettings())
	if err != nil {
		return nil, err
	}

	tips, err := parseAITipsResponse(responseText)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta de la IA: %w", err)
	}

	log.Printf("✅ Lote de %d tips válido: %d generados de %d solicitados",
		req.Count, len(tips), req.Count)
	return tips, nil
}

// parseAITipsResponse parsea el array JSON devuelto por el modelo y descarta
// los elementos inválidos.
func parseAITipsResponse(text string) ([]*GeneratedTip, error) {
	var raw []GeneratedTip
	if err := json.Unmarshal([]byte(stripCodeFence(text)), &raw); err != nil {
		return nil, fmt.Errorf("JSON inválido: %w\nRespuesta: %s", err, truncate(text, 300))
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("la IA devolvió un array vacío")
	}

	valid := make([]*GeneratedTip, 0, len(raw))
	for i := range raw {
		tip := raw[i]
		if err := validateGeneratedTip(&tip); err != nil {
			log.Printf("⚠️ Tip %d del lote descartado: %v", i+1, err)
			continue
		}
		valid = append(valid, &tip)
	}

	if len(valid) == 0 {
		return nil, fmt.Errorf("ninguno de los %d tips devueltos superó la validación", len(raw))
	}

	return valid, nil
}

// validateGeneratedTip normaliza y valida un tip devuelto por la IA.
func validateGeneratedTip(t *GeneratedTip) error {
	t.Content = strings.TrimSpace(t.Content)
	t.Topic = strings.TrimSpace(t.Topic)

	if t.Content == "" {
		return fmt.Errorf("contenido de tip vacío")
	}
	if utf8.RuneCountInString(t.Content) > maxTipContentLength {
		return fmt.Errorf("contenido de tip demasiado largo (%d caracteres)", utf8.RuneCountInString(t.Content))
	}
	if t.Topic == "" {
		return fmt.Errorf("tema del tip vacío")
	}
	return nil
}

func buildTipsUserPrompt(req TipBatchRequest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Genera EXACTAMENTE %d tips sobre los siguientes temas:\n", req.Count))
	for i, name := range req.TopicNames {
		sb.WriteString(fmt.Sprintf("- %s (slug: %s)\n", name, req.TopicSlugs[i]))
	}

	if len(req.TopicNames) > 1 {
		sb.WriteString(fmt.Sprintf(
			"Distribuye los %d tips entre los %d temas. Varía los temas para cubrirlos todos.\n",
			req.Count, len(req.TopicNames),
		))
	}

	sb.WriteString("Recuerda: tips cortos (1-3 oraciones), accionables y específicos; cada concepto DISTINTO; " +
		"\"topic\" debe ser uno de los slugs de la lista; TODO en español.\n")

	if len(req.ExistingContent) > 0 {
		existing := req.ExistingContent
		if len(existing) > maxExistingTipsInPrompt {
			existing = existing[len(existing)-maxExistingTipsInPrompt:]
		}
		sb.WriteString("\nYa existen los siguientes tips para estos temas. NO los repitas ni reformules:\n")
		for _, c := range existing {
			sb.WriteString(fmt.Sprintf("- %s\n", truncate(c, 120)))
		}
	}

	return sb.String()
}
