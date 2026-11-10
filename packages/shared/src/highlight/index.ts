import hljs from 'highlight.js/lib/core'
import type { HLJSApi, LanguageFn } from 'highlight.js'
import bash from 'highlight.js/lib/languages/bash'
import csharp from 'highlight.js/lib/languages/csharp'
import dockerfile from 'highlight.js/lib/languages/dockerfile'
import go from 'highlight.js/lib/languages/go'
import graphql from 'highlight.js/lib/languages/graphql'
import http from 'highlight.js/lib/languages/http'
import java from 'highlight.js/lib/languages/java'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import plaintext from 'highlight.js/lib/languages/plaintext'
import python from 'highlight.js/lib/languages/python'
import rust from 'highlight.js/lib/languages/rust'
import sql from 'highlight.js/lib/languages/sql'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'

/**
 * Lenguajes registrados en highlight.js. Fuente única para web y mobile: un
 * bloque con un lenguaje que no esté aquí se pinta como `plaintext`.
 *
 * Se importan uno por uno (nunca `highlight.js/lib/common`) para que el bundle
 * de mobile solo cargue lo que el contenido realmente usa.
 */
const LANGUAGE_DEFINITIONS: Record<string, LanguageFn> = {
  bash,
  csharp,
  dockerfile,
  go,
  graphql,
  http,
  java,
  javascript,
  json,
  plaintext,
  python,
  rust,
  sql,
  typescript,
  xml,
  yaml,
}

/** Nombres canónicos de los lenguajes soportados. */
export const CODE_LANGUAGES: string[] = Object.keys(LANGUAGE_DEFINITIONS).sort()

/** Lenguaje usado cuando no hay pista o el lenguaje no está registrado. */
export const DEFAULT_CODE_LANGUAGE = 'plaintext'

/**
 * Alias que highlight.js no conoce. Los que sí conoce (`js`, `ts`, `py`, `sh`,
 * `html`, `yml`, `gql`, `docker`, …) se resuelven con el índice de alias.
 */
const LANGUAGE_ALIASES: Record<string, string> = {
  'docker-compose': 'yaml',
  lang: DEFAULT_CODE_LANGUAGE,
  node: 'javascript',
  plain: DEFAULT_CODE_LANGUAGE,
  text: DEFAULT_CODE_LANGUAGE,
  txt: DEFAULT_CODE_LANGUAGE,
  vue: 'xml',
}

let highlighter: HLJSApi | null = null
let languageIndex: Map<string, string> | null = null

/**
 * Devuelve la instancia de highlight.js con todos los lenguajes registrados.
 * El registro es perezoso y memoizado: importar este módulo no cuesta nada.
 */
export function getHighlighter(): HLJSApi {
  if (highlighter) return highlighter

  for (const [name, definition] of Object.entries(LANGUAGE_DEFINITIONS)) {
    hljs.registerLanguage(name, definition)
  }

  highlighter = hljs
  return highlighter
}

/**
 * Índice `alias → nombre canónico`. highlight.js registra los alias, pero
 * `Language.name` es el nombre legible ("TypeScript"), así que se construye
 * aparte para poder devolver siempre la clave con la que se registró.
 */
function getLanguageIndex(): Map<string, string> {
  if (languageIndex) return languageIndex

  const instance = getHighlighter()
  const index = new Map<string, string>()

  for (const name of instance.listLanguages()) {
    index.set(name, name)
    for (const alias of instance.getLanguage(name)?.aliases ?? []) {
      index.set(alias, name)
    }
  }

  languageIndex = index
  return index
}

/**
 * Normaliza un lenguaje de fence (```` ```ts ````) al nombre canónico que
 * entiende highlight.js. Nunca falla: devuelve `plaintext` si no lo conoce.
 */
export function normalizeCodeLanguage(language?: string): string {
  const name = (language ?? '').trim().toLowerCase()
  if (!name) return DEFAULT_CODE_LANGUAGE

  const alias = LANGUAGE_ALIASES[name]
  if (alias) return alias

  return getLanguageIndex().get(name) ?? DEFAULT_CODE_LANGUAGE
}

/**
 * Fragmento de código ya resaltado (o no, si `classes` viene vacío).
 */
export type CodeToken = {
  /** Texto ya desescapado, listo para pintar. */
  text: string
  /**
   * Clases `hljs-*` aplicadas, de la más externa a la más interna. El
   * consumidor decide el color; vacío significa texto sin estilo propio.
   */
  classes: string[]
}

const ENTITIES: Record<string, string> = {
  '&amp;': '&',
  '&lt;': '<',
  '&gt;': '>',
  '&quot;': '"',
  '&#x27;': "'",
  '&#39;': "'",
}

function decodeEntities(text: string): string {
  if (!text.includes('&')) return text
  return text.replace(/&(?:amp|lt|gt|quot|#x27|#39);/g, (entity) => ENTITIES[entity] ?? entity)
}

function pushText(tokens: CodeToken[], raw: string, stack: string[][]): void {
  if (!raw) return

  const text = decodeEntities(raw)
  if (!text) return

  const classes = stack.flat()
  const last = tokens[tokens.length - 1]

  // Une fragmentos consecutivos con el mismo estilo para no crear un <Text>
  // por palabra.
  if (last && last.classes.join(' ') === classes.join(' ')) {
    last.text += text
    return
  }

  tokens.push({ text, classes })
}

/**
 * Convierte la salida HTML de highlight.js (`<span class="hljs-keyword">…`)
 * en tokens planos. La salida solo contiene esos `<span>` y entidades HTML,
 * así que un recorrido lineal basta y no necesita DOM.
 */
function parseHighlightedHtml(html: string): CodeToken[] {
  const tokens: CodeToken[] = []
  const stack: string[][] = []
  let index = 0

  while (index < html.length) {
    const open = html.indexOf('<span class="', index)
    const close = html.indexOf('</span>', index)
    const next = open === -1 ? close : close === -1 ? open : Math.min(open, close)

    if (next === -1) {
      pushText(tokens, html.slice(index), stack)
      break
    }

    if (next > index) pushText(tokens, html.slice(index, next), stack)

    if (next === open) {
      const classEnd = html.indexOf('"', open + 13)
      if (classEnd === -1) {
        pushText(tokens, html.slice(open), stack)
        break
      }
      stack.push(
        html
          .slice(open + 13, classEnd)
          .split(/\s+/)
          .filter(Boolean),
      )
      index = classEnd + 2
      continue
    }

    stack.pop()
    index = close + 7
  }

  return tokens
}

/**
 * Resalta `code` y devuelve los tokens con sus clases `hljs-*`. Un lenguaje
 * desconocido (o vacío) devuelve un único token sin clases.
 */
export function tokenizeCode(code: string, language?: string): CodeToken[] {
  if (!code) return []

  const instance = getHighlighter()
  const lang = normalizeCodeLanguage(language)

  if (lang === DEFAULT_CODE_LANGUAGE) return [{ text: code, classes: [] }]

  const { value } = instance.highlight(code, {
    language: lang,
    ignoreIllegals: true,
  })
  return parseHighlightedHtml(value)
}
