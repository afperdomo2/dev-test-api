import { mount } from '@vue/test-utils'
import TopicMasteryList from './TopicMasteryList.vue'
import type { TopicMasteryItem } from '@/types/progress.types'

const weakTopic: TopicMasteryItem = {
  topicId: 't1',
  topicName: 'Go',
  attempts: 5,
  correctCount: 2,
  accuracy: 0.4,
  isWeak: true,
}

const strongTopic: TopicMasteryItem = {
  topicId: 't2',
  topicName: 'Vue',
  attempts: 5,
  correctCount: 4,
  accuracy: 0.8,
  isWeak: false,
}

const items = [weakTopic, strongTopic]

describe('TopicMasteryList', () => {
  it('renderiza una fila por tema con aciertos/respuestas', () => {
    const wrapper = mount(TopicMasteryList, { props: { items } })

    expect(wrapper.text()).toContain('Go')
    expect(wrapper.text()).toContain('Vue')
    expect(wrapper.text()).toContain('2/5')
    expect(wrapper.text()).toContain('4/5')
    // Solo las barras con etiqueta accesible: `v-card` pinta además un loader
    // interno que también expone role="progressbar" (sin aria-label).
    expect(wrapper.findAll('[role="progressbar"][aria-label]')).toHaveLength(2)
  })

  it('muestra el chip "Débil" solo en los temas débiles', () => {
    const wrapper = mount(TopicMasteryList, { props: { items } })

    const chips = wrapper.findAll('.v-chip')
    expect(chips).toHaveLength(1)
    expect(chips[0]?.text()).toBe('Débil')
  })

  it('con un solo tema fuerte no muestra ningún chip', () => {
    const wrapper = mount(TopicMasteryList, { props: { items: [strongTopic] } })

    expect(wrapper.findAll('.v-chip')).toHaveLength(0)
  })

  it('muestra el estado vacío cuando no hay temas', () => {
    const wrapper = mount(TopicMasteryList, { props: { items: [] } })

    expect(wrapper.text()).toContain('Aún no hay respuestas registradas')
    expect(wrapper.findAll('[role="progressbar"][aria-label]')).toHaveLength(0)
  })

  it('mientras carga muestra el skeleton y no el estado vacío', () => {
    const wrapper = mount(TopicMasteryList, { props: { items: [], loading: true } })

    expect(wrapper.find('.v-skeleton-loader').exists()).toBe(true)
    expect(wrapper.text()).not.toContain('Aún no hay respuestas registradas')
  })

  it('expone cada barra con etiqueta accesible y porcentaje', () => {
    const wrapper = mount(TopicMasteryList, { props: { items } })

    const bars = wrapper.findAll('[role="progressbar"][aria-label]')
    expect(bars).toHaveLength(2)
    expect(bars[0]?.attributes('aria-label')).toBe('Go: 40% de acierto')
    expect(bars[0]?.attributes('aria-valuenow')).toBe('40')
    expect(bars[1]?.attributes('aria-label')).toBe('Vue: 80% de acierto')
    expect(bars[1]?.attributes('aria-valuenow')).toBe('80')
  })
})
