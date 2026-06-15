# dev-test-api

> REST API con **Go** + **Gin** · documentada con **Swagger**

## 📋 Requisitos

| Herramienta | Versión |
|------------|---------|
| Go         | 1.26+   |
| swag CLI   | latest  |

## 📁 Estructura del proyecto

```
dev-test-api/
├── main.go              # Punto de entrada, router y anotaciones Swagger
├── go.mod               # Módulo Go y dependencias
├── middleware/
│   └── logger.go        # Middleware de logging personalizado
├── docs/                # Documentación Swagger auto-generada
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── README.md
```

## ⚙️ Instalación

```bash
git clone <repo-url>
cd dev-test-api
go mod tidy
```

## 🔧 Instalar swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🚀 Ejecutar

```bash
go run main.go
```

El servidor arranca en `http://localhost:8080`.

## 📖 Documentación Swagger

Los endpoints se documentan mediante anotaciones en el código con [swaggo/swag](https://github.com/swaggo/swag).

### Generar docs

```bash
swag init -g main.go
```

Este comando genera la carpeta `docs/` con la especificación OpenAPI en JSON y YAML, más el código Go que Swagger UI consume.

### Ver documentación

Con el servidor corriendo, abrí:

```
http://localhost:8080/swagger/index.html
```

### Agregar documentación a un endpoint

```go
// @Summary      Health check
// @Description  Verifica que el servidor esté corriendo
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
```

Después de agregar o modificar anotaciones, volvé a correr `swag init -g main.go` para regenerar la documentación.

## 📦 Dependencias

| Paquete | Uso |
|---------|-----|
| `gin-gonic/gin` | Framework HTTP |
| `swaggo/swag` | Generación de docs OpenAPI desde anotaciones |
| `swaggo/gin-swagger` | Integración Swagger UI con Gin |
| `swaggo/files` | Archivos estáticos de Swagger UI |
| `google/uuid` | Generación de IDs únicos |

## 🧪 Health check

```bash
curl http://localhost:8080/health
# { "status": "ok" }
```

## 📝 Notas

- Los endpoints se documentan exclusivamente con anotaciones Swagger. No se documentan manualmente en este README.
- Gin incluye por defecto middleware de recovery (panic recovery) y logger.
- `swag init` debe ejecutarse cada vez que se agregan o modifican anotaciones en los handlers.
- El paquete `docs/` es auto-generado y no debe editarse manualmente.
