import '@testing-library/jest-dom/vitest'
import { config } from '@vue/test-utils'
import { createVuetify } from 'vuetify'

const vuetify = createVuetify({})

config.global.plugins = [vuetify]

// Stub ResizeObserver required by Vuetify
globalThis.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
} as unknown as typeof ResizeObserver

// Stub visualViewport for Vuetify VOverlay
Object.defineProperty(window, 'visualViewport', {
  writable: true,
  value: {
    width: 1024,
    height: 768,
    offsetTop: 0,
    offsetLeft: 0,
    pageTop: 0,
    pageLeft: 0,
    scale: 1,
    addEventListener: () => {},
    removeEventListener: () => {},
  },
})
Object.defineProperty(globalThis, 'visualViewport', {
  writable: true,
  value: window.visualViewport,
})

// Stub matchMedia / requestAnimationFrame for jsdom
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
})

if (!window.requestAnimationFrame) {
  window.requestAnimationFrame = (cb: FrameRequestCallback) =>
    setTimeout(cb, 0) as unknown as number
  window.cancelAnimationFrame = (id: number) => clearTimeout(id)
}
