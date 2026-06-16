<script setup lang="ts">
import { computed } from 'vue'
import { ITEMS_PER_PAGE_OPTIONS } from '@/constants'

interface Props {
  page: number
  perPage: number
  total: number
  inTable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  inTable: false,
})

const emit = defineEmits<{
  'update:page': [page: number]
  'update:perPage': [perPage: number]
}>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.perPage)))

const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.perPage + 1))
const to = computed(() => Math.min(props.page * props.perPage, props.total))
</script>

<template>
  <template v-if="inTable">
    <v-divider />
    <div class="d-flex align-center justify-space-between flex-wrap ga-2 pa-3">
      <div class="d-flex align-center ga-2" style="white-space: nowrap;">
        <v-select
          :model-value="perPage"
          :items="ITEMS_PER_PAGE_OPTIONS as unknown as Array<number>"
          density="compact"
          variant="outlined"
          hide-details
          style="min-width: 100px;"
          @update:model-value="emit('update:perPage', $event as number)"
        />
        <span class="text-body-2 text-medium-emphasis" style="white-space: nowrap;">
          registros por página
        </span>
      </div>

      <span class="text-body-2 text-medium-emphasis" style="white-space: nowrap;">
        Mostrando {{ from }}–{{ to }} de {{ total }}
      </span>

      <v-pagination
        :model-value="page"
        :length="totalPages"
        :total-visible="5"
        show-first-last-page
        density="compact"
        :disabled="totalPages <= 1"
        @update:model-value="emit('update:page', $event)"
      />
    </div>
  </template>

  <v-card v-else variant="flat" border class="mt-4">
    <div class="d-flex align-center justify-space-between flex-wrap ga-2 pa-3">
      <div class="d-flex align-center ga-2" style="white-space: nowrap;">
        <v-select
          :model-value="perPage"
          :items="ITEMS_PER_PAGE_OPTIONS as unknown as Array<number>"
          density="compact"
          variant="outlined"
          hide-details
          style="min-width: 100px;"
          @update:model-value="emit('update:perPage', $event as number)"
        />
        <span class="text-body-2 text-medium-emphasis" style="white-space: nowrap;">
          registros por página
        </span>
      </div>

      <span class="text-body-2 text-medium-emphasis" style="white-space: nowrap;">
        Mostrando {{ from }}–{{ to }} de {{ total }}
      </span>

      <v-pagination
        :model-value="page"
        :length="totalPages"
        :total-visible="5"
        show-first-last-page
        density="compact"
        :disabled="totalPages <= 1"
        @update:model-value="emit('update:page', $event)"
      />
    </div>
  </v-card>
</template>
