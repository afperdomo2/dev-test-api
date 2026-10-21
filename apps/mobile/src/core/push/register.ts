import { Platform } from 'react-native'

import {
  registerToken,
  unregisterToken,
  type PushPlatform,
} from '@/core/api/services/push.service'

import { ensurePermission, getDevicePushToken } from './notifications'

// Token en memoria del proceso: permite darlo de baja en el logout sin
// persistirlo en disco (SPEC-mobile-push).
let currentToken: string | null = null
let registering = false

function nativePlatform(): PushPlatform | null {
  if (Platform.OS === 'ios' || Platform.OS === 'android') return Platform.OS
  return null
}

/**
 * Pide permiso, obtiene el token de Expo y lo registra en el backend.
 * Nunca lanza: sin push la app funciona igual (best-effort).
 */
export async function registerForPush(): Promise<void> {
  if (registering || currentToken) return

  const platform = nativePlatform()
  if (!platform) return

  registering = true
  try {
    const granted = await ensurePermission()
    if (!granted) return

    const token = await getDevicePushToken()
    if (!token) return

    await registerToken(token, platform)
    currentToken = token
  } catch (error) {
    console.warn('⚠️ No se pudo registrar el dispositivo para push:', error)
  } finally {
    registering = false
  }
}

/**
 * Da de baja el token de esta sesión. Nunca lanza; limpia la memoria incluso
 * si el backend falla.
 */
export async function unregisterFromPush(): Promise<void> {
  const token = currentToken
  if (!token) return

  try {
    await unregisterToken(token)
  } catch (error) {
    console.warn('⚠️ No se pudo dar de baja el dispositivo de push:', error)
  } finally {
    currentToken = null
  }
}
