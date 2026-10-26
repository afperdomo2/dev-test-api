import { useCallback, useEffect, useState } from 'react'
import { useRouter } from 'expo-router'
import {
  Pressable,
  RefreshControl,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import type { BadgesResponse, InsightSummary, OverviewResponse, TopicMasteryItem } from '@devtest/shared'

import { BadgeStrip } from '@/core/components/BadgeStrip'
import { LevelLine } from '@/core/components/LevelLine'
import { SectionHeader } from '@/core/components/SectionHeader'
import { getBadges, getOverview } from '@/core/api/services/gamification.service'
import { getInsightSummary } from '@/core/api/services/insights.service'
import { getTopicMastery } from '@/core/api/services/progress.service'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'
import { getErrorDetail } from '@/core/utils/errors'
import { formatDurationMs, toPercent } from '@/core/utils/numbers'

import { ActivitySparkline } from '../components/ActivitySparkline'
import { TopicRingGrid } from '../components/TopicRingGrid'

const HOME_TOPICS = 6

type SectionState<T> = {
  data: T | null
  loading: boolean
  error: string | null
}

/**
 * Cada bloque del dashboard carga por su cuenta: si los logros fallan, los
 * temas siguen visibles. Sin esto, un endpoint caído vaciaría la pestaña.
 */
function useSection<T>(load: () => Promise<T>, fallbackMessage: string) {
  const [state, setState] = useState<SectionState<T>>({
    data: null,
    loading: true,
    error: null,
  })

  const run = useCallback(async () => {
    setState((prev) => ({ ...prev, loading: true, error: null }))
    try {
      const data = await load()
      setState({ data, loading: false, error: null })
    } catch (err) {
      setState({ data: null, loading: false, error: getErrorDetail(err, fallbackMessage) })
    }
  }, [load, fallbackMessage])

  useEffect(() => {
    void run()
  }, [run])

  return { ...state, reload: run }
}

function LoadingBlock({ height }: { height: number }) {
  return (
    <View
      accessible
      accessibilityLabel="Cargando"
      style={[styles.card, styles.loadingBlock, { height }]}
    />
  )
}

function ErrorBlock({
  message,
  section,
  onRetry,
}: {
  message: string
  section: string
  onRetry: () => void
}) {
  return (
    <View style={styles.card}>
      <Text style={styles.errorText}>{message}</Text>
      <Pressable
        onPress={onRetry}
        style={({ pressed }) => [styles.retry, pressed && styles.pressed]}
        accessibilityRole="button"
        accessibilityLabel={`Reintentar ${section}`}
      >
        <Text style={styles.retryText}>Reintentar</Text>
      </Pressable>
    </View>
  )
}

function EmptyBlock({ message }: { message: string }) {
  return (
    <View style={styles.card}>
      <Text style={styles.emptyText}>{message}</Text>
    </View>
  )
}

function ActivitySummary({ summary }: { summary: InsightSummary }) {
  const minutes = formatDurationMs(summary.timeSpentMs)

  return (
    <>
      <Text style={styles.metrics}>
        {`${summary.questionsAnswered} preguntas · ${toPercent(summary.accuracy)}% de acierto · ${minutes}`}
      </Text>
      <ActivitySparkline days={summary.dailyActivity} />
    </>
  )
}

function BadgesSummary({ data }: { data: BadgesResponse }) {
  return <BadgeStrip data={data} detailed />
}

export default function ProgressDashboardScreen() {
  const router = useRouter()

  const loadWeek = useCallback(() => getInsightSummary('week'), [])

  const overview = useSection<OverviewResponse>(getOverview, 'No se pudo cargar tu nivel')
  const topics = useSection<Array<TopicMasteryItem>>(
    getTopicMastery,
    'No se pudo cargar tu dominio por tema',
  )
  const week = useSection<InsightSummary>(loadWeek, 'No se pudo cargar tu actividad')
  const badges = useSection<BadgesResponse>(getBadges, 'No se pudieron cargar los logros')

  const [refreshing, setRefreshing] = useState(false)

  const refresh = useCallback(async () => {
    setRefreshing(true)
    await Promise.allSettled([
      overview.reload(),
      topics.reload(),
      week.reload(),
      badges.reload(),
    ])
    setRefreshing(false)
  }, [overview, topics, week, badges])

  const topicItems = topics.data ?? []

  return (
    <ScrollView
      style={styles.container}
      contentInsetAdjustmentBehavior="automatic"
      contentContainerStyle={styles.content}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={() => void refresh()}
          tintColor={colors.primary}
        />
      }
    >
      {overview.data ? (
        <LevelLine
          level={overview.data.level}
          totalXp={overview.data.totalXp}
          xpToLevel={overview.data.xpToLevel}
        />
      ) : overview.loading ? (
        <LoadingBlock height={32} />
      ) : (
        <ErrorBlock
          message={overview.error ?? ''}
          section="nivel"
          onRetry={() => void overview.reload()}
        />
      )}

      <View style={styles.section}>
        <SectionHeader
          title="Temas"
          onPress={() => router.push('/(tabs)/progreso/temas')}
        />

        {topics.loading ? (
          <LoadingBlock height={224} />
        ) : topics.error ? (
          <ErrorBlock
            message={topics.error}
            section="temas"
            onRetry={() => void topics.reload()}
          />
        ) : topicItems.length === 0 ? (
          <EmptyBlock message="Responde preguntas para ver tu dominio por tema." />
        ) : (
          <View style={styles.card}>
            <TopicRingGrid items={topicItems.slice(0, HOME_TOPICS)} />
            {topicItems.length > HOME_TOPICS ? (
              <Text style={styles.cardFooter}>
                y {topicItems.length - HOME_TOPICS} tema
                {topicItems.length - HOME_TOPICS === 1 ? '' : 's'} más
              </Text>
            ) : null}
          </View>
        )}
      </View>

      <View style={styles.section}>
        <SectionHeader
          title="Actividad"
          actionLabel="Ver detalle"
          onPress={() => router.push('/(tabs)/progreso/actividad')}
        />

        {week.loading ? (
          <LoadingBlock height={96} />
        ) : week.error ? (
          <ErrorBlock
            message={week.error}
            section="actividad"
            onRetry={() => void week.reload()}
          />
        ) : week.data && week.data.questionsAnswered > 0 ? (
          <View style={styles.card}>
            <ActivitySummary summary={week.data} />
          </View>
        ) : (
          <EmptyBlock message="Aún no hay actividad esta semana." />
        )}
      </View>

      <View style={styles.section}>
        <SectionHeader
          title="Logros"
          onPress={() => router.push('/(tabs)/progreso/logros')}
        />

        {badges.loading ? (
          <LoadingBlock height={96} />
        ) : badges.error ? (
          <ErrorBlock
            message={badges.error}
            section="logros"
            onRetry={() => void badges.reload()}
          />
        ) : badges.data && badges.data.total > 0 ? (
          <View style={styles.card}>
            <BadgesSummary data={badges.data} />
          </View>
        ) : (
          <EmptyBlock message="Aún no hay logros disponibles." />
        )}
      </View>
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
    paddingBottom: spacing.xl,
    gap: spacing.lg,
  },
  section: {
    gap: spacing.sm,
  },
  pressed: {
    opacity: 0.7,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.md,
    gap: spacing.sm,
  },
  loadingBlock: {
    backgroundColor: colors.card,
  },
  cardFooter: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  metrics: {
    fontSize: typography.size.sm,
    color: colors.foreground,
  },
  errorText: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  retry: {
    alignSelf: 'flex-start',
    minHeight: 44,
    paddingHorizontal: spacing.md,
    justifyContent: 'center',
    borderRadius: radius.sm,
    borderWidth: 1,
    borderColor: colors.primary,
  },
  retryText: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.primary,
  },
  emptyText: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
})
