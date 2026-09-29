import { StyleSheet, Text } from 'react-native'

import { colors, typography } from '@/core/theme/tokens'

type AuthServerErrorProps = {
  message: string | null
}

export function AuthServerError({ message }: AuthServerErrorProps) {
  if (!message) return null
  return <Text style={styles.error}>{message}</Text>
}

const styles = StyleSheet.create({
  error: {
    fontSize: typography.size.sm,
    color: colors.destructive,
    textAlign: 'center',
  },
})
