import { mountWithProviders } from '@/__tests__/utils/mount'
import SessionRow from './SessionRow.vue'
import type { Session } from '@/types/session.types'

function buildSession(overrides: Partial<Session> = {}): Session {
  return {
    id: 's1',
    userId: 'u1',
    name: 'Repaso de Go',
    status: 'in_progress',
    mode: 'review',
    difficulty: 'intermediate',
    questionLimit: 20,
    score: 0,
    startedAt: new Date().toISOString(),
    finishedAt: null,
    topics: ['Go', 'Concurrencia', 'Testing'],
    answerCount: 12,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    ...overrides,
  }
}

describe('SessionRow', () => {
  it('renders name, status and progress for an in-progress session', async () => {
    const wrapper = await mountWithProviders(SessionRow, {
      props: { session: buildSession() },
      routerRoutes: [{ path: '/sessions/:id/study', component: { template: '<div />' } }],
    })

    expect(wrapper.find('.session-row__name').text()).toBe('Repaso de Go')
    expect(wrapper.text()).toContain('En progreso')
    expect(wrapper.text()).toContain('12/20')
    expect(wrapper.text()).toContain('Continuar')
  })

  it('shows the score and review CTA for a completed session', async () => {
    const wrapper = await mountWithProviders(SessionRow, {
      props: {
        session: buildSession({
          status: 'completed',
          score: 87.5,
          answerCount: 20,
          finishedAt: new Date(Date.now() + 34 * 60 * 1000).toISOString(),
        }),
      },
      routerRoutes: [{ path: '/sessions/:id/review', component: { template: '<div />' } }],
    })

    expect(wrapper.text()).toContain('87.5%')
    expect(wrapper.text()).toContain('20 respuestas')
    expect(wrapper.text()).toContain('34 min')
    expect(wrapper.text()).toContain('Ver resumen')
  })

  it('emits open when the row is clicked', async () => {
    const wrapper = await mountWithProviders(SessionRow, {
      props: { session: buildSession() },
      routerRoutes: [{ path: '/sessions/:id/study', component: { template: '<div />' } }],
    })

    await wrapper.find('.session-row').trigger('click')
    expect(wrapper.emitted('open')?.[0]?.[0]).toMatchObject({ id: 's1' })
  })

  it('only offers the actions menu when an action is available', async () => {
    const withActions = await mountWithProviders(SessionRow, {
      props: { session: buildSession() },
      routerRoutes: [{ path: '/sessions/:id/study', component: { template: '<div />' } }],
    })
    expect(withActions.findComponent({ name: 'VMenu' }).exists()).toBe(true)

    const withoutActions = await mountWithProviders(SessionRow, {
      props: { session: buildSession({ status: 'completed', answerCount: 5 }) },
      routerRoutes: [{ path: '/sessions/:id/review', component: { template: '<div />' } }],
    })
    expect(withoutActions.findComponent({ name: 'VMenu' }).exists()).toBe(false)
  })
})
