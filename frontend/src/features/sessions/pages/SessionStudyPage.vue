<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import {
  sessionDetailOptions,
  nextQuestionOptions,
  sessionSummaryOptions,
  submitAnswerMutation,
  finishSessionMutation,
} from '@/queries/sessions.queries'
import { useAppStore } from '@/stores/app.store'
import { DIFFICULTY_COLORS, TYPE_ICONS } from '@/types/question.types'
import type { SessionAnswer } from '@/types/session.types'
import CodeContent from '@/components/CodeContent.vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const queryClient = useQueryClient()

const sessionId = computed(() => route.params.id as string)

const { data: session, isLoading: sessionLoading } = useQuery(
  sessionDetailOptions(() => sessionId.value),
)

const { data: summary } = useQuery(sessionSummaryOptions(() => sessionId.value))

const currentQuestionNumber = computed(() => {
  const count = summary.value?.answerCount ?? 0
  const limit = summary.value?.questionLimit
  if (limit !== null && count >= limit) return limit
  return count + 1
})

const questionLimit = computed(() => summary.value?.questionLimit ?? null)

const {
  data: currentQuestion,
  isLoading: questionLoading,
  isFetching: questionFetching,
  refetch: refetchNext,
} = useQuery(nextQuestionOptions(() => sessionId.value))

const submitMut = useMutation(submitAnswerMutation())
const finishMut = useMutation(finishSessionMutation())

type AnswerState = 'answering' | 'submitting' | 'result' | 'done'
const answerState = ref<AnswerState>('answering')
const selectedOptions = ref<Array<string>>([])
const answerText = ref('')
const lastResult = ref<SessionAnswer | null>(null)

// Reset answer when question changes
watch(currentQuestion, () => {
  selectedOptions.value = []
  answerText.value = ''
  lastResult.value = null
  answerState.value = 'answering'
})

const isSingleChoice = computed(() => currentQuestion.value?.type === 'single_choice')
const isMultipleChoice = computed(() => currentQuestion.value?.type === 'multiple_choice')
const isCodeCompletion = computed(() => currentQuestion.value?.type === 'code_completion')

const canSubmit = computed(() => {
  if (isSingleChoice.value) return selectedOptions.value.length === 1
  if (isMultipleChoice.value) return selectedOptions.value.length > 0
  if (isCodeCompletion.value) return answerText.value.trim().length > 0
  return false
})

const optionsArray = computed<Array<{ id: string; content: string }>>(() => {
  if (!currentQuestion.value) return []
  if (isCodeCompletion.value) return []
  return currentQuestion.value.options ?? []
})

function toggleOption(optionId: string) {
  if (isSingleChoice.value) {
    selectedOptions.value = [optionId]
    return
  }
  const idx = selectedOptions.value.indexOf(optionId)
  if (idx >= 0) {
    selectedOptions.value.splice(idx, 1)
  } else {
    selectedOptions.value.push(optionId)
  }
}

function isSelected(optionId: string): boolean {
  return selectedOptions.value.includes(optionId)
}

function isCorrectOption(optionId: string): boolean {
  if (!currentQuestion.value) return false
  return !!currentQuestion.value.options?.find((o) => o.id === optionId)?.isCorrect
}

async function submitAnswer() {
  if (!canSubmit.value || !currentQuestion.value) return
  answerState.value = 'submitting'
  try {
    const result = await submitMut.mutateAsync({
      sessionId: sessionId.value,
      data: {
        questionId: currentQuestion.value.id,
        selectedOptions: isCodeCompletion.value ? undefined : selectedOptions.value,
        answerText: isCodeCompletion.value ? answerText.value : undefined,
      },
    })
    lastResult.value = result
    answerState.value = 'result'
    queryClient.invalidateQueries({ queryKey: ['sessions', 'summary', sessionId] })
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al enviar respuesta'
    appStore.showSnackbar(detail, 'error')
    answerState.value = 'answering'
  }
}

async function nextQuestion() {
  lastResult.value = null
  const result = await refetchNext()
  if (!result.data) return
}

async function finishSession() {
  try {
    await finishMut.mutateAsync(sessionId.value)
    appStore.showSnackbar('Sesión finalizada')
    queryClient.invalidateQueries({ queryKey: ['sessions', 'list'] })
    router.push('/sessions')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al finalizar sesión'
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-container>
    <!-- Session header -->
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" to="/sessions" class="mr-2" />
      <div>
        <h1 class="text-h4">{{ session?.name ?? 'Cargando...' }}</h1>
        <div v-if="questionLimit" class="text-body-2 text-medium-emphasis mt-1">
          Pregunta {{ currentQuestionNumber }} de {{ questionLimit }}
        </div>
      </div>
    </div>

    <v-skeleton-loader v-if="sessionLoading" type="article" />

    <template v-else-if="session">
      <!-- Question area -->
      <v-skeleton-loader v-if="questionLoading || questionFetching" type="card" class="mb-4" />

      <template v-else-if="currentQuestion">
        <!-- Question card -->
        <v-card class="mb-4">
          <v-card-item>
            <template #prepend>
              <v-icon :icon="TYPE_ICONS[currentQuestion.type]" color="primary" size="28" />
            </template>
            <v-card-title class="text-wrap text-body-1">
              <CodeContent :text="currentQuestion.content" />
            </v-card-title>
            <v-card-subtitle>
              <v-chip
                :color="DIFFICULTY_COLORS[currentQuestion.difficulty]"
                size="x-small"
                variant="tonal"
                class="mr-1"
              >
                {{ currentQuestion.difficulty }}
              </v-chip>
              <v-chip
                v-for="topic in currentQuestion.topics"
                :key="topic.id"
                size="x-small"
                variant="outlined"
                class="mr-1"
              >
                {{ topic.name }}
              </v-chip>
            </v-card-subtitle>
          </v-card-item>
        </v-card>

        <!-- Options (single/multiple choice) -->
        <v-card v-if="!isCodeCompletion" class="mb-4">
          <v-card-text>
            <v-list density="compact">
              <v-list-item
                v-for="option in optionsArray"
                :key="option.id"
                :class="{
                  'bg-success-lighten-5': answerState === 'result' && isCorrectOption(option.id),
                  'bg-error-lighten-5':
                    answerState === 'result' &&
                    isSelected(option.id) &&
                    !isCorrectOption(option.id),
                }"
                rounded
                class="mb-1"
                :disabled="answerState === 'submitting' || answerState === 'result'"
                @click="toggleOption(option.id)"
              >
                <template #prepend>
                  <v-checkbox
                    v-if="isMultipleChoice"
                    :model-value="isSelected(option.id)"
                    :disabled="answerState === 'submitting' || answerState === 'result'"
                    hide-details
                    density="compact"
                  />
                  <v-radio
                    v-else
                    :model-value="isSelected(option.id)"
                    :disabled="answerState === 'submitting' || answerState === 'result'"
                    hide-details
                    density="compact"
                  />
                </template>
                <v-list-item-title>
                  <CodeContent :text="option.content" />
                </v-list-item-title>
                <template v-if="answerState === 'result'" #append>
                  <v-icon v-if="isCorrectOption(option.id)" color="success" size="small">
                    mdi-check-circle
                  </v-icon>
                  <v-icon v-else-if="isSelected(option.id)" color="error" size="small">
                    mdi-close-circle
                  </v-icon>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>

        <!-- Code completion -->
        <v-card v-if="isCodeCompletion" class="mb-4">
          <v-card-text>
            <div
              v-if="currentQuestion.codeChallenge"
              class="text-caption text-medium-emphasis mb-2"
            >
              Código inicial:
            </div>
            <pre
              v-if="currentQuestion.codeChallenge"
              class="rounded pa-4 overflow-auto bg-grey-lighten-4 mb-4"
            ><code
              v-highlight="currentQuestion.codeChallenge.language"
            >{{ currentQuestion.codeChallenge.starterCode }}</code></pre>
            <v-textarea
              v-model="answerText"
              label="Tu código"
              :disabled="answerState === 'submitting' || answerState === 'result'"
              rows="6"
              auto-grow
              variant="outlined"
              class="font-monospace"
            />
          </v-card-text>
        </v-card>

        <!-- Result feedback -->
        <v-card v-if="answerState === 'result' && lastResult" class="mb-4">
          <v-card-item>
            <template #prepend>
              <v-icon :color="lastResult.isCorrect ? 'success' : 'error'" size="32">
                {{ lastResult.isCorrect ? 'mdi-check-circle' : 'mdi-close-circle' }}
              </v-icon>
            </template>
            <v-card-title>
              {{ lastResult.isCorrect ? '¡Correcto!' : 'Incorrecto' }}
            </v-card-title>
            <v-card-subtitle v-if="lastResult.aiFeedback">
              <CodeContent :text="lastResult.aiFeedback" />
            </v-card-subtitle>
          </v-card-item>
          <v-card-actions>
            <v-spacer />
            <v-btn
              v-if="currentQuestionNumber >= (questionLimit ?? Infinity)"
              color="primary"
              @click="finishSession"
            >
              Finalizar sesión
            </v-btn>
            <v-btn v-else color="primary" variant="text" @click="nextQuestion">
              Siguiente pregunta
            </v-btn>
          </v-card-actions>
        </v-card>

        <!-- Action buttons (before answer) -->
        <div
          v-if="answerState === 'answering' || answerState === 'submitting'"
          class="d-flex ga-2 justify-end"
        >
          <v-btn variant="tonal" :disabled="answerState === 'submitting'" @click="finishSession">
            Finalizar sesión
          </v-btn>
          <v-btn
            color="primary"
            :loading="answerState === 'submitting'"
            :disabled="!canSubmit"
            @click="submitAnswer"
          >
            Responder
          </v-btn>
        </div>
      </template>

      <!-- No more questions -->
      <v-card v-else>
        <v-card-text class="text-center py-8">
          <template
            v-if="
              session.mode === 'generate' &&
              (summary?.questionsGenerated ?? 0) < (summary?.questionLimit ?? Infinity)
            "
          >
            <v-progress-circular indeterminate color="primary" size="48" class="mb-4" />
            <h2 class="text-h5 mb-2">Generando preguntas...</h2>
            <p class="text-body-1 text-medium-emphasis mb-4">
              La IA está creando más preguntas para esta sesión. Esto toma unos segundos.
            </p>
            <v-btn
              color="primary"
              variant="tonal"
              :loading="questionLoading"
              @click="refetchNext()"
            >
              Reintentar
            </v-btn>
          </template>
          <template v-else>
            <v-icon size="64" color="success" class="mb-4"> mdi-check-circle </v-icon>
            <h2 class="text-h5 mb-2">¡No hay más preguntas!</h2>
            <p class="text-body-1 text-medium-emphasis mb-4">
              Has completado todas las preguntas disponibles en esta sesión.
            </p>
            <v-btn color="primary" @click="finishSession"> Finalizar sesión </v-btn>
          </template>
        </v-card-text>
      </v-card>
    </template>
  </v-container>
</template>
