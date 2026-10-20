import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'

import SetupScreen from './SetupScreen'
import { useAuthStore } from '@/core/auth/auth.store'

describe('SetupScreen', () => {
  const setup = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    useAuthStore.setState({ setup, actionLoading: false })
  })

  it('valida los tres campos vacíos', async () => {
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

  it('con datos válidos llama al setup con el email recortado', async () => {
    setup.mockResolvedValue(undefined)

    await render(<SetupScreen />)
    await fireEvent.changeText(screen.getByTestId('setup-email'), '  admin@test.com  ')
    await fireEvent.changeText(screen.getByTestId('setup-password'), 'secreta123')
    await fireEvent.changeText(screen.getByTestId('setup-confirm'), 'secreta123')
    await fireEvent.press(screen.getByLabelText('Crear cuenta'))

    await waitFor(() => expect(setup).toHaveBeenCalledWith('admin@test.com', 'secreta123'))
  })
})
