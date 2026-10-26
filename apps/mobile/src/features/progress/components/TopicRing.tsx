import { StyleSheet, Text, View } from 'react-native'
import Svg, { Circle } from 'react-native-svg'

import { colors, spacing, typography } from '@/core/theme/tokens'
import { toPercent } from '@/core/utils/numbers'

const SIZE = 64
const STROKE = 6
const RADIUS = (SIZE - STROKE) / 2
const CIRCUMFERENCE = 2 * Math.PI * RADIUS

type TopicRingProps = {
  name: string
  /** 0..1 */
  accuracy: number
  isWeak?: boolean
  attempts?: number
  correctCount?: number
}

/**
 * Dominio de un tema en formato anillo (estilo Duolingo). El porcentaje va
 * escrito dentro del anillo y el tema débil se marca con el color del arco y del
 * número (`accent`), sin icono: el % ya es el dato y "tema débil" sigue en la
 * etiqueta accesible. Para comparar temas *dentro de un periodo* se usa
 * `TopicCompactRow`, donde el chip "Débil" sí lleva texto.
 */
export function TopicRing({ name, accuracy, isWeak = false, attempts, correctCount }: TopicRingProps) {
  const value = toPercent(accuracy)
  const filled = (value / 100) * CIRCUMFERENCE
  const strokeColor = isWeak ? colors.accent : colors.primary
  const counts =
    attempts !== undefined && correctCount !== undefined
      ? `, ${correctCount} de ${attempts} aciertos`
      : ''
  const label = `${name}: ${value}% de acierto${isWeak ? ', tema débil' : ''}${counts}`

  return (
    <View
      accessible
      accessibilityRole="progressbar"
      accessibilityLabel={label}
      accessibilityValue={{ min: 0, max: 100, now: value }}
      style={styles.item}
    >
      <View style={styles.ring}>
        <Svg width={SIZE} height={SIZE}>
          <Circle
            cx={SIZE / 2}
            cy={SIZE / 2}
            r={RADIUS}
            stroke={colors.muted}
            strokeWidth={STROKE}
            fill="none"
          />
          <Circle
            cx={SIZE / 2}
            cy={SIZE / 2}
            r={RADIUS}
            stroke={strokeColor}
            strokeWidth={STROKE}
            fill="none"
            strokeLinecap="round"
            strokeDasharray={`${CIRCUMFERENCE} ${CIRCUMFERENCE}`}
            strokeDashoffset={CIRCUMFERENCE - filled}
            transform={`rotate(-90 ${SIZE / 2} ${SIZE / 2})`}
          />
        </Svg>

        <Text style={[styles.value, isWeak && styles.valueWeak]}>{value}%</Text>
      </View>

      <Text style={styles.name} numberOfLines={1}>
        {name}
      </Text>
    </View>
  )
}

const styles = StyleSheet.create({
  item: {
    width: '33.333%',
    alignItems: 'center',
    gap: spacing.xs + 2,
  },
  ring: {
    width: SIZE,
    height: SIZE,
    alignItems: 'center',
    justifyContent: 'center',
  },
  value: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  valueWeak: {
    color: colors.accent,
  },
  name: {
    maxWidth: '100%',
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
    textAlign: 'center',
  },
})
