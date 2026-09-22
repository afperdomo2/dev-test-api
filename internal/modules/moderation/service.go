package moderation

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type questionStore interface {
	FindByID(id uuid.UUID) (*models.Question, error)
	Unpublish(id uuid.UUID) error
}

type tipStore interface {
	FindByID(id uuid.UUID) (*models.Tip, error)
}

type Service interface {
	Create(userID uuid.UUID, req CreateReportRequest) (*ReportResponse, error)
	List(params ListReportsParams) ([]ReportResponse, int64, error)
	Review(id uuid.UUID, req ReviewReportRequest) (*ReportResponse, error)
}

type moderationService struct {
	store         Store
	questionStore questionStore
	tipStore      tipStore
}

func NewService(store Store, questionStore questionStore, tipStore tipStore) Service {
	return &moderationService{store: store, questionStore: questionStore, tipStore: tipStore}
}

func (s *moderationService) Create(userID uuid.UUID, req CreateReportRequest) (*ReportResponse, error) {
	switch req.TargetType {
	case "question":
		q, err := s.questionStore.FindByID(req.TargetID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, apierr.ErrValidation("El contenido a reportar no existe", "")
			}
			return nil, apierr.ErrInternal("Error al validar el contenido", "")
		}
		if !q.IsPublic && q.Source != "ai_generated" {
			return nil, apierr.ErrValidation("Solo se puede reportar contenido público", "")
		}
	case "tip":
		if _, err := s.tipStore.FindByID(req.TargetID); err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, apierr.ErrValidation("El contenido a reportar no existe", "")
			}
			return nil, apierr.ErrInternal("Error al validar el contenido", "")
		}
	default:
		return nil, apierr.ErrValidation("El tipo de contenido no es válido", "")
	}

	exists, err := s.store.HasPending(userID, req.TargetType, req.TargetID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al verificar reportes existentes", "")
	}
	if exists {
		return nil, apierr.ErrConflict("Reporte duplicado", "Ya tienes un reporte pendiente para este contenido", "")
	}

	report := &models.Report{
		ReporterUserID: userID,
		TargetType:     req.TargetType,
		TargetID:       req.TargetID,
		Reason:         req.Reason,
		Status:         "pending",
	}

	if err := s.store.Create(report); err != nil {
		return nil, apierr.ErrInternal("Error al crear el reporte", "")
	}

	resp := ToReportResponse(*report)
	return &resp, nil
}

func (s *moderationService) List(params ListReportsParams) ([]ReportResponse, int64, error) {
	if params.Status != "" && params.Status != "pending" && params.Status != "reviewed" {
		return nil, 0, apierr.ErrValidation("El estado del reporte no es válido", "")
	}

	reports, total, err := s.store.FindPage(params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar los reportes", "")
	}

	result := make([]ReportResponse, len(reports))
	for i, r := range reports {
		result[i] = ToReportResponse(r)
	}
	return result, total, nil
}

func (s *moderationService) Review(id uuid.UUID, req ReviewReportRequest) (*ReportResponse, error) {
	report, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Reporte", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el reporte", "")
	}

	if report.Status == "reviewed" {
		return nil, apierr.ErrConflict("Reporte ya revisado", "El reporte ya fue revisado", "")
	}

	if req.Unpublish != nil && *req.Unpublish && report.TargetType == "question" {
		if err := s.questionStore.Unpublish(report.TargetID); err != nil {
			return nil, apierr.ErrInternal("Error al des-publicar el contenido", "")
		}
	}

	report.Status = "reviewed"
	if err := s.store.Update(report); err != nil {
		return nil, apierr.ErrInternal("Error al actualizar el reporte", "")
	}

	resp := ToReportResponse(*report)
	return &resp, nil
}
