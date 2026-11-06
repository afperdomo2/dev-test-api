import { StyleSheet, View } from 'react-native'
import type { TopicMasteryItem } from '@devtest/shared'

import { spacing } from '@/core/theme/tokens'

import { TopicRing } from './TopicRing'

type TopicRingGridProps = {
  items: Array<TopicMasteryItem>
}

/** Rejilla de 3 columnas: 9-12 temas caben en media pantalla. */
export function TopicRingGrid({ items }: TopicRingGridProps) {
  return (
    <View style={styles.grid}>
      {items.map((item) => (
        <TopicRing
          key={item.topicId}
          name={item.topicName}
          accuracy={item.accuracy}
          isWeak={item.isWeak}
          attempts={item.attempts}
          correctCount={item.correctCount}
        />
      ))}
    </View>
  )
}

const styles = StyleSheet.create({
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    rowGap: spacing.lg,
  },
})
