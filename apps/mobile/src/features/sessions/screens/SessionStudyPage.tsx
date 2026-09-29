import { useCallback, useEffect, useRef, useState } from 'react'
import { ActivityIndicator, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native'
import { useLocalSearchParams, useRouter } from 'expo-router'
import type { NextQuestionItem, SessionAnswer, Reward } from '@devtest/shared'

import { QuestionCard } from '@/features/sessions/components/QuestionCard'
import { FeedbackOverlay } from '@/features/sessions/components/FeedbackOverlay'
import {
  getSessionById,
  getNextQuestion,
  submitAnswer,
  finishSession,
} from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type StudyState = 'loading' | 'question' | 'feedback' | 'finished' | 'error'

export default function SessionStudyPage() {
  const { id } = useLocalSearchParams<{ id: string }>()
  const router = useRouter()

  const [state, setState] = useState<StudyState>('loading')
  const [question, setQuestion] = useState<NextQuestionItem | null>(null)
  const [feedback, setFeedback] = useState<{
    isCorrect: boolean
    xpGained: number
    reward?: Reward
  } | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [answerCount, setAnswerCount] = useState(0)
  const loadingRef = useRef(false)

  const loadNext = useCallback(
    async (sessionId: string) => {
      if (loadingRef.current) return
      loadingRef.current = true
      setState('loading')
      try {
        const next = await getNextQuestion(sessionId)
        setQuestion(next)
        setState('question')
      } catch {
        setState('finished')
      } finally {
        loadingRef.current = false
      }
    },
    [],
  )

  useEffect(() => {
    if (!id) return
    void loadNext(id)
  }, [id, loadNext])

  const handleAnswer = useCallback(
    async (selectedOptions: string[]) => {
      if (!id || !question) return
      setState('feedback')
      try {
        const result = await submitAnswer(id, {
          questionId: question.id,
          selectedOptions,
        })
        setFeedback({
          isCorrect: result.isCorrect,
          xpGained: 0, // el backend no retorna xpGained en SessionAnswer
        })
        setAnswerCount((prev) => prev + 1)
      } catch {
        // si falla, avanzar de todas formas
        setFeedback({ isCorrect: false, xpGained: 0 })
      }
    },
    [id, question],
  )

  const handleFeedbackDismiss = useCallback(() => {
    setFeedback(null)
    if (id) void loadNext(id)
  }, [id, loadNext])

  const handleFinish = useCallback(async () => {
    if (!id) return
    try {
      await finishSession(id)
    } finally {
      router.replace('/(tabs)/sesiones')
    }
  }, [id, router])

  if (state === 'loading') {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    )
  }

  if (state === 'error' || error) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>{error ?? 'Ocurrio un error'}</Text>
        <Pressable style={styles.button} onPress={() => id && void loadNext(id)}>
          <Text style={styles.buttonLabel}>Reintentar</Text>
        </Pressable>
        <Pressable style={styles.buttonSecondary} onPress={() => router.back()}>
          <Text style={styles.buttonSecondaryLabel}>Volver</Text>
        </Pressable>
      </View>
    )
  }

  if (state === 'finished') {
    return (
      <View style={styles.center}>
        <Text style={styles.finishedIcon}>🎉</Text>
        <Text style={styles.finishedTitle}>Sesion completada</Text>
        <Text style={styles.finishedMeta}>Respondiste {answerCount} preguntas</Text>
        <Pressable style={styles.button} onPress={() => handleFinish()}>
          <Text style={styles.buttonLabel}>Ver resumen</Text>
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
        <Text style={styles.counter}>{answerCount}</Text>
      </View>

      <ScrollView style={styles.scroll} contentContainerStyle={styles.scrollContent}>
        {question && <QuestionCard question={question} onAnswer={handleAnswer} />}
      </ScrollView>

      {feedback && (
        <FeedbackOverlay
          isCorrect={feedback.isCorrect}
          xpGained={feedback.xpGained}
          reward={feedback.reward}
          onDismiss={handleFeedbackDismiss}
        />
      )}
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
    fontSize: typography.size.md,
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
  buttonSecondary: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
  },
  buttonSecondaryLabel: {
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
