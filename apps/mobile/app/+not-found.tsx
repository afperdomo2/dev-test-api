import { Link } from 'expo-router'
import { Text, View } from 'react-native'

import { colors, spacing } from '@/core/theme/tokens'

export default function NotFound() {
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: colors.background,
        alignItems: 'center',
        justifyContent: 'center',
        gap: spacing.md,
      }}
    >
      <Text style={{ color: colors.foreground, fontSize: 20, fontWeight: '600' }}>
        Página no encontrada
      </Text>
      <Link href="/" style={{ color: colors.primary, fontSize: 16 }}>
        Volver a Inicio
      </Link>
    </View>
  )
}