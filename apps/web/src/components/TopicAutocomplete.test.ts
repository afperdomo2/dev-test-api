import { flushPromises } from '@vue/test-utils'
import { mountWithProviders } from '@/__tests__/utils/mount'
import { listTopics } from '@/api/services/topics.service'
import type { Topic } from '@/types/topic.types'
import TopicAutocomplete from './TopicAutocomplete.vue'

vi.mock('@/api/services/topics.service', () => ({
  listTopics: vi.fn(),
}))

const vueTopic: Topic = {
  id: 't1',
  slug: 'vue',
  name: 'Vue',
  category: 'frontend',
  isSystem: true,
  createdAt: '2026-01-01T00:00:00Z',
}

function pageResponse(topics: Array<Topic>, total: number) {
  return { data: topics, meta: { total, page: 1, perPage: 20 } }
}

describe('TopicAutocomplete', () => {
  beforeEach(() => vi.clearAllMocks())

  it('loads first page with server search and renders items', async () => {
    vi.mocked(listTopics).mockResolvedValue(pageResponse([vueTopic], 1))
    const wrapper = await mountWithProviders(TopicAutocomplete, {
      props: { modelValue: [] },
    })
    await flushPromises()

    expect(listTopics).toHaveBeenCalledWith(1, 20, 'name', 'asc', '')
    const vm = wrapper.vm as unknown as { items: Array<{ title: string; value: string }> }
    expect(vm.items).toEqual([{ title: 'Vue', value: 't1', props: { subtitle: 'frontend' } }])
  })

  it('emits selected topics when model value changes', async () => {
    vi.mocked(listTopics).mockResolvedValue(pageResponse([vueTopic], 1))
    const wrapper = await mountWithProviders(TopicAutocomplete, {
      props: { modelValue: [] },
    })
    await flushPromises()

    await wrapper.setProps({ modelValue: ['t1'] })
    await flushPromises()

    const emitted = wrapper.emitted('update:selectedTopics')
    expect(emitted).toBeTruthy()
    const last = emitted?.at(-1)?.[0] as Array<Topic>
    expect(last).toEqual([vueTopic])
  })

  it('coalesces empty search to empty string', async () => {
    vi.mocked(listTopics).mockResolvedValue(pageResponse([], 0))
    await mountWithProviders(TopicAutocomplete, {
      props: { modelValue: [] },
    })
    await flushPromises()

    expect(listTopics).toHaveBeenCalledWith(1, 20, 'name', 'asc', '')
  })
})
