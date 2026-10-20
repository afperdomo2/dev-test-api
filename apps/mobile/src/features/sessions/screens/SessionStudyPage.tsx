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
import type { NextQuestionItem, SessionAnswer, SessionSummary, Reward } from '@devtest/shared'

import { QuestionCard } from '@/features/sessions/components/QuestionCard'
import { DailyTipCard } from '@/features/sessions/components/DailyTipCard'
import { clearPendingTip, getPendingTip } from '@/features/sessions/dailyTipStore'
import {
  getSessionById,
  getSessionSummary,
  getNextQuestion,
  submitAnswer,
  finishSession,
} from '@/core/api/services/sessions.service'
import { getErrorDetail, isNetworkError } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type Phase = 'loading' | 'question' | 'result' | 'finished' | 'error'

export default function SessionStudyPage() {
  const { id } = useLocalSearchParams<{ id: string }>()
  const router = useRouter()

  const [phase, setPhase] = useState<Phase>('loading')
  const [sessionStatus, setSessionStatus] = useState<string | null>(null)
  const [question, setQuestion] = useState<NextQuestionItem | null>(null)
  const [summary, setSummary] = useState<SessionSummary | null>(null)
  const [selectedOptions, setSelectedOptions] = useState<string[]>([])
  const [lastResult, setLastResult] = useState<SessionAnswer | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitErrorIsNetwork, setSubmitErrorIsNetwork] = useState(false)
  const loadingRef = useRef(false)
  const submittingRef = useRef(false)
  // UUID estable por pregunta: se reutiliza en los reintentos y se renueva al
  // cargar la siguiente (idempotencia, D-35).
  const clientAnswerIdRef = useRef<string | null>(null)

  // El tip del día se guarda al crear la sesión (Home) y solo se muestra al
  // empezar. Se limpia al salir de la pantalla.
  const [tipDismissed, setTipDismissed] = useState(false)
  const pendingTip = id ? getPendingTip(id) : null

  useEffect(() => () => clearPendingTip(), [])

  const questionNumber = summary ? summary.answerCount + 1 : 1
  const questionLimit = summary?.questionLimit ?? null
  const progress = questionLimit ? questionNumber / questionLimit : 0

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
      setSelectedOptions([])
      setLastResult(null)
      setSubmitError(null)
      setSubmitErrorIsNetwork(false)
      clientAnswerIdRef.current = null
      try {
        const next = await getNextQuestion(sessionId)
        clientAnswerIdRef.current = randomUUID()
        setQuestion(next)
        setPhase('question')
      } catch {
        await loadSessionData(sessionId)
        setPhase('finished')
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
        const { session, summary: sum } = await loadSessionData(id)
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
        setPhase('result')
      } catch (err) {
        // No se avanza: el usuario ve el error y puede reintentar sin duplicar.
        const network = isNetworkError(err)
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

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Pressable
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

      <ScrollView style={styles.scroll} contentContainerStyle={styles.scrollContent}>
        {pendingTip && !tipDismissed && phase === 'question' && questionNumber === 1 && (
          <View style={styles.tipWrapper}>
            <DailyTipCard tip={pendingTip} onDismiss={() => setTipDismissed(true)} />
          </View>
        )}

        {question && (
          <QuestionCard
            key={question.id}
            question={question}
            mode={phase === 'result' ? 'result' : 'answer'}
            selectedOptions={selectedOptions}
            lastResult={lastResult}
            onAnswer={handleAnswer}
            onNext={handleNext}
            submitting={phase === 'loading'}
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
    paddingTop: spacing.xl,
    paddingBottom: spacing.sm,
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
    flexGrow: 1,
  },
  tipWrapper: {
    marginBottom: spacing.lg,
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
})
