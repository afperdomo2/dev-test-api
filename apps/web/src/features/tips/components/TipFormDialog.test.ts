import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import TipFormDialog from './TipFormDialog.vue'

describe('TipFormDialog', () => {
  function mountDialog(props: Record<string, unknown> = {}) {
    setActivePinia(createPinia())
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
    return mount(TipFormDialog, {
      props: { modelValue: true, tip: null, ...props } as never,
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
        stubs: { teleport: true },
      },
    })
  }

  it('shows nuevo tip', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Nuevo tip')
  })

  it('shows editar tip when tip provided', () => {
    const tip = { id: 't1', topicId: 'topic-1', content: 'un tip' }
    const wrapper = mountDialog({ tip })
    const vm = wrapper.vm as unknown as { dialogTitle: string; isEdit: boolean }
    expect(vm.dialogTitle).toBe('Editar tip')
    expect(vm.isEdit).toBe(true)
  })

  it('validate fails when empty', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { validate: () => boolean }
    expect(vm.validate()).toBe(false)
  })

  it('validate fails without topic', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      form: { topicId: string; content: string }
      validate: () => boolean
    }
    vm.form.content = 'contenido'
    expect(vm.validate()).toBe(false)
  })

  it('validate passes when topic and content filled', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      form: { topicId: string; content: string }
      validate: () => boolean
    }
    vm.form.topicId = 'topic-1'
    vm.form.content = 'contenido'
    expect(vm.validate()).toBe(true)
  })

  it('validate fails when content exceeds max length', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      form: { topicId: string; content: string }
      validate: () => boolean
    }
    vm.form.topicId = 'topic-1'
    vm.form.content = 'a'.repeat(501)
    expect(vm.validate()).toBe(false)
  })

  it('topicIdModel coerces null to empty string', async () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      form: { topicId: string }
      topicIdModel: string | null
    }
    vm.topicIdModel = 'topic-9'
    await wrapper.vm.$nextTick()
    expect(vm.form.topicId).toBe('topic-9')
    vm.topicIdModel = null
    await wrapper.vm.$nextTick()
    expect(vm.form.topicId).toBe('')
  })

  it('emits close', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { close: () => void }
    vm.close()
    expect(wrapper.emitted('update:modelValue')).toEqual([[false]])
  })
})
