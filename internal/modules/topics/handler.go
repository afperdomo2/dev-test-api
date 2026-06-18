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

// @Summary      Listar temas (con paginación)
// @Description  Lista los temas según el rol del usuario. Admin ve solo temas del sistema (is_system=true). Usuarios normales ven solo sus propios temas personalizados (is_system=false).
// @Tags         topics
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        perPage     query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        sortBy       query  string  false  "Campo de ordenación: name, slug, category, created_at"
// @Param        sortOrder query  string  false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Lista de temas (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/topics [get]
func (h *Handler) List(c *gin.Context) {
	params, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	isAdmin, userID, apiErr := getUserRoleAndID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	topics, total, err := h.service.List(params, isAdmin, userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, topics, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Obtener tema
// @Description  Obtiene un tema por ID. Admin solo puede ver temas del sistema (is_system=true). Usuarios normales pueden ver temas del sistema y sus propios temas.
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

	isAdmin, userID, apiErr := getUserRoleAndID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	topic, err := h.service.GetByID(id, isAdmin, userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, topic)
}

// @Summary      Crear tema
// @Description  Crea un nuevo tema. Admin crea temas del sistema (is_system=true), usuarios normales crean temas personalizados (is_system=false).
// @Tags         topics
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateTopicRequest  true  "Datos del tema"
// @Success      201   {object}  TopicResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError
// @Router       /api/v1/topics [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	isAdmin, userID, apiErr := getUserRoleAndID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	topic, err := h.service.Create(userID, req, isAdmin)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, topic)
}

// @Summary      Actualizar tema
// @Description  Actualiza un tema existente. Admin solo puede modificar temas del sistema (is_system=true). Usuarios normales solo pueden modificar sus propios temas (is_system=false).
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

	isAdmin, userID, apiErr := getUserRoleAndID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	topic, err := h.service.Update(id, req, isAdmin, userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, topic)
}

// @Summary      Eliminar tema
// @Description  Elimina un tema. Admin solo puede eliminar temas del sistema (is_system=true). Usuarios normales solo pueden eliminar sus propios temas (is_system=false).
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

	isAdmin, userID, apiErr := getUserRoleAndID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	if err := h.service.Delete(id, isAdmin, userID); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Tema eliminado"})
}

func getUserRoleAndID(c *gin.Context) (bool, uuid.UUID, *apierr.APIError) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return false, uuid.Nil, apierr.ErrUnauthorized("Usuario no autenticado", "")
	}

	mapClaims, ok := claims.(*jwt.MapClaims)
	if !ok {
		return false, uuid.Nil, apierr.ErrInternal("Error al obtener los claims del usuario", "")
	}

	sub, ok := (*mapClaims)["sub"].(string)
	if !ok {
		return false, uuid.Nil, apierr.ErrInternal("Error al obtener el ID del usuario", "")
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return false, uuid.Nil, apierr.ErrInternal("Error al obtener el ID del usuario", "")
	}

	isAdmin, _ := (*mapClaims)["is_admin"].(bool)

	return isAdmin, id, nil
}
