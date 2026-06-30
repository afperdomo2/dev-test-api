<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { questionsListOptions } from '@/queries/questions.queries'
import { useAuthStore } from '@/stores/auth.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import QuestionFilters from '../components/QuestionFilters.vue'
import QuestionTable from '../components/QuestionTable.vue'
import type { Question } from '@/types/question.types'
import type { QuestionsFilters } from '@/api/services/questions.service'

const authStore = useAuthStore()
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

watch(filters, () => {
  resetPagination()
})
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Preguntas"
      create-label="Nueva pregunta"
      :show-create="authStore.isAdmin"
      @refresh="handleRefresh"
    />

    <QuestionFilters @change="filters = $event" />

    <v-card>
      <v-card-text>
        <QuestionTable :questions="questionList" :loading="isLoading" :items-per-page="perPage">
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
  </v-container>
</template>
