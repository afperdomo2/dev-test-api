<script setup lang="ts">
import type { Question } from '@/types/question.types'
import { DIFFICULTY_COLORS, TYPE_ICONS } from '@/types/question.types'
import { formatDate } from '@/utils/format'

interface Props {
  questions: Array<Question>
  loading: boolean
  itemsPerPage: number
}

defineProps<Props>()

const headers = [
  { title: 'Tipo', key: 'type', sortable: false, align: 'center' as const, width: 72 },
  { title: 'Contenido', key: 'content', sortable: false },
  { title: 'Dificultad', key: 'difficulty', sortable: false, align: 'center' as const, width: 120 },
  { title: 'Temas', key: 'topics', sortable: false, width: 200 },
  { title: 'Creado', key: 'createdAt', sortable: false, align: 'center' as const, width: 120 },
]
</script>

<template>
  <v-data-table
    :headers="headers"
    :items="questions"
    :items-per-page="itemsPerPage"
    :loading="loading"
    no-data-text="No hay preguntas"
    loading-text="Cargando preguntas..."
    hover
  >
    <template #[`item.type`]="{ item }">
      <div class="d-flex justify-center">
        <v-icon :icon="TYPE_ICONS[item.type]" color="primary" size="small" />
      </div>
    </template>

    <template #[`item.content`]="{ item }">
      <router-link :to="`/questions/${item.id}`" class="text-decoration-none text-body-2">
        <span class="text-truncate d-inline-block" style="max-width: 400px" :title="item.content">
          {{ item.content }}
        </span>
      </router-link>
    </template>

    <template #[`item.difficulty`]="{ item }">
      <div class="d-flex justify-center">
        <v-chip :color="DIFFICULTY_COLORS[item.difficulty]" size="x-small" variant="tonal">
          {{ item.difficulty }}
        </v-chip>
      </div>
    </template>

    <template #[`item.topics`]="{ item }">
      <div v-if="item.topics.length" class="d-flex flex-wrap ga-1">
        <v-chip
          v-for="topic in item.topics.slice(0, 2)"
          :key="topic"
          size="x-small"
          variant="outlined"
        >
          {{ topic }}
        </v-chip>
        <span v-if="item.topics.length > 2" class="text-caption text-medium-emphasis ml-1">
          +{{ item.topics.length - 2 }}
        </span>
      </div>
    </template>

    <template #[`item.createdAt`]="{ item }">
      <div class="d-flex justify-center text-caption text-medium-emphasis">
        {{ formatDate(item.createdAt) }}
      </div>
    </template>

    <template #bottom>
      <slot name="footer" />
    </template>
  </v-data-table>
</template>
