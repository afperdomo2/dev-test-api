<script setup lang="ts">
import { computed } from 'vue'
import type { TopicXpItem } from '@/types/gamification.types'

const props = defineProps<{
  items: Array<TopicXpItem>
}>()

const sorted = computed<Array<TopicXpItem>>(() => [...props.items].sort((a, b) => b.xp - a.xp))

const maxXp = computed(() => sorted.value.reduce((max, item) => Math.max(max, item.xp), 0))

function barPercent(xp: number): number {
  if (maxXp.value <= 0) return 0
  return Math.max(4, Math.round((xp / maxXp.value) * 100))
}
</script>

<template>
  <div>
    <div class="d-flex align-center ga-2 mb-3">
      <v-icon icon="mdi-tag-multiple" size="20" color="primary" />
      <span class="text-subtitle-2">XP por tema</span>
    </div>

    <v-list v-if="sorted.length" density="compact" lines="one" class="py-0">
      <v-list-item v-for="item in sorted" :key="item.topicId" class="px-0">
        <v-list-item-title class="text-body-2 mb-1">
          <div class="d-flex align-center justify-space-between ga-2">
            <span class="text-truncate">{{ item.name }}</span>
            <span class="text-caption font-weight-medium">{{ item.xp }} XP</span>
          </div>
        </v-list-item-title>
        <v-progress-linear :model-value="barPercent(item.xp)" color="primary" height="8" rounded />
      </v-list-item>
    </v-list>

    <div v-else class="text-body-2 text-medium-emphasis text-center py-4">
      Aún no hay XP registrado por tema.
    </div>
  </div>
</template>
