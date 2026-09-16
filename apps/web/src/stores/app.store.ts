import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const theme = ref<'light' | 'dark'>('light')
  const sidebarOpen = ref(true)
  const snackbar = ref({
    show: false,
    message: '',
    color: 'success' as 'success' | 'error' | 'warning' | 'info',
  })

  function toggleTheme() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
  }

  function showSnackbar(
    message: string,
    color: 'success' | 'error' | 'warning' | 'info' = 'success',
  ) {
    snackbar.value = { show: true, message, color }
  }

  function hideSnackbar() {
    snackbar.value.show = false
  }

  return {
    theme,
    sidebarOpen,
    snackbar,
    toggleTheme,
    toggleSidebar,
    showSnackbar,
    hideSnackbar,
  }
})
