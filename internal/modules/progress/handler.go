package progress

import (
	"net/http"
	"strconv"

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

// @Summary      Preguntas pendientes
// @Description  Lista las preguntas marcadas para repasar cuyo próximo repaso es hoy
// @Tags         progress
// @Security     BearerAuth
// @Produce      json
// @Param        page      query  int  false  "Número de página (default: 1)"
// @Param        per_page  query  int  false  "Elementos por página (default: 20)"
// @Success      200  {object}  response.Meta  "Preguntas pendientes"
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/progress/upcoming [get]
func (h *Handler) Upcoming(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	items, total, err := h.service.Upcoming(userID, page, perPage)
	if err != nil {
		response.Problem(c, err.(*apierr.APIError))
		return
	}

	response.Paginated(c, http.StatusOK, items, response.Meta{Total: total, Page: page, PerPage: perPage})
}

// @Summary      Preguntas guardadas
// @Description  Lista todas las preguntas que el usuario ha guardado para repasar
// @Tags         progress
// @Security     BearerAuth
// @Produce      json
// @Param        page      query  int  false  "Número de página (default: 1)"
// @Param        per_page  query  int  false  "Elementos por página (default: 20)"
// @Success      200  {object}  response.Meta  "Preguntas guardadas"
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/progress/saved [get]
func (h *Handler) Saved(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	items, total, err := h.service.Saved(userID, page, perPage)
	if err != nil {
		response.Problem(c, err.(*apierr.APIError))
		return
	}

	response.Paginated(c, http.StatusOK, items, response.Meta{Total: total, Page: page, PerPage: perPage})
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
