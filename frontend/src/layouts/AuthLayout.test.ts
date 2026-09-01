import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import AuthLayout from './AuthLayout.vue'

describe('AuthLayout', () => {
  it('renders router-view container', () => {
    setActivePinia(createPinia())
    const wrapper = mount(AuthLayout, {
      global: {
        stubs: {
          RouterView: { template: '<div>child</div>' },
          'router-view': { template: '<div>child</div>' },
        },
      },
    })
    expect(wrapper.html()).toContain('child')
  })

  it('shows snackbar bound to appStore', async () => {
    setActivePinia(createPinia())
    const wrapper = mount(AuthLayout)
    const vm = wrapper.vm as unknown as { appStore: { snackbar: { show: boolean; message: string }; showSnackbar: (m: string) => void } }
    vm.appStore.showSnackbar('hello')
    await wrapper.vm.$nextTick()
    expect(vm.appStore.snackbar.message).toBe('hello')
    expect(vm.appStore.snackbar.show).toBe(true)
  })
})
