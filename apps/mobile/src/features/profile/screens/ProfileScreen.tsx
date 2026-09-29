import { Pressable, StyleSheet, Text, View } from 'react-native'

import { Screen } from '@/core/components/Screen'
import { useAuthStore } from '@/core/auth/auth.store'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

export default function ProfileScreen() {
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)

  return (
    <Screen title="Perfil">
      <View style={styles.info}>
        <Text style={styles.label}>Email</Text>
        <Text style={styles.value}>{user?.email}</Text>
      </View>

      <View style={styles.info}>
        <Text style={styles.label}>Rol</Text>
        <Text style={styles.value}>{user?.isAdmin ? 'Administrador' : 'Usuario'}</Text>
      </View>

      <Pressable style={styles.logoutButton} onPress={() => void logout()}>
        <Text style={styles.logoutLabel}>Cerrar sesion</Text>
      </Pressable>
    </Screen>
  )
}

const styles = StyleSheet.create({
  info: {
    backgroundColor: colors.card,
    borderRadius: radius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.xs,
  },
  label: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  value: {
    fontSize: typography.size.md,
    color: colors.foreground,
    fontWeight: typography.weight.medium,
  },
  logoutButton: {
    backgroundColor: colors.destructive,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    alignItems: 'center',
    marginTop: spacing.lg,
  },
  logoutLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: '#FFFFFF',
  },
})
