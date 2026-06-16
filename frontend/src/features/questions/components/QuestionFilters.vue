<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useDebounce } from '@/composables/useDebounce'
import type { QuestionsFilters } from '@/api/services/questions.service'
import { QUESTION_TYPES, QUESTION_DIFFICULTIES } from '@/types/question.types'

const emit = defineEmits<{
  change: [filters: QuestionsFilters]
}>()

const { value: searchText, debouncedValue: debouncedSearch } = useDebounce('', 500)

const selectedType = ref<string>('')
const selectedDifficulty = ref<string>('')

const hasFilters = computed(
  () => selectedType.value || selectedDifficulty.value || searchText.value,
)

function updateFilters() {
  emit('change', {
    type: selectedType.value || undefined,
    difficulty: selectedDifficulty.value || undefined,
  })
}

watch(debouncedSearch, updateFilters)
watch(selectedType, updateFilters)
watch(selectedDifficulty, updateFilters)

function clearFilters() {
  searchText.value = ''
  selectedType.value = ''
  selectedDifficulty.value = ''
}
</script>

<template>
  <v-card class="mb-4" variant="outlined">
    <v-card-text>
      <v-row align="center" dense>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="searchText"
            label="Buscar pregunta"
            prepend-inner-icon="mdi-magnify"
            clearable
            hide-details
            density="compact"
          />
        </v-col>

        <v-col cols="6" md="3">
          <v-select
            v-model="selectedType"
            label="Tipo"
            :items="QUESTION_TYPES"
            clearable
            hide-details
            density="compact"
          />
        </v-col>

        <v-col cols="6" md="3">
          <v-select
            v-model="selectedDifficulty"
            label="Dificultad"
            :items="QUESTION_DIFFICULTIES"
            clearable
            hide-details
            density="compact"
          />
        </v-col>

        <v-col v-if="hasFilters" cols="12" md="2" class="text-end">
          <v-btn size="small" variant="text" color="error" @click="clearFilters">
            Limpiar filtros
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>
