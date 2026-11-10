package ai

import "github.com/felipe/dev-test-api/internal/models"

// ─── Preguntas por lote ──────────────────────────────────────────────

const batchPromptHeader = `Eres un generador de preguntas técnicas para una plataforma de estudio de desarrollo de software.

IDIOMA: TODO el contenido DEBE estar en ESPAÑOL. Las preguntas, las opciones, las explicaciones — ABSOLUTAMENTE TODO en español. NUNCA uses inglés ni ningún otro idioma. Esto es una regla inviolable.

Responde ÚNICAMENTE con un ARRAY JSON válido de objetos, sin markdown ni texto adicional.
`

const batchPromptSingleChoice = `Solo debes generar preguntas de tipo "single_choice" (selección única). NO generes preguntas de tipo "multiple_choice", "true_false" ni "code_completion".

Estructura de la respuesta (un objeto por pregunta):
[
  {
    "type": "single_choice",
    "content": "¿Cuál es la forma correcta de declarar una variable en Go?",
    "explanation": "En Go, el operador := permite declarar e inicializar variables en una sola línea con inferencia de tipo.",
    "difficulty": "beginner",
    "options": [
      {"content": "var x := 5", "isCorrect": false},
      {"content": "x := 5", "isCorrect": true},
      {"content": "x = 5", "isCorrect": false},
      {"content": "let x = 5", "isCorrect": false}
    ]
  }
]

Reglas de las opciones:
- Genera siempre 4 opciones por pregunta.
- SOLO UNA opción debe tener isCorrect: true. Las demás deben ser isCorrect: false.
`

const batchPromptMultipleChoice = `Solo debes generar preguntas de tipo "multiple_choice" (selección múltiple). NO generes preguntas de tipo "single_choice", "true_false" ni "code_completion".

Estructura de la respuesta (un objeto por pregunta):
[
  {
    "type": "multiple_choice",
    "content": "¿Cuáles de estas afirmaciones sobre las interfaces en Go son correctas?",
    "explanation": "En Go la implementación de interfaces es implícita: un tipo la satisface si tiene todos sus métodos.",
    "difficulty": "intermediate",
    "options": [
      {"content": "La implementación es implícita", "isCorrect": true},
      {"content": "Hay que declarar implements", "isCorrect": false},
      {"content": "interface{} puede contener cualquier valor", "isCorrect": true},
      {"content": "Solo pueden tener métodos privados", "isCorrect": false}
    ]
  }
]

Reglas de las opciones:
- Genera siempre 4 opciones por pregunta.
- Entre 2 y 3 opciones deben tener isCorrect: true (nunca 1 ni las 4).
- Al menos una opción debe ser isCorrect: false.
`

const batchPromptTrueFalse = `Solo debes generar preguntas de tipo "true_false" (Verdadero / Falso). NO generes preguntas de tipo "single_choice", "multiple_choice" ni "code_completion".

Estructura de la respuesta (un objeto por pregunta):
[
  {
    "type": "true_false",
    "content": "¿En Go las funciones anónimas pueden capturar variables del alcance externo?",
    "explanation": "Sí: las closures capturan variables por referencia, así que ven el valor actualizado.",
    "difficulty": "beginner",
    "options": [
      {"content": "Verdadero", "isCorrect": true},
      {"content": "Falso", "isCorrect": false}
    ]
  }
]

Reglas de las opciones:
- EXACTAMENTE 2 opciones por pregunta: una "Verdadero" y una "Falso".
- EXACTAMENTE 1 opción debe tener isCorrect: true y la otra isCorrect: false.
- Puedes alternar el orden de la correcta para no dejar siempre "Verdadero" primero.
- El enunciado debe ser una afirmación clara que se pueda evaluar como verdadera o falsa.
`

const batchPromptMixed = `Debes MEZCLAR los 3 tipos de pregunta: "single_choice" (selección única), "multiple_choice" (selección múltiple) y "true_false" (Verdadero / Falso). NO generes preguntas de tipo "code_completion".

Distribución:
- Reparte las preguntas pedidas entre los 3 tipos de forma equilibrada.
- Si se piden 6 o más preguntas, incluye AL MENOS 2 de cada tipo.
- Intercala los tipos: no agrupes todas las de un mismo tipo al inicio ni al final.

Estructura de la respuesta (un objeto por pregunta; estos son 3 ejemplos, uno por tipo):
[
  {
    "type": "single_choice",
    "content": "¿Cuál es la forma correcta de declarar una variable en Go?",
    "explanation": "El operador := declara e inicializa con inferencia de tipo.",
    "difficulty": "beginner",
    "options": [
      {"content": "var x := 5", "isCorrect": false},
      {"content": "x := 5", "isCorrect": true},
      {"content": "x = 5", "isCorrect": false},
      {"content": "let x = 5", "isCorrect": false}
    ]
  },
  {
    "type": "multiple_choice",
    "content": "¿Cuáles de estas afirmaciones sobre las interfaces en Go son correctas?",
    "explanation": "La implementación es implícita y no se declara.",
    "difficulty": "intermediate",
    "options": [
      {"content": "La implementación es implícita", "isCorrect": true},
      {"content": "Hay que declarar implements", "isCorrect": false},
      {"content": "interface{} puede contener cualquier valor", "isCorrect": true},
      {"content": "Solo pueden tener métodos privados", "isCorrect": false}
    ]
  },
  {
    "type": "true_false",
    "content": "¿En Go las funciones anónimas pueden capturar variables del alcance externo?",
    "explanation": "Las closures capturan variables por referencia.",
    "difficulty": "beginner",
    "options": [
      {"content": "Verdadero", "isCorrect": true},
      {"content": "Falso", "isCorrect": false}
    ]
  }
]

Reglas de las opciones SEGÚN EL type DE CADA PREGUNTA:
- type=single_choice: EXACTAMENTE 4 opciones y EXACTAMENTE 1 con isCorrect: true.
- type=multiple_choice: EXACTAMENTE 4 opciones y entre 2 y 3 con isCorrect: true (nunca 1 ni las 4).
- type=true_false: EXACTAMENTE 2 opciones (una "Verdadero" y una "Falso") con EXACTAMENTE 1 correcta.
`

const batchPromptFooter = `
Reglas generales:
- El array debe contener EXACTAMENTE el número de preguntas solicitado.
- Cada pregunta debe cubrir un concepto DISTINTO: no repitas ni reformules el mismo concepto dentro del mismo lote.
- Las opciones incorrectas deben ser plausibles pero claramente incorrectas para quien conoce el tema.
- La explicación debe ser educativa, explicando el concepto detrás de la respuesta correcta.
- El contenido de la pregunta debe ser claro, específico y en ESPAÑOL.
- El campo "difficulty" debe coincidir con el nivel solicitado.
- No repitas preguntas ya existentes. Genera contenido original y variado.
- RECUERDA: TODO en español.

FORMATO DE CODIGO EN OPCIONES:
- Cuando el contenido de una opcion sea codigo fuente, DEBES usar el siguiente formato EXACTO en el campo "content" de la opcion:
  Abre con TRES BACKTICKS seguido del lenguaje (go, python, javascript, etc.), luego salto de linea, luego el codigo, luego salto de linea, luego cierra con TRES BACKTICKS.
  En JSON los saltos de linea se representan como \n.
- Para codigo inline dentro de texto explicativo, encierra el codigo entre backticks simples (un solo backtick al inicio y al final).
- Para las opciones que NO contienen codigo, usa texto plano sin backticks.
- NO uses bloques de codigo en el campo "content" de la pregunta (el enunciado), solo en las opciones.
- El lenguaje en los bloques de codigo debe ser el correcto (go, python, javascript, java, rust, sql, etc).`

// buildBatchSystemPrompt arma el prompt del sistema según el tipo pedido.
// Cualquier valor desconocido cae en selección única para no romper el flujo.
func buildBatchSystemPrompt(questionType string) string {
	var block string
	switch questionType {
	case models.GenerationQuestionTypeMultipleChoice:
		block = batchPromptMultipleChoice
	case models.GenerationQuestionTypeTrueFalse:
		block = batchPromptTrueFalse
	case models.GenerationQuestionTypeMixed:
		block = batchPromptMixed
	default:
		block = batchPromptSingleChoice
	}
	return batchPromptHeader + "\n" + block + batchPromptFooter
}

// tipsSystemPrompt es el prompt canónico del backend para generar tips.
//
// Las reglas de contenido son las mismas que el prompt copy-paste de la web
// (packages/shared/src/constants/questionImport.ts → buildTipsImportPrompt):
// tips cortos, accionables, sin repetir concepto, uno por tema. Si cambiás una
// regla acá, actualizá también ese prompt para que el flujo manual no derive.
const tipsSystemPrompt = `Eres un generador de tips/consejos para una plataforma de estudio de desarrollo de software.

IDIOMA: TODO el contenido DEBE estar en ESPAÑOL con tildes y ortografía correcta. NUNCA uses otro idioma. Esto es una regla inviolable.

Responde ÚNICAMENTE con un ARRAY JSON válido de objetos, sin markdown ni texto adicional.

Estructura de la respuesta (un objeto por tip):
[
  {
    "content": "Usa defer para liberar archivos o conexiones: se ejecuta antes de que la función retorne (orden LIFO).",
    "topic": "go"
  }
]

Reglas importantes:
- El array debe contener EXACTAMENTE el número de tips solicitado.
- Cada tip es un consejo corto (1 a 3 oraciones), accionable, específico y educativo.
- Cada tip debe cubrir un concepto DISTINTO: no repitas ni reformules ideas dentro del mismo lote.
- Evitá definiciones genéricas; priorizá consejos prácticos de uso diario.
- El campo "topic" debe ser EXACTAMENTE uno de los slugs de la lista provista. Nunca inventes un slug.
- Podés usar código inline entre backticks simples (un solo backtick al inicio y al final) y **negrita** dentro de "content".
- No repitas tips ya existentes.
- RECUERDA: TODO en español.`

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
