import { ActivityIndicator, Pressable, StyleSheet, Text } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type AuthButtonProps = {
  title: string
  onPress: () => void
  loading?: boolean
  disabled?: boolean
}

export function AuthButton({ title, onPress, loading, disabled }: AuthButtonProps) {
  const inactive = loading === true || disabled === true
  return (
    <Pressable
      style={[styles.button, inactive ? styles.buttonInactive : null]}
      onPress={onPress}
      disabled={inactive}
      accessibilityRole="button"
    >
      {loading === true ? (
        <ActivityIndicator color={colors.onPrimary} />
      ) : (
        <Text style={styles.label}>{title}</Text>
      )}
    </Pressable>
  )
}

const styles = StyleSheet.create({
  button: {
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 52,
  },
  buttonInactive: {
    opacity: 0.6,
  },
  label: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
})
