package middleware

import "github.com/gin-gonic/gin"

// securityHeadersAPICSP es la CSP mínima para las respuestas de la API: no debe
// cargar ni ejecutar nada. La CSP de la SPA la sirve el hosting estático, no
// esta API (ver specs/SPEC-client-hardening.md).
const securityHeadersAPICSP = "default-src 'none'; frame-ancestors 'none'; base-uri 'none'"

// SecurityHeaders añade las cabeceras de seguridad a todas las respuestas.
// HSTS solo se emite si el borde sirve HTTPS (hsts=true).
func SecurityHeaders(hsts bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		h.Set("Content-Security-Policy", securityHeadersAPICSP)
		if hsts {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}
