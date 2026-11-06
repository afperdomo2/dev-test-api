<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { getStatus, setup as setupApi } from '@/api/services/auth.service'
import type { SetupRequest } from '@/types/auth.types'
import {
  confirmPasswordRule,
  emailRule,
  passwordRule,
  requiredRule,
  validateRules,
} from '@/utils/validators'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const form = ref<SetupRequest & { confirmPassword: string }>({
  email: '',
  password: '',
  setupToken: '',
  confirmPassword: '',
})
const errors = ref<Record<string, Array<string>>>({})
const loading = ref(false)
const statusLoaded = ref(false)
const setupRequiresToken = ref(false)

onMounted(async () => {
  try {
    const status = await getStatus()
    setupRequiresToken.value = status.setupRequiresToken
  } catch {
    // Fail closed: sin estado no se permite enviar el formulario.
    setupRequiresToken.value = false
  } finally {
    statusLoaded.value = true
  }
})

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.email = validateRules([requiredRule(), emailRule()], form.value.email)
  newErrors.password = validateRules(
    [requiredRule(), passwordRule(8, 'Mínimo 8 caracteres')],
    form.value.password,
  )
  newErrors.confirmPassword = validateRules(
    [requiredRule('Confirma la contraseña'), confirmPasswordRule(form.value.password)],
    form.value.confirmPassword,
  )
  if (setupRequiresToken.value) {
    newErrors.setupToken = validateRules(
      [requiredRule('El token de instalación es obligatorio')],
      form.value.setupToken ?? '',
    )
  }
  errors.value = newErrors
  return Object.values(newErrors).every((e) => e.length === 0)
}

async function submit() {
  if (!statusLoaded.value || !setupRequiresToken.value) return
  if (!validate()) return
  loading.value = true
  try {
    const { confirmPassword: _confirm, ...payload } = form.value
    await setupApi(payload)
    authStore.needsSetup = false
    appStore.showSnackbar('Cuenta creada correctamente. Inicia sesión.', 'success')
    router.push('/login')
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
      <v-card-subtitle class="text-center mb-4"> Crea la cuenta de administrador </v-card-subtitle>

      <v-alert
        v-if="statusLoaded && !setupRequiresToken"
        type="warning"
        variant="tonal"
        class="mb-4"
        data-test="setup-disabled"
      >
        La configuración inicial está deshabilitada en este servidor. Pide al administrador que
        configure <code>SETUP_TOKEN</code>.
      </v-alert>

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
        <v-text-field
          v-model="form.confirmPassword"
          label="Confirmar contraseña"
          type="password"
          :error-messages="errors.confirmPassword"
          :disabled="loading"
          required
        />
        <v-text-field
          v-if="setupRequiresToken"
          v-model="form.setupToken"
          label="Token de instalación"
          type="password"
          autocomplete="off"
          :error-messages="errors.setupToken"
          :disabled="loading"
          data-test="setup-token"
          required
        />
        <v-btn
          color="primary"
          block
          type="submit"
          :loading="loading"
          :disabled="!statusLoaded || !setupRequiresToken"
          size="large"
          class="mt-2"
        >
          Crear cuenta
        </v-btn>
      </v-form>
    </v-card>
  </v-container>
</template>
