<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth.store'
import { getAiQuota, getImportQuota } from '@/api/services/questions.service'
import { getOpenCodeUsage } from '@/api/services/opencode.service'
import { formatShortDateTime, formatTimeUntil } from '@/utils/format'
import type { ApiError } from '@/types/api.types'

const authStore = useAuthStore()

const formattedDate = computed(() => {
  const date = authStore.user?.createdAt
  if (!date) return ''
  return new Date(date).toLocaleDateString('es-CO', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

const isNotAdmin = computed(() => !authStore.isAdmin)
const isAdmin = computed(() => authStore.isAdmin)

const { data: aiQuota } = useQuery({
  queryKey: ['questions', 'ai-quota'],
  queryFn: () => getAiQuota(),
  staleTime: 30 * 1000,
  enabled: isNotAdmin,
})

const { data: importQuota } = useQuery({
  queryKey: ['questions', 'import-quota'],
  queryFn: () => getImportQuota(),
  staleTime: 30 * 1000,
  enabled: isNotAdmin,
})

const {
  data: openCodeUsage,
  isLoading: isOpenCodeLoading,
  error: openCodeError,
} = useQuery({
  queryKey: ['opencode', 'usage'],
  queryFn: () => getOpenCodeUsage(),
  staleTime: 60 * 1000,
  enabled: isAdmin,
})

function quotaPercent(used: number, limit: number): number {
  if (limit === 0) return 0
  return Math.min(100, Math.round((used / limit) * 100))
}

function quotaColor(percent: number): string {
  if (percent >= 90) return 'error'
  if (percent >= 70) return 'warning'
  return 'success'
}

function openCodeErrorMessage(err: unknown): string {
  const e = err as ApiError
  if (e?.detail) return e.detail
  if (e?.title) return e.title
  return 'No se pudo cargar el consumo de OpenCode'
}
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Dashboard</h1>

    <v-row>
      <v-col cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="primary" size="40">mdi-account-circle</v-icon>
            </template>
            <v-card-title class="text-h5">
              {{ authStore.user?.email }}
            </v-card-title>
            <v-card-subtitle class="mt-1"> Miembro desde {{ formattedDate }} </v-card-subtitle>
          </v-card-item>
          <v-card-text>
            <v-chip
              :color="authStore.isAdmin ? 'error' : 'info'"
              size="small"
              variant="tonal"
              class="font-weight-bold"
            >
              <v-icon start size="18">
                {{ authStore.isAdmin ? 'mdi-shield-account' : 'mdi-account' }}
              </v-icon>
              {{ authStore.isAdmin ? 'Administrador' : 'Usuario' }}
            </v-chip>
            <div class="d-flex align-center ga-2 mt-3 text-caption text-medium-emphasis">
              <v-icon size="16" color="success">mdi-identifier</v-icon>
              <span class="text-truncate">{{ authStore.user?.id }}</span>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col v-if="isNotAdmin" cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="deep-purple" size="40">mdi-robot</v-icon>
            </template>
            <v-card-title>Cupo IA hoy</v-card-title>
            <v-card-subtitle>Preguntas generadas con IA</v-card-subtitle>
          </v-card-item>
          <v-card-text v-if="aiQuota">
            <div class="d-flex align-center justify-space-between mb-2">
              <span class="text-body-2 font-weight-medium">
                {{ aiQuota.usedToday }} / {{ aiQuota.dailyLimit }}
                <span class="text-medium-emphasis">usadas</span>
              </span>
              <v-chip
                :color="quotaColor(quotaPercent(aiQuota.usedToday, aiQuota.dailyLimit))"
                size="small"
                variant="tonal"
              >
                Restan {{ aiQuota.remaining }}
              </v-chip>
            </div>
            <v-progress-linear
              :model-value="quotaPercent(aiQuota.usedToday, aiQuota.dailyLimit)"
              :color="quotaColor(quotaPercent(aiQuota.usedToday, aiQuota.dailyLimit))"
              height="20"
              rounded
            >
              <template #default>
                <span class="text-caption font-weight-medium">
                  {{ quotaPercent(aiQuota.usedToday, aiQuota.dailyLimit) }}%
                </span>
              </template>
            </v-progress-linear>
            <div class="text-caption text-medium-emphasis mt-2">
              Límite diario configurable por el administrador (1–50)
            </div>
          </v-card-text>
          <v-card-text v-else>
            <v-skeleton-loader type="text" />
          </v-card-text>
        </v-card>
      </v-col>

      <v-col v-if="isNotAdmin" cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="orange" size="40">mdi-file-upload</v-icon>
            </template>
            <v-card-title>Cupo importación hoy</v-card-title>
            <v-card-subtitle>Preguntas subidas por CSV</v-card-subtitle>
          </v-card-item>
          <v-card-text v-if="importQuota">
            <div class="d-flex align-center justify-space-between mb-2">
              <span class="text-body-2 font-weight-medium">
                {{ importQuota.usedToday }} / {{ importQuota.dailyLimit }}
                <span class="text-medium-emphasis">usadas</span>
              </span>
              <v-chip
                :color="quotaColor(quotaPercent(importQuota.usedToday, importQuota.dailyLimit))"
                size="small"
                variant="tonal"
              >
                Restan {{ importQuota.remaining }}
              </v-chip>
            </div>
            <v-progress-linear
              :model-value="quotaPercent(importQuota.usedToday, importQuota.dailyLimit)"
              :color="quotaColor(quotaPercent(importQuota.usedToday, importQuota.dailyLimit))"
              height="20"
              rounded
            >
              <template #default>
                <span class="text-caption font-weight-medium">
                  {{ quotaPercent(importQuota.usedToday, importQuota.dailyLimit) }}%
                </span>
              </template>
            </v-progress-linear>
            <div class="text-caption text-medium-emphasis mt-2">
              Máximo por archivo: 50 preguntas
            </div>
          </v-card-text>
          <v-card-text v-else>
            <v-skeleton-loader type="text" />
          </v-card-text>
        </v-card>
      </v-col>

      <v-col v-if="isAdmin" cols="12" md="6" lg="6">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="primary" size="40">mdi-cloud-percent</v-icon>
            </template>
            <v-card-title>Consumo OpenCode Go</v-card-title>
            <v-card-subtitle>Rolling (5h) · Semanal · Mensual — vía API key</v-card-subtitle>
          </v-card-item>

          <v-card-text v-if="isOpenCodeLoading">
            <v-skeleton-loader type="article" />
          </v-card-text>

          <v-card-text v-else-if="openCodeError">
            <v-alert type="warning" variant="tonal" density="compact">
              {{ openCodeErrorMessage(openCodeError) }}
            </v-alert>
          </v-card-text>

          <v-card-text v-else-if="openCodeUsage">
            <div class="d-flex flex-column ga-4">
              <div>
                <div class="d-flex align-center justify-space-between mb-1">
                  <span class="text-body-2 font-weight-medium">Rolling (5 h)</span>
                  <v-chip
                    :color="quotaColor(openCodeUsage.rolling.percent)"
                    size="small"
                    variant="tonal"
                  >
                    {{ Math.round(openCodeUsage.rolling.percent) }}% usado
                  </v-chip>
                </div>
                <v-progress-linear
                  :model-value="openCodeUsage.rolling.percent"
                  :color="quotaColor(openCodeUsage.rolling.percent)"
                  height="18"
                  rounded
                />
                <div
                  class="text-caption text-medium-emphasis mt-1 d-flex align-center ga-1 flex-wrap"
                >
                  <v-tooltip
                    :text="formatShortDateTime(openCodeUsage.rolling.resetsAt)"
                    location="top"
                  >
                    <template #activator="{ props }">
                      <span
                        v-bind="props"
                        style="cursor: help; border-bottom: 1px dotted currentColor"
                      >
                        {{ formatTimeUntil(openCodeUsage.rolling.resetsAt) }}
                      </span>
                    </template>
                  </v-tooltip>
                  <span v-if="openCodeUsage.rolling.status">
                    · {{ openCodeUsage.rolling.status }}</span
                  >
                </div>
              </div>

              <div>
                <div class="d-flex align-center justify-space-between mb-1">
                  <span class="text-body-2 font-weight-medium">Semanal</span>
                  <v-chip
                    :color="quotaColor(openCodeUsage.weekly.percent)"
                    size="small"
                    variant="tonal"
                  >
                    {{ Math.round(openCodeUsage.weekly.percent) }}% usado
                  </v-chip>
                </div>
                <v-progress-linear
                  :model-value="openCodeUsage.weekly.percent"
                  :color="quotaColor(openCodeUsage.weekly.percent)"
                  height="18"
                  rounded
                />
                <div
                  class="text-caption text-medium-emphasis mt-1 d-flex align-center ga-1 flex-wrap"
                >
                  <v-tooltip
                    :text="formatShortDateTime(openCodeUsage.weekly.resetsAt)"
                    location="top"
                  >
                    <template #activator="{ props }">
                      <span
                        v-bind="props"
                        style="cursor: help; border-bottom: 1px dotted currentColor"
                      >
                        {{ formatTimeUntil(openCodeUsage.weekly.resetsAt) }}
                      </span>
                    </template>
                  </v-tooltip>
                  <span v-if="openCodeUsage.weekly.status">
                    · {{ openCodeUsage.weekly.status }}</span
                  >
                </div>
              </div>

              <div>
                <div class="d-flex align-center justify-space-between mb-1">
                  <span class="text-body-2 font-weight-medium">Mensual</span>
                  <v-chip
                    :color="quotaColor(openCodeUsage.monthly.percent)"
                    size="small"
                    variant="tonal"
                  >
                    {{ Math.round(openCodeUsage.monthly.percent) }}% usado
                  </v-chip>
                </div>
                <v-progress-linear
                  :model-value="openCodeUsage.monthly.percent"
                  :color="quotaColor(openCodeUsage.monthly.percent)"
                  height="18"
                  rounded
                />
                <div
                  class="text-caption text-medium-emphasis mt-1 d-flex align-center ga-1 flex-wrap"
                >
                  <v-tooltip
                    :text="formatShortDateTime(openCodeUsage.monthly.resetsAt)"
                    location="top"
                  >
                    <template #activator="{ props }">
                      <span
                        v-bind="props"
                        style="cursor: help; border-bottom: 1px dotted currentColor"
                      >
                        {{ formatTimeUntil(openCodeUsage.monthly.resetsAt) }}
                      </span>
                    </template>
                  </v-tooltip>
                  <span v-if="openCodeUsage.monthly.status">
                    · {{ openCodeUsage.monthly.status }}</span
                  >
                </div>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
