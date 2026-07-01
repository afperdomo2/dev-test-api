<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { questionsListOptions, deleteQuestionMutation } from '@/queries/questions.queries'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import QuestionFilters from '../components/QuestionFilters.vue'
import QuestionTable from '../components/QuestionTable.vue'
import QuestionFormDialog from '../components/QuestionFormDialog.vue'
import type { Question } from '@/types/question.types'
import type { QuestionsFilters } from '@/api/services/questions.service'

const authStore = useAuthStore()
const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const filters = ref<QuestionsFilters>({})
const queryFilters = computed(() => filters.value)

const { data, isLoading } = useQuery(
  questionsListOptions(
    () => page.value,
    () => perPage.value,
    () => queryFilters.value,
  ),
)

const deleteMut = useMutation(deleteQuestionMutation())
const deleteTarget = ref<Question | null>(null)
const deleteDialog = ref(false)

const formDialogOpen = ref(false)
const editingQuestion = ref<Question | null>(null)

const questionList = computed<Array<Question>>(() => {
  return data.value?.data ?? []
})

const totalQuestions = computed(() => data.value?.meta?.total ?? 0)

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['questions', 'list'] })
}

function openCreate() {
  editingQuestion.value = null
  formDialogOpen.value = true
}

function openEdit(question: Question) {
  editingQuestion.value = question
  formDialogOpen.value = true
}

function confirmDelete(question: Question) {
  deleteTarget.value = question
  deleteDialog.value = true
}

async function executeDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteMut.mutateAsync(deleteTarget.value.id)
    appStore.showSnackbar('Pregunta eliminada')
    deleteDialog.value = false
    deleteTarget.value = null
    queryClient.invalidateQueries({ queryKey: ['questions', 'list'] })
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al eliminar pregunta'
    appStore.showSnackbar(detail, 'error')
  }
}

watch(filters, () => {
  resetPagination()
})
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Preguntas"
      create-label="Nueva pregunta"
      :show-create="true"
      @refresh="handleRefresh"
      @create="openCreate"
    />

    <QuestionFilters @change="filters = $event" />

    <v-card>
      <v-card-text>
        <QuestionTable
          :questions="questionList"
          :loading="isLoading"
          :items-per-page="perPage"
          :current-user-id="authStore.user?.id"
          @edit="openEdit"
          @delete="confirmDelete"
        >
          <template #footer>
            <PaginatedFooter
              :page="page"
              :per-page="perPage"
              :total="totalQuestions"
              :in-table="true"
              @update:page="page = $event"
              @update:per-page="onPerPageChange"
            />
          </template>
        </QuestionTable>
      </v-card-text>
    </v-card>

    <QuestionFormDialog
      v-model="formDialogOpen"
      :question="editingQuestion"
      @saved="formDialogOpen = false"
    />

    <v-dialog v-model="deleteDialog" max-width="420">
      <v-card>
        <v-card-title>Eliminar pregunta</v-card-title>
        <v-card-text>
          ¿Estás seguro de eliminar esta pregunta?
          <template v-if="deleteTarget">
            <br /><br />
            <strong>{{ deleteTarget.content }}</strong>
          </template>
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
