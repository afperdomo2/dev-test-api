import { StyleSheet, Text, View } from 'react-native'

import { colors, spacing, typography } from '@/core/theme/tokens'

type StatRowItem = {
  label: string
  value: string
}

type StatRowProps = {
  items: Array<StatRowItem>
}

/** Fila de métricas separadas por divisores (sustituye a las cards por KPI). */
export function StatRow({ items }: StatRowProps) {
  return (
    <View style={styles.row}>
      {items.map((item, index) => (
        <View key={item.label} style={styles.item}>
          {index > 0 ? <View style={styles.divider} /> : null}
          <Text style={styles.value}>{item.value}</Text>
          <Text style={styles.label}>{item.label}</Text>
        </View>
      ))}
    </View>
  )
}

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    alignItems: 'stretch',
  },
  item: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 2,
    paddingVertical: spacing.sm,
  },
  divider: {
    position: 'absolute',
    left: 0,
    top: spacing.sm,
    bottom: spacing.sm,
    width: 1,
    backgroundColor: colors.border,
  },
  value: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  label: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
