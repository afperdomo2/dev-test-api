# dev-test-api

> REST API con **Go** + **Gin** · **JWT** · **PostgreSQL** · **Swagger** · live reload con **Air**

## 📋 Requisitos

| Herramienta | Versión |
|------------|---------|
| Go         | 1.26+   |
| PostgreSQL | 14+     |
| make       | —       |

## 📁 Estructura del proyecto

```
dev-test-api/
├── main.go                    # Punto de entrada, wiring de dependencias
├── config/
│   └── config.go              # Configuración desde env vars
├── database/
│   └── database.go            # Conexión GORM + AutoMigrate
├── internal/                  # Módulos de negocio (privados)
│   ├── auth/                  # Autenticación (setup, login, JWT)
│   ├── users/                 # Gestión de usuarios (CRUD)
│   └── questions/             # Preguntas (placeholder)
├── middleware/
│   ├── logger.go              # Logging de requests
│   └── auth.go                # Validación JWT + admin check
├── pkg/
│   ├── response/              # Helpers: Success(), Problem() RFC 9457
│   └── apierr/                # Errores tipados de API
├── docs/                      # Swagger auto-generado
├── .env.example
├── .air.toml
├── Makefile
└── README.md
```

Cada módulo en `internal/` contiene todo lo que necesita: modelo, store (DB), servicio (lógica) y handler (HTTP). Esto permite extraer un módulo como microservicio sin desarmar el proyecto.

## 🔐 Flujo de autenticación

### Bootstrap — primer usuario admin

```
POST /api/v1/auth/setup
```

Solo funciona si **no hay usuarios en la base de datos**. Crea el primer usuario con `is_admin: true` y devuelve un JWT.

Si ya existen usuarios, devuelve `409 Conflict`.

### Login

```
POST /api/v1/auth/login
```

Devuelve un JWT que se usa como `Authorization: Bearer <token>` en los endpoints protegidos.

### Usuarios (admin)

```
GET    /api/v1/users          # listar
POST   /api/v1/users          # crear
GET    /api/v1/users/:id      # obtener
DELETE /api/v1/users/:id      # soft-delete
```

## ⚙️ Variables de entorno

Copiá `.env.example` a `.env` y ajustá los valores:

```bash
cp .env.example .env
```

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto del servidor |
| `GIN_MODE` | `debug` | `debug` o `release` |
| `DB_HOST` | `localhost` | Host de PostgreSQL |
| `DB_PORT` | `5432` | Puerto de PostgreSQL |
| `DB_USER` | `postgres` | Usuario de la DB |
| `DB_PASSWORD` | `secret` | Contraseña de la DB |
| `DB_NAME` | `dev_test_api` | Nombre de la base de datos |
| `DB_SSL_MODE` | `disable` | SSL mode |
| `JWT_SECRET` | — | Secreto para firmar JWTs |
| `JWT_EXPIRY_HOURS` | `24` | Expiración del token |

## 🛠️ Scripts (Makefile)

| Comando | Descripción |
|---------|-------------|
| `make install` | Instala dependencias (`go mod tidy`) |
| `make dev` | Live reload con Air |
| `make build` | Compila binario en `./tmp/main` |
| `make run` | Corre el servidor con `go run` |
| `make swagger` | Regenera docs Swagger |
| `make clean` | Elimina `tmp/` y `docs/` |

## 📖 Documentación Swagger

Los endpoints se documentan mediante anotaciones en el código.

```bash
make swagger
```

Con el servidor corriendo: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 📦 Dependencias

| Paquete | Uso |
|---------|-----|
| `gin-gonic/gin` | Framework HTTP |
| `gorm.io/gorm` | ORM |
| `gorm.io/driver/postgres` | Driver PostgreSQL |
| `golang-jwt/jwt/v5` | JWT tokens |
| `golang.org/x/crypto` | bcrypt |
| `joho/godotenv` | Variables de entorno |
| `swaggo/swag` | Docs OpenAPI |
| `swaggo/gin-swagger` | Swagger UI |
| `air-verse/air` | Live reload (project tool) |

## 🧪 Primer uso

```bash
# 1. Instalar dependencias
make install

# 2. Configurar .env con los datos de tu PostgreSQL
cp .env.example .env

# 3. Crear el primer usuario admin
make dev
curl -X POST http://localhost:8080/api/v1/auth/setup \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123456"}'

# 4. Usar el token devuelto para los endpoints protegidos
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <tu-token>"
```

## 📝 Respuestas

**Error** — RFC 9457 Problem Detail:

```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 422,
  "detail": "email: required",
  "instance": "/api/v1/users"
}
```

**Éxito** — envelope `{ "data": ... }`:

```json
{
  "data": {
    "token": "...",
    "user": { "id": "...", "email": "...", "is_admin": true }
  }
}
```

## 📝 Notas

- Las tablas se crean/actualizan automáticamente al arrancar (AutoMigrate de GORM).
- `docs/` y `tmp/` son auto-generados y no se editan manualmente.
- Air se instaló como project tool (`go get -tool`), no requiere instalación global.
- `internal/` es privado al módulo Go — no puede ser importado desde fuera del proyecto.
