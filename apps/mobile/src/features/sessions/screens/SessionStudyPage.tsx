import { useCallback, useEffect, useRef, useState } from 'react'
import {
  ActivityIndicator,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { useLocalSearchParams, useRouter } from 'expo-router'
import { randomUUID } from 'expo-crypto'
import { useSafeAreaInsets } from 'react-native-safe-area-context'
import type { NextQuestionItem, SessionAnswer, SessionSummary } from '@devtest/shared'

import { QuestionCard } from '@/features/sessions/components/QuestionCard'
import { useScrollToTopOnChange } from '@/features/sessions/useScrollToTopOnChange'
import {
  getSessionById,
  getSessionSummary,
  getNextQuestion,
  submitAnswer,
  finishSession,
  deleteSession,
} from '@/core/api/services/sessions.service'
import { getErrorDetail, getErrorStatus, isNetworkError } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

// 'empty' = la sesión no tiene preguntas (nació vacía). No es 'finished':
// no se respondió nada, así que no debe mostrarse como "Sesión completada".
type Phase = 'loading' | 'question' | 'result' | 'finished' | 'empty' | 'error'

export default function SessionStudyPage() {
  const { id } = useLocalSearchParams<{ id: string }>()
  const router = useRouter()
  const insets = useSafeAreaInsets()

  const [phase, setPhase] = useState<Phase>('loading')
  const [sessionStatus, setSessionStatus] = useState<string | null>(null)
  const [question, setQuestion] = useState<NextQuestionItem | null>(null)
  const [summary, setSummary] = useState<SessionSummary | null>(null)
  const [selectedOptions, setSelectedOptions] = useState<string[]>([])
  const [lastResult, setLastResult] = useState<SessionAnswer | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitErrorIsNetwork, setSubmitErrorIsNetwork] = useState(false)
  // Qué acción muestra la barra fija. Se mantiene durante el `loading` para que
  // el botón no cambie de etiqueta mientras espera la red.
  const [barAction, setBarAction] = useState<'confirm' | 'next'>('confirm')
  const loadingRef = useRef(false)
  const submittingRef = useRef(false)
  const scrollRef = useRef<ScrollView>(null)
  // UUID estable por pregunta: se reutiliza en los reintentos y se renueva al
  // cargar la siguiente (idempotencia, D-35).
  const clientAnswerIdRef = useRef<string | null>(null)

  const questionNumber = summary ? summary.answerCount + 1 : 1
  const questionLimit = summary?.questionLimit ?? null
  const progress = questionLimit ? questionNumber / questionLimit : 0

  // La siguiente pregunta debe empezar arriba, no donde quedó el scroll de la
  // anterior (acción al final de una pregunta/ explicación larga).
  useScrollToTopOnChange(scrollRef, question?.id)

  const loadSessionData = useCallback(async (sessionId: string) => {
    const [sess, sum] = await Promise.all([
      getSessionById(sessionId).catch(() => null),
      getSessionSummary(sessionId).catch(() => null),
    ])
    if (sess) setSessionStatus(sess.status)
    if (sum) setSummary(sum)
    return { session: sess, summary: sum }
  }, [])

  const loadNext = useCallback(
    async (sessionId: string) => {
      if (loadingRef.current) return
      loadingRef.current = true
      setPhase('loading')
      setError(null)
      setSelectedOptions([])
      setLastResult(null)
      setSubmitError(null)
      setSubmitErrorIsNetwork(false)
      clientAnswerIdRef.current = null
      try {
        const next = await getNextQuestion(sessionId)
        clientAnswerIdRef.current = randomUUID()
        setQuestion(next)
        setBarAction('confirm')
        setPhase('question')
      } catch (err) {
        const { session, summary: sum } = await loadSessionData(sessionId)
        if (!session) {
          setError(getErrorDetail(err, 'No se pudo cargar la sesión'))
          setPhase('error')
          return
        }
        if (getErrorStatus(err) === 404) {
          // 404 = no quedan preguntas. Solo es "completada" si ya respondiste
          // algo; si no, la sesión está vacía y hay que ofrecer una salida.
          setPhase((sum?.answerCount ?? 0) > 0 ? 'finished' : 'empty')
          return
        }
        setError(getErrorDetail(err, 'No se pudo cargar la siguiente pregunta'))
        setPhase('error')
      } finally {
        loadingRef.current = false
      }
    },
    [loadSessionData],
  )

  useEffect(() => {
    if (!id) return
    let cancelled = false
    ;(async () => {
      try {
        const { session } = await loadSessionData(id)
        if (cancelled) return

        if (session?.status === 'completed') {
          setPhase('finished')
          return
        }
        await loadNext(id)
      } catch {
        if (!cancelled) setPhase('error')
      }
    })()
    return () => {
      cancelled = true
    }
  }, [id, loadSessionData, loadNext])

  const handleAnswer = useCallback(
    async (options: string[]) => {
      if (!id || !question || submittingRef.current) return

      // Se genera una vez por pregunta; los reintentos reutilizan el mismo id.
      const clientAnswerId = clientAnswerIdRef.current ?? randomUUID()
      clientAnswerIdRef.current = clientAnswerId

      submittingRef.current = true
      setSelectedOptions(options)
      setSubmitError(null)
      setSubmitErrorIsNetwork(false)
      setPhase('loading')
      try {
        const result = await submitAnswer(id, {
          questionId: question.id,
          selectedOptions: options,
          clientAnswerId,
        })
        setLastResult(result)
        setSummary((prev) =>
          prev ? { ...prev, answerCount: prev.answerCount + 1 } : prev,
        )
        setBarAction('next')
        setPhase('result')
      } catch (err) {
        // No se avanza: el usuario ve el error y puede reintentar sin duplicar.
        const network = isNetworkError(err)
        setBarAction('confirm')
        setPhase('question')
        setSubmitErrorIsNetwork(network)
        setSubmitError(
          network
            ? 'Sin conexión. Revisa tu red e inténtalo de nuevo.'
            : getErrorDetail(err, 'No se pudo enviar la respuesta'),
        )
      } finally {
        submittingRef.current = false
      }
    },
    [id, question],
  )

  const handleNext = useCallback(() => {
    if (!id) return
    if (questionLimit !== null && questionNumber >= questionLimit) {
      setPhase('finished')
      return
    }
    void loadNext(id)
  }, [id, questionNumber, questionLimit, loadNext])

  const handleToggleOption = useCallback(
    (optionId: string) => {
      if (!question || phase === 'result' || phase === 'loading') return
      setSelectedOptions((prev) => {
        if (question.type === 'multiple_choice') {
          return prev.includes(optionId)
            ? prev.filter((id) => id !== optionId)
            : [...prev, optionId]
        }
        return [optionId]
      })
    },
    [question, phase],
  )

  const handleFinish = useCallback(async () => {
    if (!id) return
    if (sessionStatus === 'completed') {
      router.replace(`/(tabs)/sesiones/${id}/summary`)
      return
    }
    try {
      await finishSession(id)
      router.replace(`/(tabs)/sesiones/${id}/summary`)
    } catch (err) {
      const detail = getErrorDetail(err, 'Error al finalizar')
      if (detail.includes('completada')) {
        router.replace(`/(tabs)/sesiones/${id}/summary`)
      }
    }
  }, [id, sessionStatus, router])

  // Una sesión sin preguntas no tiene valor: se descarta (best effort) al salir
  // para no dejarla "en progreso" para siempre. Si el borrado falla, la salida
  // no se bloquea; el usuario puede eliminarla desde Sesiones.
  const handleLeaveEmpty = useCallback(
    async (target: '/(tabs)/sesiones' | '/(tabs)/sesiones/nueva') => {
      if (id) {
        try {
          await deleteSession(id)
        } catch {
          // Silencioso a propósito: es limpieza de una sesión sin respuestas.
        }
      }
      router.replace(target)
    },
    [id, router],
  )

  if (phase === 'loading' && !question) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    )
  }

  if (phase === 'error' || error) {
    return (
      <View style={styles.center}>
        <MaterialCommunityIcons name="alert-circle-outline" size={48} color={colors.destructive} />
        <Text style={styles.error}>{error ?? 'Ocurrio un error'}</Text>
        <Pressable
          style={styles.button}
          onPress={() => id && void loadNext(id)}
          accessibilityRole="button"
        >
          <MaterialCommunityIcons name="refresh" size={18} color={colors.onPrimary} />
          <Text style={styles.buttonLabel}>Reintentar</Text>
        </Pressable>
        <Pressable
          style={styles.buttonGhost}
          onPress={() => router.back()}
          accessibilityRole="button"
        >
          <Text style={styles.buttonGhostLabel}>Volver</Text>
        </Pressable>
      </View>
    )
  }

  if (phase === 'finished') {
    return (
      <View style={styles.center}>
        <MaterialCommunityIcons name="check-circle" size={64} color={colors.primary} />
        <Text style={styles.finishedTitle}>Sesion completada</Text>
        <Text style={styles.finishedMeta}>
          Respondiste {summary?.answerCount ?? 0} preguntas
        </Text>
        <Pressable
          style={styles.button}
          onPress={() => void handleFinish()}
          accessibilityRole="button"
          accessibilityLabel="Ver resumen de la sesion"
        >
          <MaterialCommunityIcons name="file-document-outline" size={18} color={colors.onPrimary} />
          <Text style={styles.buttonLabel}>Ver resumen</Text>
        </Pressable>
        <Pressable
          style={styles.buttonGhost}
          onPress={() => router.replace('/(tabs)/sesiones')}
          accessibilityRole="button"
        >
          <Text style={styles.buttonGhostLabel}>Volver a sesiones</Text>
        </Pressable>
      </View>
    )
  }

  if (phase === 'empty') {
    return (
      <View style={styles.center}>
        <MaterialCommunityIcons
          name="inbox-outline"
          size={56}
          color={colors.mutedForeground}
        />
        <Text style={styles.finishedTitle}>No hay preguntas para repasar</Text>
        <Text style={[styles.finishedMeta, styles.emptyMeta]}>
          No tienes repasos vencidos ni preguntas disponibles en este momento.
        </Text>
        <Pressable
          style={styles.button}
          onPress={() => void handleLeaveEmpty('/(tabs)/sesiones/nueva')}
          accessibilityRole="button"
          accessibilityLabel="Crear sesión de estudio"
        >
          <MaterialCommunityIcons name="plus-circle-outline" size={18} color={colors.onPrimary} />
          <Text style={styles.buttonLabel}>Crear sesión de estudio</Text>
        </Pressable>
        <Pressable
          style={styles.buttonGhost}
          onPress={() => void handleLeaveEmpty('/(tabs)/sesiones')}
          accessibilityRole="button"
        >
          <Text style={styles.buttonGhostLabel}>Volver a sesiones</Text>
        </Pressable>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <View style={[styles.header, { paddingTop: insets.top + spacing.sm }]}>
        <Pressable
          style={({ pressed }) => [styles.iconButton, pressed && styles.pressed]}
          onPress={() => router.back()}
          accessibilityRole="button"
          accessibilityLabel="Volver"
        >
          <MaterialCommunityIcons name="arrow-left" size={24} color={colors.primary} />
        </Pressable>
        <Text style={styles.counter}>
          {questionLimit !== null
            ? `${questionNumber} / ${questionLimit}`
            : `Pregunta ${questionNumber}`}
        </Text>
      </View>

      {questionLimit && (
        <View style={styles.progressContainer}>
          <View style={styles.progressTrack}>
            <View style={[styles.progressFill, { width: `${progress * 100}%` }]} />
          </View>
        </View>
      )}

      <ScrollView
        ref={scrollRef}
        style={styles.scroll}
        contentContainerStyle={styles.scrollContent}
      >
        {question && (
          <QuestionCard
            key={question.id}
            question={question}
            mode={phase === 'result' ? 'result' : 'answer'}
            selected={selectedOptions}
            lastResult={lastResult}
            onToggleOption={handleToggleOption}
            disabled={phase === 'loading'}
          />
        )}

        {phase === 'question' && submitError ? (
          <View style={styles.submitError}>
            <View style={styles.submitErrorHeader}>
              <MaterialCommunityIcons
                name={submitErrorIsNetwork ? 'cloud-off-outline' : 'alert-circle-outline'}
                size={18}
                color={colors.destructive}
              />
              <Text style={styles.submitErrorText}>{submitError}</Text>
            </View>
            <Pressable
              onPress={() => void handleAnswer(selectedOptions)}
              style={({ pressed }) => [styles.retryButton, pressed && styles.retryButtonPressed]}
              accessibilityRole="button"
              accessibilityLabel="Reintentar"
            >
              <MaterialCommunityIcons name="refresh" size={18} color={colors.onPrimary} />
              <Text style={styles.retryButtonLabel}>Reintentar</Text>
            </Pressable>
          </View>
        ) : null}
      </ScrollView>

      <View style={[styles.actionBar, { paddingBottom: Math.max(insets.bottom, spacing.sm) }]}>
        {barAction === 'next' ? (
          <Pressable
            style={({ pressed }) => [
              styles.actionButton,
              phase === 'loading' && styles.actionButtonDisabled,
              pressed && styles.actionButtonPressed,
            ]}
            onPress={handleNext}
            disabled={phase === 'loading'}
            accessibilityRole="button"
            accessibilityLabel="Siguiente pregunta"
          >
            <Text style={styles.actionLabel}>Siguiente</Text>
            <MaterialCommunityIcons name="arrow-right" size={18} color={colors.onPrimary} />
          </Pressable>
        ) : (
          <Pressable
            style={({ pressed }) => [
              styles.actionButton,
              (selectedOptions.length === 0 || phase === 'loading') && styles.actionButtonDisabled,
              pressed && styles.actionButtonPressed,
            ]}
            onPress={() => void handleAnswer(selectedOptions)}
            disabled={selectedOptions.length === 0 || phase === 'loading'}
            accessibilityRole="button"
            accessibilityLabel="Confirmar respuesta"
          >
            <MaterialCommunityIcons name="check" size={18} color={colors.onPrimary} />
            <Text style={styles.actionLabel}>Confirmar</Text>
          </Pressable>
        )}
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  center: {
    flex: 1,
    backgroundColor: colors.background,
    justifyContent: 'center',
    alignItems: 'center',
    padding: spacing.xl,
    gap: spacing.md,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.sm,
  },
  iconButton: {
    width: 44,
    height: 44,
    alignItems: 'center',
    justifyContent: 'center',
    marginLeft: -spacing.sm,
  },
  pressed: {
    opacity: 0.7,
  },
  counter: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
  },
  progressContainer: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.sm,
  },
  progressTrack: {
    height: 6,
    backgroundColor: colors.muted,
    borderRadius: 3,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: colors.primary,
    borderRadius: 3,
  },
  scroll: {
    flex: 1,
  },
  scrollContent: {
    padding: spacing.lg,
    // Deja aire para la barra fija de acción (no tapa la explicación).
    paddingBottom: spacing['3xl'],
    flexGrow: 1,
  },
  actionBar: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.sm,
    borderTopWidth: 1,
    borderTopColor: colors.border,
    backgroundColor: colors.background,
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    gap: spacing.sm,
    minHeight: 48,
  },
  actionButtonPressed: {
    opacity: 0.85,
  },
  actionButtonDisabled: {
    opacity: 0.5,
  },
  actionLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  submitError: {
    marginTop: spacing.lg,
    gap: spacing.md,
    padding: spacing.md,
    backgroundColor: `${colors.destructive}12`,
    borderWidth: 1,
    borderColor: `${colors.destructive}55`,
    borderRadius: radius.md,
  },
  submitErrorHeader: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: spacing.sm,
  },
  submitErrorText: {
    flex: 1,
    fontSize: typography.size.sm,
    lineHeight: 20,
    color: colors.foreground,
  },
  retryButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
    minHeight: 48,
    borderRadius: radius.md,
    backgroundColor: colors.primary,
  },
  retryButtonPressed: {
    opacity: 0.85,
  },
  retryButtonLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
  },
  button: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
    gap: spacing.sm,
    minHeight: 44,
  },
  buttonLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  buttonGhost: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
    minHeight: 44,
    justifyContent: 'center',
  },
  buttonGhostLabel: {
    fontSize: typography.size.md,
    color: colors.mutedForeground,
  },
  finishedTitle: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  finishedMeta: {
    fontSize: typography.size.md,
    color: colors.mutedForeground,
  },
  emptyMeta: {
    textAlign: 'center',
  },
})
