import { Stack } from 'expo-router'

import { colors } from '@/core/theme/tokens'

export default function ProgressLayout() {
  return (
    <Stack
      screenOptions={{
        headerStyle: { backgroundColor: colors.background },
        headerTintColor: colors.foreground,
        contentStyle: { backgroundColor: colors.background },
      }}
    >
      <Stack.Screen name="index" options={{ title: 'Progreso' }} />
      <Stack.Screen name="temas" options={{ title: 'Temas' }} />
      <Stack.Screen name="semana" options={{ title: 'Mi semana' }} />
    </Stack>
  )
}
