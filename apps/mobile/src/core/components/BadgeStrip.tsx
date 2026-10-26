import { MaterialCommunityIcons } from '@expo/vector-icons'
import { StyleSheet, Text, View } from 'react-native'
import type { Badge, BadgesResponse } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'
import { toMcIcon } from '@/core/utils/icons'

const DEFAULT_MAX = 5

function badgeProgress(badge: Badge): number {
  if (badge.target <= 0) return 0
  return Math.min(100, Math.round((badge.current / badge.target) * 100))
}

/** Conseguidos primero y, dentro de los pendientes, el más avanzado. */
function orderBadges(badges: Array<Badge>): Array<Badge> {
  return [...badges].sort((a, b) => {
    if (a.earned !== b.earned) return a.earned ? -1 : 1
    return badgeProgress(b) - badgeProgress(a)
  })
}

type BadgeStripProps = {
  data: BadgesResponse
  /** Iconos visibles antes del chip "+N". */
  max?: number
  /** Añade "Siguiente: …" y la barra de progreso (Progreso). */
  detailed?: boolean
}

/**
 * Colección de logros en una fila. En Inicio va en modo compacto (solo iconos y
 * contador); en Progreso añade el siguiente logro y la barra. La navegación la
 * pone la cabecera de sección, no la tira.
 */
export function BadgeStrip({ data, max = DEFAULT_MAX, detailed = false }: BadgeStripProps) {
  if (data.total <= 0) return null

  const percent = Math.min(100, Math.round((data.earned / data.total) * 100))
  const ordered = orderBadges(data.badges)
  const visible = ordered.slice(0, max)
  const remaining = Math.max(0, ordered.length - visible.length)
  const next = ordered.find((badge) => !badge.earned)

  const counter =
    detailed && next
      ? `${data.earned} de ${data.total} · Siguiente: ${next.name} (${next.current}/${next.target})`
      : `${data.earned} de ${data.total}`

  return (
    <View style={styles.container}>
      <View style={styles.row}>
        {visible.map((badge) => (
          <View
            key={badge.id}
            accessible
            accessibilityLabel={`${badge.name}${badge.earned ? ', conseguido' : ''}`}
            style={[styles.icon, badge.earned ? styles.earned : styles.locked]}
          >
            <MaterialCommunityIcons
              name={toMcIcon(badge.icon)}
              size={16}
              color={badge.earned ? colors.primary : colors.mutedForeground}
            />
          </View>
        ))}

        {remaining > 0 ? (
          <View style={styles.more} accessibilityLabel={`${remaining} logros más`}>
            <Text style={styles.moreText}>+{remaining}</Text>
          </View>
        ) : null}

        <Text style={styles.counter} numberOfLines={detailed ? 2 : 1}>
          {counter}
        </Text>
      </View>

      {detailed ? (
        <View
          accessible
          accessibilityRole="progressbar"
          accessibilityLabel={`${percent}% de logros`}
          accessibilityValue={{ min: 0, max: 100, now: percent }}
          style={styles.track}
        >
          <View style={[styles.fill, { width: `${percent}%` }]} />
        </View>
      ) : null}
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    gap: spacing.sm,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  icon: {
    width: 28,
    height: 28,
    borderRadius: 14,
    alignItems: 'center',
    justifyContent: 'center',
  },
  earned: {
    backgroundColor: `${colors.primary}1F`,
  },
  locked: {
    backgroundColor: colors.muted,
  },
  more: {
    width: 28,
    height: 28,
    borderRadius: 14,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.muted,
  },
  moreText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
    color: colors.mutedForeground,
  },
  counter: {
    flex: 1,
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  track: {
    height: 4,
    borderRadius: radius.sm / 2,
    backgroundColor: colors.muted,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    borderRadius: radius.sm / 2,
    backgroundColor: colors.primary,
  },
})
