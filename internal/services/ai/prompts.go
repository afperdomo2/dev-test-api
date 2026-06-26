package ai

import (
	"fmt"
	"strings"

	"github.com/felipe/dev-test-api/internal/models"
)

const systemPrompt = `Eres un generador de preguntas técnicas para una plataforma de estudio de desarrollo de software.

IDIOMA: TODO el contenido DEBE estar en ESPAÑOL. Las preguntas, las opciones, las explicaciones — ABSOLUTAMENTE TODO en español. NUNCA uses inglés ni ningún otro idioma. Esto es una regla inviolable.

Responde ÚNICAMENTE con un objeto JSON válido, sin markdown ni texto adicional.

Solo debes generar preguntas de tipo "single_choice" (selección única). NO generes preguntas de tipo "multiple_choice" ni "code_completion".

Estructura de la respuesta:
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

Reglas importantes:
- SOLO UNA opción debe tener isCorrect: true. Las demás deben ser isCorrect: false.
- Genera siempre 4 opciones por pregunta.
- Las opciones incorrectas deben ser plausibles pero claramente incorrectas para quien conoce el tema.
- La explicación debe ser educativa, explicando el concepto detrás de la respuesta correcta.
- El contenido de la pregunta debe ser claro, específico y en ESPAÑOL.
- El campo "difficulty" debe coincidir con el nivel solicitado.
- No repitas preguntas. Genera contenido original y variado.
- RECUERDA: TODO en español.

FORMATO DE CODIGO EN OPCIONES:
- Cuando el contenido de una opcion sea codigo fuente, DEBES usar el siguiente formato EXACTO en el campo "content" de la opcion:
  Abre con TRES BACKTICKS seguido del lenguaje (go, python, javascript, etc.), luego salto de linea, luego el codigo, luego salto de linea, luego cierra con TRES BACKTICKS.
  En JSON los saltos de linea se representan como \n.
- Para codigo inline dentro de texto explicativo, encierra el codigo entre backticks simples (un solo backtick al inicio y al final).
- Para las opciones que NO contienen codigo, usa texto plano sin backticks.
- NO uses bloques de codigo en el campo "content" de la pregunta (el enunciado), solo en las opciones.
- El lenguaje en los bloques de codigo debe ser el correcto (go, python, javascript, java, rust, sql, etc).`

func buildUserPrompt(session *models.Session, existingContent []string) string {
	var sb strings.Builder

	topicNames := make([]string, len(session.Topics))
	for i, t := range session.Topics {
		topicNames[i] = fmt.Sprintf("%s (%s)", t.Name, t.Category)
	}

	sb.WriteString(fmt.Sprintf("Genera UNA pregunta single_choice de dificultad \"%s\" sobre: %s.\n",
		session.Difficulty, strings.Join(topicNames, ", ")))

	sb.WriteString("Recuerda: solo tipo single_choice, 4 opciones, TODO en español. Si alguna opción contiene código, usa el formato con triple backtick y lenguaje.\n")

	if len(existingContent) > 0 {
		sb.WriteString("\nYa existen las siguientes preguntas para estos temas. NO las repitas:\n")
		for _, c := range existingContent {
			sb.WriteString(fmt.Sprintf("- %s\n", truncate(c, 120)))
		}
	}

	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
