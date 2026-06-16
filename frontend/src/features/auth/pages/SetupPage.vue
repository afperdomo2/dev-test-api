<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { setup as setupApi } from '@/api/services/auth.service'
import type { SetupRequest } from '@/types/auth.types'
import { emailRule, passwordRule, requiredRule, validateRules } from '@/utils/validators'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const form = ref<SetupRequest>({ email: '', password: '' })
const errors = ref<Record<string, Array<string>>>({})
const loading = ref(false)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.email = validateRules([requiredRule(), emailRule()], form.value.email)
  newErrors.password = validateRules(
    [requiredRule(), passwordRule(8, 'Mínimo 8 caracteres')],
    form.value.password,
  )
  errors.value = newErrors
  return Object.values(newErrors).every((e) => e.length === 0)
}

async function submit() {
  if (!validate()) return
  loading.value = true
  try {
    const res = await setupApi(form.value)
    authStore.setSession(res.token, res.user)
    router.push('/')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error en la configuración inicial'
    appStore.showSnackbar(detail, 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card max-width="480" class="pa-sm-6 pa-4 w-100">
      <v-card-title class="text-center text-h5 mb-2">Configuración inicial</v-card-title>
      <v-card-subtitle class="text-center mb-4">
        Crea la cuenta de administrador
      </v-card-subtitle>

      <v-form @submit.prevent="submit">
        <v-text-field
          v-model="form.email"
          label="Email"
          type="email"
          :error-messages="errors.email"
          :disabled="loading"
          required
        />
        <v-text-field
          v-model="form.password"
          label="Contraseña"
          type="password"
          :error-messages="errors.password"
          :disabled="loading"
          required
        />
        <v-btn
          color="primary"
          block
          type="submit"
          :loading="loading"
          size="large"
          class="mt-2"
        >
          Crear cuenta
        </v-btn>
      </v-form>
    </v-card>
  </v-container>
</template>
