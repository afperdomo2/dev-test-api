<!--
Título: <tipo>(<alcance>): resumen corto
Ejemplo: fix(sessions): enforce ownership on session detail
-->

## Qué cambia

<!-- 1-3 líneas. Si hay contexto de una tarea/spec, enlázalo. -->

## Cómo se probó

- [ ] `make check` / `make test` (backend)
- [ ] `make fe-check` / `make fe-test` (web)
- [ ] `pnpm --filter mobile test:run` (mobile, si aplica)
- [ ] Prueba manual: <!-- qué hiciste y qué esperabas -->

## Impacto

- **DB/migraciones**: <!-- ninguna / nombre de la migración -->
- **API**: <!-- sin cambios / endpoints afectados -->
- **Variables de entorno**: <!-- ninguna / cuáles -->

## Rollback

<!-- "revert del merge" salvo que haya algo especial (migración, cambio de contrato). -->

## Checklist

- [ ] Sin secretos ni `.env` en el diff.
- [ ] Specs/tareas actualizadas si cambió el contrato.
- [ ] Si hay migración: es compatible con la versión anterior de la app.
