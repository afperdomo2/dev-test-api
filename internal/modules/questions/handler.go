package questions

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

// @Summary      Listar preguntas
// @Description  Lista todas las preguntas disponibles, con filtros opcionales
// @Tags         questions
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        per_page   query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        type       query  string  false  "Filtrar por tipo (single_choice, multiple_choice, code_completion)"
// @Param        difficulty query  string  false  "Filtrar por dificultad (beginner, intermediate, advanced)"
// @Param        sort_by    query  string  false  "Campo de ordenación: type, difficulty, created_at, updated_at"
// @Param        sort_order query  string  false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Lista de preguntas"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/questions [get]
func (h *Handler) List(c *gin.Context) {
	params, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	filters := QuestionFilters{
		Type:       c.Query("type"),
		Difficulty: c.Query("difficulty"),
	}

	questions, total, err := h.service.List(params.Page, params.PerPage, params.SortBy, params.SortOrder, filters)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, questions, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Obtener pregunta
// @Description  Obtiene una pregunta por ID con todos sus detalles (opciones, código, temas)
// @Tags         questions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Question ID"
// @Success      200  {object}  QuestionResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/questions/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Pregunta", c.Request.URL.Path)
		return
	}

	question, err := h.service.GetByID(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, question)
}

// @Summary      Crear pregunta (Admin)
// @Description  Crea una nueva pregunta con opciones o código
// @Tags         questions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateQuestionRequest  true  "Datos de la pregunta"
// @Success      201   {object}  QuestionResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Router       /api/v1/questions [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateQuestionRequest
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

	question, err := h.service.Create(userID, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, question)
}

// @Summary      Actualizar pregunta (Admin)
// @Description  Actualiza una pregunta existente
// @Tags         questions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                true  "Question ID"
// @Param        body  body  UpdateQuestionRequest  true  "Datos a actualizar"
// @Success      200   {object}  QuestionResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Router       /api/v1/questions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Pregunta", c.Request.URL.Path)
		return
	}

	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	question, err := h.service.Update(id, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, question)
}

// @Summary      Eliminar pregunta (Admin)
// @Description  Elimina una pregunta (soft delete)
// @Tags         questions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Question ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/questions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Pregunta", c.Request.URL.Path)
		return
	}

	if err := h.service.Delete(id); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Pregunta eliminada"})
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
