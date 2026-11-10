import type { Directive, DirectiveBinding } from 'vue'
import { getHighlighter, normalizeCodeLanguage } from '@devtest/shared'

// El registro de lenguajes y los alias viven en `@devtest/shared` para que web
// y mobile resalten exactamente lo mismo.
function getLanguage(el: HTMLElement, binding: DirectiveBinding<string | undefined>): string {
  return normalizeCodeLanguage(binding.value ?? el.dataset.lang)
}

export const vHighlight: Directive<HTMLElement, string | undefined> = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | undefined>) {
    const lang = getLanguage(el, binding)
    el.classList.add(`language-${lang}`)
    getHighlighter().highlightElement(el)
  },
  updated(el: HTMLElement, binding: DirectiveBinding<string | undefined>) {
    const lang = getLanguage(el, binding)

    if (!el.classList.contains(`language-${lang}`)) {
      el.classList.remove(...Array.from(el.classList).filter((c) => c.startsWith('language-')))
      el.classList.add(`language-${lang}`)
    }

    // highlight.js marca el elemento; hay que limpiar la marca para repintar.
    el.removeAttribute('data-highlighted')
    getHighlighter().highlightElement(el)
  },
}
