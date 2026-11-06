import { StyleSheet, View } from 'react-native'
import type { DailyActivityItem } from '@devtest/shared'

import { colors } from '@/core/theme/tokens'

const MAX_HEIGHT = 28

type ActivitySparklineProps = {
  days: Array<DailyActivityItem>
}

/** Mini gráfico de barras de los últimos días (sin etiquetas). */
export function ActivitySparkline({ days }: ActivitySparklineProps) {
  if (days.length === 0) return null

  const max = Math.max(1, ...days.map((day) => day.questionsAnswered))
  const total = days.reduce((sum, day) => sum + day.questionsAnswered, 0)
  const label = `Actividad de los últimos ${days.length} días: ${total} preguntas`

  return (
    <View
      accessible
      accessibilityRole="image"
      accessibilityLabel={label}
      style={styles.row}
    >
      {days.map((day) => {
        const active = day.questionsAnswered > 0
        const height = active
          ? Math.max(6, Math.round((day.questionsAnswered / max) * MAX_HEIGHT))
          : 3

        return (
          <View
            key={day.date}
            testID={`spark-bar-${day.date}`}
            style={[
              styles.bar,
              { height, backgroundColor: active ? colors.primary : colors.muted },
            ]}
          />
        )
      })}
    </View>
  )
}

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    height: MAX_HEIGHT,
    gap: 4,
  },
  bar: {
    flex: 1,
    borderRadius: 2,
  },
})
