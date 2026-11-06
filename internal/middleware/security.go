package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// securityHeadersAPICSP es la CSP mínima para las respuestas de la API: no debe
	// cargar ni ejecutar nada. La CSP de la SPA la sirve el hosting estático, no
	// esta API (ver specs/SPEC-client-hardening.md).
	securityHeadersAPICSP = "default-src 'none'; frame-ancestors 'none'; base-uri 'none'"

	// securityHeadersSwaggerCSP permite que Swagger UI cargue sus assets locales
	// (JS/CSS) y consulte su propio doc.json. Solo se aplica a /swagger/*, que
	// únicamente se registra cuando SWAGGER_ENABLED está activo.
	securityHeadersSwaggerCSP = "default-src 'self'; script-src 'self'; " +
		"style-src 'self' 'unsafe-inline'; img-src 'self' data:; " +
		"font-src 'self' data:; connect-src 'self'; " +
		"frame-ancestors 'none'; base-uri 'none'"
)

// SecurityHeaders añade las cabeceras de seguridad a todas las respuestas.
// HSTS solo se emite si el borde sirve HTTPS (hsts=true).
func SecurityHeaders(hsts bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		if isSwaggerPath(c.Request.URL.Path) {
			h.Set("Content-Security-Policy", securityHeadersSwaggerCSP)
		} else {
			h.Set("Content-Security-Policy", securityHeadersAPICSP)
		}
		if hsts {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

// isSwaggerPath cubre /swagger y /swagger/*; el primero solo redirige.
func isSwaggerPath(path string) bool {
	return path == "/swagger" || strings.HasPrefix(path, "/swagger/")
}
