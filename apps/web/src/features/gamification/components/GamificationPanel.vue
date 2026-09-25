<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { overviewOptions, freezeStreakMutation } from '@/queries/gamification.queries'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import DailyGoalDialog from './DailyGoalDialog.vue'
import TopicXpList from './TopicXpList.vue'
import type { ApiError } from '@/types/api.types'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { extractFieldErrors } = useFormErrors()

const { data, isLoading, error, refetch } = useQuery(overviewOptions())

const freezeMut = useMutation(freezeStreakMutation())
const goalDialog = ref(false)

const overview = computed(() => data.value ?? null)

const dailyGoalColor = computed(() => {
  const goal = overview.value?.dailyGoal
  if (!goal) return 'primary'
  if (goal.reached) return 'success'
  if (goal.progress >= 70) return 'warning'
  return 'primary'
})

const freezeDisabled = computed(() => {
  const streak = overview.value?.streak
  if (!streak) return true
  return streak.freezesAvailable <= 0
})

/**
 * El backend deriva el nivel de `totalXp` con `1 + sqrt(xp/100)`, es decir
 * el umbral del nivel N es `(N-1)^2 * 100`. Replicamos la misma fórmula para
 * pintar el progreso dentro del nivel actual sin pedir otro endpoint.
 */
const levelPercent = computed(() => {
  const current = overview.value
  if (!current || current.xpToLevel <= 0) return 100
  const levelSpan = (2 * current.level - 1) * 100
  if (levelSpan <= 0) return 100
  const currentLevelXp = levelSpan - current.xpToLevel
  return Math.max(0, Math.min(100, Math.round((currentLevelXp / levelSpan) * 100)))
})

function errorMessage(err: unknown): string {
  const e = err as ApiError
  if (e?.detail) return e.detail
  if (e?.title) return e.title
  return 'No se pudo cargar la gamificación'
}

function invalidate() {
  queryClient.invalidateQueries({ queryKey: ['gamification'] })
}

async function handleFreeze() {
  try {
    await freezeMut.mutateAsync()
    appStore.showSnackbar('Freeze aplicado a tu racha')
    invalidate()
  } catch (err: unknown) {
    const fieldErrors = extractFieldErrors(err)
    const detail =
      Object.values(fieldErrors)[0] ??
      (err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al aplicar el freeze')
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-card variant="flat" border>
    <v-card-title class="d-flex align-center ga-2">
      <v-icon icon="mdi-trophy-variant" color="primary" />
      Gamificación
      <v-spacer />
      <v-btn
        icon="mdi-refresh"
        variant="text"
        size="small"
        :loading="isLoading"
        aria-label="Refrescar gamificación"
        @click="refetch()"
      />
    </v-card-title>

    <v-card-text v-if="error" class="py-4">
      <v-alert type="warning" variant="tonal" density="compact">
        {{ errorMessage(error) }}
      </v-alert>
    </v-card-text>

    <v-card-text v-else-if="isLoading || !overview" class="py-4">
      <v-skeleton-loader type="list-item-three-line@4" />
    </v-card-text>

    <v-card-text v-else>
      <v-row>
        <v-col cols="12" md="6">
          <v-row dense>
            <v-col cols="6">
              <div class="text-caption text-medium-emphasis">XP total</div>
              <div class="text-h5 font-weight-bold text-primary">
                {{ overview.totalXp.toLocaleString('es-CO') }}
              </div>
            </v-col>
            <v-col cols="6">
              <div class="text-caption text-medium-emphasis">Nivel</div>
              <div class="text-h5 font-weight-bold">{{ overview.level }}</div>
            </v-col>
            <v-col cols="12">
              <v-progress-linear :model-value="levelPercent" color="primary" height="10" rounded>
                <template #default>
                  <span class="text-caption font-weight-medium">
                    {{ overview.xpToLevel }} XP para el nivel {{ overview.level + 1 }}
                  </span>
                </template>
              </v-progress-linear>
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12" md="6">
          <div class="d-flex align-center ga-2 mb-1">
            <v-icon icon="mdi-fire" size="20" color="orange" />
            <span class="text-subtitle-2">Racha</span>
          </div>
          <div class="d-flex align-center ga-2 flex-wrap mb-2">
            <v-chip color="orange" variant="tonal" size="small">
              {{ overview.streak.current }} días
            </v-chip>
            <v-chip variant="tonal" size="small">Récord: {{ overview.streak.longest }}</v-chip>
            <v-chip :color="freezeDisabled ? undefined : 'info'" variant="tonal" size="small">
              Freezes: {{ overview.streak.freezesAvailable }}
            </v-chip>
          </div>
          <v-btn
            size="small"
            variant="tonal"
            color="info"
            prepend-icon="mdi-snowflake"
            :disabled="freezeDisabled"
            :loading="freezeMut.isPending.value"
            @click="handleFreeze"
          >
            Usar freeze
          </v-btn>
        </v-col>
      </v-row>

      <v-divider class="my-4" />

      <div class="d-flex align-center ga-2 mb-1">
        <v-icon icon="mdi-target" size="20" color="primary" />
        <span class="text-subtitle-2">Meta diaria</span>
        <v-spacer />
        <v-chip
          v-if="overview.dailyGoal.reached"
          color="success"
          variant="tonal"
          size="small"
          prepend-icon="mdi-check-circle"
        >
          Completada
        </v-chip>
        <v-btn
          size="small"
          variant="text"
          color="primary"
          prepend-icon="mdi-pencil"
          @click="goalDialog = true"
        >
          Editar
        </v-btn>
      </div>
      <v-progress-linear
        :model-value="overview.dailyGoal.progress"
        :color="dailyGoalColor"
        height="14"
        rounded
      >
        <template #default>
          <span class="text-caption font-weight-medium">
            {{ overview.dailyGoal.current }} / {{ overview.dailyGoal.target }} XP
          </span>
        </template>
      </v-progress-linear>

      <v-divider class="my-4" />

      <TopicXpList :items="overview.topicXp" />
    </v-card-text>

    <DailyGoalDialog
      v-if="overview"
      v-model="goalDialog"
      :current-target="overview.dailyGoal.target"
      @saved="invalidate()"
    />
  </v-card>
</template>
