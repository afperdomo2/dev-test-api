import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'
import type { ComponentProps } from 'react'
import type { Badge } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type IconName = ComponentProps<typeof MaterialCommunityIcons>['name']

const FALLBACK_ICON: IconName = 'trophy-outline'

/**
 * El backend manda nombres `mdi-*`; MaterialCommunityIcons los usa sin prefijo.
 * Si el nombre no existe en la fuente, se cae al trofeo en vez de romper la pantalla.
 */
export function toMcIcon(name: string): IconName {
  const stripped = name.replace(/^mdi-/, '')
  const glyphMap = MaterialCommunityIcons.glyphMap as Record<string, number> | undefined

  if (glyphMap && !(stripped in glyphMap)) return FALLBACK_ICON

  return stripped as IconName
}

function percent(badge: Badge): number {
  if (badge.target <= 0) return 0
  return Math.min(100, Math.round((badge.current / badge.target) * 100))
}

type BadgeCardProps = {
  badge: Badge
}

export function BadgeCard({ badge }: BadgeCardProps) {
  const value = percent(badge)
  const label = badge.earned
    ? `${badge.name}. Conseguido`
    : `${badge.name}. ${badge.current} de ${badge.target}`

  return (
    <View
      accessible
      accessibilityLabel={label}
      style={[styles.card, badge.earned ? styles.cardEarned : styles.cardLocked]}
    >
      <View style={styles.top}>
        <View style={[styles.iconWrap, badge.earned ? styles.iconEarned : styles.iconLocked]}>
          <MaterialCommunityIcons
            name={toMcIcon(badge.icon)}
            size={24}
            color={badge.earned ? colors.primary : colors.mutedForeground}
          />
        </View>
        <Text style={styles.name} numberOfLines={2}>
          {badge.name}
        </Text>
      </View>

      <Text style={styles.description} numberOfLines={2}>
        {badge.description}
      </Text>

      {badge.earned ? (
        <View style={styles.earnedRow}>
          <MaterialCommunityIcons name="check-circle" size={14} color={colors.primary} />
          <Text style={styles.earnedText}>Conseguido</Text>
        </View>
      ) : (
        <View style={styles.progress}>
          <View
            accessibilityElementsHidden
            importantForAccessibility="no-hide-descendants"
            style={styles.track}
          >
            <View style={[styles.fill, { width: `${value}%` }]} />
          </View>
          <Text style={styles.counts}>
            {badge.current}/{badge.target}
          </Text>
        </View>
      )}
    </View>
  )
}

const styles = StyleSheet.create({
  card: {
    flex: 1,
    minHeight: 172,
    gap: spacing.sm,
    padding: spacing.md,
    borderWidth: 1,
    borderRadius: radius.md,
  },
  cardEarned: {
    backgroundColor: colors.card,
    borderColor: `${colors.primary}55`,
  },
  cardLocked: {
    backgroundColor: colors.muted,
    borderColor: colors.border,
  },
  top: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  iconWrap: {
    width: 40,
    height: 40,
    borderRadius: radius.sm,
    alignItems: 'center',
    justifyContent: 'center',
  },
  iconEarned: {
    backgroundColor: `${colors.primary}1F`,
  },
  iconLocked: {
    backgroundColor: `${colors.mutedForeground}1F`,
  },
  name: {
    flex: 1,
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  description: {
    flex: 1,
    fontSize: typography.size.xs,
    lineHeight: 18,
    color: colors.mutedForeground,
  },
  earnedRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  earnedText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.primary,
  },
  progress: {
    gap: spacing.xs,
  },
  track: {
    height: 6,
    borderRadius: radius.sm,
    backgroundColor: colors.background,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    borderRadius: radius.sm,
    backgroundColor: colors.primary,
  },
  counts: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
