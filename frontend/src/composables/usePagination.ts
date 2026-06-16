import { computed, ref } from 'vue'
import { DEFAULT_PER_PAGE } from '@/constants'

export function usePagination(initialPage = 1, initialPerPage = DEFAULT_PER_PAGE) {
  const page = ref(initialPage)
  const perPage = ref(initialPerPage)
  const total = ref(0)

  const totalPages = computed(() => Math.ceil(total.value / perPage.value))

  function reset() {
    page.value = 1
  }

  function setTotal(newTotal: number) {
    total.value = newTotal
  }

  return {
    page,
    perPage,
    total,
    totalPages,
    reset,
    setTotal,
  }
}
