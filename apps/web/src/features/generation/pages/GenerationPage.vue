<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import ListPageHeader from '@/components/ListPageHeader.vue'
import GenerationPanel from '../components/GenerationPanel.vue'
import type { GenerationKind } from '@/types/generation.types'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()

const activeTab = ref<GenerationKind>(route.query.type === 'tip' ? 'tip' : 'question')

const tabItems: Array<{ title: string; value: GenerationKind; icon: string }> = [
  { title: 'Preguntas', value: 'question', icon: 'mdi-help-circle-outline' },
  { title: 'Tips', value: 'tip', icon: 'mdi-lightbulb-on-outline' },
]

const panelKey = computed(() => activeTab.value)

watch(activeTab, (value) => {
  const current = route.query.type
  if (current === value) return
  router.replace({ query: { ...route.query, type: value } })
})

// Mantiene el tab sincronizado si el usuario navega con atrás/adelante.
watch(
  () => route.query.type,
  (value) => {
    const next: GenerationKind = value === 'tip' ? 'tip' : 'question'
    if (next !== activeTab.value) activeTab.value = next
  },
)

function handleRefresh() {
  queryClient.invalidateQueries({ queryKey: ['generation'] })
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Generación con IA"
      create-label=""
      :show-create="false"
      @refresh="handleRefresh"
    />

    <v-tabs v-model="activeTab" color="primary" class="mb-4">
      <v-tab v-for="tab in tabItems" :key="tab.value" :value="tab.value">
        <v-icon start :icon="tab.icon" />
        {{ tab.title }}
      </v-tab>
    </v-tabs>

    <GenerationPanel :key="panelKey" :kind="activeTab" />
  </v-container>
</template>
