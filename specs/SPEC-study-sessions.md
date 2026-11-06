# SPEC — study-sessions (mobile)

| Campo | Valor |
|---|---|
| Module id | `study-sessions` |
| Origen | Necesidad real: preparar una entrevista esta semana (Java, Spring, Angular, microservicios) sin perder el SRS |
| Estado | `done` |
| Depende de | — (`POST /api/v1/sessions` ya existía) |
| Migración goose | Ninguna |
| Dependencia nueva | **Ninguna** |

## Objective

Que desde **mobile** se pueda crear una sesión de estudio a medida: elegir uno o varios temas,
una dificultad y cuántas preguntas (`10 / 25 / 50 / sin límite`), y empezar a responder de
inmediato. Hoy mobile solo puede lanzar el repaso diario (vencidas + warmup + tip) y el usuario
no tiene forma de decir «quiero 50 preguntas de Java ahora».

**Historias**
- Quiero un repaso intensivo de **un solo tema** (p. ej. React) con 50 preguntas.
- Quiero estudiar **sin límite** hasta aburrirme y finalizar la sesión yo mismo.
- Quiero mezclar **varios temas** (Java + Spring + Angular) en una sola sesión.
- Si pido más preguntas de las que existen, la app me ofrece continuar sin límite en vez de dejarme en un error.

**Éxito**: crear una sesión de estudio en ≤ 4 toques desde Inicio, sin cambios de backend.

## Contexto verificado

- `POST /api/v1/sessions` (`CreateSessionRequest`) ya acepta
  `{ name, mode: 'review', difficulty, topicIds, questionLimit? }`.
- `questionLimit` es `min=1,max=50`; **omitirlo = ilimitado** (`nil` en Go).
- `mode` solo admite `review` (no existe `generate` ni modos nuevos).
- Si el límite supera las preguntas disponibles, el backend responde **422** con
  `detail: "Solo hay N preguntas disponibles para estos temas, pero pediste M…"`.
- `review` selecciona preguntas **aleatorias** de esos temas (sin filtro SRS) y cada respuesta sí
  actualiza el SRS vía `progressService`.
- `GET /api/v1/topics?hasQuestions=true&search=` ya existe y pagina (máx. 100 por página).
- La pantalla de estudio ya tiene «Finalizar» y trata el 404 «No hay más preguntas» como fin.

## Project structure

```
apps/mobile/
├── app/(tabs)/sesiones/nueva.tsx                         # NUEVO (thin re-export)
└── src/
    ├── core/api/services/topics.service.ts               # NUEVO  listTopics()
    ├── core/api/services/topics.service.test.ts          # NUEVO
    └── features/
        ├── sessions/components/TopicPicker.tsx           # NUEVO  chips + modal + búsqueda
        ├── sessions/components/TopicPicker.test.tsx      # NUEVO
        ├── sessions/screens/StudySessionScreen.tsx       # NUEVO
        ├── sessions/screens/StudySessionScreen.test.tsx  # NUEVO
        ├── sessions/screens/SessionsListPage.tsx         # + botón ＋ en el header
        └── home/screens/HomeScreen.tsx                   # + CTA secundario
```

## Contrato (UI)

### `TopicPicker`

- Props: `selected: Topic[]`, `onChange(topics)`, `max?` (por defecto **8**).
- Chips de los temas elegidos con `×` (target ≥ 44 vía `hitSlop`) y botón «Añadir temas».
- `Modal` propio con buscador (`debounce` 400 ms), lista con `accessibilityRole="checkbox"`,
  `name` + `category` por fila (≥ 56 pt) y contador `n/8`.
- La selección es **borrador**: se aplica a `onChange` solo al pulsar «Listo».
- `listTopics({ hasQuestions: true, perPage: 50 })` — nunca ofrece temas sin preguntas.
- Al llegar al máximo, el resto de filas se deshabilita y se explica por qué.

### `StudySessionScreen`

| Control | Regla |
|---|---|
| Temas | Obligatorio ≥ 1; sin temas, «Empezar» deshabilitado |
| Dificultad | Rejilla 2×2 con `SESSION_DIFFICULTY_LABELS` (shared); por defecto `random` |
| Preguntas | `10 / 25 / 50 / Sin límite`; por defecto **25**; «Sin límite» ⇒ `questionLimit: undefined` |
| Nombre | Autogenerado `Estudio · Tema1, Tema2 +N`; deja de regenerarse al editarlo (`nameDirty`) |

- CTA «Empezar» (52 pt) → `createSession(...)` → `router.replace('/(tabs)/sesiones/{id}')`.
- Error 422 → mensaje del backend + botón **«Empezar sin límite»** (reintento con `questionLimit: undefined`).
- Otros errores → mensaje legible y el formulario sigue usable.

## Boundaries

- **Always**: reutilizar `POST /sessions` (no inventar endpoint ni modo); `hasQuestions=true`;
  tokens de diseño; targets ≥ 44; `accessibilityLabel` en todo control; sin dependencias nuevas.
- **Ask first**: subir el máximo de `questionLimit` (> 50); persistir temas recientes; añadir
  «solo vencidas» como filtro de la sesión de estudio.
- **Never**: crear un `mode` nuevo; duplicar tipos (`Topic`, `CreateSessionRequest` vienen de
  `@devtest/shared`); llamar a `apiClient` desde la pantalla; dejar el 422 como callejón sin salida.

## Testing

| Archivo | Casos |
|---|---|
| `topics.service.test.ts` | paginación por defecto + `hasQuestions`; reenvía `search`/`perPage` |
| `TopicPicker.test.tsx` | abre el modal y lista temas; selección en borrador hasta «Listo»; máximo + filas deshabilitadas; quitar chip; búsqueda con debounce |
| `StudySessionScreen.test.tsx` | CTA deshabilitado sin temas; payload exacto con 10; `questionLimit: undefined` con «Sin límite»; 422 → reintento sin límite; nombre editado se respeta; «Volver» |

## Success Criteria

1. Desde Inicio (`Sesión de estudio`) y desde el `＋` del listado se llega a la pantalla.
2. Se puede crear una sesión con 1–8 temas, cualquier dificultad y 10/25/50/∞.
3. La sesión abre directamente en la pantalla de estudio.
4. Un 422 por disponibilidad se resuelve dentro de la pantalla.
5. `typecheck` y `jest` en verde; sin targets < 44.

## Fuera de alcance (siguientes entregas)

- **Entrega 2 — Enfoque**: hasta 8 temas persistidos en servidor (`user_focus_topics`,
  `GET/PUT /me/focus`), editor en mobile, tarjeta en Inicio y pre-selección en `TopicPicker`.
- **Entrega 3 — Enfoque → diario**: `CreateDaily` guarda los temas en `session_topics`,
  `NextQuestion` de una `daily` usa los de **esa** sesión y las queries de vencidas/warmup aceptan
  `topicIDs` opcional. Estado vacío con CTA a sesión de estudio. Notificaciones con el mismo filtro.
