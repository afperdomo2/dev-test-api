package settings

import (
	"net/http"

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

// @Summary      Listar configuración
// @Description  Devuelve los ajustes tunables con su valor efectivo (el de la base de datos o el default), su tipo y rango. Solo administradores. No incluye secretos.
// @Tags         settings
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Meta  "Lista de ajustes"
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Router       /api/v1/settings [get]
func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		problem(c, err)
		return
	}

	response.Success(c, http.StatusOK, items)
}

// @Summary      Actualizar configuración
// @Description  Actualiza uno o más ajustes tunables (solo administradores). Valida tipo y rango, y que ai.default_quantity no supere ai.max_per_job.
// @Tags         settings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateSettingsRequest  true  "Ajustes a actualizar"
// @Success      200   {object}  response.Meta  "Configuración actualizada"
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      422   {object}  apierr.APIError
// @Router       /api/v1/settings [put]
func (h *Handler) Update(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	items, err := h.service.Update(req.Settings, userID)
	if err != nil {
		problem(c, err)
		return
	}

	response.Success(c, http.StatusOK, items)
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
