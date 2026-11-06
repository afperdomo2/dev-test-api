# SPEC — mobile-progress-ux

| Campo | Valor |
|---|---|
| Module id | `mobile-progress-ux` |
| Origen | Revisión de UX de Progreso + Inicio (2026-10-25) |
| Estado | `ready` |
| Depende de | `mobile-progress` ✅, `mobile-weekly` ✅, `mobile-badges` ✅ |
| Consumido por | — |
| Migración goose | Ninguna |
| Dependencia nueva | **`react-native-svg`** (solo mobile; incluida en Expo Go) |

## Objective

Que la pestaña **Progreso** muestre datos (no un menú) y que **Inicio** responda solo a
"¿qué hago ahora?", sin repetir la información de Progreso.

**Problemas medidos antes del cambio**:

| Vista | Antes |
|---|---|
| Progreso (hub) | 3 filas de navegación, 0 datos, 2 taps para el primer dato |
| Temas | ~100 px por tema con card + barra → 10 temas ≈ 1,3 pantallas |
| Mi semana | ~1.100 px: 6 KPI cards + 2 listas de cards de 92 px + chart 120 px; KPI "Racha" no es del periodo |
| Inicio | 4 secciones (~600-700 px), 2 de ellas duplican Progreso (nivel y XP por tema) |

**Éxito**: el primer dato aparece al abrir la pestaña; Inicio cabe sin scroll en un
teléfono estándar; ningún dato se repite con el mismo peso en dos pantallas.

## Decisiones de esta spec

| # | Decisión | Motivo |
|---|---|---|
| D1 | **Dashboard seccionado**, no hub de botones | Un hub es un menú de ajustes: añade un tap y no informa. Las pantallas de detalle siguen existiendo para "el todo" (lista completa, filtros, mapa de calor) |
| D2 | **Anillos** para temas (estilo Duolingo) | 3 columnas × ~100 px: 9-12 temas en media pantalla (antes 5-6) |
| D3 | `react-native-svg` para el anillo | `strokeDasharray` da un círculo nítido; es la vía estándar y no añade librería de charts |
| D4 | **Actividad de 7 días en Inicio** | La racha se entiende con los días; es *momentum*, no análisis. El heatmap mensual y los KPIs quedan en Progreso → Actividad |
| D5 | `semana` → **`actividad`** | La pantalla ya cubre semana y mes; "Mi semana" se quedaba corto |
| D6 | **Anillo = estado global por tema · barra = ranking del periodo** | Dos lenguajes con semántica distinta, nunca mezclados al azar |

**Regla de propiedad** (evita volver a duplicar):

| Nivel | Pregunta | Contenido |
|---|---|---|
| Inicio | ¿qué hago ahora? | CTA + hoy (racha, 7 días, meta) + teaser de logros |
| Progreso | ¿cómo voy? | nivel, anillos de temas, resumen de actividad, logros |
| Detalle | dame el todo | lista completa, filtros, heatmap mensual |

## Project structure

```
apps/mobile/
├── app/(tabs)/progreso/
│   ├── _layout.tsx                  # títulos: Progreso · Temas · Actividad · Logros
│   ├── index.tsx                    # → ProgressDashboardScreen
│   ├── temas.tsx                    # → TopicMasteryScreen
│   ├── actividad.tsx                # → ActivityScreen (antes semana.tsx)
│   └── logros.tsx                   # → AchievementsScreen
└── src/
    ├── core/
    │   ├── components/LevelLine.tsx          # NUEVO: nivel + barra 4px (Inicio y Progreso)
    │   ├── components/SectionHeader.tsx      # NUEVO: título + "Ver todos" (compartido)
    │   ├── components/BadgeStrip.tsx         # NUEVO: iconos de logros (compacto | detailed)
    │   ├── theme/charts.ts                   # NUEVO: HEAT_COLORS + heatLevel()
    │   └── utils/
    │       ├── icons.ts                      # NUEVO: toMcIcon (movido desde BadgeCard)
    │       └── numbers.ts                    # NUEVO: toPercent + formatDurationMs
    ├── features/home/
    │   ├── screens/HomeScreen.tsx            # reescrito (slim)
    │   └── components/
    │       ├── TodayCard.tsx                 # NUEVO: racha + 7 días + meta + freeze
    │       └── DayStreakStrip.tsx            # NUEVO: 7 casillas
    └── features/progress/
        ├── screens/
        │   ├── ProgressDashboardScreen.tsx   # NUEVO (ex ProgressHubScreen)
        │   ├── TopicMasteryScreen.tsx        # reescrito: rejilla de anillos + filtro
        │   ├── ActivityScreen.tsx            # NUEVO (ex WeeklySummaryScreen)
        │   └── AchievementsScreen.tsx        # reescrito: compacto + agrupado
        └── components/
            ├── TopicRing.tsx                 # NUEVO (react-native-svg)
            ├── TopicRingGrid.tsx             # NUEVO (3 columnas)
            ├── TopicCompactRow.tsx           # NUEVO (barra, ~46 px)
            ├── ActivityHeatmap.tsx           # NUEVO (30 días, View puro)
            ├── ActivitySparkline.tsx         # NUEVO (7 barras, 28 px)
            ├── StatRow.tsx                   # NUEVO (métricas con divisores)
            └── BadgeCard.tsx                 # 172 → 140 px, importa toMcIcon de core
```

**Borrados**: `ProgressHubScreen`, `ProgressHubRow`, `WeeklySummaryScreen`, `KpiCard`,
`TopicMasteryRow`, `AchievementsTeaser` (+ sus tests).

## Contrato

Sin cambios de backend. Servicios existentes:

```ts
getOverview()                        // GET /gamification/overview
getBadges()                          // GET /gamification/badges
getTopicMastery()                    // GET /progress/topics-mastery
getInsightSummary('week' | 'month')  // GET /insights/summary
```

`InsightSummary.streakDays` **no se usa** en la UI: es la racha actual, no la del periodo
(la racha vive en Inicio).

### Dashboard (`ProgressDashboardScreen`)

| Sección | Datos | Detalle |
|---|---|---|
| Nivel | `overview` | `LevelLine` (Nv + XP que falta + barra 4px) |
| Temas | `topics-mastery` | 6 anillos + "y N temas más" si hay más |
| Actividad | `insights?period=week` | `N preguntas · % acierto · tiempo` + sparkline |
| Logros | `badges` | hasta 5 iconos + `earned de total` + siguiente + barra |

- **Cada sección carga y falla por su cuenta** (hook local `useSection`): un endpoint caído
  no vacía la pestaña; el error se muestra en su bloque con `Reintentar <sección>`.
- Pull-to-refresh recarga las 4 en paralelo (`Promise.allSettled`).
- Cabeceras de sección con acción "Ver todos" / "Ver detalle" (`accessibilityLabel`
  `"Ver todos temas"`, `"Ver detalle actividad"`, `"Ver todos logros"`).

### Temas

- `TopicRingGrid` de 3 columnas; el backend devuelve los temas ordenados (débiles primero)
  y **la app no reordena**.
- Filtro `Todos / Débiles` (chips compactos 36 px + `hitSlop`); resumen "X de Y temas flojos".
- `TopicRing`: anillo **64 px**, `strokeWidth` 6, `rowGap` de 24 px entre filas (el aire separa,
  sin líneas de tabla), % centrado (escrito, nunca solo color), `accessibilityValue
  {min:0,max:100,now}` y etiqueta `"Go: 82% de acierto, 8 de 10 aciertos"` / `"…, tema débil"`.
- Tema débil: arco y % en `accent`, **sin icono** (decisión de producto: con el color
  y el % ya se entiende; "tema débil" sigue en la etiqueta accesible y en el filtro Débiles).

### Actividad

| Bloque | Contenido |
|---|---|
| Periodo | `PeriodToggle` (Semana / Mes) |
| Hero | `N preguntas` + `en M sesiones` |
| Métricas | `StatRow`: acierto · tiempo · XP ganado |
| Actividad | Semana → `DailyActivityBar`; Mes → `ActivityHeatmap` |
| Temas | `TopicCompactRow` (máx. 5 "Más practicados", máx. 3 "Puntos a reforzar") |
| Vacío | "No hay actividad en este periodo." (sin secciones vacías) |

`ActivityHeatmap`: 30 celdas de 18 px alineadas por día de semana (lunes primero, días en
UTC), 4 niveles de `HEAT_COLORS` (`colors.primary` al 20/50/100 % sobre `colors.muted`),
etiqueta por celda `"lun 19: 8 preguntas"` y leyenda Menos/Más.

### Logros

- `BadgeStrip` (en `core`, compartido): iconos de 28 px ordenados (conseguidos primero y, en los
  pendientes, el más avanzado), chip `+N` y contador `earned de total`.
  - **Inicio**: modo compacto (`max={4}`, sin barra ni "Siguiente") bajo el header
    "LOGROS — Ver todos →". Es una colección, no un resumen: el detalle vive en Progreso.
  - **Progreso**: `detailed` → añade `· Siguiente: X (2/5)` y la barra de progreso.
- Detalle: cards de 140 px agrupadas en **Conseguidos** y **En progreso** (orden por % desc).
  Es orden presentacional; no se recalcula `earned`.

### Inicio

- `LevelLine` sustituye el badge "Nv. X" **y** la card "Progreso" (un solo dato, una vez).
- `TodayCard`: racha + mejor + 7 casillas (pulsables → Actividad) + meta diaria + freeze.
- Sección **Logros** con `SectionHeader` + `BadgeStrip` compacto → Logros.
- **Eliminadas** las cards "Progreso" (nivel) y "Temas" (`overview.topicXp`): el XP por tema
  vive en Progreso.
- Si `overview` falla: fila "No se pudo cargar tu progreso · Reintentar" y el CTA sigue
  disponible (la acción no depende del resumen).
- Si `insights` o `badges` fallan: el bloque desaparece en silencio (no se ensucia Inicio).

## Tokens y a11y

- Todo sale de `colors/spacing/radius/typography`; el heatmap de `core/theme/charts.ts`.
- Touch targets ≥ 44 pt en CTA, freeze, "Ver todos", casillas de la semana (contenedor) y chips.
- `Pressable` con feedback `opacity 0.7-0.85`; `accessibilityRole`/`accessibilityState` en
  chips, acciones y freeze.
- Estado = icono + color + texto (anillo con %, chip "Débil", "Conseguido", "Alcanzada").
- Inicio: la tira de días no expone 7 nodos al lector; el contenedor anuncia
  "Ver tu actividad: practicaste N de los últimos 7 días".

## Testing strategy

| Archivo | Casos |
|---|---|
| `TopicRing.test.tsx` | % y nombre, `accessibilityValue`, etiqueta débil, etiqueta sin conteos |
| `TopicCompactRow.test.tsx` | nombre + `N resp · %`, chip Débil, etiqueta con aciertos |
| `ActivityHeatmap.test.tsx` | una celda etiquetada por día; sin días no pinta |
| `ActivitySparkline.tsx` | (cubierto por el dashboard) |
| `DayStreakStrip.test.tsx` | rellena solo días con actividad; sin días no pinta |
| `TodayCard.test.tsx` | racha/meta, tira pulsable → Actividad, freeze deshabilitado, sin semana |
| `AchievementsTeaser.test.tsx` | (sustituido por `BadgeStrip.test.tsx`: compacto con `+N`, detallado con "Siguiente" + barra, sin catálogo) |
| `ProgressDashboardScreen.test.tsx` | 4 secciones, navegación de las 3 acciones, fallo aislado de temas, vacíos |
| `TopicMasteryScreen.test.tsx` | anillos + resumen, filtro Débiles, sin débiles, vacío, error |
| `ActivityScreen.test.tsx` | hero + métricas + temas + barras; Mes → heatmap; vacío; error |
| `AchievementsScreen.test.tsx` | contador, ambos estados, vacío, error |
| `HomeScreen.test.tsx` | nivel, CTA, "Tu día", logros, sin sección Temas, tira → Actividad, reintento de overview |
| `numbers.test.ts` / `charts.test.ts` | `toPercent`, `formatDurationMs`, tramos de `heatLevel` |

**Verificación**: `pnpm --filter mobile typecheck` · `pnpm --filter mobile test:run` ·
`pnpm --filter mobile lint` · pasada manual en dispositivo (los anillos son SVG).

## Boundaries

- **Always**: consumir los endpoints existentes (cero backend); anillo con % escrito;
  reutilizar `TopicCompactRow` en Actividad; tokens; estados vacío/error por sección.
- **Ask first**: añadir otra dependencia de gráficas; una 5ª pestaña; tocar el contrato de
  `insights` (p. ej. `streakDays`); mover KPIs de periodo a Inicio.
- **Never**: recalcular `isWeak`/`earned` (D-08, D-16); duplicar el mismo dato con el mismo
  peso en Inicio y Progreso; reordenar los temas (los ordena el backend); hex fuera de tokens;
  volver al hub de botones.

## Success Criteria

1. Progreso abre con **al menos un dato visible** (nivel + anillos) sin taps extra.
2. 6 anillos + "Ver todos" en el dashboard; 9-12 temas visibles sin scroll en el detalle.
3. Un endpoint caído deja el resto de secciones usable y ofrece reintento propio.
4. Inicio cabe sin scroll en un teléfono de 6,1" y no muestra la sección "Temas".
5. La tira de 7 días lleva a Progreso → Actividad.
6. Actividad alterna Semana (barras) / Mes (heatmap de 30 días) y no muestra la racha.
7. Logros se agrupa en Conseguidos / En progreso sin recalcular nada.
8. `typecheck`, `test:run` (158 tests / 36 archivos) y `lint` (0 errores) en verde.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Animar el arco del anillo al montar? | **Sí, en una segunda pasada** (reanimated + `useAnimatedProps`); ahora es estático para no meter riesgo en tests |
| OQ-2 | ¿Pantalla de detalle por tema (preguntas del tema)? | **No** en esta spec: el anillo es estado, no navegación |
| OQ-3 | ¿Caché de `insights/week` en Inicio? | **No** por ahora: las pestañas quedan montadas y el `useEffect` corre una vez; pull-to-refresh explícito si hace falta |
| OQ-4 | ¿Rango de fechas en Actividad? | Sí, `"10 – 16 oct"` bajo el selector, en es-CO y UTC |
