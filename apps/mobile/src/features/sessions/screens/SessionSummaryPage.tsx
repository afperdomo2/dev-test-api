import { useCallback, useEffect, useState } from 'react'
import { ActivityIndicator, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
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
        <MaterialCommunityIcons name="alert-circle-outline" size={48} color={colors.destructive} />
        <Text style={styles.error}>{error ?? 'No se pudo cargar el resumen'}</Text>
        <Pressable style={styles.stickyButton} onPress={handleBack} accessibilityRole="button">
          <Text style={styles.stickyButtonLabel}>Volver</Text>
        </Pressable>
      </View>
    )
  }

  const { session, answers } = detail
  const correctCount = answers.filter((a) => a.isCorrect).length
  const total = answers.length
  const percentage = total > 0 ? Math.round((correctCount / total) * 100) : 0

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Pressable
          onPress={handleBack}
          accessibilityRole="button"
          accessibilityLabel="Volver a sesiones"
        >
          <MaterialCommunityIcons name="arrow-left" size={24} color={colors.primary} />
        </Pressable>
        <Text style={styles.headerTitle} numberOfLines={1}>
          {session.name}
        </Text>
      </View>

      <ScrollView contentContainerStyle={styles.content}>
        <View style={styles.scoreBanner}>
          <MaterialCommunityIcons name="check-circle" size={48} color={colors.primary} />
          <Text style={styles.scoreTitle}>Sesion completada</Text>

          <Text style={styles.scorePercentage}>{percentage}%</Text>
          <Text style={styles.scoreBreakdown}>
            {correctCount} de {total} correctas
          </Text>

          <View style={styles.progressTrack}>
            <View style={[styles.progressFill, { width: `${percentage}%` }]} />
          </View>

          {session.score != null && session.score > 0 && (
            <View style={styles.xpRow}>
              <MaterialCommunityIcons name="lightning-bolt" size={16} color={colors.accent} />
              <Text style={styles.scoreXp}>{session.score} puntos</Text>
            </View>
          )}
        </View>

        <Text style={styles.sectionTitle}>Respuestas</Text>

        {answers.map((answer, idx) => (
          <AnswerCard key={answer.id} answer={answer} index={idx} />
        ))}
      </ScrollView>

      <View style={styles.stickyFooter}>
        <Pressable
          style={styles.stickyButton}
          onPress={handleBack}
          accessibilityRole="button"
          accessibilityLabel="Volver a sesiones"
        >
          <MaterialCommunityIcons name="format-list-bulleted" size={18} color={colors.onPrimary} />
          <Text style={styles.stickyButtonLabel}>Volver a sesiones</Text>
        </Pressable>
      </View>
    </View>
  )
}

function AnswerCard({ answer, index }: { answer: SessionAnswerDetail; index: number }) {
  const q = answer.question
  const selectedSet = new Set(answer.selectedOptions)

  return (
    <View style={answerStyles.card}>
      <View style={answerStyles.header}>
        <View style={answerStyles.headerLeft}>
          <MaterialCommunityIcons
            name={answer.isCorrect ? 'check-circle' : 'close-circle'}
            size={20}
            color={answer.isCorrect ? colors.primary : colors.destructive}
          />
          <Text style={answerStyles.number}>#{index + 1}</Text>
        </View>
        <View
          style={[
            answerStyles.badge,
            answer.isCorrect ? answerStyles.badgeOk : answerStyles.badgeFail,
          ]}
        >
          <Text
            style={[
              answerStyles.badgeText,
              answer.isCorrect ? answerStyles.badgeTextOk : answerStyles.badgeTextFail,
            ]}
          >
            {answer.isCorrect ? 'Correcta' : 'Incorrecta'}
          </Text>
        </View>
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
            <MaterialCommunityIcons
              name={isCorrect ? 'check' : wasSelected ? 'close' : 'minus'}
              size={14}
              color={isCorrect ? colors.primary : wasSelected ? colors.destructive : colors.mutedForeground}
            />
            <MarkdownText style={answerStyles.optionText}>
              {opt.content}
            </MarkdownText>
          </View>
        )
      })}

      {answer.explanation ? (
        <View style={answerStyles.explanation}>
          <View style={answerStyles.explanationHeader}>
            <MaterialCommunityIcons name="lightbulb-outline" size={14} color={colors.accent} />
            <Text style={answerStyles.explanationLabel}>Explicacion</Text>
          </View>
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
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
    paddingBottom: spacing.sm,
    gap: spacing.md,
  },
  headerTitle: {
    flex: 1,
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  content: {
    padding: spacing.lg,
    paddingBottom: 100,
    gap: spacing.md,
  },
  scoreBanner: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.lg,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.sm,
  },
  scoreTitle: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  scorePercentage: {
    fontSize: 48,
    fontWeight: typography.weight.bold,
    color: colors.primary,
    marginTop: spacing.sm,
  },
  scoreBreakdown: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  progressTrack: {
    width: '100%',
    height: 8,
    backgroundColor: colors.muted,
    borderRadius: 4,
    overflow: 'hidden',
    marginTop: spacing.sm,
  },
  progressFill: {
    height: '100%',
    backgroundColor: colors.primary,
    borderRadius: 4,
  },
  xpRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    marginTop: spacing.xs,
  },
  scoreXp: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.accent,
  },
  sectionTitle: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
    paddingHorizontal: spacing.xs,
  },
  error: {
    fontSize: typography.size.md,
    color: colors.destructive,
    textAlign: 'center',
  },
  stickyFooter: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    padding: spacing.lg,
    paddingBottom: spacing.xl,
    backgroundColor: colors.background,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  stickyButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    gap: spacing.sm,
    minHeight: 48,
  },
  stickyButtonLabel: {
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
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  number: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
  },
  badge: {
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
  },
  badgeOk: {
    backgroundColor: `${colors.primary}15`,
  },
  badgeFail: {
    backgroundColor: `${colors.destructive}15`,
  },
  badgeText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
  },
  badgeTextOk: {
    color: colors.primary,
  },
  badgeTextFail: {
    color: colors.destructive,
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
    gap: spacing.sm,
  },
  optionCorrect: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}10`,
  },
  optionWrong: {
    borderColor: colors.destructive,
    backgroundColor: `${colors.destructive}10`,
  },
  optionText: {
    flex: 1,
    fontSize: typography.size.sm,
    color: colors.foreground,
  },
  explanation: {
    backgroundColor: colors.background,
    borderRadius: radius.sm,
    padding: spacing.sm,
    gap: spacing.xs,
  },
  explanationHeader: {
    flexDirection: 'row',
    alignItems: 'center',
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
