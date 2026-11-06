import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type LevelLineProps = {
  level: number
  totalXp: number
  xpToLevel: number
}

/**
 * Progreso dentro del nivel actual. El backend manda el XP que falta
 * (`xpToLevel`), así que el rango del nivel se reconstruye con la misma
 * fórmula del servidor: nivel n empieza en (n-1)² * 100 XP.
 */
export function getLevelProgress(totalXp: number, level: number): number {
  const xpForCurrent = (level - 1) ** 2 * 100
  const xpForNext = level ** 2 * 100
  if (xpForNext <= xpForCurrent) return 1

  return Math.min(1, Math.max(0, (totalXp - xpForCurrent) / (xpForNext - xpForCurrent)))
}

export function LevelLine({ level, totalXp, xpToLevel }: LevelLineProps) {
  const value = Math.round(getLevelProgress(totalXp, level) * 100)
  const label = `Nivel ${level}. ${xpToLevel} XP para el nivel ${level + 1}`

  return (
    <View style={styles.container}>
      <View style={styles.row}>
        <View style={styles.badge}>
          <MaterialCommunityIcons name="lightning-bolt" size={12} color={colors.primary} />
          <Text style={styles.badgeText}>Nv. {level}</Text>
        </View>
        <Text style={styles.hint}>{xpToLevel} XP para el nivel {level + 1}</Text>
      </View>

      <View
        accessible
        accessibilityRole="progressbar"
        accessibilityLabel={label}
        accessibilityValue={{ min: 0, max: 100, now: value }}
        style={styles.track}
      >
        <View style={[styles.fill, { width: `${value}%` }]} />
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    gap: spacing.xs,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: radius.sm,
    backgroundColor: `${colors.primary}20`,
  },
  badgeText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.primary,
  },
  hint: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  track: {
    height: 4,
    borderRadius: 2,
    backgroundColor: colors.muted,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    borderRadius: 2,
    backgroundColor: colors.primary,
  },
})
