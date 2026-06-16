<script setup lang="ts">
import type { Question } from '@/types/question.types'
import { DIFFICULTY_COLORS, TYPE_ICONS } from '@/types/question.types'
import { formatDate } from '@/utils/format'

interface Props {
  question: Question
}

defineProps<Props>()

function typeLabel(type: Question['type']): string {
  const labels: Record<Question['type'], string> = {
    single_choice: 'Sel. única',
    multiple_choice: 'Sel. múltiple',
    code_completion: 'Código',
  }
  return labels[type]
}

function truncate(text: string, max = 120): string {
  return text.length > max ? text.slice(0, max) + '...' : text
}
</script>

<template>
  <v-card :to="`/questions/${question.id}`" hover>
    <v-card-item>
      <template #prepend>
        <v-icon :icon="TYPE_ICONS[question.type]" color="primary" />
      </template>

      <v-card-title class="text-body-1 font-weight-bold">
        {{ truncate(question.content) }}
      </v-card-title>

      <v-card-subtitle class="mt-1">
        <v-chip
          :color="DIFFICULTY_COLORS[question.difficulty]"
          size="x-small"
          variant="tonal"
          class="mr-1"
        >
          {{ question.difficulty }}
        </v-chip>
        <v-chip size="x-small" variant="text">
          {{ typeLabel(question.type) }}
        </v-chip>
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <div v-if="question.topics.length" class="d-flex flex-wrap ga-1">
        <v-chip v-for="topic in question.topics" :key="topic.id" size="x-small" variant="outlined">
          {{ topic.name }}
        </v-chip>
      </div>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <span class="text-caption text-medium-emphasis">
        {{ formatDate(question.createdAt) }}
      </span>
    </v-card-actions>
  </v-card>
</template>
