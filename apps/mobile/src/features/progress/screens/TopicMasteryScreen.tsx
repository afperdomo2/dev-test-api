import { useCallback, useEffect, useState } from 'react'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import {
  ActivityIndicator,
  Pressable,
  RefreshControl,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import type { TopicMasteryItem } from '@devtest/shared'

import { getTopicMastery } from '@/core/api/services/progress.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

import { TopicRingGrid } from '../components/TopicRingGrid'

type Filter = 'all' | 'weak'

const FILTERS: Array<{ value: Filter; title: string }> = [
  { value: 'all', title: 'Todos' },
  { value: 'weak', title: 'Débiles' },
]

export default function TopicMasteryScreen() {
  const [items, setItems] = useState<Array<TopicMasteryItem> | null>(null)
  const [filter, setFilter] = useState<Filter>('all')
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async (mode: 'initial' | 'refresh') => {
    if (mode === 'initial') setLoading(true)
    else setRefreshing(true)
    setError(null)

    try {
      setItems(await getTopicMastery())
    } catch (err) {
      setError(getErrorDetail(err, 'No se pudo cargar tu dominio por tema'))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [])

  useEffect(() => {
    void load('initial')
  }, [load])

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator color={colors.primary} />
      </View>
    )
  }

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
          onPress={() => void load('initial')}
          style={({ pressed }) => [styles.retry, pressed && styles.retryPressed]}
          accessibilityRole="button"
          accessibilityLabel="Reintentar"
        >
          <Text style={styles.retryText}>Reintentar</Text>
        </Pressable>
      </View>
    )
  }

  const list = items ?? []
  const weakCount = list.filter((item) => item.isWeak).length
  const visible = filter === 'weak' ? list.filter((item) => item.isWeak) : list

  return (
    <ScrollView
      style={styles.container}
      contentInsetAdjustmentBehavior="automatic"
      contentContainerStyle={styles.content}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={() => void load('refresh')}
          tintColor={colors.primary}
        />
      }
    >
      {list.length === 0 ? (
        <View style={styles.emptyCard}>
          <MaterialCommunityIcons
            name="chart-timeline-variant"
            size={40}
            color={colors.mutedForeground}
          />
          <Text style={styles.emptyTitle}>Aún no hay respuestas registradas</Text>
          <Text style={styles.emptyText}>Responde preguntas para ver tu dominio por tema.</Text>
        </View>
      ) : (
        <>
          <View style={styles.toolbar}>
            <Text style={styles.summary}>
              {weakCount} de {list.length} {list.length === 1 ? 'tema flojo' : 'temas flojos'}
            </Text>

            <View style={styles.filters}>
              {FILTERS.map((option) => {
                const active = option.value === filter
                return (
                  <Pressable
                    key={option.value}
                    onPress={() => setFilter(option.value)}
                    hitSlop={{ top: 4, bottom: 4 }}
                    style={({ pressed }) => [
                      styles.chip,
                      active && styles.chipActive,
                      pressed && styles.chipPressed,
                    ]}
                    accessibilityRole="button"
                    accessibilityState={{ selected: active }}
                    accessibilityLabel={`Filtrar por ${option.title}`}
                  >
                    <Text style={[styles.chipText, active && styles.chipTextActive]}>
                      {option.title}
                    </Text>
                  </Pressable>
                )
              })}
            </View>
          </View>

          {visible.length === 0 ? (
            <View style={styles.emptyCard}>
              <MaterialCommunityIcons name="check-decagram" size={32} color={colors.primary} />
              <Text style={styles.emptyTitle}>Sin temas débiles</Text>
              <Text style={styles.emptyText}>Ningún tema baja del 60% de acierto. Buen trabajo.</Text>
            </View>
          ) : (
            <View style={styles.card}>
              <TopicRingGrid items={visible} />
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
  content: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  toolbar: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: spacing.sm,
  },
  summary: {
    flex: 1,
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  filters: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  chip: {
    minHeight: 36,
    paddingHorizontal: spacing.sm + 4,
    justifyContent: 'center',
    borderRadius: radius.sm,
    borderWidth: 1,
    borderColor: colors.border,
  },
  chipActive: {
    backgroundColor: colors.primary,
    borderColor: colors.primary,
  },
  chipPressed: {
    opacity: 0.7,
  },
  chipText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
    color: colors.mutedForeground,
  },
  chipTextActive: {
    color: colors.onPrimary,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.md,
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
