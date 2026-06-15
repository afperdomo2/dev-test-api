package users

import (
	"net/http"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary      Listar usuarios
// @Description  Lista todos los usuarios (solo admin)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Meta  "Lista de usuarios"
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Router       /api/v1/users [get]
func (h *Handler) List(c *gin.Context) {
	users, err := h.service.List()
	if err != nil {
		response.Problem(c, err.(*apierr.APIError))
		return
	}

	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = u.ToResponse()
	}

	response.Success(c, http.StatusOK, result)
}

// @Summary      Crear usuario
// @Description  Crea un nuevo usuario (solo admin)
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateUserRequest  true  "Datos del usuario"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError
// @Router       /api/v1/users [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	user, err := h.service.Create(req.Email, req.Password, req.IsAdmin)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, user.ToResponse())
}

// @Summary      Obtener usuario
// @Description  Obtiene un usuario por ID (admin o el propio usuario)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "User ID"
// @Success      200  {object}  UserResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/users/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "user", c.Request.URL.Path)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, user.ToResponse())
}

// @Summary      Eliminar usuario
// @Description  Soft-delete de un usuario (solo admin)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "user", c.Request.URL.Path)
		return
	}

	if err := h.service.Delete(id); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "user deleted"})
}
