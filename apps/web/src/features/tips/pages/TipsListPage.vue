<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { tipsListOptions, deleteTipMutation } from '@/queries/tips.queries'
import { useAppStore } from '@/stores/app.store'
import { useAuthStore } from '@/stores/auth.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import TipFormDialog from '../components/TipFormDialog.vue'
import TopicSelect from '@/components/TopicSelect.vue'
import type { Tip } from '@/types/tip.types'

const appStore = useAppStore()
const authStore = useAuthStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const topicFilter = ref<string | null>(null)

const { data, isLoading } = useQuery(
  tipsListOptions(
    () => page.value,
    () => perPage.value,
    () => ({ topicId: topicFilter.value ?? undefined }),
  ),
)

const deleteMut = useMutation(deleteTipMutation())

const dialogOpen = ref(false)
const editingTip = ref<Tip | null>(null)
const deleteTarget = ref<Tip | null>(null)
const deleteDialog = ref(false)

const tipList = computed<Array<Tip>>(() => data.value?.data ?? [])
const totalTips = computed(() => data.value?.meta?.total ?? 0)

const headers = computed(() => {
  const cols = [
    { title: 'Contenido', key: 'content', sortable: false },
    { title: 'Tema', key: 'topic', sortable: false },
    { title: 'Origen', key: 'source', sortable: false, align: 'center' as const },
  ]
  if (authStore.isAdmin) {
    cols.push({ title: 'Acciones', key: 'actions', sortable: false, align: 'center' as const })
  }
  return cols
})

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['tips', 'list'] })
}

watch(topicFilter, () => {
  resetPagination()
})

function openCreate() {
  editingTip.value = null
  dialogOpen.value = true
}

function openEdit(tip: Tip) {
  editingTip.value = tip
  dialogOpen.value = true
}

function confirmDelete(tip: Tip) {
  deleteTarget.value = tip
  deleteDialog.value = true
}

async function executeDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteMut.mutateAsync(deleteTarget.value.id)
    appStore.showSnackbar('Tip eliminado')
    deleteDialog.value = false
    deleteTarget.value = null
    queryClient.invalidateQueries({ queryKey: ['tips'] })
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al eliminar tip'
    appStore.showSnackbar(detail, 'error')
  }
}

function preview(tip: Tip): string {
  return tip.content.length > 90 ? `${tip.content.slice(0, 90)}…` : tip.content
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Tips"
      create-label="Nuevo tip"
      :show-create="authStore.isAdmin"
      @refresh="handleRefresh"
      @create="openCreate"
    />

    <v-alert type="info" variant="tonal" density="compact" class="mb-4">
      Los tips son públicos para todos los usuarios y se muestran como consejo del día durante la
      sesión diaria. Solo un administrador puede crearlos o editarlos.
    </v-alert>

    <v-row class="mb-4" dense>
      <v-col cols="12" md="5" lg="4">
        <TopicSelect
          v-model="topicFilter"
          label="Filtrar por tema"
          density="compact"
          variant="outlined"
          hide-details
        />
      </v-col>
      <v-col v-if="topicFilter" cols="auto" class="d-flex align-center">
        <v-btn size="small" variant="text" color="error" @click="topicFilter = null">
          Limpiar
        </v-btn>
      </v-col>
    </v-row>

    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="tipList"
          :items-per-page="perPage"
          :loading="isLoading"
          no-data-text="No hay tips"
          loading-text="Cargando tips..."
          hover
        >
          <template #[`item.content`]="{ item }">
            <span class="text-body-2">{{ preview(item) }}</span>
          </template>

          <template #[`item.source`]="{ item }">
            <v-chip
              :color="item.source === 'imported' ? 'info' : 'success'"
              size="small"
              variant="tonal"
            >
              {{ item.source === 'imported' ? 'Importado' : 'Manual' }}
            </v-chip>
          </template>

          <template #[`item.actions`]="{ item }">
            <div class="d-flex ga-1 justify-center">
              <v-btn
                icon="mdi-pencil"
                variant="text"
                size="small"
                color="primary"
                @click="openEdit(item)"
              />
              <v-btn
                icon="mdi-delete"
                variant="text"
                size="small"
                color="error"
                @click="confirmDelete(item)"
              />
            </div>
          </template>

          <template #bottom>
            <PaginatedFooter
              :page="page"
              :per-page="perPage"
              :total="totalTips"
              :in-table="true"
              @update:page="page = $event"
              @update:per-page="onPerPageChange"
            />
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>

    <TipFormDialog v-model="dialogOpen" :tip="editingTip" @saved="dialogOpen = false" />

    <v-dialog v-model="deleteDialog" max-width="420">
      <v-card>
        <v-card-title>Eliminar tip</v-card-title>
        <v-card-text>
          ¿Estás seguro de eliminar este tip?
          <div class="text-caption text-medium-emphasis mt-2">
            {{ deleteTarget ? preview(deleteTarget) : '' }}
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false"> Cancelar </v-btn>
          <v-btn color="error" :loading="deleteMut.isPending.value" @click="executeDelete">
            Eliminar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
