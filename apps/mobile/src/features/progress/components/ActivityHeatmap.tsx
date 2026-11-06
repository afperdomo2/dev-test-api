import { StyleSheet, Text, View } from 'react-native'
import type { DailyActivityItem } from '@devtest/shared'

import { HEAT_COLORS, heatLevel } from '@/core/theme/charts'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const WEEKDAY_LABELS = ['L', 'M', 'X', 'J', 'V', 'S', 'D'] as const
const CELL = 18
const GAP = 6

/** Lunes = 0 … domingo = 6, en UTC para no desplazar el día. */
function weekdayIndex(date: string): number {
  const day = new Date(`${date}T00:00:00Z`).getUTCDay()
  return (day + 6) % 7
}

function dayLabel(date: string): string {
  const value = new Date(`${date}T00:00:00Z`)
  const weekday = value.toLocaleDateString('es-ES', { weekday: 'short', timeZone: 'UTC' })
  return `${weekday} ${value.getUTCDate()}`
}

type ActivityHeatmapProps = {
  days: Array<DailyActivityItem>
}

/**
 * Mapa de calor tipo GitHub para 30 días: columnas = días de la semana,
 * filas = semanas. Se dibuja con `View` (sin librería de gráficas).
 */
export function ActivityHeatmap({ days }: ActivityHeatmapProps) {
  if (days.length === 0) return null

  const max = Math.max(1, ...days.map((day) => day.questionsAnswered))
  const leading = weekdayIndex(days[0]!.date)

  return (
    <View style={styles.container}>
      <View style={styles.grid}>
        {Array.from({ length: leading }).map((_, index) => (
          <View key={`blank-${index}`} style={styles.blank} />
        ))}

        {days.map((day) => (
          <View
            key={day.date}
            accessible
            accessibilityRole="image"
            accessibilityLabel={`${dayLabel(day.date)}: ${day.questionsAnswered} preguntas`}
            testID={`heat-cell-${day.date}`}
            style={[
              styles.cell,
              { backgroundColor: HEAT_COLORS[heatLevel(day.questionsAnswered, max)] },
            ]}
          />
        ))}
      </View>

      <View style={styles.weekdays}>
        {WEEKDAY_LABELS.map((label, index) => (
          <Text key={`${label}-${index}`} style={styles.weekday}>
            {label}
          </Text>
        ))}
      </View>

      <View style={styles.legend}>
        <Text style={styles.legendText}>Menos</Text>
        {HEAT_COLORS.map((color) => (
          <View key={color} style={[styles.legendCell, { backgroundColor: color }]} />
        ))}
        <Text style={styles.legendText}>Más</Text>
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    gap: spacing.sm,
    alignItems: 'center',
  },
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    width: CELL * 7 + GAP * 6,
    gap: GAP,
  },
  blank: {
    width: CELL,
    height: CELL,
  },
  cell: {
    width: CELL,
    height: CELL,
    borderRadius: 4,
  },
  weekdays: {
    flexDirection: 'row',
    width: CELL * 7 + GAP * 6,
    gap: GAP,
  },
  weekday: {
    width: CELL,
    textAlign: 'center',
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  legend: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  legendCell: {
    width: 10,
    height: 10,
    borderRadius: radius.sm / 2,
  },
  legendText: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
