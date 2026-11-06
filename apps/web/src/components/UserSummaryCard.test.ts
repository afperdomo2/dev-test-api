import { mount } from '@vue/test-utils'
import UserSummaryCard from './UserSummaryCard.vue'

const baseUser = {
  id: '3f1c2a9e-1234-4abc-9def-567890abcdef',
  email: 'a@b.com',
  isAdmin: false,
  dailyImportLimit: 200,
  createdAt: '2026-08-28T12:00:00Z',
  updatedAt: '2026-08-28T12:00:00Z',
}

describe('UserSummaryCard', () => {
  it('renders email, member since, role and id', () => {
    const wrapper = mount(UserSummaryCard, { props: { user: baseUser } })

    expect(wrapper.text()).toContain('a@b.com')
    expect(wrapper.text()).toContain('Miembro desde')
    expect(wrapper.text()).toContain('Usuario')
    expect(wrapper.text()).toContain(baseUser.id)
  })

  it('shows the admin role for admins', () => {
    const wrapper = mount(UserSummaryCard, { props: { user: { ...baseUser, isAdmin: true } } })

    expect(wrapper.text()).toContain('Administrador')
  })

  it('renders safely without a user', () => {
    const wrapper = mount(UserSummaryCard, { props: { user: null } })

    expect(wrapper.text()).toContain('Usuario')
    expect(wrapper.text()).not.toContain('Administrador')
  })
})
