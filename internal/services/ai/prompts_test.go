package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	assert.Equal(t, "hi", truncate("hi", 10))
	assert.Equal(t, "hello...", truncate("hello world", 5))
	assert.Equal(t, "", truncate("", 5))
}
