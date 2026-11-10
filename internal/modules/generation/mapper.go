package generation

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
)

// toJobResponse proyecta un job. topicNames se usa cuando el store ya trajo
// los temas precargados; si viene vacío se omiten.
func toJobResponse(job models.GenerationJob, topicNames []string) JobResponse {
	topicIDs := make([]string, 0, len(job.JobTopics)+len(job.Topics))
	topics := make([]string, 0, len(topicNames))

	for _, t := range job.JobTopics {
		topicIDs = append(topicIDs, t.TopicID.String())
	}
	for _, t := range job.Topics {
		topicIDs = append(topicIDs, t.ID.String())
		if len(topicNames) == 0 {
			topics = append(topics, t.Name)
		}
	}
	if len(topicNames) > 0 {
		topics = topicNames
	}

	kind := job.Kind
	if kind == "" {
		kind = models.GenerationKindQuestion
	}

	// Los trabajos anteriores a la migración 00011 se leen con el default de la
	// columna, pero los tests y los jobs en memoria pueden traerlo vacío.
	questionType := job.QuestionType
	if questionType == "" {
		questionType = models.GenerationQuestionTypeSingleChoice
	}

	return JobResponse{
		ID:           job.ID.String(),
		RequestedBy:  job.RequestedBy.String(),
		Kind:         kind,
		TopicIDs:     topicIDs,
		Topics:       topics,
		Difficulty:   job.Difficulty,
		QuestionType: questionType,
		Quantity:     job.Quantity,
		Status:       job.Status,
		Generated:    job.Generated,
		Error:        job.Error,
		CreatedAt:    job.CreatedAt,
		StartedAt:    job.StartedAt,
		FinishedAt:   job.FinishedAt,
	}
}

func uniqueTopicIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
