package gamification

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

// @Summary      Ver resumen de gamificación
// @Description  Devuelve XP total, nivel, racha, meta diaria y XP por tema del usuario autenticado.
// @Tags         gamification
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  OverviewResponse
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/gamification/overview [get]
func (h *Handler) Overview(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	overview, err := h.service.GetOverview(userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, overview)
}

// @Summary      Cambiar meta diaria de XP
// @Description  Fija la meta diaria de XP del usuario autenticado (entre 10 y 200). Devuelve el reward con la meta actualizada.
// @Tags         gamification
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  DailyGoalRequest  true  "Meta diaria"
// @Success      200   {object}  Reward
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      422   {object}  apierr.APIError
// @Router       /api/v1/gamification/daily-goal [put]
func (h *Handler) SetDailyGoal(c *gin.Context) {
	var req DailyGoalRequest
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

	reward, err := h.service.SetDailyGoal(userID, req.DailyGoalXp)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, reward)
}

// @Summary      Congelar racha
// @Description  Consume un freeze disponible para cubrir el día de hoy sin actividad y mantener la racha.
// @Tags         gamification
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  Reward
// @Failure      401  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Router       /api/v1/gamification/streak/freeze [post]
func (h *Handler) Freeze(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	reward, err := h.service.Freeze(userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, reward)
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
