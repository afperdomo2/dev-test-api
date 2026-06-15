# dev-test-api

> REST API con **Go** + **Gin** · documentada con **Swagger** · live reload con **Air**

## 📋 Requisitos

| Herramienta | Versión |
|------------|---------|
| Go         | 1.26+   |
| make       | —       |

## 📁 Estructura del proyecto

```
dev-test-api/
├── main.go              # Punto de entrada, router y anotaciones Swagger
├── go.mod               # Módulo Go y dependencias
├── Makefile             # Comandos: dev, build, swagger, clean
├── .air.toml            # Configuración de live reload (Air)
├── middleware/
│   └── logger.go        # Middleware de logging personalizado
├── docs/                # Documentación Swagger auto-generada
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── README.md
```

## 🛠️ Scripts (Makefile)

| Comando | Descripción |
|---------|-------------|
| `make install` | Instala dependencias (`go mod tidy`) |
| `make dev` | Levanta el servidor con live reload (Air) |
| `make build` | Compila el binario en `./tmp/main` |
| `make run` | Corre el servidor directamente con `go run` |
| `make swagger` | Genera/regenera la documentación Swagger |
| `make clean` | Elimina `tmp/` y `docs/` |

## 🔥 Live reload con Air

El proyecto usa [Air](https://github.com/air-verse/air) para hot reload. Cuando guardás cambios en archivos `.go`, Air recompila y reinicia el servidor automáticamente.

### Configuración

`.air.toml` define qué archivos y carpetas monitorear, excluir y el comando de build. Ya viene pre-configurado para este proyecto.

### Uso

```bash
make dev
```

Esto ejecuta `go tool air` y arranca el servidor con live reload.

## 🚀 Ejecutar sin live reload

```bash
make run
# o directo:
go run main.go
```

El servidor arranca en `http://localhost:8080`.

## 📖 Documentación Swagger

Los endpoints se documentan mediante anotaciones en el código con [swaggo/swag](https://github.com/swaggo/swag).

### Generar docs

```bash
make swagger
# o directo:
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

Después de agregar o modificar anotaciones, volvé a correr `make swagger` para regenerar la documentación.

## 📦 Dependencias

| Paquete | Uso |
|---------|-----|
| `gin-gonic/gin` | Framework HTTP |
| `air-verse/air` | Live reload (project tool) |
| `swaggo/swag` | Generación de docs OpenAPI desde anotaciones |
| `swaggo/gin-swagger` | Integración Swagger UI con Gin |
| `swaggo/files` | Archivos estáticos de Swagger UI |

## 🧪 Health check

```bash
curl http://localhost:8080/health
# { "status": "ok" }
```

## 📝 Notas

- Los endpoints se documentan exclusivamente con anotaciones Swagger. No se documentan manualmente en este README.
- Gin incluye por defecto middleware de recovery (panic recovery) y logger.
- `make swagger` debe ejecutarse cada vez que se agregan o modifican anotaciones en los handlers.
- El paquete `docs/` es auto-generado y no debe editarse manualmente.
- Air se instaló como project tool (`go get -tool`), no requiere instalación global.
