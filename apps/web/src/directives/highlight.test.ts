/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { vHighlight } from './highlight'

const { highlightElement, getHighlighter } = vi.hoisted(() => {
  const highlightElement = vi.fn()
  return {
    highlightElement,
    getHighlighter: vi.fn(() => ({ highlightElement })),
  }
})

vi.mock('@devtest/shared', () => ({
  getHighlighter,
  // El registro real vive en `@devtest/shared`; aquí solo interesa que la
  // directiva pida la normalización y use lo que devuelve.
  normalizeCodeLanguage: (lang?: string) => (lang && lang !== 'unknown' ? lang : 'plaintext'),
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
    ;(
      vHighlight as unknown as {
        mounted: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void
      }
    ).mounted(el, binding, null as never, null as never)
    expect(el.classList.contains('language-javascript')).toBe(true)
    expect(highlightElement).toHaveBeenCalledWith(el)
  })

  it('mounted falls back to plaintext when language not found', () => {
    const el = createEl('unknown')
    const binding = { value: 'unknown' } as never
    ;(
      vHighlight as unknown as {
        mounted: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void
      }
    ).mounted(el, binding, null as never, null as never)
    expect(el.classList.contains('language-plaintext')).toBe(true)
  })

  it('mounted uses dataset lang when binding undefined', () => {
    const el = createEl('go')
    const binding = { value: undefined } as never
    ;(
      vHighlight as unknown as {
        mounted: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void
      }
    ).mounted(el, binding, null as never, null as never)
    expect(el.classList.contains('language-go')).toBe(true)
  })

  it('updated highlights same language', () => {
    const el = createEl('javascript')
    el.classList.add('language-javascript')
    const binding = { value: 'javascript' } as never
    ;(
      vHighlight as unknown as {
        updated: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void
      }
    ).updated(el, binding, null as never, null as never)
    expect(highlightElement).toHaveBeenCalled()
  })

  it('updated switches language class', () => {
    const el = createEl('javascript')
    el.classList.add('language-javascript')
    const binding = { value: 'python' } as never
    ;(
      vHighlight as unknown as {
        updated: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void
      }
    ).updated(el, binding, null as never, null as never)
    expect(el.classList.contains('language-python')).toBe(true)
    expect(el.classList.contains('language-javascript')).toBe(false)
    expect(highlightElement).toHaveBeenCalled()
  })

  it('updated repaints when the language did not change', () => {
    const el = createEl('go')
    el.classList.add('language-go')
    el.setAttribute('data-highlighted', 'yes')
    const binding = { value: 'go' } as never
    ;(
      vHighlight as unknown as {
        updated: (el: HTMLElement, binding: unknown, a: unknown, b: unknown) => void
      }
    ).updated(el, binding, null as never, null as never)
    expect(el.getAttribute('data-highlighted')).toBeNull()
    expect(highlightElement).toHaveBeenCalledWith(el)
  })
})
