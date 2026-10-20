import * as SecureStore from 'expo-secure-store'

import { getToken, removeToken, setToken } from './storage'

describe('auth/storage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('getToken lee con la clave auth_token', async () => {
    jest.mocked(SecureStore.getItemAsync).mockResolvedValue('abc')

    await expect(getToken()).resolves.toBe('abc')
    expect(SecureStore.getItemAsync).toHaveBeenCalledWith('auth_token')
  })

  it('getToken devuelve null si no hay nada guardado', async () => {
    jest.mocked(SecureStore.getItemAsync).mockResolvedValue(null)

    await expect(getToken()).resolves.toBeNull()
  })

  it('setToken guarda con la clave auth_token', async () => {
    await setToken('nuevo')

    expect(SecureStore.setItemAsync).toHaveBeenCalledWith('auth_token', 'nuevo')
  })

  it('removeToken borra con la clave auth_token', async () => {
    await removeToken()

    expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith('auth_token')
  })
})
