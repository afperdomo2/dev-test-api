import { StyleSheet, Text, View } from 'react-native'
import type { ReactNode } from 'react'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type AuthCardProps = {
  title: string
  subtitle?: string
  children: ReactNode
}

export function AuthCard({ title, subtitle, children }: AuthCardProps) {
  return (
    <View style={styles.card}>
      <Text style={styles.title}>{title}</Text>
      {subtitle ? <Text style={styles.subtitle}>{subtitle}</Text> : null}
      <View style={styles.body}>{children}</View>
    </View>
  )
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.lg,
    gap: spacing.sm,
  },
  title: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
    textAlign: 'center',
  },
  subtitle: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
    textAlign: 'center',
  },
  body: {
    marginTop: spacing.sm,
    gap: spacing.md,
  },
})
