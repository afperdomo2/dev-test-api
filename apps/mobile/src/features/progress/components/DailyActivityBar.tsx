import { StyleSheet, Text, View } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type DailyActivityBarProps = {
  /** YYYY-MM-DD (UTC) */
  date: string
  value: number
  max: number
}

/** Día de la semana en UTC, para que no se desplace un día. */
function weekday(date: string): string {
  return new Date(`${date}T00:00:00Z`).toLocaleDateString('es-ES', {
    weekday: 'short',
    timeZone: 'UTC',
  })
}

export function DailyActivityBar({ date, value, max }: DailyActivityBarProps) {
  const height = value > 0 ? Math.max(4, Math.round((value / max) * 80)) : 2

  return (
    <View
      style={styles.item}
      accessible
      accessibilityRole="image"
      accessibilityLabel={`${date}: ${value} preguntas`}
    >
      <View
        testID={`daily-bar-${date}`}
        style={[
          styles.bar,
          {
            height,
            backgroundColor: value > 0 ? colors.primary : colors.muted,
          },
        ]}
      />
      <Text style={styles.label}>{weekday(date)}</Text>
    </View>
  )
}

const styles = StyleSheet.create({
  item: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'flex-end',
    gap: spacing.xs,
    minWidth: 18,
  },
  bar: {
    width: '100%',
    borderRadius: radius.sm,
  },
  label: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
