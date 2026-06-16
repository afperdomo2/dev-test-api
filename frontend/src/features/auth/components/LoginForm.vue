<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { login as loginApi } from '@/api/services/auth.service'
import type { LoginRequest } from '@/types/auth.types'
import { emailRule, requiredRule, validateRules } from '@/utils/validators'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

const form = ref<LoginRequest>({ email: '', password: '' })
const errors = ref<Record<string, Array<string>>>({})
const loading = ref(false)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.email = validateRules([requiredRule(), emailRule()], form.value.email)
  newErrors.password = validateRules([requiredRule()], form.value.password)
  errors.value = newErrors
  return Object.values(newErrors).every((e) => e.length === 0)
}

async function submit() {
  if (!validate()) return
  loading.value = true
  try {
    const res = await loginApi(form.value)
    authStore.setSession(res.token, res.user)
    const redirect = (route.query.redirect as string) ?? '/'
    router.push(redirect)
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al iniciar sesión'
    appStore.showSnackbar(detail, 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
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
      Iniciar sesión
    </v-btn>
  </v-form>
</template>
