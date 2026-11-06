# Mobile design system — dev-test

> **Source of truth** para el diseño de la app mobile. Generada con la skill `ui-ux-pro-max` (producto "Coding Challenge & Practice") y curada a mano. Espejo en código: `apps/mobile/src/core/theme/tokens.ts`. Estilo: **Dark Mode (OLED)**, dark-only.

## Style

- **Dark Mode (OLED)** + Minimalism & Swiss Style (base limpia; sin cyberpunk/neón excesivo).
- Light mode: **no recomendado** (`app.json` → `userInterfaceStyle: "dark"`).
- Efectos: glow mínimo, transiciones dark→light, foco visible, **sin fondos blancos puros**.
- Tipografía UI: fuente del sistema (apple-design: default a la fuente de la plataforma). Mono (JetBrains Mono / SF Mono / Menlo) para bloques de código.

## Color palette

| Token | Hex | Uso |
|-------|-----|-----|
| primary | `#22C55E` | Verde código / correcto / CTA |
| onPrimary | `#0F172A` | Texto sobre primary |
| secondary | `#059669` | Verde secundario |
| accent | `#D97706` | Ámbar / dificultad media |
| background | `#0F172A` | Fondo app |
| foreground | `#FFFFFF` | Texto principal |
| card | `#192134` | Superficies de tarjeta |
| cardForeground | `#FFFFFF` | Texto sobre card |
| muted | `#10242E` | Fondos secundarios |
| mutedForeground | `#94A3B8` | Texto secundario |
| border | `rgba(255,255,255,0.08)` | Bordes |
| destructive | `#DC2626` | Error / dificultad alta |
| onDestructive | `#FFFFFF` | Texto sobre destructive |
| ring | `#22C55E` | Foco |

**Gradiente de dificultad** (ya mapeado en `@devtest/shared` → `DIFFICULTY_COLORS`):
- `beginner` → success (verde `#22C55E`)
- `intermediate` → warning (ámbar `#D97706`)
- `advanced` → error (rojo `#DC2626`)

## Spacing

`xs=4 · sm=8 · md=16 · lg=24 · xl=32 · 2xl=48 · 3xl=64`

## Radius

`sm=8 (inputs/buttons) · md=12 (cards) · lg=16 (modals/sheets)`

## Typography

- UI: fuente del sistema (SF Pro en iOS, Roboto en Android).
- Código (`code_completion`, code challenges): monospace (`Menlo` iOS / `monospace` Android en tokens; JetBrains Mono si se cargan fuentes custom).
- Títulos grandes: tracking negativo (p.ej. `-0.5`), leading compacto.

## Data viz

- **Anillo** = dominio global de un tema (contra 100 %). `TopicRing`, SVG con `strokeDasharray`, **64 px**, `strokeWidth` 6, % escrito en el centro. Si el tema es débil: arco y número en `accent`, **sin icono** (el % es el dato y "tema débil" queda en la etiqueta accesible y en el filtro Débiles). Rejilla de 3 columnas con `rowGap: spacing.lg` (24 px): el aire separa las filas, sin líneas.
- **Barra compacta** = comparación dentro de un periodo. `TopicCompactRow`, 4 px de alto, misma barra para "Más practicados" y "Puntos a reforzar".
- **Heatmap** = 30 días, celdas 18 px, 4 niveles (`HEAT_COLORS` en `src/core/theme/charts.ts`) derivados de `primary` sobre `muted`; `heatLevel(valor, máximo)` decide el tramo.
- **Sparkline** = 7 barras de 28 px (sin etiquetas) para el resumen del dashboard.
- Reglas: el dato siempre se puede leer sin color (%, número o texto); nada de librería de charts; alturas proporcionales calculadas con `Math.max`/redondeo, nunca valores fijos por dato.

## Anti-patterns (no usar)

- ❌ Fondos blancos puros.
- ❌ Emojis como iconos → usar SF Symbols (`sf`) / Material Symbols (`md`).
- ❌ Texto con contraste < 4.5:1.
- ❌ Cambios de estado instantáneos → transiciones 150-300ms.
- ❌ Foco invisible.

## Regenerar propuesta (opcional)

La skill puede proponer ajustes; el archivo canónico es este:

```bash
python .agents/skills/ui-ux-pro-max/scripts/search.py "coding challenge leetcode practice developer tool" --design-system
```