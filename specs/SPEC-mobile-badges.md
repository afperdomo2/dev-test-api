# SPEC — mobile-badges

| Campo | Valor |
|---|---|
| Module id | `mobile-badges` |
| Origen | Alcance nuevo mobile (D-26, D-29, D-33) |
| Estado | `ready` |
| Depende de | `badges` (backend: `GET /gamification/badges`) |
| Migración goose | Ninguna |
| Dependencia nueva | **Ninguna** |

## Objective

Que el usuario vea en el móvil sus **logros**: cuáles tiene, cuáles le faltan y cuánto le queda para el siguiente.

**Historias**:
- Desde el hub de Progreso entro en **Logros**.
- Veo los 8 logros; los conseguidos con su icono y la etiqueta "Conseguido"; los pendientes atenuados y con una barra `12/100`.
- Si todavía no tengo ninguno, veo un mensaje que me anima a empezar.

**Éxito**: el estado de cada logro coincide con `GET /gamification/badges` y no se promete nada que el backend no sepa (no hay fecha de consecución, D-16).

## Contexto

- `SPEC-badges.md` define el catálogo y el endpoint. **Aún no implementados** (`BD-1…BD-5`); esta spec se ejecuta en la Ola 3, después.
- El DTO `Badge` / `BadgesResponse` se define en `@devtest/shared` en la spec de badges. **No se duplica aquí**.
- La pestaña y el hub los crea `mobile-progress`.

## Commands

```bash
pnpm --filter mobile typecheck
pnpm --filter mobile test:run
pnpm --filter mobile lint
```

## Project structure

```
apps/mobile/app/(tabs)/progreso/logros.tsx            # NUEVO (thin re-export)
apps/mobile/src/features/progress/
├── screens/AchievementsScreen.tsx                     # NUEVO
├── components/BadgeCard.tsx                           # NUEVO
└── screens/AchievementsScreen.test.tsx                # NUEVO
apps/mobile/src/features/progress/screens/ProgressHubScreen.tsx  # + fila "Logros"
apps/mobile/src/core/api/services/gamification.service.ts        # + getBadges
apps/mobile/src/core/api/services/gamification.service.test.ts   # + caso
```

## Contrato

### Servicio

```ts
// src/core/api/services/gamification.service.ts
export async function getBadges(): Promise<BadgesResponse> {
  return apiClient.get<BadgesResponse>('/api/v1/gamification/badges')
}
```

(`BadgesResponse` = `{ badges: Badge[]; earned: number; total: number }`, definido en shared.)

### Pantalla

- Cabecera: **`earned` de `total`** ("3 de 8 logros") con una barra de progreso global.
- **Cuadrícula de 2 columnas** (`FlatList numColumns={2}`) con `BadgeCard`.
- `BadgeCard`:
  - icono del catálogo (`icon` viene como nombre `mdi-*`; se mapea a `MaterialCommunityIcons`),
  - nombre y descripción,
  - si `earned` → etiqueta **"Conseguido"** con color y texto,
  - si no → barra `current/target` + texto `current/target`,
  - `accessibilityLabel` completo: `"Primera sesión. Conseguido"` o `"100 respuestas. 12 de 100"`.
- **Sin `earnedAt`**: no existe (D-16). No se muestra fecha ni se ordena por ella.
- Estado vacío: si `badges.length === 0` → mensaje (el backend siempre devuelve 8, pero se cubre).
- Estados: loading (skeletons), error (+ "Reintentar").

### Mapeo de iconos

El backend manda `mdi-flag-checkered`, `mdi-fire`, etc. Se define un helper local:

```ts
// mdi-flag-checkered -> "flag-checkered" (MaterialCommunityIcons)
function toMcIcon(name: string): string { return name.replace(/^mdi-/, '') }
```

Si el nombre no existe en la fuente, el icono cae a un fallback (`trophy-outline`) y **no** rompe la pantalla. (Los nombres reales se fijan en `SPEC-badges.md`; validarlos allí contra el catálogo de MaterialCommunityIcons.)

## Testing strategy

| Archivo | Casos |
|---|---|
| `gamification.service.test.ts` | `getBadges` URL y respuesta |
| `AchievementsScreen.test.tsx` | pinta los 8; contador `earned/total`; un conseguido muestra "Conseguido"; un pendiente muestra `current/target`; estado vacío; error con reintento |
| `BadgeCard.test.tsx` | `accessibilityLabel` correcto en ambos estados; el icono se mapea |

## Boundaries

- **Always**: consumir el endpoint; tipos de shared; iconos con texto (nunca solo color); tokens de diseño.
- **Ask first**: persistir logros en el dispositivo; ordenar por fecha; añadir animaciones de desbloqueo.
- **Never**: inventar `earnedAt`; crear una tabla de logros (D-16); duplicar el catálogo de ids en mobile (viene del backend); recalcular criterios en el cliente.

## Success Criteria

1. Desde el hub se abre **Logros** y se listan los 8.
2. El contador coincide con `earned`/`total`.
3. Un logro conseguido se distingue por **icono + etiqueta**, no solo por color.
4. Un pendiente muestra su progreso `current/target`.
5. No se muestra ninguna fecha de consecución.
6. `typecheck`, `test:run` y `lint` en verde.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Cuadrícula o lista? | **Cuadrícula de 2 columnas**: los logros son visuales y caben bien. |
| OQ-2 | ¿Agrupar por "conseguidos / pendientes"? | No en la v1: el orden del backend ya es estable y mezclar permite ver lo que falta. |
| OQ-3 | ¿Animación al abrir? | No; primero que funcione. Se puede añadir después sin cambiar el contrato. |
