<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { questionsListOptions } from '@/queries/questions.queries'
import { usePagination } from '@/composables/usePagination'
import QuestionFilters from '../components/QuestionFilters.vue'
import QuestionCard from '../components/QuestionCard.vue'
import type { Question } from '@/types/question.types'
import type { QuestionsFilters } from '@/api/services/questions.service'

const { page, perPage } = usePagination()

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

function totalQuestions(): number {
  return data.value?.meta?.total ?? 0
}
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Preguntas</h1>

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

    <div
      v-if="totalQuestions() > perPage"
      class="d-flex justify-center mt-4"
    >
      <v-pagination
        v-model="page"
        :length="Math.ceil(totalQuestions() / perPage)"
        :total-visible="5"
      />
    </div>
  </v-container>
</template>
