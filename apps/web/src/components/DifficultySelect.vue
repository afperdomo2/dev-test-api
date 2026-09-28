<script setup lang="ts">
import { SESSION_DIFFICULTIES } from '@/types/session.types'

interface Props {
  modelValue: string
  label?: string
  disabled?: boolean
  required?: boolean
}

withDefaults(defineProps<Props>(), {
  label: 'Dificultad',
  disabled: false,
  required: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<template>
  <v-select
    :model-value="modelValue"
    :items="SESSION_DIFFICULTIES"
    item-props="props"
    :label="label"
    :disabled="disabled"
    density="compact"
    variant="outlined"
    hide-details
    @update:model-value="emit('update:modelValue', $event)"
  />
</template>
