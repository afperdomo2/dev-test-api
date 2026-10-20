import { useCallback, useEffect, useState } from 'react'
import {
  ActivityIndicator,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import type { InsightPeriod, InsightSummary } from '@devtest/shared'

import { getInsightSummary } from '@/core/api/services/insights.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

import { DailyActivityBar } from '../components/DailyActivityBar'
import { KpiCard } from '../components/KpiCard'
import { PeriodToggle } from '../components/PeriodToggle'
import { TopicMasteryRow } from '../components/TopicMasteryRow'

function percent(value: number): number {
  return Math.round(value * 100)
}

function formatDuration(ms: number): string {
  const totalMinutes = Math.round(ms / 60000)
  if (totalMinutes < 60) return `${totalMinutes} min`

  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
}

export default function WeeklySummaryScreen() {
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
          <Text style={styles.range}>
            {summary.from.slice(0, 10)} — {summary.to.slice(0, 10)}
          </Text>

          {isEmpty ? <Text style={styles.empty}>No hay actividad en este periodo.</Text> : null}

          <View style={styles.grid}>
            <KpiCard icon="play-circle-outline" label="Sesiones" value={String(summary.sessions)} />
            <KpiCard
              icon="help-circle-outline"
              label="Preguntas"
              value={String(summary.questionsAnswered)}
            />
            <KpiCard icon="target" label="Acierto" value={`${percent(summary.accuracy)}%`} />
            <KpiCard
              icon="timer-outline"
              label="Tiempo"
              value={formatDuration(summary.timeSpentMs)}
            />
            <KpiCard icon="star-outline" label="XP ganado" value={String(summary.xpGained)} />
            <KpiCard icon="fire" label="Racha" value={`${summary.streakDays} días`} />
          </View>

          <Text style={styles.section}>Más practicados</Text>
          {summary.topTopics.length === 0 ? (
            <Text style={styles.muted}>Sin datos en este periodo.</Text>
          ) : (
            summary.topTopics.map((topic) => (
              <TopicMasteryRow
                key={topic.topicId}
                topicName={topic.topicName}
                attempts={topic.attempts}
                accuracy={topic.accuracy}
              />
            ))
          )}

          {summary.weakTopics.length > 0 ? (
            <>
              <Text style={styles.section}>Temas débiles</Text>
              {summary.weakTopics.map((topic) => (
                <TopicMasteryRow
                  key={topic.topicId}
                  topicName={topic.topicName}
                  attempts={topic.attempts}
                  accuracy={topic.accuracy}
                  isWeak
                  correctCount={topic.correctCount}
                />
              ))}
            </>
          ) : null}

          <Text style={styles.section}>Actividad diaria</Text>
          <View style={styles.chart}>
            {summary.dailyActivity.map((day) => (
              <DailyActivityBar
                key={day.date}
                date={day.date}
                value={day.questionsAnswered}
                max={maxDaily}
              />
            ))}
          </View>
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
  empty: {
    fontSize: typography.size.sm,
    color: colors.accent,
  },
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    rowGap: spacing.sm,
  },
  section: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
    marginTop: spacing.sm,
  },
  muted: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  chart: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    height: 120,
    gap: spacing.xs,
  },
})
