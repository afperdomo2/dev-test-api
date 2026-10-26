import { Platform } from 'react-native'
import Constants, { ExecutionEnvironment } from 'expo-constants'
import * as Device from 'expo-device'

type NotificationsModule = typeof import('expo-notifications')

const CHANNEL_ID = 'default'

let cachedModule: NotificationsModule | null | undefined

/**
 * Expo Go en Android (SDK 57) no registra el nativo `ExpoTopicSubscriptionModule`
 * y `expo-notifications` lo resuelve al importarse, así que un import estático
 * deja la app en pantalla roja antes de renderizar (expo/expo#52175). En ese
 * entorno el push remoto tampoco está soportado desde SDK 53, por eso ni
 * siquiera cargamos el módulo.
 */
function isExpoGoAndroid(): boolean {
  return (
    Platform.OS === 'android' &&
    Constants.executionEnvironment === ExecutionEnvironment.StoreClient
  )
}

/**
 * Carga `expo-notifications` de forma diferida y opcional: si el entorno no lo
 * soporta (Expo Go en Android) o el import falla, el push se desactiva sin
 * romper el arranque de la app.
 */
function loadNotifications(): NotificationsModule | null {
  if (cachedModule !== undefined) return cachedModule

  if (isExpoGoAndroid()) {
    cachedModule = null
    return cachedModule
  }

  try {
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    cachedModule = require('expo-notifications') as NotificationsModule
  } catch (error) {
    console.warn('⚠️ Push deshabilitado: expo-notifications no está disponible:', error)
    cachedModule = null
  }

  return cachedModule
}

export async function ensurePermission(): Promise<boolean> {
  const Notifications = loadNotifications()
  if (!Notifications) return false

  const current = await Notifications.getPermissionsAsync()
  if (current.granted) return true

  const requested = await Notifications.requestPermissionsAsync()
  return requested.granted
}

async function ensureAndroidChannel(Notifications: NotificationsModule): Promise<void> {
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

  const Notifications = loadNotifications()
  if (!Notifications) {
    console.warn('⚠️ Push deshabilitado: expo-notifications no está disponible en este entorno')
    return null
  }

  if (!Device.isDevice) {
    console.warn('⚠️ Push deshabilitado: no es un dispositivo físico')
    return null
  }

  const projectId = Constants.expoConfig?.extra?.eas?.projectId as string | undefined
  if (!projectId) {
    console.warn('⚠️ Push deshabilitado: falta eas.projectId')
    return null
  }

  await ensureAndroidChannel(Notifications)

  const token = await Notifications.getExpoPushTokenAsync({ projectId })
  return token.data
}

/**
 * Configura cómo se muestran las notificaciones y navega a Inicio al tocar una
 * (incluido el arranque en frío). Devuelve la limpieza de los listeners.
 */
export function configureNotificationHandling(onOpen: () => void): () => void {
  const Notifications = loadNotifications()
  if (!Notifications) return () => undefined

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
