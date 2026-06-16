<script setup lang="ts">
import LoginForm from '../components/LoginForm.vue'
</script>

<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card max-width="480" class="pa-sm-6 pa-4 w-100">
      <v-card-title class="text-center text-h5 mb-2">Iniciar sesión</v-card-title>
      <LoginForm />
    </v-card>
  </v-container>
</template>
