import { codeTheme } from '@/core/theme/tokens'

/**
 * markdown-it agrega un `\n` extra al final del contenido de los bloques de
 * código; se recorta para no dejar una línea vacía dentro del bloque.
 */
export function stripTrailingNewline(content: string | undefined): string {
  if (!content) return ''
  return content.endsWith('\n') ? content.slice(0, -1) : content
}

/**
 * Resuelve el color de un token de `tokenizeCode`. Prueba de la clase más
 * interna a la más externa porque highlight.js emite clases combinadas
 * (`hljs-title function_`). Devuelve `undefined` si ninguna está en el tema
 * (el texto hereda `colors.foreground`).
 */
export function resolveCodeColor(classes: readonly string[]): string | undefined {
  for (let index = classes.length - 1; index >= 0; index -= 1) {
    const className = classes[index]
    if (!className) continue

    const color = codeTheme[className]
    if (color) return color
  }

  return undefined
}
