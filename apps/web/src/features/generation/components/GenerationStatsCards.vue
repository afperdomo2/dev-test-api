<script setup lang="ts">
import { computed } from 'vue'
import type { GenerationStats } from '@/types/generation.types'

interface Props {
  stats?: GenerationStats | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), { stats: null, loading: false })

interface StatusMeta {
  label: string
  color: string
  icon: string
}

type CounterKey = 'pending' | 'running' | 'failed'

const statusMeta: Record<CounterKey, StatusMeta> = {
  pending: { label: 'En cola', color: 'warning', icon: 'mdi-clock-outline' },
  running: { label: 'Generando', color: 'info', icon: 'mdi-robot-outline' },
  failed: { label: 'Fallidos', color: 'error', icon: 'mdi-alert-circle-outline' },
}

const usedPercent = computed(() => {
  if (!props.stats || props.stats.dailyCap <= 0) return 0
  return Math.min(100, Math.round((props.stats.reservedToday / props.stats.dailyCap) * 100))
})

const budgetColor = computed(() => {
  if (usedPercent.value >= 100) return 'error'
  if (usedPercent.value >= 80) return 'warning'
  return 'primary'
})

const counters = computed(() => {
  const stats = props.stats
  if (!stats) return []
  return (['pending', 'running', 'failed'] as const).map((key) => ({
    key,
    label: statusMeta[key].label,
    value: stats[key],
    color: statusMeta[key].color,
    icon: statusMeta[key].icon,
  }))
})
</script>

<template>
  <v-card>
    <v-card-item>
      <v-card-title class="text-subtitle-1">Consumo de hoy</v-card-title>
      <v-card-subtitle v-if="stats">
        {{ stats.reservedToday }} de {{ stats.dailyCap }} reservadas · {{ stats.generatedToday }}
        generadas
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <v-progress-linear
        :model-value="usedPercent"
        :color="budgetColor"
        height="8"
        rounded
        class="mb-4"
      />

      <v-row dense>
        <v-col v-for="counter in counters" :key="counter.key" cols="4">
          <div class="d-flex align-center ga-2">
            <v-icon :color="counter.color" :icon="counter.icon" size="small" />
            <div>
              <div class="text-body-2 font-weight-medium">
                {{ loading ? '—' : counter.value }}
              </div>
              <div class="text-caption text-medium-emphasis">{{ counter.label }}</div>
            </div>
          </div>
        </v-col>
      </v-row>

      <v-alert v-if="stats && !stats.aiConfigured" type="warning" variant="tonal" density="compact">
        La IA no está configurada. Define <code>AI_API_URL</code> y <code>AI_API_KEY</code> en el
        servidor para encolar trabajos.
      </v-alert>
    </v-card-text>
  </v-card>
</template>
