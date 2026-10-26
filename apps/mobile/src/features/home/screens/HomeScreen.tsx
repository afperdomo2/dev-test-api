import { useCallback, useEffect, useState } from 'react'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { useRouter } from 'expo-router'
import {
  ActivityIndicator,
  Alert,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import type {
  BadgesResponse,
  DailyActivityItem,
  InsightSummary,
  OverviewResponse,
} from '@devtest/shared'

import { useAuthStore } from '@/core/auth/auth.store'
import { BadgeStrip } from '@/core/components/BadgeStrip'
import { LevelLine } from '@/core/components/LevelLine'
import { SectionHeader } from '@/core/components/SectionHeader'
import { freezeStreak, getBadges, getOverview } from '@/core/api/services/gamification.service'
import { getInsightSummary } from '@/core/api/services/insights.service'
import { createDailySession } from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { setPendingTip } from '@/features/sessions/dailyTipStore'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

import { TodayCard } from '../components/TodayCard'

export default function HomeScreen() {
  const router = useRouter()
  const user = useAuthStore((s) => s.user)

  const [overview, setOverview] = useState<OverviewResponse | null>(null)
  const [overviewError, setOverviewError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [week, setWeek] = useState<Array<DailyActivityItem> | null>(null)
  const [badges, setBadges] = useState<BadgesResponse | null>(null)

  const [creatingDaily, setCreatingDaily] = useState(false)
  const [dailyError, setDailyError] = useState<string | null>(null)
  const [freezing, setFreezing] = useState(false)

  const loadOverview = useCallback(async () => {
    setLoading(true)
    setOverviewError(null)

    try {
      setOverview(await getOverview())
    } catch {
      setOverview(null)
      setOverviewError('No se pudo cargar tu progreso')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadOverview()
  }, [loadOverview])

  // La actividad de 7 días es contexto de la racha: si falla, la card sigue
  // mostrando racha y meta. No bloquea ni ensucia la pantalla con un error.
  useEffect(() => {
    let cancelled = false

    getInsightSummary('week')
      .then((summary: InsightSummary) => {
        if (!cancelled) setWeek(summary.dailyActivity)
      })
      .catch(() => {})

    getBadges()
      .then((response) => {
        if (!cancelled) setBadges(response)
      })
      .catch(() => {})

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
    Alert.alert('Usar freeze', 'Esto congelara tu racha actual por 1 dia. ¿Continuar?', [
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
    ])
  }, [overview])

  const userName = user?.email?.split('@')[0] ?? 'Usuario'

  if (loading && !overview) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    )
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      contentInsetAdjustmentBehavior="automatic"
    >
      <Text style={styles.greeting}>Hola, {userName}</Text>

      {overview ? (
        <LevelLine
          level={overview.level}
          totalXp={overview.totalXp}
          xpToLevel={overview.xpToLevel}
        />
      ) : (
        <Pressable
          onPress={() => void loadOverview()}
          style={({ pressed }) => [styles.overviewError, pressed && styles.pressed]}
          accessibilityRole="button"
          accessibilityLabel="Reintentar carga de progreso"
        >
          <MaterialCommunityIcons
            name="cloud-off-outline"
            size={16}
            color={colors.mutedForeground}
          />
          <Text style={styles.overviewErrorText}>{overviewError}</Text>
          <Text style={styles.overviewErrorAction}>Reintentar</Text>
        </Pressable>
      )}

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
            <MaterialCommunityIcons name="book-open-variant" size={26} color={colors.onPrimary} />
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

      <Pressable
        style={({ pressed }) => [styles.secondaryCta, pressed && styles.pressed]}
        onPress={() => router.push('/(tabs)/sesiones/nueva')}
        accessibilityRole="button"
        accessibilityLabel="Crear sesión de estudio"
      >
        <MaterialCommunityIcons name="plus-circle-outline" size={24} color={colors.primary} />
        <View style={styles.ctaText}>
          <Text style={styles.secondaryTitle}>Sesión de estudio</Text>
          <Text style={styles.secondarySubtitle}>Elige temas y cuántas preguntas</Text>
        </View>
        <MaterialCommunityIcons name="chevron-right" size={20} color={colors.mutedForeground} />
      </Pressable>

      {overview ? (
        <View style={styles.section}>
          <SectionHeader title="Tu día" />
          <TodayCard
            streak={overview.streak}
            dailyGoal={overview.dailyGoal}
            week={week}
            freezing={freezing}
            onFreeze={() => void handleFreeze()}
            onOpenActivity={() => router.push('/(tabs)/progreso/actividad')}
          />
        </View>
      ) : null}

      {badges && badges.total > 0 ? (
        <View style={styles.section}>
          <SectionHeader
            title="Logros"
            onPress={() => router.push('/(tabs)/progreso/logros')}
          />
          <View style={styles.card}>
            <BadgeStrip data={badges} max={4} />
          </View>
        </View>
      ) : null}
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
    backgroundColor: colors.background,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    padding: spacing.lg,
    paddingBottom: spacing.xl,
    gap: spacing.lg,
  },
  greeting: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
  },
  overviewError: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    minHeight: 44,
  },
  overviewErrorText: {
    flex: 1,
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  overviewErrorAction: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.primary,
  },
  pressed: {
    opacity: 0.7,
  },
  ctaCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    padding: spacing.md,
    gap: spacing.md,
    minHeight: 64,
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
    gap: 2,
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
  secondaryCta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    minHeight: 64,
    padding: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
  },
  secondaryTitle: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  secondarySubtitle: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  section: {
    gap: spacing.sm,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.md,
  },
})
