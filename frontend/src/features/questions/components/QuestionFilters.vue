<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { listTopics } from '@/api/services/topics.service'
import type { QuestionsFilters } from '@/api/services/questions.service'
import { QUESTION_TYPES, QUESTION_DIFFICULTIES } from '@/types/question.types'

const props = defineProps<{
  initialTopicIds?: string
}>()

const emit = defineEmits<{
  change: [filters: QuestionsFilters]
}>()

const selectedType = ref<string>('')
const selectedDifficulty = ref<string>('')
const selectedTopicIds = ref<Array<string>>([])

const { data: topicsData } = useQuery({
  queryKey: ['topics', 'list', 1, 100, 'name', 'asc'],
  queryFn: () => listTopics(1, 100, 'name', 'asc'),
  staleTime: 60 * 1000,
})

const topicItems = computed(() =>
  (topicsData.value?.data ?? []).map((t) => ({
    title: t.name,
    value: t.id,
    props: { subtitle: t.category },
  })),
)

const hasFilters = computed(
  () => !!selectedType.value || !!selectedDifficulty.value || selectedTopicIds.value.length > 0,
)

function updateFilters() {
  emit('change', {
    type: selectedType.value || undefined,
    difficulty: selectedDifficulty.value || undefined,
    topicIds: selectedTopicIds.value.length ? selectedTopicIds.value.join(',') : undefined,
  })
}

watch(selectedType, updateFilters)
watch(selectedDifficulty, updateFilters)
watch(selectedTopicIds, updateFilters)

watch(
  () => props.initialTopicIds,
  (val) => {
    if (val) selectedTopicIds.value = val.split(',').filter(Boolean)
    else if (val === '') selectedTopicIds.value = []
  },
  { immediate: true },
)

function clearFilters() {
  selectedType.value = ''
  selectedDifficulty.value = ''
  selectedTopicIds.value = []
}
</script>

<template>
  <v-card class="mb-4" variant="flat" border>
    <v-card-text>
      <v-row align="center" dense>
        <v-col cols="12" md="4">
          <v-autocomplete
            v-model="selectedTopicIds"
            label="Tema"
            :items="topicItems"
            multiple
            chips
            closable-chips
            clearable
            hide-details
            density="compact"
            placeholder="Filtrar por tema"
            no-data-text="Sin temas"
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
