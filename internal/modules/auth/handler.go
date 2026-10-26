package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
)

const setupDisabledDetail = "Configuración inicial deshabilitada o token inválido"

type Handler struct {
	service    Service
	setupToken string
}

func NewHandler(service Service, setupToken string) *Handler {
	return &Handler{service: service, setupToken: setupToken}
}

// setupTokenMatches compara en tiempo constante sobre SHA-256. El hash evita
// filtrar la longitud del token; el vacío siempre falla (fail closed, D-38).
func setupTokenMatches(got, want string) bool {
	if want == "" {
		return false
	}
	a, b := sha256.Sum256([]byte(got)), sha256.Sum256([]byte(want))
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}

// @Summary      Inicializar sistema (primer usuario admin)
// @Description  Crea el primer usuario admin del sistema. Solo funciona si no hay usuarios creados y el body incluye un `setupToken` válido (SETUP_TOKEN del servidor). Sin token o con token incorrecto responde 403; si ya hay usuarios, 409.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SetupRequest  true  "Credenciales del admin inicial + setupToken"
// @Success      201   {object}  AuthResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      409   {object}  apierr.APIError  "System Already Initialized"
// @Failure      429   {object}  apierr.APIError
// @Router       /api/v1/auth/setup [post]
func (h *Handler) Setup(c *gin.Context) {
	var req SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	if !setupTokenMatches(req.SetupToken, h.setupToken) {
		response.Problem(c, apierr.ErrForbidden(setupDisabledDetail, c.Request.URL.Path))
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

// @Summary      Estado del sistema
// @Description  Indica si el sistema ya fue inicializado (tiene al menos un usuario)
// @Tags         auth
// @Produce      json
// @Success      200  {object}  StatusResponse
// @Router       /api/v1/auth/status [get]
func (h *Handler) Status(c *gin.Context) {
	res, err := h.service.Initialized()
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	res.SetupRequiresToken = h.setupToken != ""
	response.Success(c, http.StatusOK, res)
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
// @Failure      429   {object}  apierr.APIError
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
