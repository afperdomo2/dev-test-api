package questions

import (
	"unicode/utf8"

	"github.com/felipe/dev-test-api/internal/models"
)

func truncate(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen]) + "..."
}

func topicsToStrings(topics []models.Topic) []string {
	if topics == nil {
		return nil
	}
	names := make([]string, len(topics))
	for i, t := range topics {
		names[i] = t.Name
	}
	return names
}

func ToQuestionResponse(q models.Question) QuestionResponse {
	resp := QuestionResponse{
		ID:          q.ID,
		UserID:      q.UserID,
		Type:        q.Type,
		Content:     q.Content,
		Explanation: q.Explanation,
		Difficulty:  q.Difficulty,
		IsPublic:    q.IsPublic,
		Source:      q.Source,
		Topics:      topicsToStrings(q.Topics),
		CreatedAt:   q.CreatedAt,
		UpdatedAt:   q.UpdatedAt,
	}

	if q.Options != nil {
		resp.Options = make([]OptionResponse, len(q.Options))
		for i, o := range q.Options {
			resp.Options[i] = OptionResponse{
				ID:        o.ID,
				Content:   o.Content,
				IsCorrect: o.IsCorrect,
			}
		}
	}

	if q.CodeChallenge != nil {
		resp.CodeChallenge = &CodeChallengeResponse{
			ID:             q.CodeChallenge.ID,
			StarterCode:    q.CodeChallenge.StarterCode,
			ExpectedOutput: q.CodeChallenge.ExpectedOutput,
			Language:       q.CodeChallenge.Language,
			TestCasesJSON:  q.CodeChallenge.TestCasesJSON,
		}
	}

	return resp
}

func ToQuestionListResponse(q models.Question) QuestionListResponse {
	resp := QuestionListResponse{
		ID:         q.ID,
		UserID:     q.UserID,
		Type:       q.Type,
		Content:    truncate(q.Content, 70),
		Difficulty: q.Difficulty,
		IsPublic:   q.IsPublic,
		Source:     q.Source,
		Topics:     topicsToStrings(q.Topics),
		CreatedAt:  q.CreatedAt,
	}

	if q.CodeChallenge != nil {
		resp.CodeChallenge = &CodeChallengeResponse{
			ID:             q.CodeChallenge.ID,
			StarterCode:    q.CodeChallenge.StarterCode,
			ExpectedOutput: q.CodeChallenge.ExpectedOutput,
			Language:       q.CodeChallenge.Language,
			TestCasesJSON:  q.CodeChallenge.TestCasesJSON,
		}
	}

	return resp
}
