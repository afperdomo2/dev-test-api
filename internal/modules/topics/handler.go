package topics

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

// @Summary      Listar temas
// @Description  Lista todos los temas disponibles
// @Tags         topics
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        per_page   query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        sort_by    query  string  false  "Campo de ordenación: name, slug, category, created_at"
// @Param        sort_order query  string  false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Lista de temas"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/topics [get]
func (h *Handler) List(c *gin.Context) {
	params, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	topics, total, err := h.service.List(params.Page, params.PerPage, params.SortBy, params.SortOrder)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, topics, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Obtener tema
// @Description  Obtiene un tema por ID
// @Tags         topics
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Topic ID"
// @Success      200  {object}  TopicResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/topics/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Tema", c.Request.URL.Path)
		return
	}

	topic, err := h.service.GetByID(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, topic)
}

// @Summary      Crear tema (Admin)
// @Description  Crea un nuevo tema personalizado
// @Tags         topics
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateTopicRequest  true  "Datos del tema"
// @Success      201   {object}  TopicResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError
// @Router       /api/v1/topics [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateTopicRequest
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

	topic, err := h.service.Create(userID, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, topic)
}

// @Summary      Actualizar tema (Admin)
// @Description  Actualiza un tema existente
// @Tags         topics
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string              true  "Topic ID"
// @Param        body  body  UpdateTopicRequest  true  "Datos a actualizar"
// @Success      200   {object}  TopicResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Router       /api/v1/topics/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Tema", c.Request.URL.Path)
		return
	}

	var req UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	topic, err := h.service.Update(id, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, topic)
}

// @Summary      Eliminar tema (Admin)
// @Description  Elimina un tema
// @Tags         topics
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Topic ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/topics/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Tema", c.Request.URL.Path)
		return
	}

	if err := h.service.Delete(id); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Tema eliminado"})
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
