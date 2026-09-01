/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { vHighlight } from './highlight'
import hljs from 'highlight.js/lib/core'

vi.mock('highlight.js/lib/core', () => ({
  default: {
    highlightElement: vi.fn(),
    getLanguage: vi.fn(() => true),
    registerLanguage: vi.fn(),
  },
}))

describe('vHighlight directive', () => {
  function createEl(lang?: string) {
    const el = document.createElement('code')
    if (lang) el.dataset.lang = lang
    el.textContent = 'const x = 1'
    return el
  }

  it('mounted adds language class and highlights', () => {
    const el = createEl('javascript')
    const binding = { value: 'javascript' } as never
    ;(vHighlight as unknown as { mounted: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void }).mounted(el, binding, null as never, null as never)
    expect(el.classList.contains('language-javascript')).toBe(true)
    expect(hljs.highlightElement).toHaveBeenCalledWith(el)
  })

  it('mounted falls back to plaintext when language not found', () => {
    vi.mocked(hljs.getLanguage).mockReturnValueOnce(false as never)
    const el = createEl('unknown')
    const binding = { value: 'unknown' } as never
    ;(vHighlight as unknown as { mounted: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void }).mounted(el, binding, null as never, null as never)
    expect(el.classList.contains('language-plaintext')).toBe(true)
  })

  it('mounted uses dataset lang when binding undefined', () => {
    const el = createEl('go')
    const binding = { value: undefined } as never
    ;(vHighlight as unknown as { mounted: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void }).mounted(el, binding, null as never, null as never)
    expect(el.classList.contains('language-go')).toBe(true)
  })

  it('updated highlights same language', () => {
    const el = createEl('javascript')
    el.classList.add('language-javascript')
    const binding = { value: 'javascript' } as never
    ;(vHighlight as unknown as { updated: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void }).updated(el, binding, null as never, null as never)
    expect(hljs.highlightElement).toHaveBeenCalled()
  })

  it('updated switches language class', () => {
    const el = createEl('javascript')
    el.classList.add('language-javascript')
    const binding = { value: 'python' } as never
    vi.mocked(hljs.getLanguage).mockReturnValue(true as never)
    ;(vHighlight as unknown as { updated: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void }).updated(el, binding, null as never, null as never)
    expect(el.classList.contains('language-python')).toBe(true)
    expect(hljs.highlightElement).toHaveBeenCalled()
  })
})
