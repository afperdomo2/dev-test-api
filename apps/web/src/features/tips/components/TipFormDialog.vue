<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createTipMutation, updateTipMutation } from '@/queries/tips.queries'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import { maxLengthRule, requiredRule, validateRules } from '@/utils/validators'
import TopicSelect from '@/components/TopicSelect.vue'
import type { CreateTipRequest, Tip, UpdateTipRequest } from '@/types/tip.types'

const props = defineProps<{
  modelValue: boolean
  tip: Tip | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()
const { extractFieldErrors } = useFormErrors()

const isEdit = computed(() => !!props.tip)

const dialogTitle = computed(() => (isEdit.value ? 'Editar tip' : 'Nuevo tip'))

const form = ref<CreateTipRequest>({
  topicId: '',
  content: '',
})

const validationErrors = ref<Record<string, Array<string>>>({})
const serverErrors = ref<Record<string, string>>({})
const saving = ref(false)

const createMut = useMutation(createTipMutation())
const updateMut = useMutation(updateTipMutation())

/** `TopicSelect` es clearable y emite `null`; el request exige un tema. */
const topicIdModel = computed<string | null>({
  get: () => form.value.topicId || null,
  set: (value) => {
    form.value.topicId = value ?? ''
  },
})

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return
    if (props.tip) {
      form.value = { topicId: props.tip.topicId, content: props.tip.content }
    } else {
      form.value = { topicId: '', content: '' }
    }
    validationErrors.value = {}
    serverErrors.value = {}
  },
)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.topicId = validateRules(
    [requiredRule('Selecciona un tema')],
    form.value.topicId ? 'filled' : '',
  )
  newErrors.content = validateRules(
    [requiredRule('El contenido es requerido'), maxLengthRule(500, 'Máximo 500 caracteres')],
    form.value.content,
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
  saving.value = true
  try {
    if (isEdit.value && props.tip) {
      const data: UpdateTipRequest = {
        topicId: form.value.topicId,
        content: form.value.content,
      }
      await updateMut.mutateAsync({ id: props.tip.id, data })
      appStore.showSnackbar('Tip actualizado')
    } else {
      await createMut.mutateAsync({ ...form.value })
      appStore.showSnackbar('Tip creado')
    }
    queryClient.invalidateQueries({ queryKey: ['tips'] })
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
        : 'Error al guardar el tip'
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
    max-width="560"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>

      <v-card-text>
        <v-form @submit.prevent="submit">
          <TopicSelect
            v-model="topicIdModel"
            :error-messages="fieldError('topicId')"
            :disabled="saving"
            :clearable="false"
            required
          />

          <v-textarea
            v-model="form.content"
            label="Contenido"
            :error-messages="fieldError('content')"
            :disabled="saving"
            rows="4"
            auto-grow
            counter="500"
            required
            hint="Los tips son públicos para todos los usuarios"
            persistent-hint
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
