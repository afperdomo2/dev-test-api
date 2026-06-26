import type { Directive, DirectiveBinding } from 'vue'
import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import python from 'highlight.js/lib/languages/python'
import java from 'highlight.js/lib/languages/java'
import csharp from 'highlight.js/lib/languages/csharp'
import rust from 'highlight.js/lib/languages/rust'
import plaintext from 'highlight.js/lib/languages/plaintext'

hljs.registerLanguage('go', go)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('java', java)
hljs.registerLanguage('csharp', csharp)
hljs.registerLanguage('rust', rust)
hljs.registerLanguage('plaintext', plaintext)

function getLanguage(el: HTMLElement, binding: DirectiveBinding<string | undefined>): string {
  const lang = binding.value ?? (el.dataset.lang || 'plaintext')
  if (hljs.getLanguage(lang)) {
    return lang
  }
  return 'plaintext'
}

export const vHighlight: Directive<HTMLElement, string | undefined> = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | undefined>) {
    const lang = getLanguage(el, binding)
    el.classList.add(`language-${lang}`)
    hljs.highlightElement(el)
  },
  updated(el: HTMLElement, binding: DirectiveBinding<string | undefined>) {
    const lang = getLanguage(el, binding)
    if (el.classList.contains(`language-${lang}`)) {
      el.removeAttribute('data-highlighted')
      hljs.highlightElement(el)
      return
    }
    el.classList.remove(...Array.from(el.classList).filter((c) => c.startsWith('language-')))
    el.classList.add(`language-${lang}`)
    el.removeAttribute('data-highlighted')
    hljs.highlightElement(el)
  },
}
