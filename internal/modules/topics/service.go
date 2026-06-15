package topics

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(page, perPage int) ([]TopicResponse, int64, error)
	GetByID(id uuid.UUID) (*TopicResponse, error)
	Create(userID uuid.UUID, input CreateTopicRequest) (*TopicResponse, error)
	Update(id uuid.UUID, input UpdateTopicRequest) (*TopicResponse, error)
	Delete(id uuid.UUID) error
}

type topicService struct {
	store Store
}

func NewService(store Store) Service {
	return &topicService{store: store}
}

func (s *topicService) List(page, perPage int) ([]TopicResponse, int64, error) {
	topics, total, err := s.store.FindPage(page, perPage)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar los temas", "")
	}

	result := make([]TopicResponse, len(topics))
	for i, t := range topics {
		result[i] = ToTopicResponse(t)
	}
	return result, total, nil
}

func (s *topicService) GetByID(id uuid.UUID) (*TopicResponse, error) {
	topic, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Tema", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el tema", "")
	}

	resp := ToTopicResponse(*topic)
	return &resp, nil
}

func (s *topicService) Create(userID uuid.UUID, input CreateTopicRequest) (*TopicResponse, error) {
	existing, _ := s.store.FindBySlugAndUser(input.Slug, nil)
	if existing != nil && existing.IsSystem {
		return nil, apierr.ErrConflict(
			"Slug Already Exists",
			"Ya existe un tema del sistema con este slug",
			"",
		)
	}

	existing, _ = s.store.FindBySlugAndUser(input.Slug, &userID)
	if existing != nil {
		return nil, apierr.ErrConflict(
			"Slug Already Exists",
			"Ya tienes un tema personalizado con este slug",
			"",
		)
	}

	topic := &models.Topic{
		Slug:      input.Slug,
		Name:      input.Name,
		Category:  input.Category,
		IsSystem:  false,
		CreatedBy: &userID,
	}

	if err := s.store.Create(topic); err != nil {
		return nil, apierr.ErrInternal("Error al crear el tema", "")
	}

	resp := ToTopicResponse(*topic)
	return &resp, nil
}

func (s *topicService) Update(id uuid.UUID, input UpdateTopicRequest) (*TopicResponse, error) {
	topic, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Tema", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el tema", "")
	}

	if input.Name != "" {
		topic.Name = input.Name
	}
	if input.Category != "" {
		topic.Category = input.Category
	}

	if err := s.store.Update(topic); err != nil {
		return nil, apierr.ErrInternal("Error al actualizar el tema", "")
	}

	resp := ToTopicResponse(*topic)
	return &resp, nil
}

func (s *topicService) Delete(id uuid.UUID) error {
	_, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierr.ErrNotFound("Tema", "")
		}
		return apierr.ErrInternal("Error al obtener el tema", "")
	}

	if err := s.store.Delete(id); err != nil {
		return apierr.ErrInternal("Error al eliminar el tema", "")
	}

	return nil
}
