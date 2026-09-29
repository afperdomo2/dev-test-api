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

  it('requires selecting the correct answer for true_false', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      selectedType: string
      content: string
      topicIds: Array<string>
      trueFalseAnswer: boolean | null
      validate: () => boolean
      validationErrors: Record<string, Array<string>>
    }
    vm.selectedType = 'true_false'
    vm.content = '¿Go permite múltiples valores de retorno?'
    vm.topicIds = ['t1']

    expect(vm.validate()).toBe(false)
    expect(vm.validationErrors.trueFalseAnswer).toBeTruthy()

    vm.trueFalseAnswer = true
    expect(vm.validate()).toBe(true)
  })

  it('buildTrueFalseOptions returns exactly 2 options with one correct', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      trueFalseAnswer: boolean | null
      buildTrueFalseOptions: () => Array<{ content: string; isCorrect: boolean }>
    }
    vm.trueFalseAnswer = true
    expect(vm.buildTrueFalseOptions()).toEqual([
      { content: 'Verdadero', isCorrect: true },
      { content: 'Falso', isCorrect: false },
    ])
    vm.trueFalseAnswer = false
    expect(vm.buildTrueFalseOptions()).toEqual([
      { content: 'Verdadero', isCorrect: false },
      { content: 'Falso', isCorrect: true },
    ])
  })

  it('derives the true_false answer when editing an existing question', async () => {
    const question = {
      id: '1',
      type: 'true_false',
      content: 'Q',
      difficulty: 'easy',
      explanation: '',
      topics: [],
      options: [
        { id: 'o1', content: 'Verdadero', isCorrect: false },
        { id: 'o2', content: 'Falso', isCorrect: true },
      ],
    }
    const wrapper = mountDialog({ modelValue: false, question })
    await wrapper.setProps({ modelValue: true })
    const vm = wrapper.vm as unknown as { trueFalseAnswer: boolean | null }
    expect(vm.trueFalseAnswer).toBe(false)
  })
})
