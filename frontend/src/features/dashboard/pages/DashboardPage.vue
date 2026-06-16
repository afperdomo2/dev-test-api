<script setup lang="ts">
import { useAuthStore } from '@/stores/auth.store'
import { computed } from 'vue'

const authStore = useAuthStore()

const formattedDate = computed(() => {
  const date = authStore.user?.createdAt
  if (!date) return ''
  return new Date(date).toLocaleDateString('es-CO', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Dashboard</h1>

    <v-row>
      <v-col cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="primary" size="40">mdi-account-circle</v-icon>
            </template>
            <v-card-title class="text-h5">
              {{ authStore.user?.email }}
            </v-card-title>
            <v-card-subtitle class="mt-1"> Miembro desde {{ formattedDate }} </v-card-subtitle>
          </v-card-item>
          <v-card-text>
            <v-chip
              :color="authStore.isAdmin ? 'error' : 'info'"
              size="small"
              variant="tonal"
              class="font-weight-bold"
            >
              <v-icon start size="18">
                {{ authStore.isAdmin ? 'mdi-shield-account' : 'mdi-account' }}
              </v-icon>
              {{ authStore.isAdmin ? 'Administrador' : 'Usuario' }}
            </v-chip>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="warning" size="40">mdi-calendar-clock</v-icon>
            </template>
            <v-card-title>Pendientes de repaso</v-card-title>
            <v-card-subtitle>Carga desde el backend</v-card-subtitle>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="success" size="40">mdi-information-outline</v-icon>
            </template>
            <v-card-title>ID de usuario</v-card-title>
            <v-card-subtitle class="text-caption text-truncate">
              {{ authStore.user?.id }}
            </v-card-subtitle>
          </v-card-item>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
