import { flushPromises } from '@vue/test-utils'
import InsightsPage from './InsightsPage.vue'
import { mountWithProviders } from '@/__tests__/utils/mount'
import { getInsightSummary } from '@/api/services/insights.service'
import type { InsightSummary } from '@/types/insight.types'

vi.mock('@/api/services/insights.service', () => ({
  getInsightSummary: vi.fn(),
}))

const summary: InsightSummary = {
  period: 'week',
  from: '2026-10-13T00:00:00Z',
  to: '2026-10-19T12:00:00Z',
  sessions: 2,
  questionsAnswered: 5,
  accuracy: 0.6,
  timeSpentMs: 12 * 60 * 1000,
  topTopics: [{ topicId: 't1', topicName: 'Go', attempts: 3, accuracy: 0.6667 }],
  weakTopics: [
    {
      topicId: 't2',
      topicName: 'SQL',
      attempts: 5,
      correctCount: 2,
      accuracy: 0.4,
      isWeak: true,
    },
  ],
  xpGained: 75,
  streakDays: 3,
  dailyActivity: [{ date: '2026-10-19', questionsAnswered: 5 }],
}

async function mountPage() {
  const wrapper = await mountWithProviders(InsightsPage)
  await flushPromises()
  return wrapper
}

describe('InsightsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getInsightSummary).mockResolvedValue(summary)
  })

  it('muestra los KPIs del periodo', async () => {
    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('Informes')
    expect(wrapper.text()).toContain('Sesiones')
    expect(wrapper.text()).toContain('Preguntas')
    expect(wrapper.text()).toContain('60%')
    expect(wrapper.text()).toContain('12 min')
    expect(wrapper.text()).toContain('75')
    expect(wrapper.text()).toContain('3 días')
  })

  it('muestra los temas practicados y el tema débil', async () => {
    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('Más practicados')
    expect(wrapper.text()).toContain('Go')
    expect(wrapper.text()).toContain('Temas débiles')
    expect(wrapper.text()).toContain('SQL')
    expect(wrapper.text()).toContain('Débil')
  })

  it('al cambiar a Mes vuelve a consultar con period=month', async () => {
    const wrapper = await mountPage()

    const monthButton = wrapper.findAll('button').find((button) => button.text() === 'Mes')
    if (!monthButton) throw new Error('No se encontró el botón Mes')

    await monthButton.trigger('click')
    await flushPromises()

    expect(getInsightSummary).toHaveBeenCalledWith('month')
  })

  it('sin actividad muestra el aviso', async () => {
    vi.mocked(getInsightSummary).mockResolvedValue({
      ...summary,
      sessions: 0,
      questionsAnswered: 0,
    })

    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('No hay actividad en este periodo')
  })

  it('oculta temas débiles cuando no hay', async () => {
    vi.mocked(getInsightSummary).mockResolvedValue({ ...summary, weakTopics: [] })

    const wrapper = await mountPage()

    expect(wrapper.text()).not.toContain('Temas débiles')
  })

  it('si falla la carga permite reintentar', async () => {
    vi.mocked(getInsightSummary).mockRejectedValueOnce(new Error('boom'))

    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('No se pudo cargar el informe')
    expect(wrapper.text()).toContain('Reintentar')
  })
})
