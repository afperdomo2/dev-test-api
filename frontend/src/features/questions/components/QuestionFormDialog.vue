<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { createQuestionMutation, updateQuestionMutation } from '@/queries/questions.queries'
import { listTopics } from '@/api/services/topics.service'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import { requiredRule, validateRules } from '@/utils/validators'
import { QUESTION_TYPES, QUESTION_DIFFICULTIES } from '@/types/question.types'
import type {
  Question,
  QuestionType,
  QuestionDifficulty,
  CreateQuestionOption,
  CreateQuestionRequest,
  UpdateQuestionRequest,
} from '@/types/question.types'

interface Props {
  modelValue: boolean
  question: Question | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()
const { extractFieldErrors } = useFormErrors()

const isEdit = computed(() => !!props.question)
const dialogTitle = computed(() => (isEdit.value ? 'Editar pregunta' : 'Nueva pregunta'))

const { data: topicsData, isLoading: topicsLoading } = useQuery({
  queryKey: ['topics', 'list', 1, 100, 'name', 'asc'],
  queryFn: () => listTopics(1, 100, 'name', 'asc'),
  staleTime: 60 * 1000,
  enabled: computed(() => props.modelValue),
})

const topicItems = computed(() => {
  return topicsData.value?.data?.map((t) => ({ title: t.name, value: t.id })) ?? []
})

const selectedType = ref<QuestionType>('single_choice')
const content = ref('')
const difficulty = ref<QuestionDifficulty>('intermediate')
const explanation = ref('')
const topicIds = ref<Array<string>>([])

const options = ref<Array<{ content: string; isCorrect: boolean }>>([
  { content: '', isCorrect: false },
  { content: '', isCorrect: false },
])

const language = ref('')
const starterCode = ref('')
const expectedOutput = ref('')
const testCases = ref('')

const isChoiceType = computed(
  () => selectedType.value === 'single_choice' || selectedType.value === 'multiple_choice',
)
const isCodeType = computed(() => selectedType.value === 'code_completion')

const validationErrors = ref<Record<string, Array<string>>>({})
const serverErrors = ref<Record<string, string>>({})
const saving = ref(false)

const createMut = useMutation(createQuestionMutation())
const updateMut = useMutation(updateQuestionMutation())

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      if (props.question) {
        selectedType.value = props.question.type
        content.value = props.question.content
        difficulty.value = props.question.difficulty
        explanation.value = props.question.explanation ?? ''
        topicIds.value = [...props.question.topics]
        if (props.question.options?.length) {
          options.value = props.question.options.map((o) => ({
            content: o.content,
            isCorrect: o.isCorrect ?? false,
          }))
        }
        if (props.question.codeChallenge) {
          language.value = props.question.codeChallenge.language
          starterCode.value = props.question.codeChallenge.starterCode
          expectedOutput.value = props.question.codeChallenge.expectedOutput
          testCases.value = props.question.codeChallenge.testCases
        }
      } else {
        selectedType.value = 'single_choice'
        content.value = ''
        difficulty.value = 'intermediate'
        explanation.value = ''
        topicIds.value = []
        options.value = [
          { content: '', isCorrect: false },
          { content: '', isCorrect: false },
        ]
        language.value = ''
        starterCode.value = ''
        expectedOutput.value = ''
        testCases.value = ''
      }
      validationErrors.value = {}
      serverErrors.value = {}
    }
  },
)

function addOption() {
  options.value.push({ content: '', isCorrect: false })
}

function removeOption(index: number) {
  options.value.splice(index, 1)
}

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.content = validateRules([requiredRule()], content.value)
  newErrors.topicIds = validateRules([requiredRule()], topicIds.value.length > 0 ? 'filled' : '')
  if (isChoiceType.value) {
    const filledOptions = options.value.filter((o) => o.content.trim())
    if (filledOptions.length < 2) {
      newErrors.options = ['Debe tener al menos 2 opciones']
    }
    if (!filledOptions.some((o) => o.isCorrect)) {
      newErrors.options = [
        ...(newErrors.options ?? []),
        'Debe marcar al menos una opción como correcta',
      ]
    }
  }
  if (isCodeType.value) {
    newErrors.language = validateRules([requiredRule()], language.value)
  }
  validationErrors.value = newErrors
  serverErrors.value = {}
  return Object.values(newErrors).every((e) => e.length === 0)
}

function fieldError(field: string): Array<string> {
  const server = serverErrors.value[field]
  const client = validationErrors.value[field] ?? []
  return server ? [...client, server] : client
}

async function submit() {
  if (!validate()) return
  saving.value = true
  try {
    if (isEdit.value && props.question) {
      const data: UpdateQuestionRequest = {
        type: selectedType.value,
        content: content.value,
        difficulty: difficulty.value,
        explanation: explanation.value || undefined,
        topicIds: topicIds.value,
      }
      if (isChoiceType.value) {
        data.options = options.value.filter((o) => o.content.trim()) as Array<CreateQuestionOption>
      } else {
        data.language = language.value
        data.starterCode = starterCode.value || undefined
        data.expectedOutput = expectedOutput.value || undefined
        data.testCases = testCases.value || undefined
      }
      await updateMut.mutateAsync({ id: props.question.id, data })
      appStore.showSnackbar('Pregunta actualizada')
    } else {
      const data: CreateQuestionRequest = {
        type: selectedType.value,
        content: content.value,
        difficulty: difficulty.value,
        explanation: explanation.value || undefined,
        topicIds: topicIds.value,
      }
      if (isChoiceType.value) {
        data.options = options.value.filter((o) => o.content.trim()) as Array<CreateQuestionOption>
      } else {
        data.language = language.value
        data.starterCode = starterCode.value || undefined
        data.expectedOutput = expectedOutput.value || undefined
        data.testCases = testCases.value || undefined
      }
      await createMut.mutateAsync(data)
      appStore.showSnackbar('Pregunta creada')
    }
    queryClient.invalidateQueries({ queryKey: ['questions', 'list'] })
    emit('saved')
    close()
  } catch (err: unknown) {
    const fieldErrors = extractFieldErrors(err)
    if (Object.keys(fieldErrors).length > 0) {
      serverErrors.value = fieldErrors
    }
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al guardar pregunta'
    appStore.showSnackbar(detail, 'error')
  } finally {
    saving.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    max-width="640"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>

      <v-card-text>
        <v-form @submit.prevent="submit">
          <v-select
            v-model="selectedType"
            label="Tipo"
            :items="QUESTION_TYPES"
            :disabled="saving || isEdit"
            hide-details
            density="compact"
            class="mb-3"
          />

          <v-textarea
            v-model="content"
            label="Contenido"
            :error-messages="fieldError('content')"
            :disabled="saving"
            auto-grow
            rows="3"
            required
          />

          <v-select
            v-model="difficulty"
            label="Dificultad"
            :items="QUESTION_DIFFICULTIES"
            hide-details
            density="compact"
            class="mb-3"
          />

          <v-textarea
            v-model="explanation"
            label="Explicación (opcional)"
            :disabled="saving"
            auto-grow
            rows="2"
            hide-details
            class="mb-3"
          />

          <v-autocomplete
            v-model="topicIds"
            label="Temas"
            :items="topicItems"
            :loading="topicsLoading"
            :error-messages="fieldError('topicIds')"
            :disabled="saving || !topicItems.length"
            multiple
            chips
            closable-chips
            hide-details
            density="compact"
            class="mb-3"
          />

          <!-- Choice options -->
          <template v-if="isChoiceType">
            <v-divider class="mb-3" />
            <div class="text-subtitle-2 mb-2">Opciones</div>
            <v-card v-for="(opt, idx) in options" :key="idx" variant="outlined" class="mb-2 pa-2">
              <v-row dense align="center">
                <v-col cols="8">
                  <v-text-field
                    v-model="opt.content"
                    :label="`Opción ${idx + 1}`"
                    :disabled="saving"
                    hide-details
                    density="compact"
                  />
                </v-col>
                <v-col cols="2" class="text-center">
                  <v-switch
                    v-model="opt.isCorrect"
                    :label="opt.isCorrect ? 'Correcta' : ''"
                    color="success"
                    hide-details
                    density="compact"
                  />
                </v-col>
                <v-col cols="2" class="text-center">
                  <v-btn
                    v-if="options.length > 2"
                    icon="mdi-close"
                    variant="text"
                    size="small"
                    color="error"
                    :disabled="saving"
                    @click="removeOption(idx)"
                  />
                </v-col>
              </v-row>
            </v-card>
            <div v-if="fieldError('options').length" class="text-caption text-error mb-2">
              {{ fieldError('options')[0] }}
            </div>
            <v-btn
              variant="text"
              color="primary"
              size="small"
              prepend-icon="mdi-plus"
              :disabled="saving"
              @click="addOption"
            >
              Agregar opción
            </v-btn>
          </template>

          <!-- Code challenge -->
          <template v-if="isCodeType">
            <v-divider class="mb-3" />
            <div class="text-subtitle-2 mb-2">Código</div>
            <v-text-field
              v-model="language"
              label="Lenguaje"
              :error-messages="fieldError('language')"
              :disabled="saving"
              hide-details
              density="compact"
              class="mb-3"
            />
            <v-textarea
              v-model="starterCode"
              label="Código base (opcional)"
              :disabled="saving"
              auto-grow
              rows="4"
              hide-details
              class="mb-3"
            />
            <v-text-field
              v-model="expectedOutput"
              label="Salida esperada (opcional)"
              :disabled="saving"
              hide-details
              density="compact"
              class="mb-3"
            />
          </template>
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" :disabled="saving" @click="close"> Cancelar </v-btn>
        <v-btn color="primary" :loading="saving" @click="submit">
          {{ isEdit ? 'Guardar' : 'Crear' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
