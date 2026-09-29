import { useCallback, useEffect, useState } from 'react'
import { ActivityIndicator, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native'
import { useLocalSearchParams, useRouter } from 'expo-router'
import type { SessionDetail, SessionAnswerDetail } from '@devtest/shared'

import { MarkdownText } from '@/core/components/MarkdownText'
import { getSessionDetail } from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

export default function SessionSummaryPage() {
  const { id } = useLocalSearchParams<{ id: string }>()
  const router = useRouter()
  const [detail, setDetail] = useState<SessionDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    let cancelled = false
    getSessionDetail(id)
      .then((res) => {
        if (!cancelled) setDetail(res)
      })
      .catch((err) => {
        if (!cancelled) setError(getErrorDetail(err, 'No se pudo cargar el resumen'))
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [id])

  const handleBack = useCallback(() => {
    router.replace('/(tabs)/sesiones')
  }, [router])

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    )
  }

  if (error || !detail) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>{error ?? 'No se pudo cargar el resumen'}</Text>
        <Pressable style={styles.button} onPress={handleBack}>
          <Text style={styles.buttonLabel}>Volver</Text>
        </Pressable>
      </View>
    )
  }

  const { session, answers } = detail
  const correctCount = answers.filter((a) => a.isCorrect).length
  const total = answers.length

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Pressable onPress={handleBack}>
          <Text style={styles.backLabel}>← Sesiones</Text>
        </Pressable>
        <Text style={styles.headerTitle}>{session.name}</Text>
      </View>

      <ScrollView contentContainerStyle={styles.content}>
        <View style={styles.scoreBanner}>
          <Text style={styles.scoreValue}>
            {correctCount}/{total}
          </Text>
          <Text style={styles.scoreLabel}>respuestas correctas</Text>
          {session.score != null && session.score > 0 && (
            <Text style={styles.scoreXp}>{session.score} puntos</Text>
          )}
        </View>

        {answers.map((answer, idx) => (
          <AnswerCard key={answer.id} answer={answer} index={idx} />
        ))}

        <Pressable style={styles.button} onPress={handleBack}>
          <Text style={styles.buttonLabel}>Volver a sesiones</Text>
        </Pressable>
      </ScrollView>
    </View>
  )
}

function AnswerCard({ answer, index }: { answer: SessionAnswerDetail; index: number }) {
  const q = answer.question
  const selectedSet = new Set(answer.selectedOptions)

  return (
    <View style={answerStyles.card}>
      <View style={answerStyles.header}>
        <Text style={answerStyles.number}>#{index + 1}</Text>
        <Text style={[answerStyles.badge, answer.isCorrect ? answerStyles.badgeOk : answerStyles.badgeFail]}>
          {answer.isCorrect ? 'Correcto' : 'Incorrecto'}
        </Text>
      </View>

      <MarkdownText style={answerStyles.question}>{q.content}</MarkdownText>

      {q.options.map((opt) => {
        const wasSelected = selectedSet.has(opt.id)
        const isCorrect = opt.isCorrect
        return (
          <View
            key={opt.id}
            style={[
              answerStyles.option,
              isCorrect && answerStyles.optionCorrect,
              wasSelected && !isCorrect && answerStyles.optionWrong,
            ]}
          >
            <Text style={answerStyles.optionIcon}>
              {isCorrect ? '✓' : wasSelected ? '✗' : ''}
            </Text>
            <MarkdownText>
              {opt.content}
            </MarkdownText>
          </View>
        )
      })}

      {answer.explanation ? (
        <View style={answerStyles.explanation}>
          <Text style={answerStyles.explanationLabel}>Explicacion</Text>
          <MarkdownText style={answerStyles.explanationText}>{answer.explanation}</MarkdownText>
        </View>
      ) : null}
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
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
    paddingBottom: spacing.sm,
    gap: spacing.xs,
  },
  backLabel: {
    fontSize: typography.size.md,
    color: colors.primary,
  },
  headerTitle: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  content: {
    padding: spacing.lg,
    gap: spacing.md,
    paddingBottom: spacing.xl,
  },
  scoreBanner: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.lg,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.xs,
  },
  scoreValue: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    color: colors.primary,
  },
  scoreLabel: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  scoreXp: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.accent,
  },
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
  },
  button: {
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    alignItems: 'center',
  },
  buttonLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
})

const answerStyles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.sm,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  number: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
  },
  badge: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    overflow: 'hidden',
  },
  badgeOk: {
    color: colors.primary,
    backgroundColor: `${colors.primary}15`,
  },
  badgeFail: {
    color: colors.destructive,
    backgroundColor: `${colors.destructive}15`,
  },
  question: {
    fontSize: typography.size.md,
    color: colors.foreground,
    lineHeight: 22,
  },
  option: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.sm,
    padding: spacing.sm,
    gap: spacing.xs,
  },
  optionCorrect: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}10`,
  },
  optionWrong: {
    borderColor: colors.destructive,
    backgroundColor: `${colors.destructive}10`,
  },
  optionIcon: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.bold,
    width: 16,
    textAlign: 'center',
  },
  optionText: {
    flex: 1,
    fontSize: typography.size.sm,
    color: colors.foreground,
  },
  optionTextCorrect: {
    color: colors.primary,
    fontWeight: typography.weight.medium,
  },
  optionTextWrong: {
    color: colors.destructive,
  },
  explanation: {
    backgroundColor: colors.background,
    borderRadius: radius.sm,
    padding: spacing.sm,
    gap: spacing.xs,
  },
  explanationLabel: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  explanationText: {
    fontSize: typography.size.sm,
    color: colors.foreground,
    lineHeight: 20,
  },
})
