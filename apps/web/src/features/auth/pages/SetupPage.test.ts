import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import SetupPage from './SetupPage.vue'
import { getStatus, setup } from '@/api/services/auth.service'

vi.mock('@/api/services/auth.service', () => ({
  setup: vi.fn(),
  getStatus: vi.fn(),
}))

describe('SetupPage', () => {
  async function mountPage(
    status: { initialized: boolean; setupRequiresToken: boolean } = {
      initialized: false,
      setupRequiresToken: true,
    },
  ) {
    vi.mocked(getStatus).mockResolvedValue(status)
    setActivePinia(createPinia())
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/login', component: { template: '<div>login</div>' } }],
    })
    const wrapper = mount(SetupPage, {
      global: { plugins: [router] },
    })
    await flushPromises()
    return wrapper
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders configuración inicial', async () => {
    const wrapper = await mountPage()
    expect(wrapper.text()).toContain('Configuración inicial')
    expect(wrapper.text()).toContain('Crea la cuenta de administrador')
  })

  it('has email, password, confirm fields', async () => {
    const wrapper = await mountPage()
    expect(wrapper.text()).toContain('Email')
    expect(wrapper.text()).toContain('Contraseña')
    expect(wrapper.text()).toContain('Confirmar contraseña')
  })

  it('shows the setup token field when the backend requires it', async () => {
    const wrapper = await mountPage()
    expect(wrapper.find('[data-test="setup-token"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="setup-disabled"]').exists()).toBe(false)
  })

  it('shows a disabled notice when the backend has no token configured', async () => {
    const wrapper = await mountPage({ initialized: false, setupRequiresToken: false })
    expect(wrapper.find('[data-test="setup-token"]').exists()).toBe(false)
    expect(wrapper.find('[data-test="setup-disabled"]').exists()).toBe(true)
  })

  it('validates empty submit', async () => {
    const wrapper = await mountPage()
    const form = wrapper.find('form')
    await form.trigger('submit.prevent')
    expect(vi.mocked(setup)).not.toHaveBeenCalled()
  })

  it('calls setup with the setup token on valid data', async () => {
    const wrapper = await mountPage()
    const vm = wrapper.vm as unknown as {
      form: {
        email: string
        password: string
        confirmPassword: string
        setupToken: string
      }
      submit: () => Promise<void>
    }
    vm.form.email = 'admin@example.com'
    vm.form.password = 'secret123'
    vm.form.confirmPassword = 'secret123'
    vm.form.setupToken = 'install-token'
    await vm.submit()
    expect(vi.mocked(setup)).toHaveBeenCalledWith({
      email: 'admin@example.com',
      password: 'secret123',
      setupToken: 'install-token',
    })
  })

  it('does not submit when setup is disabled', async () => {
    const wrapper = await mountPage({ initialized: false, setupRequiresToken: false })
    const vm = wrapper.vm as unknown as {
      form: { email: string; password: string; confirmPassword: string }
      submit: () => Promise<void>
    }
    vm.form.email = 'admin@example.com'
    vm.form.password = 'secret123'
    vm.form.confirmPassword = 'secret123'
    await vm.submit()
    expect(vi.mocked(setup)).not.toHaveBeenCalled()
  })
})
