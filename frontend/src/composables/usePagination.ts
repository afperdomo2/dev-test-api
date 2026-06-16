import { computed, ref } from 'vue'

export function usePagination(initialPage = 1, initialPerPage = 20) {
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
