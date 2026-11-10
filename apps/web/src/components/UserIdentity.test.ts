import { mount } from '@vue/test-utils'
import UserIdentity from './UserIdentity.vue'
import type { User } from '@/types/user.types'

function makeUser(overrides: Partial<User> = {}): User {
  return {
    id: 'user-1',
    email: 'juan@empresa.com',
    isAdmin: false,
    dailyImportLimit: 200,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    ...overrides,
  }
}

describe('UserIdentity', () => {
  describe('variant icon', () => {
    it('muestra solo el avatar con la inicial', () => {
      const wrapper = mount(UserIdentity, {
        props: { user: makeUser(), variant: 'icon' },
      })

      expect(wrapper.find('.v-avatar').text()).toBe('J')
      expect(wrapper.text()).not.toContain('juan@empresa.com')
      expect(wrapper.text()).not.toContain('Usuario')
    })

    it('no muestra el rol de administrador', () => {
      const wrapper = mount(UserIdentity, {
        props: { user: makeUser({ isAdmin: true }), variant: 'icon' },
      })

      expect(wrapper.text()).not.toContain('Administrador')
    })
  })

  describe('variant sidebar', () => {
    it('muestra el email y el rol', () => {
      const wrapper = mount(UserIdentity, {
        props: { user: makeUser({ isAdmin: true }) },
      })

      expect(wrapper.text()).toContain('juan@empresa.com')
      expect(wrapper.text()).toContain('Administrador')
    })
  })

  describe('sin usuario', () => {
    it('no muestra rol y deja el avatar con ?', () => {
      const wrapper = mount(UserIdentity, {
        props: { user: null, variant: 'icon' },
      })

      expect(wrapper.find('.v-avatar').text()).toBe('?')
      expect(wrapper.text()).not.toContain('Usuario')
      expect(wrapper.text()).not.toContain('Administrador')
    })
  })
})
