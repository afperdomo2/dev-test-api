package questions

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func buildCSV(rows int, delimiter string) string {
	header := strings.Join([]string{"type", "content", "difficulty", "topics", "explanation", "options"}, delimiter)
	var sb strings.Builder
	sb.WriteString(header + "\n")
	for i := 1; i <= rows; i++ {
		content := fmt.Sprintf("Pregunta %d", i)
		opts := "[v] correcta | [ ] incorrecta"
		line := strings.Join([]string{"single_choice", content, "beginner", "go", "explicacion", opts}, delimiter)
		if delimiter == "," {
			line = fmt.Sprintf(`single_choice,"%s",beginner,go,"explicacion","%s"`, content, opts)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

func TestParseCSV_HeaderNotCounted(t *testing.T) {
	csv := buildCSV(2, ",")
	rows, err := parseCSV(strings.NewReader(csv))
	require.NoError(t, err)
	assert.Len(t, rows, 2, "header no debe contarse")
	assert.Equal(t, 2, rows[0].RowNum)
	assert.Equal(t, 3, rows[1].RowNum)
}

func TestParseCSV_SemicolonDelimiter(t *testing.T) {
	csv := buildCSV(1, ";")
	rows, err := parseCSV(strings.NewReader(csv))
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "single_choice", rows[0].Type)
}

func TestParseCSV_EmptyInput(t *testing.T) {
	_, err := parseCSV(strings.NewReader("   "))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vacío")
}

func TestParseCSV_MissingColumn(t *testing.T) {
	csv := "type,content\nsingle_choice,pregunta\n"
	_, err := parseCSV(strings.NewReader(csv))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "falta la columna")
}

func TestParseCSV_BlankLinesSkipped(t *testing.T) {
	header := "type,content,difficulty,topics,explanation,options"
	csv := header + "\n" +
		"single_choice,\"p1\",beginner,go,\"e\",\"[v] a | [ ] b\"\n" +
		"\n" +
		"single_choice,\"p2\",beginner,go,\"e\",\"[v] a | [ ] b\"\n"
	rows, err := parseCSV(strings.NewReader(csv))
	require.NoError(t, err)
	assert.Len(t, rows, 2)
}

func TestParseCSV_QuotedCommas(t *testing.T) {
	csv := "type,content,difficulty,topics,explanation,options\n" +
		"single_choice,\"Pregunta, con coma\",beginner,go,\"exp\",\"[v] a | [ ] b\"\n"
	rows, err := parseCSV(strings.NewReader(csv))
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "Pregunta, con coma", rows[0].Content)
}

func TestImport_FileLimitPartial(t *testing.T) {
	store := newQuestionMock()
	topicStore := newTopicMock()
	topicID := uuid.New()
	topicStore.findBySlugAndUserFn = func(slug string, createdBy *uuid.UUID) (*models.Topic, error) {
		return &models.Topic{ID: topicID, Slug: slug}, nil
	}
	topicStore.createFn = func(*models.Topic) error { return nil }

	userStore := newUserMock(200)
	svc := NewService(store, topicStore, userStore)

	store.bulkCreateFn = func(qs []*models.Question) error {
		assert.Len(t, qs, 50, "debe importar solo 50")
		return nil
	}

	csv := buildCSV(60, ",")
	result, err := svc.Import(uuid.New(), strings.NewReader(csv))
	require.NoError(t, err)
	assert.Equal(t, 60, result.Total)
	assert.Equal(t, 50, result.Imported)
	assert.Equal(t, 10, result.Failed)
	assert.Len(t, result.Errors, 10)
	for _, e := range result.Errors {
		assert.Contains(t, e.Reason, "Supera el límite de 50")
	}
}

func TestImport_DailyLimitPartial(t *testing.T) {
	store := newQuestionMock()
	topicStore := newTopicMock()
	topicID := uuid.New()
	topicStore.findBySlugAndUserFn = func(string, *uuid.UUID) (*models.Topic, error) {
		return &models.Topic{ID: topicID}, nil
	}
	userStore := &mockUserStore{
		findByIDFn: func(uuid.UUID) (*models.User, error) {
			return &models.User{DailyImportLimit: 200}, nil
		},
	}
	store.countImportedSinceFn = func(uuid.UUID, time.Time) (int64, error) {
		return 190, nil // remaining 10
	}
	svc := NewService(store, topicStore, userStore)

	store.bulkCreateFn = func(qs []*models.Question) error {
		assert.Len(t, qs, 10)
		return nil
	}

	csv := buildCSV(20, ",")
	result, err := svc.Import(uuid.New(), strings.NewReader(csv))
	require.NoError(t, err)
	assert.Equal(t, 20, result.Total)
	assert.Equal(t, 10, result.Imported)
	assert.Equal(t, 10, result.Failed)
}

func TestImport_DailyExhausted(t *testing.T) {
	store := newQuestionMock()
	topicStore := newTopicMock()
	userStore := &mockUserStore{
		findByIDFn: func(uuid.UUID) (*models.User, error) {
			return &models.User{DailyImportLimit: 200}, nil
		},
	}
	store.countImportedSinceFn = func(uuid.UUID, time.Time) (int64, error) {
		return 200, nil
	}
	svc := NewService(store, topicStore, userStore)
	csv := buildCSV(1, ",")
	_, err := svc.Import(uuid.New(), strings.NewReader(csv))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Has alcanzado")
}

func TestImport_FileAndDailyCombined(t *testing.T) {
	store := newQuestionMock()
	topicStore := newTopicMock()
	topicID := uuid.New()
	topicStore.findBySlugAndUserFn = func(string, *uuid.UUID) (*models.Topic, error) {
		return &models.Topic{ID: topicID}, nil
	}
	userStore := &mockUserStore{
		findByIDFn: func(uuid.UUID) (*models.User, error) {
			return &models.User{DailyImportLimit: 200}, nil
		},
	}
	store.countImportedSinceFn = func(uuid.UUID, time.Time) (int64, error) {
		return 180, nil // remaining 20
	}
	svc := NewService(store, topicStore, userStore)
	store.bulkCreateFn = func(qs []*models.Question) error {
		assert.Len(t, qs, 20)
		return nil
	}
	csv := buildCSV(100, ",")
	result, err := svc.Import(uuid.New(), strings.NewReader(csv))
	require.NoError(t, err)
	assert.Equal(t, 100, result.Total)
	assert.Equal(t, 20, result.Imported)
	assert.Equal(t, 80, result.Failed)
}

func TestImport_ValidationErrorsAlongsideOverflow(t *testing.T) {
	store := newQuestionMock()
	topicStore := newTopicMock()
	topicID := uuid.New()
	topicStore.findBySlugAndUserFn = func(string, *uuid.UUID) (*models.Topic, error) {
		return &models.Topic{ID: topicID}, nil
	}
	userStore := newUserMock(200)
	svc := NewService(store, topicStore, userStore)
	store.bulkCreateFn = func(qs []*models.Question) error {
		// 50 first rows: first one is invalid, so 49 imported
		assert.Len(t, qs, 49)
		return nil
	}
	header := "type,content,difficulty,topics,explanation,options"
	var sb strings.Builder
	sb.WriteString(header + "\n")
	// invalid row (no options)
	sb.WriteString("single_choice,\"pregunta invalida\",beginner,go,\"e\",\"\"\n")
	for i := 2; i <= 60; i++ {
		sb.WriteString(fmt.Sprintf("single_choice,\"p%d\",beginner,go,\"e\",\"[v] a | [ ] b\"\n", i))
	}
	result, err := svc.Import(uuid.New(), strings.NewReader(sb.String()))
	require.NoError(t, err)
	assert.Equal(t, 60, result.Total)
	// 10 overflow + 1 validation = 11 failed
	assert.Equal(t, 11, result.Failed)
	assert.Equal(t, 49, result.Imported)
	_ = gorm.ErrRecordNotFound // ensure import used
}
