# SPEC — mobile-daily-tip

| Campo | Valor |
|---|---|
| Module id | `mobile-daily-tip` |
| Origen | Alcance nuevo mobile (D-26, D-30) |
| Estado | `ready` |
| Depende de | — (el backend ya lo envía en `POST /sessions/daily`) |
| Migración goose | Ninguna |
| Dependencia nueva | **Ninguna** |

## Objective

Mostrar el **tip del día** al empezar la sesión diaria. Hoy el backend lo devuelve y la app lo descarta: `handleDailySession` navega a la sesión sin mirar `res.tip`.

**Historias**:
- Al pulsar "Repasar hoy", la sesión diaria arranca mostrando el tip de un tema, arriba de la primera pregunta.
- Si el tip molesta, se puede cerrar y seguir respondiendo sin que reaparezca en esa sesión.
- Si el usuario no tiene tip ese día (o la sesión no es diaria), no aparece ninguna tarjeta y todo funciona igual.

**Éxito**: el tip que ya viaja en la respuesta se muestra **sin llamadas nuevas** al backend y sin efectos secundarios.

## Contexto verificado

```ts
// @devtest/shared — DailySessionResponse
{ session: Session; tip?: Tip; dueCount: number; warmupCount: number }
```

- El tip **solo** existe en la respuesta de `POST /sessions/daily`. No hay `GET` de tip, y llamar al POST solo para pintarlo crearía una sesión (efecto secundario). Por eso D-30: se muestra al entrar a la sesión.
- `Tip` trae `content` y `topic`.
- `HomeScreen.handleDailySession()` ya hace `createDailySession()` y navega: ese es el único punto que hay que tocar para no perderlo.

## Tech stack

React Native + Expo Router · `MarkdownText` (ya existe en `src/core/components`) · design tokens.

## Commands

```bash
pnpm --filter mobile test:run
pnpm --filter mobile typecheck
pnpm --filter mobile lint
```

## Project structure

```
apps/mobile/src/features/sessions/
├── dailyTipStore.ts                  # NUEVO: memoria sessionId -> Tip
├── dailyTipStore.test.ts             # NUEVO
├── components/DailyTipCard.tsx       # NUEVO
├── components/DailyTipCard.test.tsx  # NUEVO
└── screens/SessionStudyPage.tsx      # + leer y pintar el tip

apps/mobile/src/features/home/screens/HomeScreen.tsx    # + guardar el tip antes de navegar
apps/mobile/src/features/home/screens/HomeScreen.test.tsx # + caso
```

## Contrato

### Store en memoria (`dailyTipStore.ts`)

```ts
import type { Tip } from '@devtest/shared'

let pending: { sessionId: string; tip: Tip } | null = null

/** Guarda el tip que vino con la sesión diaria. No persiste. */
export function setPendingTip(sessionId: string, tip?: Tip): void

/** Devuelve el tip de esa sesión, o null. NO lo consume (se puede releer). */
export function getPendingTip(sessionId: string): Tip | null

/** Limpia el tip pendiente (al salir de la sesión o al crear otra). */
export function clearPendingTip(): void
```

**Por qué memoria y no router params**: el contenido de un tip es texto largo; meterlo en la URL de Expo Router obliga a serializarlo y ensucia el historial. Un store de 3 funciones es más simple y no persiste nada sensible.

**Límite asumido**: si el proceso de la app muere y se reabre a mitad de la sesión diaria, el tip ya no se muestra. Es aceptable (D-30) y queda documentado.

### `HomeScreen`

```ts
const res = await createDailySession()
setPendingTip(res.session.id, res.tip)   // ← única línea nueva
router.push(`/(tabs)/sesiones/${res.session.id}`)
```

Si `res.tip` es `undefined`, el store no guarda nada y la tarjeta no aparece.

### `SessionStudyPage`

- Al montar (con `sessionId` disponible), `const tip = getPendingTip(sessionId)`.
- Si hay tip → `<DailyTipCard tip={tip} onDismiss={...} />` **arriba del `QuestionCard`**, dentro del contenido de la fase `question`.
- Al desmontar la pantalla o al terminar la sesión → `clearPendingTip()`.
- La tarjeta **no** aparece en las fases `result`, `finished` ni `error` (para no tapar el feedback ni el resumen).

### `DailyTipCard`

- Título "Tip del día" + tema (`tip.topic`) como etiqueta.
- Contenido con `MarkdownText` (el tip puede traer markdown).
- Botón de cerrar con `accessibilityLabel="Cerrar tip del día"` y touch target ≥ 44 pt.
- Estilo con tokens: fondo `colors.surface`, borde `colors.border`, texto `colors.text`. **Sin hex hardcodeados**.
- El estado "cerrado" es local (`useState`); no se persiste.

## Code style

- Componentes presentacionales: reciben `tip` y emiten `onDismiss`; no leen el store.
- `SessionStudyPage` es la que habla con el store (igual que ya hace con los servicios).
- Sin emojis como icono: `MaterialCommunityIcons` para el cierre (`close`).
- Nombres en el idioma del módulo: `DailyTipCard`, `setPendingTip`.

## Testing strategy

| Archivo | Casos |
|---|---|
| `dailyTipStore.test.ts` | guardar y leer por sessionId; `get` no consume (dos lecturas seguidas); otra sesión devuelve `null`; `clear` limpia; `set` con tip `undefined` no guarda |
| `DailyTipCard.test.tsx` | pinta título, tema y contenido; el botón llama `onDismiss`; tiene `accessibilityLabel` |
| `HomeScreen.test.tsx` | tras `createDailySession` con tip, el store queda cargado (mock del store o del servicio) |
| `SessionStudyPage.test.tsx` | con tip pendiente pinta la tarjeta; sin tip no la pinta; al cerrar desaparece; en fase `result` no aparece |

## Boundaries

- **Always**: el tip es opcional (sin tip, todo funciona); usar `MarkdownText`; tokens de diseño; touch target ≥ 44 pt; limpiar el store al salir.
- **Ask first**: persistir el tip (AsyncStorage/SecureStore); añadir un `GET /tips/today` al backend; mostrarlo también en Home; mostrarlo en sesiones que no sean `daily`.
- **Never**: llamar a `POST /sessions/daily` solo para leer el tip; meter el texto del tip en params de router; mostrar la tarjeta durante el feedback de una respuesta; duplicar tipos de `Tip` (viene de `@devtest/shared`).

## Success Criteria

1. Al crear una sesión diaria con tip, la pantalla de la sesión muestra la tarjeta **antes de la primera pregunta**.
2. Sin tip, la pantalla es idéntica a hoy.
3. Cerrar la tarjeta la oculta en esa sesión y no reaparece al responder ni al pasar a la siguiente pregunta.
4. La tarjeta no se pinta en las fases `result`, `finished` ni `error`.
5. `getPendingTip` se puede llamar dos veces con el mismo resultado.
6. `test:run`, `typecheck` y `lint` en verde.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Persistir el tip para sobrevivir a un reinicio? | **No** (D-30). Un tip es efímero; persistirlo exige AsyncStorage y limpieza. |
| OQ-2 | ¿Mostrar también `dueCount`/`warmupCount` en la tarjeta? | **No** en esta spec: la tarjeta es del tip. Si se quiere contexto, va en otra tarjeta de la Home. |
| OQ-3 | ¿El botón de cerrar o un "entendido"? | **Icono de cerrar** con `accessibilityLabel`; patrón habitual en tarjetas. |
