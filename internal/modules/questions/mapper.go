package questions

import "github.com/felipe/dev-test-api/internal/models"

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

	if q.Topics != nil {
		resp.Topics = make([]TopicInfo, len(q.Topics))
		for i, t := range q.Topics {
			resp.Topics[i] = TopicInfo{
				ID:   t.ID,
				Slug: t.Slug,
				Name: t.Name,
			}
		}
	}

	return resp
}

func ToQuestionListResponse(q models.Question) QuestionListResponse {
	resp := QuestionListResponse{
		ID:          q.ID,
		UserID:      q.UserID,
		Type:        q.Type,
		Content:     q.Content,
		Explanation: q.Explanation,
		Difficulty:  q.Difficulty,
		IsPublic:    q.IsPublic,
		Source:      q.Source,
		CreatedAt:   q.CreatedAt,
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

	if q.Topics != nil {
		resp.Topics = make([]TopicInfo, len(q.Topics))
		for i, t := range q.Topics {
			resp.Topics[i] = TopicInfo{
				ID:   t.ID,
				Slug: t.Slug,
				Name: t.Name,
			}
		}
	}

	return resp
}
