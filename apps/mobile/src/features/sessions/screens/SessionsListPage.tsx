import { useCallback, useEffect, useState } from 'react'
import {
  ActivityIndicator,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import { useRouter } from 'expo-router'
import type { Session } from '@devtest/shared'
import { SESSION_STATUS_LABELS, SESSION_MODE_LABELS } from '@devtest/shared'

import { Screen } from '@/core/components/Screen'
import { listSessions, deleteSession } from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

export default function SessionsListPage() {
  const router = useRouter()
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async (isRefresh = false) => {
    if (isRefresh) setRefreshing(true)
    else setLoading(true)
    setError(null)
    try {
      const res = await listSessions(1, 50)
      setSessions(res.data)
    } catch (err) {
      setError(getErrorDetail(err, 'No se pudieron cargar las sesiones'))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const handleDelete = useCallback(
    async (id: string) => {
      try {
        await deleteSession(id)
        setSessions((prev) => prev.filter((s) => s.id !== id))
      } catch {
        // silenciar — el usuario puede reintentar
      }
    },
    [],
  )

  const handleNewSession = useCallback(() => {
    // TODO: Fase 7.4 — abrir diálogo de creación
  }, [])

  if (loading) {
    return (
      <Screen title="Sesiones">
        <ActivityIndicator size="large" color={colors.primary} />
      </Screen>
    )
  }

  if (error) {
    return (
      <Screen title="Sesiones">
        <Text style={styles.error}>{error}</Text>
        <Pressable style={styles.retryButton} onPress={() => void load()}>
          <Text style={styles.retryLabel}>Reintentar</Text>
        </Pressable>
      </Screen>
    )
  }

  return (
    <Screen title="Sesiones">
      <FlatList
        data={sessions}
        keyExtractor={(item) => item.id}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={() => void load(true)}
            tintColor={colors.primary}
          />
        }
        contentContainerStyle={sessions.length === 0 ? styles.emptyContainer : undefined}
        ListEmptyComponent={
          <View style={styles.empty}>
            <Text style={styles.emptyTitle}>No hay sesiones</Text>
            <Text style={styles.emptySubtitle}>Crea tu primera sesion para comenzar a estudiar</Text>
          </View>
        }
        renderItem={({ item }) => (
          <Pressable
            style={styles.card}
            onPress={() => router.push(`/(tabs)/sesiones/${item.id}`)}
          >
            <View style={styles.cardHeader}>
              <Text style={styles.cardName} numberOfLines={1}>
                {item.name}
              </Text>
              <Text style={styles.cardMode}>{SESSION_MODE_LABELS[item.mode]}</Text>
            </View>
            <View style={styles.cardMeta}>
              <Text style={styles.cardStatus}>{SESSION_STATUS_LABELS[item.status]}</Text>
              <Text style={styles.cardScore}>{item.score} pts</Text>
              <Text style={styles.cardQuestions}>
                {item.answerCount}/{item.questionLimit ?? '∞'}
              </Text>
            </View>
          </Pressable>
        )}
      />
    </Screen>
  )
}

const styles = StyleSheet.create({
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
  },
  retryButton: {
    alignSelf: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
  },
  retryLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: spacing.xl,
  },
  empty: {
    alignItems: 'center',
    gap: spacing.xs,
  },
  emptyTitle: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  emptySubtitle: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
    textAlign: 'center',
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    marginBottom: spacing.sm,
    borderWidth: 1,
    borderColor: colors.border,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.xs,
  },
  cardName: {
    flex: 1,
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
    marginRight: spacing.sm,
  },
  cardMode: {
    fontSize: typography.size.xs,
    color: colors.primary,
    fontWeight: typography.weight.medium,
  },
  cardMeta: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  cardStatus: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  cardScore: {
    fontSize: typography.size.xs,
    color: colors.accent,
    fontWeight: typography.weight.medium,
  },
  cardQuestions: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
