import type { Topic } from '../types/topic.types'

export const IMPORT_CSV_HEADER = 'type,content,difficulty,topics,explanation,options'

export const IMPORT_CSV_TEMPLATE = `${IMPORT_CSV_HEADER}
single_choice,"¿Qué hace \`defer\` en Go?

\`\`\`go
defer fmt.Println(""adios"")
fmt.Println(""hola"")
\`\`\`",beginner,go,"**defer** programa la llamada para ejecutarse justo antes de que la función retorne, en orden **LIFO**. Ideal para liberar recursos como \`file.Close()\`.","[v] Se ejecuta al retornar la función (LIFO) | [ ] Se ejecuta al inicio de la iteración | [ ] Cancela la goroutine | [ ] Convierte la función en asíncrona"
multiple_choice,"¿Cuáles de estas afirmaciones sobre **interfaces** en Go son correctas?",intermediate,go,"En Go las interfaces se implementan **implícitamente**: un tipo la satisface si posee todos sus métodos. La interfaz vacía \`interface{}\` puede contener cualquier valor.","[v] La implementación es implícita, no se declara \`implements\` | [ ] Un tipo debe declarar explícitamente qué interfaces implementa | [v] \`interface{}\` puede contener cualquier valor | [ ] Las interfaces solo pueden tener métodos privados"
true_false,"¿En Go las funciones anónimas pueden capturar variables del alcance externo?",beginner,go,"**Sí**, las closures capturan variables por referencia: si la variable cambia después, la closure ve el valor actualizado.","[v] Verdadero | [ ] Falso"
`

export const IMPORT_CSV_TEMPLATE_TRUE_FALSE = `${IMPORT_CSV_HEADER}
true_false,"¿Go permite múltiples return values?",intermediate,go,"En Go es habitual devolver \`(resultado, error)\` en una sola función.","[v] Verdadero | [ ] Falso"
`

export const IMPORT_CSV_TEMPLATE_MULTIPLE = `${IMPORT_CSV_HEADER}
multiple_choice,"¿Cuáles de estas afirmaciones sobre **goroutines** en Go son correctas?",intermediate,go,"Las goroutines son green threads multiplexadas por el runtime. Se comunican con channels.","[v] Son green threads multiplexadas por el runtime | [ ] Se comunican con \`shared memory\` sin sync | [v] Se comunican principalmente con channels | [ ] Cada goroutine tiene su propio PID"
`

/** Devuelve la plantilla CSV de ejemplo para el tipo de pregunta dado. */
export function getImportCsvTemplate(kind: PromptQuestionKind): string {
  switch (kind) {
    case 'true_false':
      return IMPORT_CSV_TEMPLATE_TRUE_FALSE
    case 'multiple_choice':
      return IMPORT_CSV_TEMPLATE_MULTIPLE
    default:
      return IMPORT_CSV_TEMPLATE
  }
}

export type ImportQuestionKind = 'single_choice' | 'multiple_choice' | 'true_false'

export const IMPORT_QUESTION_KINDS: Array<{ title: string; value: ImportQuestionKind }> = [
  { title: 'Selección única', value: 'single_choice' },
  { title: 'Selección múltiple', value: 'multiple_choice' },
  { title: 'Verdadero / Falso', value: 'true_false' },
]

/**
 * Tipo elegido para generar el prompt con IA. `mixed` reparte la cantidad pedida
 * entre los 3 tipos (single_choice, multiple_choice y true_false).
 */
export type PromptQuestionKind = ImportQuestionKind | 'mixed'

export const IMPORT_PROMPT_KINDS: Array<{
  title: string
  value: PromptQuestionKind
}> = [{ title: 'Mixta (3 tipos)', value: 'mixed' }, ...IMPORT_QUESTION_KINDS]

export type PromptDifficulty = 'mixed' | 'beginner' | 'intermediate' | 'advanced'

export const PROMPT_DIFFICULTIES: Array<{ title: string; value: PromptDifficulty }> = [
  { title: 'Mixta (variada)', value: 'mixed' },
  { title: 'Principiante', value: 'beginner' },
  { title: 'Intermedio', value: 'intermediate' },
  { title: 'Avanzado', value: 'advanced' },
]

export type PromptObjective = 'review' | 'interview' | 'learning' | 'exam' | 'daily' | 'custom'

export const PROMPT_OBJECTIVES: Array<{
  title: string
  value: PromptObjective
  description: string
}> = [
  {
    title: 'Repaso / refuerzo',
    value: 'review',
    description: 'Reforzar conceptos ya conocidos y errores comunes',
  },
  {
    title: 'Entrevista técnica',
    value: 'interview',
    description: 'Estilo entrevista real, casos prácticos y profundidad',
  },
  {
    title: 'Aprendizaje',
    value: 'learning',
    description: 'Pedagógico, de lo fundamental a lo intermedio',
  },
  {
    title: 'Autoevaluación / examen',
    value: 'exam',
    description: 'Con trampas sutiles, escenarios y distractores',
  },
  {
    title: 'Práctica diaria',
    value: 'daily',
    description: 'Cortas y variadas para mantener el ritmo',
  },
  {
    title: 'Personalizado',
    value: 'custom',
    description: 'Define tu propio contexto abajo',
  },
]

const DIFFICULTY_PROMPT: Record<PromptDifficulty, string> = {
  mixed:
    'Dificultad: MIXTA — distribuye las preguntas entre beginner, intermediate y advanced de forma equilibrada. Varía la profundidad para cubrir distintos niveles.',
  beginner:
    'Dificultad: PRINCIPIANTE (beginner) — enfócate en conceptos fundamentales, definiciones y uso básico. Evita detalles avanzados.',
  intermediate:
    'Dificultad: INTERMEDIA (intermediate) — enfócate en aplicación práctica, casos comunes, diferencias sutiles y buenas prácticas.',
  advanced:
    'Dificultad: AVANZADA (advanced) — enfócate en profundidad, edge cases, rendimiento, trade-offs y decisiones de diseño.',
}

const OBJECTIVE_PROMPT: Record<Exclude<PromptObjective, 'custom'>, string> = {
  review: `OBJETIVO: REPASO / REFUERZO.
Las preguntas deben servir para repasar y reforzar conceptos ya conocidos.
Enfócate en los conceptos clave del tema, errores comunes y confusiones frecuentes.
Evita preguntas triviales o demasiado genéricas: profundiza en el "por qué" y en el uso correcto.
El tono debe ser de refuerzo de memoria y comprensión.`,
  interview: `OBJETIVO: ENTREVISTA TÉCNICA.
Las preguntas deben tener el estilo de una entrevista técnica real.
Prioriza preguntas que suelen aparecer en entrevistas de trabajo para ese stack/tema,
incluyendo casos prácticos, gotchas, trade-offs y preguntas de profundidad.
Evita definiciones triviales o genéricas; busca que cada pregunta discrimine nivel.`,
  learning: `OBJETIVO: APRENDIZAJE.
Las preguntas deben servir para aprender el tema de forma progresiva.
Parte de conceptos fundamentales hacia intermedios, con explicaciones pedagógicas
que enseñen el porqué, no solo el qué. Secuencia las preguntas de lo básico a lo aplicado.`,
  exam: `OBJETIVO: AUTOEVALUACIÓN / EXAMEN.
Las preguntas deben evaluar el nivel como en un examen o certificación.
Incluye trampas sutiles, escenarios y distractores plausibles.
La explicación debe indicar por qué cada alternativa incorrecta es incorrecta.`,
  daily: `OBJETIVO: PRÁCTICA DIARIA.
Preguntas cortas, variadas y directas para mantener la práctica diaria y la retención.
Mezcla ángulos distintos del tema y niveles de dificultad. Sé conciso pero preciso.`,
}

const OPTION_RULES: Record<ImportQuestionKind, string> = {
  single_choice: `Reglas de options:
- EXACTAMENTE 4 opciones por pregunta.
- single_choice: EXACTAMENTE 1 opción con [v].
- Usa: [v] para correctas, [ ] para incorrectas.`,
  multiple_choice: `Reglas de options:
- EXACTAMENTE 4 opciones por pregunta.
- multiple_choice: entre 2 y 3 opciones con [v] (nunca 1 ni las 4).
- Usa: [v] para correctas, [ ] para incorrectas.`,
  true_false: `Reglas de options (Verdadero / Falso):
- EXACTAMENTE 2 opciones por pregunta: una "Verdadero" y una "Falso".
- EXACTAMENTE 1 opción con [v] (la correcta), la otra con [ ].
- NO uses más de 2 opciones.
- Puedes reordenar cuál es Verdadero y cuál es Falso para evitar que siempre sea la primera.
- Puedes usar la forma "Verdadero: ..." o "Falso: ..." seguido de una breve razón.`,
}

/** Fila CSV de ejemplo por tipo de pregunta (usa el primer slug seleccionado). */
function exampleRow(kind: ImportQuestionKind, slugsExample: string): string {
  switch (kind) {
    case 'multiple_choice':
      return `multiple_choice,"¿Cuáles de estas afirmaciones sobre **goroutines** en Go son correctas?",intermediate,${slugsExample},"Las goroutines son hilos ligeros gestionados por el runtime; se comunican principalmente con channels.","[v] Son hilos ligeros gestionados por el runtime | [ ] Cada goroutine tiene su propio proceso | [v] Se comunican principalmente con channels | [ ] Solo puede haber una goroutine por programa"`
    case 'true_false':
      return `true_false,"¿En Go las funciones anónimas pueden capturar variables del alcance externo?",beginner,${slugsExample},"**Sí**, las closures en Go capturan variables por referencia. La variable capturada puede modificar su valor después.","[v] Verdadero | [ ] Falso"`
    default:
      return `single_choice,"¿Qué imprime este código en Go?

\`\`\`go
fmt.Println(len(""hola""))
\`\`\`",beginner,${slugsExample},"**len** devuelve la longitud en bytes de la cadena. Para ""hola"" son 4 bytes.","[v] 4 | [ ] 5 | [ ] 0 | [ ] error"`
  }
}

const ENRICHED_FORMAT_BLOCK = `FORMATO ENRIQUECIDO (soportado por la UI en content, explanation y options):
- Usa \`código inline\` para identificadores, funciones, comandos y valores literales (ej: \`defer\`, \`Array<T>\`, \`useState\`).
- Usa bloques de código con fences cuando el enunciado o la explicación lo requiera:
  \`\`\`go
  defer file.Close()
  \`\`\`
  Idiomas soportados: go, javascript, typescript, python, java, csharp, rust, sql, bash. Si el ejemplo no es código ejecutable usa plaintext.
- Usa **negrita** para resaltar conceptos clave.
- Aprovecha el formato para que la pregunta sea atractiva y legible: si hay código, muéstralo con formato; si hay un concepto clave, destácalo.`

const CSV_RULES_BLOCK = `REGLAS DE CSV:
- Si un campo contiene comas, saltos de línea o comillas dobles, enciérralo entre comillas dobles.
- Escapa comillas dobles internas duplicándolas ("").
- Los campos multilínea (con bloques de código) DEBEN ir entre comillas dobles y conservar los saltos de línea dentro del campo.
- Solo texto, comas y ; están permitidos en los slugs.`

export interface BuildImportPromptOptions {
  kind?: PromptQuestionKind
  difficulty?: PromptDifficulty
  objective?: PromptObjective
  customContext?: string
}

export function buildImportPrompt(
  topics: Array<Topic>,
  quantity: number = 20,
  opts: BuildImportPromptOptions = {},
): string {
  const qty = Math.min(50, Math.max(1, Math.floor(quantity) || 20))
  const kind: PromptQuestionKind = opts.kind ?? 'mixed'
  const difficulty: PromptDifficulty = opts.difficulty ?? 'mixed'
  const objective: PromptObjective = opts.objective ?? 'review'
  const customContext = opts.customContext?.trim() ?? ''

  const topicLines =
    topics.length > 0
      ? topics.map((t) => `- ${t.name} (slug: ${t.slug})`).join('\n')
      : '- (no se seleccionaron temas — usa slugs genéricos como "javascript", "python", "go")'

  const slugsExample = topics.length > 0 && topics[0] ? topics[0].slug : 'javascript'
  const topicsNote =
    topics.length > 1
      ? `Distribuye las ${qty} preguntas entre los ${topics.length} temas listados. Cada pregunta debe estar asociada a al menos uno de esos slugs (puede llevar uno o varios separados por ;, pero siempre de esta lista). Varía los temas para cubrirlos todos.`
      : topics.length === 1
        ? `Todas las preguntas deben usar el slug "${slugsExample}" en la columna topics.`
        : 'Cada pregunta debe usar un slug coherente en la columna topics.'

  const difficultyBlock = DIFFICULTY_PROMPT[difficulty]

  const objectiveBlock =
    objective === 'custom'
      ? customContext
        ? `OBJETIVO PERSONALIZADO:\n${customContext}`
        : 'OBJETIVO PERSONALIZADO: (sin contexto adicional — genera preguntas equilibradas y específicas, evitando definiciones genéricas).'
      : OBJECTIVE_PROMPT[objective as Exclude<PromptObjective, 'custom'>]

  const customBlock =
    objective !== 'custom' && customContext
      ? `\nCONTEXTO ADICIONAL DEL USUARIO:\n${customContext}\n`
      : ''

  const commonArgs = [
    qty,
    difficultyBlock,
    objectiveBlock,
    customBlock,
    topicLines,
    topicsNote,
    slugsExample,
  ] as const

  if (kind === 'mixed') {
    return buildMixedPrompt(...commonArgs)
  }

  if (kind === 'true_false') {
    return buildTrueFalsePrompt(...commonArgs)
  }

  const typeLine: ImportQuestionKind = kind

  return `Eres un generador de preguntas para una plataforma de estudio de desarrollo de software.

Genera EXACTAMENTE ${qty} preguntas en formato CSV de tipo ${typeLine}. NO uses code_completion.

${difficultyBlock}

${objectiveBlock}${customBlock}
Evita preguntas genéricas o triviales. Cada pregunta debe ser específica del tema, con enunciado concreto y, cuando aporte valor, apoyada en un ejemplo de código o escenario práctico.

FORMATO OBLIGATORIO:
- Encabezado: ${IMPORT_CSV_HEADER}
- Delimitador: coma (,) con comillas dobles si el campo contiene comas o saltos de línea.
- Codificación UTF-8.

Columnas:
- type: ${typeLine}
- content: enunciado en ESPAÑOL (claro y específico). Puede incluir formato enriquecido (ver FORMATO ENRIQUECIDO).
- difficulty: beginner | intermediate | advanced
- topics: slugs separados por ; — USA EXACTAMENTE ESTOS SLUGS, sin traducir ni cambiar mayúsculas:
${topicLines}
  ${topicsNote}
  Ejemplo topics: ${slugsExample}  (si varias: ${slugsExample};otro-slug)
- explanation: explicación educativa en español. Puede incluir formato enriquecido.
- options: opciones separadas por | . Cada opción con prefijo [v] si es CORRECTA o [ ] si es INCORRECTA.

${OPTION_RULES[typeLine]}
- Las opciones también pueden usar \`código inline\` y **negrita**, pero NO bloques de código multilínea ni el carácter | dentro del texto de una opción (es el separador).

${ENRICHED_FORMAT_BLOCK}

${CSV_RULES_BLOCK}

REGLAS IMPORTANTES:
- TODO el contenido debe estar en ESPAÑOL (con tildes y ortografía correcta).
- NO añadas texto, comentarios, explicaciones ni fences markdown FUERA del CSV. Devuelve SOLO el CSV (encabezado + ${qty} filas).
- Dentro de los campos SÍ usa markdown inline/bloques según las reglas de FORMATO ENRIQUECIDO.
- Genera exactamente ${qty} filas de datos + la fila de encabezado (total ${qty + 1} filas).
- No repitas preguntas.
- No pongas comillas simples alrededor de [v]/[ ].

Ejemplo de fila:
${exampleRow(typeLine, slugsExample)}

Ahora genera el CSV solicitado.`
}

function buildTrueFalsePrompt(
  qty: number,
  difficultyBlock: string,
  objectiveBlock: string,
  customBlock: string,
  topicLines: string,
  topicsNote: string,
  slugsExample: string,
): string {
  return `Eres un generador de preguntas para una plataforma de estudio de desarrollo de software.

Genera EXACTAMENTE ${qty} preguntas en formato CSV de tipo true_false (Verdadero / Falso). NO uses code_completion ni single_choice ni multiple_choice.

${difficultyBlock}

${objectiveBlock}${customBlock}
Evita preguntas genéricas o triviales. Cada pregunta debe ser específica del tema, con enunciado concreto y, cuando aporte valor, apoyada en un ejemplo de código o escenario práctico.

FORMATO OBLIGATORIO:
- Encabezado: ${IMPORT_CSV_HEADER}
- Delimitador: coma (,) con comillas dobles si el campo contiene comas o saltos de línea.
- Codificación UTF-8.

Columnas:
- type: true_false
- content: enunciado en ESPAÑOL (claro y específico). Puede incluir formato enriquecido (ver FORMATO ENRIQUECIDO).
- difficulty: beginner | intermediate | advanced
- topics: slugs separados por ; — USA EXACTAMENTE ESTOS SLUGS, sin traducir ni cambiar mayúsculas:
${topicLines}
  ${topicsNote}
  Ejemplo topics: ${slugsExample}  (si varias: ${slugsExample};otro-slug)
- explanation: explicación educativa en español. Puede incluir formato enriquecido.
- options: EXACTAMENTE 2 opciones separadas por |. Una marcada con [v] (correcta) y otra con [ ] (incorrecta).

${OPTION_RULES.true_false}

${ENRICHED_FORMAT_BLOCK}

${CSV_RULES_BLOCK}

REGLAS IMPORTANTES:
- TODO el contenido debe estar en ESPAÑOL (con tildes y ortografía correcta).
- NO añadas texto, comentarios, explicaciones ni fences markdown FUERA del CSV. Devuelve SOLO el CSV (encabezado + ${qty} filas).
- Dentro de los campos SÍ usa markdown inline/bloques según las reglas de FORMATO ENRIQUECIDO.
- Genera exactamente ${qty} filas de datos + la fila de encabezado (total ${qty + 1} filas).
- No repitas preguntas.
- No pongas comillas simples alrededor de [v]/[ ].

Ejemplo de fila:
${exampleRow('true_false', slugsExample)}

Ahora genera el CSV solicitado.`
}

function buildMixedPrompt(
  qty: number,
  difficultyBlock: string,
  objectiveBlock: string,
  customBlock: string,
  topicLines: string,
  topicsNote: string,
  slugsExample: string,
): string {
  const perKind = Math.max(1, Math.round(qty / 3))

  return `Eres un generador de preguntas para una plataforma de estudio de desarrollo de software.

Genera EXACTAMENTE ${qty} preguntas en formato CSV MEZCLANDO los 3 tipos: single_choice (selección única), multiple_choice (selección múltiple) y true_false (Verdadero / Falso). NO uses code_completion.

${difficultyBlock}

${objectiveBlock}${customBlock}
Evita preguntas genéricas o triviales. Cada pregunta debe ser específica del tema, con enunciado concreto y, cuando aporte valor, apoyada en un ejemplo de código o escenario práctico.

DISTRIBUCIÓN OBLIGATORIA DE TIPOS:
- Reparte las ${qty} preguntas entre los 3 tipos de forma equilibrada (~${perKind} de cada uno).
- Si son 6 o más preguntas, incluye AL MENOS 2 de cada tipo.
- Intercala los tipos: no agrupes todas las de un mismo tipo al inicio ni al final.

FORMATO OBLIGATORIO:
- Encabezado: ${IMPORT_CSV_HEADER}
- Delimitador: coma (,) con comillas dobles si el campo contiene comas o saltos de línea.
- Codificación UTF-8.

Columnas:
- type: single_choice | multiple_choice | true_false
- content: enunciado en ESPAÑOL (claro y específico). Puede incluir formato enriquecido (ver FORMATO ENRIQUECIDO).
- difficulty: beginner | intermediate | advanced
- topics: slugs separados por ; — USA EXACTAMENTE ESTOS SLUGS, sin traducir ni cambiar mayúsculas:
${topicLines}
  ${topicsNote}
  Ejemplo topics: ${slugsExample}  (si varias: ${slugsExample};otro-slug)
- explanation: explicación educativa en español. Puede incluir formato enriquecido.
- options: opciones separadas por | . Cada opción con prefijo [v] si es CORRECTA o [ ] si es INCORRECTA.

Reglas de options SEGÚN EL type DE CADA FILA (cada fila usa solo las de su tipo, no mezcles formatos):
- type=single_choice: EXACTAMENTE 4 opciones y EXACTAMENTE 1 con [v].
- type=multiple_choice: EXACTAMENTE 4 opciones y entre 2 y 3 con [v] (nunca 1 ni las 4).
- type=true_false: EXACTAMENTE 2 opciones ("Verdadero" y "Falso"), EXACTAMENTE 1 con [v] y la otra con [ ]. Puedes reordenarlas.
- Las opciones también pueden usar \`código inline\` y **negrita**, pero NO bloques de código multilínea ni el carácter | dentro del texto de una opción (es el separador).

${ENRICHED_FORMAT_BLOCK}

${CSV_RULES_BLOCK}

REGLAS IMPORTANTES:
- TODO el contenido debe estar en ESPAÑOL (con tildes y ortografía correcta).
- NO añadas texto, comentarios, explicaciones ni fences markdown FUERA del CSV. Devuelve SOLO el CSV (encabezado + ${qty} filas).
- Dentro de los campos SÍ usa markdown inline/bloques según las reglas de FORMATO ENRIQUECIDO.
- Genera exactamente ${qty} filas de datos + la fila de encabezado (total ${qty + 1} filas).
- No repitas preguntas.
- No pongas comillas simples alrededor de [v]/[ ].

Ejemplos de fila (uno por tipo):
${exampleRow('single_choice', slugsExample)}
${exampleRow('multiple_choice', slugsExample)}
${exampleRow('true_false', slugsExample)}

Ahora genera el CSV solicitado.`
}

// ─── Tips ────────────────────────────────────────────────────────────

export const TIPS_CSV_HEADER = 'content,topic'

/** Prompt externo para generar tips con IA (copy-paste a ChatGPT/Claude). */
export function buildTipsImportPrompt(
  topics: Array<Topic>,
  quantity: number = 20,
  opts: { customContext?: string } = {},
): string {
  const qty = Math.min(50, Math.max(1, Math.floor(quantity) || 20))
  const customContext = opts.customContext?.trim() ?? ''

  const topicLines =
    topics.length > 0
      ? topics.map((t) => `- ${t.name} (slug: ${t.slug})`).join('\n')
      : '- (no se seleccionaron temas — usa slugs genéricos como "javascript", "python", "go")'

  const slugsExample = topics.length > 0 && topics[0] ? topics[0].slug : 'javascript'
  const topicsNote =
    topics.length > 1
      ? `Distribuye los ${qty} tips entre los ${topics.length} temas listados. Cada tip debe estar asociado a exactamente UN slug de esta lista. Varía los temas para cubrirlos todos.`
      : topics.length === 1
        ? `Todos los tips deben usar el slug "${slugsExample}" en la columna topic.`
        : 'Cada tip debe usar un slug coherente en la columna topic.'

  const customBlock = customContext
    ? `\nCONTEXTO ADICIONAL DEL USUARIO:\n${customContext}\n`
    : ''

  return `Eres un generador de tips/consejos para una plataforma de estudio de desarrollo de software.

Genera EXACTAMENTE ${qty} tips en formato CSV. Cada tip es un consejo corto, accionable y específico sobre un tema técnico.

${customContext ? `CONTEXTO: ${customContext}\n` : ''}
${topicLines}

FORMATO OBLIGATORIO:
- Encabezado: ${TIPS_CSV_HEADER}
- Delimitador: coma (,).
- Codificación UTF-8.

Columnas:
- content: consejo conciso en ESPAÑOL (1-3 oraciones). Debe ser específico, accionable y educativo. Puede incluir \`código inline\` y **negrita**.
- topic: slug exacto de un tema existente. Debe pertenecer a esta lista:
${topicLines}
  ${topicsNote}
  Ejemplo topic: ${slugsExample}

REGLAS:
- TODO el contenido debe estar en ESPAÑOL.
- Cada tip debe cubrir un concepto DISTINTO (no repitas).
- Sé conciso: un tip no debería superar 2 oraciones.
- Evita definiciones genéricas; prioriza consejos prácticos.
- No incluyas texto FUERA del CSV. Devuelve SOLO el CSV (encabezado + ${qty} filas).

Ejemplo de fila:
"Usa \`defer\` para liberar archivos o conexiones: se ejecuta antes de que la función retorne (orden LIFO).",${slugsExample}

Ahora genera el CSV solicitado.`
}
