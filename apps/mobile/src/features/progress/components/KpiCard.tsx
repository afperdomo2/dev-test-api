import type { ComponentProps } from 'react'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type IconName = ComponentProps<typeof MaterialCommunityIcons>['name']

type KpiCardProps = {
  icon: IconName
  label: string
  value: string
}

export function KpiCard({ icon, label, value }: KpiCardProps) {
  return (
    <View style={styles.card}>
      <View style={styles.header}>
        <MaterialCommunityIcons name={icon} size={16} color={colors.primary} />
        <Text style={styles.label}>{label}</Text>
      </View>
      <Text style={styles.value}>{value}</Text>
    </View>
  )
}

const styles = StyleSheet.create({
  card: {
    width: '48%',
    gap: spacing.xs,
    padding: spacing.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  label: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  value: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
})
