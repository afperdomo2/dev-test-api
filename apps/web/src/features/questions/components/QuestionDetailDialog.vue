<script setup lang="ts">
import { ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { toggleSaveMutation } from '@/queries/progress.queries'
import { useAppStore } from '@/stores/app.store'
import {
  DIFFICULTY_COLORS,
  DIFFICULTY_LABELS,
  TYPE_ICONS,
  SOURCE_LABELS,
  SOURCE_COLORS,
} from '@/types/question.types'
import { formatDateTime, formatDate } from '@/utils/format'
import CodeContent from '@/components/CodeContent.vue'
import type { Question } from '@/types/question.types'
import type { Progress } from '@/types/progress.types'

interface Props {
  modelValue: boolean
  question: Question | null
  progress?: Progress | null
  loading?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()

const appStore = useAppStore()
const queryClient = useQueryClient()
const toggleMut = useMutation(toggleSaveMutation())
const toggling = ref(false)

const isSaved = ref(false)
watch(
  () => props.progress?.isSaved,
  (val) => {
    isSaved.value = val ?? false
  },
  { immediate: true },
)

async function handleToggleSave() {
  if (!props.question) return
  toggling.value = true
  try {
    const updated = await toggleMut.mutateAsync(props.question.id)
    isSaved.value = updated.isSaved
    await queryClient.invalidateQueries({ queryKey: ['progress'] })
    appStore.showSnackbar(updated.isSaved ? 'Pregunta guardada' : 'Pregunta desguardada')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al guardar'
    appStore.showSnackbar(detail, 'error')
  } finally {
    toggling.value = false
  }
}

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
  <v-dialog
    :model-value="modelValue"
    max-width="800"
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center py-3">
        <v-icon v-if="question" :icon="TYPE_ICONS[question.type]" color="primary" class="mr-2" />
        <span class="text-body-1 flex-grow-1">Detalle de pregunta</span>
        <v-btn
          icon="mdi-close"
          variant="text"
          size="small"
          @click="emit('update:modelValue', false)"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="overflow-y-auto" style="max-height: 70vh">
        <v-skeleton-loader v-if="loading && !question" type="article" />

        <template v-else-if="question">
          <!-- Content -->
          <div class="mb-4">
            <CodeContent :text="question.content" />
          </div>

          <!-- Chips -->
          <div class="d-flex flex-wrap ga-1 mb-4">
            <v-chip :color="DIFFICULTY_COLORS[question.difficulty]" size="small" variant="tonal">
              {{ DIFFICULTY_LABELS[question.difficulty] }}
            </v-chip>
            <v-chip size="small" variant="tonal">
              {{ typeLabel(question.type) }}
            </v-chip>
            <v-chip :color="SOURCE_COLORS[question.source]" size="small" variant="tonal">
              {{ SOURCE_LABELS[question.source] }}
            </v-chip>
            <v-chip v-if="question.isPublic" size="small" variant="text" color="success">
              Pública
            </v-chip>
          </div>

          <!-- Options -->
          <v-card
            v-if="question.options?.length"
            variant="flat"
            bg-color="grey-lighten-5"
            class="mb-4"
          >
            <v-card-title class="text-subtitle-1">Opciones</v-card-title>
            <v-card-text class="pt-0">
              <v-list density="compact">
                <v-list-item v-for="option in question.options" :key="option.id">
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
          <v-card
            v-if="question.codeChallenge"
            variant="flat"
            bg-color="grey-lighten-5"
            class="mb-4"
          >
            <v-card-title class="text-subtitle-1">Código</v-card-title>
            <v-card-text class="pt-0">
              <v-chip size="small" variant="tonal" class="mb-2">
                {{ question.codeChallenge.language }}
              </v-chip>
              <pre
                class="rounded pa-4 overflow-auto bg-grey-lighten-4"
              ><code v-highlight="question.codeChallenge.language">{{ question.codeChallenge.starterCode }}</code></pre>
            </v-card-text>
          </v-card>

          <!-- Explanation -->
          <v-card v-if="question.explanation" variant="flat" bg-color="grey-lighten-5" class="mb-4">
            <v-card-title class="text-subtitle-1">Explicación</v-card-title>
            <v-card-text class="pt-0">
              <CodeContent :text="question.explanation" />
            </v-card-text>
          </v-card>

          <!-- Topics -->
          <div v-if="question.topics.length" class="mb-4">
            <p class="text-subtitle-2 mb-1">Temas</p>
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
          </div>

          <!-- Progress info -->
          <v-card v-if="progress" variant="flat" bg-color="grey-lighten-5" class="mb-4">
            <v-card-title class="text-subtitle-1">Progreso SM-2</v-card-title>
            <v-card-text class="pt-0">
              <v-row dense>
                <v-col cols="6" sm="3">
                  <p class="text-caption text-medium-emphasis">Repeticiones</p>
                  <p class="text-body-2">{{ progress.repetitions }}</p>
                </v-col>
                <v-col cols="6" sm="3">
                  <p class="text-caption text-medium-emphasis">Intervalo</p>
                  <p class="text-body-2">{{ progress.intervalDays }} días</p>
                </v-col>
                <v-col cols="6" sm="3">
                  <p class="text-caption text-medium-emphasis">Factor ease</p>
                  <p class="text-body-2">{{ progress.easeFactor.toFixed(1) }}</p>
                </v-col>
                <v-col cols="6" sm="3">
                  <p class="text-caption text-medium-emphasis">Estado</p>
                  <v-chip v-if="progress.isMastered" size="x-small" color="success" variant="tonal">
                    Dominada
                  </v-chip>
                  <v-chip v-else size="x-small" variant="tonal"> En progreso </v-chip>
                </v-col>
              </v-row>
              <div v-if="progress.nextReviewAt" class="mt-2 text-caption text-medium-emphasis">
                Próximo repaso: {{ formatDate(progress.nextReviewAt) }}
              </div>
            </v-card-text>
          </v-card>

          <!-- Metadata -->
          <div class="text-caption text-medium-emphasis">
            Creado: {{ formatDateTime(question.createdAt) }}
          </div>
        </template>
      </v-card-text>

      <v-divider />

      <v-card-actions class="pa-4">
        <v-spacer />
        <v-btn
          :color="isSaved ? 'warning' : 'primary'"
          :loading="toggling"
          variant="tonal"
          @click="handleToggleSave"
        >
          <v-icon :icon="isSaved ? 'mdi-bookmark' : 'mdi-bookmark-outline'" start />
          {{ isSaved ? 'Guardada' : 'Guardar' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
