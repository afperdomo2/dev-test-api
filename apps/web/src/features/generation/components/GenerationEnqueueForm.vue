<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import TopicAutocomplete from '@/components/TopicAutocomplete.vue'
import DifficultySelect from '@/components/DifficultySelect.vue'
import { IMPORT_PROMPT_KINDS, type PromptQuestionKind } from '@/constants/questionImport'
import type {
  GenerationDifficulty,
  GenerationKind,
  GenerationQuestionType,
  GenerationStats,
} from '@/types/generation.types'

interface Props {
  kind?: GenerationKind
  stats?: GenerationStats | null
  loading?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  kind: 'question',
  stats: null,
  loading: false,
  disabled: false,
})

const emit = defineEmits<{
  submit: [
    payload: {
      topicIds: Array<string>
      difficulty: GenerationDifficulty
      questionType: GenerationQuestionType
      quantity: number
    },
  ]
}>()

const topicIds = ref<Array<string>>([])
const difficulty = ref<GenerationDifficulty>('random')
const questionType = ref<PromptQuestionKind>('mixed')
const quantity = ref<number>(20)
const touched = ref(false)

const isTips = computed(() => props.kind === 'tip')
const noun = computed(() => (isTips.value ? 'tips' : 'preguntas'))

const perJobCap = computed(() => props.stats?.perJobCap ?? 20)
const defaultQuantity = computed(() => props.stats?.defaultQuantity ?? 20)
const reservedToday = computed(() => props.stats?.reservedToday ?? 0)
const dailyCap = computed(() => props.stats?.dailyCap ?? 200)
const remainingToday = computed(() => Math.max(0, dailyCap.value - reservedToday.value))

/** El tope efectivo nunca supera lo que queda de presupuesto diario. */
const maxSelectable = computed(() => Math.min(perJobCap.value, remainingToday.value))
const budgetExhausted = computed(() => remainingToday.value === 0)

// La cantidad inicial sigue al default del servidor y se recorta si queda poco
// presupuesto hoy. No pisa lo que el usuario ya escribió.
watch(
  [defaultQuantity, maxSelectable],
  ([value, cap]) => {
    if (touched.value) return
    quantity.value = Math.max(1, Math.min(value, cap || value))
  },
  { immediate: true },
)

watch([difficulty, quantity, topicIds], () => {
  touched.value = true
})

const canSubmit = computed(
  () =>
    !props.disabled &&
    !props.loading &&
    topicIds.value.length > 0 &&
    quantity.value > 0 &&
    quantity.value <= maxSelectable.value,
)

const quantityHint = computed(() => {
  if (budgetExhausted.value) return `Tope diario alcanzado (${dailyCap.value}/${dailyCap.value})`
  if (maxSelectable.value < perJobCap.value) {
    return `Máx. ${maxSelectable.value} (quedan ${remainingToday.value} de ${dailyCap.value} hoy)`
  }
  return `Máx. ${perJobCap.value} ${noun.value} por trabajo`
})

/** Resumen legible de lo que se va a generar, para que no haya sorpresas. */
const questionTypeHint = computed(() => {
  const n = quantity.value > 0 ? quantity.value : 0
  const prefix = n > 0 ? `Se generarán ${n} preguntas` : 'Se generarán preguntas'
  switch (questionType.value) {
    case 'single_choice':
      return `${prefix} de selección única: 4 opciones y 1 correcta.`
    case 'multiple_choice':
      return `${prefix} de selección múltiple: 4 opciones y 2 o 3 correctas.`
    case 'true_false':
      return `${prefix} de verdadero/falso: 2 opciones y 1 correcta.`
    default:
      return `${prefix} mezclando selección única, selección múltiple y verdadero/falso.`
  }
})

function submit() {
  if (!canSubmit.value) return
  emit('submit', {
    topicIds: topicIds.value,
    difficulty: difficulty.value,
    questionType: questionType.value,
    quantity: quantity.value,
  })
}

function reset() {
  topicIds.value = []
  difficulty.value = 'random'
  questionType.value = 'mixed'
  quantity.value = defaultQuantity.value
  touched.value = false
}

defineExpose({ reset })
</script>

<template>
  <v-card>
    <v-card-item>
      <v-card-title class="text-subtitle-1">
        {{ isTips ? 'Encolar generación de tips' : 'Encolar generación de preguntas' }}
      </v-card-title>
      <v-card-subtitle>
        {{
          isTips
            ? 'Un trabajo = una llamada a la IA. Cada tip se asigna a un tema.'
            : 'Un trabajo = una llamada a la IA. Las preguntas se guardan como públicas.'
        }}
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <v-row dense>
        <v-col cols="12" :md="isTips ? 8 : 6">
          <TopicAutocomplete
            v-model="topicIds"
            :label="isTips ? 'Temas para los tips' : 'Temas'"
            required
          />
        </v-col>
        <v-col v-if="!isTips" cols="12" md="3">
          <DifficultySelect v-model="difficulty" />
        </v-col>
        <v-col cols="12" :md="isTips ? 4 : 3">
          <v-text-field
            v-model.number="quantity"
            :min="1"
            :max="maxSelectable"
            type="number"
            label="Cantidad"
            :hint="quantityHint"
            persistent-hint
            variant="outlined"
            density="compact"
            :disabled="budgetExhausted"
          />
        </v-col>
        <v-col v-if="!isTips" cols="12" md="6">
          <v-select
            v-model="questionType"
            label="Tipo de pregunta"
            :items="IMPORT_PROMPT_KINDS"
            item-title="title"
            item-value="value"
            density="compact"
            variant="outlined"
            hide-details
          />
        </v-col>
      </v-row>

      <v-alert v-if="!isTips" type="info" variant="tonal" density="compact" class="mt-4">
        {{ questionTypeHint }}
        <div class="text-caption mt-1">
          Completar código no se genera con IA: se crea a mano desde Preguntas.
        </div>
      </v-alert>
    </v-card-text>

    <v-card-actions>
      <v-btn variant="text" prepend-icon="mdi-eraser" @click="reset"> Limpiar </v-btn>
      <v-spacer />
      <v-btn
        color="primary"
        prepend-icon="mdi-auto-fix"
        :loading="loading"
        :disabled="!canSubmit"
        @click="submit"
      >
        Generar {{ isTips ? 'tips' : 'preguntas' }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
