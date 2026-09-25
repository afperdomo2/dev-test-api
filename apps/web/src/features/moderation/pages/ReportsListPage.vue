<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { reportsListOptions } from '@/queries/moderation.queries'
import { usePagination } from '@/composables/usePagination'
import { formatShortDateTime } from '@/utils/format'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import ReportReviewDialog from '../components/ReportReviewDialog.vue'
import type { Report, ReportStatus } from '@/types/moderation.types'

const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

type StatusFilter = ReportStatus | 'all'

const statusFilter = ref<StatusFilter>('all')

const statusItems: Array<{ title: string; value: StatusFilter }> = [
  { title: 'Todos', value: 'all' },
  { title: 'Pendientes', value: 'pending' },
  { title: 'Revisados', value: 'reviewed' },
]

const { data, isLoading } = useQuery(
  reportsListOptions(
    () => page.value,
    () => perPage.value,
    () => (statusFilter.value === 'all' ? undefined : statusFilter.value),
  ),
)

const reviewDialog = ref(false)
const selectedReport = ref<Report | null>(null)

const reportList = computed<Array<Report>>(() => data.value?.data ?? [])
const totalReports = computed(() => data.value?.meta?.total ?? 0)

const headers = [
  { title: 'Motivo', key: 'reason', sortable: false },
  { title: 'Tipo', key: 'targetType', sortable: false, align: 'center' as const },
  { title: 'Estado', key: 'status', sortable: false, align: 'center' as const },
  { title: 'Creado', key: 'createdAt', sortable: false },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'center' as const },
]

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['reports', 'list'] })
}

watch(statusFilter, () => {
  resetPagination()
})

function openReview(report: Report) {
  selectedReport.value = report
  reviewDialog.value = true
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Moderación"
      create-label=""
      :show-create="false"
      @refresh="handleRefresh"
    />

    <v-alert type="info" variant="tonal" density="compact" class="mb-4">
      Revisa los reportes enviados por los usuarios. Des-publicar solo aplica a preguntas públicas;
      el contenido generado con IA requiere un mecanismo adicional fuera de alcance.
    </v-alert>

    <v-row class="mb-4" dense>
      <v-col cols="12" md="4" lg="3">
        <v-select
          v-model="statusFilter"
          label="Estado"
          :items="statusItems"
          hide-details
          density="compact"
          variant="outlined"
        />
      </v-col>
    </v-row>

    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="reportList"
          :items-per-page="perPage"
          :loading="isLoading"
          no-data-text="No hay reportes"
          loading-text="Cargando reportes..."
          hover
        >
          <template #[`item.reason`]="{ item }">
            <span class="text-body-2">{{ item.reason }}</span>
          </template>

          <template #[`item.targetType`]="{ item }">
            <v-chip
              size="small"
              variant="tonal"
              :color="item.targetType === 'question' ? 'primary' : 'info'"
            >
              {{ item.targetType === 'question' ? 'Pregunta' : 'Tip' }}
            </v-chip>
          </template>

          <template #[`item.status`]="{ item }">
            <v-chip
              size="small"
              variant="tonal"
              :color="item.status === 'pending' ? 'warning' : 'success'"
            >
              {{ item.status === 'pending' ? 'Pendiente' : 'Revisado' }}
            </v-chip>
          </template>

          <template #[`item.createdAt`]="{ item }">
            <span class="text-caption">{{ formatShortDateTime(item.createdAt) }}</span>
          </template>

          <template #[`item.actions`]="{ item }">
            <v-btn
              v-if="item.status === 'pending'"
              size="small"
              variant="tonal"
              color="primary"
              prepend-icon="mdi-check"
              @click="openReview(item)"
            >
              Revisar
            </v-btn>
            <span v-else class="text-caption text-medium-emphasis">Atendido</span>
          </template>

          <template #bottom>
            <PaginatedFooter
              :page="page"
              :per-page="perPage"
              :total="totalReports"
              :in-table="true"
              @update:page="page = $event"
              @update:per-page="onPerPageChange"
            />
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>

    <ReportReviewDialog
      v-model="reviewDialog"
      :report="selectedReport"
      @reviewed="reviewDialog = false"
    />
  </v-container>
</template>
