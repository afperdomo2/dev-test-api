import { useCallback, useEffect, useRef, useState } from 'react'
import {
  ActivityIndicator,
  Alert,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import { GestureHandlerRootView } from 'react-native-gesture-handler'
import { useRouter } from 'expo-router'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import type BottomSheet from '@gorhom/bottom-sheet'
import type { Session } from '@devtest/shared'
import { SESSION_STATUS_LABELS, SESSION_MODE_LABELS } from '@devtest/shared'

import { listSessions, deleteSession, finishSession } from '@/core/api/services/sessions.service'
import { SessionActionsSheet } from '@/features/sessions/components/SessionActionsSheet'
import { formatRelativeTime } from '@/core/utils/dates'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const STATUS_COLORS: Record<string, string> = {
  in_progress: colors.accent,
  completed: colors.primary,
  cancelled: colors.destructive,
}

const STATUS_ICON_NAMES: Record<string, keyof typeof MaterialCommunityIcons.glyphMap> = {
  in_progress: 'progress-clock',
  completed: 'check-circle',
  cancelled: 'close-circle',
}

const MODE_ICON_NAMES: Record<string, keyof typeof MaterialCommunityIcons.glyphMap> = {
  review: 'book-open-variant',
  daily: 'calendar-star',
}

export default function SessionsListPage() {
  const router = useRouter()
  const bottomSheetRef = useRef<BottomSheet>(null)
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)

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

  const openActions = useCallback((session: Session) => {
    setSelectedSession(session)
    bottomSheetRef.current?.expand()
  }, [])

  const handleContinue = useCallback(
    (id: string) => {
      router.push(`/(tabs)/sesiones/${id}`)
    },
    [router],
  )

  const handleFinish = useCallback(
    async (id: string) => {
      try {
        await finishSession(id)
        setSessions((prev) =>
          prev.map((s) => (s.id === id ? { ...s, status: 'completed' as const } : s)),
        )
        router.push(`/(tabs)/sesiones/${id}/summary`)
      } catch (err) {
        Alert.alert('Error', getErrorDetail(err, 'No se pudo finalizar la sesion'))
      }
    },
    [router],
  )

  const handleReview = useCallback(
    (id: string) => {
      router.push(`/(tabs)/sesiones/${id}/summary`)
    },
    [router],
  )

  const handleDelete = useCallback(async (id: string) => {
    try {
      await deleteSession(id)
      setSessions((prev) => prev.filter((s) => s.id !== id))
    } catch (err) {
      Alert.alert('Error', getErrorDetail(err, 'No se pudo eliminar la sesion'))
    }
  }, [])

  if (loading) {
    return (
      <GestureHandlerRootView style={styles.container}>
        <Text style={styles.title}>Sesiones</Text>
        <ActivityIndicator size="large" color={colors.primary} />
      </GestureHandlerRootView>
    )
  }

  if (error) {
    return (
      <GestureHandlerRootView style={styles.container}>
        <Text style={styles.title}>Sesiones</Text>
        <Text style={styles.error}>{error}</Text>
        <Pressable
          style={styles.retryButton}
          onPress={() => void load()}
          accessibilityRole="button"
          accessibilityLabel="Reintentar carga de sesiones"
        >
          <Text style={styles.retryLabel}>Reintentar</Text>
        </Pressable>
      </GestureHandlerRootView>
    )
  }

  return (
    <GestureHandlerRootView style={styles.container}>
      <Text style={styles.title}>Sesiones</Text>
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
        contentContainerStyle={sessions.length === 0 ? styles.emptyContainer : styles.list}
        ListEmptyComponent={
          <View style={styles.empty}>
            <Text style={styles.emptyTitle}>No hay sesiones</Text>
            <Text style={styles.emptySubtitle}>Crea tu primera sesion para comenzar a estudiar</Text>
          </View>
        }
        renderItem={({ item }) => {
          const statusColor = STATUS_COLORS[item.status] ?? colors.mutedForeground
          const statusIconName = STATUS_ICON_NAMES[item.status] ?? 'help-circle'
          const modeIconName = MODE_ICON_NAMES[item.mode] ?? 'book-open-variant'
          const statusLabel = SESSION_STATUS_LABELS[item.status] ?? item.status
          const modeLabel = SESSION_MODE_LABELS[item.mode] ?? item.mode
          return (
            <Pressable
              style={({ pressed }) => [styles.card, pressed && styles.cardPressed]}
              onPress={() =>
                router.push(
                  item.status === 'completed'
                    ? `/(tabs)/sesiones/${item.id}/summary`
                    : `/(tabs)/sesiones/${item.id}`,
                )
              }
              onLongPress={() => openActions(item)}
              accessibilityRole="button"
              accessibilityLabel={`Sesion ${item.name}, ${statusLabel}, ${modeLabel}, ${item.answerCount} respuestas${item.score ? `, ${item.score} puntos` : ''}`}
              accessibilityState={{
                busy: false,
              }}
            >
              <View style={[styles.accentDot, { backgroundColor: statusColor }]} />

              <View style={styles.cardBody}>
                <View style={styles.cardHeader}>
                  <Text style={styles.cardName} numberOfLines={1}>
                    {item.name}
                  </Text>
                  <Pressable
                    style={styles.menuButton}
                    onPress={() => openActions(item)}
                    hitSlop={12}
                    accessibilityRole="button"
                    accessibilityLabel={`Abrir acciones de ${item.name}`}
                  >
                    <MaterialCommunityIcons
                      name="dots-vertical"
                      size={20}
                      color={colors.mutedForeground}
                    />
                  </Pressable>
                </View>

                <View style={styles.cardBadges}>
                  <View style={[styles.statusBadge, { backgroundColor: `${statusColor}20` }]}>
                    <MaterialCommunityIcons
                      name={statusIconName}
                      size={12}
                      color={statusColor}
                    />
                    <Text style={[styles.statusText, { color: statusColor }]}>
                      {statusLabel}
                    </Text>
                  </View>
                  <View style={styles.modeBadge}>
                    <MaterialCommunityIcons
                      name={modeIconName}
                      size={12}
                      color={colors.primary}
                    />
                    <Text style={styles.modeText}>{modeLabel}</Text>
                  </View>
                </View>

                <View style={styles.cardMeta}>
                  <Text style={styles.metaDate}>{formatRelativeTime(item.createdAt)}</Text>
                  <Text style={styles.metaDot}>·</Text>
                  <Text style={styles.metaText}>{item.answerCount} respuestas</Text>
                  {item.score != null && item.score > 0 && (
                    <>
                      <Text style={styles.metaDot}>·</Text>
                      <Text style={styles.metaScore}>{item.score} pts</Text>
                    </>
                  )}
                </View>
              </View>
            </Pressable>
          )
        }}
      />

      <SessionActionsSheet
        ref={bottomSheetRef}
        session={selectedSession}
        onContinue={handleContinue}
        onFinish={handleFinish}
        onReview={handleReview}
        onDelete={handleDelete}
      />
    </GestureHandlerRootView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    paddingTop: spacing.xl,
  },
  title: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.md,
  },
  list: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
    marginBottom: spacing.md,
  },
  retryButton: {
    alignSelf: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
    minHeight: 44,
    justifyContent: 'center',
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
    flexDirection: 'row',
    backgroundColor: colors.card,
    borderRadius: radius.md,
    marginBottom: spacing.sm,
    borderWidth: 1,
    borderColor: colors.border,
    overflow: 'hidden',
  },
  cardPressed: {
    opacity: 0.7,
    borderColor: colors.primary,
  },
  accentDot: {
    width: 4,
  },
  cardBody: {
    flex: 1,
    padding: spacing.md,
    gap: spacing.xs,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  cardName: {
    flex: 1,
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
    marginRight: spacing.sm,
  },
  menuButton: {
    padding: spacing.xs,
    minWidth: 44,
    minHeight: 44,
    justifyContent: 'center',
    alignItems: 'center',
  },
  cardBadges: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    gap: spacing.xs,
  },
  statusText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
  },
  modeBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    backgroundColor: `${colors.primary}15`,
    gap: spacing.xs,
  },
  modeText: {
    fontSize: typography.size.xs,
    color: colors.primary,
    fontWeight: typography.weight.medium,
  },
  cardMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: spacing.xs,
  },
  metaDate: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  metaDot: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
    marginHorizontal: spacing.xs,
  },
  metaText: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  metaScore: {
    fontSize: typography.size.xs,
    color: colors.accent,
    fontWeight: typography.weight.medium,
  },
})
