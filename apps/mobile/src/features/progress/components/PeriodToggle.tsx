import { Pressable, StyleSheet, Text, View } from 'react-native'
import type { InsightPeriod } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const OPTIONS: Array<{ title: string; value: InsightPeriod }> = [
  { title: 'Semana', value: 'week' },
  { title: 'Mes', value: 'month' },
]

type PeriodToggleProps = {
  value: InsightPeriod
  onChange: (period: InsightPeriod) => void
}

export function PeriodToggle({ value, onChange }: PeriodToggleProps) {
  return (
    <View style={styles.container}>
      {OPTIONS.map((option) => {
        const selected = option.value === value

        return (
          <Pressable
            key={option.value}
            onPress={() => onChange(option.value)}
            style={[styles.option, selected && styles.optionSelected]}
            accessibilityRole="button"
            accessibilityState={{ selected }}
            accessibilityLabel={option.title}
          >
            <Text style={[styles.label, selected && styles.labelSelected]}>{option.title}</Text>
          </Pressable>
        )
      })}
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  option: {
    minHeight: 44,
    paddingHorizontal: spacing.lg,
    justifyContent: 'center',
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
  },
  optionSelected: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}22`,
  },
  label: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
  },
  labelSelected: {
    color: colors.primary,
  },
})
