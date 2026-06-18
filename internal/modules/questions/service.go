package questions

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(params common.PaginationParams, filters QuestionFilters) ([]QuestionListResponse, int64, error)
	GetByID(id uuid.UUID) (*QuestionResponse, error)
	Create(userID uuid.UUID, input CreateQuestionRequest) (*QuestionResponse, error)
	Update(id uuid.UUID, input UpdateQuestionRequest) (*QuestionResponse, error)
	Delete(id uuid.UUID) error
}

type questionService struct {
	store Store
}

func NewService(store Store) Service {
	return &questionService{store: store}
}

func (s *questionService) List(params common.PaginationParams, filters QuestionFilters) ([]QuestionListResponse, int64, error) {
	questions, total, err := s.store.FindPage(params, filters)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las preguntas", "")
	}

	result := make([]QuestionListResponse, len(questions))
	for i, q := range questions {
		result[i] = ToQuestionListResponse(q)
	}
	return result, total, nil
}

func (s *questionService) GetByID(id uuid.UUID) (*QuestionResponse, error) {
	question, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Pregunta", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la pregunta", "")
	}

	resp := ToQuestionResponse(*question)
	return &resp, nil
}

func (s *questionService) Create(userID uuid.UUID, input CreateQuestionRequest) (*QuestionResponse, error) {
	if input.Type == "single_choice" || input.Type == "multiple_choice" {
		if len(input.Options) == 0 {
			return nil, apierr.ErrValidation("Se requieren opciones para preguntas de tipo choice", "")
		}
	}
	if input.Type == "code_completion" && input.Language == "" {
		return nil, apierr.ErrValidation("Se requiere el lenguaje para preguntas code_completion", "")
	}

	question := &models.Question{
		UserID:      userID,
		Type:        input.Type,
		Content:     input.Content,
		Explanation: input.Explanation,
		Difficulty:  input.Difficulty,
		Source:      input.Source,
	}

	if question.Source == "" {
		question.Source = "manual"
	}

	if input.Type == "code_completion" {
		question.CodeChallenge = &models.CodeChallenge{
			StarterCode:    input.StarterCode,
			ExpectedOutput: input.ExpectedOutput,
			Language:       input.Language,
			TestCasesJSON:  input.TestCasesJSON,
		}
	} else {
		for _, opt := range input.Options {
			question.Options = append(question.Options, models.QuestionOption{
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
			})
		}
	}

	if err := s.store.Create(question); err != nil {
		return nil, apierr.ErrInternal("Error al crear la pregunta", "")
	}

	if err := s.store.AddQuestionTopics(question.ID, input.TopicIDs); err != nil {
		return nil, apierr.ErrInternal("Error al asociar los temas", "")
	}

	created, err := s.store.FindByID(question.ID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la pregunta creada", "")
	}
	question = created

	resp := ToQuestionResponse(*question)
	return &resp, nil
}

func (s *questionService) Update(id uuid.UUID, input UpdateQuestionRequest) (*QuestionResponse, error) {
	question, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Pregunta", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la pregunta", "")
	}

	if input.Content != "" {
		question.Content = input.Content
	}
	if input.Explanation != "" {
		question.Explanation = input.Explanation
	}
	if input.Difficulty != "" {
		question.Difficulty = input.Difficulty
	}

	if input.TopicIDs != nil {
		if err := s.store.ReplaceQuestionTopics(question.ID, input.TopicIDs); err != nil {
			return nil, apierr.ErrInternal("Error al actualizar los temas", "")
		}
		question.Topics = nil
	}

	if input.Options != nil {
		options := make([]models.QuestionOption, len(input.Options))
		for i, opt := range input.Options {
			options[i] = models.QuestionOption{
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
			}
		}
		if err := s.store.ReplaceQuestionOptions(question.ID, options); err != nil {
			return nil, apierr.ErrInternal("Error al actualizar las opciones", "")
		}
		question.Options = nil
	}

	if input.StarterCode != "" || input.Language != "" {
		if question.CodeChallenge == nil {
			question.CodeChallenge = &models.CodeChallenge{
				QuestionID: question.ID,
			}
		}
		if input.StarterCode != "" {
			question.CodeChallenge.StarterCode = input.StarterCode
		}
		if input.ExpectedOutput != "" {
			question.CodeChallenge.ExpectedOutput = input.ExpectedOutput
		}
		if input.Language != "" {
			question.CodeChallenge.Language = input.Language
		}
		if input.TestCasesJSON != "" {
			question.CodeChallenge.TestCasesJSON = input.TestCasesJSON
		}
	}

	if err := s.store.Update(question); err != nil {
		return nil, apierr.ErrInternal("Error al actualizar la pregunta", "")
	}

	question, err = s.store.FindByID(question.ID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la pregunta actualizada", "")
	}

	resp := ToQuestionResponse(*question)
	return &resp, nil
}

func (s *questionService) Delete(id uuid.UUID) error {
	_, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierr.ErrNotFound("Pregunta", "")
		}
		return apierr.ErrInternal("Error al obtener la pregunta", "")
	}

	if err := s.store.Delete(id); err != nil {
		return apierr.ErrInternal("Error al eliminar la pregunta", "")
	}

	return nil
}
