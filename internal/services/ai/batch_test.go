package ai

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validBatchJSON = `[
  {
    "type": "single_choice",
    "content": "¿Cómo se declara una variable en Go?",
    "explanation": "El operador := declara e inicializa.",
    "difficulty": "beginner",
    "options": [
      {"content": "var x = 5", "isCorrect": false},
      {"content": "x := 5", "isCorrect": true},
      {"content": "x = 5", "isCorrect": false},
      {"content": "let x = 5", "isCorrect": false}
    ]
  },
  {
    "type": "single_choice",
    "content": "¿Qué devuelve una función sin return explícito en Go?",
    "explanation": "Devuelve el valor cero del tipo de retorno.",
    "difficulty": "beginner",
    "options": [
      {"content": "nil", "isCorrect": false},
      {"content": "El valor cero del tipo", "isCorrect": true},
      {"content": "error", "isCorrect": false},
      {"content": "panic", "isCorrect": false}
    ]
  }
]`

func TestParseAIBatchResponseReturnsAllValidQuestions(t *testing.T) {
	questions, err := parseAIBatchResponse(validBatchJSON, "beginner")
	require.NoError(t, err)
	require.Len(t, questions, 2)
	assert.Equal(t, "¿Cómo se declara una variable en Go?", questions[0].Content)
	assert.True(t, questions[0].Options[1].IsCorrect)
}

func TestParseAIBatchResponseStripsCodeFence(t *testing.T) {
	fenced := "```json\n" + validBatchJSON + "\n```"

	questions, err := parseAIBatchResponse(fenced, "beginner")
	require.NoError(t, err)
	assert.Len(t, questions, 2)
}

func TestParseAIBatchResponseSkipsInvalidItems(t *testing.T) {
	// La segunda pregunta tiene dos opciones correctas: es inválida pero no debe
	// tumbar el lote entero.
	mixed := `[
      {
        "type": "single_choice",
        "content": "Pregunta válida",
        "explanation": "Explicación",
        "options": [
          {"content": "a", "isCorrect": true},
          {"content": "b", "isCorrect": false}
        ]
      },
      {
        "type": "single_choice",
        "content": "Pregunta inválida",
        "options": [
          {"content": "a", "isCorrect": true},
          {"content": "b", "isCorrect": true}
        ]
      }
    ]`

	questions, err := parseAIBatchResponse(mixed, "beginner")
	require.NoError(t, err)
	require.Len(t, questions, 1, "solo sobrevive la pregunta válida")
	assert.Equal(t, "Pregunta válida", questions[0].Content)
}

func TestParseAIBatchResponseFailsWhenNothingIsValid(t *testing.T) {
	allInvalid := `[
      {"type": "inventado", "content": "x", "options": []},
      {"type": "single_choice", "content": "", "options": []}
    ]`

	_, err := parseAIBatchResponse(allInvalid, "beginner")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ninguna")
}

func TestParseAIBatchResponseRejectsNonArray(t *testing.T) {
	_, err := parseAIBatchResponse(`{"type": "single_choice"}`, "beginner")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JSON inválido")
}

func TestParseAIBatchResponseRejectsEmptyArray(t *testing.T) {
	_, err := parseAIBatchResponse("[]", "beginner")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "array vacío")
}

func TestParseAIBatchResponseFillsMissingDifficulty(t *testing.T) {
	noDifficulty := `[
      {
        "type": "single_choice",
        "content": "¿Qué es un slice?",
        "explanation": "Una estructura de tamaño variable",
        "options": [
          {"content": "a", "isCorrect": true},
          {"content": "b", "isCorrect": false}
        ]
      }
    ]`

	questions, err := parseAIBatchResponse(noDifficulty, "advanced")
	require.NoError(t, err)
	require.Len(t, questions, 1)
	assert.Equal(t, "advanced", questions[0].Difficulty, "usa la dificultad solicitada")
}

func TestParseAIBatchResponseMapsRandomToIntermediate(t *testing.T) {
	noDifficulty := `[
      {
        "type": "single_choice",
        "content": "¿Qué es un slice?",
        "options": [
          {"content": "a", "isCorrect": true},
          {"content": "b", "isCorrect": false}
        ]
      }
    ]`

	questions, err := parseAIBatchResponse(noDifficulty, "random")
	require.NoError(t, err)
	assert.Equal(t, "intermediate", questions[0].Difficulty)
}

func TestBuildBatchUserPromptAsksForExactCount(t *testing.T) {
	prompt := buildBatchUserPrompt(BatchRequest{
		TopicNames: []string{"Go", "Gin"},
		Difficulty: "beginner",
		Count:      7,
	})

	assert.Contains(t, prompt, "EXACTAMENTE 7")
	assert.Contains(t, prompt, "Go, Gin")
	assert.Contains(t, prompt, "beginner")
	assert.Contains(t, prompt, "no repitas ni reformules")
}

func TestBuildBatchUserPromptExpandsRandomDifficulty(t *testing.T) {
	prompt := buildBatchUserPrompt(BatchRequest{
		TopicNames: []string{"Go"},
		Difficulty: "random",
		Count:      3,
	})

	assert.Contains(t, prompt, "variada")
	assert.NotContains(t, prompt, `"random"`)
}

func TestBuildBatchUserPromptCapsExistingContent(t *testing.T) {
	existing := make([]string, 0, maxExistingInPrompt+40)
	for range maxExistingInPrompt + 40 {
		existing = append(existing, "pregunta previa con contenido suficientemente largo")
	}

	prompt := buildBatchUserPrompt(BatchRequest{
		TopicNames:      []string{"Go"},
		Difficulty:      "beginner",
		Count:           3,
		ExistingContent: existing,
	})

	assert.Equal(t, 1, strings.Count(prompt, "Ya existen"))
	assert.Equal(t, maxExistingInPrompt, strings.Count(prompt, "pregunta previa"))
}

func TestBuildBatchUserPromptOmitsHistoryWhenEmpty(t *testing.T) {
	prompt := buildBatchUserPrompt(BatchRequest{
		TopicNames: []string{"Go"},
		Difficulty: "beginner",
		Count:      1,
	})

	assert.NotContains(t, prompt, "Ya existen")
}

func TestBatchSystemPromptDemandsArrayAndSpanish(t *testing.T) {
	assert.Contains(t, batchSystemPrompt, "ARRAY JSON")
	assert.Contains(t, batchSystemPrompt, "ESPAÑOL")
	assert.Contains(t, batchSystemPrompt, "single_choice")
	assert.Contains(t, batchSystemPrompt, "EXACTAMENTE")
}

func TestValidateGeneratedQuestionRejectsBadInput(t *testing.T) {
	tests := []struct {
		name string
		q    GeneratedQuestion
	}{
		{"sin tipo", GeneratedQuestion{Content: "x"}},
		{"sin contenido", GeneratedQuestion{Type: "single_choice"}},
		{"tipo desconocido", GeneratedQuestion{Type: "inventado", Content: "x"}},
		{"sin opciones", GeneratedQuestion{Type: "single_choice", Content: "x"}},
		{
			"sin correcta",
			GeneratedQuestion{Type: "single_choice", Content: "x", Options: []GeneratedOption{
				{Content: "a"}, {Content: "b"},
			}},
		},
		{
			"dos correctas",
			GeneratedQuestion{Type: "single_choice", Content: "x", Options: []GeneratedOption{
				{Content: "a", IsCorrect: true}, {Content: "b", IsCorrect: true},
			}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Error(t, validateGeneratedQuestion(&tc.q))
		})
	}
}

func TestValidateGeneratedQuestionDefaultsCodeLanguage(t *testing.T) {
	q := GeneratedQuestion{Type: "code_completion", Content: "Completa", StarterCode: "x"}

	require.NoError(t, validateGeneratedQuestion(&q))
	assert.Equal(t, "plaintext", q.Language)
}

func TestBatchRequestJSONRoundTrip(t *testing.T) {
	// El shape del batch es parte del contrato con la IA: si cambia, el prompt
	// deja de pedir un array. Este test avisa de cualquier cambio accidental.
	raw, err := json.Marshal(BatchRequest{TopicNames: []string{"Go"}, Difficulty: "beginner", Count: 3})
	require.NoError(t, err)

	var back map[string]any
	require.NoError(t, json.Unmarshal(raw, &back))
	assert.Equal(t, []any{"Go"}, back["TopicNames"])
	assert.Equal(t, "beginner", back["Difficulty"])
	assert.Equal(t, float64(3), back["Count"])
}
