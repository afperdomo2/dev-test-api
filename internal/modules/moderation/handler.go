package moderation

import (
	"net/http"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary      Reportar contenido
// @Description  Reporta contenido público (pregunta o tip) por contenido inapropiado. Solo contenido público es reportable.
// @Tags         moderation
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateReportRequest  true  "Datos del reporte"
// @Success      201   {object}  ReportResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError
// @Failure      422   {object}  apierr.APIError
// @Router       /api/v1/reports [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	report, err := h.service.Create(userID, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, report)
}

// @Summary      Listar reportes
// @Description  Lista los reportes con paginación y filtro opcional por estado (solo administradores)
// @Tags         moderation
// @Security     BearerAuth
// @Produce      json
// @Param        page      query  int    false  "Número de página (default: 1)"
// @Param        perPage   query  int    false  "Elementos por página (default: 20, max: 100)"
// @Param        status    query  string false  "Filtrar por estado: pending, reviewed"
// @Param        sortBy    query  string false  "Campo de ordenación: created_at, updated_at"
// @Param        sortOrder query  string false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Lista de reportes (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/reports [get]
func (h *Handler) List(c *gin.Context) {
	pagination, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	params := ListReportsParams{
		PaginationParams: pagination,
		Status:           c.Query("status"),
	}

	reports, total, err := h.service.List(params)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, reports, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Revisar reporte
// @Description  Marca un reporte como revisado. Opcionalmente des-publica el contenido reportado si es pregunta de admin (`unpublish: true`) (solo administradores)
// @Tags         moderation
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                true  "Reporte ID"
// @Param        body  body  ReviewReportRequest   true  "Acción de revisión"
// @Success      200   {object}  ReportResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError
// @Router       /api/v1/reports/{id}/review [put]
func (h *Handler) Review(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Reporte", c.Request.URL.Path)
		return
	}

	var req ReviewReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	report, err := h.service.Review(id, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, report)
}

func getUserID(c *gin.Context) (uuid.UUID, *apierr.APIError) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return uuid.Nil, apierr.ErrUnauthorized("Usuario no autenticado", "")
	}

	mapClaims, ok := claims.(*jwt.MapClaims)
	if !ok {
		return uuid.Nil, apierr.ErrInternal("Error al obtener los claims del usuario", "")
	}

	sub, ok := (*mapClaims)["sub"].(string)
	if !ok {
		return uuid.Nil, apierr.ErrInternal("Error al obtener el ID del usuario", "")
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, apierr.ErrInternal("Error al obtener el ID del usuario", "")
	}

	return id, nil
}
