<script setup lang="ts">
import { ref } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { changePasswordMutation } from '@/queries/users.queries'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import {
  confirmPasswordRule,
  maxLengthRule,
  passwordRule,
  requiredRule,
  validateRules,
} from '@/utils/validators'

const appStore = useAppStore()
const { extractFieldErrors } = useFormErrors()
const changeMut = useMutation(changePasswordMutation())

const form = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const validationErrors = ref<Record<string, Array<string>>>({})
const serverErrors = ref<Record<string, string>>({})
const serverError = ref<string | null>(null)
const showCurrent = ref(false)
const showNew = ref(false)
const showConfirm = ref(false)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}

  newErrors.currentPassword = validateRules(
    [requiredRule('Ingresa tu contraseña actual'), maxLengthRule(72)],
    form.value.currentPassword,
  )

  newErrors.newPassword = validateRules(
    [
      requiredRule('Ingresa la nueva contraseña'),
      passwordRule(8, 'Mínimo 8 caracteres'),
      maxLengthRule(72),
    ],
    form.value.newPassword,
  )
  if (form.value.newPassword && form.value.newPassword === form.value.currentPassword) {
    newErrors.newPassword.push('La nueva contraseña debe ser diferente a la actual')
  }

  newErrors.confirmPassword = validateRules(
    [requiredRule('Confirma la nueva contraseña'), confirmPasswordRule(form.value.newPassword)],
    form.value.confirmPassword,
  )

  validationErrors.value = newErrors
  serverErrors.value = {}
  serverError.value = null

  return Object.values(newErrors).every((errors) => errors.length === 0)
}

function fieldError(field: string): Array<string> {
  const server = serverErrors.value[field]
  const client = validationErrors.value[field] ?? []
  return server ? [...client, server] : client
}

function resetForm() {
  form.value = { currentPassword: '', newPassword: '', confirmPassword: '' }
  validationErrors.value = {}
  serverErrors.value = {}
  showCurrent.value = false
  showNew.value = false
  showConfirm.value = false
}

async function submit() {
  if (!validate()) return

  try {
    await changeMut.mutateAsync({
      currentPassword: form.value.currentPassword,
      newPassword: form.value.newPassword,
    })
    resetForm()
    appStore.showSnackbar('Contraseña actualizada correctamente')
  } catch (err: unknown) {
    const fieldErrors = extractFieldErrors(err)
    if (Object.keys(fieldErrors).length > 0) {
      serverErrors.value = fieldErrors
    }
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al actualizar la contraseña'
    serverError.value = detail
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-card>
    <v-card-item>
      <template #prepend>
        <v-icon color="primary" size="40">mdi-lock-reset</v-icon>
      </template>
      <v-card-title>Cambiar contraseña</v-card-title>
      <v-card-subtitle>Mínimo 8 caracteres. Tu sesión se mantiene activa.</v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <v-alert
        v-if="serverError"
        type="error"
        variant="tonal"
        density="compact"
        class="mb-4"
        :text="serverError"
      />

      <v-form @submit.prevent="submit">
        <v-text-field
          v-model="form.currentPassword"
          label="Contraseña actual"
          :type="showCurrent ? 'text' : 'password'"
          :append-inner-icon="showCurrent ? 'mdi-eye-off' : 'mdi-eye'"
          :error-messages="fieldError('currentPassword')"
          :disabled="changeMut.isPending.value"
          autocomplete="current-password"
          @click:append-inner="showCurrent = !showCurrent"
        />

        <v-text-field
          v-model="form.newPassword"
          label="Nueva contraseña"
          :type="showNew ? 'text' : 'password'"
          :append-inner-icon="showNew ? 'mdi-eye-off' : 'mdi-eye'"
          :error-messages="fieldError('newPassword')"
          :disabled="changeMut.isPending.value"
          autocomplete="new-password"
          @click:append-inner="showNew = !showNew"
        />

        <v-text-field
          v-model="form.confirmPassword"
          label="Confirmar nueva contraseña"
          :type="showConfirm ? 'text' : 'password'"
          :append-inner-icon="showConfirm ? 'mdi-eye-off' : 'mdi-eye'"
          :error-messages="fieldError('confirmPassword')"
          :disabled="changeMut.isPending.value"
          autocomplete="new-password"
          @click:append-inner="showConfirm = !showConfirm"
        />
      </v-form>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <v-btn
        color="primary"
        :loading="changeMut.isPending.value"
        :disabled="changeMut.isPending.value"
        @click="submit"
      >
        Actualizar contraseña
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
