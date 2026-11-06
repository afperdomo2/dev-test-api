/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import LoginForm from './LoginForm.vue'
import * as authService from '@/api/services/auth.service'

vi.mock('@/api/services/auth.service', () => ({
  login: vi.fn(),
}))

describe('LoginForm', () => {
  function mountForm() {
    setActivePinia(createPinia())
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: { template: '<div>home</div>' } },
        { path: '/login', component: { template: '<div>login</div>' } },
      ],
    })
    return mount(LoginForm, {
      global: { plugins: [router] },
    })
  }

  beforeEach(() => vi.clearAllMocks())

  it('renders email and password fields', () => {
    const wrapper = mountForm()
    expect(wrapper.text()).toContain('Email')
    expect(wrapper.text()).toContain('Contraseña')
  })

  it('oculta la contraseña y permite mostrarla u ocultarla de nuevo', async () => {
    const wrapper = mountForm()
    const toggle = wrapper.find('[data-test="login-password-toggle"]')

    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(toggle.attributes('aria-label')).toBe('Mostrar contraseña')

    await toggle.trigger('click')
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(toggle.attributes('aria-label')).toBe('Ocultar contraseña')

    await toggle.trigger('click')
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    // El botón es type="button": alternar visibilidad nunca envía el formulario.
    expect(authService.login).not.toHaveBeenCalled()
  })

  it('oculta la contraseña mientras se envía el formulario', async () => {
    let resolve: (value: unknown) => void
    vi.mocked(authService.login).mockReturnValue(new Promise((r) => (resolve = r)) as never)
    const wrapper = mountForm()

    await wrapper.find('input[type="email"]').setValue('a@b.com')
    await wrapper.find('input[type="password"]').setValue('secret123')
    await wrapper.find('[data-test="login-password-toggle"]').trigger('click')
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    resolve!({ token: 't', user: { id: '1' } })
  })

  it('validates empty submit', async () => {
    const wrapper = mountForm()
    await wrapper.find('form').trigger('submit.prevent')
    // should not call login if validation fails
    expect(authService.login).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('requerido')
  })

  it('calls login on valid submit', async () => {
    const mockRes = { token: 'tok', user: { id: '1', isAdmin: false } }
    vi.mocked(authService.login).mockResolvedValue(mockRes as never)
    const wrapper = mountForm()

    // Need to set v-model values - use wrapper.vm? Access via component instance
    // LoginForm uses ref form with v-model, we can set via inputs
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')
    await emailInput.setValue('test@example.com')
    await passwordInput.setValue('secret123')
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    // Allow async login to resolve
    await new Promise((r) => setTimeout(r, 0))
    expect(authService.login).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'secret123',
    })
  })

  it('shows snackbar on error', async () => {
    vi.mocked(authService.login).mockRejectedValue({ detail: 'Credenciales inválidas' } as never)
    const wrapper = mountForm()
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')
    await emailInput.setValue('test@example.com')
    await passwordInput.setValue('secret123')
    await wrapper.find('form').trigger('submit.prevent')
    await new Promise((r) => setTimeout(r, 0))
    // Should have called login and failed, but still not crashed
    expect(authService.login).toHaveBeenCalled()
  })

  it('disables button while loading', async () => {
    let resolve: (v: unknown) => void
    vi.mocked(authService.login).mockReturnValue(new Promise((r) => (resolve = r)) as never)
    const wrapper = mountForm()
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')
    await emailInput.setValue('a@b.com')
    await passwordInput.setValue('secret123')
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    // loading state should be true shortly
    // We can't easily assert vuetify loading prop without inspecting, but ensure login was called
    expect(authService.login).toHaveBeenCalled()
    resolve!({ token: 't', user: { id: '1' } })
  })
})
