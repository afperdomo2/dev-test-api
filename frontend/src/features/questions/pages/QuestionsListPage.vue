<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { questionsListOptions } from '@/queries/questions.queries'
import { useAuthStore } from '@/stores/auth.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import QuestionFilters from '../components/QuestionFilters.vue'
import QuestionCard from '../components/QuestionCard.vue'
import type { Question } from '@/types/question.types'
import type { QuestionsFilters } from '@/api/services/questions.service'

const authStore = useAuthStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const filters = ref<QuestionsFilters>({})
const queryFilters = computed(() => filters.value)

const { data, isLoading, refetch } = useQuery(
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

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['questions', 'list'] })
  refetch()
}
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

    <v-row v-if="isLoading">
      <v-col v-for="n in 4" :key="n" cols="12" sm="6" lg="4">
        <v-skeleton-loader type="card" />
      </v-col>
    </v-row>

    <v-row v-else-if="questionList.length">
      <v-col
        v-for="question in questionList"
        :key="question.id"
        cols="12"
        sm="6"
        lg="4"
      >
        <QuestionCard :question="question" />
      </v-col>
    </v-row>

    <v-card v-else>
      <v-card-text class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1" class="mb-2">
          mdi-help-circle-outline
        </v-icon>
        <p class="text-body-1 text-medium-emphasis">No se encontraron preguntas</p>
      </v-card-text>
    </v-card>

    <PaginatedFooter
      :page="page"
      :per-page="perPage"
      :total="totalQuestions"
      class="mt-4"
      @update:page="page = $event"
      @update:per-page="perPage = $event"
    />
  </v-container>
</template>
