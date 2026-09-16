<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { sessionDetailReviewOptions } from '@/queries/sessions.queries'
import { formatDate, formatScore } from '@/utils/format'
import type { SessionAnswerDetail } from '@/types/session.types'
import CodeContent from '@/components/CodeContent.vue'

const route = useRoute()
const sessionId = computed(() => route.params.id as string)

const { data: sessionData, isLoading } = useQuery(sessionDetailReviewOptions(() => sessionId.value))

const session = computed(() => sessionData.value?.session ?? null)
const answers = computed<Array<SessionAnswerDetail>>(() => sessionData.value?.answers ?? [])

const correctCount = computed(() => answers.value.filter((a) => a.isCorrect).length)
const totalCount = computed(() => answers.value.length)
const scoreValue = computed(() => {
  if (totalCount.value === 0) return 0
  return Math.round((correctCount.value / totalCount.value) * 100)
})

function isSelected(answer: SessionAnswerDetail, optionId: string): boolean {
  return answer.selectedOptions.includes(optionId)
}

function optionClass(
  answer: SessionAnswerDetail,
  opt: { id: string; isCorrect: boolean },
): Record<string, boolean> {
  const selected = isSelected(answer, opt.id)
  return {
    'bg-success-lighten-5': opt.isCorrect,
    'bg-error-lighten-5': selected && !opt.isCorrect,
  }
}

function optionIcon(answer: SessionAnswerDetail, opt: { id: string; isCorrect: boolean }): string {
  const selected = isSelected(answer, opt.id)
  if (opt.isCorrect) return 'mdi-check-circle'
  if (selected) return 'mdi-close-circle'
  return 'mdi-circle-outline'
}

function optionIconColor(
  answer: SessionAnswerDetail,
  opt: { id: string; isCorrect: boolean },
): string {
  const selected = isSelected(answer, opt.id)
  if (opt.isCorrect) return 'success'
  if (selected) return 'error'
  return 'grey-lighten-1'
}
</script>

<template>
  <v-container>
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" to="/sessions" class="mr-2" />
      <div>
        <h1 class="text-h4">{{ session?.name ?? 'Cargando...' }}</h1>
        <div v-if="session" class="text-body-2 text-medium-emphasis mt-1">Sesión completada</div>
      </div>
    </div>

    <v-skeleton-loader v-if="isLoading" type="card" />

    <template v-else-if="session">
      <v-card class="mb-4">
        <v-card-text class="text-center py-6">
          <v-icon size="64" color="success" class="mb-4">mdi-check-circle</v-icon>
          <h2 class="text-h5 mb-2">Sesión completada</h2>
          <div class="d-flex justify-center ga-6 flex-wrap mt-4">
            <div class="text-center">
              <div class="text-h4 font-weight-bold">{{ formatScore(scoreValue) }}</div>
              <div class="text-caption text-medium-emphasis">Puntaje</div>
            </div>
            <div class="text-center">
              <div class="text-h4 font-weight-bold">{{ correctCount }}/{{ totalCount }}</div>
              <div class="text-caption text-medium-emphasis">Correctas</div>
            </div>
            <div class="text-center">
              <div class="text-h4 font-weight-bold">{{ totalCount }}</div>
              <div class="text-caption text-medium-emphasis">Preguntas</div>
            </div>
          </div>
          <p class="text-caption text-medium-emphasis mt-4">
            Finalizada el {{ session.finishedAt ? formatDate(session.finishedAt) : '—' }}
          </p>
        </v-card-text>
      </v-card>

      <div class="d-flex align-center mb-3">
        <h3 class="text-h6">Respuestas</h3>
        <v-spacer />
        <v-chip
          v-if="session.score != null"
          size="small"
          :color="(session.score ?? 0) >= 60 ? 'success' : 'error'"
          variant="tonal"
        >
          {{ formatScore(session.score) }}
        </v-chip>
      </div>

      <v-card v-for="(answer, idx) in answers" :key="answer.id" class="mb-3">
        <v-card-item>
          <template #prepend>
            <v-icon :color="answer.isCorrect ? 'success' : 'error'" size="28">
              {{ answer.isCorrect ? 'mdi-check-circle' : 'mdi-close-circle' }}
            </v-icon>
          </template>
          <v-card-title class="text-body-2 font-weight-medium">
            Pregunta {{ idx + 1 }}
          </v-card-title>
          <v-card-subtitle>
            <v-chip size="x-small" :color="answer.isCorrect ? 'success' : 'error'" variant="tonal">
              {{ answer.isCorrect ? 'Correcta' : 'Incorrecta' }}
            </v-chip>
          </v-card-subtitle>
        </v-card-item>

        <v-card-text>
          <div class="text-body-1 mb-3">
            <CodeContent :text="answer.question.content" />
          </div>

          <v-list v-if="answer.question.options?.length" density="compact">
            <v-list-item
              v-for="opt in answer.question.options"
              :key="opt.id"
              :class="optionClass(answer, opt)"
              rounded
              class="mb-1"
            >
              <template #prepend>
                <v-icon :color="optionIconColor(answer, opt)" size="small">
                  {{ optionIcon(answer, opt) }}
                </v-icon>
              </template>
              <v-list-item-title>
                <CodeContent :text="opt.content" />
              </v-list-item-title>
            </v-list-item>
          </v-list>

          <div v-if="answer.answerText" class="mt-2">
            <span class="text-caption text-medium-emphasis">Tu respuesta:</span>
            <pre
              class="bg-grey-lighten-4 rounded pa-3 mt-1 overflow-auto"
            ><code>{{ answer.answerText }}</code></pre>
          </div>

          <v-card v-if="answer.explanation" variant="tonal" color="info" class="mt-3">
            <v-card-text class="text-body-2">
              <span class="font-weight-medium">Explicación:</span>
              <CodeContent :text="answer.explanation" />
            </v-card-text>
          </v-card>
        </v-card-text>
      </v-card>

      <div class="text-center py-4">
        <v-btn color="primary" to="/sessions">Volver a sesiones</v-btn>
      </div>
    </template>
  </v-container>
</template>
