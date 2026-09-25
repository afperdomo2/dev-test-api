package generation

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

// @Summary      Encolar generación de preguntas con IA
// @Description  Encola un trabajo de generación por lotes (solo administradores). Resuelve N preguntas en una única llamada a la IA. La cantidad es opcional y por defecto usa el ajuste ai.default_quantity. Rechaza si se supera el tope diario.
// @Tags         generation
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateJobRequest  true  "Trabajo a encolar"
// @Success      201   {object}  response.Meta  "Trabajo encolado"
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      422   {object}  apierr.APIError
// @Failure      429   {object}  apierr.APIError
// @Router       /api/v1/generation [post]
func (h *Handler) CreateJob(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	job, err := h.service.CreateJob(c.Request.Context(), userID, req)
	if err != nil {
		problem(c, err)
		return
	}

	response.Success(c, http.StatusCreated, job)
}

// @Summary      Listar trabajos de generación
// @Description  Lista paginada de los trabajos de generación con IA (solo administradores), filtrable por estado.
// @Tags         generation
// @Security     BearerAuth
// @Produce      json
// @Param        status     query     string  false  "Filtrar por estado: pending, running, completed o failed"
// @Param        page       query     int     false  "Página (default 1)"
// @Param        perPage    query     int     false  "Elementos por página (default 20, max 100)"
// @Param        sortBy     query     string  false  "Campo de orden: created_at, started_at o finished_at"
// @Param        sortOrder  query     string  false  "Orden: asc o desc"
// @Success      200        {object}  response.Meta  "Lista de trabajos"
// @Failure      401        {object}  apierr.APIError
// @Failure      403        {object}  apierr.APIError
// @Failure      422        {object}  apierr.APIError
// @Router       /api/v1/generation [get]
func (h *Handler) ListJobs(c *gin.Context) {
	pagination, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	params := ListJobsParams{PaginationParams: pagination, Status: c.Query("status")}

	jobs, total, err := h.service.ListJobs(params)
	if err != nil {
		problem(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, jobs, response.Meta{
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	})
}

// @Summary      Ver un trabajo de generación
// @Description  Devuelve el detalle de un trabajo de generación con IA (solo administradores).
// @Tags         generation
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "ID del trabajo (UUID)"
// @Success      200  {object}  response.Meta  "Detalle del trabajo"
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/generation/{id} [get]
func (h *Handler) GetJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "El ID del trabajo no es un UUID válido", c.Request.URL.Path)
		return
	}

	job, err := h.service.GetJob(id)
	if err != nil {
		problem(c, err)
		return
	}

	response.Success(c, http.StatusOK, job)
}

// @Summary      Estadísticas de generación con IA
// @Description  Consumo del día (reservado y generado), topes vigentes y colas por estado (solo administradores).
// @Tags         generation
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Meta  "Estadísticas"
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Router       /api/v1/generation/stats [get]
func (h *Handler) Stats(c *gin.Context) {
	stats, err := h.service.Stats()
	if err != nil {
		problem(c, err)
		return
	}

	response.Success(c, http.StatusOK, stats)
}

func problem(c *gin.Context, err error) {
	e, ok := err.(*apierr.APIError)
	if !ok {
		e = apierr.ErrInternal("Error inesperado", c.Request.URL.Path)
	}
	e.Instance = c.Request.URL.Path
	response.Problem(c, e)
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
