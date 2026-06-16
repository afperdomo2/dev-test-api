<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ITEMS_PER_PAGE_OPTIONS } from '@/constants'

interface Props {
  page: number
  perPage: number
  total: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:page': [page: number]
  'update:perPage': [perPage: number]
}>()

const localPerPage = ref(props.perPage)

watch(
  () => props.perPage,
  (val) => {
    localPerPage.value = val
  },
)

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.perPage)))

const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.perPage + 1))
const to = computed(() => Math.min(props.page * props.perPage, props.total))
</script>

<template>
  <div class="d-flex align-center justify-space-between flex-wrap ga-2">
    <div class="d-flex align-center ga-2">
      <v-select
        :model-value="perPage"
        :items="ITEMS_PER_PAGE_OPTIONS as unknown as Array<number>"
        density="compact"
        variant="outlined"
        hide-details
        style="width: 80px"
        @update:model-value="emit('update:perPage', $event as number)"
      />
      <span class="text-body-2 text-medium-emphasis">
        registros por página
      </span>
    </div>

    <span class="text-body-2 text-medium-emphasis">
      Mostrando {{ from }}–{{ to }} de {{ total }}
    </span>

    <v-pagination
      v-if="totalPages > 1"
      :model-value="page"
      :length="totalPages"
      :total-visible="5"
      show-first-last-page
      density="compact"
      variant="outlined"
      @update:model-value="emit('update:page', $event)"
    />
  </div>
</template>
