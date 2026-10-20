import { MaterialCommunityIcons } from '@expo/vector-icons'
import { Pressable, StyleSheet, Text, View } from 'react-native'
import type { ComponentProps } from 'react'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type IconName = ComponentProps<typeof MaterialCommunityIcons>['name']

type ProgressHubRowProps = {
  icon: IconName
  title: string
  subtitle: string
  onPress: () => void
}

export function ProgressHubRow({ icon, title, subtitle, onPress }: ProgressHubRowProps) {
  return (
    <Pressable
      onPress={onPress}
      style={({ pressed }) => [styles.row, pressed && styles.rowPressed]}
      accessibilityRole="button"
      accessibilityLabel={title}
    >
      <View style={styles.iconWrap}>
        <MaterialCommunityIcons name={icon} size={22} color={colors.primary} />
      </View>

      <View style={styles.texts}>
        <Text style={styles.title}>{title}</Text>
        <Text style={styles.subtitle} numberOfLines={2}>
          {subtitle}
        </Text>
      </View>

      <MaterialCommunityIcons name="chevron-right" size={22} color={colors.mutedForeground} />
    </Pressable>
  )
}

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    minHeight: 64,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
  },
  rowPressed: {
    opacity: 0.85,
  },
  iconWrap: {
    width: 40,
    height: 40,
    borderRadius: radius.sm,
    backgroundColor: `${colors.primary}15`,
    alignItems: 'center',
    justifyContent: 'center',
  },
  texts: {
    flex: 1,
    gap: 2,
  },
  title: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  subtitle: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
})
