---
description: Genera un mensaje de commit convencional
subagent: true
---
Cambios staged:
!`git diff --staged --stat`
!`git diff --staged`

Genera un mensaje con Conventional Commits (`type(scope): subject`).
- Scope: dominio/feature afectado (ej. `auth`, `sessions`, `ai`); sin scope si es transversal.
- Idioma: subject en inglés, en minúsculas después del tipo, sin punto final.
- Solo devuelve el mensaje. Nunca ejecutes `git commit` — los commits son manuales en este repo.
