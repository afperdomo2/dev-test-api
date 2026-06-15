package auth

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

// @Summary      Inicializar sistema (primer usuario admin)
// @Description  Crea el primer usuario admin del sistema. Solo funciona si no hay usuarios creados.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SetupRequest  true  "Credenciales del admin inicial"
// @Success      201   {object}  AuthResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError  "System Already Initialized"
// @Router       /api/v1/auth/setup [post]
func (h *Handler) Setup(c *gin.Context) {
	var req SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	res, err := h.service.Setup(req.Email, req.Password)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, res)
}

// @Summary      Iniciar sesión
// @Description  Autentica un usuario y devuelve un token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "Credenciales de login"
// @Success      200   {object}  AuthResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Router       /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	res, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, res)
}
