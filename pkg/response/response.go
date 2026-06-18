package response

import (
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
)

type Meta struct {
	Total   int64 `json:"total,omitempty"`
	Page    int   `json:"page,omitempty"`
	PerPage int   `json:"perPage,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func Paginated(c *gin.Context, status int, data any, meta Meta) {
	c.JSON(status, gin.H{
		"data": data,
		"meta": meta,
	})
}

func Problem(c *gin.Context, err *apierr.APIError) {
	c.AbortWithStatusJSON(err.Status, err)
}

func ValidationError(c *gin.Context, detail string, instance string) {
	Problem(c, apierr.ErrValidation(detail, instance))
}

func NotFound(c *gin.Context, entity string, instance string) {
	Problem(c, apierr.ErrNotFound(entity, instance))
}
