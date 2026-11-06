import { render, screen } from '@testing-library/react-native'

import ProfileScreen from './ProfileScreen'
import { useAuthStore } from '@/core/auth/auth.store'

describe('ProfileScreen', () => {
  beforeEach(() => {
    useAuthStore.setState({
      token: 'tok',
      user: {
        id: '3f1c2a9e-1234-4abc-9def-567890abcdef',
        email: 'a@b.com',
        isAdmin: false,
        dailyImportLimit: 200,
        createdAt: '2026-08-28T12:00:00Z',
        updatedAt: '2026-08-28T12:00:00Z',
      },
    })
  })

  it('muestra los datos del usuario', async () => {
    await render(<ProfileScreen />)

    expect(screen.getByText('Perfil')).toBeOnTheScreen()
    expect(screen.getAllByText('a@b.com').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Usuario').length).toBeGreaterThan(0)
    expect(screen.getByText('Miembro desde')).toBeOnTheScreen()
    expect(screen.getByText('ID de usuario')).toBeOnTheScreen()
  })

  it('incluye la sección de seguridad con el cambio de contraseña', async () => {
    await render(<ProfileScreen />)

    expect(screen.getByText('Seguridad')).toBeOnTheScreen()
    expect(screen.getByText('Cambiar contraseña')).toBeOnTheScreen()
    expect(screen.getByLabelText('Actualizar contraseña')).toBeOnTheScreen()
  })
})
