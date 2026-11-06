import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'
import { toPercent } from '@/core/utils/numbers'

type TopicCompactRowProps = {
  topicName: string
  attempts: number
  /** 0..1 */
  accuracy: number
  isWeak?: boolean
  correctCount?: number
}

/**
 * Fila compacta (~46px) para comparar temas **dentro de un periodo**.
 * El estado global de un tema usa `TopicRing`.
 */
export function TopicCompactRow({
  topicName,
  attempts,
  accuracy,
  isWeak = false,
  correctCount,
}: TopicCompactRowProps) {
  const value = toPercent(accuracy)
  const counts =
    correctCount !== undefined ? `, ${correctCount} de ${attempts} aciertos` : ''
  const label = `${topicName}: ${value}% de acierto en ${attempts} respuestas${
    isWeak ? ', tema débil' : ''
  }${counts}`

  return (
    <View
      accessible
      accessibilityRole="progressbar"
      accessibilityLabel={label}
      accessibilityValue={{ min: 0, max: 100, now: value }}
      style={styles.row}
    >
      <View style={styles.header}>
        <Text style={styles.name} numberOfLines={1}>
          {topicName}
        </Text>

        {isWeak ? (
          <View style={styles.chip}>
            <MaterialCommunityIcons name="alert-outline" size={12} color={colors.accent} />
            <Text style={styles.chipText}>Débil</Text>
          </View>
        ) : null}

        <Text style={styles.meta}>
          {attempts} resp · {value}%
        </Text>
      </View>

      <View style={styles.track}>
        <View
          style={[
            styles.fill,
            {
              width: `${value}%`,
              backgroundColor: isWeak ? colors.accent : colors.primary,
            },
          ]}
        />
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  row: {
    gap: spacing.xs + 2,
    paddingVertical: spacing.sm + 2,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  name: {
    flex: 1,
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.foreground,
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    paddingHorizontal: spacing.sm,
    borderRadius: radius.sm,
    backgroundColor: `${colors.accent}22`,
  },
  chipText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.accent,
  },
  meta: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  track: {
    height: 4,
    borderRadius: 2,
    backgroundColor: colors.muted,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    borderRadius: 2,
  },
})
