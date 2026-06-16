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
  is_correct: boolean
}

export interface CodeChallenge {
  id: string
  starter_code: string
  expected_output: string
  language: string
  test_cases: string
}

export interface QuestionTopic {
  id: string
  slug: string
  name: string
}

export interface Question {
  id: string
  user_id: string
  type: QuestionType
  content: string
  explanation: string
  difficulty: QuestionDifficulty
  is_public: boolean
  source: QuestionSource
  options: Array<QuestionOption>
  code_challenge: CodeChallenge | null
  topics: Array<QuestionTopic>
  created_at: string
  updated_at: string
}
