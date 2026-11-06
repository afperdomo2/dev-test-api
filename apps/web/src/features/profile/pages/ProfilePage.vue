<script setup lang="ts">
import { useAuthStore } from '@/stores/auth.store'
import UserSummaryCard from '@/components/UserSummaryCard.vue'
import ChangePasswordCard from '../components/ChangePasswordCard.vue'

const authStore = useAuthStore()
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Mi perfil</h1>

    <v-row>
      <v-col cols="12" md="6">
        <UserSummaryCard :user="authStore.user" />
      </v-col>

      <v-col cols="12" md="6">
        <ChangePasswordCard />
      </v-col>
    </v-row>
  </v-container>
</template>
