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

export const DIFFICULTY_LABELS: Record<QuestionDifficulty, string> = {
  beginner: 'Principiante',
  intermediate: 'Intermedio',
  advanced: 'Avanzado',
}

export const SOURCE_LABELS: Record<QuestionSource, string> = {
  ai_generated: 'IA',
  manual: 'Manual',
  imported: 'Importada',
}

export const SOURCE_COLORS: Record<QuestionSource, string> = {
  ai_generated: 'purple',
  manual: 'blue',
  imported: 'orange',
}

export const SOURCE_ICONS: Record<QuestionSource, string> = {
  ai_generated: 'mdi-robot',
  manual: 'mdi-pencil',
  imported: 'mdi-file-excel',
}

export const TYPE_ICONS: Record<QuestionType, string> = {
  single_choice: 'mdi-radiobox-marked',
  multiple_choice: 'mdi-checkbox-multiple-marked',
  code_completion: 'mdi-code-braces',
}

export interface QuestionOption {
  id: string
  content: string
  isCorrect?: boolean
}

export interface CodeChallenge {
  id: string
  starterCode: string
  expectedOutput: string
  language: string
  testCases: string
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
  options?: Array<QuestionOption>
  codeChallenge: CodeChallenge | null
  topics: Array<string>
  createdAt: string
  updatedAt: string
}

export interface CreateQuestionOption {
  content: string
  isCorrect: boolean
}

export interface CreateQuestionRequest {
  type: QuestionType
  content: string
  difficulty: QuestionDifficulty
  explanation?: string
  topicIds: Array<string>
  options?: Array<CreateQuestionOption>
  starterCode?: string
  expectedOutput?: string
  language?: string
  testCases?: string
}

export interface UpdateQuestionRequest {
  type?: QuestionType
  content?: string
  difficulty?: QuestionDifficulty
  explanation?: string
  topicIds?: Array<string>
  options?: Array<CreateQuestionOption>
  starterCode?: string
  expectedOutput?: string
  language?: string
  testCases?: string
}

export interface ImportRowError {
  row: number
  reason: string
}

export interface ImportResult {
  imported: number
  failed: number
  total: number
  errors: Array<ImportRowError>
}

export interface ImportQuota {
  dailyLimit: number
  usedToday: number
  remaining: number
}

export interface AiQuota {
  dailyLimit: number
  usedToday: number
  remaining: number
}

export const MAX_QUESTIONS_PER_FILE = 50

export interface TopicCount {
  topicId: string
  name: string
  slug: string
  category: string
  count: number
}

export interface CategoryCount {
  category: string
  count: number
}

export interface DifficultyCount {
  difficulty: string
  count: number
}

export interface TypeCount {
  type: string
  count: number
}

export interface QuestionStats {
  total: number
  byTopic: Array<TopicCount>
  byCategory: Array<CategoryCount>
  byDifficulty: Array<DifficultyCount>
  byType: Array<TypeCount>
}
