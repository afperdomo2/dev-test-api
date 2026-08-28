import { usePagination } from './usePagination'
import { DEFAULT_PER_PAGE } from '@/constants'

describe('usePagination', () => {
  it('defaults', () => {
    const { page, perPage, total, totalPages } = usePagination()
    expect(page.value).toBe(1)
    expect(perPage.value).toBe(DEFAULT_PER_PAGE)
    expect(total.value).toBe(0)
    expect(totalPages.value).toBe(0)
  })

  it('custom initial values', () => {
    const { page, perPage } = usePagination(3, 20)
    expect(page.value).toBe(3)
    expect(perPage.value).toBe(20)
  })

  it('totalPages computes correctly', () => {
    const { totalPages, setTotal, perPage } = usePagination(1, 10)
    setTotal(25)
    expect(totalPages.value).toBe(3)
    setTotal(0)
    expect(totalPages.value).toBe(0)
    setTotal(10)
    expect(totalPages.value).toBe(1)
    perPage.value = 5
    expect(totalPages.value).toBe(2)
  })

  it('reset sets page to 1', () => {
    const { page, reset } = usePagination(5)
    page.value = 7
    reset()
    expect(page.value).toBe(1)
  })

  it('setTotal updates total', () => {
    const { total, setTotal } = usePagination()
    setTotal(42)
    expect(total.value).toBe(42)
  })
})
