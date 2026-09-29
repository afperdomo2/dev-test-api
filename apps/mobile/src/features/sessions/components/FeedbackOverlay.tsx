import { useEffect, useRef } from 'react'
import { Animated, StyleSheet, Text, View } from 'react-native'
import type { Reward } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type FeedbackOverlayProps = {
  isCorrect: boolean
  xpGained: number
  reward?: Reward
  onDismiss: () => void
}

export function FeedbackOverlay({ isCorrect, xpGained, onDismiss }: FeedbackOverlayProps) {
  const opacity = useRef(new Animated.Value(0)).current
  const scale = useRef(new Animated.Value(0.8)).current

  useEffect(() => {
    Animated.parallel([
      Animated.timing(opacity, {
        toValue: 1,
        duration: 200,
        useNativeDriver: true,
      }),
      Animated.spring(scale, {
        toValue: 1,
        friction: 8,
        tension: 40,
        useNativeDriver: true,
      }),
    ]).start()

    const timer = setTimeout(() => {
      Animated.timing(opacity, {
        toValue: 0,
        duration: 300,
        useNativeDriver: true,
      }).start(() => onDismiss())
    }, 1500)

    return () => clearTimeout(timer)
  }, [opacity, scale, onDismiss])

  return (
    <Animated.View style={[styles.overlay, { opacity }]}>
      <Animated.View style={[styles.card, { transform: [{ scale }] }]}>
        <Text style={[styles.icon, isCorrect ? styles.iconCorrect : styles.iconIncorrect]}>
          {isCorrect ? '✓' : '✗'}
        </Text>
        <Text style={styles.title}>{isCorrect ? 'Correcto' : 'Incorrecto'}</Text>
        {xpGained > 0 && <Text style={styles.xp}>+{xpGained} XP</Text>}
      </Animated.View>
    </Animated.View>
  )
}

const styles = StyleSheet.create({
  overlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.7)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 100,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    padding: spacing.xl,
    alignItems: 'center',
    gap: spacing.md,
    minWidth: 200,
  },
  icon: {
    fontSize: 48,
    fontWeight: typography.weight.bold,
  },
  iconCorrect: {
    color: colors.primary,
  },
  iconIncorrect: {
    color: colors.destructive,
  },
  title: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  xp: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.semibold,
    color: colors.accent,
  },
})
