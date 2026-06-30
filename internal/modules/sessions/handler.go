package sessions

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

// @Summary      Listar sesiones (con paginación)
// @Description  Lista las sesiones de estudio del usuario autenticado
// @Tags         sessions
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "Número de página (default: 1)"
// @Param        perPage     query  int     false  "Elementos por página (default: 20, max: 100)"
// @Param        sortBy       query  string  false  "Campo de ordenación: status, score, started_at, created_at"
// @Param        sortOrder query  string  false  "Dirección: asc o desc (default: desc)"
// @Param        status    query  string  false  "Filtro por estado: in_progress, completed, cancelled"
// @Success      200  {object}  response.Meta  "Lista de sesiones (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/sessions [get]
func (h *Handler) List(c *gin.Context) {
	params, err := common.ParsePagination(c, sortConfig)
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

	listParams := ListSessionsParams{
		PaginationParams: params,
		Status:           c.Query("status"),
	}

	sessions, total, err := h.service.List(userID, listParams)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, sessions, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Crear sesión
// @Description  Crea una nueva sesión de estudio
// @Tags         sessions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateSessionRequest  true  "Datos de la sesión"
// @Success      201   {object}  SessionResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Router       /api/v1/sessions [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateSessionRequest
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

	session, err := h.service.Create(userID, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, session)
}

// @Summary      Obtener sesión
// @Description  Obtiene el detalle de una sesión con sus respuestas
// @Tags         sessions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Session ID"
// @Success      200  {object}  SessionDetailResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/sessions/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Sesion", c.Request.URL.Path)
		return
	}

	session, err := h.service.GetByID(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, session)
}

// @Summary      Finalizar sesión
// @Description  Finaliza una sesión — el score se calcula automáticamente
// @Tags         sessions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Session ID"
// @Success      200  {object}  SessionResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Router       /api/v1/sessions/{id}/finish [put]
func (h *Handler) Finish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Sesion", c.Request.URL.Path)
		return
	}

	session, err := h.service.Finish(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, session)
}

// @Summary      Siguiente pregunta
// @Description  Obtiene la siguiente pregunta sin responder de la sesión. Si no hay preguntas disponibles y la sesión es modo "generate", genera una nueva pregunta por IA en tiempo real. Puede devolver 404 si se alcanzó el límite de preguntas o no hay más disponibles, y 409 si la sesión ya fue completada.
// @Tags         sessions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Session ID"
// @Success      200  {object}  NextQuestionResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /api/v1/sessions/{id}/next [get]
func (h *Handler) NextQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Sesion", c.Request.URL.Path)
		return
	}

	next, err := h.service.NextQuestion(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, next)
}

// @Summary      Responder pregunta
// @Description  Registra la respuesta del usuario y actualiza el progreso SM-2
// @Tags         sessions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string         true  "Session ID"
// @Param        body  body  AnswerRequest  true  "Respuesta"
// @Success      200  {object}  SessionAnswerResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Router       /api/v1/sessions/{id}/answer [post]
func (h *Handler) Answer(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Sesion", c.Request.URL.Path)
		return
	}

	var req AnswerRequest
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

	resp, err := h.service.Answer(sessionID, userID, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// @Summary      Resumen de sesión
// @Description  Obtiene un resumen liviano de la sesión (contadores, estado)
// @Tags         sessions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Session ID"
// @Success      200  {object}  SessionSummaryResponse
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/sessions/{id}/summary [get]
func (h *Handler) Summary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Sesion", c.Request.URL.Path)
		return
	}

	summary, err := h.service.Summary(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, summary)
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
