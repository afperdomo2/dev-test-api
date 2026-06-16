<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createTopicMutation, updateTopicMutation } from '@/queries/topics.queries'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import { requiredRule, validateRules } from '@/utils/validators'
import { TOPIC_CATEGORIES } from '@/types/topic.types'
import type { Topic, CreateTopicRequest, UpdateTopicRequest } from '@/types/topic.types'

const props = defineProps<{
  modelValue: boolean
  topic: Topic | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()
const { extractFieldErrors } = useFormErrors()

const isEdit = computed(() => !!props.topic)

const dialogTitle = computed(() =>
  isEdit.value ? 'Editar tema' : 'Nuevo tema',
)

const form = ref<CreateTopicRequest>({
  slug: '',
  name: '',
  category: '',
})

const validationErrors = ref<Record<string, Array<string>>>({})
const serverErrors = ref<Record<string, string>>({})
const saving = ref(false)

const createMut = useMutation(createTopicMutation())
const updateMut = useMutation(updateTopicMutation())

const categoryItems: Array<string> = [...TOPIC_CATEGORIES]

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      if (props.topic) {
        form.value = {
          slug: props.topic.slug,
          name: props.topic.name,
          category: props.topic.category,
        }
      } else {
        form.value = { slug: '', name: '', category: '' }
      }
      validationErrors.value = {}
      serverErrors.value = {}
    }
  },
)

function validate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.slug = validateRules([requiredRule()], form.value.slug)
  newErrors.name = validateRules([requiredRule()], form.value.name)
  newErrors.category = validateRules([requiredRule()], form.value.category)
  validationErrors.value = newErrors
  serverErrors.value = {}
  return Object.values(newErrors).every((e) => e.length === 0)
}

function fieldError(field: string): Array<string> {
  const server = serverErrors.value[field]
  const client = validationErrors.value[field] ?? []
  return server ? [...client, server] : client
}

function slugify(value: string): string {
  return value
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^a-z0-9-]/g, '')
}

function onNameInput() {
  if (!isEdit.value && form.value.slug === slugify(form.value.name)) {
    return
  }
  if (!isEdit.value && !form.value.slug) {
    form.value.slug = slugify(form.value.name)
  }
}

async function submit() {
  if (!validate()) return
  saving.value = true
  try {
    if (isEdit.value && props.topic) {
      const data: UpdateTopicRequest = {
        slug: form.value.slug,
        name: form.value.name,
        category: form.value.category,
      }
      await updateMut.mutateAsync({ id: props.topic.id, data })
      appStore.showSnackbar('Tema actualizado')
    } else {
      await createMut.mutateAsync(form.value)
      appStore.showSnackbar('Tema creado')
    }
    queryClient.invalidateQueries({ queryKey: ['topics', 'list'] })
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
        : 'Error al guardar tema'
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
  <v-dialog :model-value="modelValue" max-width="480" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>

      <v-card-text>
        <v-form @submit.prevent="submit">
          <v-text-field
            v-model="form.name"
            label="Nombre"
            :error-messages="fieldError('name')"
            :disabled="saving"
            required
            @input="onNameInput"
          />

          <v-text-field
            v-model="form.slug"
            label="Slug"
            :error-messages="fieldError('slug')"
            :disabled="saving || isEdit"
            required
            hint="Identificador único (solo letras, números y guiones)"
            persistent-hint
          />

          <v-select
            v-model="form.category"
            label="Categoría"
            :items="categoryItems"
            :error-messages="fieldError('category')"
            :disabled="saving"
            required
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" :disabled="saving" @click="close"> Cancelar </v-btn>
        <v-btn
          color="primary"
          :loading="saving"
          @click="submit"
        >
          {{ isEdit ? 'Guardar' : 'Crear' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
