/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import ErrorState from './ErrorState.vue'

function createWrapper(error: unknown, retry?: () => void) {
  setActivePinia(createPinia())
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div>home</div>' } },
      { path: '/login', component: { template: '<div>login</div>' } },
    ],
  })
  return mount(ErrorState, {
    props: { error, retry } as never,
    global: {
      plugins: [router],
    },
  })
}

describe('ErrorState', () => {
  it('shows 401 sesión expirada', () => {
    const wrapper = createWrapper({ status: 401, detail: 'expired' })
    expect(wrapper.text()).toContain('Sesión expirada')
    expect(wrapper.text()).toContain('Ir al login')
  })

  it('shows 403 acceso denegado', () => {
    const wrapper = createWrapper({ status: 403 })
    expect(wrapper.text()).toContain('Acceso denegado')
  })

  it('shows 404 no encontrado with message', () => {
    const wrapper = createWrapper({ status: 404, detail: 'No existe' })
    expect(wrapper.text()).toContain('No encontrado')
    expect(wrapper.text()).toContain('No existe')
  })

  it('shows default error with message and retry button when provided', () => {
    const retry = vi.fn()
    const wrapper = createWrapper({ status: 500, detail: 'Fallo' }, retry)
    expect(wrapper.text()).toContain('Error')
    expect(wrapper.text()).toContain('Fallo')
    expect(wrapper.text()).toContain('Reintentar')
  })

  it('shows generic message when error is not object', () => {
    const wrapper = createWrapper('string error')
    expect(wrapper.text()).toContain('Ocurrió un error inesperado')
  })

  it('calls retry when button clicked', async () => {
    const retry = vi.fn()
    const wrapper = createWrapper({ status: 500, detail: 'oops' }, retry)
    const btn = wrapper.findAll('button').find((b) => b.text().includes('Reintentar'))
    expect(btn).toBeDefined()
    await btn!.trigger('click')
    expect(retry).toHaveBeenCalled()
  })

  it('does not show retry button when no retry prop', () => {
    const wrapper = createWrapper({ status: 500, detail: 'oops' })
    expect(wrapper.text()).not.toContain('Reintentar')
  })
})
