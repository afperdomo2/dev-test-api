import { ScrollView, StyleSheet, Text, View } from 'react-native'
import type { ReactNode } from 'react'

import { colors, spacing, typography } from '@/core/theme/tokens'

type ScreenProps = {
  title: string
  children?: ReactNode
}

export function Screen({ title, children }: ScreenProps) {
  return (
    <ScrollView
      style={styles.container}
      contentInsetAdjustmentBehavior="automatic"
      contentContainerStyle={styles.content}
    >
      <Text style={styles.title}>{title}</Text>
      {children ? <View style={styles.body}>{children}</View> : null}
    </ScrollView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing.lg,
    gap: spacing.lg,
  },
  title: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
  },
  body: {
    gap: spacing.lg,
  },
})