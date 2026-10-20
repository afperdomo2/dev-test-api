import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'
import type { TopicMasteryItem } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type TopicMasteryRowProps = {
  item: TopicMasteryItem
}

/** `accuracy` llega como 0..1; la barra trabaja en porcentaje entero. */
function percent(accuracy: number): number {
  return Math.round(accuracy * 100)
}

export function TopicMasteryRow({ item }: TopicMasteryRowProps) {
  const value = percent(item.accuracy)
  const isWeak = item.isWeak

  return (
    <View style={styles.row}>
      <View style={styles.header}>
        <Text style={styles.name} numberOfLines={1}>
          {item.topicName}
        </Text>

        {isWeak ? (
          <View style={styles.chip} accessibilityLabel="Tema débil">
            <MaterialCommunityIcons name="alert-outline" size={12} color={colors.accent} />
            <Text style={styles.chipText}>Débil</Text>
          </View>
        ) : null}
      </View>

      <View
        accessible
        accessibilityRole="progressbar"
        accessibilityLabel={`${item.topicName}: ${value}% de acierto`}
        accessibilityValue={{ min: 0, max: 100, now: value }}
        style={styles.track}
      >
        <View
          style={[
            styles.fill,
            { width: `${value}%`, backgroundColor: isWeak ? colors.accent : colors.primary },
          ]}
        />
      </View>

      <Text style={styles.counts}>
        {item.correctCount}/{item.attempts} aciertos
      </Text>
    </View>
  )
}

const styles = StyleSheet.create({
  row: {
    gap: spacing.sm,
    padding: spacing.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: spacing.sm,
  },
  name: {
    flex: 1,
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    backgroundColor: `${colors.accent}22`,
  },
  chipText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.accent,
  },
  track: {
    height: 8,
    borderRadius: radius.sm,
    backgroundColor: colors.muted,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    borderRadius: radius.sm,
  },
  counts: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
