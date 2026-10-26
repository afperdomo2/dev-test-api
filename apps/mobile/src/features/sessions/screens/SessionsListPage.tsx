import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  ActivityIndicator,
  Alert,
  FlatList,
  Pressable,
  RefreshControl,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import { useRouter } from 'expo-router'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import type { BottomSheetModal } from '@gorhom/bottom-sheet'
import type { Session, SessionDifficulty, SessionStatus } from '@devtest/shared'
import {
  SESSION_DIFFICULTY_LABELS,
  SESSION_MODE_LABELS,
  SESSION_STATUS_LABELS,
  SESSION_STATUS_FILTERS,
} from '@devtest/shared'

import { listSessions, deleteSession, finishSession } from '@/core/api/services/sessions.service'
import { SessionActionsSheet } from '@/features/sessions/components/SessionActionsSheet'
import { formatRelativeTime } from '@/core/utils/dates'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const PER_PAGE = 20

/**
 * Los nombres de color viven en `SESSION_STATUS_COLORS` (shared); aquí solo se
 * traducen a design tokens del tema dark. `Record<SessionStatus, string>` hace
 * que el compilador avise si aparece un estado nuevo.
 */
const STATUS_COLORS: Record<SessionStatus, string> = {
  in_progress: colors.accent,
  completed: colors.primary,
  cancelled: colors.mutedForeground,
}

const STATUS_ICON_NAMES: Record<SessionStatus, keyof typeof MaterialCommunityIcons.glyphMap> = {
  in_progress: 'progress-clock',
  completed: 'check-circle',
  cancelled: 'close-circle',
}

const MODE_ICON_NAMES: Record<string, keyof typeof MaterialCommunityIcons.glyphMap> = {
  review: 'book-open-variant',
  daily: 'calendar-star',
}

const DIFFICULTY_ICON_NAMES: Record<SessionDifficulty, keyof typeof MaterialCommunityIcons.glyphMap> =
  {
    random: 'shuffle-variant',
    beginner: 'signal-cellular-1',
    intermediate: 'signal-cellular-2',
    advanced: 'signal-cellular-3',
  }

const DIFFICULTY_COLORS: Record<SessionDifficulty, string> = {
  random: colors.mutedForeground,
  beginner: colors.primary,
  intermediate: colors.accent,
  advanced: colors.destructive,
}

type LoadMode = 'initial' | 'refresh' | 'more'

export default function SessionsListPage() {
  const router = useRouter()
  const bottomSheetRef = useRef<BottomSheetModal>(null)
  const [sessions, setSessions] = useState<Session[]>([])
  const [status, setStatus] = useState<SessionStatus | undefined>(undefined)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(false)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)

  const load = useCallback(
    async (targetPage: number, targetStatus: SessionStatus | undefined, mode: LoadMode) => {
      if (mode === 'refresh') setRefreshing(true)
      else if (mode === 'more') setLoadingMore(true)
      else setLoading(true)
      setError(null)
      try {
        const res = await listSessions(targetPage, PER_PAGE, targetStatus)
        setSessions((prev) => (targetPage === 1 ? res.data : [...prev, ...res.data]))
        setPage(targetPage)
        setHasMore(res.data.length > 0 && targetPage * PER_PAGE < res.meta.total)
      } catch (err) {
        setError(getErrorDetail(err, 'No se pudieron cargar las sesiones'))
      } finally {
        setLoading(false)
        setLoadingMore(false)
        setRefreshing(false)
      }
    },
    [],
  )

  useEffect(() => {
    void load(1, status, 'initial')
  }, [load, status])

  const resumeSession = useMemo(
    () => sessions.find((session) => session.status === 'in_progress') ?? null,
    [sessions],
  )

  const openActions = useCallback((session: Session) => {
    setSelectedSession(session)
    bottomSheetRef.current?.present()
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

  const handleEndReached = useCallback(() => {
    if (!hasMore || loading || loadingMore) return
    void load(page + 1, status, 'more')
  }, [hasMore, loading, loadingMore, page, status, load])

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Sesiones</Text>

      {resumeSession && (
        <Pressable
          style={({ pressed }) => [styles.resumeBanner, pressed && styles.pressed]}
          onPress={() => handleContinue(resumeSession.id)}
          accessibilityRole="button"
          accessibilityLabel={`Continuar sesion ${resumeSession.name}, ${resumeSession.answerCount} de ${resumeSession.questionLimit ?? 'ilimitadas'} preguntas respondidas`}
        >
          <MaterialCommunityIcons name="play-circle" size={22} color={colors.accent} />
          <View style={styles.resumeBody}>
            <Text style={styles.resumeTitle} numberOfLines={1}>
              Continuar: {resumeSession.name}
            </Text>
            <Text style={styles.resumeMeta}>
              {resumeSession.answerCount}
              {resumeSession.questionLimit ? `/${resumeSession.questionLimit}` : ''} preguntas
              respondidas
            </Text>
          </View>
          <View style={styles.resumeButton}>
            <Text style={styles.resumeButtonLabel}>Continuar</Text>
          </View>
        </Pressable>
      )}

      <ScrollView
        horizontal
        style={styles.filters}
        contentContainerStyle={styles.filtersContent}
        showsHorizontalScrollIndicator={false}
        alwaysBounceHorizontal={false}
        accessibilityLabel="Filtros por estado"
      >
        {SESSION_STATUS_FILTERS.map((filter) => {
          const active = status === filter.value
          return (
            <Pressable
              key={filter.value ?? 'all'}
              style={({ pressed }) => [
                styles.chip,
                active && styles.chipActive,
                pressed && styles.chipPressed,
              ]}
              onPress={() => setStatus(filter.value)}
              hitSlop={{ top: 4, bottom: 4 }}
              accessibilityRole="button"
              accessibilityState={{ selected: active }}
              accessibilityLabel={`Filtrar por ${filter.title}`}
            >
              <Text
                style={[styles.chipText, active && styles.chipTextActive]}
                numberOfLines={1}
              >
                {filter.title}
              </Text>
            </Pressable>
          )
        })}
      </ScrollView>

      {loading ? (
        <View style={styles.center}>
          <ActivityIndicator size="large" color={colors.primary} />
        </View>
      ) : error ? (
        <View style={styles.center}>
          <Text style={styles.error}>{error}</Text>
          <Pressable
            style={({ pressed }) => [styles.retryButton, pressed && styles.retryButtonPressed]}
            onPress={() => void load(1, status, 'initial')}
            accessibilityRole="button"
            accessibilityLabel="Reintentar carga de sesiones"
          >
            <Text style={styles.retryLabel}>Reintentar</Text>
          </Pressable>
        </View>
      ) : (
        <FlatList
          data={sessions}
          keyExtractor={(item) => item.id}
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={() => void load(1, status, 'refresh')}
              tintColor={colors.primary}
            />
          }
          contentContainerStyle={sessions.length === 0 ? styles.emptyContainer : styles.list}
          onEndReached={handleEndReached}
          onEndReachedThreshold={0.4}
          ListFooterComponent={
            loadingMore ? (
              <ActivityIndicator
                style={styles.footerSpinner}
                size="small"
                color={colors.primary}
              />
            ) : null
          }
          ListEmptyComponent={
            <View style={styles.empty}>
              <MaterialCommunityIcons
                name="play-circle-outline"
                size={48}
                color={colors.mutedForeground}
              />
              <Text style={styles.emptyTitle}>
                {status ? 'Sin sesiones en este estado' : 'No hay sesiones'}
              </Text>
              <Text style={styles.emptySubtitle}>
                {status
                  ? 'Prueba con otro filtro o crea una nueva sesion'
                  : 'Crea tu primera sesion para comenzar a estudiar'}
              </Text>
            </View>
          }
          renderItem={({ item }) => {
            const statusColor = STATUS_COLORS[item.status] ?? colors.mutedForeground
            const statusIconName = STATUS_ICON_NAMES[item.status] ?? 'help-circle'
            const modeIconName = MODE_ICON_NAMES[item.mode] ?? 'book-open-variant'
            const difficultyIconName = DIFFICULTY_ICON_NAMES[item.difficulty] ?? 'shuffle-variant'
            const difficultyColor = DIFFICULTY_COLORS[item.difficulty] ?? colors.mutedForeground
            const difficultyLabel = SESSION_DIFFICULTY_LABELS[item.difficulty] ?? item.difficulty
            const statusLabel = SESSION_STATUS_LABELS[item.status] ?? item.status
            const modeLabel = SESSION_MODE_LABELS[item.mode] ?? item.mode
            const showScore = item.status === 'completed'
            const progressLabel = item.questionLimit
              ? `${item.answerCount}/${item.questionLimit}`
              : `${item.answerCount}`

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
                accessibilityLabel={`Sesion ${item.name}, ${statusLabel}, ${modeLabel}, ${difficultyLabel}, ${item.answerCount} respuestas`}
              >
                <View style={[styles.accentDot, { backgroundColor: statusColor }]} />

                <View style={styles.cardBody}>
                  <View style={styles.cardHeader}>
                    <Text style={styles.cardName} numberOfLines={1}>
                      {item.name}
                    </Text>
                    <Pressable
                      style={({ pressed }) => [styles.menuButton, pressed && styles.menuButtonPressed]}
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
                      <Text style={[styles.statusText, { color: statusColor }]}>{statusLabel}</Text>
                    </View>
                    <View style={styles.modeBadge}>
                      <MaterialCommunityIcons
                        name={modeIconName}
                        size={12}
                        color={colors.primary}
                      />
                      <Text style={styles.modeText}>{modeLabel}</Text>
                    </View>
                    <View style={styles.difficultyBadge}>
                      <MaterialCommunityIcons
                        name={difficultyIconName}
                        size={12}
                        color={difficultyColor}
                      />
                      <Text style={[styles.difficultyText, { color: difficultyColor }]}>
                        {difficultyLabel}
                      </Text>
                    </View>
                  </View>

                  <View style={styles.cardMeta}>
                    <View style={styles.cardMetaLeft}>
                      <Text style={styles.metaText}>{formatRelativeTime(item.createdAt)}</Text>
                      <Text style={styles.metaDot}>·</Text>
                      <Text style={styles.metaText}>{progressLabel} respuestas</Text>
                      {showScore && (
                        <>
                          <Text style={styles.metaDot}>·</Text>
                          <Text style={styles.metaScore}>{item.score} pts</Text>
                        </>
                      )}
                    </View>
                    <MaterialCommunityIcons
                      name="chevron-right"
                      size={18}
                      color={colors.mutedForeground}
                    />
                  </View>
                </View>
              </Pressable>
            )
          }}
        />
      )}

      <SessionActionsSheet
        ref={bottomSheetRef}
        session={selectedSession}
        onContinue={handleContinue}
        onFinish={(id) => void handleFinish(id)}
        onReview={handleReview}
        onDelete={(id) => void handleDelete(id)}
      />
    </View>
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
  pressed: {
    opacity: 0.7,
  },
  resumeBanner: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    marginHorizontal: spacing.lg,
    marginBottom: spacing.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    minHeight: 56,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: `${colors.accent}40`,
    backgroundColor: `${colors.accent}14`,
  },
  resumeBody: {
    flex: 1,
    gap: 2,
  },
  resumeTitle: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  resumeMeta: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  resumeButton: {
    minHeight: 44,
    minWidth: 88,
    paddingHorizontal: spacing.md,
    borderRadius: radius.sm,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  resumeButtonLabel: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  filters: {
    flexGrow: 0,
    flexShrink: 0,
    height: 36,
    marginBottom: spacing.md,
  },
  filtersContent: {
    gap: spacing.sm,
    paddingHorizontal: spacing.lg,
    paddingRight: spacing.sm,
    alignItems: 'center',
  },
  chip: {
    flexShrink: 0,
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.sm + 4,
    borderRadius: radius.sm,
    borderWidth: 1,
    borderColor: colors.border,
    height: 36,
    justifyContent: 'center',
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
  center: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: spacing.xl,
    gap: spacing.md,
  },
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
  },
  retryButton: {
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
    minHeight: 44,
    justifyContent: 'center',
  },
  retryButtonPressed: {
    opacity: 0.85,
  },
  retryLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  list: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  footerSpinner: {
    paddingVertical: spacing.md,
  },
  emptyContainer: {
    flexGrow: 1,
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
    paddingVertical: spacing.sm + 2,
    paddingHorizontal: spacing.sm + 4,
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
    marginRight: -spacing.sm,
    marginVertical: -spacing.md,
  },
  menuButtonPressed: {
    opacity: 0.6,
  },
  cardBadges: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
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
  difficultyBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    backgroundColor: colors.muted,
    gap: spacing.xs,
  },
  difficultyText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
  },
  cardMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: spacing.sm,
  },
  cardMetaLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flexShrink: 1,
  },
  metaText: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  metaDot: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
    marginHorizontal: spacing.xs,
  },
  metaScore: {
    fontSize: typography.size.xs,
    color: colors.accent,
    fontWeight: typography.weight.medium,
  },
})
