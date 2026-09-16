import { setActivePinia, createPinia } from 'pinia'
import { useAppStore } from './app.store'

describe('app.store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('defaults', () => {
    const store = useAppStore()
    expect(store.theme).toBe('light')
    expect(store.sidebarOpen).toBe(true)
    expect(store.snackbar).toEqual({ show: false, message: '', color: 'success' })
  })

  it('toggleTheme switches', () => {
    const store = useAppStore()
    store.toggleTheme()
    expect(store.theme).toBe('dark')
    store.toggleTheme()
    expect(store.theme).toBe('light')
  })

  it('toggleSidebar flips', () => {
    const store = useAppStore()
    const before = store.sidebarOpen
    store.toggleSidebar()
    expect(store.sidebarOpen).toBe(!before)
    store.toggleSidebar()
    expect(store.sidebarOpen).toBe(before)
  })

  it('showSnackbar sets message/color and show', () => {
    const store = useAppStore()
    store.showSnackbar('hello', 'error')
    expect(store.snackbar).toEqual({ show: true, message: 'hello', color: 'error' })
  })

  it('showSnackbar default color success', () => {
    const store = useAppStore()
    store.showSnackbar('msg')
    expect(store.snackbar.color).toBe('success')
    expect(store.snackbar.show).toBe(true)
  })

  it('hideSnackbar sets show false', () => {
    const store = useAppStore()
    store.showSnackbar('x')
    store.hideSnackbar()
    expect(store.snackbar.show).toBe(false)
    expect(store.snackbar.message).toBe('x') // keeps message
  })
})
