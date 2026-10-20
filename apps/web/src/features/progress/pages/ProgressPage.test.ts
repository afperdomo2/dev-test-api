import { flushPromises } from '@vue/test-utils'
import type { VueWrapper } from '@vue/test-utils'
import ProgressPage from './ProgressPage.vue'
import { mountWithProviders } from '@/__tests__/utils/mount'

// La fixture va DENTRO del factory: vi.mock se hoistea y no puede capturar
// constantes declaradas después de los imports.
vi.mock('@/queries/progress.queries', () => ({
  upcomingQuestionsOptions: () => ({
    queryKey: ['progress', 'upcoming'],
    queryFn: async () => ({ data: [], meta: { total: 0 } }),
  }),
  savedQuestionsOptions: () => ({
    queryKey: ['progress', 'saved'],
    queryFn: async () => ({ data: [], meta: { total: 0 } }),
  }),
  topicMasteryOptions: () => ({
    queryKey: ['progress', 'topics-mastery'],
    queryFn: async () => [
      { topicId: 't1', topicName: 'Go', attempts: 5, correctCount: 2, accuracy: 0.4, isWeak: true },
    ],
  }),
  toggleSaveMutation: () => ({
    mutationKey: ['progress', 'toggle-save'],
    mutationFn: vi.fn(),
  }),
}))

/** Busca una pestaña por su texto sin recurrir a aserciones no-null. */
function tabByLabel(wrapper: VueWrapper<unknown>, label: string) {
  const tab = wrapper.findAll('[role="tab"]').find((t) => t.text() === label)
  if (!tab) throw new Error(`No se encontró la pestaña "${label}"`)
  return tab
}

describe('ProgressPage — pestaña Temas', () => {
  async function mountPage() {
    const wrapper = await mountWithProviders(ProgressPage, {
      stubs: { GamificationPanel: true, QuestionDetailDialog: true },
    })
    await flushPromises()
    return wrapper
  }

  it('muestra las cuatro pestañas, incluida Temas', async () => {
    const wrapper = await mountPage()

    const labels = wrapper.findAll('[role="tab"]').map((tab) => tab.text())
    expect(labels).toEqual(['Pendientes de repaso', 'Guardadas', 'Gamificación', 'Temas'])
  })

  it('arranca en la pestaña de repaso, no en Temas', async () => {
    const wrapper = await mountPage()

    expect(wrapper.text()).not.toContain('Dominio por tema')
    expect(wrapper.text()).toContain('No hay repasos pendientes')
  })

  it('al pulsar Temas renderiza el dominio por tema', async () => {
    const wrapper = await mountPage()

    await tabByLabel(wrapper, 'Temas').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Dominio por tema')
    expect(wrapper.text()).toContain('Go')
    expect(wrapper.text()).toContain('2/5')
    expect(wrapper.text()).toContain('Débil')
  })

  it('al volver a Guardadas deja de mostrar Temas', async () => {
    const wrapper = await mountPage()

    await tabByLabel(wrapper, 'Temas').trigger('click')
    await flushPromises()
    expect(wrapper.text()).toContain('Dominio por tema')

    await tabByLabel(wrapper, 'Guardadas').trigger('click')
    await flushPromises()

    expect(wrapper.text()).not.toContain('Dominio por tema')
  })
})
