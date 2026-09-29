import { useCallback, useEffect, useState } from 'react'
import { ActivityIndicator, Pressable, StyleSheet, Text, View } from 'react-native'
import { useRouter } from 'expo-router'
import type { OverviewResponse } from '@devtest/shared'

import { Screen } from '@/core/components/Screen'
import { useAuthStore } from '@/core/auth/auth.store'
import { getOverview } from '@/core/api/services/gamification.service'
import { createDailySession } from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

export default function HomeScreen() {
  const router = useRouter()
  const user = useAuthStore((s) => s.user)
  const [overview, setOverview] = useState<OverviewResponse | null>(null)
  const [loadingOverview, setLoadingOverview] = useState(true)
  const [creatingDaily, setCreatingDaily] = useState(false)
  const [dailyError, setDailyError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    getOverview()
      .then((res) => {
        if (!cancelled) setOverview(res)
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setLoadingOverview(false)
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
      router.push(`/(tabs)/sesiones/${res.session.id}`)
    } catch (err) {
      setDailyError(getErrorDetail(err, 'No se pudo crear la sesion diaria'))
    } finally {
      setCreatingDaily(false)
    }
  }, [router])

  return (
    <Screen title={`Hola, ${user?.email?.split('@')[0] ?? 'Usuario'}`}>
      <Pressable
        style={[styles.dailyCard, creatingDaily && styles.dailyCardDisabled]}
        onPress={() => void handleDailySession()}
        disabled={creatingDaily}
      >
        {creatingDaily ? (
          <ActivityIndicator color={colors.primary} />
        ) : (
          <>
            <Text style={styles.dailyIcon}>📅</Text>
            <View style={styles.dailyText}>
              <Text style={styles.dailyTitle}>Repasar hoy</Text>
              <Text style={styles.dailySubtitle}>
                Resuelve las preguntas vencidas y mantén tu racha
              </Text>
            </View>
          </>
        )}
      </Pressable>

      {dailyError && <Text style={styles.error}>{dailyError}</Text>}

      {loadingOverview ? (
        <ActivityIndicator color={colors.primary} style={styles.loader} />
      ) : overview ? (
        <View style={styles.statsGrid}>
          <View style={styles.statCard}>
            <Text style={styles.statValue}>{overview.totalXp}</Text>
            <Text style={styles.statLabel}>XP Total</Text>
          </View>
          <View style={styles.statCard}>
            <Text style={styles.statValue}>Nv. {overview.level}</Text>
            <Text style={styles.statLabel}>Nivel</Text>
          </View>
          <View style={styles.statCard}>
            <Text style={styles.statValue}>{overview.streak.current}</Text>
            <Text style={styles.statLabel}>Racha (dias)</Text>
          </View>
          <View style={styles.statCard}>
            <Text style={styles.statValue}>
              {overview.dailyGoal.current}/{overview.dailyGoal.target}
            </Text>
            <Text style={styles.statLabel}>Meta diaria</Text>
          </View>
        </View>
      ) : null}
    </Screen>
  )
}

const styles = StyleSheet.create({
  dailyCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.lg,
    borderWidth: 1,
    borderColor: colors.primary,
    gap: spacing.md,
  },
  dailyCardDisabled: {
    opacity: 0.6,
    justifyContent: 'center',
  },
  dailyIcon: {
    fontSize: 32,
  },
  dailyText: {
    flex: 1,
    gap: spacing.xs,
  },
  dailyTitle: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  dailySubtitle: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  error: {
    fontSize: typography.size.sm,
    color: colors.destructive,
  },
  loader: {
    marginVertical: spacing.xl,
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  statCard: {
    width: '48%',
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    alignItems: 'center',
    gap: spacing.xs,
  },
  statValue: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.primary,
  },
  statLabel: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
