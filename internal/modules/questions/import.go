package questions

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

const (
	maxQuestionsPerFile     = 50
	defaultImportDifficulty = "intermediate"
)

var validImportTypes = map[string]bool{
	"single_choice":   true,
	"multiple_choice": true,
}

var validImportDifficulties = map[string]bool{
	"beginner":     true,
	"intermediate": true,
	"advanced":     true,
}

type importRow struct {
	RowNum      int
	Type        string
	Content     string
	Difficulty  string
	TopicsRaw   string
	Explanation string
	OptionsRaw  string
}

type parsedOption struct {
	Content   string
	IsCorrect bool
}

func detectDelimiter(header string) rune {
	commaCount := strings.Count(header, ",")
	semiCount := strings.Count(header, ";")
	if semiCount > commaCount {
		return ';'
	}
	return ','
}

func parseCSV(r io.Reader) ([]importRow, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo: %w", err)
	}
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return nil, fmt.Errorf("el archivo está vacío")
	}

	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("el archivo está vacío")
	}
	delimiter := detectDelimiter(lines[0])

	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al parsear el CSV: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("el archivo está vacío")
	}

	header := records[0]
	headerMap := map[string]int{}
	for i, h := range header {
		key := strings.ToLower(strings.TrimSpace(h))
		key = strings.Trim(key, "\"'")
		headerMap[key] = i
	}

	required := []string{"type", "content", "difficulty", "topics", "explanation", "options"}
	for _, col := range required {
		if _, ok := headerMap[col]; !ok {
			return nil, fmt.Errorf("falta la columna requerida: %s", col)
		}
	}

	var rows []importRow
	for idx, rec := range records[1:] {
		if len(rec) == 0 {
			continue
		}
		allEmpty := true
		for _, f := range rec {
			if strings.TrimSpace(f) != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			continue
		}

		get := func(col string) string {
			pos := headerMap[col]
			if pos < len(rec) {
				return strings.TrimSpace(rec[pos])
			}
			return ""
		}

		rows = append(rows, importRow{
			RowNum:      idx + 2,
			Type:        get("type"),
			Content:     get("content"),
			Difficulty:  get("difficulty"),
			TopicsRaw:   get("topics"),
			Explanation: get("explanation"),
			OptionsRaw:  get("options"),
		})
	}

	return rows, nil
}

func parseOptions(raw string) ([]parsedOption, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("se requieren opciones")
	}

	parts := strings.Split(raw, "|")
	var opts []parsedOption
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		isCorrect := false
		content := p
		lower := strings.ToLower(strings.TrimSpace(p))
		if strings.HasPrefix(lower, "[v]") || strings.HasPrefix(lower, "[x]") {
			isCorrect = true
			content = strings.TrimSpace(p[3:])
		} else if strings.HasPrefix(lower, "[ ]") || strings.HasPrefix(lower, "[f]") {
			isCorrect = false
			content = strings.TrimSpace(p[3:])
		}
		if strings.TrimSpace(content) == "" {
			continue
		}
		opts = append(opts, parsedOption{Content: content, IsCorrect: isCorrect})
	}

	if len(opts) < 2 {
		return nil, fmt.Errorf("se requieren al menos 2 opciones")
	}

	hasCorrect := false
	for _, o := range opts {
		if o.IsCorrect {
			hasCorrect = true
			break
		}
	}
	if !hasCorrect {
		return nil, fmt.Errorf("debe marcar al menos una opción como correcta ([v])")
	}

	return opts, nil
}

type validatedImportRow struct {
	Type       string
	Difficulty string
	Opts       []parsedOption
	TopicSlugs []string
}

func validateImportRow(row importRow) (*validatedImportRow, error) {
	typeNorm := strings.ToLower(strings.TrimSpace(row.Type))
	if !validImportTypes[typeNorm] {
		return nil, fmt.Errorf("tipo no válido: %s (usa single_choice o multiple_choice)", row.Type)
	}

	if strings.TrimSpace(row.Content) == "" {
		return nil, fmt.Errorf("el contenido es obligatorio")
	}

	diff := strings.ToLower(strings.TrimSpace(row.Difficulty))
	if diff == "" {
		diff = defaultImportDifficulty
	}
	if !validImportDifficulties[diff] {
		return nil, fmt.Errorf("dificultad no válida: %s", row.Difficulty)
	}

	opts, err := parseOptions(row.OptionsRaw)
	if err != nil {
		return nil, err
	}

	correctCount := 0
	for _, o := range opts {
		if o.IsCorrect {
			correctCount++
		}
	}
	if typeNorm == "single_choice" && correctCount != 1 {
		return nil, fmt.Errorf("single_choice debe tener exactamente 1 opción correcta, tiene %d", correctCount)
	}

	topicSlugs := parseTopicSlugs(row.TopicsRaw)
	if len(topicSlugs) == 0 {
		return nil, fmt.Errorf("se requiere al menos un tema")
	}

	return &validatedImportRow{
		Type:       typeNorm,
		Difficulty: diff,
		Opts:       opts,
		TopicSlugs: topicSlugs,
	}, nil
}

func parseTopicSlugs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ";")
	var slugs []string
	seen := map[string]bool{}
	for _, p := range parts {
		s := strings.ToLower(strings.TrimSpace(p))
		s = strings.Trim(s, "\"'")
		if s == "" {
			continue
		}
		if !seen[s] {
			seen[s] = true
			slugs = append(slugs, s)
		}
	}
	return slugs
}
