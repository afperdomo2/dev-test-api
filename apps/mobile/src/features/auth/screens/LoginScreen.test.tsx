import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'

import LoginScreen from './LoginScreen'
import { useAuthStore } from '@/core/auth/auth.store'

describe('LoginScreen', () => {
  const login = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    useAuthStore.setState({ login, actionLoading: false })
  })

  it('muestra los campos y el botón', async () => {
    await render(<LoginScreen />)

    expect(screen.getByText('Email')).toBeOnTheScreen()
    expect(screen.getByText('Contraseña')).toBeOnTheScreen()
    expect(screen.getByLabelText('Iniciar sesión')).toBeOnTheScreen()
    // El título de la tarjeta y el botón comparten texto.
    expect(screen.getAllByText('Iniciar sesión')).toHaveLength(2)
  })

  it('valida los campos vacíos sin llamar al backend', async () => {
    await render(<LoginScreen />)

    await fireEvent.press(screen.getByLabelText('Iniciar sesión'))

    expect(screen.getByText('El email es obligatorio')).toBeOnTheScreen()
    expect(screen.getByText('La contraseña es obligatoria')).toBeOnTheScreen()
    expect(login).not.toHaveBeenCalled()
  })

  it('valida el formato del email', async () => {
    await render(<LoginScreen />)

    await fireEvent.changeText(screen.getByTestId('login-email'), 'no-es-email')
    await fireEvent.changeText(screen.getByTestId('login-password'), 'x')
    await fireEvent.press(screen.getByLabelText('Iniciar sesión'))

    expect(screen.getByText('Email no válido')).toBeOnTheScreen()
    expect(login).not.toHaveBeenCalled()
  })

  it('con datos válidos llama al login con el email recortado', async () => {
    login.mockResolvedValue(undefined)

    await render(<LoginScreen />)
    await fireEvent.changeText(screen.getByTestId('login-email'), '  ana@test.com  ')
    await fireEvent.changeText(screen.getByTestId('login-password'), 'secreta123')
    await fireEvent.press(screen.getByLabelText('Iniciar sesión'))

    await waitFor(() => expect(login).toHaveBeenCalledWith('ana@test.com', 'secreta123'))
  })
})
