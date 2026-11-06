import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'

import { ChangePasswordSection } from './ChangePasswordSection'
import { changePassword } from '@/core/api/services/auth.service'

jest.mock('@/core/api/services/auth.service', () => ({
  changePassword: jest.fn(),
}))

const mockedChangePassword = jest.mocked(changePassword)

async function fillForm(current: string, next: string, confirm: string) {
  await fireEvent.changeText(screen.getByTestId('change-password-current'), current)
  await fireEvent.changeText(screen.getByTestId('change-password-new'), next)
  await fireEvent.changeText(screen.getByTestId('change-password-confirm'), confirm)
}

describe('ChangePasswordSection', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('muestra los campos y el botón', async () => {
    await render(<ChangePasswordSection />)

    expect(screen.getByText('Contraseña actual')).toBeOnTheScreen()
    expect(screen.getByText('Nueva contraseña')).toBeOnTheScreen()
    expect(screen.getByText('Confirmar nueva contraseña')).toBeOnTheScreen()
    expect(screen.getByLabelText('Actualizar contraseña')).toBeOnTheScreen()
  })

  it('valida los campos vacíos sin llamar al backend', async () => {
    await render(<ChangePasswordSection />)

    await fireEvent.press(screen.getByLabelText('Actualizar contraseña'))

    expect(screen.getByText('La contraseña actual es obligatoria')).toBeOnTheScreen()
    expect(screen.getByText('La nueva contraseña es obligatoria')).toBeOnTheScreen()
    expect(screen.getByText('Confirma la nueva contraseña')).toBeOnTheScreen()
    expect(mockedChangePassword).not.toHaveBeenCalled()
  })

  it('rechaza una nueva contraseña igual a la actual', async () => {
    await render(<ChangePasswordSection />)

    await fillForm('current123', 'current123', 'current123')
    await fireEvent.press(screen.getByLabelText('Actualizar contraseña'))

    expect(screen.getByText('La nueva contraseña debe ser diferente a la actual')).toBeOnTheScreen()
    expect(mockedChangePassword).not.toHaveBeenCalled()
  })

  it('rechaza contraseñas cortas o que no coinciden', async () => {
    await render(<ChangePasswordSection />)

    await fillForm('current123', 'short', 'other')
    await fireEvent.press(screen.getByLabelText('Actualizar contraseña'))

    expect(screen.getByText('Mínimo 8 caracteres')).toBeOnTheScreen()
    expect(screen.getByText('Las contraseñas no coinciden')).toBeOnTheScreen()
    expect(mockedChangePassword).not.toHaveBeenCalled()
  })

  it('cambia la contraseña, limpia el formulario y avisa del éxito', async () => {
    mockedChangePassword.mockResolvedValue(undefined)
    await render(<ChangePasswordSection />)

    await fillForm('current123', 'newpass123', 'newpass123')
    await fireEvent.press(screen.getByLabelText('Actualizar contraseña'))

    await waitFor(() => {
      expect(mockedChangePassword).toHaveBeenCalledWith({
        currentPassword: 'current123',
        newPassword: 'newpass123',
      })
    })
    expect(screen.getByText('Contraseña actualizada correctamente')).toBeOnTheScreen()
    expect(screen.getByTestId('change-password-current')).toHaveProp('value', '')
    expect(screen.getByTestId('change-password-new')).toHaveProp('value', '')
    expect(screen.getByTestId('change-password-confirm')).toHaveProp('value', '')
  })

  it('muestra el detalle del backend cuando la contraseña actual es incorrecta', async () => {
    mockedChangePassword.mockRejectedValue({
      type: 'about:blank',
      title: 'Validation Error',
      status: 422,
      detail: 'La contraseña actual es incorrecta',
      instance: '/api/v1/profile/password',
    })
    await render(<ChangePasswordSection />)

    await fillForm('incorrecta', 'newpass123', 'newpass123')
    await fireEvent.press(screen.getByLabelText('Actualizar contraseña'))

    await waitFor(() => {
      expect(screen.getByText('La contraseña actual es incorrecta')).toBeOnTheScreen()
    })
    // El formulario no se limpia para que el usuario corrija.
    expect(screen.getByTestId('change-password-current')).toHaveProp('value', 'incorrecta')
  })
})
