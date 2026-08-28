package ai

import (
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildUserPrompt(t *testing.T) {
	sess := &models.Session{
		Difficulty: "beginner",
		Topics:     []models.Topic{{ID: uuid.New(), Name: "Go", Category: "lang"}, {ID: uuid.New(), Name: "Gin", Category: "framework"}},
	}
	prompt := buildUserPrompt(sess, nil)
	assert.Contains(t, prompt, "beginner")
	assert.Contains(t, prompt, "Go")
	assert.Contains(t, prompt, "Gin")

	withExisting := buildUserPrompt(sess, []string{"¿Qué es Go?", "Otra pregunta"})
	assert.Contains(t, withExisting, "Ya existen")
	assert.Contains(t, withExisting, "¿Qué es Go?")

	emptyTopics := buildUserPrompt(&models.Session{Difficulty: "advanced", Topics: nil}, nil)
	assert.Contains(t, emptyTopics, "advanced")
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "hi", truncate("hi", 10))
	assert.Equal(t, "hello...", truncate("hello world", 5))
	assert.Equal(t, "", truncate("", 5))
}
