import { useCallback, useEffect, useState } from 'react'
import { ActivityIndicator, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native'
import type { InsightPeriod, InsightSummary } from '@devtest/shared'

import { getInsightSummary } from '@/core/api/services/insights.service'
import { getErrorDetail } from '@/core/utils/errors'
import { formatDurationMs, toPercent } from '@/core/utils/numbers'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

import { ActivityHeatmap } from '../components/ActivityHeatmap'
import { DailyActivityBar } from '../components/DailyActivityBar'
import { PeriodToggle } from '../components/PeriodToggle'
import { StatRow } from '../components/StatRow'
import { TopicCompactRow } from '../components/TopicCompactRow'

const TOP_TOPICS = 5
const WEAK_TOPICS = 3

function formatRange(from: string, to: string): string {
  const format = (value: string) =>
    new Date(value).toLocaleDateString('es-CO', {
      day: 'numeric',
      month: 'short',
      timeZone: 'UTC',
    })

  return `${format(from)} – ${format(to)}`
}

export default function ActivityScreen() {
  const [period, setPeriod] = useState<InsightPeriod>('week')
  const [summary, setSummary] = useState<InsightSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async (value: InsightPeriod) => {
    setLoading(true)
    setError(null)

    try {
      setSummary(await getInsightSummary(value))
    } catch (err) {
      setError(getErrorDetail(err, 'No se pudo cargar el informe'))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load(period)
  }, [load, period])

  const isEmpty = summary !== null && summary.sessions === 0 && summary.questionsAnswered === 0
  const maxDaily = Math.max(
    1,
    ...(summary?.dailyActivity.map((day) => day.questionsAnswered) ?? []),
  )
  const topTopics = summary?.topTopics.slice(0, TOP_TOPICS) ?? []
  const weakTopics = summary?.weakTopics.slice(0, WEAK_TOPICS) ?? []

  return (
    <ScrollView
      style={styles.container}
      contentInsetAdjustmentBehavior="automatic"
      contentContainerStyle={styles.content}
    >
      <PeriodToggle value={period} onChange={setPeriod} />

      {loading ? (
        <View style={styles.center}>
          <ActivityIndicator color={colors.primary} />
        </View>
      ) : error ? (
        <View style={styles.center}>
          <Text style={styles.error}>{error}</Text>
          <Pressable
            onPress={() => void load(period)}
            style={({ pressed }) => [styles.retry, pressed && styles.retryPressed]}
            accessibilityRole="button"
            accessibilityLabel="Reintentar"
          >
            <Text style={styles.retryText}>Reintentar</Text>
          </Pressable>
        </View>
      ) : summary ? (
        <>
          <Text style={styles.range}>{formatRange(summary.from, summary.to)}</Text>

          {isEmpty ? (
            <Text style={styles.empty}>No hay actividad en este periodo.</Text>
          ) : (
            <>
              <View style={styles.hero}>
                <Text style={styles.heroValue}>{summary.questionsAnswered}</Text>
                <Text style={styles.heroLabel}>
                  preguntas en {summary.sessions} {summary.sessions === 1 ? 'sesión' : 'sesiones'}
                </Text>
              </View>

              <View style={styles.card}>
                <StatRow
                  items={[
                    { value: `${toPercent(summary.accuracy)}%`, label: 'acierto' },
                    { value: formatDurationMs(summary.timeSpentMs), label: 'tiempo' },
                    { value: `+${summary.xpGained}`, label: 'XP ganado' },
                  ]}
                />
              </View>

              <View style={styles.section}>
                <Text style={styles.sectionTitle}>
                  {period === 'week' ? 'Últimos 7 días' : 'Últimos 30 días'}
                </Text>
                {period === 'week' ? (
                  <View style={[styles.card, styles.chart]}>
                    {summary.dailyActivity.map((day) => (
                      <DailyActivityBar
                        key={day.date}
                        date={day.date}
                        value={day.questionsAnswered}
                        max={maxDaily}
                      />
                    ))}
                  </View>
                ) : (
                  <View style={styles.card}>
                    <ActivityHeatmap days={summary.dailyActivity} />
                  </View>
                )}
              </View>

              <View style={styles.section}>
                <Text style={styles.sectionTitle}>Más practicados</Text>
                {topTopics.length === 0 ? (
                  <Text style={styles.muted}>Sin datos en este periodo.</Text>
                ) : (
                  <View style={styles.card}>
                    {topTopics.map((topic, index) => (
                      <View
                        key={topic.topicId}
                        style={index > 0 ? styles.rowDivider : undefined}
                      >
                        <TopicCompactRow
                          topicName={topic.topicName}
                          attempts={topic.attempts}
                          accuracy={topic.accuracy}
                        />
                      </View>
                    ))}
                  </View>
                )}
              </View>

              {weakTopics.length > 0 ? (
                <View style={styles.section}>
                  <Text style={styles.sectionTitle}>Puntos a reforzar</Text>
                  <View style={styles.card}>
                    {weakTopics.map((topic, index) => (
                      <View
                        key={topic.topicId}
                        style={index > 0 ? styles.rowDivider : undefined}
                      >
                        <TopicCompactRow
                          topicName={topic.topicName}
                          attempts={topic.attempts}
                          accuracy={topic.accuracy}
                          isWeak
                          correctCount={topic.correctCount}
                        />
                      </View>
                    ))}
                  </View>
                </View>
              ) : null}
            </>
          )}
        </>
      ) : null}
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
    gap: spacing.md,
  },
  center: {
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.md,
    padding: spacing.xl,
  },
  error: {
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
  range: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  hero: {
    gap: 2,
  },
  heroValue: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
  },
  heroLabel: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  empty: {
    fontSize: typography.size.sm,
    color: colors.accent,
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
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  chart: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    height: 104,
    gap: spacing.xs,
    paddingTop: spacing.md,
  },
  muted: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  rowDivider: {
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
})
