<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { QuestionsFilters } from '@/api/services/questions.service'
import { QUESTION_TYPES, QUESTION_DIFFICULTIES } from '@/types/question.types'

const emit = defineEmits<{
  change: [filters: QuestionsFilters]
}>()

const selectedType = ref<string>('')
const selectedDifficulty = ref<string>('')

const hasFilters = computed(
  () => !!selectedType.value || !!selectedDifficulty.value,
)

function updateFilters() {
  emit('change', {
    type: selectedType.value || undefined,
    difficulty: selectedDifficulty.value || undefined,
  })
}

watch(selectedType, updateFilters)
watch(selectedDifficulty, updateFilters)

function clearFilters() {
  selectedType.value = ''
  selectedDifficulty.value = ''
}
</script>

<template>
  <v-card class="mb-4" variant="flat" border>
    <v-card-text>
      <v-row align="center" dense>
        <v-col cols="6" md="4">
          <v-select
            v-model="selectedType"
            label="Tipo"
            :items="QUESTION_TYPES"
            clearable
            hide-details
            density="compact"
          />
        </v-col>

        <v-col cols="6" md="4">
          <v-select
            v-model="selectedDifficulty"
            label="Dificultad"
            :items="QUESTION_DIFFICULTIES"
            clearable
            hide-details
            density="compact"
          />
        </v-col>

        <v-col v-if="hasFilters" cols="12" md="4" class="text-end">
          <v-btn size="small" variant="text" color="error" @click="clearFilters">
            Limpiar filtros
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>
