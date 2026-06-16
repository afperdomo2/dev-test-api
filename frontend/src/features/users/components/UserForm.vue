<script setup lang="ts">
import { ref } from 'vue'
import type { CreateUserRequest } from '@/types/user.types'
import { emailRule, passwordRule, requiredRule, validateRules } from '@/utils/validators'
import { useFormErrors } from '@/composables/useFormErrors'

const emit = defineEmits<{
  submit: [data: CreateUserRequest]
}>()

const { extractFieldErrors } = useFormErrors()

const form = ref<CreateUserRequest>({
  email: '',
  password: '',
  is_admin: false,
})

const validationErrors = ref<Record<string, Array<string>>>({})
const serverErrors = ref<Record<string, string>>({})
const loading = ref(false)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.email = validateRules([requiredRule(), emailRule()], form.value.email)
  newErrors.password = validateRules(
    [requiredRule(), passwordRule(8, 'Mínimo 8 caracteres')],
    form.value.password,
  )
  validationErrors.value = newErrors
  serverErrors.value = {}
  return Object.values(newErrors).every((e) => e.length === 0)
}

function fieldError(field: string): Array<string> {
  const server = serverErrors.value[field]
  const client = validationErrors.value[field] ?? []
  return server ? [...client, server] : client
}

async function submit() {
  if (!validate()) return
  loading.value = true
  try {
    await emit('submit', { ...form.value })
  } catch (err: unknown) {
    const fieldErrors = extractFieldErrors(err)
    if (Object.keys(fieldErrors).length > 0) {
      serverErrors.value = fieldErrors
    }
    throw err
  } finally {
    loading.value = false
  }
}

function reset() {
  form.value = { email: '', password: '', is_admin: false }
  validationErrors.value = {}
  serverErrors.value = {}
}
</script>

<template>
  <v-form @submit.prevent="submit">
    <v-text-field
      v-model="form.email"
      label="Email"
      type="email"
      :error-messages="fieldError('email')"
      :disabled="loading"
      required
    />

    <v-text-field
      v-model="form.password"
      label="Contraseña"
      type="password"
      :error-messages="fieldError('password')"
      :disabled="loading"
      required
    />

    <v-switch
      v-model="form.is_admin"
      label="Administrador"
      color="primary"
      :disabled="loading"
      hide-details
      class="mb-4"
    />

    <div class="d-flex ga-2">
      <v-btn color="primary" type="submit" :loading="loading"> Crear usuario </v-btn>
      <v-btn variant="text" :disabled="loading" @click="reset"> Limpiar </v-btn>
    </div>
  </v-form>
</template>
