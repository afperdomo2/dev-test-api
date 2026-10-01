---
description: Crea una migración goose y guía su escritura
---
Descripción: $ARGUMENTS

1. `make migrate-status` — verifica el estado actual de las migraciones.
2. `make migrate-create name=$ARGUMENTS` — crea `NNNNN_<desc>.sql` (up+down)
   en `internal/database/migrations/`.
3. Escribe el SQL up/down siguiendo `.agents/backend/data-model.md`.
4. Reglas: las migraciones son la única fuente de verdad del schema (nunca
   AutoMigrate); up y down deben ser consistentes entre sí.
5. Valida: `make migrate-up` y comprueba que `make migrate-down` revierte limpio.
