<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth.store'
import { getAiQuota, getImportQuota } from '@/api/services/questions.service'

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

function quotaPercent(used: number, limit: number): number {
  if (limit === 0) return 0
  return Math.min(100, Math.round((used / limit) * 100))
}

function quotaColor(percent: number): string {
  if (percent >= 90) return 'error'
  if (percent >= 70) return 'warning'
  return 'success'
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

      <v-col v-if="!isNotAdmin" cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="warning" size="40">mdi-calendar-clock</v-icon>
            </template>
            <v-card-title>Pendientes de repaso</v-card-title>
            <v-card-subtitle>Carga desde el backend</v-card-subtitle>
          </v-card-item>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
