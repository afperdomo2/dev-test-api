package push

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

// @Summary      Registrar dispositivo para push
// @Description  Registra o refresca el token de Expo del dispositivo. Es idempotente y reasigna el token si ya existía con otro usuario.
// @Tags         push
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  RegisterTokenRequest  true  "Token y plataforma"
// @Success      200  {object}  map[string]bool
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/push/tokens [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterTokenRequest
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

	if err := h.service.Register(userID, req.Token, req.Platform); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"registered": true})
}

// @Summary      Dar de baja el dispositivo
// @Description  Elimina el token si pertenece al usuario autenticado. Idempotente: si no existe o es ajeno responde removed=false sin error.
// @Tags         push
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  UnregisterTokenRequest  true  "Token a eliminar"
// @Success      200  {object}  map[string]bool
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/push/tokens/unregister [post]
func (h *Handler) Unregister(c *gin.Context) {
	var req UnregisterTokenRequest
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

	removed, err := h.service.Unregister(userID, req.Token)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"removed": removed})
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
