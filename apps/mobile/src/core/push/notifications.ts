import { Platform } from 'react-native'
import Constants from 'expo-constants'
import * as Device from 'expo-device'
import * as Notifications from 'expo-notifications'

const CHANNEL_ID = 'default'

export async function ensurePermission(): Promise<boolean> {
  const current = await Notifications.getPermissionsAsync()
  if (current.granted) return true

  const requested = await Notifications.requestPermissionsAsync()
  return requested.granted
}

async function ensureAndroidChannel(): Promise<void> {
  if (Platform.OS !== 'android') return

  await Notifications.setNotificationChannelAsync(CHANNEL_ID, {
    name: 'Avisos',
    importance: Notifications.AndroidImportance.DEFAULT,
  })
}

/**
 * Devuelve el token de Expo o null cuando no se puede usar push: emulador, sin
 * permisos, sin eas.projectId o plataforma no nativa. La app sigue funcionando.
 */
export async function getDevicePushToken(): Promise<string | null> {
  if (Platform.OS !== 'ios' && Platform.OS !== 'android') return null

  if (!Device.isDevice) {
    console.warn('⚠️ Push deshabilitado: no es un dispositivo físico')
    return null
  }

  const projectId = Constants.expoConfig?.extra?.eas?.projectId as string | undefined
  if (!projectId) {
    console.warn('⚠️ Push deshabilitado: falta eas.projectId')
    return null
  }

  await ensureAndroidChannel()

  const token = await Notifications.getExpoPushTokenAsync({ projectId })
  return token.data
}

/**
 * Configura cómo se muestran las notificaciones y navega a Inicio al tocar una
 * (incluido el arranque en frío). Devuelve la limpieza de los listeners.
 */
export function configureNotificationHandling(onOpen: () => void): () => void {
  Notifications.setNotificationHandler({
    handleNotification: async () => ({
      shouldShowBanner: true,
      shouldShowList: true,
      shouldPlaySound: false,
      shouldSetBadge: false,
    }),
  })

  let openedByListener = false
  const subscription = Notifications.addNotificationResponseReceivedListener(() => {
    openedByListener = true
    onOpen()
  })

  void Notifications.getLastNotificationResponseAsync()
    .then((response) => {
      if (response && !openedByListener) onOpen()
    })
    .catch(() => {
      // Sin respuesta previa: arranque normal.
    })

  return () => subscription.remove()
}
