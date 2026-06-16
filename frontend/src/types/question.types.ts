export type QuestionType = 'single_choice' | 'multiple_choice' | 'code_completion'
export type QuestionDifficulty = 'beginner' | 'intermediate' | 'advanced'
export type QuestionSource = 'ai_generated' | 'manual' | 'imported'

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
