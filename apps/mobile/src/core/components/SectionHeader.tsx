import { MaterialCommunityIcons } from '@expo/vector-icons'
import { Pressable, StyleSheet, Text, View } from 'react-native'

import { colors, spacing, typography } from '@/core/theme/tokens'

type SectionHeaderProps = {
  title: string
  /** Texto de la acción; por defecto "Ver todos". */
  actionLabel?: string
  onPress?: () => void
}

/**
 * Cabecera de sección (mayúsculas + acción opcional). Compartida por Inicio y
 * Progreso para que un bloque con "Ver todos" se vea igual en las dos pantallas.
 */
export function SectionHeader({ title, actionLabel = 'Ver todos', onPress }: SectionHeaderProps) {
  return (
    <View style={styles.header}>
      <Text style={styles.title}>{title}</Text>

      {onPress ? (
        <Pressable
          onPress={onPress}
          hitSlop={{ top: 10, bottom: 10, left: 8, right: 8 }}
          style={({ pressed }) => [styles.action, pressed && styles.pressed]}
          accessibilityRole="button"
          accessibilityLabel={`${actionLabel} ${title.toLowerCase()}`}
        >
          <Text style={styles.actionText}>{actionLabel}</Text>
          <MaterialCommunityIcons name="chevron-right" size={16} color={colors.mutedForeground} />
        </Pressable>
      ) : null}
    </View>
  )
}

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    minHeight: 24,
    paddingHorizontal: spacing.xs,
  },
  title: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  action: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 2,
  },
  pressed: {
    opacity: 0.7,
  },
  actionText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
    color: colors.mutedForeground,
  },
})
