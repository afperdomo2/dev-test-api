import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import LoginPage from './LoginPage.vue'

vi.mock('@/api/services/auth.service', () => ({
  login: vi.fn(),
}))

describe('LoginPage', () => {
  it('renders card with title', () => {
    setActivePinia(createPinia())
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/', component: { template: '<div>home</div>' } }],
    })
    const wrapper = mount(LoginPage, {
      global: { plugins: [router] },
    })
    expect(wrapper.text()).toContain('Iniciar sesión')
  })

  it('contains LoginForm', () => {
    setActivePinia(createPinia())
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/', component: { template: '<div>home</div>' } }],
    })
    const wrapper = mount(LoginPage, {
      global: { plugins: [router] },
    })
    expect(wrapper.html()).toContain('v-text-field')
  })
})
