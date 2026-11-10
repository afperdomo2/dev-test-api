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
  /** Opciones seleccionadas (controlado por la pantalla). */
  selected: string[]
  lastResult: SessionAnswer | null
  onToggleOption: (optionId: string) => void
  /** Bloquea la interacción (p. ej. mientras se envía la respuesta). */
  disabled?: boolean
}

/**
 * Tarjeta presentacional: pregunta, opciones y feedback. Las acciones
 * ("Confirmar" / "Siguiente") viven en la barra fija de `SessionStudyPage`, así
 * el botón principal no cambia de sitio entre fases.
 */
export function QuestionCard({
  question,
  mode,
  selected,
  lastResult,
  onToggleOption,
  disabled = false,
}: QuestionCardProps) {
  const isResult = mode === 'result'
  const locked = isResult || disabled
  const selectedSet = new Set(selected)

  // El contrato debería mandar siempre arrays, pero un backend viejo (o una
  // pregunta rara) no debe tumbar la pantalla con `.map is not a function`.
  const options = question.options ?? []
  const topics = question.topics ?? []

  // `lastResult.correctOptions` es la fuente de verdad del backend. El
  // fallback cubre backends viejos: si la respuesta fue correcta, las
  // opciones elegidas son las correctas; si fue incorrecta no sabemos cuál
  // era, así que solo marcamos las elegidas.
  const correctIds =
    isResult && lastResult?.correctOptions && lastResult.correctOptions.length > 0
      ? new Set(lastResult.correctOptions)
      : null
  const isCorrectOption = (optionId: string, isSelected: boolean) => {
    if (!isResult) return false
    if (correctIds) return correctIds.has(optionId)
    return isSelected && lastResult?.isCorrect === true
  }

  return (
    <View style={styles.container}>
      <MarkdownText style={styles.content}>{question.content}</MarkdownText>

      {topics.length > 0 && (
        <View style={styles.topics}>
          {topics.map((t) => (
            <Text key={t} style={styles.topic}>
              {t}
            </Text>
          ))}
        </View>
      )}

      <View style={styles.options}>
        {options.map((opt, idx) => {
          const isSelected = selectedSet.has(opt.id)
          const optionIsCorrect = isCorrectOption(opt.id, isSelected)
          const showWrongHighlight = isResult && isSelected && !optionIsCorrect
          const letter = LETTERS[idx] ?? String(idx + 1)
          // Sin banner visible: el estado de cada opción viaja en su etiqueta
          // accesible para que un lector de pantalla no dependa del color.
          const optionLabel = !isResult
            ? opt.content
            : optionIsCorrect
              ? `${opt.content}, correcta`
              : showWrongHighlight
                ? `${opt.content}, incorrecta`
                : opt.content

          return (
            <Pressable
              key={opt.id}
              onPress={() => onToggleOption(opt.id)}
              disabled={locked}
              style={({ pressed }) => [
                styles.option,
                isSelected && !isResult && styles.optionSelected,
                isResult && optionIsCorrect && styles.optionCorrect,
                showWrongHighlight && styles.optionWrong,
                pressed && !locked && styles.optionPressed,
              ]}
              accessibilityRole="button"
              accessibilityState={{ selected: isSelected, disabled: locked }}
              accessibilityLabel={optionLabel}
            >
              {isResult ? (
                <MaterialCommunityIcons
                  name={
                    optionIsCorrect
                      ? 'check-circle'
                      : showWrongHighlight
                        ? 'close-circle'
                        : 'minus'
                  }
                  size={18}
                  color={
                    optionIsCorrect
                      ? colors.primary
                      : showWrongHighlight
                        ? colors.destructive
                        : colors.mutedForeground
                  }
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
            </Pressable>
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
  optionPressed: {
    opacity: 0.85,
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
})
