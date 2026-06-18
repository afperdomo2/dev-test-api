<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createUserMutation, updateUserMutation } from '@/queries/users.queries'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import { emailRule, passwordRule, requiredRule, validateRules } from '@/utils/validators'
import type { User, CreateUserRequest } from '@/types/user.types'

const props = defineProps<{
  modelValue: boolean
  user: User | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()
const { extractFieldErrors } = useFormErrors()

const isEdit = computed(() => !!props.user)

const dialogTitle = computed(() => (isEdit.value ? 'Editar usuario' : 'Nuevo usuario'))

const form = ref<CreateUserRequest>({
  email: '',
  password: '',
  isAdmin: false,
})

const validationErrors = ref<Record<string, Array<string>>>({})
const serverErrors = ref<Record<string, string>>({})
const saving = ref(false)

const createMut = useMutation(createUserMutation())
const updateMut = useMutation(updateUserMutation())

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      if (props.user) {
        form.value = {
          email: props.user.email,
          password: '',
          isAdmin: props.user.isAdmin,
        }
      } else {
        form.value = { email: '', password: '', isAdmin: false }
      }
      validationErrors.value = {}
      serverErrors.value = {}
    }
  },
)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  if (!isEdit.value) {
    newErrors.email = validateRules([requiredRule(), emailRule()], form.value.email)
    newErrors.password = validateRules(
      [requiredRule(), passwordRule(8, 'Mínimo 8 caracteres')],
      form.value.password,
    )
  } else {
    if (form.value.password) {
      newErrors.password = validateRules(
        [passwordRule(8, 'Mínimo 8 caracteres')],
        form.value.password,
      )
    }
  }
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
  saving.value = true
  try {
    if (isEdit.value && props.user) {
      const data: Record<string, unknown> = { isAdmin: form.value.isAdmin }
      if (form.value.password) {
        data.password = form.value.password
      }
      await updateMut.mutateAsync({ id: props.user.id, data })
      appStore.showSnackbar('Usuario actualizado correctamente')
    } else {
      await createMut.mutateAsync({
        email: form.value.email,
        password: form.value.password,
        isAdmin: form.value.isAdmin,
      })
      appStore.showSnackbar('Usuario creado correctamente')
    }
    queryClient.invalidateQueries({ queryKey: ['users', 'list'] })
    emit('saved')
    close()
  } catch (err: unknown) {
    const fieldErrors = extractFieldErrors(err)
    if (Object.keys(fieldErrors).length > 0) {
      serverErrors.value = fieldErrors
    }
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al guardar usuario'
    appStore.showSnackbar(detail, 'error')
  } finally {
    saving.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    max-width="480"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>

      <v-card-text>
        <v-form @submit.prevent="submit">
          <v-text-field
            v-if="!isEdit"
            v-model="form.email"
            label="Email"
            type="email"
            :error-messages="fieldError('email')"
            :disabled="saving"
            required
          />

          <v-text-field
            v-model="form.password"
            :label="isEdit ? 'Nueva contraseña (opcional)' : 'Contraseña'"
            type="password"
            :error-messages="fieldError('password')"
            :disabled="saving"
            :required="!isEdit"
            :hint="isEdit ? 'Dejar vacío para mantener la actual' : ''"
            persistent-hint
          />

          <v-switch
            v-model="form.isAdmin"
            label="Administrador"
            color="primary"
            :disabled="saving"
            hide-details
            class="mb-4"
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" :disabled="saving" @click="close"> Cancelar </v-btn>
        <v-btn color="primary" :loading="saving" @click="submit">
          {{ isEdit ? 'Guardar' : 'Crear' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
