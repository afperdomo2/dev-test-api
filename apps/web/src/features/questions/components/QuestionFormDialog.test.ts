import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import QuestionFormDialog from './QuestionFormDialog.vue'

vi.mock('@/api/services/topics.service', () => ({
  listTopics: vi.fn().mockResolvedValue({ data: [] }),
}))

describe('QuestionFormDialog', () => {
  function mountDialog(props: Record<string, unknown> = {}) {
    setActivePinia(createPinia())
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
    return mount(QuestionFormDialog, {
      props: { modelValue: true, question: null, ...props } as never,
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
        stubs: { teleport: true },
      },
    })
  }

  it('shows nueva pregunta title', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Nueva pregunta')
  })

  it('shows editar pregunta when editing', () => {
    const question = {
      id: '1',
      type: 'single_choice',
      content: 'Q',
      difficulty: 'easy',
      explanation: '',
      topics: [],
      options: [],
    }
    const wrapper = mountDialog({ question })
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Editar pregunta')
  })

  it('validates content required', async () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      validate: () => boolean
      validationErrors: Record<string, unknown>
    }
    const valid = vm.validate()
    expect(valid).toBe(false)
  })

  it('has addOption and removeOption', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      options: Array<{ content: string }>
      addOption: () => void
      removeOption: (i: number) => void
    }
    const before = vm.options.length
    vm.addOption()
    expect(vm.options.length).toBe(before + 1)
    vm.removeOption(0)
    expect(vm.options.length).toBe(before)
  })

  it('emits close', async () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { close: () => void }
    vm.close()
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })
})
