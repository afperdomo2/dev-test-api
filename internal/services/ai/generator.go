package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
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

type DailyLimitError struct {
	Limit int
}

func (e *DailyLimitError) Error() string {
	return fmt.Sprintf("límite diario de IA alcanzado (%d)", e.Limit)
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

// evaluateAiQuota decide si una generación con IA está permitida. Los administradores
// están exentos del cupo diario. Ante cualquier error de base de datos el error se
// propaga (fail-closed) para no generar IA sin poder verificar el límite.
func evaluateAiQuota(isAdmin bool, limit int, used int64) error {
	if isAdmin {
		return nil
	}
	if int(used) >= limit {
		return &DailyLimitError{Limit: limit}
	}
	return nil
}

func (g *Generator) checkDailyAiLimit(userID uuid.UUID) error {
	var user models.User
	if err := g.db.Select("is_admin", "daily_ai_limit").First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("error al verificar el cupo diario de IA: %w", err)
	}
	// Salida temprana: exime al admin del conteo (y por tanto de cualquier error en esa consulta).
	if user.IsAdmin {
		return nil
	}
	y, m, d := time.Now().UTC().Date()
	since := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	var used int64
	if err := g.db.Model(&models.Question{}).
		Where("user_id = ?", userID).
		Where("source = ?", "ai_generated").
		Where("created_at >= ?", since).
		Count(&used).Error; err != nil {
		return fmt.Errorf("error al contar el consumo de IA: %w", err)
	}
	return evaluateAiQuota(user.IsAdmin, user.DailyAiLimit, used)
}

func (g *Generator) GenerateQuestion(session *models.Session) error {
	if !g.IsConfigured() {
		log.Println("⚠️ AI no configurada — omitiendo generación")
		return nil
	}

	if len(session.Topics) == 0 {
		return fmt.Errorf("la sesión no tiene temas asociados")
	}

	if err := g.checkDailyAiLimit(session.UserID); err != nil {
		if _, ok := err.(*DailyLimitError); ok {
			log.Printf("🚫 Límite diario de IA alcanzado para usuario %s — omitiendo generación", session.UserID)
		}
		return err
	}

	start := time.Now()

	existingContent, err := g.existingContent(session)
	if err != nil {
		log.Printf("⚠️ No se pudieron obtener preguntas existentes: %v", err)
	}

	userPrompt := buildUserPrompt(session, existingContent)
	log.Printf("🤖 Generando pregunta para sesión %s (tema: %s, dificultad: %s)", session.ID, safeTopicName(session.Topics), session.Difficulty)

	responseText, err := g.client.Chat(context.Background(), systemPrompt, userPrompt, g.runtimeSettings())
	if err != nil {
		return err
	}

	generated, err := parseAIResponse(responseText)
	if err != nil {
		return fmt.Errorf("error al parsear respuesta de IA: %w", err)
	}

	question, err := g.saveQuestion(session, generated)
	if err != nil {
		return fmt.Errorf("error al guardar pregunta: %w", err)
	}

	result := g.db.Model(&models.Session{}).
		Where("id = ? AND (question_limit IS NULL OR questions_generated < question_limit)", session.ID).
		UpdateColumn("questions_generated", gorm.Expr("questions_generated + 1"))
	if result.Error != nil {
		return fmt.Errorf("error al incrementar contador: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("🛢️ Sesión %s: límite de preguntas alcanzado (%d/%s), pregunta huérfana generada", session.ID, session.QuestionsGenerated, limitStr(session))
	}
	session.QuestionsGenerated++

	log.Printf("✅ Pregunta generada para sesión %s en %v: %s (%d/%s)",
		session.ID,
		time.Since(start).Round(time.Millisecond),
		truncate(question.Content, 80),
		session.QuestionsGenerated,
		limitStr(session),
	)

	return nil
}

func (g *Generator) existingContent(session *models.Session) ([]string, error) {
	topicIDs := make([]uuid.UUID, len(session.Topics))
	for i, t := range session.Topics {
		topicIDs[i] = t.ID
	}

	var contents []string
	err := g.db.Model(&models.Question{}).
		Joins("JOIN question_topics ON question_topics.question_id = questions.id").
		Where("question_topics.topic_id IN ?", topicIDs).
		Where("questions.deleted_at IS NULL").
		Pluck("questions.content", &contents).Error
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (g *Generator) saveQuestion(session *models.Session, gen *GeneratedQuestion) (*models.Question, error) {
	source := "ai_generated"
	difficulty := gen.Difficulty
	if difficulty == "" || difficulty == "random" {
		if session.Difficulty == "random" {
			difficulty = "intermediate"
		} else {
			difficulty = session.Difficulty
		}
	}

	topicIDs := make([]uuid.UUID, len(session.Topics))
	for i, t := range session.Topics {
		topicIDs[i] = t.ID
	}

	question := &models.Question{
		UserID:      session.UserID,
		Type:        gen.Type,
		Content:     gen.Content,
		Explanation: gen.Explanation,
		Difficulty:  difficulty,
		Source:      source,
	}

	if gen.Type == "code_completion" {
		question.CodeChallenge = &models.CodeChallenge{
			StarterCode:    gen.StarterCode,
			ExpectedOutput: gen.ExpectedOutput,
			Language:       gen.Language,
		}
	} else {
		for _, opt := range gen.Options {
			question.Options = append(question.Options, models.QuestionOption{
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
			})
		}
	}

	if err := g.db.Create(question).Error; err != nil {
		return nil, err
	}

	for _, topicID := range topicIDs {
		if err := g.db.Create(&models.QuestionTopic{
			QuestionID: question.ID,
			TopicID:    topicID,
		}).Error; err != nil {
			return question, fmt.Errorf("error al asociar tema %s: %w", topicID, err)
		}
	}

	return question, nil
}

func parseAIResponse(text string) (*GeneratedQuestion, error) {
	var q GeneratedQuestion
	if err := json.Unmarshal([]byte(stripCodeFence(text)), &q); err != nil {
		return nil, fmt.Errorf("JSON inválido: %w\nRespuesta: %s", err, truncate(text, 300))
	}
	if err := validateGeneratedQuestion(&q); err != nil {
		return nil, err
	}
	return &q, nil
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

func safeTopicName(topics []models.Topic) string {
	if len(topics) == 0 {
		return "general"
	}
	return topics[0].Name
}

func limitStr(session *models.Session) string {
	if session.QuestionLimit == nil {
		return "∞"
	}
	return fmt.Sprintf("%d", *session.QuestionLimit)
}
