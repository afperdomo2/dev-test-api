import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type TopicMasteryRowProps = {
  topicName: string
  attempts: number
  accuracy: number
  /** Solo en la pantalla de Temas; en informes no se pinta el chip. */
  isWeak?: boolean
  /** Si se omite, no se pinta la línea "x/y aciertos". */
  correctCount?: number
}

/** `accuracy` llega como 0..1; la barra trabaja en porcentaje entero. */
function percent(accuracy: number): number {
  return Math.round(accuracy * 100)
}

export function TopicMasteryRow({
  topicName,
  attempts,
  accuracy,
  isWeak = false,
  correctCount,
}: TopicMasteryRowProps) {
  const value = percent(accuracy)

  return (
    <View style={styles.row}>
      <View style={styles.header}>
        <Text style={styles.name} numberOfLines={1}>
          {topicName}
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
        accessibilityLabel={`${topicName}: ${value}% de acierto`}
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

      {correctCount !== undefined ? (
        <Text style={styles.counts}>
          {correctCount}/{attempts} aciertos
        </Text>
      ) : null}
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
