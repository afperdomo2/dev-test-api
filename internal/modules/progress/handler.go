package progress

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

// @Summary      Responder pregunta
// @Description  Registra una respuesta correcta/incorrecta y actualiza el progreso SM-2
// @Tags         progress
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        question_id  path  string         true  "Question ID"
// @Param        body         body  AnswerRequest  true  "Respuesta"
// @Success      200  {object}  ProgressResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/progress/{question_id}/answer [post]
func (h *Handler) Answer(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("question_id"))
	if err != nil {
		response.NotFound(c, "Pregunta", c.Request.URL.Path)
		return
	}

	var req AnswerRequest
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

	resp, err := h.service.Answer(userID, questionID, req.IsCorrect)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// @Summary      Preguntas pendientes (con paginación)
// @Description  Lista las preguntas marcadas para repasar cuyo próximo repaso es hoy
// @Tags         progress
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        perPage     query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        sortBy       query  string  false  "Campo de ordenación: next_review_at, repetitions, ease_factor"
// @Param        sortOrder query  string  false  "Dirección: asc o desc (default: asc)"
// @Success      200  {object}  response.Meta  "Preguntas pendientes (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/progress/upcoming [get]
func (h *Handler) Upcoming(c *gin.Context) {
	params, err := common.ParsePagination(c, upcomingSortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	items, total, err := h.service.Upcoming(userID, params)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, items, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Preguntas guardadas (con paginación)
// @Description  Lista todas las preguntas que el usuario ha guardado para repasar
// @Tags         progress
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        perPage     query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        sortBy       query  string  false  "Campo de ordenación: updated_at, repetitions, ease_factor"
// @Param        sortOrder query  string  false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Preguntas guardadas (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/progress/saved [get]
func (h *Handler) Saved(c *gin.Context) {
	params, err := common.ParsePagination(c, savedSortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	items, total, err := h.service.Saved(userID, params)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, items, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Marcar/desmarcar pregunta
// @Description  Alterna el estado is_saved de una pregunta para el usuario
// @Tags         progress
// @Security     BearerAuth
// @Produce      json
// @Param        question_id  path  string  true  "Question ID"
// @Success      200  {object}  ProgressResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/progress/{question_id}/toggle-save [post]
func (h *Handler) ToggleSave(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("question_id"))
	if err != nil {
		response.NotFound(c, "Pregunta", c.Request.URL.Path)
		return
	}

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	resp, err := h.service.ToggleSave(userID, questionID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, resp)
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
