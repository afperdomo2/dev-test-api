<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { unreadCountOptions } from '@/queries/notifications.queries'

const { data } = useQuery(unreadCountOptions())
const count = computed(() => data.value?.count ?? 0)
</script>

<template>
  <v-btn icon variant="text" to="/notifications" aria-label="Notificaciones">
    <v-badge :content="count" :model-value="count > 0" :max="99" color="error">
      <v-icon icon="mdi-bell-outline" />
    </v-badge>
  </v-btn>
</template>
