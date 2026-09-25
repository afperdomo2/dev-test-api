<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { topicsListOptions } from '@/queries/topics.queries'
import { useDebounce } from '@/composables/useDebounce'
import type { Topic } from '@/types/topic.types'

interface Props {
  label?: string
  placeholder?: string
  hint?: string
  persistentHint?: boolean
  hideDetails?: boolean
  density?: 'default' | 'comfortable' | 'compact'
  errorMessages?: string | Array<string>
  required?: boolean
  disabled?: boolean
  clearable?: boolean
  itemTitle?: (topic: Topic) => string
}

const props = withDefaults(defineProps<Props>(), {
  label: 'Tema',
  placeholder: undefined,
  hint: undefined,
  persistentHint: false,
  hideDetails: false,
  density: 'default',
  errorMessages: undefined,
  required: false,
  disabled: false,
  clearable: true,
  itemTitle: (topic: Topic) => `${topic.name} (${topic.slug})`,
})

const model = defineModel<string | null>({ required: true })

const { value: search, debouncedValue: debouncedSearch } = useDebounce('', 300)

const { data, isLoading } = useQuery(
  topicsListOptions(
    () => 1,
    () => 100,
    () => 'name',
    () => 'asc',
    () => debouncedSearch.value ?? '',
  ),
)

const items = computed<Array<{ title: string; value: string }>>(() =>
  (data.value?.data ?? []).map((topic) => ({
    title: props.itemTitle(topic),
    value: topic.id,
  })),
)
</script>

<template>
  <v-autocomplete
    v-model="model"
    v-model:search="search"
    :items="items"
    item-title="title"
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
    :clearable="clearable"
    no-filter
  />
</template>
