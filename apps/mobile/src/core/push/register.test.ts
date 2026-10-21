import { registerForPush, unregisterFromPush } from './register'
import { ensurePermission, getDevicePushToken } from './notifications'
import { registerToken, unregisterToken } from '@/core/api/services/push.service'

jest.mock('./notifications', () => ({
  ensurePermission: jest.fn(),
  getDevicePushToken: jest.fn(),
}))

jest.mock('@/core/api/services/push.service', () => ({
  registerToken: jest.fn(),
  unregisterToken: jest.fn(),
}))

const mockedPermission = jest.mocked(ensurePermission)
const mockedToken = jest.mocked(getDevicePushToken)
const mockedRegister = jest.mocked(registerToken)
const mockedUnregister = jest.mocked(unregisterToken)

describe('registerForPush', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    jest.spyOn(console, 'warn').mockImplementation(() => undefined)
    mockedPermission.mockResolvedValue(true)
    mockedToken.mockResolvedValue('ExponentPushToken[abc]')
    mockedRegister.mockResolvedValue()
    mockedUnregister.mockResolvedValue()
  })

  afterEach(async () => {
    await unregisterFromPush()
    jest.clearAllMocks()
    jest.restoreAllMocks()
  })

  it('registra el token con la plataforma', async () => {
    await registerForPush()

    expect(mockedRegister).toHaveBeenCalledWith('ExponentPushToken[abc]', 'ios')
  })

  it('no registra sin permiso', async () => {
    mockedPermission.mockResolvedValue(false)

    await registerForPush()

    expect(mockedToken).not.toHaveBeenCalled()
    expect(mockedRegister).not.toHaveBeenCalled()
  })

  it('no registra si no hay token (emulador o sin projectId)', async () => {
    mockedToken.mockResolvedValue(null)

    await registerForPush()

    expect(mockedRegister).not.toHaveBeenCalled()
  })

  it('no lanza si el backend falla', async () => {
    mockedRegister.mockRejectedValue(new Error('sin red'))

    await expect(registerForPush()).resolves.toBeUndefined()
  })

  it('no registra dos veces en la misma sesión', async () => {
    await registerForPush()
    await registerForPush()

    expect(mockedRegister).toHaveBeenCalledTimes(1)
  })

  it('unregister da de baja y limpia el token en memoria', async () => {
    await registerForPush()
    await unregisterFromPush()

    expect(mockedUnregister).toHaveBeenCalledWith('ExponentPushToken[abc]')

    // Tras la baja, una nueva sesión puede volver a registrarse.
    mockedRegister.mockClear()
    await registerForPush()
    expect(mockedRegister).toHaveBeenCalledTimes(1)
  })

  it('unregister sin token no llama al backend', async () => {
    await unregisterFromPush()

    expect(mockedUnregister).not.toHaveBeenCalled()
  })
})
