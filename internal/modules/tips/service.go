package tips

import (
	"fmt"
	"io"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type topicStore interface {
	FindByID(id uuid.UUID) (*models.Topic, error)
	FindBySlugAndUser(slug string, createdBy *uuid.UUID) (*models.Topic, error)
}

type Service interface {
	List(params ListTipsParams) ([]TipResponse, int64, error)
	GetByID(id uuid.UUID) (*TipResponse, error)
	GetTipOfTheDay(userID uuid.UUID) (*TipResponse, error)
	Create(userID uuid.UUID, input CreateTipRequest) (*TipResponse, error)
	Update(id uuid.UUID, input UpdateTipRequest) (*TipResponse, error)
	Delete(id uuid.UUID) error
	Import(userID uuid.UUID, r io.Reader) (*ImportResult, error)
}

type tipService struct {
	store      Store
	topicStore topicStore
}

func NewService(store Store, topicStore topicStore) Service {
	return &tipService{store: store, topicStore: topicStore}
}

func (s *tipService) List(params ListTipsParams) ([]TipResponse, int64, error) {
	tips, total, err := s.store.FindPage(params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar los tips", "")
	}

	result := make([]TipResponse, len(tips))
	for i, t := range tips {
		result[i] = ToTipResponse(t)
	}
	return result, total, nil
}

func (s *tipService) GetByID(id uuid.UUID) (*TipResponse, error) {
	tip, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Tip", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el tip", "")
	}

	resp := ToTipResponse(*tip)
	return &resp, nil
}

func (s *tipService) GetTipOfTheDay(userID uuid.UUID) (*TipResponse, error) {
	tip, err := s.store.FindRandom()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Tip", "No hay tips disponibles")
		}
		return nil, apierr.ErrInternal("Error al obtener el tip del día", "")
	}

	resp := ToTipResponse(*tip)
	return &resp, nil
}

func (s *tipService) Create(userID uuid.UUID, input CreateTipRequest) (*TipResponse, error) {
	if _, err := s.topicStore.FindByID(input.TopicID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("El tema seleccionado no existe", "")
		}
		return nil, apierr.ErrInternal("Error al validar el tema", "")
	}

	tip := &models.Tip{
		UserID:  userID,
		TopicID: input.TopicID,
		Content: input.Content,
		Source:  "manual",
	}

	if err := s.store.Create(tip); err != nil {
		return nil, apierr.ErrInternal("Error al crear el tip", "")
	}

	return s.GetByID(tip.ID)
}

func (s *tipService) Update(id uuid.UUID, input UpdateTipRequest) (*TipResponse, error) {
	tip, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Tip", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el tip", "")
	}

	if input.Content != "" {
		tip.Content = input.Content
	}
	if input.TopicID != uuid.Nil {
		if _, tErr := s.topicStore.FindByID(input.TopicID); tErr != nil {
			if tErr == gorm.ErrRecordNotFound {
				return nil, apierr.ErrValidation("El tema seleccionado no existe", "")
			}
			return nil, apierr.ErrInternal("Error al validar el tema", "")
		}
		tip.TopicID = input.TopicID
	}

	if err := s.store.Update(tip); err != nil {
		return nil, apierr.ErrInternal("Error al actualizar el tip", "")
	}

	return s.GetByID(tip.ID)
}

func (s *tipService) Delete(id uuid.UUID) error {
	if _, err := s.store.FindByID(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierr.ErrNotFound("Tip", "")
		}
		return apierr.ErrInternal("Error al obtener el tip", "")
	}

	if err := s.store.Delete(id); err != nil {
		return apierr.ErrInternal("Error al eliminar el tip", "")
	}

	return nil
}

func (s *tipService) Import(userID uuid.UUID, r io.Reader) (*ImportResult, error) {
	rows, err := parseCSV(r)
	if err != nil {
		return nil, apierr.ErrValidation(err.Error(), "")
	}

	if len(rows) == 0 {
		return nil, apierr.ErrValidation("El archivo no contiene tips", "")
	}

	totalRows := len(rows)
	var overflowErrors []ImportRowError

	if len(rows) > maxTipsPerFile {
		for _, row := range rows[maxTipsPerFile:] {
			overflowErrors = append(overflowErrors, ImportRowError{
				Row:    row.RowNum,
				Reason: fmt.Sprintf("Supera el límite de %d tips por archivo", maxTipsPerFile),
			})
		}
		rows = rows[:maxTipsPerFile]
	}

	topicCache := make(map[string]uuid.UUID)
	validTips := make([]*models.Tip, 0, len(rows))
	rowErrors := overflowErrors

	for _, row := range rows {
		validated, vErr := validateImportRow(row)
		if vErr != nil {
			rowErrors = append(rowErrors, ImportRowError{Row: row.RowNum, Reason: vErr.Error()})
			continue
		}

		topicID, ok := topicCache[validated.TopicSlug]
		if !ok {
			topic, tErr := s.topicStore.FindBySlugAndUser(validated.TopicSlug, nil)
			if tErr != nil {
				rowErrors = append(rowErrors, ImportRowError{Row: row.RowNum, Reason: fmt.Sprintf("tema no encontrado: %s", validated.TopicSlug)})
				continue
			}
			topicID = topic.ID
			topicCache[validated.TopicSlug] = topicID
		}

		validTips = append(validTips, &models.Tip{
			UserID:  userID,
			TopicID: topicID,
			Content: validated.Content,
			Source:  "imported",
		})
	}

	result := &ImportResult{
		Total:  totalRows,
		Failed: len(rowErrors),
		Errors: rowErrors,
	}

	if len(validTips) > 0 {
		if err := s.store.BulkCreate(validTips); err != nil {
			return nil, apierr.ErrInternal("Error al importar los tips", "")
		}
		result.Imported = len(validTips)
	}

	if result.Errors == nil {
		result.Errors = []ImportRowError{}
	}

	return result, nil
}
