<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { questionDetailOptions } from '@/queries/questions.queries'
import {
  DIFFICULTY_COLORS,
  DIFFICULTY_LABELS,
  TYPE_ICONS,
  SOURCE_LABELS,
  SOURCE_COLORS,
} from '@/types/question.types'
import { formatDateTime } from '@/utils/format'
import CodeContent from '@/components/CodeContent.vue'

const route = useRoute()
const questionId = computed(() => route.params.id as string)

const { data: question, isLoading } = useQuery(questionDetailOptions(() => questionId.value))

function typeLabel(type: string): string {
  const labels: Record<string, string> = {
    single_choice: 'Selección única',
    multiple_choice: 'Selección múltiple',
    code_completion: 'Completar código',
  }
  return labels[type] ?? type
}
</script>

<template>
  <v-container>
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" to="/questions" class="mr-2" />
      <h1 class="text-h4">Detalle de pregunta</h1>
    </div>

    <v-skeleton-loader v-if="isLoading" type="article" />

    <template v-else-if="question">
      <!-- Header -->
      <v-card class="mb-4">
        <v-card-item>
          <template #prepend>
            <v-icon :icon="TYPE_ICONS[question.type]" color="primary" size="32" />
          </template>
          <v-card-title class="text-wrap">
            <CodeContent :text="question.content" />
          </v-card-title>
          <v-card-subtitle>
            <v-chip
              :color="DIFFICULTY_COLORS[question.difficulty]"
              size="small"
              variant="tonal"
              class="mr-2"
            >
              {{ DIFFICULTY_LABELS[question.difficulty] }}
            </v-chip>
            <v-chip size="small" variant="tonal" class="mr-2">
              {{ typeLabel(question.type) }}
            </v-chip>
            <v-chip
              :color="SOURCE_COLORS[question.source]"
              size="small"
              variant="tonal"
              class="mr-2"
            >
              {{ SOURCE_LABELS[question.source] }}
            </v-chip>
            <v-chip v-if="question.isPublic" size="small" variant="text" color="success">
              Pública
            </v-chip>
          </v-card-subtitle>
        </v-card-item>
      </v-card>

      <!-- Options -->
      <v-card v-if="question.options?.length" class="mb-4">
        <v-card-title class="text-h6">Opciones</v-card-title>
        <v-card-text>
          <v-list density="compact">
            <v-list-item v-for="option in question.options ?? []" :key="option.id">
              <template #prepend>
                <v-icon v-if="option.isCorrect" color="success" size="small">
                  mdi-check-circle
                </v-icon>
                <v-icon v-else color="grey" size="small"> mdi-circle-outline </v-icon>
              </template>
              <v-list-item-title>
                <CodeContent :text="option.content" />
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>

      <!-- Code Challenge -->
      <v-card v-if="question.codeChallenge" class="mb-4">
        <v-card-title class="text-h6">Código</v-card-title>
        <v-card-text>
          <v-chip size="small" variant="tonal" class="mb-2">
            {{ question.codeChallenge.language }}
          </v-chip>
          <pre class="rounded pa-4 overflow-auto bg-grey-lighten-4"><code
            v-highlight="question.codeChallenge.language"
          >{{ question.codeChallenge.starterCode }}</code></pre>
        </v-card-text>
      </v-card>

      <!-- Explanation -->
      <v-card v-if="question.explanation" class="mb-4">
        <v-card-title class="text-h6">Explicación</v-card-title>
        <v-card-text>
          <p class="text-body-1">
            <CodeContent :text="question.explanation" />
          </p>
        </v-card-text>
      </v-card>

      <!-- Topics -->
      <v-card v-if="question.topics.length" class="mb-4">
        <v-card-title class="text-h6">Temas</v-card-title>
        <v-card-text>
          <div class="d-flex flex-wrap ga-1">
            <v-chip
              v-for="topic in question.topics"
              :key="topic"
              size="small"
              variant="outlined"
              color="primary"
            >
              {{ topic }}
            </v-chip>
          </div>
        </v-card-text>
      </v-card>

      <!-- Metadata footer -->
      <v-card variant="text">
        <v-card-text class="text-caption text-medium-emphasis">
          Creado: {{ formatDateTime(question.createdAt) }}
        </v-card-text>
      </v-card>
    </template>

    <v-card v-else>
      <v-card-text class="text-center py-8">
        <v-icon size="48" color="error" class="mb-2"> mdi-alert-circle </v-icon>
        <p class="text-body-1">Pregunta no encontrada</p>
        <v-btn color="primary" to="/questions"> Volver al listado </v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>
