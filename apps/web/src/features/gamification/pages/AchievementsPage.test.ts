import { flushPromises } from '@vue/test-utils'
import AchievementsPage from './AchievementsPage.vue'
import { mountWithProviders } from '@/__tests__/utils/mount'
import { getBadges } from '@/api/services/gamification.service'
import type { BadgesResponse } from '@/types/gamification.types'

vi.mock('@/api/services/gamification.service', () => ({
  getBadges: vi.fn(),
}))

const response: BadgesResponse = {
  earned: 2,
  total: 3,
  badges: [
    {
      id: 'first_session',
      name: 'Primera sesión',
      description: 'Completa tu primera sesión',
      icon: 'mdi-flag-checkered',
      earned: true,
      current: 1,
      target: 1,
    },
    {
      id: 'answers_100',
      name: '100 respuestas',
      description: 'Responde 100 preguntas',
      icon: 'mdi-check-all',
      earned: false,
      current: 12,
      target: 100,
    },
    {
      id: 'streak_7',
      name: 'Racha de 7',
      description: '7 días seguidos de estudio',
      icon: 'mdi-fire',
      earned: true,
      current: 7,
      target: 7,
    },
  ],
}

async function mountPage() {
  const wrapper = await mountWithProviders(AchievementsPage)
  await flushPromises()
  return wrapper
}

describe('AchievementsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getBadges).mockResolvedValue(response)
  })

  it('muestra el contador de logros', async () => {
    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('Logros')
    expect(wrapper.text()).toContain('2 de 3 logros')
  })

  it('distingue conseguidos y pendientes', async () => {
    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('Primera sesión')
    expect(wrapper.text()).toContain('Conseguido')
    expect(wrapper.text()).toContain('100 respuestas')
    expect(wrapper.text()).toContain('12/100')
    // Un pendiente no se pinta como conseguido
    expect(wrapper.text().match(/Conseguido/g)?.length).toBe(2)
  })

  it('si falla la carga permite reintentar', async () => {
    vi.mocked(getBadges).mockRejectedValueOnce(new Error('boom'))

    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('No se pudieron cargar los logros')
    expect(wrapper.text()).toContain('Reintentar')
  })
})
