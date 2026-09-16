<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'

const props = withDefaults(
  defineProps<{
    error: unknown
    retry?: () => void
  }>(),
  {
    retry: undefined,
  },
)

interface ApiErr {
  status?: number
  title?: string
  detail?: string
}

const router = useRouter()
const authStore = useAuthStore()

function getError(): ApiErr {
  if (props.error && typeof props.error === 'object') {
    return props.error as ApiErr
  }
  return { detail: 'Ocurrió un error inesperado' }
}

function getMessage(): string {
  const err = getError()
  return err.detail ?? err.title ?? 'Error desconocido'
}

function getStatus(): number {
  return getError().status ?? 500
}

function goLogin() {
  authStore.clearSession()
  router.push('/login')
}
</script>

<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card max-width="480" class="pa-6 text-center">
      <v-icon size="64" color="error" class="mb-4">mdi-alert-circle</v-icon>

      <template v-if="getStatus() === 401">
        <h2 class="text-h5 mb-2">Sesión expirada</h2>
        <p class="text-body-1 text-medium-emphasis mb-4">
          Tu sesión ha expirado. Por favor inicia sesión nuevamente.
        </p>
        <v-btn color="primary" @click="goLogin"> Ir al login </v-btn>
      </template>

      <template v-else-if="getStatus() === 403">
        <h2 class="text-h5 mb-2">Acceso denegado</h2>
        <p class="text-body-1 text-medium-emphasis mb-4">
          No tienes permisos para acceder a este recurso.
        </p>
        <v-btn color="primary" to="/"> Volver al inicio </v-btn>
      </template>

      <template v-else-if="getStatus() === 404">
        <h2 class="text-h5 mb-2">No encontrado</h2>
        <p class="text-body-1 text-medium-emphasis mb-4">{{ getMessage() }}</p>
        <v-btn color="primary" to="/"> Volver al inicio </v-btn>
      </template>

      <template v-else>
        <h2 class="text-h5 mb-2">Error</h2>
        <p class="text-body-1 text-medium-emphasis mb-4">{{ getMessage() }}</p>
        <v-btn v-if="retry" color="primary" class="mr-2" @click="retry"> Reintentar </v-btn>
        <v-btn variant="text" to="/"> Volver al inicio </v-btn>
      </template>
    </v-card>
  </v-container>
</template>
