import { Platform } from 'react-native'
import Constants from 'expo-constants'
import * as Device from 'expo-device'
import * as Notifications from 'expo-notifications'

import {
  configureNotificationHandling,
  ensurePermission,
  getDevicePushToken,
} from './notifications'

jest.mock('expo-notifications', () => ({
  getPermissionsAsync: jest.fn(),
  requestPermissionsAsync: jest.fn(),
  getExpoPushTokenAsync: jest.fn(),
  setNotificationChannelAsync: jest.fn(),
  setNotificationHandler: jest.fn(),
  addNotificationResponseReceivedListener: jest.fn(),
  getLastNotificationResponseAsync: jest.fn(),
  AndroidImportance: { DEFAULT: 3 },
}))

jest.mock('expo-device', () => ({ isDevice: true }))

jest.mock('expo-constants', () => ({
  __esModule: true,
  default: { expoConfig: { extra: {} } },
}))

const constants = Constants as unknown as {
  expoConfig: { extra: Record<string, unknown> }
}

function setProjectId(projectId: string | undefined) {
  constants.expoConfig.extra = projectId ? { eas: { projectId } } : {}
}

describe('ensurePermission', () => {
  beforeEach(() => jest.clearAllMocks())

  it('no vuelve a pedir permiso si ya está concedido', async () => {
    jest.mocked(Notifications.getPermissionsAsync).mockResolvedValue({ granted: true } as never)

    await expect(ensurePermission()).resolves.toBe(true)
    expect(Notifications.requestPermissionsAsync).not.toHaveBeenCalled()
  })

  it('pide permiso cuando falta', async () => {
    jest.mocked(Notifications.getPermissionsAsync).mockResolvedValue({ granted: false } as never)
    jest.mocked(Notifications.requestPermissionsAsync).mockResolvedValue({ granted: true } as never)

    await expect(ensurePermission()).resolves.toBe(true)
    expect(Notifications.requestPermissionsAsync).toHaveBeenCalled()
  })

  it('devuelve false si el usuario deniega', async () => {
    jest.mocked(Notifications.getPermissionsAsync).mockResolvedValue({ granted: false } as never)
    jest.mocked(Notifications.requestPermissionsAsync).mockResolvedValue({ granted: false } as never)

    await expect(ensurePermission()).resolves.toBe(false)
  })
})

describe('getDevicePushToken', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    ;(Device as { isDevice: boolean }).isDevice = true
    setProjectId(undefined)
    jest.spyOn(console, 'warn').mockImplementation(() => undefined)
    jest.mocked(Notifications.getExpoPushTokenAsync).mockResolvedValue({ data: 'tok' } as never)
    jest.mocked(Notifications.setNotificationChannelAsync).mockResolvedValue(null)
  })

  afterEach(() => {
    jest.restoreAllMocks()
  })

  it('devuelve null en emulador', async () => {
    ;(Device as { isDevice: boolean }).isDevice = false

    await expect(getDevicePushToken()).resolves.toBeNull()
    expect(Notifications.getExpoPushTokenAsync).not.toHaveBeenCalled()
  })

  it('devuelve null y avisa si falta eas.projectId', async () => {
    await expect(getDevicePushToken()).resolves.toBeNull()

    expect(Notifications.getExpoPushTokenAsync).not.toHaveBeenCalled()
    expect(console.warn).toHaveBeenCalledWith(expect.stringContaining('eas.projectId'))
  })

  it('crea el canal de Android antes de pedir el token', async () => {
    jest.replaceProperty(Platform, 'OS', 'android')
    setProjectId('proj-1')

    await expect(getDevicePushToken()).resolves.toBe('tok')

    expect(Notifications.setNotificationChannelAsync).toHaveBeenCalledWith(
      'default',
      expect.objectContaining({ importance: Notifications.AndroidImportance.DEFAULT }),
    )
    expect(Notifications.getExpoPushTokenAsync).toHaveBeenCalledWith({ projectId: 'proj-1' })
  })

  it('devuelve el token con projectId en iOS', async () => {
    setProjectId('proj-1')

    await expect(getDevicePushToken()).resolves.toBe('tok')
    expect(Notifications.setNotificationChannelAsync).not.toHaveBeenCalled()
  })
})

describe('configureNotificationHandling', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    jest.mocked(Notifications.getLastNotificationResponseAsync).mockResolvedValue(null)
  })

  it('configura el handler y limpia el listener', () => {
    const remove = jest.fn()
    jest
      .mocked(Notifications.addNotificationResponseReceivedListener)
      .mockReturnValue({ remove } as never)

    const cleanup = configureNotificationHandling(() => undefined)

    expect(Notifications.setNotificationHandler).toHaveBeenCalled()
    expect(Notifications.addNotificationResponseReceivedListener).toHaveBeenCalled()

    cleanup()
    expect(remove).toHaveBeenCalled()
  })

  it('abre Inicio al arrancar en frío desde una notificación', async () => {
    jest
      .mocked(Notifications.addNotificationResponseReceivedListener)
      .mockReturnValue({ remove: jest.fn() } as never)
    jest.mocked(Notifications.getLastNotificationResponseAsync).mockResolvedValue({} as never)

    const onOpen = jest.fn()
    configureNotificationHandling(onOpen)

    await Promise.resolve()
    await Promise.resolve()

    expect(onOpen).toHaveBeenCalledTimes(1)
  })
})
