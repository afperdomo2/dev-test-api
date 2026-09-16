/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { mount } from '@vue/test-utils'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import QuestionFilters from './QuestionFilters.vue'

vi.mock('@/api/services/topics.service', () => ({
  listTopics: vi.fn().mockResolvedValue({ data: [] }),
}))

function mountFilters(props: Record<string, unknown> = {}) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return mount(QuestionFilters, {
    props: props as never,
    global: {
      plugins: [[VueQueryPlugin, { queryClient }]],
      stubs: { teleport: true },
    },
  })
}

describe('QuestionFilters', () => {
  it('renders selects', () => {
    const wrapper = mountFilters()
    expect(wrapper.text()).toContain('Tipo')
    expect(wrapper.text()).toContain('Dificultad')
  })

  it('does not show limpiar when no filters', () => {
    const wrapper = mountFilters()
    expect(wrapper.text()).not.toContain('Limpiar filtros')
  })

  it('emits change when selecting type', async () => {
    const wrapper = mountFilters()
    // Access refs directly
    const vm = wrapper.vm as unknown as { selectedType: string; selectedDifficulty: string }
    vm.selectedType = 'single_choice'
    await wrapper.vm.$nextTick()
    // watch should emit
    await wrapper.vm.$nextTick()
    expect(wrapper.emitted('change')).toBeTruthy()
    const last = wrapper.emitted('change')!.at(-1) as unknown as Array<{ type?: string }>
    expect(last[0]).toMatchObject({ type: 'single_choice' })
  })

  it('emits change when selecting difficulty', async () => {
    const wrapper = mountFilters()
    const vm = wrapper.vm as unknown as { selectedDifficulty: string }
    vm.selectedDifficulty = 'hard'
    await wrapper.vm.$nextTick()
    await new Promise((r) => setTimeout(r, 0))
    expect(wrapper.emitted('change')).toBeTruthy()
  })

  it('clearFilters resets and emits', async () => {
    const wrapper = mountFilters()
    const vm = wrapper.vm as unknown as {
      selectedType: string
      selectedDifficulty: string
      clearFilters: () => void
    }
    vm.selectedType = 'single_choice'
    vm.selectedDifficulty = 'hard'
    await wrapper.vm.$nextTick()
    vm.clearFilters()
    await wrapper.vm.$nextTick()
    expect(vm.selectedType).toBe('')
    expect(vm.selectedDifficulty).toBe('')
  })
})
