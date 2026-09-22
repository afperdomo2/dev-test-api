package tips

import (
	"bytes"
	"io"
	"net/http"
	"strings"

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

// @Summary      Listar tips (con paginación)
// @Description  Lista todos los tips públicos con paginación y filtro opcional por tema
// @Tags         tips
// @Security     BearerAuth
// @Produce      json
// @Param        page     query  int    false  "Número de página (default: 1)"
// @Param        perPage  query  int    false  "Elementos por página (default: 20, max: 100)"
// @Param        topicId  query  string false  "Filtrar por tema (UUID)"
// @Param        sortBy   query  string false  "Campo de ordenación: created_at, updated_at"
// @Param        sortOrder query string false  "Dirección: asc o desc (default: desc)"
// @Success      200  {object}  response.Meta  "Lista de tips (con paginación)"
// @Failure      401  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/tips [get]
func (h *Handler) List(c *gin.Context) {
	pagination, err := common.ParsePagination(c, sortConfig)
	if err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	params := ListTipsParams{
		PaginationParams: pagination,
	}

	if raw := c.Query("topicId"); raw != "" {
		if id, err2 := uuid.Parse(raw); err2 == nil {
			params.TopicID = id
		}
	}

	tips, total, err := h.service.List(params)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Paginated(c, http.StatusOK, tips, response.Meta{Total: total, Page: params.Page, PerPage: params.PerPage})
}

// @Summary      Obtener tip
// @Description  Obtiene un tip por ID con los detalles del tema asociado
// @Tags         tips
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Tip ID"
// @Success      200  {object}  TipResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/tips/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Tip", c.Request.URL.Path)
		return
	}

	tip, err := h.service.GetByID(id)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, tip)
}

// @Summary      Crear tip manual
// @Description  Crea un nuevo tip (solo administradores). Los tips son siempre públicos.
// @Tags         tips
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  CreateTipRequest  true  "Datos del tip"
// @Success      201   {object}  TipResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      422   {object}  apierr.APIError
// @Router       /api/v1/tips [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateTipRequest
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

	tip, err := h.service.Create(userID, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusCreated, tip)
}

// @Summary      Actualizar tip
// @Description  Actualiza el contenido o el tema de un tip (solo administradores)
// @Tags         tips
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string             true  "Tip ID"
// @Param        body  body  UpdateTipRequest   true  "Datos a actualizar"
// @Success      200   {object}  TipResponse
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Failure      422   {object}  apierr.APIError
// @Router       /api/v1/tips/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Tip", c.Request.URL.Path)
		return
	}

	var req UpdateTipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error(), c.Request.URL.Path)
		return
	}

	tip, err := h.service.Update(id, req)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, tip)
}

// @Summary      Eliminar tip
// @Description  Elimina un tip (solo administradores)
// @Tags         tips
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  string  true  "Tip ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Router       /api/v1/tips/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NotFound(c, "Tip", c.Request.URL.Path)
		return
	}

	if err := h.service.Delete(id); err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Tip eliminado"})
}

// @Summary      Importar tips por CSV
// @Description  Importa tips masivos desde CSV (solo administradores, sin cupo diario). Columnas requeridas: content, topic (slug del tema del sistema). Los slugs desconocidos generan error por fila.
// @Tags         tips
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData  file    false  "Archivo CSV"
// @Param        content  formData  string  false  "Texto CSV pegado (alternativa a file)"
// @Success      200  {object}  ImportResult
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      422  {object}  apierr.APIError
// @Router       /api/v1/tips/import [post]
func (h *Handler) Import(c *gin.Context) {
	userID, apiErr := getUserID(c)
	if apiErr != nil {
		apiErr.Instance = c.Request.URL.Path
		response.Problem(c, apiErr)
		return
	}

	var reader io.Reader

	if fh, err := c.FormFile("file"); err == nil {
		f, err2 := fh.Open()
		if err2 != nil {
			response.ValidationError(c, "No se pudo leer el archivo", c.Request.URL.Path)
			return
		}
		defer f.Close()
		buf := new(bytes.Buffer)
		if _, err2 := io.Copy(buf, f); err2 != nil {
			response.ValidationError(c, "No se pudo leer el archivo", c.Request.URL.Path)
			return
		}
		reader = buf
	}

	if reader == nil {
		content := c.PostForm("content")
		if strings.TrimSpace(content) == "" {
			response.ValidationError(c, "Se requiere un archivo (file) o contenido pegado (content)", c.Request.URL.Path)
			return
		}
		reader = strings.NewReader(content)
	}

	result, err := h.service.Import(userID, reader)
	if err != nil {
		e := err.(*apierr.APIError)
		e.Instance = c.Request.URL.Path
		response.Problem(c, e)
		return
	}

	response.Success(c, http.StatusOK, result)
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
