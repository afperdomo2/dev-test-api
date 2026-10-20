# SPEC — mobile-weekly

| Campo | Valor |
|---|---|
| Module id | `mobile-weekly` |
| Origen | Alcance nuevo mobile (D-26, D-33) |
| Estado | `ready` |
| Depende de | `insights` (backend, `GET /insights/summary`) |
| Migración goose | Ninguna |
| Dependencia nueva | **Ninguna** |

## Objective

Que el usuario vea en el móvil un resumen de **su última semana o mes**: cuánto practicó, con qué acierto, en qué temas, cuánto XP ganó y cómo va su racha.

**Historias**:
- Desde el hub de Progreso entro en **Mi semana**.
- Alterno **Semana / Mes** y los datos cambian.
- Veo 6 cifras, mis temas más practicados, mis temas débiles y la actividad día a día.
- Si no he hecho nada en el periodo, veo un mensaje claro en vez de ceros sin contexto.

**Éxito**: los números coinciden con `GET /insights/summary` y con lo que muestra la web.

## Contexto

- El backend y la web de insights están especificados en `SPEC-insights.md` (`IN-1…IN-8`) y **aún no implementados**. Esta spec consume su contrato.
- `weakTopics` usa la definición única de `review-mastery` (D-08): la app **no** la recalcula.
- La pestaña y el hub los crea `mobile-progress`.

## Commands

```bash
pnpm --filter mobile typecheck
pnpm --filter mobile test:run
pnpm --filter mobile lint
```

## Project structure

```
apps/mobile/app/(tabs)/progreso/semana.tsx           # NUEVO (thin re-export)
apps/mobile/src/features/progress/
├── screens/WeeklySummaryScreen.tsx                 # NUEVO
├── components/KpiCard.tsx                          # NUEVO
├── components/PeriodToggle.tsx                     # NUEVO
├── components/DailyActivityBar.tsx                 # NUEVO
└── screens/WeeklySummaryScreen.test.tsx            # NUEVO
apps/mobile/src/features/progress/screens/ProgressHubScreen.tsx  # + fila "Mi semana"
apps/mobile/src/core/api/services/insights.service.ts            # NUEVO
apps/mobile/src/core/api/services/insights.service.test.ts       # NUEVO
```

## Contrato

### Servicio

```ts
// src/core/api/services/insights.service.ts
export async function getInsightSummary(period: InsightPeriod): Promise<InsightSummary> {
  return apiClient.get<InsightSummary>('/api/v1/insights/summary', { period })
}
```

Tipos: `InsightPeriod`, `InsightSummary`, `TopicInsightItem`, `DailyActivityItem` (todos en `@devtest/shared`, ya creados por la spec de insights).

### Pantalla

- **Selector de periodo** (`PeriodToggle`): dos botones "Semana" / "Mes"; el activo se distingue por color **y** por `accessibilityState={{ selected: true }}`.
- **6 KPI cards** (`KpiCard`), en 2 columnas:
  `Sesiones` · `Preguntas` · `Acierto` (% redondeado) · `Tiempo` (`Xh Ym` o `Y min`) · `XP ganado` · `Racha` (días).
- **Más practicados**: `topTopics`, reutilizando `TopicMasteryRow` de `mobile-progress` (misma barra y a11y).
- **Temas débiles**: `weakTopics` (mismos items que `TopicMasteryItem`); si viene vacío, se oculta la sección (no un bloque vacío).
- **Actividad diaria**: una barra por día (`DailyActivityBar`), altura proporcional al máximo del periodo; etiqueta accesible `"lun 12: 8 preguntas"`.
- Fechas: el backend manda `YYYY-MM-DD`; se muestran como día de la semana abreviado + número.

**Sin librería de gráficas**: barras con `View` y porcentajes.

### Estados

| Estado | UI |
|---|---|
| Cargando | skeletons de las KPI + lista |
| Error | mensaje + "Reintentar" |
| Sin actividad (`questionsAnswered === 0` y `sessions === 0`) | mensaje "No hay actividad en este periodo" y selector de periodo visible |
| Con datos | todo lo anterior |

## Code style

- Tokens de diseño; sin hex.
- `InsightSummary` se pasa por props a los componentes presentacionales.
- Cambiar de periodo cambia la `key`/estado y vuelve a pedir datos (no cachea ambos).

## Testing strategy

| Archivo | Casos |
|---|---|
| `insights.service.test.ts` | URL y `params.period` |
| `WeeklySummaryScreen.test.tsx` | pinta las 6 KPIs; cambia de periodo y vuelve a llamar al servicio; estado vacío; error con reintento; oculta "Temas débiles" si viene vacío |
| `DailyActivityBar.test.tsx` | altura proporcional y `accessibilityLabel` correcto |
| `PeriodToggle.test.tsx` | emite el periodo al pulsar y marca el activo |

## Boundaries

- **Always**: consumir `GET /insights/summary` (no calcular métricas en el cliente); reutilizar `TopicMasteryRow`; tokens; estados vacío/error.
- **Ask first**: añadir gráficas de líneas o librería de charts; cachear ambos periodos; exportar el informe; mostrar periodos distintos a semana/mes.
- **Never**: recalcular `weakTopics` (D-08); duplicar `InsightSummary`; asumir que el endpoint existe antes de la Ola 2 (esta spec se implementa después).

## Success Criteria

1. Desde el hub se abre **Mi semana** y se carga el periodo `week` por defecto.
2. Pulsar **Mes** vuelve a consultar con `period=month`.
3. Las 6 cifras coinciden con la respuesta del backend.
4. `weakTopics` coincide con lo que muestra Progreso → Temas (mismos items).
5. Sin actividad en el periodo → mensaje, no una pantalla en blanco.
6. Error de red → "Reintentar" funciona.
7. `typecheck`, `test:run` y `lint` en verde.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Qué 6 KPIs exactamente? | Las de `InsightSummary`: `sessions`, `questionsAnswered`, `accuracy`, `timeSpentMs`, `xpGained`, `streakDays`. |
| OQ-2 | ¿Se muestra `from`/`to` del periodo? | Sí, como subtítulo bajo el selector ("2–8 oct"). |
| OQ-3 | ¿Se puede deslizar entre periodos? | No: dos botones. Menos gestos que mantener. |
