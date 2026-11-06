import { MaterialCommunityIcons } from '@expo/vector-icons'
import { ActivityIndicator, Pressable, StyleSheet, Text, View } from 'react-native'
import type { DailyActivityItem, DailyGoalInfo, StreakInfo } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

import { DayStreakStrip } from './DayStreakStrip'

type TodayCardProps = {
  streak: StreakInfo
  dailyGoal: DailyGoalInfo
  /** Actividad de los últimos 7 días; `null` mientras carga o si falló. */
  week: Array<DailyActivityItem> | null
  freezing: boolean
  onFreeze: () => void
  onOpenActivity: () => void
}

/** "Tu día": racha + actividad de la semana + meta diaria en una sola card. */
export function TodayCard({
  streak,
  dailyGoal,
  week,
  freezing,
  onFreeze,
  onOpenActivity,
}: TodayCardProps) {
  const practiced = week?.filter((day) => day.questionsAnswered > 0).length ?? 0
  const goalValue = Math.min(100, Math.round(dailyGoal.progress))
  const canFreeze = streak.freezesAvailable > 0 && !freezing

  return (
    <View style={styles.card}>
      <View style={styles.streakRow}>
        <View style={styles.streakMain}>
          <MaterialCommunityIcons name="fire" size={22} color={colors.accent} />
          <Text style={styles.streakValue}>{streak.current}</Text>
          <Text style={styles.streakLabel}>{streak.current === 1 ? 'día' : 'días'}</Text>
        </View>

        <View style={styles.streakRecord}>
          <Text style={styles.streakRecordLabel}>Mejor</Text>
          <Text style={styles.streakRecordValue}>{streak.longest}</Text>
        </View>
      </View>

      {week ? (
        <Pressable
          onPress={onOpenActivity}
          style={({ pressed }) => [styles.week, pressed && styles.pressed]}
          accessibilityRole="button"
          accessibilityLabel={`Ver tu actividad: practicaste ${practiced} de los últimos ${week.length} días`}
        >
          <DayStreakStrip days={week} />
        </Pressable>
      ) : null}

      <View style={styles.divider} />

      <View style={styles.goalHeader}>
        <Text style={styles.goalTitle}>Meta diaria</Text>
        <View style={styles.goalValueRow}>
          {dailyGoal.reached ? (
            <View style={styles.goalBadge}>
              <MaterialCommunityIcons name="check-circle" size={12} color={colors.primary} />
              <Text style={styles.goalBadgeText}>Alcanzada</Text>
            </View>
          ) : null}
          <Text style={styles.goalValue}>
            {dailyGoal.current} / {dailyGoal.target} XP
          </Text>
        </View>
      </View>

      <View
        accessible
        accessibilityRole="progressbar"
        accessibilityLabel={`Meta diaria: ${dailyGoal.current} de ${dailyGoal.target} XP`}
        accessibilityValue={{ min: 0, max: 100, now: goalValue }}
        style={styles.track}
      >
        <View
          style={[
            styles.fill,
            {
              width: `${goalValue}%`,
              backgroundColor: dailyGoal.reached ? colors.primary : colors.accent,
            },
          ]}
        />
      </View>

      <View style={styles.freezeRow}>
        <View style={styles.freezeInfo}>
          <MaterialCommunityIcons name="snowflake" size={14} color={colors.mutedForeground} />
          <Text style={styles.freezeText}>
            {streak.freezesAvailable} {streak.freezesAvailable === 1 ? 'freeze' : 'freezes'}
          </Text>
        </View>

        <Pressable
          style={({ pressed }) => [
            styles.freezeButton,
            pressed && styles.freezeButtonPressed,
            !canFreeze && styles.freezeButtonDisabled,
          ]}
          onPress={onFreeze}
          disabled={!canFreeze}
          accessibilityRole="button"
          accessibilityState={{ disabled: !canFreeze, busy: freezing }}
          accessibilityLabel="Usar freeze para conservar racha"
        >
          {freezing ? (
            <ActivityIndicator size="small" color={colors.primary} />
          ) : (
            <Text style={styles.freezeButtonText}>Usar</Text>
          )}
        </Pressable>
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.md,
  },
  streakRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  streakMain: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  streakValue: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  streakLabel: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  streakRecord: {
    alignItems: 'flex-end',
  },
  streakRecordLabel: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  streakRecordValue: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  week: {
    minHeight: 44,
    justifyContent: 'center',
    borderRadius: radius.sm,
  },
  pressed: {
    opacity: 0.7,
  },
  divider: {
    height: 1,
    backgroundColor: colors.border,
  },
  goalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  goalTitle: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  goalValueRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  goalBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  goalBadgeText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
    color: colors.primary,
  },
  goalValue: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  track: {
    height: 6,
    borderRadius: 3,
    backgroundColor: colors.muted,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    borderRadius: 3,
  },
  freezeRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  freezeInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  freezeText: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  freezeButton: {
    paddingHorizontal: spacing.md,
    borderRadius: radius.sm,
    borderWidth: 1,
    borderColor: colors.primary,
    minHeight: 44,
    minWidth: 72,
    justifyContent: 'center',
    alignItems: 'center',
  },
  freezeButtonPressed: {
    backgroundColor: `${colors.primary}15`,
  },
  freezeButtonDisabled: {
    borderColor: colors.mutedForeground,
    opacity: 0.5,
  },
  freezeButtonText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
    color: colors.primary,
  },
})
