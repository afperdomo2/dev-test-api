package gamification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXpHelpers(t *testing.T) {
	assert.Equal(t, 5, XpForDifficulty("beginner"))
	assert.Equal(t, 10, XpForDifficulty("intermediate"))
	assert.Equal(t, 15, XpForDifficulty("advanced"))
	assert.Equal(t, 0, XpForDifficulty("unknown"))
	assert.Equal(t, 20, SessionBonusXp())
}
