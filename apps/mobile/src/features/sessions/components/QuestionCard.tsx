import { useState } from 'react'
import { Pressable, StyleSheet, Text, View } from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import type { NextQuestionItem, SessionAnswer } from '@devtest/shared'

import { MarkdownText } from '@/core/components/MarkdownText'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const LETTERS = ['A', 'B', 'C', 'D', 'E', 'F']

type QuestionCardMode = 'answer' | 'result'

type QuestionCardProps = {
  question: NextQuestionItem
  mode: QuestionCardMode
  selectedOptions: string[]
  lastResult: SessionAnswer | null
  onAnswer: (selectedOptions: string[]) => void
  onNext: () => void
}

export function QuestionCard({
  question,
  mode,
  selectedOptions,
  lastResult,
  onAnswer,
  onNext,
}: QuestionCardProps) {
  const [localSelected, setLocalSelected] = useState<Set<string>>(new Set())
  const isMultiple = question.type === 'multiple_choice'
  const isResult = mode === 'result'

  const toggleOption = (id: string) => {
    if (isResult) return
    setLocalSelected((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        if (!isMultiple) next.clear()
        next.add(id)
      }
      return next
    })
  }

  const handleSubmit = () => {
    if (localSelected.size === 0) return
    onAnswer(Array.from(localSelected))
  }

  const activeSelection = isResult ? new Set(selectedOptions) : localSelected

  return (
    <View style={styles.container}>
      <MarkdownText style={styles.content}>{question.content}</MarkdownText>

      {question.topics.length > 0 && (
        <View style={styles.topics}>
          {question.topics.map((t) => (
            <Text key={t} style={styles.topic}>
              {t}
            </Text>
          ))}
        </View>
      )}

      <View style={styles.options}>
        {question.options.map((opt, idx) => {
          const isSelected = activeSelection.has(opt.id)
          const isCorrect = opt.isCorrect === true
          const showCorrectHighlight = isResult && isCorrect
          const showWrongHighlight = isResult && isSelected && !isCorrect
          const letter = LETTERS[idx] ?? String(idx + 1)

          return (
            <View
              key={opt.id}
              style={[
                styles.option,
                isSelected && !isResult && styles.optionSelected,
                showCorrectHighlight && styles.optionCorrect,
                showWrongHighlight && styles.optionWrong,
              ]}
            >
              {isResult ? (
                <MaterialCommunityIcons
                  name={isCorrect ? 'check-circle' : isSelected ? 'close-circle' : 'circle-outline'}
                  size={18}
                  color={isCorrect ? colors.primary : isSelected ? colors.destructive : colors.mutedForeground}
                />
              ) : (
                <View style={[styles.letterBadge, isSelected && styles.letterBadgeSelected]}>
                  <Text style={[styles.letterText, isSelected && styles.letterTextSelected]}>
                    {letter}
                  </Text>
                </View>
              )}
              <MarkdownText style={styles.optionText}>
                {opt.content}
              </MarkdownText>
            </View>
          )
        })}
      </View>

      {isResult && lastResult?.explanation ? (
        <View style={styles.explanationBox}>
          <View style={styles.explanationHeader}>
            <MaterialCommunityIcons name="lightbulb-outline" size={16} color={colors.accent} />
            <Text style={styles.explanationLabel}>Explicacion</Text>
          </View>
          <MarkdownText style={styles.explanationText}>{lastResult.explanation}</MarkdownText>
        </View>
      ) : null}

      {isResult ? (
        <Pressable
          style={({ pressed }) => [styles.submit, pressed && styles.submitPressed]}
          onPress={onNext}
          accessibilityRole="button"
          accessibilityLabel="Siguiente pregunta"
        >
          <Text style={styles.submitLabel}>Siguiente</Text>
          <MaterialCommunityIcons name="arrow-right" size={18} color={colors.onPrimary} />
        </Pressable>
      ) : (
        <Pressable
          style={({ pressed }) => [
            styles.submit,
            activeSelection.size === 0 && styles.submitDisabled,
            pressed && styles.submitPressed,
          ]}
          onPress={() => handleSubmit()}
          disabled={activeSelection.size === 0}
          accessibilityRole="button"
          accessibilityLabel="Confirmar respuesta"
        >
          <MaterialCommunityIcons name="check" size={18} color={colors.onPrimary} />
          <Text style={styles.submitLabel}>Confirmar</Text>
        </Pressable>
      )}
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    gap: spacing.lg,
  },
  content: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.medium,
    color: colors.foreground,
    lineHeight: 28,
  },
  topics: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.xs,
  },
  topic: {
    fontSize: typography.size.xs,
    color: colors.primary,
    backgroundColor: `${colors.primary}15`,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    overflow: 'hidden',
  },
  options: {
    gap: spacing.sm,
  },
  option: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
    padding: spacing.md,
    backgroundColor: colors.card,
    gap: spacing.md,
    minHeight: 48,
  },
  optionSelected: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}15`,
  },
  optionCorrect: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}10`,
  },
  optionWrong: {
    borderColor: colors.destructive,
    backgroundColor: `${colors.destructive}10`,
  },
  letterBadge: {
    width: 28,
    height: 28,
    borderRadius: 14,
    backgroundColor: colors.muted,
    justifyContent: 'center',
    alignItems: 'center',
  },
  letterBadgeSelected: {
    backgroundColor: colors.primary,
  },
  letterText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.bold,
    color: colors.mutedForeground,
  },
  letterTextSelected: {
    color: colors.onPrimary,
  },
  optionText: {
    flex: 1,
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  explanationBox: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.sm,
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
  submit: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    gap: spacing.sm,
    minHeight: 48,
  },
  submitPressed: {
    opacity: 0.85,
  },
  submitDisabled: {
    opacity: 0.5,
  },
  submitLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
})
