import { Alert, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import Constants from 'expo-constants'

import { useAuthStore } from '@/core/auth/auth.store'
import { unregisterFromPush } from '@/core/push/register'
import { formatDate } from '@devtest/shared'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

function getInitials(email: string): string {
  return email.slice(0, 2).toUpperCase()
}

export default function ProfileScreen() {
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)

  const handleLogout = () => {
    Alert.alert(
      'Cerrar sesion',
      '¿Seguro que deseas cerrar sesion?',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Cerrar sesion',
          style: 'destructive',
          onPress: () => {
            void (async () => {
              await unregisterFromPush()
              await logout()
            })()
          },
        },
      ],
    )
  }

  if (!user) return null

  const isAdmin = user.isAdmin
  const avatarColor = isAdmin ? colors.primary : colors.mutedForeground
  const initials = getInitials(user.email)
  const version = Constants.expoConfig?.version ?? '1.0.0'

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      contentInsetAdjustmentBehavior="automatic"
    >
      <Text style={styles.screenTitle}>Perfil</Text>

      <View style={styles.header}>
        <View style={[styles.avatar, { backgroundColor: avatarColor }]}>
          <Text style={styles.avatarText}>{initials}</Text>
        </View>
        <Text style={styles.headerEmail}>{user.email}</Text>
        <View style={[styles.roleBadge, { backgroundColor: `${avatarColor}20` }]}>
          <MaterialCommunityIcons
            name={isAdmin ? 'shield-account' : 'account'}
            size={12}
            color={avatarColor}
          />
          <Text style={[styles.roleBadgeText, { color: avatarColor }]}>
            {isAdmin ? 'Administrador' : 'Usuario'}
          </Text>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Cuenta</Text>

        <InfoRow
          icon="email"
          label="Email"
          value={user.email}
        />
        <InfoRow
          icon="shield-account"
          label="Rol"
          value={isAdmin ? 'Administrador' : 'Usuario'}
        />
        <InfoRow
          icon="calendar"
          label="Miembro desde"
          value={formatDate(user.createdAt)}
        />
        <InfoRow
          icon="identifier"
          label="ID de usuario"
          value={user.id.length > 20 ? `${user.id.slice(0, 20)}...` : user.id}
        />
        <InfoRow
          icon="cellphone"
          label="Version"
          value={version}
          isLast
        />
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Sesion</Text>

        <Pressable
          style={({ pressed }) => [styles.infoRow, pressed && styles.infoRowPressed]}
          onPress={handleLogout}
          accessibilityRole="button"
          accessibilityLabel="Cerrar sesion"
        >
          <MaterialCommunityIcons name="logout" size={20} color={colors.destructive} />
          <Text style={styles.logoutLabel}>Cerrar sesion</Text>
          <MaterialCommunityIcons
            name="chevron-right"
            size={18}
            color={colors.mutedForeground}
          />
        </Pressable>
      </View>
    </ScrollView>
  )
}

function InfoRow({
  icon,
  label,
  value,
  isLast,
}: {
  icon: keyof typeof MaterialCommunityIcons.glyphMap
  label: string
  value: string
  isLast?: boolean
}) {
  return (
    <View style={[styles.infoRow, !isLast && styles.infoRowBorder]}>
      <View style={styles.infoRowLeft}>
        <MaterialCommunityIcons name={icon} size={20} color={colors.mutedForeground} />
        <Text style={styles.infoLabel}>{label}</Text>
      </View>
      <Text style={styles.infoValue} numberOfLines={1}>
        {value}
      </Text>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing.lg,
    paddingBottom: spacing.xl,
    gap: spacing.lg,
  },
  screenTitle: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    letterSpacing: -0.5,
    color: colors.foreground,
  },
  header: {
    alignItems: 'center',
    gap: spacing.sm,
    paddingVertical: spacing.lg,
  },
  avatar: {
    width: 72,
    height: 72,
    borderRadius: 36,
    justifyContent: 'center',
    alignItems: 'center',
  },
  avatarText: {
    fontSize: typography.size['2xl'],
    fontWeight: typography.weight.bold,
    color: '#FFFFFF',
  },
  headerEmail: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  roleBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.sm,
    paddingVertical: 4,
    borderRadius: radius.sm,
    gap: spacing.xs,
  },
  roleBadgeText: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.medium,
  },
  section: {
    gap: spacing.xs,
  },
  sectionTitle: {
    fontSize: typography.size.xs,
    fontWeight: typography.weight.semibold,
    color: colors.mutedForeground,
    textTransform: 'uppercase',
    letterSpacing: 1,
    paddingHorizontal: spacing.xs,
    marginBottom: spacing.xs,
  },
  infoRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: colors.card,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    minHeight: 48,
    borderWidth: 1,
    borderColor: colors.border,
  },
  infoRowBorder: {
    borderBottomWidth: 0,
  },
  infoRowLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    flex: 1,
  },
  infoRowPressed: {
    opacity: 0.7,
  },
  infoLabel: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  infoValue: {
    fontSize: typography.size.sm,
    color: colors.foreground,
    fontWeight: typography.weight.medium,
    maxWidth: '60%',
    textAlign: 'right',
  },
  logoutLabel: {
    flex: 1,
    fontSize: typography.size.md,
    fontWeight: typography.weight.medium,
    color: colors.destructive,
  },
})
