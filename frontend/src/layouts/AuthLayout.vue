<script setup lang="ts">
import { useAppStore } from '@/stores/app.store'

const appStore = useAppStore()
</script>

<template>
  <v-app :theme="appStore.theme">
    <v-main>
      <v-container class="fill-height pa-0" fluid>
        <router-view />
      </v-container>
    </v-main>

    <v-snackbar
      v-model="appStore.snackbar.show"
      :color="appStore.snackbar.color"
      :timeout="4000"
      location="bottom right"
    >
      {{ appStore.snackbar.message }}
      <template #actions>
        <v-btn variant="text" icon="mdi-close" @click="appStore.hideSnackbar" />
      </template>
    </v-snackbar>
  </v-app>
</template>
