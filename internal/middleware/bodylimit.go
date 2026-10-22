package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit acota el tamaño del body de la petición. Si se supera, las lecturas
// (ShouldBindJSON, FormFile, io.Copy…) fallan con *http.MaxBytesError: el body
// queda acotado en memoria, que es el objetivo de seguridad.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}

// IsPayloadTooLarge informa si el error viene de superar el BodyLimit.
func IsPayloadTooLarge(err error) bool {
	var maxErr *http.MaxBytesError
	return errors.As(err, &maxErr)
}
