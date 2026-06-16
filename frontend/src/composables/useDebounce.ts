import { ref, watch } from 'vue'

export function useDebounce<T>(initialValue: T, delay = 500) {
  const value = ref<T>(initialValue) as ReturnType<typeof ref<T>>
  const debouncedValue = ref<T>(initialValue) as ReturnType<typeof ref<T>>

  let timer: ReturnType<typeof setTimeout> | null = null

  watch(value, (newVal) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      debouncedValue.value = newVal
    }, delay)
  })

  return { value, debouncedValue }
}
