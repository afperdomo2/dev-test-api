import { useCallback, useEffect, useState } from 'react'
import {
  ActivityIndicator,
  Alert,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { useRouter } from 'expo-router'
import type { OverviewResponse, TopicXpItem } from '@devtest/shared'

import { useAuthStore } from '@/core/auth/auth.store'
import { getOverview, freezeStreak } from '@/core/api/services/gamification.service'
import { createDailySession } from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { setPendingTip } from '@/features/sessions/dailyTipStore'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

function getLevelProgress(totalXp: number, level: number): number {
  const xpForCurrent = (level - 1) ** 2 * 100
  const xpForNext = level ** 2 * 100
  if (xpForNext <= xpForCurrent) return 1
  return Math.min(1, Math.max(0, (totalXp - xpForCurrent) / (xpForNext - xpForCurrent)))
}

function getMaxTopicXp(topics: TopicXpItem[]): number {
  if (topics.length === 0) return 1
  return Math.max(...topics.map((t) => t.xp))
}

export default function HomeScreen() {
  const router = useRouter()
  const user = useAuthStore((s) => s.user)
  const [overview, setOverview] = useState<OverviewResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [creatingDaily, setCreatingDaily] = useState(false)
  const [dailyError, setDailyError] = useState<string | null>(null)
  const [freezing, setFreezing] = useState(false)

  useEffect(() => {
    let cancelled = false
    getOverview()
      .then((res) => {
        if (!cancelled) setOverview(res)
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const handleDailySession = useCallback(async () => {
    setCreatingDaily(true)
    setDailyError(null)
    try {
      const res = await createDailySession()
      setPendingTip(res.session.id, res.tip)
      router.push(`/(tabs)/sesiones/${res.session.id}`)
    } catch (err) {
      setDailyError(getErrorDetail(err, 'No se pudo crear la sesion diaria'))
    } finally {
      setCreatingDaily(false)
    }
  }, [router])

  const handleFreeze = useCallback(async () => {
    if (!overview || overview.streak.freezesAvailable <= 0) return
    Alert.alert(
      'Usar freeze',
      'Esto congelara tu racha actual por 1 dia. ¿Continuar?',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Usar',
          onPress: async () => {
            setFreezing(true)
            try {
              const reward = await freezeStreak()
              setOverview((prev) =>
                prev
                  ? {
                      ...prev,
                      totalXp: reward.totalXp,
                      level: reward.level,
                      xpToLevel: reward.xpToLevel,
                      streak: reward.streak,
                      dailyGoal: reward.dailyGoal,
                    }
                  : prev,
              )
            } catch (err) {
              Alert.alert('Error', getErrorDetail(err, 'No se pudo usar el freeze'))
            } finally {
              setFreezing(false)
            }
          },
        },
      ],
    )
  }, [overview])

  const userName = user?.email?.split('@')[0] ?? 'Usuario'
  const levelProgress = overview ? getLevelProgress(overview.totalXp, overview.level) : 0
  const maxTopicXp = overview ? getMaxTopicXp(overview.topicXp) : 1

  if (loading) {
    return (
      <ScrollView style={styles.container} contentContainerStyle={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </ScrollView>
    )
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      contentInsetAdjustmentBehavior="automatic"
    >
      <View style={styles.header}>
        <Text style={styles.greeting}>Hola, {userName}</Text>
        {overview && (
          <View style={styles.levelBadge}>
            <MaterialCommunityIcons name="lightning-bolt" size={14} color={colors.primary} />
            <Text style={styles.levelBadgeText}>Nv. {overview.level}</Text>
          </View>
        )}
      </View>

      <Pressable
        style={({ pressed }) => [
          styles.ctaCard,
          pressed && styles.ctaCardPressed,
          creatingDaily && styles.ctaCardDisabled,
        ]}
        onPress={() => void handleDailySession()}
        disabled={creatingDaily}
        accessibilityRole="button"
        accessibilityLabel="Repasar hoy"
      >
        {creatingDaily ? (
          <ActivityIndicator color={colors.onPrimary} />
        ) : (
          <>
            <MaterialCommunityIcons name="book-open-variant" size={28} color={colors.onPrimary} />
            <View style={styles.ctaText}>
              <Text style={styles.ctaTitle}>Repasar hoy</Text>
              <Text style={styles.ctaSubtitle}>
                Resuelve las preguntas vencidas y mantén tu racha
              </Text>
            </View>
            <MaterialCommunityIcons name="arrow-right" size={20} color={colors.onPrimary} />
          </>
        )}
      </Pressable>

      {dailyError && <Text style={styles.error}>{dailyError}</Text>}

      {overview && (
        <>
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Progreso</Text>
            <View style={styles.card}>
              <View style={styles.levelRow}>
                <Text style={styles.levelLabel}>Nivel {overview.level}</Text>
                <Text style={styles.levelXp}>
                  {overview.totalXp} / {overview.totalXp + overview.xpToLevel} XP
                </Text>
              </View>
              <View style={styles.progressTrack}>
                <View style={[styles.progressFill, { width: `${levelProgress * 100}%` }]} />
              </View>
            </View>
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Racha</Text>
            <View style={styles.card}>
              <View style={styles.streakRow}>
                <View style={styles.streakMain}>
                  <MaterialCommunityIcons name="fire" size={24} color={colors.accent} />
                  <Text style={styles.streakValue}>{overview.streak.current}</Text>
                  <Text style={styles.streakLabel}>dias</Text>
                </View>
                <View style={styles.streakRecord}>
                  <Text style={styles.streakRecordLabel}>Mejor</Text>
                  <Text style={styles.streakRecordValue}>{overview.streak.longest}</Text>
                </View>
              </View>

              <View style={styles.streakActions}>
                <View style={styles.freezeInfo}>
                  <MaterialCommunityIcons name="snowflake" size={14} color={colors.mutedForeground} />
                  <Text style={styles.freezeText}>
                    {overview.streak.freezesAvailable} freezes disponibles
                  </Text>
                </View>
                <Pressable
                  style={({ pressed }) => [
                    styles.freezeButton,
                    pressed && styles.freezeButtonPressed,
                    (overview.streak.freezesAvailable <= 0 || freezing) &&
                      styles.freezeButtonDisabled,
                  ]}
                  onPress={() => void handleFreeze()}
                  disabled={overview.streak.freezesAvailable <= 0 || freezing}
                  accessibilityRole="button"
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
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Meta diaria</Text>
            <View style={styles.card}>
              <View style={styles.goalHeader}>
                <Text style={styles.goalValue}>
                  {overview.dailyGoal.current} / {overview.dailyGoal.target} XP
                </Text>
                {overview.dailyGoal.reached && (
                  <View style={styles.goalReachedBadge}>
                    <MaterialCommunityIcons name="check-circle" size={12} color={colors.primary} />
                    <Text style={styles.goalReachedText}>Meta alcanzada</Text>
                  </View>
                )}
              </View>
              <View style={styles.progressTrack}>
                <View
                  style={[
                    styles.progressFill,
                    {
                      width: `${Math.min(100, overview.dailyGoal.progress)}%`,
                      backgroundColor: overview.dailyGoal.reached ? colors.primary : colors.accent,
                    },
                  ]}
                />
              </View>
            </View>
          </View>

          {overview.topicXp.length > 0 && (
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Temas</Text>
              <View style={styles.card}>
                {overview.topicXp.map((topic) => (
                  <View key={topic.topicId} style={styles.topicRow}>
                    <Text style={styles.topicName} numberOfLines={1}>
                      {topic.name}
                    </Text>
                    <View style={styles.topicBarContainer}>
                      <View style={styles.topicBarTrack}>
                        <View
                          style={[
                            styles.topicBarFill,
                            { width: `${(topic.xp / maxTopicXp) * 100}%` },
                          ]}
                        />
                      </View>
                      <Text style={styles.topicXp}>{topic.xp} XP</Text>
                    </View>
                  </View>
                ))}
              </View>
            </View>
          )}
        </>
      )}
    </ScrollView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    padding: spacing.lg,
    paddingBottom: spacing.xl,
    gap: spacing.lg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  greeting: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
  },
  levelBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: `${colors.primary}20`,
    paddingHorizontal: spacing.sm,
    paddingVertical: 4,
    borderRadius: radius.sm,
    gap: spacing.xs,
  },
  levelBadgeText: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.primary,
  },
  ctaCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    padding: spacing.lg,
    gap: spacing.md,
    minHeight: 56,
  },
  ctaCardPressed: {
    opacity: 0.85,
  },
  ctaCardDisabled: {
    opacity: 0.6,
    justifyContent: 'center',
  },
  ctaText: {
    flex: 1,
    gap: spacing.xs,
  },
  ctaTitle: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.bold,
    color: colors.onPrimary,
  },
  ctaSubtitle: {
    fontSize: typography.size.sm,
    color: `${colors.onPrimary}CC`,
  },
  error: {
    fontSize: typography.size.sm,
    color: colors.destructive,
  },
  section: {
    gap: spacing.xs,
  },
  sectionTitle: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
    paddingHorizontal: spacing.xs,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.md,
  },
  levelRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  levelLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  levelXp: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  progressTrack: {
    height: 8,
    backgroundColor: colors.muted,
    borderRadius: 4,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: colors.primary,
    borderRadius: 4,
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
  streakActions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
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
    paddingVertical: spacing.xs,
    borderRadius: radius.sm,
    borderWidth: 1,
    borderColor: colors.primary,
    minHeight: 32,
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
  goalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  goalValue: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  goalReachedBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: `${colors.primary}15`,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    gap: spacing.xs,
  },
  goalReachedText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
    color: colors.primary,
  },
  topicRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  topicName: {
    width: 100,
    fontSize: typography.size.sm,
    color: colors.foreground,
  },
  topicBarContainer: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  topicBarTrack: {
    flex: 1,
    height: 6,
    backgroundColor: colors.muted,
    borderRadius: 3,
    overflow: 'hidden',
  },
  topicBarFill: {
    height: '100%',
    backgroundColor: colors.primary,
    borderRadius: 3,
  },
  topicXp: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
    minWidth: 45,
    textAlign: 'right',
  },
})
