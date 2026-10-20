# SPEC — mobile-progress

| Campo | Valor |
|---|---|
| Module id | `mobile-progress` |
| Origen | Alcance nuevo mobile (D-26, D-33) |
| Estado | `ready` |
| Depende de | `review-mastery` ✅ (endpoint ya hecho) |
| Consumido por | `mobile-weekly` y `mobile-badges` (reutilizan la pestaña y el hub) |
| Migración goose | Ninguna |
| Dependencia nueva | **Ninguna** |

## Objective

Dar acceso desde la app al **dominio por tema** que ya existe en la web, dentro de una nueva pestaña **Progreso**.

**Historias**:
- Como usuario, veo una 4ª pestaña **Progreso** en la barra inferior.
- Al entrar, veo un índice con **Temas** (y, más adelante, Mi semana y Logros).
- En **Temas** veo una barra por tema con mi porcentaje, mis aciertos/respuestas y el aviso **"Débil"** cuando corresponde.
- Si aún no he respondido nada, veo un mensaje claro, no una pantalla vacía.

**Éxito**: el usuario puede saber **qué repasar** desde el móvil, con los mismos números que la web (mismo endpoint, misma definición de "débil").

## Contexto verificado

- `NativeTabs` hoy tiene **3** triggers: `index`, `sesiones`, `perfil` (`app/(tabs)/_layout.tsx`).
- mobile ya tiene `src/core/api/services/progress.service.ts` con `getUpcomingReviews`, `getSavedQuestions`, `toggleSaveQuestion` — **falta** `getTopicMastery`.
- `@devtest/shared` ya exporta `TopicMasteryItem` (hecho en RM-6).
- `GET /progress/topics-mastery` devuelve `{ data: TopicMasteryItem[] }`; el cliente mobile desenvuelve el envelope.

## Tech stack

Expo Router (`NativeTabs` + `Stack`) · React Native · tokens de `src/core/theme/tokens.ts` · `@devtest/shared`.

## Commands

```bash
pnpm --filter mobile typecheck
pnpm --filter mobile test:run
pnpm --filter mobile lint
pnpm --filter mobile start   # verificación manual en dispositivo
```

## Project structure

```
apps/mobile/
├── app/(tabs)/_layout.tsx                  # + trigger "progreso"
├── app/(tabs)/progreso/
│   ├── _layout.tsx                         # NUEVO: Stack (header + back)
│   ├── index.tsx                           # NUEVO: hub (thin re-export)
│   └── temas.tsx                           # NUEVO: thin re-export
└── src/features/progress/                  # NUEVO dominio
    ├── screens/ProgressHubScreen.tsx
    ├── screens/TopicMasteryScreen.tsx
    ├── components/TopicMasteryRow.tsx
    ├── components/ProgressHubRow.tsx
    └── *.test.tsx

apps/mobile/src/core/api/services/progress.service.ts   # + getTopicMastery
apps/mobile/src/core/api/services/progress.service.test.ts # + caso
```

## Contrato

### Servicio

```ts
// src/core/api/services/progress.service.ts
export async function getTopicMastery(): Promise<Array<TopicMasteryItem>> {
  return apiClient.get<Array<TopicMasteryItem>>('/api/v1/progress/topics-mastery')
}
```

Sin parámetros y sin paginación (el backend devuelve todos los temas con respuestas).

### Navegación

```tsx
// app/(tabs)/_layout.tsx
<NativeTabs.Trigger name="progreso">
  <NativeTabs.Trigger.Icon sf="chart.line.uptrend.xyaxis" md="insights" />
  <NativeTabs.Trigger.Label>Progreso</NativeTabs.Trigger.Label>
</NativeTabs.Trigger>
```

Orden final: **Inicio · Sesiones · Progreso · Perfil**.

`progreso/_layout.tsx` monta un `<Stack screenOptions={{ headerShown: true }}>` para que `temas` tenga botón atrás. El hub (`index`) oculta el header (lo pinta la propia pantalla).

### Hub (`ProgressHubScreen`)

Lista vertical de filas (`ProgressHubRow`): icono + título + subtítulo + chevron.

| Fila | Estado en esta ola | Navega a |
|---|---|---|
| **Temas** | ✅ se crea ahora | `/progreso/temas` |
| **Mi semana** | ⬜ se añade en `mobile-weekly` | `/progreso/semana` |
| **Logros** | ⬜ se añade en `mobile-badges` | `/progreso/logros` |

> El hub se crea con **solo la fila de Temas**. Las otras dos filas las añaden sus specs, para que cada ola deje la app coherente (nada de enlaces a pantallas que no existen).

### Pantalla `TopicMasteryScreen`

- `useEffect` → `getTopicMastery()`; estados `loading | error | empty | data`.
- **Lista** (`FlatList` o `ScrollView`): una `TopicMasteryRow` por tema.
- Cada fila:
  - nombre del tema,
  - barra horizontal: `View` con `width: accuracy * 100 + '%'` (sin librería de gráficas),
  - texto `correctCount/attempts`,
  - chip **"Débil"** cuando `isWeak` (color + **texto**, nunca solo color),
  - `accessibilityLabel` con el porcentaje, p. ej. `"Go: 40% de acierto"`.
- **Vacío**: mensaje "Aún no hay respuestas registradas. Responde preguntas para ver tu dominio por tema."
- **Error**: mensaje + botón "Reintentar".
- **Refresh**: `RefreshControl` que vuelve a llamar al servicio.
- Orden: el que devuelve el backend (débiles primero); la app **no** reordena.

### Tokens y a11y

- Colores: `colors.primary` (barra normal), `colors.warning` (barra y chip de débil), `colors.text`, `colors.textMuted`, `colors.surface`, `colors.border`.
- Touch targets ≥ 44 pt en las filas del hub y en el chip (si es pulsable, no lo es).
- `accessibilityRole="button"` en las filas del hub; `accessibilityLabel` en las barras.
- Los iconos son `MaterialCommunityIcons`; nada de emojis.

## Code style

- `app/` contiene **solo** rutas thin (re-export) y layouts (regla del repo).
- Pantallas con la lógica; componentes presentacionales.
- Sin `any`; tipos desde `@devtest/shared`.

## Testing strategy

| Archivo | Casos |
|---|---|
| `progress.service.test.ts` | `getTopicMastery` llama la URL correcta y devuelve el array |
| `TopicMasteryScreen.test.tsx` | pinta una fila por tema; chip solo en `isWeak`; estado vacío; estado de error con "Reintentar"; el refresh vuelve a llamar al servicio |
| `ProgressHubScreen.test.tsx` | pinta la fila Temas y navega a `/progreso/temas` |
| `TopicMasteryRow.test.tsx` | porcentaje correcto en el ancho de la barra y en el `accessibilityLabel` |
| `(tabs)/_layout` | (manual) las 4 pestañas aparecen y las 3 existentes siguen funcionando |

## Boundaries

- **Always**: endpoint existente (no crear backend); tokens de diseño; textos de estado además del color; reexport thin en `app/`; probar en dispositivo que la barra inferior no se rompe.
- **Ask first**: añadir una 5ª pestaña; meter gráficas o una librería de charts; paginar el endpoint; añadir filtros por tema.
- **Never**: duplicar `TopicMasteryItem`; recalcular `isWeak` en el cliente (el backend es la única fuente, D-08); ordenar distinto al backend; hex hardcodeados; tocar web o backend.

## Success Criteria

1. La barra inferior muestra **4** pestañas y las 3 previas siguen funcionando.
2. `Progreso` abre un hub con la fila **Temas**.
3. **Temas** lista los mismos temas, con los mismos números, que `GET /progress/topics-mastery`.
4. Un tema `isWeak: true` muestra el chip "Débil" con texto; uno `false` no lo muestra.
5. Usuario sin respuestas → estado vacío con el mensaje; nunca una pantalla en blanco.
6. Fallo de red → mensaje + "Reintentar" que reintenta de verdad.
7. `typecheck`, `test:run` y `lint` en verde.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿La pantalla de Temas muestra también los temas con 0 respuestas? | **No**: el backend no los devuelve y el vacío global ya se explica. |
| OQ-2 | ¿"Progreso" debería incluir también los repasos pendientes (hoy solo en sesiones)? | **No** en esta spec: `mobile-progress` es dominio por tema. Se puede abrir otra spec si se quiere. |
| OQ-3 | ¿Pull-to-refresh o botón? | **Pull-to-refresh** (nativo y sin ocupar pantalla). |
