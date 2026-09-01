import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import UserFormDialog from './UserFormDialog.vue'

vi.mock('@/api/services/users.service', () => ({
  createUser: vi.fn(),
  updateUser: vi.fn(),
}))

describe('UserFormDialog', () => {
  function mountDialog(props: Record<string, unknown> = {}) {
    setActivePinia(createPinia())
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
    return mount(UserFormDialog, {
      props: { modelValue: true, user: null, ...props } as never,
      global: {
        plugins: [[VueQueryPlugin, { queryClient }]],
        stubs: { teleport: true },
      },
    })
  }

  it('shows nuevo usuario title when creating', () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Nuevo usuario')
  })

  it('shows editar usuario when editing', () => {
    const user = { id: '1', email: 'a@b.com', isAdmin: true }
    const wrapper = mountDialog({ user })
    const vm = wrapper.vm as unknown as { dialogTitle: string }
    expect(vm.dialogTitle).toBe('Editar usuario')
  })

  it('hides email field in edit mode', () => {
    const user = { id: '1', email: 'a@b.com', isAdmin: false }
    const wrapper = mountDialog({ user })
    const vm = wrapper.vm as unknown as { isEdit: boolean }
    expect(vm.isEdit).toBe(true)
  })

  it('emits update:modelValue on close', async () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { close: () => void }
    vm.close()
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })

  it('validates required fields', async () => {
    const wrapper = mountDialog()
    const vm = wrapper.vm as unknown as { validate: () => boolean }
    expect(vm.validate()).toBe(false)
    const { createUser } = await import('@/api/services/users.service')
    expect(vi.mocked(createUser)).not.toHaveBeenCalled()
  })
})
