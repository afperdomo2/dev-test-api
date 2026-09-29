<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { settingsListOptions, updateSettingsMutation } from '@/queries/generation.queries'
import { useAppStore } from '@/stores/app.store'
import ListPageHeader from '@/components/ListPageHeader.vue'
import type { Setting } from '@/types/generation.types'

const appStore = useAppStore()
const queryClient = useQueryClient()

const { data, isLoading } = useQuery(settingsListOptions())
const updateSettings = useMutation(updateSettingsMutation())

/** Borrador editable: el servidor no se toca hasta pulsar Guardar. */
const draft = ref<Record<string, string>>({})
const initialized = ref(false)

const settings = computed<Array<Setting>>(() => data.value ?? [])

watch(
  settings,
  (list) => {
    if (initialized.value || list.length === 0) return
    draft.value = Object.fromEntries(list.map((s) => [s.key, s.value]))
    initialized.value = true
  },
  { immediate: true },
)

/** Solo los ajustes que el admin editó de verdad, para no reenviar el resto. */
const changedSettings = computed(() => {
  const out: Array<{ key: string; value: string }> = []
  for (const setting of settings.value) {
    const value = draft.value[setting.key]
    if (value === undefined || value === setting.value) continue
    out.push({ key: setting.key, value })
  }
  return out
})

const hasChanges = computed(() => changedSettings.value.length > 0)

function validate(setting: Setting, raw: string): string | null {
  if (raw.trim() === '') return 'Obligatorio'

  if (setting.type !== 'string') {
    const parsed = Number(raw)
    if (Number.isNaN(parsed)) return 'Debe ser un número'
    if (setting.min !== undefined && parsed < setting.min) return `Mínimo ${setting.min}`
    if (setting.max !== undefined && parsed > setting.max) return `Máximo ${setting.max}`
    if (setting.type === 'int' && !Number.isInteger(parsed)) return 'Debe ser un entero'
  }
  return null
}

const invalidKeys = computed(() => {
  const out: Record<string, string> = {}
  for (const setting of settings.value) {
    const raw = draft.value[setting.key]
    if (raw === undefined) continue
    const error = validate(setting, raw)
    if (error) out[setting.key] = error
  }
  return out
})

const canSave = computed(() => hasChanges.value && Object.keys(invalidKeys.value).length === 0)

function fieldProps(setting: Setting) {
  if (setting.type === 'string') {
    return { type: 'text' as const }
  }
  return {
    type: 'number' as const,
    min: setting.min,
    max: setting.max,
    step: setting.type === 'int' ? 1 : 0.1,
  }
}

function reset() {
  draft.value = Object.fromEntries(settings.value.map((s) => [s.key, s.value]))
}

async function save() {
  if (!canSave.value) return
  try {
    await updateSettings.mutateAsync({ settings: changedSettings.value })
    appStore.showSnackbar('Configuración actualizada')
    await queryClient.invalidateQueries({ queryKey: ['settings'] })
  } catch (error) {
    const detail =
      error && typeof error === 'object' && 'detail' in error
        ? (error as { detail: string }).detail
        : 'Error al actualizar la configuración'
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Configuración de IA"
      create-label=""
      :show-create="false"
      @refresh="reset"
    >
      <template #actions>
        <v-btn
          color="primary"
          prepend-icon="mdi-content-save"
          :loading="updateSettings.isPending.value"
          :disabled="!canSave"
          @click="save"
        >
          Guardar
        </v-btn>
      </template>
    </ListPageHeader>

    <v-alert type="info" variant="tonal" density="compact" class="mb-4">
      Las claves de API y la URL del proveedor se siguen leyendo del entorno del servidor
      (<code>AI_API_KEY</code>, <code>AI_API_URL</code>); no se guardan aquí.
    </v-alert>

    <v-card>
      <v-card-text>
        <v-row dense>
          <v-col v-for="setting in settings" :key="setting.key" cols="12" md="6">
            <v-text-field
              v-model="draft[setting.key]"
              :label="setting.key"
              :hint="setting.description"
              persistent-hint
              :error-messages="invalidKeys[setting.key]"
              :disabled="updateSettings.isPending.value"
              variant="outlined"
              density="compact"
              v-bind="fieldProps(setting)"
            />
          </v-col>
        </v-row>

        <v-progress-linear v-if="isLoading" indeterminate color="primary" class="mt-2" />
      </v-card-text>
    </v-card>
  </v-container>
</template>
