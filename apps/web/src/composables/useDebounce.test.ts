import { nextTick } from 'vue'
import { useDebounce } from './useDebounce'

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('initial values match', () => {
    const { value, debouncedValue } = useDebounce('hello', 500)
    expect(value.value).toBe('hello')
    expect(debouncedValue.value).toBe('hello')
  })

  it('debounces after delay', async () => {
    const { value, debouncedValue } = useDebounce('a', 300)
    value.value = 'b'
    await nextTick()
    expect(debouncedValue.value).toBe('a')
    vi.advanceTimersByTime(299)
    expect(debouncedValue.value).toBe('a')
    vi.advanceTimersByTime(1)
    expect(debouncedValue.value).toBe('b')
  })

  it('resets timer on rapid changes (only last value)', async () => {
    const { value, debouncedValue } = useDebounce('x', 500)
    value.value = 'y'
    await nextTick()
    vi.advanceTimersByTime(200)
    value.value = 'z'
    await nextTick()
    vi.advanceTimersByTime(200)
    expect(debouncedValue.value).toBe('x')
    vi.advanceTimersByTime(300)
    expect(debouncedValue.value).toBe('z')
  })

  it('uses custom delay', async () => {
    const { value, debouncedValue } = useDebounce(0, 100)
    value.value = 1
    await nextTick()
    vi.advanceTimersByTime(99)
    expect(debouncedValue.value).toBe(0)
    vi.advanceTimersByTime(1)
    expect(debouncedValue.value).toBe(1)
  })

  it('default delay is 500', async () => {
    const { value, debouncedValue } = useDebounce('init')
    value.value = 'next'
    await nextTick()
    vi.advanceTimersByTime(499)
    expect(debouncedValue.value).toBe('init')
    vi.advanceTimersByTime(1)
    expect(debouncedValue.value).toBe('next')
  })
})
