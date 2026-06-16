<script setup lang="ts">
import LoginForm from '../components/LoginForm.vue'
</script>

<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card max-width="420" class="pa-6">
      <v-card-title class="text-center">Iniciar sesión</v-card-title>
      <LoginForm />
    </v-card>
  </v-container>
</template>
