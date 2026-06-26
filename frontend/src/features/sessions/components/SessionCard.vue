<script setup lang="ts">
import { computed } from 'vue'
import type { Session } from '@/types/session.types'
import {
  SESSION_STATUS_LABELS,
  SESSION_STATUS_COLORS,
  SESSION_MODE_LABELS,
} from '@/types/session.types'
import { formatDate, formatScore } from '@/utils/format'

interface Props {
  session: Session
}

const props = defineProps<Props>()

const isGenerating = computed(
  () =>
    props.session.mode === 'generate' &&
    props.session.status === 'in_progress' &&
    props.session.questionsGenerated < 2,
)
</script>

<template>
  <v-card :to="`/sessions/${session.id}/study`" hover>
    <v-card-item>
      <template #prepend>
        <v-icon color="primary" size="32">
          {{ session.mode === 'generate' ? 'mdi-play-circle' : 'mdi-refresh-circle' }}
        </v-icon>
      </template>
      <v-card-title>{{ session.name }}</v-card-title>
      <v-card-subtitle>
        <v-chip
          :color="SESSION_STATUS_COLORS[session.status]"
          size="x-small"
          variant="tonal"
          class="mr-1"
        >
          {{ SESSION_STATUS_LABELS[session.status] }}
        </v-chip>
        <v-chip size="x-small" variant="tonal" class="mr-1">
          {{ SESSION_MODE_LABELS[session.mode] }}
        </v-chip>
        <v-chip size="x-small" variant="text">
          {{ session.difficulty }}
        </v-chip>
        <v-chip v-if="isGenerating" size="x-small" variant="flat" color="warning" class="ml-1">
          <v-icon start size="14">mdi-cog</v-icon>
          Generando...
        </v-chip>
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <div class="d-flex align-center ga-4">
        <div>
          <span class="text-caption text-medium-emphasis">Puntuación</span>
          <div class="text-body-2 font-weight-medium">
            {{ session.status === 'completed' ? formatScore(session.score) : '—' }}
          </div>
        </div>
        <div>
          <span class="text-caption text-medium-emphasis">Respuestas</span>
          <div class="text-body-2 font-weight-medium">{{ session.answerCount }}</div>
        </div>
        <div v-if="session.topics.length">
          <span class="text-caption text-medium-emphasis">Temas</span>
          <div class="text-body-2">{{ session.topics.length }}</div>
        </div>
      </div>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <span class="text-caption text-medium-emphasis">
        {{ formatDate(session.createdAt) }}
      </span>
    </v-card-actions>
  </v-card>
</template>
