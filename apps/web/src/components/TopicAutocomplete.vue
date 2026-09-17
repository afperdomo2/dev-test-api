<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useInfiniteQuery } from '@tanstack/vue-query'
import { topicsInfiniteOptions } from '@/queries/topics.queries'
import { useDebounce } from '@/composables/useDebounce'
import type { Topic } from '@/types/topic.types'

interface Props {
  modelValue: Array<string>
  active?: boolean
  disabled?: boolean
  label?: string
  placeholder?: string
  hint?: string
  persistentHint?: boolean
  hideDetails?: boolean
  density?: 'default' | 'comfortable' | 'compact'
  errorMessages?: string | Array<string>
  required?: boolean
  noDataText?: string
  itemTitle?: (topic: Topic) => string
  itemSubtitle?: (topic: Topic) => string
}

const props = withDefaults(defineProps<Props>(), {
  active: true,
  disabled: false,
  label: 'Temas',
  placeholder: undefined,
  hint: undefined,
  persistentHint: false,
  hideDetails: false,
  density: 'default',
  errorMessages: undefined,
  required: false,
  noDataText: 'No se encontraron temas',
  itemTitle: (topic: Topic) => topic.name,
  itemSubtitle: (topic: Topic) => topic.category,
})

const model = defineModel<Array<string>>({ required: true })
const selectedTopics = defineModel<Array<Topic>>('selectedTopics', { default: () => [] })

const { value: search, debouncedValue: debouncedSearch } = useDebounce('', 400)

const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useInfiniteQuery(
  topicsInfiniteOptions(
    () => debouncedSearch.value ?? '',
    () => props.active,
  ),
)

const topics = computed<Array<Topic>>(() => data.value?.pages.flatMap((p) => p.data) ?? [])

const items = computed(() =>
  topics.value.map((t) => ({
    title: props.itemTitle(t),
    value: t.id,
    props: { subtitle: props.itemSubtitle(t) },
  })),
)

const topicMap = ref<Map<string, Topic>>(new Map())

watch(
  topics,
  (list) => {
    const next = new Map(topicMap.value)
    for (const t of list) next.set(t.id, t)
    topicMap.value = next
    syncSelected()
  },
  { immediate: true },
)

function syncSelected() {
  selectedTopics.value = model.value
    .map((id) => topicMap.value.get(id))
    .filter((t): t is Topic => Boolean(t))
}

function onModelValueUpdate(value: Array<string>) {
  model.value = value
  syncSelected()
}

watch(() => model.value, syncSelected, { deep: true })

function onIntersect(isIntersecting: boolean) {
  if (isIntersecting && hasNextPage.value && !isFetchingNextPage.value) {
    fetchNextPage()
  }
}
</script>

<template>
  <v-autocomplete
    v-model="model"
    v-model:search="search"
    :items="items"
    item-title="title"
    item-props="props"
    :label="label"
    :placeholder="placeholder"
    :hint="hint"
    :persistent-hint="persistentHint"
    :hide-details="hideDetails"
    :density="density"
    :error-messages="errorMessages"
    :disabled="disabled"
    :required="required"
    :loading="isLoading"
    multiple
    chips
    closable-chips
    clearable
    no-filter
    @update:model-value="onModelValueUpdate"
  >
    <template #append-item>
      <div v-intersect="onIntersect" class="pa-2 text-center">
        <v-progress-circular v-if="isFetchingNextPage" indeterminate size="20" />
        <span v-else-if="!hasNextPage && topics.length" class="text-caption text-medium-emphasis">
          No hay más temas
        </span>
      </div>
    </template>
    <template #no-data>
      <div v-if="!isLoading" class="pa-4 text-center">
        <p class="text-body-2 text-medium-emphasis mb-1">{{ noDataText }}</p>
        <p class="text-caption text-medium-emphasis">Prueba con otro término de búsqueda</p>
      </div>
    </template>
  </v-autocomplete>
</template>
