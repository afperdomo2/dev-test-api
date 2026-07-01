<script setup lang="ts">
import type { Question } from '@/types/question.types'
import {
  DIFFICULTY_COLORS,
  DIFFICULTY_LABELS,
  TYPE_ICONS,
  SOURCE_LABELS,
  SOURCE_COLORS,
} from '@/types/question.types'
import { formatDate } from '@/utils/format'

interface Props {
  questions: Array<Question>
  loading: boolean
  itemsPerPage: number
  currentUserId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  edit: [question: Question]
  delete: [question: Question]
}>()

function canModify(question: Question): boolean {
  if (question.source !== 'manual' && question.source !== 'imported') return false
  return question.userId === props.currentUserId
}

function lockReason(question: Question): string {
  if (question.source !== 'manual' && question.source !== 'imported') {
    return 'Las preguntas generadas por IA no se pueden modificar'
  }
  return 'No tienes permiso para modificar esta pregunta'
}

const headers = [
  { title: 'Tipo', key: 'type', sortable: false, align: 'center' as const, width: 72 },
  { title: 'Origen', key: 'source', sortable: false, align: 'center' as const, width: 88 },
  { title: 'Contenido', key: 'content', sortable: false },
  { title: 'Dificultad', key: 'difficulty', sortable: false, align: 'center' as const, width: 120 },
  { title: 'Temas', key: 'topics', sortable: false, width: 200 },
  { title: 'Creado', key: 'createdAt', sortable: false, align: 'center' as const, width: 120 },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'center' as const, width: 100 },
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

    <template #[`item.source`]="{ item }">
      <div class="d-flex justify-center">
        <v-chip :color="SOURCE_COLORS[item.source]" size="x-small" variant="tonal">
          {{ SOURCE_LABELS[item.source] }}
        </v-chip>
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
          {{ DIFFICULTY_LABELS[item.difficulty] }}
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

    <template #[`item.actions`]="{ item }">
      <div class="d-flex ga-1 justify-center">
        <template v-if="canModify(item)">
          <v-btn
            icon="mdi-pencil"
            variant="text"
            size="small"
            color="primary"
            @click="emit('edit', item)"
          />
          <v-btn
            icon="mdi-delete"
            variant="text"
            size="small"
            color="error"
            @click="emit('delete', item)"
          />
        </template>
        <v-tooltip v-else :text="lockReason(item)">
          <template #activator="{ props }">
            <span v-bind="props">
              <v-btn
                icon="mdi-lock-outline"
                variant="text"
                size="small"
                color="disabled"
                disabled
              />
            </span>
          </template>
        </v-tooltip>
      </div>
    </template>

    <template #bottom>
      <slot name="footer" />
    </template>
  </v-data-table>
</template>
