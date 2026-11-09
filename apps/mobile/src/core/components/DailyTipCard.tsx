import { MaterialCommunityIcons } from '@expo/vector-icons'
import { Pressable, StyleSheet, Text, View } from 'react-native'
import type { Tip } from '@devtest/shared'

import { MarkdownText } from '@/core/components/MarkdownText'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type DailyTipCardProps = {
  tip: Tip
  onDismiss: () => void
}

export function DailyTipCard({ tip, onDismiss }: DailyTipCardProps) {
  return (
    <View style={styles.card}>
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <MaterialCommunityIcons name="lightbulb-on-outline" size={18} color={colors.accent} />
          <Text style={styles.title}>Tip del día</Text>
        </View>

        <Pressable
          onPress={onDismiss}
          style={({ pressed }) => [styles.close, pressed && styles.closePressed]}
          accessibilityRole="button"
          accessibilityLabel="Cerrar tip del día"
        >
          <MaterialCommunityIcons name="close" size={18} color={colors.mutedForeground} />
        </Pressable>
      </View>

      <Text style={styles.topic} numberOfLines={1}>
        {tip.topic}
      </Text>

      <MarkdownText style={styles.content}>{tip.content}</MarkdownText>
    </View>
  )
}

const styles = StyleSheet.create({
  card: {
    gap: spacing.sm,
    padding: spacing.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  title: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.accent,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  close: {
    width: 44,
    height: 44,
    alignItems: 'center',
    justifyContent: 'center',
  },
  closePressed: {
    opacity: 0.7,
  },
  topic: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  content: {
    fontSize: typography.size.md,
    color: colors.foreground,
    lineHeight: 22,
  },
})
