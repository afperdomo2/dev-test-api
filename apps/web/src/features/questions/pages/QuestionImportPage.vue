<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { importQuestions, getImportQuota } from '@/api/services/questions.service'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import TopicAutocomplete from '@/components/TopicAutocomplete.vue'
import { MAX_QUESTIONS_PER_FILE, type ImportResult } from '@/types/question.types'
import {
  IMPORT_CSV_TEMPLATE,
  buildImportPrompt,
  PROMPT_DIFFICULTIES,
  PROMPT_OBJECTIVES,
  type PromptDifficulty,
  type PromptObjective,
} from '@/constants/questionImport'
import { requiredRule, validateRules } from '@/utils/validators'
import type { Topic } from '@/types/topic.types'

const appStore = useAppStore()
const authStore = useAuthStore()
const queryClient = useQueryClient()

const selectedTopicIds = ref<Array<string>>([])
const questionCount = ref<number>(10)
const selectedDifficulty = ref<PromptDifficulty>('mixed')
const selectedObjective = ref<PromptObjective>('review')
const customContext = ref('')
const promptText = ref('')
const csvText = ref('')
const importMode = ref<'file' | 'paste'>('file')
const selectedFile = ref<File | null>(null)
const fileInputKey = ref(0)
const importResult = ref<ImportResult | null>(null)
const importing = ref(false)
const promptValidated = ref(false)

const { data: quota } = useQuery({
  queryKey: ['questions', 'import-quota'],
  queryFn: () => getImportQuota(),
  staleTime: 30 * 1000,
  enabled: computed(() => !authStore.isAdmin),
})

const selectedTopics = ref<Array<Topic>>([])

const topicTitle = (topic: Topic) => `${topic.name} (${topic.slug})`

const quotaPercent = computed(() => {
  if (!quota.value || quota.value.dailyLimit === 0) return 0
  return Math.min(100, Math.round((quota.value.usedToday / quota.value.dailyLimit) * 100))
})

const quotaColor = computed(() => {
  const p = quotaPercent.value
  if (p >= 90) return 'error'
  if (p >= 70) return 'warning'
  return 'success'
})

const questionCountError = computed(() => {
  const n = Number(questionCount.value)
  if (!Number.isInteger(n) || n < 1) return 'Mínimo 1'
  if (n > 50) return 'Máximo 50'
  return ''
})

const topicsPromptError = computed<Array<string>>(() => {
  if (!promptValidated.value) return []
  return validateRules([requiredRule()], selectedTopicIds.value.length > 0 ? 'filled' : '')
})

function rebuildPrompt() {
  const n = Number(questionCount.value)
  const qty = Number.isInteger(n) && n >= 1 && n <= 50 ? n : 5
  promptText.value = buildImportPrompt(selectedTopics.value, qty, {
    difficulty: selectedDifficulty.value,
    objective: selectedObjective.value,
    customContext: customContext.value,
  })
}

watch(selectedTopicIds, rebuildPrompt, { immediate: true })
watch(questionCount, rebuildPrompt)
watch(selectedDifficulty, rebuildPrompt)
watch(selectedObjective, rebuildPrompt)
watch(customContext, rebuildPrompt)

async function copyPrompt() {
  promptValidated.value = true
  if (topicsPromptError.value.length > 0 || questionCountError.value) return
  if (!promptText.value) return
  await navigator.clipboard.writeText(promptText.value)
  appStore.showSnackbar('Prompt copiado al portapapeles')
}

function downloadTemplate() {
  const blob = new Blob([IMPORT_CSV_TEMPLATE], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'plantilla-preguntas.csv'
  a.click()
  URL.revokeObjectURL(url)
}

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  selectedFile.value = input.files?.[0] ?? null
}

async function doImport() {
  importing.value = true
  importResult.value = null
  try {
    const formData = new FormData()
    if (importMode.value === 'file') {
      if (!selectedFile.value) {
        appStore.showSnackbar('Selecciona un archivo CSV', 'error')
        importing.value = false
        return
      }
      formData.append('file', selectedFile.value)
    } else {
      if (!csvText.value.trim()) {
        appStore.showSnackbar('Pega el contenido CSV', 'error')
        importing.value = false
        return
      }
      formData.append('content', csvText.value)
    }

    const result = await importQuestions(formData)
    importResult.value = result
    queryClient.invalidateQueries({ queryKey: ['questions', 'list'] })
    queryClient.invalidateQueries({ queryKey: ['questions', 'import-quota'] })
    if (result.imported > 0) {
      appStore.showSnackbar(`Se importaron ${result.imported} preguntas`)
      // Limpiar para evitar re-importar accidentalmente
      if (importMode.value === 'paste') {
        csvText.value = ''
      } else {
        selectedFile.value = null
        fileInputKey.value++
      }
    }
    if (result.failed > 0) {
      appStore.showSnackbar(`${result.failed} filas con errores`, 'error')
    }
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al importar'
    appStore.showSnackbar(detail, 'error')
  } finally {
    importing.value = false
  }
}
</script>

<template>
  <v-container>
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" to="/questions" class="mr-2" />
      <h1 class="text-h4">Importar preguntas</h1>
    </div>

    <v-card v-if="quota && !authStore.isAdmin" class="mb-4">
      <v-card-text>
        <div class="d-flex align-center justify-space-between flex-wrap ga-2 mb-2">
          <span class="text-body-2 font-weight-medium">
            Cupo diario: {{ quota.usedToday }} / {{ quota.dailyLimit }} usadas
            <span class="text-medium-emphasis">— Restantes: {{ quota.remaining }}</span>
          </span>
          <span class="text-caption"
            >Máximo por archivo: {{ MAX_QUESTIONS_PER_FILE }} preguntas</span
          >
        </div>
        <v-progress-linear :model-value="quotaPercent" :color="quotaColor" height="22" rounded>
          <template #default>
            <span class="text-caption font-weight-medium">{{ quotaPercent }}%</span>
          </template>
        </v-progress-linear>
        <div class="text-caption text-medium-emphasis mt-2">
          Los slugs deben coincidir con temas ya existentes; un slug inexistente o una fila sin
          temas genera error.
        </div>
      </v-card-text>
    </v-card>
    <v-card v-else-if="authStore.isAdmin" class="mb-4">
      <v-card-text>
        <div class="d-flex align-center justify-space-between flex-wrap ga-2">
          <span class="text-body-2 font-weight-medium"
            >Importación de administrador: sin límite diario</span
          >
          <span class="text-caption"
            >Máximo por archivo: {{ MAX_QUESTIONS_PER_FILE }} preguntas</span
          >
        </div>
        <div class="text-caption text-medium-emphasis mt-2">
          Las preguntas importadas como admin son públicas (visibles para todos los usuarios). Los
          slugs deben coincidir con temas ya existentes.
        </div>
      </v-card-text>
    </v-card>

    <v-row>
      <!-- Columna: Generar con IA -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title class="d-flex align-center ga-2">
            <v-icon icon="mdi-robot" color="primary" />
            Generar con IA
          </v-card-title>
          <v-card-subtitle> Prepara el prompt para tu IA </v-card-subtitle>
          <v-card-text>
            <p class="text-caption mb-3">
              Selecciona los temas para incluir sus slugs exactos en el prompt. Así la IA generará
              el CSV sin errores de nombres.
            </p>
            <v-alert type="info" variant="tonal" density="compact" class="mb-3">
              El contenido soporta formato enriquecido:
              <code>`código inline`</code>, bloques <code>```lang```</code> y
              <strong>**negrita**</strong>. Úsalo en enunciado, explicación y opciones para que las
              preguntas sean más atractivas. El prompt de la IA ya incluye estas instrucciones.
            </v-alert>
            <v-row dense class="mb-3">
              <v-col style="flex: 0 0 60%; max-width: 60%">
                <TopicAutocomplete
                  v-model="selectedTopicIds"
                  v-model:selected-topics="selectedTopics"
                  label="Temas para el prompt *"
                  :item-title="topicTitle"
                  :error-messages="topicsPromptError"
                  required
                  density="compact"
                  placeholder="Elige temas para el prompt"
                />
              </v-col>
              <v-col style="flex: 0 0 40%; max-width: 40%">
                <v-text-field
                  v-model.number="questionCount"
                  label="Cantidad *"
                  type="number"
                  :min="1"
                  :max="50"
                  :error-messages="questionCountError || undefined"
                  density="compact"
                  variant="outlined"
                  hide-details="auto"
                  hint="1 a 50"
                  persistent-hint
                  required
                />
              </v-col>
            </v-row>
            <v-row dense class="mb-3">
              <v-col cols="12" sm="6">
                <v-select
                  v-model="selectedObjective"
                  label="Objetivo"
                  :items="PROMPT_OBJECTIVES"
                  item-title="title"
                  item-value="value"
                  density="compact"
                  variant="outlined"
                  hide-details="auto"
                >
                  <template #item="{ props: itemProps, item }">
                    <v-list-item v-bind="itemProps">
                      <template #subtitle>
                        <span class="text-caption">{{ item.raw.description }}</span>
                      </template>
                    </v-list-item>
                  </template>
                </v-select>
              </v-col>
              <v-col cols="12" sm="6">
                <v-select
                  v-model="selectedDifficulty"
                  label="Dificultad"
                  :items="PROMPT_DIFFICULTIES"
                  item-title="title"
                  item-value="value"
                  density="compact"
                  variant="outlined"
                  hide-details="auto"
                />
              </v-col>
            </v-row>
            <v-textarea
              v-model="customContext"
              label="Contexto adicional (opcional)"
              placeholder="Ej: preguntas enfocadas a Go con concurrencia para entrevista senior, con ejemplos de canales y goroutines…"
              rows="2"
              auto-grow
              density="compact"
              variant="outlined"
              hide-details="auto"
              class="mb-3"
            />
            <div
              v-if="selectedObjective === 'custom' && !customContext.trim()"
              class="text-caption text-warning mb-3"
            >
              Con objetivo “Personalizado”, escribe tu contexto arriba para orientar a la IA.
            </div>
            <v-textarea
              v-model="promptText"
              label="Prompt para la IA"
              no-auto-grow
              rows="12"
              readonly
              density="compact"
              variant="outlined"
              class="mb-3 fixed-textarea"
            />
            <div class="d-flex ga-2">
              <v-btn
                color="primary"
                variant="tonal"
                prepend-icon="mdi-content-copy"
                @click="copyPrompt"
              >
                Copiar prompt
              </v-btn>
              <v-btn variant="tonal" prepend-icon="mdi-download" @click="downloadTemplate">
                Descargar plantilla CSV
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Columna: Importar -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title class="d-flex align-center ga-2">
            <v-icon icon="mdi-file-upload" color="primary" />
            Importar CSV
          </v-card-title>
          <v-card-subtitle> Sube un archivo o pega el contenido </v-card-subtitle>
          <v-card-text>
            <v-btn-toggle
              v-model="importMode"
              mandatory
              density="compact"
              class="mb-3"
              color="primary"
            >
              <v-btn value="file" prepend-icon="mdi-file-delimited">Archivo</v-btn>
              <v-btn value="paste" prepend-icon="mdi-clipboard-text">Pegar texto</v-btn>
            </v-btn-toggle>

            <div v-if="importMode === 'file'" class="mb-3">
              <v-file-input
                :key="fileInputKey"
                label="Archivo CSV"
                accept=".csv,text/csv"
                :disabled="importing"
                hide-details
                density="compact"
                variant="outlined"
                @change="onFileChange"
              />
              <div class="text-caption text-medium-emphasis mt-1">
                Encabezado requerido: type,content,difficulty,topics,explanation,options — el
                contenido admite <code>`código`</code>, bloques <code>```</code> y
                <strong>**negrita**</strong>.
              </div>
            </div>
            <div v-else class="mb-3">
              <v-textarea
                v-model="csvText"
                label="Pega aquí el CSV"
                :disabled="importing"
                no-auto-grow
                rows="10"
                density="compact"
                variant="outlined"
                class="fixed-textarea"
                placeholder='type,content,difficulty,topics,explanation,options&#10;single_choice,"¿...?",beginner,go,"...","[v] a | [ ] b"'
              />
            </div>

            <v-btn
              color="primary"
              :loading="importing"
              block
              size="large"
              prepend-icon="mdi-upload"
              @click="doImport"
            >
              Importar
            </v-btn>

            <div v-if="importResult" class="mt-4">
              <v-alert
                :type="importResult.failed === 0 ? 'success' : 'warning'"
                variant="tonal"
                class="mb-3"
              >
                Importadas: {{ importResult.imported }} / {{ importResult.total }} — Fallidas:
                {{ importResult.failed }}
              </v-alert>
              <v-data-table
                v-if="importResult.errors.length"
                :headers="[
                  { title: 'Fila', key: 'row' },
                  { title: 'Motivo', key: 'reason' },
                ]"
                :items="importResult.errors"
                density="compact"
                hide-default-footer
                no-data-text="Sin errores"
              />
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.fixed-textarea :deep(textarea) {
  overflow-y: auto !important;
}
</style>
