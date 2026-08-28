package ai

import (
	"testing"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIsConfigured(t *testing.T) {
	assert.False(t, newAIClient(config.AIConfig{}).IsConfigured())
	assert.False(t, newAIClient(config.AIConfig{APIURL: "http://x"}).IsConfigured())
	assert.False(t, newAIClient(config.AIConfig{APIKey: "k"}).IsConfigured())
	assert.True(t, newAIClient(config.AIConfig{APIURL: "http://x", APIKey: "k"}).IsConfigured())
}

func TestBuildChatURL(t *testing.T) {
	tests := []struct{ in, want string }{
		{"https://api.com/v1", "https://api.com/v1/chat/completions"},
		{"https://api.com/v1/", "https://api.com/v1/chat/completions"},
		{"https://api.com/v1/chat/completions", "https://api.com/v1/chat/completions"},
		{"https://api.com/v1/chat/completions/", "https://api.com/v1/chat/completions"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, buildChatURL(tt.in))
	}
}
