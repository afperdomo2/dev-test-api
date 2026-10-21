import { DarkTheme, SplashScreen, Stack, ThemeProvider, useRouter, useSegments } from 'expo-router'
import { StatusBar } from 'expo-status-bar'
import { useEffect } from 'react'
import { Pressable, StyleSheet, Text, View } from 'react-native'
import { GestureHandlerRootView } from 'react-native-gesture-handler'
import { BottomSheetModalProvider } from '@gorhom/bottom-sheet'

import { useAuthStore } from '@/core/auth/auth.store'
import { API_BASE_URL } from '@/core/config'
import { configureNotificationHandling } from '@/core/push/notifications'
import { registerForPush } from '@/core/push/register'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

void SplashScreen.preventAutoHideAsync()

function AuthGate() {
  const router = useRouter()
  const segments = useSegments()
  const ready = useAuthStore((s) => s.ready)
  const needsSetup = useAuthStore((s) => s.needsSetup)
  const token = useAuthStore((s) => s.token)
  const bootstrap = useAuthStore((s) => s.bootstrap)

  useEffect(() => {
    void bootstrap()
  }, [bootstrap])

  // Registro de push una vez por sesión autenticada (nunca bloquea el arranque).
  useEffect(() => {
    if (token) void registerForPush()
  }, [token])

  // Tocar una notificación (o abrir la app desde una) lleva a Inicio.
  useEffect(() => configureNotificationHandling(() => router.push('/')), [router])

  useEffect(() => {
    if (!ready || needsSetup === null) return
    const inAuthGroup = segments[0] === 'login' || segments[0] === 'setup'
    if (needsSetup && segments[0] !== 'setup') {
      router.replace('/setup')
    } else if (!needsSetup && !token && segments[0] !== 'login') {
      router.replace('/login')
    } else if (!needsSetup && token && inAuthGroup) {
      router.replace('/(tabs)')
    }
  }, [ready, needsSetup, token, segments, router])

  useEffect(() => {
    if (ready) {
      void SplashScreen.hideAsync()
    }
  }, [ready])

  if (!ready) return null

  if (needsSetup === null) {
    return (
      <View style={styles.retry}>
        <Text style={styles.retryTitle}>No se pudo conectar con el servidor</Text>
        <Text style={styles.retryUrl}>{API_BASE_URL}</Text>
        <Pressable
          style={styles.retryButton}
          onPress={() => void bootstrap()}
          accessibilityRole="button"
        >
          <Text style={styles.retryButtonLabel}>Reintentar</Text>
        </Pressable>
      </View>
    )
  }

  return <Stack screenOptions={{ headerShown: false }} />
}

export default function RootLayout() {
  return (
    <GestureHandlerRootView style={styles.root}>
      <ThemeProvider value={DarkTheme}>
        <StatusBar style="light" />
        <BottomSheetModalProvider>
          <AuthGate />
        </BottomSheetModalProvider>
      </ThemeProvider>
    </GestureHandlerRootView>
  )
}

const styles = StyleSheet.create({
  root: {
    flex: 1,
  },
  retry: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    padding: spacing.lg,
    gap: spacing.sm,
  },
  retryTitle: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
    textAlign: 'center',
  },
  retryUrl: {
    fontSize: typography.size.sm,
    fontFamily: typography.family.mono,
    color: colors.mutedForeground,
    textAlign: 'center',
  },
  retryButton: {
    marginTop: spacing.sm,
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
  },
  retryButtonLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
})
