export type QuestionType = 'single_choice' | 'multiple_choice' | 'code_completion'
export type QuestionDifficulty = 'beginner' | 'intermediate' | 'advanced'
export type QuestionSource = 'ai_generated' | 'manual' | 'imported'

export const QUESTION_TYPES: Array<{ title: string; value: QuestionType }> = [
  { title: 'Selección única', value: 'single_choice' },
  { title: 'Selección múltiple', value: 'multiple_choice' },
  { title: 'Completar código', value: 'code_completion' },
]

export const QUESTION_DIFFICULTIES: Array<{ title: string; value: QuestionDifficulty }> = [
  { title: 'Principiante', value: 'beginner' },
  { title: 'Intermedio', value: 'intermediate' },
  { title: 'Avanzado', value: 'advanced' },
]

export const DIFFICULTY_COLORS: Record<QuestionDifficulty, string> = {
  beginner: 'success',
  intermediate: 'warning',
  advanced: 'error',
}

export const TYPE_ICONS: Record<QuestionType, string> = {
  single_choice: 'mdi-radiobox-marked',
  multiple_choice: 'mdi-checkbox-multiple-marked',
  code_completion: 'mdi-code-braces',
}

export interface QuestionOption {
  id: string
  content: string
  isCorrect: boolean
}

export interface CodeChallenge {
  id: string
  starterCode: string
  expectedOutput: string
  language: string
  testCases: string
}

export interface QuestionTopic {
  id: string
  slug: string
  name: string
}

export interface Question {
  id: string
  userId: string
  type: QuestionType
  content: string
  explanation: string
  difficulty: QuestionDifficulty
  isPublic: boolean
  source: QuestionSource
  options: Array<QuestionOption>
  codeChallenge: CodeChallenge | null
  topics: Array<QuestionTopic>
  createdAt: string
  updatedAt: string
}
