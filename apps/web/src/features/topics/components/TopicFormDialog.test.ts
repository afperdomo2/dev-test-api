import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import TopicFormDialog from './TopicFormDialog.vue'

describe('TopicFormDialog', () => {
  function mountDialog(props: Record<string, unknown> = {}) {
    setActivePinia(createPinia())
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
    return mount(TopicFormDialog, {
      props: { modelValue: true, topic: null, ...props } as never,
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
        stubs: { teleport: true },
      },
    })
  }

  it('shows nuevo tema', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Nuevo tema')
  })

  it('shows editar tema when topic provided', () => {
    const topic = { id: '1', slug: 'js', name: 'JavaScript', category: 'programming' }
    const wrapper = mountDialog({ topic })
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Editar tema')
  })

  it('slugifies name on input when creating', async () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      form: { name: string; slug: string }
      onNameInput: () => void
    }
    vm.form.name = 'My Topic Name'
    vm.onNameInput()
    expect(vm.form.slug).toBe('my-topic-name')
  })

  it('does not slugify in edit mode', () => {
    const topic = { id: '1', slug: 'existing', name: 'Old', category: 'cat' }
    const wrapper = mountDialog({ topic })
    const vm = wrapper.vm as unknown as {
      form: { slug: string; name: string }
      onNameInput: () => void
      isEdit: boolean
    }
    // isEdit should be true
    expect(vm.isEdit).toBe(true)
    vm.form.name = 'New Name'
    vm.onNameInput()
    // should remain existing slug (but form initially empty due to watch not triggered, so slug is still '' - we set it first)
    vm.form.slug = 'existing'
    vm.form.name = 'New Name 2'
    vm.onNameInput()
    expect(vm.form.slug).toBe('existing')
  })

  it('validate fails when empty', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { validate: () => boolean }
    expect(vm.validate()).toBe(false)
  })

  it('validate passes when filled', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as {
      form: { slug: string; name: string; category: string }
      validate: () => boolean
    }
    vm.form.slug = 'slug'
    vm.form.name = 'name'
    vm.form.category = 'cat'
    expect(vm.validate()).toBe(true)
  })

  it('emits close', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { close: () => void }
    vm.close()
    expect(wrapper.emitted('update:modelValue')).toEqual([[false]])
  })

  it('onSlugInput sets manuallyEdited', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { slugManuallyEdited: boolean; onSlugInput: () => void }
    expect(vm.slugManuallyEdited).toBe(false)
    vm.onSlugInput()
    expect(vm.slugManuallyEdited).toBe(true)
  })
})
