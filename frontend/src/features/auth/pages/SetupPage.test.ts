import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import SetupPage from './SetupPage.vue'

vi.mock('@/api/services/auth.service', () => ({
  setup: vi.fn(),
}))

describe('SetupPage', () => {
  function mountPage() {
    setActivePinia(createPinia())
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/login', component: { template: '<div>login</div>' } }],
    })
    return mount(SetupPage, {
      global: { plugins: [router] },
    })
  }

  it('renders configuración inicial', () => {
    const wrapper = mountPage()
    expect(wrapper.text()).toContain('Configuración inicial')
    expect(wrapper.text()).toContain('Crea la cuenta de administrador')
  })

  it('has email, password, confirm fields', () => {
    const wrapper = mountPage()
    expect(wrapper.text()).toContain('Email')
    expect(wrapper.text()).toContain('Contraseña')
    expect(wrapper.text()).toContain('Confirmar contraseña')
  })

  it('validates empty submit', async () => {
    const wrapper = mountPage()
    const form = wrapper.find('form')
    await form.trigger('submit.prevent')
    const { setup } = await import('@/api/services/auth.service')
    expect(vi.mocked(setup)).not.toHaveBeenCalled()
  })

  it('calls setup on valid data', async () => {
    const wrapper = mountPage()
    const vm = wrapper.vm as unknown as {
      form: { email: string; password: string; confirmPassword: string }
      submit: () => Promise<void>
    }
    vm.form.email = 'admin@example.com'
    vm.form.password = 'secret123'
    vm.form.confirmPassword = 'secret123'
    await vm.submit()
    const { setup } = await import('@/api/services/auth.service')
    expect(vi.mocked(setup)).toHaveBeenCalledWith({ email: 'admin@example.com', password: 'secret123' })
  })
})
