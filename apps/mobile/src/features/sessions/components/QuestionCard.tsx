import { useState } from 'react'
import { Pressable, StyleSheet, Text, View } from 'react-native'
import type { NextQuestionItem, SessionAnswer } from '@devtest/shared'

import { MarkdownText } from '@/core/components/MarkdownText'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

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
        {question.options.map((opt) => {
          const isSelected = activeSelection.has(opt.id)
          const isCorrect = opt.isCorrect === true
          const showCorrectHighlight = isResult && isCorrect
          const showWrongHighlight = isResult && isSelected && !isCorrect

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
              {isResult && (
                <Text style={[styles.optionIcon, isCorrect ? styles.iconCorrect : styles.iconWrong]}>
                  {isCorrect ? '✓' : isSelected ? '✗' : ''}
                </Text>
              )}
              <MarkdownText>
                {opt.content}
              </MarkdownText>
            </View>
          )
        })}
      </View>

      {isResult && lastResult?.explanation ? (
        <View style={styles.explanationBox}>
          <Text style={styles.explanationLabel}>Explicacion</Text>
          <MarkdownText style={styles.explanationText}>{lastResult.explanation}</MarkdownText>
        </View>
      ) : null}

      {isResult ? (
        <Pressable style={styles.submit} onPress={onNext}>
          <Text style={styles.submitLabel}>Siguiente</Text>
        </Pressable>
      ) : (
        <Pressable
          style={[styles.submit, activeSelection.size === 0 && styles.submitDisabled]}
          onPress={() => handleSubmit()}
          disabled={activeSelection.size === 0}
        >
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
    gap: spacing.sm,
  },
  optionSelected: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}15`,
  },
  optionCorrect: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}15`,
  },
  optionWrong: {
    borderColor: colors.destructive,
    backgroundColor: `${colors.destructive}15`,
  },
  optionIcon: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.bold,
    width: 20,
    textAlign: 'center',
  },
  iconCorrect: {
    color: colors.primary,
  },
  iconWrong: {
    color: colors.destructive,
  },
  optionText: {
    flex: 1,
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  optionTextSelected: {
    color: colors.primary,
    fontWeight: typography.weight.medium,
  },
  optionTextCorrect: {
    color: colors.primary,
    fontWeight: typography.weight.medium,
  },
  optionTextWrong: {
    color: colors.destructive,
  },
  explanationBox: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
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
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    alignItems: 'center',
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
