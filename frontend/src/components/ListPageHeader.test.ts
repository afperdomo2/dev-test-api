/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { mount } from '@vue/test-utils'
import ListPageHeader from './ListPageHeader.vue'

describe('ListPageHeader', () => {
  it('renders title and create label', () => {
    const wrapper = mount(ListPageHeader, {
      props: { title: 'Usuarios', createLabel: 'Nuevo usuario' },
    })
    expect(wrapper.text()).toContain('Usuarios')
    expect(wrapper.text()).toContain('Nuevo usuario')
  })

  it('does not show create button when showCreate false', () => {
    const wrapper = mount(ListPageHeader, {
      props: { title: 'Usuarios', createLabel: 'Nuevo', showCreate: false },
    })
    expect(wrapper.text()).not.toContain('Nuevo')
    expect(wrapper.text()).toContain('Refrescar')
  })

  it('emits refresh on click', async () => {
    const wrapper = mount(ListPageHeader, {
      props: { title: 'Test', createLabel: 'Crear' },
    })
    const refreshBtn = wrapper.findAll('button').find((b) => b.text().includes('Refrescar'))
    await refreshBtn!.trigger('click')
    expect(wrapper.emitted('refresh')).toHaveLength(1)
  })

  it('emits create when no createTo', async () => {
    const wrapper = mount(ListPageHeader, {
      props: { title: 'Test', createLabel: 'Crear' },
    })
    const createBtn = wrapper.findAll('button').find((b) => b.text().includes('Crear'))
    await createBtn!.trigger('click')
    expect(wrapper.emitted('create')).toHaveLength(1)
  })

  it('respects refresh cooldown', async () => {
    vi.useFakeTimers()
    const wrapper = mount(ListPageHeader, {
      props: { title: 'Test', createLabel: 'Crear' },
    })
    const refreshBtn = wrapper.findAll('button').find((b) => b.text().includes('Refrescar'))!
    await refreshBtn.trigger('click')
    expect(wrapper.emitted('refresh')).toHaveLength(1)
    // second click during cooldown should not emit
    await refreshBtn.trigger('click')
    expect(wrapper.emitted('refresh')).toHaveLength(1)

    vi.advanceTimersByTime(1000)
    await vi.runOnlyPendingTimersAsync()
    await refreshBtn.trigger('click')
    expect(wrapper.emitted('refresh')).toHaveLength(2)
    vi.useRealTimers()
  })
})
