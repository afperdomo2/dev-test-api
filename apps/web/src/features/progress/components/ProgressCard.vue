<script setup lang="ts">
import type { ProgressItem } from '@/types/progress.types'
import { DIFFICULTY_COLORS, TYPE_ICONS } from '@/types/question.types'
import { formatDate } from '@/utils/format'

interface Props {
  item: ProgressItem
  showToggle?: boolean
  toggling?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  select: [item: ProgressItem]
  toggle: [questionId: string]
}>()
</script>

<template>
  <v-card hover class="cursor-pointer" @click="emit('select', item)">
    <v-card-item>
      <template #prepend>
        <v-icon :icon="TYPE_ICONS[item.question.type]" color="primary" />
      </template>
      <v-card-title class="text-body-1">
        {{ item.question.content }}
      </v-card-title>
      <v-card-subtitle>
        <v-chip
          :color="DIFFICULTY_COLORS[item.question.difficulty]"
          size="x-small"
          variant="tonal"
          class="mr-1"
        >
          {{ item.question.difficulty }}
        </v-chip>
        <v-chip
          v-for="topic in item.question.topics"
          :key="topic"
          size="x-small"
          variant="outlined"
          class="mr-1"
        >
          {{ topic }}
        </v-chip>
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <div v-if="item.progress.nextReviewAt" class="text-caption text-medium-emphasis">
        Próximo repaso: {{ formatDate(item.progress.nextReviewAt) }}
      </div>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <v-btn
        v-if="showToggle"
        :icon="item.progress.isSaved ? 'mdi-bookmark' : 'mdi-bookmark-outline'"
        :color="item.progress.isSaved ? 'warning' : ''"
        variant="text"
        size="small"
        :loading="toggling"
        @click.stop="emit('toggle', item.question.id)"
      />
    </v-card-actions>
  </v-card>
</template>
