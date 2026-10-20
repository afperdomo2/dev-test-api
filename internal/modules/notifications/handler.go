package notifications

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

// @Summary      Listar notificaciones
// @Description  Lista paginada de las notificaciones del usuario autenticado, ordenadas por fecha descendente.
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Param        page     query  int  false  "Número de página (default: 1)"
// @Param        perPage  query  int  false  "Elementos por página (default: 20, max: 100)"
// @Success      200  {object}  response.Meta  "Lista de notificaciones (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/notifications [get]
func (h *Handler) List(c *gin.Context) {
	pagination, err := common.ParsePagination(c, sortConfig)
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

	params := ListNotificationsParams{PaginationParams: pagination}
	items, total, listErr := h.service.List(userID, params)
	if listErr != nil {
		e := listErr.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, items, response.Meta{
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	})
}

// @Summary      Contar no leídas
// @Description  Devuelve cuántas notificaciones sin leer tiene el usuario autenticado.
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  UnreadCountResponse
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/notifications/unread-count [get]
func (h *Handler) UnreadCount(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	count, err := h.service.UnreadCount(userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, UnreadCountResponse{Count: count})
}

// @Summary      Marcar notificación como leída
// @Description  Marca una notificación propia como leída. Idempotente; una notificación ajena o inexistente devuelve 404.
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Notification ID"
// @Success      200  {object}  NotificationResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/notifications/{id}/read [post]
func (h *Handler) MarkRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Notificacion", c.Request.URL.Path)
		return
	}

	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	notification, markErr := h.service.MarkRead(id, userID)
	if markErr != nil {
		e := markErr.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, notification)
}

// @Summary      Marcar todas como leídas
// @Description  Marca como leídas todas las notificaciones no leídas del usuario autenticado.
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  MarkAllReadResponse
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/notifications/read-all [post]
func (h *Handler) MarkAllRead(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	updated, err := h.service.MarkAllRead(userID)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, MarkAllReadResponse{Updated: updated})
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
