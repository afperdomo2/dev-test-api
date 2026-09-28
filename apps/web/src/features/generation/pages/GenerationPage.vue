<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  createGenerationJobMutation,
  generationJobsListOptions,
  generationStatsOptions,
} from '@/queries/generation.queries'
import { usePagination } from '@/composables/usePagination'
import { useAppStore } from '@/stores/app.store'
import { formatShortDateTime } from '@/utils/format'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import GenerationEnqueueForm from '../components/GenerationEnqueueForm.vue'
import GenerationStatsCards from '../components/GenerationStatsCards.vue'
import type {
  CreateGenerationJobRequest,
  GenerationDifficulty,
  GenerationJob,
  GenerationJobStatus,
} from '@/types/generation.types'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

type StatusFilter = GenerationJobStatus | 'all'

const statusFilter = ref<StatusFilter>('all')

const statusItems: Array<{ title: string; value: StatusFilter }> = [
  { title: 'Todos', value: 'all' },
  { title: 'En cola', value: 'pending' },
  { title: 'Generando', value: 'running' },
  { title: 'Completados', value: 'completed' },
  { title: 'Fallidos', value: 'failed' },
]

const statusMeta: Record<GenerationJobStatus, { label: string; color: string }> = {
  pending: { label: 'En cola', color: 'warning' },
  running: { label: 'Generando', color: 'info' },
  completed: { label: 'Completado', color: 'success' },
  failed: { label: 'Fallido', color: 'error' },
}

const { data: stats, isLoading: statsLoading } = useQuery(generationStatsOptions())

const { data, isLoading } = useQuery(
  generationJobsListOptions(
    () => page.value,
    () => perPage.value,
    () => (statusFilter.value === 'all' ? undefined : statusFilter.value),
  ),
)

const createJob = useMutation(createGenerationJobMutation())
const formRef = ref<InstanceType<typeof GenerationEnqueueForm>>()

const jobList = computed<Array<GenerationJob>>(() => data.value?.data ?? [])
const totalJobs = computed(() => data.value?.meta?.total ?? 0)

const headers = [
  { title: 'Temas', key: 'topics', sortable: false },
  { title: 'Dificultad', key: 'difficulty', sortable: false },
  { title: 'Pedido / generado', key: 'quantity', sortable: false, align: 'center' as const },
  { title: 'Estado', key: 'status', sortable: false, align: 'center' as const },
  { title: 'Encolado', key: 'createdAt', sortable: false },
  { title: 'Detalle', key: 'error', sortable: false },
]

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  queryClient.invalidateQueries({ queryKey: ['generation'] })
}

watch(statusFilter, () => resetPagination())

function difficultyLabel(value: string) {
  const labels: Record<string, string> = {
    beginner: 'Principiante',
    intermediate: 'Intermedio',
    advanced: 'Avanzado',
    random: 'Aleatoria',
  }
  return labels[value] ?? value
}

async function onSubmit(payload: {
  topicIds: Array<string>
  difficulty: GenerationDifficulty
  quantity: number
}) {
  const body: CreateGenerationJobRequest = payload
  try {
    await createJob.mutateAsync(body)
    appStore.showSnackbar('Trabajo encolado')
    queryClient.invalidateQueries({ queryKey: ['generation'] })
    formRef.value?.reset()
  } catch (error) {
    const detail =
      error && typeof error === 'object' && 'detail' in error
        ? (error as { detail: string }).detail
        : 'Error al encolar el trabajo'
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Generación con IA"
      create-label=""
      :show-create="false"
      @refresh="handleRefresh"
    />

    <v-row class="mb-4" dense>
      <v-col cols="12" md="7">
        <GenerationEnqueueForm
          ref="formRef"
          :stats="stats"
          :loading="createJob.isPending.value"
          @submit="onSubmit"
        />
      </v-col>
      <v-col cols="12" md="5">
        <GenerationStatsCards :stats="stats" :loading="statsLoading" />
      </v-col>
    </v-row>

    <v-card>
      <v-card-text>
        <v-row class="mb-2" dense>
          <v-col cols="12" sm="4" md="3">
            <v-select
              v-model="statusFilter"
              :items="statusItems"
              label="Estado"
              hide-details
              density="compact"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <v-data-table
          :headers="headers"
          :items="jobList"
          :items-per-page="perPage"
          :loading="isLoading"
          no-data-text="No hay trabajos de generación"
          loading-text="Cargando trabajos..."
          hover
        >
          <template #[`item.topics`]="{ item }">
            <div class="d-flex flex-wrap ga-1">
              <v-chip v-for="topic in item.topics" :key="topic" size="small" variant="outlined">
                {{ topic }}
              </v-chip>
            </div>
          </template>

          <template #[`item.difficulty`]="{ item }">
            <span class="text-body-2">{{ difficultyLabel(item.difficulty) }}</span>
          </template>

          <template #[`item.quantity`]="{ item }">
            <span class="text-body-2">
              {{ item.quantity }} / <strong>{{ item.generated }}</strong>
            </span>
          </template>

          <template #[`item.status`]="{ item }">
            <v-chip size="small" variant="tonal" :color="statusMeta[item.status].color">
              {{ statusMeta[item.status].label }}
            </v-chip>
          </template>

          <template #[`item.createdAt`]="{ item }">
            <span class="text-caption">{{ formatShortDateTime(item.createdAt) }}</span>
          </template>

          <template #[`item.error`]="{ item }">
            <span class="text-caption text-error">{{ item.error || '—' }}</span>
          </template>

          <template #bottom>
            <PaginatedFooter
              :page="page"
              :per-page="perPage"
              :total="totalJobs"
              :in-table="true"
              @update:page="page = $event"
              @update:per-page="onPerPageChange"
            />
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>
</template>
