# Tasks — study-sessions (mobile)

> Spec: `specs/SPEC-study-sessions.md` · Estado: **done**

## Tareas

- [x] Servicio `listTopics()` con `hasQuestions`, `search`, paginación.
      **Hecho:** `src/core/api/services/topics.service.ts` (+ test con los params exactos).
- [x] `TopicPicker`: chips + `Modal` con buscador (debounce 400 ms desde el handler), máximo 8,
      filas ≥ 56 pt con `accessibilityRole="checkbox"`, confirmación con «Listo».
      **Hecho:** `features/sessions/components/TopicPicker.tsx` (+ 5 tests).
- [x] `StudySessionScreen`: temas, dificultad 2×2, cantidad `10/25/50/Sin límite`, nombre
      autogenerado editable, CTA y manejo del 422 con reintento sin límite.
      **Hecho:** `features/sessions/screens/StudySessionScreen.tsx` (+ 6 tests).
- [x] Ruta `app/(tabs)/sesiones/nueva.tsx` (thin re-export).
      **Hecho:** archivo de 1 línea.
- [x] Entradas: CTA secundario en Inicio y botón `＋` en el header del listado.
      **Hecho:** `HomeScreen.tsx` (`Sesión de estudio`) y `SessionsListPage.tsx` (`Nueva sesión de estudio`).
- [x] Docs: rutas y convenciones de formularios de estudio en `.agents/mobile/architecture.md`;
      línea de rutas en `apps/mobile/AGENTS.md`.
      **Hecho:** sección «Formularios de estudio» + árbol de rutas actualizado.
- [x] Verificación.
      **Hecho:** `pnpm --filter mobile typecheck` limpio; jest relacionados
      (`TopicPicker`, `StudySessionScreen`, `topics.service`, `HomeScreen`, `SessionsListPage`) en verde.

## Pendiente (entregas siguientes)

- [ ] **Enfoque** persistido en servidor + editor + tarjeta en Inicio + prefill del picker.
- [ ] **Enfoque → repaso diario**: snapshot de temas en la sesión `daily` y filtro opcional en
      `FindNextDailyQuestion` / `CountDueReviews` / `CountWarmupCandidates`.
- [ ] Notificaciones respetando el enfoque.
- [ ] Atajo «Usar mis temas flojos como enfoque» desde Progreso → Temas.
