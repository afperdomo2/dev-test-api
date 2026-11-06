import { StyleSheet, Text, View } from 'react-native'
import type { DailyActivityItem } from '@devtest/shared'

import { colors, spacing, typography } from '@/core/theme/tokens'

/** Lunes = 0 … domingo = 6, en UTC para no desplazar el día. */
function weekdayInitial(date: string): string {
  return new Date(`${date}T00:00:00Z`)
    .toLocaleDateString('es-ES', { weekday: 'narrow', timeZone: 'UTC' })
    .toUpperCase()
}

type DayStreakStripProps = {
  days: Array<DailyActivityItem>
}

/**
 * Casillas de los últimos 7 días: rellenas si hubo práctica.
 * Es un elemento visual; la etiqueta accesible la pone el contenedor pulsable.
 */
export function DayStreakStrip({ days }: DayStreakStripProps) {
  if (days.length === 0) return null

  return (
    <View style={styles.row} importantForAccessibility="no-hide-descendants">
      {days.map((day) => {
        const active = day.questionsAnswered > 0

        return (
          <View key={day.date} style={styles.item}>
            <View
              testID={`day-cell-${day.date}`}
              style={[styles.cell, active ? styles.cellActive : styles.cellIdle]}
            />
            <Text style={[styles.label, active && styles.labelActive]}>
              {weekdayInitial(day.date)}
            </Text>
          </View>
        )
      })}
    </View>
  )
}

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    gap: spacing.xs + 2,
  },
  item: {
    flex: 1,
    alignItems: 'center',
    gap: spacing.xs,
  },
  cell: {
    width: '100%',
    height: 26,
    borderRadius: 6,
  },
  cellActive: {
    backgroundColor: colors.primary,
  },
  cellIdle: {
    backgroundColor: colors.muted,
  },
  label: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  labelActive: {
    color: colors.foreground,
  },
})
