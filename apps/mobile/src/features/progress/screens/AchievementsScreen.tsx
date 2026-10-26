import { useCallback, useEffect, useState } from 'react'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native'
import type { Badge, BadgesResponse } from '@devtest/shared'

import { getBadges } from '@/core/api/services/gamification.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

import { BadgeCard } from '../components/BadgeCard'

function progressOf(badge: Badge): number {
  if (badge.target <= 0) return 0
  return Math.min(100, Math.round((badge.current / badge.target) * 100))
}

function Summary({ earned, total }: { earned: number; total: number }) {
  const value = total > 0 ? Math.round((earned / total) * 100) : 0
  const label = `${earned} de ${total} logros`

  return (
    <View style={styles.summary}>
      <View style={styles.summaryHeader}>
        <Text style={styles.summaryText}>{label}</Text>
        <Text style={styles.summaryPercent}>{value}%</Text>
      </View>

      <View
        accessible
        accessibilityRole="progressbar"
        accessibilityLabel={label}
        accessibilityValue={{ min: 0, max: 100, now: value }}
        style={styles.track}
      >
        <View style={[styles.fill, { width: `${value}%` }]} />
      </View>
    </View>
  )
}

function Skeleton() {
  return (
    <View
      accessible
      accessibilityLabel="Cargando logros"
      style={[styles.container, styles.content]}
    >
      <View style={styles.summarySkeleton} />
      <View style={styles.grid}>
        {[0, 1, 2, 3].map((key) => (
          <View key={key} style={styles.cardSkeleton} />
        ))}
      </View>
    </View>
  )
}

/** Dos columnas; las cards se agrupan por estado para que lo pendiente no se pierda. */
function BadgeGrid({ badges }: { badges: Array<Badge> }) {
  return (
    <View style={styles.grid}>
      {badges.map((badge) => (
        <View key={badge.id} style={styles.gridCell}>
          <BadgeCard badge={badge} />
        </View>
      ))}
    </View>
  )
}

export default function AchievementsScreen() {
  const [data, setData] = useState<BadgesResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      setData(await getBadges())
    } catch (err) {
      setError(getErrorDetail(err, 'No se pudieron cargar los logros'))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  if (loading) return <Skeleton />

  if (error) {
    return (
      <View style={styles.center}>
        <MaterialCommunityIcons
          name="cloud-off-outline"
          size={40}
          color={colors.mutedForeground}
        />
        <Text style={styles.errorText}>{error}</Text>
        <Pressable
          onPress={() => void load()}
          style={({ pressed }) => [styles.retry, pressed && styles.retryPressed]}
          accessibilityRole="button"
          accessibilityLabel="Reintentar"
        >
          <Text style={styles.retryText}>Reintentar</Text>
        </Pressable>
      </View>
    )
  }

  const badges = data?.badges ?? []
  const earned = badges.filter((badge) => badge.earned)
  const pending = badges
    .filter((badge) => !badge.earned)
    .sort((a, b) => progressOf(b) - progressOf(a))

  return (
    <ScrollView
      style={styles.container}
      contentInsetAdjustmentBehavior="automatic"
      contentContainerStyle={styles.content}
    >
      {badges.length === 0 ? (
        <View style={styles.emptyCard}>
          <MaterialCommunityIcons
            name="trophy-outline"
            size={40}
            color={colors.mutedForeground}
          />
          <Text style={styles.emptyTitle}>Aún no hay logros disponibles</Text>
          <Text style={styles.emptyText}>
            Sigue practicando: tus logros aparecerán aquí en cuanto el backend los publique.
          </Text>
        </View>
      ) : (
        <>
          <Summary earned={data?.earned ?? earned.length} total={data?.total ?? badges.length} />

          {earned.length > 0 ? (
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Conseguidos</Text>
              <BadgeGrid badges={earned} />
            </View>
          ) : null}

          {pending.length > 0 ? (
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>En progreso</Text>
              <BadgeGrid badges={pending} />
            </View>
          ) : null}
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
  content: {
    padding: spacing.lg,
    gap: spacing.lg,
  },
  section: {
    gap: spacing.sm,
  },
  sectionTitle: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  gridCell: {
    width: '48.5%',
  },
  summary: {
    gap: spacing.sm,
    padding: spacing.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
  },
  summaryHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  summaryText: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.foreground,
  },
  summaryPercent: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
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
    backgroundColor: colors.primary,
  },
  summarySkeleton: {
    height: 62,
    borderRadius: radius.md,
    backgroundColor: colors.card,
  },
  cardSkeleton: {
    width: '48.5%',
    height: 140,
    borderRadius: radius.md,
    backgroundColor: colors.card,
  },
  center: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.md,
    padding: spacing.lg,
  },
  errorText: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
    textAlign: 'center',
  },
  retry: {
    minHeight: 48,
    paddingHorizontal: spacing.lg,
    justifyContent: 'center',
    borderRadius: radius.md,
    backgroundColor: colors.primary,
  },
  retryPressed: {
    opacity: 0.85,
  },
  retryText: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  emptyCard: {
    alignItems: 'center',
    gap: spacing.sm,
    padding: spacing.xl,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
  },
  emptyTitle: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
    textAlign: 'center',
  },
  emptyText: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
    textAlign: 'center',
  },
})
