<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { topicsListOptions, deleteTopicMutation } from '@/queries/topics.queries'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import { useDebounce } from '@/composables/useDebounce'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import TopicFormDialog from '../components/TopicFormDialog.vue'
import type { Topic } from '@/types/topic.types'

const authStore = useAuthStore()
const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const { value: searchText, debouncedValue: debouncedSearch } = useDebounce('', 500)
const myOnly = ref(false)

const { data, isLoading, refetch } = useQuery(
  topicsListOptions(
    () => page.value,
    () => perPage.value,
    () => 'name',
    () => 'asc',
    () => debouncedSearch.value || '',
    () => myOnly.value,
  ),
)

const deleteMut = useMutation(deleteTopicMutation())
const deleteTarget = ref<Topic | null>(null)
const deleteDialog = ref(false)

const dialogOpen = ref(false)
const editingTopic = ref<Topic | null>(null)

const topicList = computed<Array<Topic>>(() => {
  return data.value?.data ?? []
})

const totalTopics = computed(() => data.value?.meta?.total ?? 0)

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['topics', 'list'] })
  refetch()
}

watch(debouncedSearch, () => {
  resetPagination()
})

function onMyOnlyChange(val: boolean | null) {
  if (val === null) return
  myOnly.value = val
  resetPagination()
}

const headers = [
  { title: 'Nombre', key: 'name', sortable: false },
  { title: 'Slug', key: 'slug', sortable: false },
  { title: 'Categoría', key: 'category', sortable: false, align: 'center' as const },
  { title: 'Tipo', key: 'isSystem', sortable: false, align: 'center' as const },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'center' as const },
]

function openCreate() {
  editingTopic.value = null
  dialogOpen.value = true
}

function openEdit(topic: Topic) {
  editingTopic.value = topic
  dialogOpen.value = true
}

function confirmDelete(topic: Topic) {
  deleteTarget.value = topic
  deleteDialog.value = true
}

async function executeDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteMut.mutateAsync(deleteTarget.value.id)
    appStore.showSnackbar('Tema eliminado')
    deleteDialog.value = false
    deleteTarget.value = null
    queryClient.invalidateQueries({ queryKey: ['topics', 'list'] })
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al eliminar tema'
    appStore.showSnackbar(detail, 'error')
  }
}

function canModify(topic: Topic): boolean {
  if (authStore.isAdmin) return topic.isSystem
  return !topic.isSystem
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Temas"
      create-label="Nuevo tema"
      :show-create="true"
      @refresh="handleRefresh"
      @create="openCreate"
    />

    <v-row class="mb-4" dense>
      <v-col cols="12" md="5" lg="4">
        <v-text-field
          v-model="searchText"
          prepend-inner-icon="mdi-magnify"
          label="Buscar"
          clearable
          hide-details
          density="compact"
          @click:clear="searchText = ''"
        />
      </v-col>
      <v-col v-if="!authStore.isAdmin" cols="auto">
        <v-switch
          :model-value="myOnly"
          label="Mis temas"
          color="primary"
          hide-details
          density="compact"
          inset
          @update:model-value="onMyOnlyChange"
        />
      </v-col>
    </v-row>

    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="topicList"
          :items-per-page="perPage"
          :loading="isLoading"
          no-data-text="No hay temas"
          loading-text="Cargando temas..."
          hover
        >
          <template #[`item.isSystem`]="{ item }">
            <v-chip :color="item.isSystem ? 'info' : 'success'" size="small" variant="tonal">
              {{ item.isSystem ? 'Sistema' : 'Custom' }}
            </v-chip>
          </template>

          <template #[`item.actions`]="{ item }">
            <div class="d-flex ga-1 justify-center">
              <template v-if="canModify(item)">
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
              </template>
              <v-tooltip
                v-else-if="item.isSystem"
                text="Los temas del sistema no se pueden modificar"
              >
                <template #activator="{ props }">
                  <span v-bind="props">
                    <v-btn
                      icon="mdi-lock-outline"
                      variant="text"
                      size="small"
                      color="disabled"
                      disabled
                    />
                  </span>
                </template>
              </v-tooltip>
            </div>
          </template>

          <template #bottom>
            <PaginatedFooter
              :page="page"
              :per-page="perPage"
              :total="totalTopics"
              :in-table="true"
              @update:page="page = $event"
              @update:per-page="onPerPageChange"
            />
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>

    <TopicFormDialog v-model="dialogOpen" :topic="editingTopic" @saved="dialogOpen = false" />

    <v-dialog v-model="deleteDialog" max-width="420">
      <v-card>
        <v-card-title>Eliminar tema</v-card-title>
        <v-card-text>
          ¿Estás seguro de eliminar
          <strong>{{ deleteTarget?.name }}</strong
          >?
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
