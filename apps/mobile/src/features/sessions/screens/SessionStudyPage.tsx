import { useCallback, useEffect, useRef, useState } from 'react'
import {
  ActivityIndicator,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native'
import { useLocalSearchParams, useRouter } from 'expo-router'
import type { NextQuestionItem, SessionAnswer, SessionSummary, Reward } from '@devtest/shared'

import { QuestionCard } from '@/features/sessions/components/QuestionCard'
import {
  getSessionById,
  getSessionSummary,
  getNextQuestion,
  submitAnswer,
  finishSession,
} from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
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
  const loadingRef = useRef(false)

  const questionNumber = summary ? summary.answerCount + 1 : 1
  const questionLimit = summary?.questionLimit ?? null

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
      try {
        const next = await getNextQuestion(sessionId)
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
      if (!id || !question) return
      setSelectedOptions(options)
      setPhase('loading')
      try {
        const result = await submitAnswer(id, {
          questionId: question.id,
          selectedOptions: options,
        })
        setLastResult(result)
        setSummary((prev) =>
          prev ? { ...prev, answerCount: prev.answerCount + 1 } : prev,
        )
        setPhase('result')
      } catch {
        setPhase('result')
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
        <Text style={styles.error}>{error ?? 'Ocurrio un error'}</Text>
        <Pressable style={styles.button} onPress={() => id && void loadNext(id)}>
          <Text style={styles.buttonLabel}>Reintentar</Text>
        </Pressable>
        <Pressable style={styles.buttonGhost} onPress={() => router.back()}>
          <Text style={styles.buttonGhostLabel}>Volver</Text>
        </Pressable>
      </View>
    )
  }

  if (phase === 'finished') {
    return (
      <View style={styles.center}>
        <Text style={styles.finishedIcon}>🎉</Text>
        <Text style={styles.finishedTitle}>Sesion completada</Text>
        <Text style={styles.finishedMeta}>
          Respondiste {summary?.answerCount ?? 0} preguntas
        </Text>
        <Pressable style={styles.button} onPress={() => void handleFinish()}>
          <Text style={styles.buttonLabel}>Ver resumen</Text>
        </Pressable>
        <Pressable style={styles.buttonGhost} onPress={() => router.replace('/(tabs)/sesiones')}>
          <Text style={styles.buttonGhostLabel}>Volver a sesiones</Text>
        </Pressable>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Pressable onPress={() => router.back()}>
          <Text style={styles.backLabel}>← Volver</Text>
        </Pressable>
        <Text style={styles.counter}>
          {questionLimit !== null
            ? `${questionNumber} / ${questionLimit}`
            : `Pregunta ${questionNumber}`}
        </Text>
      </View>

      <ScrollView style={styles.scroll} contentContainerStyle={styles.scrollContent}>
        {question && (
          <QuestionCard
            question={question}
            mode={phase === 'result' ? 'result' : 'answer'}
            selectedOptions={selectedOptions}
            lastResult={lastResult}
            onAnswer={handleAnswer}
            onNext={handleNext}
          />
        )}
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
  backLabel: {
    fontSize: typography.size.md,
    color: colors.primary,
  },
  counter: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
  },
  scroll: {
    flex: 1,
  },
  scrollContent: {
    padding: spacing.lg,
    flexGrow: 1,
  },
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
  },
  button: {
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
  },
  buttonLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
  buttonGhost: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
  },
  buttonGhostLabel: {
    fontSize: typography.size.md,
    color: colors.mutedForeground,
  },
  finishedIcon: {
    fontSize: 48,
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
