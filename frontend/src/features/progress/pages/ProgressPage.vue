<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { upcomingQuestionsOptions, savedQuestionsOptions, toggleSaveMutation } from '@/queries/progress.queries'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import ProgressCard from '../components/ProgressCard.vue'
import type { UpcomingQuestion } from '@/types/progress.types'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const activeTab = ref<'upcoming' | 'saved'>('upcoming')

const { data: upcomingData, isLoading: upcomingLoading, refetch: refetchUpcoming } = useQuery(
  upcomingQuestionsOptions({ page: page.value, perPage: perPage.value }),
)

const { data: savedData, isLoading: savedLoading, refetch: refetchSaved } = useQuery(
  savedQuestionsOptions({ page: page.value, perPage: perPage.value }),
)

const toggleMut = useMutation(toggleSaveMutation())
const togglingId = ref<string | null>(null)

const upcomingList = computed<Array<UpcomingQuestion>>(() => {
  return upcomingData.value?.data ?? []
})

const savedList = computed<Array<UpcomingQuestion>>(() => {
  return savedData.value?.data ?? []
})

const currentList = computed<Array<UpcomingQuestion>>(() => {
  return activeTab.value === 'upcoming' ? upcomingList.value : savedList.value
})

const isLoading = computed(() => {
  return activeTab.value === 'upcoming' ? upcomingLoading.value : savedLoading.value
})

const currentTotal = computed(() => {
  const d = activeTab.value === 'upcoming' ? upcomingData.value : savedData.value
  return d?.meta?.total ?? 0
})

function switchTab(tab: 'upcoming' | 'saved') {
  activeTab.value = tab
  resetPagination()
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['progress', 'upcoming'] })
  queryClient.invalidateQueries({ queryKey: ['progress', 'saved'] })
  if (activeTab.value === 'upcoming') {
    refetchUpcoming()
  } else {
    refetchSaved()
  }
}

async function handleToggle(questionId: string) {
  togglingId.value = questionId
  try {
    await toggleMut.mutateAsync(questionId)
    queryClient.invalidateQueries({ queryKey: ['progress', 'saved'] })
    queryClient.invalidateQueries({ queryKey: ['progress', 'upcoming'] })
    appStore.showSnackbar('Guardado actualizado')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al guardar'
    appStore.showSnackbar(detail, 'error')
  } finally {
    togglingId.value = null
  }
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Progreso"
      create-label=""
      :show-create="false"
      @refresh="handleRefresh"
    />

    <v-tabs v-model="activeTab" color="primary" class="mb-4" @update:model-value="switchTab">
      <v-tab value="upcoming"> Pendientes de repaso </v-tab>
      <v-tab value="saved"> Guardadas </v-tab>
    </v-tabs>

    <v-card variant="flat" border>
      <v-row v-if="isLoading" class="ma-0">
        <v-col v-for="n in 4" :key="n" cols="12" sm="6">
          <v-skeleton-loader type="card" />
        </v-col>
      </v-row>

      <v-row v-else-if="currentList.length" class="ma-0">
        <v-col
          v-for="question in currentList"
          :key="question.id"
          cols="12"
          sm="6"
          lg="4"
        >
          <ProgressCard
            :question="question"
            :show-toggle="true"
            :is-saved="activeTab === 'saved'"
            :toggling="togglingId === question.id"
            @toggle="handleToggle"
          />
        </v-col>
      </v-row>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1" class="mb-2">
          {{ activeTab === 'upcoming' ? 'mdi-calendar-check' : 'mdi-bookmark-outline' }}
        </v-icon>
        <p class="text-body-1 text-medium-emphasis">
          {{ activeTab === 'upcoming' ? 'No hay repasos pendientes' : 'No tienes preguntas guardadas' }}
        </p>
      </v-card-text>

      <template v-if="!isLoading && currentList.length">
        <PaginatedFooter
          :page="page"
          :per-page="perPage"
          :total="currentTotal"
          :in-table="true"
          @update:page="page = $event"
          @update:per-page="perPage = $event"
        />
      </template>
    </v-card>
  </v-container>
</template>
