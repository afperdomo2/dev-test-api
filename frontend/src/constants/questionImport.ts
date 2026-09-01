import type { Topic } from '@/types/topic.types'

export const IMPORT_CSV_HEADER = 'type,content,difficulty,topics,explanation,options'

export const IMPORT_CSV_TEMPLATE = `${IMPORT_CSV_HEADER}
single_choice,"¿Cuál es la forma correcta de declarar una variable en Go? (usa :=)",beginner,go,"El operador := declara e inicializa con inferencia de tipo","[v] x := 5 | [ ] var x := 5 | [ ] x = 5 | [ ] let x = 5"
multiple_choice,"¿Cuáles son tipos válidos en Go?",beginner,go,"int, string y bool son tipos primitivos de Go","[v] int | [v] string | [ ] var | [ ] class"
`

export function buildImportPrompt(topics: Array<Topic>, quantity: number = 5): string {
  const qty = Math.min(50, Math.max(1, Math.floor(quantity) || 5))
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

  return `Eres un generador de preguntas para una plataforma de estudio de desarrollo de software.

Genera EXACTAMENTE ${qty} preguntas en formato CSV de tipo single_choice o multiple_choice. NO uses code_completion.

FORMATO OBLIGATORIO:
- Encabezado: ${IMPORT_CSV_HEADER}
- Delimitador: coma (,) con comillas dobles si el campo contiene comas.
- Codificación UTF-8.

Columnas:
- type: single_choice o multiple_choice
- content: enunciado en ESPAÑOL (claro y específico)
- difficulty: beginner | intermediate | advanced
- topics: slugs separados por ; — USA EXACTAMENTE ESTOS SLUGS, sin traducir ni cambiar mayúsculas:
${topicLines}
  ${topicsNote}
  Ejemplo topics: ${slugsExample}  (si varias: ${slugsExample};otro-slug)
- explanation: explicación educativa opcional en español
- options: opciones separadas por | . Cada opción con prefijo [v] si es CORRECTA o [ ] si es INCORRECTA.

Reglas de options:
- Entre 2 y 5 opciones por pregunta.
- single_choice: EXACTAMENTE 1 opción con [v].
- multiple_choice: 2 o más opciones con [v] permitidas.
- Usa: [v] para correctas, [ ] para incorrectas.

REGLAS IMPORTANTES:
- TODO el contenido debe estar en ESPAÑOL.
- No añadas markdown, comentarios ni texto fuera del CSV.
- Genera exactamente ${qty} filas de datos + la fila de encabezado (total ${qty + 1} filas).
- No repitas preguntas.
- No pongas comillas simples alrededor de [v]/[ ].

Ejemplo de fila:
single_choice,"¿Qué imprime fmt.Println(len(\\"hola\\")) en Go?",beginner,${slugsExample},"len devuelve la longitud de la cadena","[v] 4 | [ ] 5 | [ ] 0 | [ ] error"

Ahora genera el CSV solicitado.`
}
