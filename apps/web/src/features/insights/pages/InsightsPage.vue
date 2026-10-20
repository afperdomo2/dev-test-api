<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { insightSummaryOptions } from '@/queries/insights.queries'
import ListPageHeader from '@/components/ListPageHeader.vue'
import { PERIOD_OPTIONS } from '@/types/insight.types'
import type { InsightPeriod } from '@/types/insight.types'

const queryClient = useQueryClient()
const period = ref<InsightPeriod>('week')

const { data, isLoading, error } = useQuery(insightSummaryOptions(() => period.value))

const summary = computed(() => data.value ?? null)

const kpis = computed(() => {
  const s = summary.value
  if (!s) return []

  return [
    { label: 'Sesiones', value: String(s.sessions), icon: 'mdi-play-circle-outline' },
    { label: 'Preguntas', value: String(s.questionsAnswered), icon: 'mdi-help-circle-outline' },
    { label: 'Acierto', value: `${percent(s.accuracy)}%`, icon: 'mdi-target' },
    { label: 'Tiempo', value: formatDuration(s.timeSpentMs), icon: 'mdi-timer-outline' },
    { label: 'XP ganado', value: String(s.xpGained), icon: 'mdi-star-outline' },
    { label: 'Racha', value: `${s.streakDays} días`, icon: 'mdi-fire' },
  ]
})

const maxDaily = computed(() => {
  const values = summary.value?.dailyActivity.map((day) => day.questionsAnswered) ?? []
  return Math.max(1, ...values)
})

/** La API devuelve `accuracy` en 0..1. */
function percent(value: number): number {
  return Math.round(value * 100)
}

function formatDuration(ms: number): string {
  const totalMinutes = Math.round(ms / 60000)
  if (totalMinutes < 60) return `${totalMinutes} min`

  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
}

function dayLabel(date: string): string {
  // `date` es YYYY-MM-DD (UTC); forzamos UTC para que no se desplace un día.
  return new Date(`${date}T00:00:00Z`).toLocaleDateString('es-ES', {
    weekday: 'short',
    timeZone: 'UTC',
  })
}

function handleRefresh() {
  queryClient.invalidateQueries({ queryKey: ['insights', 'summary'] })
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Informes"
      create-label=""
      :show-create="false"
      @refresh="handleRefresh"
    />

    <div class="d-flex flex-wrap align-center ga-3 mb-4">
      <v-btn-toggle v-model="period" mandatory color="primary" density="comfortable">
        <v-btn v-for="option in PERIOD_OPTIONS" :key="option.value" :value="option.value">
          {{ option.title }}
        </v-btn>
      </v-btn-toggle>

      <span v-if="summary" class="text-caption text-medium-emphasis">
        {{ summary.from.slice(0, 10) }} — {{ summary.to.slice(0, 10) }}
      </span>
    </div>

    <v-skeleton-loader v-if="isLoading" type="card@2" />

    <v-alert v-else-if="error" type="error" variant="tonal" class="mb-4">
      No se pudo cargar el informe.
      <template #append>
        <v-btn variant="text" @click="handleRefresh">Reintentar</v-btn>
      </template>
    </v-alert>

    <template v-else-if="summary">
      <v-alert
        v-if="summary.sessions === 0 && summary.questionsAnswered === 0"
        type="info"
        variant="tonal"
        class="mb-4"
      >
        No hay actividad en este periodo.
      </v-alert>

      <v-row>
        <v-col v-for="kpi in kpis" :key="kpi.label" cols="6" md="4">
          <v-card variant="flat" border class="h-100">
            <v-card-text>
              <div class="d-flex align-center ga-2 mb-1">
                <v-icon :icon="kpi.icon" color="primary" size="18" />
                <span class="text-caption text-medium-emphasis">{{ kpi.label }}</span>
              </div>
              <div class="text-h6">{{ kpi.value }}</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row class="mt-2">
        <v-col cols="12" md="6">
          <v-card variant="flat" border class="h-100">
            <v-card-title>Más practicados</v-card-title>
            <v-card-text>
              <p v-if="!summary.topTopics.length" class="text-body-2 text-medium-emphasis">
                Sin datos en este periodo.
              </p>

              <div v-for="topic in summary.topTopics" :key="topic.topicId" class="mb-3">
                <div class="d-flex justify-space-between text-body-2 mb-1">
                  <span>{{ topic.topicName }}</span>
                  <span class="text-medium-emphasis">
                    {{ percent(topic.accuracy) }}% · {{ topic.attempts }}
                  </span>
                </div>
                <v-progress-linear
                  :model-value="percent(topic.accuracy)"
                  height="8"
                  rounded
                  color="primary"
                  role="progressbar"
                  :aria-label="`${topic.topicName}: ${percent(topic.accuracy)}% de acierto`"
                />
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col v-if="summary.weakTopics.length" cols="12" md="6">
          <v-card variant="flat" border class="h-100">
            <v-card-title>Temas débiles</v-card-title>
            <v-card-text>
              <div v-for="topic in summary.weakTopics" :key="topic.topicId" class="mb-3">
                <div class="d-flex justify-space-between align-center text-body-2 mb-1">
                  <span>{{ topic.topicName }}</span>
                  <v-chip
                    size="small"
                    color="warning"
                    variant="tonal"
                    prepend-icon="mdi-alert-outline"
                  >
                    Débil
                  </v-chip>
                </div>
                <v-progress-linear
                  :model-value="percent(topic.accuracy)"
                  height="8"
                  rounded
                  color="warning"
                  role="progressbar"
                  :aria-label="`${topic.topicName}: ${percent(topic.accuracy)}% de acierto`"
                />
                <div class="text-caption text-medium-emphasis mt-1">
                  {{ topic.correctCount }}/{{ topic.attempts }} aciertos
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-card variant="flat" border class="mt-4">
        <v-card-title>Actividad diaria</v-card-title>
        <v-card-text>
          <div class="d-flex align-end ga-2" style="height: 140px">
            <div
              v-for="day in summary.dailyActivity"
              :key="day.date"
              class="d-flex flex-column align-center flex-grow-1"
            >
              <div
                class="rounded w-100"
                role="img"
                :aria-label="`${day.date}: ${day.questionsAnswered} preguntas`"
                :style="{
                  height: `${(day.questionsAnswered / maxDaily) * 90}px`,
                  backgroundColor: day.questionsAnswered
                    ? 'rgb(var(--v-theme-primary))'
                    : 'rgba(255, 255, 255, 0.08)',
                }"
              />
              <span class="text-caption text-medium-emphasis mt-1">{{ dayLabel(day.date) }}</span>
            </div>
          </div>
        </v-card-text>
      </v-card>
    </template>
  </v-container>
</template>
