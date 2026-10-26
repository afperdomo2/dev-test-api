import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'

import SetupScreen from './SetupScreen'
import { useAuthStore } from '@/core/auth/auth.store'

describe('SetupScreen', () => {
  const setup = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    useAuthStore.setState({ setup, actionLoading: false, setupRequiresToken: true })
  })

  it('valida los campos vacíos', async () => {
    await render(<SetupScreen />)

    await fireEvent.press(screen.getByLabelText('Crear cuenta'))

    expect(screen.getByText('El email es obligatorio')).toBeOnTheScreen()
    expect(screen.getByText('La contraseña es obligatoria')).toBeOnTheScreen()
    expect(screen.getByText('Confirma la contraseña')).toBeOnTheScreen()
    expect(setup).not.toHaveBeenCalled()
  })

  it('avisa si las contraseñas no coinciden', async () => {
    await render(<SetupScreen />)

    await fireEvent.changeText(screen.getByTestId('setup-email'), 'admin@test.com')
    await fireEvent.changeText(screen.getByTestId('setup-password'), 'secreta123')
    await fireEvent.changeText(screen.getByTestId('setup-confirm'), 'otra12345')
    await fireEvent.press(screen.getByLabelText('Crear cuenta'))

    expect(screen.getByText('Las contraseñas no coinciden')).toBeOnTheScreen()
    expect(setup).not.toHaveBeenCalled()
  })

  it('exige el token de instalación cuando el backend lo pide', async () => {
    await render(<SetupScreen />)

    await fireEvent.changeText(screen.getByTestId('setup-email'), 'admin@test.com')
    await fireEvent.changeText(screen.getByTestId('setup-password'), 'secreta123')
    await fireEvent.changeText(screen.getByTestId('setup-confirm'), 'secreta123')
    await fireEvent.press(screen.getByLabelText('Crear cuenta'))

    expect(screen.getByText('El token de instalación es obligatorio')).toBeOnTheScreen()
    expect(setup).not.toHaveBeenCalled()
  })

  it('con datos válidos llama al setup con email recortado y token', async () => {
    setup.mockResolvedValue(undefined)

    await render(<SetupScreen />)
    await fireEvent.changeText(screen.getByTestId('setup-email'), '  admin@test.com  ')
    await fireEvent.changeText(screen.getByTestId('setup-password'), 'secreta123')
    await fireEvent.changeText(screen.getByTestId('setup-confirm'), 'secreta123')
    await fireEvent.changeText(screen.getByTestId('setup-token'), 'install-token')
    await fireEvent.press(screen.getByLabelText('Crear cuenta'))

    await waitFor(() =>
      expect(setup).toHaveBeenCalledWith('admin@test.com', 'secreta123', 'install-token'),
    )
  })

  it('muestra el aviso y bloquea el submit cuando el setup está deshabilitado', async () => {
    useAuthStore.setState({ setup, actionLoading: false, setupRequiresToken: false })

    await render(<SetupScreen />)

    expect(screen.getByTestId('setup-disabled')).toBeOnTheScreen()
    expect(screen.queryByTestId('setup-token')).toBeNull()

    await fireEvent.press(screen.getByLabelText('Crear cuenta'))
    expect(setup).not.toHaveBeenCalled()
  })
})
