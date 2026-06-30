package ai

import (
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

type aiGeneratedQuestion struct {
	Type           string              `json:"type"`
	Content        string              `json:"content"`
	Explanation    string              `json:"explanation"`
	Difficulty     string              `json:"difficulty"`
	Options        []aiGeneratedOption `json:"options,omitempty"`
	StarterCode    string              `json:"starterCode,omitempty"`
	ExpectedOutput string              `json:"expectedOutput,omitempty"`
	Language       string              `json:"language,omitempty"`
}

type aiGeneratedOption struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect"`
}

type Generator struct {
	db     *gorm.DB
	client *aiClient
}

func NewGenerator(db *gorm.DB, cfg config.AIConfig) *Generator {
	return &Generator{
		db:     db,
		client: newAIClient(cfg),
	}
}

func (g *Generator) IsConfigured() bool {
	return g.client.IsConfigured()
}

func (g *Generator) GenerateQuestion(session *models.Session) error {
	if !g.IsConfigured() {
		log.Println("⚠️ AI no configurada — omitiendo generación")
		return nil
	}

	if len(session.Topics) == 0 {
		return fmt.Errorf("la sesión no tiene temas asociados")
	}

	start := time.Now()

	existingContent, err := g.existingContent(session)
	if err != nil {
		log.Printf("⚠️ No se pudieron obtener preguntas existentes: %v", err)
	}

	userPrompt := buildUserPrompt(session, existingContent)
	log.Printf("🤖 Generando pregunta para sesión %s (tema: %s, dificultad: %s)", session.ID, safeTopicName(session.Topics), session.Difficulty)

	responseText, err := g.client.Chat(systemPrompt, userPrompt)
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

func (g *Generator) saveQuestion(session *models.Session, gen *aiGeneratedQuestion) (*models.Question, error) {
	source := "ai_generated"
	difficulty := gen.Difficulty
	if difficulty == "" {
		difficulty = session.Difficulty
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

func parseAIResponse(text string) (*aiGeneratedQuestion, error) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}

	var q aiGeneratedQuestion
	if err := json.Unmarshal([]byte(text), &q); err != nil {
		return nil, fmt.Errorf("JSON inválido: %w\nRespuesta: %s", err, truncate(text, 300))
	}

	if q.Type == "" {
		return nil, fmt.Errorf("tipo de pregunta vacío")
	}
	if q.Content == "" {
		return nil, fmt.Errorf("contenido de pregunta vacío")
	}

	validTypes := map[string]bool{"single_choice": true, "multiple_choice": true, "code_completion": true}
	if !validTypes[q.Type] {
		return nil, fmt.Errorf("tipo de pregunta no válido: %s", q.Type)
	}

	if q.Type == "code_completion" && q.Language == "" {
		q.Language = "plaintext"
	}

	if (q.Type == "single_choice" || q.Type == "multiple_choice") && len(q.Options) < 2 {
		return nil, fmt.Errorf("se requieren al menos 2 opciones para preguntas de tipo choice")
	}

	return &q, nil
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
