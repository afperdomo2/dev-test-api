package opencode

import (
	"net/http"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetUsage godoc
// @Summary      Consumo de OpenCode Go (solo admin)
// @Description  Consulta el endpoint oficial de OpenCode Go y devuelve el consumo rolling (5h), semanal y mensual con % usado y fecha de reset. Requiere rol de administrador; el backend reutiliza AI_API_KEY.
// @Tags         opencode
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  opencode.UsageResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      502  {object}  apierr.APIError
// @Failure      503  {object}  apierr.APIError
// @Router       /api/v1/opencode/usage [get]
func (h *Handler) GetUsage(c *gin.Context) {
	data, err := h.service.GetUsage(c.Request.Context())
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, data)
}
