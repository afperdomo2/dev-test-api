---
description: Estado de la BD y migraciones pendientes
---
Ejecuta `make migrate-status` y resume: migraciones aplicadas y pendientes.
Si hay pendientes, sugiere `make migrate-up`. Si el usuario necesita datos de
prueba, recuerda `make db-seed`. Ante cualquier duda del schema, consulta
`internal/database/migrations/*.sql` y `.agents/backend/data-model.md`.
