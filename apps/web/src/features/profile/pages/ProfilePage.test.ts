import { createPinia, setActivePinia } from 'pinia'
import { flushPromises } from '@vue/test-utils'
import ProfilePage from './ProfilePage.vue'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { mountWithProviders } from '@/__tests__/utils/mount'

vi.mock('@/api/services/users.service', () => ({
  changePassword: vi.fn(),
}))

import * as usersService from '@/api/services/users.service'

async function mountProfile() {
  const pinia = createPinia()
  setActivePinia(pinia)
  const auth = useAuthStore()
  auth.user = {
    id: '3f1c2a9e-1234-4abc-9def-567890abcdef',
    email: 'a@b.com',
    isAdmin: false,
    dailyImportLimit: 200,
    createdAt: '2026-08-28T12:00:00Z',
    updatedAt: '2026-08-28T12:00:00Z',
  }
  auth.token = 'tok'

  const wrapper = await mountWithProviders(ProfilePage, { pinia })
  return { wrapper, auth }
}

function passwordInputs(wrapper: Awaited<ReturnType<typeof mountProfile>>['wrapper']) {
  return wrapper.findAll('input[type="password"]')
}

async function fillForm(
  wrapper: Awaited<ReturnType<typeof mountProfile>>['wrapper'],
  current: string,
  next: string,
  confirm: string,
) {
  const [currentInput, newInput, confirmInput] = passwordInputs(wrapper)
  if (!currentInput || !newInput || !confirmInput) {
    throw new Error('Password inputs not found')
  }
  await currentInput.setValue(current)
  await newInput.setValue(next)
  await confirmInput.setValue(confirm)
}

function submitButton(wrapper: Awaited<ReturnType<typeof mountProfile>>['wrapper']) {
  const button = wrapper.findAll('button').find((b) => b.text().includes('Actualizar contraseña'))
  if (!button) throw new Error('Submit button not found')
  return button
}

describe('ProfilePage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders the user summary and the password form', async () => {
    const { wrapper } = await mountProfile()

    expect(wrapper.text()).toContain('Mi perfil')
    expect(wrapper.text()).toContain('a@b.com')
    expect(wrapper.text()).toContain('Usuario')
    expect(wrapper.text()).toContain('Cambiar contraseña')
    expect(wrapper.text()).toContain('Contraseña actual')
    expect(wrapper.text()).toContain('Nueva contraseña')
    expect(wrapper.text()).toContain('Confirmar nueva contraseña')
  })

  it('validates required fields without calling the backend', async () => {
    const { wrapper } = await mountProfile()

    await submitButton(wrapper).trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Ingresa tu contraseña actual')
    expect(wrapper.text()).toContain('Ingresa la nueva contraseña')
    expect(wrapper.text()).toContain('Confirma la nueva contraseña')
    expect(usersService.changePassword).not.toHaveBeenCalled()
  })

  it('rejects a new password equal to the current one', async () => {
    const { wrapper } = await mountProfile()

    await fillForm(wrapper, 'current123', 'current123', 'current123')
    await submitButton(wrapper).trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('La nueva contraseña debe ser diferente a la actual')
    expect(usersService.changePassword).not.toHaveBeenCalled()
  })

  it('rejects a short new password and mismatched confirmation', async () => {
    const { wrapper } = await mountProfile()

    await fillForm(wrapper, 'current123', 'short', 'other')
    await submitButton(wrapper).trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Mínimo 8 caracteres')
    expect(wrapper.text()).toContain('Las contraseñas no coinciden')
    expect(usersService.changePassword).not.toHaveBeenCalled()
  })

  it('changes the password and keeps the session', async () => {
    vi.mocked(usersService.changePassword).mockResolvedValue(undefined)
    const { wrapper, auth } = await mountProfile()
    const appStore = useAppStore()

    await fillForm(wrapper, 'current123', 'newpass123', 'newpass123')
    await submitButton(wrapper).trigger('click')
    await flushPromises()

    expect(usersService.changePassword).toHaveBeenCalledWith(
      {
        currentPassword: 'current123',
        newPassword: 'newpass123',
      },
      expect.anything(),
    )
    expect(appStore.snackbar.message).toBe('Contraseña actualizada correctamente')
    expect(auth.token).toBe('tok')
    expect(
      passwordInputs(wrapper).every((input) => (input.element as HTMLInputElement).value === ''),
    ).toBe(true)
  })

  it('shows the backend detail when the current password is wrong', async () => {
    vi.mocked(usersService.changePassword).mockRejectedValue({
      type: 'about:blank',
      title: 'Validation Error',
      status: 422,
      detail: 'La contraseña actual es incorrecta',
      instance: '/api/v1/profile/password',
    })
    const { wrapper, auth } = await mountProfile()

    await fillForm(wrapper, 'incorrecta', 'newpass123', 'newpass123')
    await submitButton(wrapper).trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('La contraseña actual es incorrecta')
    expect(auth.token).toBe('tok')
  })
})
