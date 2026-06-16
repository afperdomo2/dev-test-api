<script setup lang="ts">
import type { UpcomingQuestion } from '@/types/progress.types'
import { DIFFICULTY_COLORS, TYPE_ICONS } from '@/types/question.types'
import { formatDate } from '@/utils/format'

interface Props {
  question: UpcomingQuestion
  showToggle?: boolean
  isSaved?: boolean
  toggling?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  toggle: [questionId: string]
}>()
</script>

<template>
  <v-card hover>
    <v-card-item>
      <template #prepend>
        <v-icon :icon="TYPE_ICONS[question.type]" color="primary" />
      </template>
      <v-card-title class="text-body-1">
        {{ question.content }}
      </v-card-title>
      <v-card-subtitle>
        <v-chip
          :color="DIFFICULTY_COLORS[question.difficulty]"
          size="x-small"
          variant="tonal"
          class="mr-1"
        >
          {{ question.difficulty }}
        </v-chip>
        <v-chip
          v-for="topic in question.topics"
          :key="topic.id"
          size="x-small"
          variant="outlined"
          class="mr-1"
        >
          {{ topic.name }}
        </v-chip>
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <div v-if="question.nextReviewAt" class="text-caption text-medium-emphasis">
        Próximo repaso: {{ formatDate(question.nextReviewAt) }}
      </div>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <v-btn
        v-if="showToggle"
        :icon="isSaved ? 'mdi-bookmark' : 'mdi-bookmark-outline'"
        :color="isSaved ? 'warning' : ''"
        variant="text"
        size="small"
        :loading="toggling"
        @click="emit('toggle', question.id)"
      />
    </v-card-actions>
  </v-card>
</template>
