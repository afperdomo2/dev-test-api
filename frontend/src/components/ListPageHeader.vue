<script setup lang="ts">
import { ref } from 'vue'
import { REFRESH_COOLDOWN_MS } from '@/constants'

interface Props {
  title: string
  createLabel: string
  createTo?: string
  showCreate?: boolean
}

withDefaults(defineProps<Props>(), {
  createTo: undefined,
  showCreate: true,
})

const emit = defineEmits<{
  refresh: []
  create: []
}>()

const refreshCooldown = ref(false)

function handleRefresh() {
  if (refreshCooldown.value) return
  refreshCooldown.value = true
  emit('refresh')
  setTimeout(() => {
    refreshCooldown.value = false
  }, REFRESH_COOLDOWN_MS)
}
</script>

<template>
  <div class="d-flex align-center justify-space-between mb-4">
    <h1 class="text-h4">{{ title }}</h1>
    <div class="d-flex ga-2">
      <v-btn
        variant="tonal"
        color="secondary"
        icon="mdi-refresh"
        :disabled="refreshCooldown"
        @click="handleRefresh"
      />
      <v-btn
        v-if="showCreate"
        color="primary"
        prepend-icon="mdi-plus"
        :to="createTo"
        @click="!createTo && emit('create')"
      >
        {{ createLabel }}
      </v-btn>
    </div>
  </div>
</template>
