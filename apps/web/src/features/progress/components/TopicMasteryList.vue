<script setup lang="ts">
import type { TopicMasteryItem } from '@/types/progress.types'

interface Props {
  items: Array<TopicMasteryItem>
  loading?: boolean
}

defineProps<Props>()

/** `accuracy` llega como 0..1; la barra trabaja en porcentaje entero. */
function percent(accuracy: number): number {
  return Math.round(accuracy * 100)
}
</script>

<template>
  <v-card variant="flat" border>
    <v-card-title class="d-flex align-center ga-2">
      <v-icon icon="mdi-chart-timeline-variant" color="primary" />
      Dominio por tema
    </v-card-title>

    <v-card-text v-if="loading">
      <v-skeleton-loader type="list-item-two-line@4" />
    </v-card-text>

    <v-card-text v-else-if="!items.length" class="text-center py-8">
      <v-icon size="48" color="grey-lighten-1" class="mb-2">mdi-chart-timeline-variant</v-icon>
      <p class="text-body-1 text-medium-emphasis">
        Aún no hay respuestas registradas. Responde preguntas para ver tu dominio por tema.
      </p>
    </v-card-text>

    <v-list v-else density="comfortable">
      <v-list-item v-for="item in items" :key="item.topicId">
        <div class="d-flex align-center justify-space-between ga-2">
          <span class="text-body-1">{{ item.topicName }}</span>
          <v-chip
            v-if="item.isWeak"
            size="small"
            color="warning"
            variant="tonal"
            prepend-icon="mdi-alert-outline"
          >
            Débil
          </v-chip>
        </div>

        <div class="d-flex align-center ga-3 mt-2">
          <v-progress-linear
            :model-value="percent(item.accuracy)"
            :color="item.isWeak ? 'warning' : 'primary'"
            height="8"
            rounded
            class="flex-grow-1"
            role="progressbar"
            :aria-label="`${item.topicName}: ${percent(item.accuracy)}% de acierto`"
            :aria-valuenow="percent(item.accuracy)"
            aria-valuemin="0"
            aria-valuemax="100"
          />
          <span class="text-caption text-medium-emphasis text-no-wrap">
            {{ item.correctCount }}/{{ item.attempts }}
          </span>
        </div>
      </v-list-item>
    </v-list>
  </v-card>
</template>
