<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import TopicAutocomplete from '@/components/TopicAutocomplete.vue'
import type { GenerationDifficulty, GenerationStats } from '@/types/generation.types'

interface Props {
  stats?: GenerationStats | null
  loading?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  stats: null,
  loading: false,
  disabled: false,
})

const emit = defineEmits<{
  submit: [payload: { topicIds: Array<string>; difficulty: GenerationDifficulty; quantity: number }]
}>()

const topicIds = ref<Array<string>>([])
const difficulty = ref<GenerationDifficulty>('intermediate')
const quantity = ref<number>(5)
const touched = ref(false)

const difficultyOptions: Array<{ title: string; value: GenerationDifficulty }> = [
  { title: 'Principiante', value: 'beginner' },
  { title: 'Intermedio', value: 'intermediate' },
  { title: 'Avanzado', value: 'advanced' },
  { title: 'Aleatoria', value: 'random' },
]

const perJobCap = computed(() => props.stats?.perJobCap ?? 10)
const defaultQuantity = computed(() => props.stats?.defaultQuantity ?? 5)
const reservedToday = computed(() => props.stats?.reservedToday ?? 0)
const dailyCap = computed(() => props.stats?.dailyCap ?? 200)
const remainingToday = computed(() => Math.max(0, dailyCap.value - reservedToday.value))

/** El tope efectivo nunca supera lo que queda de presupuesto diario. */
const maxSelectable = computed(() => Math.min(perJobCap.value, remainingToday.value))
const budgetExhausted = computed(() => remainingToday.value === 0)

watch(defaultQuantity, (value) => {
  if (!touched.value) quantity.value = Math.min(value, maxSelectable.value || value)
})

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
  return `Máx. ${perJobCap.value} por trabajo`
})

function submit() {
  if (!canSubmit.value) return
  emit('submit', {
    topicIds: topicIds.value,
    difficulty: difficulty.value,
    quantity: quantity.value,
  })
}
</script>

<template>
  <v-card>
    <v-card-item>
      <v-card-title class="text-subtitle-1">Encolar generación</v-card-title>
      <v-card-subtitle>
        Un trabajo = una llamada a la IA. Las preguntas se guardan como públicas.
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <v-row dense>
        <v-col cols="12" md="6">
          <TopicAutocomplete v-model="topicIds" label="Temas" required />
        </v-col>
        <v-col cols="12" md="3">
          <v-select
            v-model="difficulty"
            :items="difficultyOptions"
            label="Dificultad"
            hide-details
            variant="outlined"
            density="compact"
          />
        </v-col>
        <v-col cols="12" md="3">
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
      </v-row>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <v-btn
        color="primary"
        prepend-icon="mdi-auto-fix"
        :loading="loading"
        :disabled="!canSubmit"
        @click="submit"
      >
        Generar
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
