<script setup lang="ts">
import { computed } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { badgesOptions } from '@/queries/gamification.queries'
import ListPageHeader from '@/components/ListPageHeader.vue'

const queryClient = useQueryClient()
const { data, isLoading, error } = useQuery(badgesOptions())

const badges = computed(() => data.value?.badges ?? [])
const earned = computed(() => data.value?.earned ?? 0)
const total = computed(() => data.value?.total ?? 0)
const progress = computed(() =>
  total.value > 0 ? Math.round((earned.value / total.value) * 100) : 0,
)

function handleRefresh() {
  queryClient.invalidateQueries({ queryKey: ['gamification', 'badges'] })
}
</script>

<template>
  <v-container>
    <ListPageHeader title="Logros" create-label="" :show-create="false" @refresh="handleRefresh" />

    <v-skeleton-loader v-if="isLoading" type="card@3" />

    <v-alert v-else-if="error" type="error" variant="tonal">
      No se pudieron cargar los logros.
      <template #append>
        <v-btn variant="text" @click="handleRefresh">Reintentar</v-btn>
      </template>
    </v-alert>

    <template v-else>
      <v-card variant="flat" border class="mb-4">
        <v-card-text>
          <div class="d-flex justify-space-between mb-2">
            <span class="text-body-2">{{ earned }} de {{ total }} logros</span>
            <span class="text-caption text-medium-emphasis">{{ progress }}%</span>
          </div>
          <v-progress-linear
            :model-value="progress"
            height="8"
            rounded
            color="primary"
            role="progressbar"
            :aria-label="`${earned} de ${total} logros`"
          />
        </v-card-text>
      </v-card>

      <v-row>
        <v-col v-for="badge in badges" :key="badge.id" cols="12" sm="6" md="4">
          <v-card variant="flat" border class="h-100" :style="{ opacity: badge.earned ? 1 : 0.75 }">
            <v-card-text>
              <div class="d-flex align-start ga-3 mb-3">
                <v-icon :icon="badge.icon" :color="badge.earned ? 'primary' : 'grey'" size="32" />
                <div>
                  <div class="text-body-1 font-weight-medium">{{ badge.name }}</div>
                  <div class="text-caption text-medium-emphasis">{{ badge.description }}</div>
                </div>
              </div>

              <div v-if="badge.earned" class="d-flex align-center ga-1 text-caption">
                <v-icon icon="mdi-check-circle" color="primary" size="16" />
                <span>Conseguido</span>
              </div>

              <template v-else>
                <v-progress-linear
                  :model-value="badge.target > 0 ? (badge.current / badge.target) * 100 : 0"
                  height="6"
                  rounded
                  color="grey"
                  class="mb-1"
                  role="progressbar"
                  :aria-label="`${badge.name}: ${badge.current} de ${badge.target}`"
                />
                <div class="text-caption text-medium-emphasis">
                  {{ badge.current }}/{{ badge.target }}
                </div>
              </template>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>
