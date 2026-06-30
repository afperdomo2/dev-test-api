package sessions

import (
	"encoding/json"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

func ToSessionResponse(s models.Session) SessionResponse {
	topics := make([]string, len(s.Topics))
	for i, t := range s.Topics {
		topics[i] = t.Name
	}
	return SessionResponse{
		ID:                 s.ID,
		UserID:             s.UserID,
		Name:               s.Name,
		Status:             s.Status,
		Mode:               s.Mode,
		Difficulty:         s.Difficulty,
		QuestionLimit:      s.QuestionLimit,
		QuestionsGenerated: s.QuestionsGenerated,
		Score:              s.Score,
		StartedAt:          s.StartedAt,
		FinishedAt:         s.FinishedAt,
		Topics:             topics,
		AnswerCount:        len(s.Answers),
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}

func ToSessionListResponse(s models.Session) SessionListResponse {
	topics := make([]string, len(s.Topics))
	for i, t := range s.Topics {
		topics[i] = t.Name
	}
	return SessionListResponse{
		ID:                 s.ID,
		UserID:             s.UserID,
		Name:               s.Name,
		Status:             s.Status,
		Mode:               s.Mode,
		Difficulty:         s.Difficulty,
		QuestionLimit:      s.QuestionLimit,
		QuestionsGenerated: s.QuestionsGenerated,
		Score:              s.Score,
		StartedAt:          s.StartedAt,
		FinishedAt:         s.FinishedAt,
		Topics:             topics,
		AnswerCount:        len(s.Answers),
		CreatedAt:          s.CreatedAt,
	}
}

func ToAnswerResponse(a models.SessionAnswer) SessionAnswerResponse {
	var selected []uuid.UUID
	if a.SelectedOptions != "" {
		json.Unmarshal([]byte(a.SelectedOptions), &selected)
	}

	resp := SessionAnswerResponse{
		ID:              a.ID,
		QuestionID:      a.QuestionID,
		AnswerText:      a.AnswerText,
		SelectedOptions: selected,
		IsCorrect:       a.IsCorrect,
		AiFeedback:      a.AiFeedback,
		ResponseTimeMs:  a.ResponseTimeMs,
		CreatedAt:       a.CreatedAt,
	}

	if a.Question != nil {
		resp.Explanation = a.Question.Explanation
	}

	return resp
}

func ToAnswerDetailResponse(a models.SessionAnswer) SessionAnswerDetailResponse {
	var selected []uuid.UUID
	if a.SelectedOptions != "" {
		json.Unmarshal([]byte(a.SelectedOptions), &selected)
	}

	resp := SessionAnswerDetailResponse{
		ID:              a.ID,
		QuestionID:      a.QuestionID,
		AnswerText:      a.AnswerText,
		SelectedOptions: selected,
		IsCorrect:       a.IsCorrect,
		AiFeedback:      a.AiFeedback,
		ResponseTimeMs:  a.ResponseTimeMs,
		CreatedAt:       a.CreatedAt,
	}

	if a.Question != nil {
		resp.Explanation = a.Question.Explanation

		q := a.Question
		resp.Question = QuestionSnapshot{
			Content:    q.Content,
			Type:       q.Type,
			Difficulty: q.Difficulty,
		}

		if q.CodeChallenge != nil {
			resp.Question.CodeChallenge = &questions.CodeChallengeResponse{
				ID:             q.CodeChallenge.ID,
				StarterCode:    q.CodeChallenge.StarterCode,
				ExpectedOutput: q.CodeChallenge.ExpectedOutput,
				Language:       q.CodeChallenge.Language,
				TestCasesJSON:  q.CodeChallenge.TestCasesJSON,
			}
		}

		if q.Topics != nil {
			resp.Question.Topics = make([]string, len(q.Topics))
			for i, t := range q.Topics {
				resp.Question.Topics[i] = t.Name
			}
		}

		if q.Options != nil {
			resp.Question.Options = make([]QuestionSnapshotOption, len(q.Options))
			for i, o := range q.Options {
				resp.Question.Options[i] = QuestionSnapshotOption{
					ID:        o.ID,
					Content:   o.Content,
					IsCorrect: o.IsCorrect,
				}
			}
		}
	}

	return resp
}

func toNextQuestionItem(q models.Question) NextQuestionItem {
	item := NextQuestionItem{
		ID:         q.ID,
		Type:       q.Type,
		Content:    q.Content,
		Difficulty: q.Difficulty,
	}

	if q.Options != nil {
		item.Options = make([]NextQuestionOption, len(q.Options))
		for i, o := range q.Options {
			item.Options[i] = NextQuestionOption{
				ID:      o.ID,
				Content: o.Content,
			}
		}
	}

	if q.CodeChallenge != nil {
		item.CodeChallenge = &questions.CodeChallengeResponse{
			ID:             q.CodeChallenge.ID,
			StarterCode:    q.CodeChallenge.StarterCode,
			ExpectedOutput: q.CodeChallenge.ExpectedOutput,
			Language:       q.CodeChallenge.Language,
			TestCasesJSON:  q.CodeChallenge.TestCasesJSON,
		}
	}

	if q.Topics != nil {
		item.Topics = make([]string, len(q.Topics))
		for i, t := range q.Topics {
			item.Topics[i] = t.Name
		}
	}

	return item
}
