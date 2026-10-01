---
description: Review de código contra las convenciones del proyecto
subagent: true
---
Área o archivos (opcional): $ARGUMENTS
!`git diff --staged --stat`
!`git diff --staged`

Si no hay diff staged, revisa los archivos indicados en $ARGUMENTS.
Evalúa contra las reglas del proyecto (lee los archivos; no dupliques sus reglas):
- Backend: `.agents/backend/responses.md` (envelope RFC 9457, error language,
  camelCase) + `.agents/backend/architecture.md` (anatomía de módulo, wiring).
- Web: `.agents/web/patterns.md` (TanStack Query, debounce, Vue/TS) +
  `.agents/web/api-client.md` (interceptores, unwrap del envelope).
- Mobile: `.agents/mobile/architecture.md` (UX conventions) +
  `.agents/mobile/design-system.md` (tokens, iconos, touch targets, dark mode).
- Shared: `.agents/shared/architecture.md` (qué va en shared vs app; sin APIs
  de plataforma).

Devuelve hallazgos ordenados por severidad con `archivo:línea`, la regla
violada y una sugerencia concreta. No reportes nada sin regla aplicable.
