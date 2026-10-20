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

import { TopicMasteryRow } from '../components/TopicMasteryRow'

export default function TopicMasteryScreen() {
  const [items, setItems] = useState<Array<TopicMasteryItem> | null>(null)
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
        <View style={styles.list}>
          {list.map((item) => (
            <TopicMasteryRow key={item.topicId} item={item} />
          ))}
        </View>
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
    gap: spacing.sm,
  },
  list: {
    gap: spacing.sm,
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
