package users

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

// @Summary      Listar usuarios (Admin)
// @Description  Lista todos los usuarios
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        per_page   query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        sort_by    query  string  false  "Campo de ordenación: email, created_at, updated_at"
// @Param        sort_order query  string  false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Lista de usuarios"
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/users [get]
func (h *Handler) List(c *gin.Context) {
	params, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	users, total, err := h.service.List(params.Page, params.PerPage, params.SortBy, params.SortOrder)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = ToUserResponse(u)
	}

	response.Paginated(c, http.StatusOK, result, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Crear usuario (Admin)
// @Description  Crea un nuevo usuario
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

	response.Success(c, http.StatusCreated, ToUserResponse(*user))
}

// @Summary      Obtener usuario (Admin)
// @Description  Obtiene un usuario por ID
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
		response.NotFound(c, "Usuario", c.Request.URL.Path)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, ToUserResponse(*user))
}

// @Summary      Eliminar usuario (Admin)
// @Description  Soft-delete de un usuario
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
		response.NotFound(c, "Usuario", c.Request.URL.Path)
		return
	}

	if err := h.service.Delete(id); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Usuario eliminado"})
}

// @Summary      Obtener perfil
// @Description  Obtiene los datos del usuario autenticado
// @Tags         profile
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  UserResponse
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	user, err := h.service.GetByID(userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, ToUserResponse(*user))
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
