import { useState } from 'react'
import { Pressable, StyleSheet, Text, View } from 'react-native'
import type { NextQuestionItem } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type QuestionCardProps = {
  question: NextQuestionItem
  onAnswer: (selectedOptions: string[]) => void
}

export function QuestionCard({ question, onAnswer }: QuestionCardProps) {
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const isMultiple = question.type === 'multiple_choice'

  const toggleOption = (id: string) => {
    setSelected((prev) => {
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
    if (selected.size === 0) return
    onAnswer(Array.from(selected))
  }

  return (
    <View style={styles.container}>
      <Text style={styles.content}>{question.content}</Text>

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
          const isSelected = selected.has(opt.id)
          return (
            <Pressable
              key={opt.id}
              style={[styles.option, isSelected && styles.optionSelected]}
              onPress={() => toggleOption(opt.id)}
            >
              <Text style={[styles.optionText, isSelected && styles.optionTextSelected]}>
                {opt.content}
              </Text>
            </Pressable>
          )
        })}
      </View>

      <Pressable
        style={[styles.submit, selected.size === 0 && styles.submitDisabled]}
        onPress={handleSubmit}
        disabled={selected.size === 0}
      >
        <Text style={styles.submitLabel}>Confirmar</Text>
      </Pressable>
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
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
    padding: spacing.md,
    backgroundColor: colors.card,
  },
  optionSelected: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}15`,
  },
  optionText: {
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  optionTextSelected: {
    color: colors.primary,
    fontWeight: typography.weight.medium,
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
