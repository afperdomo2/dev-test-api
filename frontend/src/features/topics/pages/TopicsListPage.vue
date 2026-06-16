<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { topicsListOptions, deleteTopicMutation } from '@/queries/topics.queries'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import TopicFormDialog from '../components/TopicFormDialog.vue'
import type { Topic } from '@/types/topic.types'

const authStore = useAuthStore()
const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage } = usePagination()

const { data, isLoading } = useQuery(
  topicsListOptions(() => page.value, () => perPage.value),
)

const deleteMut = useMutation(deleteTopicMutation())
const deleteTarget = ref<Topic | null>(null)
const deleteDialog = ref(false)

const dialogOpen = ref(false)
const editingTopic = ref<Topic | null>(null)

const topicList = computed<Array<Topic>>(() => {
  return data.value?.data ?? []
})

function totalTopics(): number {
  return data.value?.meta?.total ?? 0
}

const headers = [
  { title: 'Nombre', key: 'name', sortable: false },
  { title: 'Slug', key: 'slug', sortable: false },
  { title: 'Categoría', key: 'category', sortable: false },
  { title: 'Tipo', key: 'is_system', sortable: false },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'end' as const },
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
  return authStore.isAdmin && !topic.is_system
}
</script>

<template>
  <v-container>
    <div class="d-flex align-center justify-space-between mb-4">
      <h1 class="text-h4">Temas</h1>
      <v-btn
        v-if="authStore.isAdmin"
        color="primary"
        prepend-icon="mdi-plus"
        @click="openCreate"
      >
        Nuevo tema
      </v-btn>
    </div>

    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="topicList"
          :loading="isLoading"
          items-per-page-text="Temas por página"
          no-data-text="No hay temas"
          loading-text="Cargando temas..."
          hover
        >
          <template #[`item.is_system`]="{ item }">
            <v-chip
              :color="item.is_system ? 'info' : 'success'"
              size="small"
              variant="tonal"
            >
              {{ item.is_system ? 'Sistema' : 'Custom' }}
            </v-chip>
          </template>

          <template #[`item.actions`]="{ item }">
            <div class="d-flex ga-1 justify-end">
              <v-btn
                v-if="canModify(item)"
                icon="mdi-pencil"
                variant="text"
                size="small"
                color="primary"
                @click="openEdit(item)"
              />
              <v-btn
                v-if="canModify(item)"
                icon="mdi-delete"
                variant="text"
                size="small"
                color="error"
                @click="confirmDelete(item)"
              />
            </div>
          </template>
        </v-data-table>

        <div
          v-if="totalTopics() > perPage"
          class="d-flex justify-center mt-4"
        >
          <v-pagination
            v-model="page"
            :length="Math.ceil(totalTopics() / perPage)"
            :total-visible="5"
          />
        </div>
      </v-card-text>
    </v-card>

    <TopicFormDialog
      v-model="dialogOpen"
      :topic="editingTopic"
      @saved="dialogOpen = false"
    />

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
          <v-btn
            color="error"
            :loading="deleteMut.isPending.value"
            @click="executeDelete"
          >
            Eliminar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
