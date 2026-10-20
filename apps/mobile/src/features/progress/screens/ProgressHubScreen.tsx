import { useRouter } from 'expo-router'
import { ScrollView, StyleSheet, Text, View } from 'react-native'

import { colors, spacing, typography } from '@/core/theme/tokens'

import { ProgressHubRow } from '../components/ProgressHubRow'

export default function ProgressHubScreen() {
  const router = useRouter()

  return (
    <ScrollView
      style={styles.container}
      contentInsetAdjustmentBehavior="automatic"
      contentContainerStyle={styles.content}
    >
      <Text style={styles.sectionTitle}>Tu avance</Text>

      <View style={styles.list}>
        <ProgressHubRow
          icon="chart-timeline-variant"
          title="Temas"
          subtitle="Tu porcentaje de acierto en cada tema"
          onPress={() => router.push('/(tabs)/progreso/temas')}
        />
      </View>
    </ScrollView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  sectionTitle: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  list: {
    gap: spacing.sm,
  },
})
