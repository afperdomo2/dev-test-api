package tips

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

const maxTipsPerFile = 50

type importRow struct {
	RowNum  int
	Content string
	Topic   string
}

type validatedImportRow struct {
	Content   string
	TopicSlug string
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

	commaCount := strings.Count(lines[0], ",")
	semiCount := strings.Count(lines[0], ";")
	delimiter := ','
	if semiCount > commaCount {
		delimiter = ';'
	}

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

	required := []string{"content", "topic"}
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
			RowNum:  idx + 2,
			Content: get("content"),
			Topic:   get("topic"),
		})
	}

	return rows, nil
}

func validateImportRow(row importRow) (*validatedImportRow, error) {
	content := strings.TrimSpace(row.Content)
	if content == "" {
		return nil, fmt.Errorf("el contenido es obligatorio")
	}

	topicSlug := strings.Trim(strings.ToLower(strings.TrimSpace(row.Topic)), "\"'")
	if topicSlug == "" {
		return nil, fmt.Errorf("el tema es obligatorio")
	}

	return &validatedImportRow{
		Content:   content,
		TopicSlug: topicSlug,
	}, nil
}
