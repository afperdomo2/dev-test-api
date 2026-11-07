package ai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAITipsResponseValid(t *testing.T) {
	tips, err := parseAITipsResponse("```json\n" + `[
		{"content": "Usa defer para liberar recursos.", "topic": "go"},
		{"content": "Evitá mutar slices compartidos.", "topic": "go"}
	]` + "\n```")
	require.NoError(t, err)
	require.Len(t, tips, 2)
	assert.Equal(t, "go", tips[0].Topic)
}

func TestParseAITipsResponseDiscardsInvalid(t *testing.T) {
	tips, err := parseAITipsResponse(`[
		{"content": "Tip válido", "topic": "go"},
		{"content": "   ", "topic": "go"},
		{"content": "Sin tema", "topic": ""},
		{"content": "` + strings.Repeat("x", maxTipContentLength+1) + `", "topic": "go"}
	]`)
	require.NoError(t, err)
	require.Len(t, tips, 1)
	assert.Equal(t, "Tip válido", tips[0].Content)
}

func TestParseAITipsResponseAllInvalid(t *testing.T) {
	_, err := parseAITipsResponse(`[{"content": "", "topic": ""}]`)
	require.Error(t, err)
}

func TestParseAITipsResponseEmptyArray(t *testing.T) {
	_, err := parseAITipsResponse(`[]`)
	require.Error(t, err)
}

func TestParseAITipsResponseInvalidJSON(t *testing.T) {
	_, err := parseAITipsResponse("no es json")
	require.Error(t, err)
}

func TestValidateGeneratedTipTrimsFields(t *testing.T) {
	tip := &GeneratedTip{Content: "  usa constantes  ", Topic: "  go  "}
	require.NoError(t, validateGeneratedTip(tip))
	assert.Equal(t, "usa constantes", tip.Content)
	assert.Equal(t, "go", tip.Topic)
}

func TestBuildTipsUserPromptDistributesTopicsAndExcludesExisting(t *testing.T) {
	prompt := buildTipsUserPrompt(TipBatchRequest{
		TopicNames:      []string{"Go", "Rust"},
		TopicSlugs:      []string{"go", "rust"},
		Count:           4,
		ExistingContent: []string{"Tip viejo"},
	})

	assert.Contains(t, prompt, "EXACTAMENTE 4 tips")
	assert.Contains(t, prompt, "Go (slug: go)")
	assert.Contains(t, prompt, "Rust (slug: rust)")
	assert.Contains(t, prompt, "Distribuye los 4 tips")
	assert.Contains(t, prompt, "NO los repitas")
	assert.Contains(t, prompt, "Tip viejo")
}
