<script setup lang="ts">
import { useAuthStore } from '@/stores/auth.store'

const authStore = useAuthStore()
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Dashboard</h1>

    <v-row>
      <v-col cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="primary" size="36">mdi-account</v-icon>
            </template>
            <v-card-title>{{ authStore.user?.email }}</v-card-title>
            <v-card-subtitle>
              {{ authStore.isAdmin ? 'Administrador' : 'Usuario' }}
            </v-card-subtitle>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" lg="4">
        <v-card>
          <v-card-item>
            <template #prepend>
              <v-icon color="warning" size="36">mdi-calendar-clock</v-icon>
            </template>
            <v-card-title>Pendientes de repaso</v-card-title>
            <v-card-subtitle>Carga desde el backend</v-card-subtitle>
          </v-card-item>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
