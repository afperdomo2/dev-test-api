<script setup lang="ts">
import { computed } from 'vue'
import type { User } from '@/types/user.types'
import { formatDate } from '@/utils/format'

const props = defineProps<{
  user: User | null
}>()

const formattedDate = computed(() =>
  props.user?.createdAt ? formatDate(props.user.createdAt) : '',
)
</script>

<template>
  <v-card>
    <v-card-item>
      <template #prepend>
        <v-icon color="primary" size="40">mdi-account-circle</v-icon>
      </template>
      <v-card-title class="text-h5">
        {{ user?.email }}
      </v-card-title>
      <v-card-subtitle class="mt-1"> Miembro desde {{ formattedDate }} </v-card-subtitle>
    </v-card-item>
    <v-card-text>
      <v-chip
        :color="user?.isAdmin ? 'error' : 'info'"
        size="small"
        variant="tonal"
        class="font-weight-bold"
      >
        <v-icon start size="18">
          {{ user?.isAdmin ? 'mdi-shield-account' : 'mdi-account' }}
        </v-icon>
        {{ user?.isAdmin ? 'Administrador' : 'Usuario' }}
      </v-chip>
      <div class="d-flex align-center ga-2 mt-3 text-caption text-medium-emphasis">
        <v-icon size="16" color="success">mdi-identifier</v-icon>
        <span class="text-truncate">{{ user?.id }}</span>
      </div>
    </v-card-text>
  </v-card>
</template>
