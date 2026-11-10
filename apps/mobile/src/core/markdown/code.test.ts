import { normalizeCodeLanguage, tokenizeCode } from '@devtest/shared'

import { resolveCodeColor, stripTrailingNewline } from './code'

describe('stripTrailingNewline', () => {
  it('recorta un único salto de línea final', () => {
    expect(stripTrailingNewline('func main() {}\n')).toBe('func main() {}')
  })

  it('deja intacto el contenido sin salto final', () => {
    expect(stripTrailingNewline('func main() {}')).toBe('func main() {}')
  })

  it('devuelve cadena vacía para contenido indefinido', () => {
    expect(stripTrailingNewline(undefined)).toBe('')
  })
})

describe('resolveCodeColor', () => {
  it('usa la clase más interna cuando hay varias', () => {
    // highlight.js emite `hljs-function` > `hljs-keyword` para `func`.
    expect(resolveCodeColor(['hljs-function', 'hljs-keyword'])).toBe('#C792EA')
  })

  it('devuelve undefined para clases fuera del tema', () => {
    expect(resolveCodeColor(['hljs-unknown'])).toBeUndefined()
  })

  it('devuelve undefined sin clases (texto plano)', () => {
    expect(resolveCodeColor([])).toBeUndefined()
  })
})

describe('normalizeCodeLanguage', () => {
  it('resuelve alias conocidos por highlight.js', () => {
    expect(normalizeCodeLanguage('ts')).toBe('typescript')
    expect(normalizeCodeLanguage('js')).toBe('javascript')
    expect(normalizeCodeLanguage('py')).toBe('python')
  })

  it('resuelve los alias propios del proyecto', () => {
    expect(normalizeCodeLanguage('vue')).toBe('xml')
    expect(normalizeCodeLanguage('lang')).toBe('plaintext')
  })

  it('cae en plaintext cuando el lenguaje no existe', () => {
    expect(normalizeCodeLanguage('no-existe')).toBe('plaintext')
    expect(normalizeCodeLanguage('')).toBe('plaintext')
    expect(normalizeCodeLanguage(undefined)).toBe('plaintext')
  })
})

describe('tokenizeCode', () => {
  it('devuelve tokens con clases hljs para un lenguaje soportado', () => {
    const tokens = tokenizeCode('func main() {}', 'go')

    expect(tokens.some((token) => token.classes.includes('hljs-keyword'))).toBe(true)
    expect(tokens.map((token) => token.text).join('')).toBe('func main() {}')
  })

  it('desescapa las entidades HTML del código', () => {
    const tokens = tokenizeCode('if (a < b && c > d) {}', 'javascript')

    expect(tokens.map((token) => token.text).join('')).toBe('if (a < b && c > d) {}')
  })

  it('devuelve un único token sin clases para lenguajes desconocidos', () => {
    expect(tokenizeCode('x = 1', 'no-existe')).toEqual([{ text: 'x = 1', classes: [] }])
  })

  it('devuelve una lista vacía sin código', () => {
    expect(tokenizeCode('', 'go')).toEqual([])
  })
})
