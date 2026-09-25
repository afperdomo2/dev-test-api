<script setup lang="ts">
import { computed, ref } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { importTipsMutation } from '@/queries/tips.queries'
import { useAppStore } from '@/stores/app.store'
import { MAX_TIPS_PER_FILE, TIP_IMPORT_CSV_TEMPLATE } from '../constants'
import type { ImportResult } from '@/types/tip.types'

const emit = defineEmits<{
  imported: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()

const importMode = ref<'file' | 'paste'>('file')
const selectedFile = ref<File | null>(null)
const fileInputKey = ref(0)
const csvText = ref('')
const importing = ref(false)
const importResult = ref<ImportResult | null>(null)

const importMut = useMutation(importTipsMutation())

const canImport = computed(() => {
  if (importMode.value === 'file') return selectedFile.value !== null
  return csvText.value.trim().length > 0
})

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  selectedFile.value = input.files?.[0] ?? null
}

function downloadTemplate() {
  const blob = new Blob([TIP_IMPORT_CSV_TEMPLATE], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'plantilla-tips.csv'
  a.click()
  URL.revokeObjectURL(url)
}

function resetInput() {
  csvText.value = ''
  selectedFile.value = null
  fileInputKey.value++
}

async function doImport() {
  if (!canImport.value) {
    appStore.showSnackbar(
      importMode.value === 'file' ? 'Selecciona un archivo CSV' : 'Pega el contenido CSV',
      'error',
    )
    return
  }

  importing.value = true
  importResult.value = null
  try {
    const formData = new FormData()
    if (importMode.value === 'file' && selectedFile.value) {
      formData.append('file', selectedFile.value)
    } else {
      formData.append('content', csvText.value)
    }

    const result = await importMut.mutateAsync(formData)
    importResult.value = result
    queryClient.invalidateQueries({ queryKey: ['tips'] })
    emit('imported')

    if (result.imported > 0) {
      appStore.showSnackbar(`Se importaron ${result.imported} tips`)
      resetInput()
    }
    if (result.failed > 0) {
      appStore.showSnackbar(`${result.failed} filas con errores`, 'error')
    }
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al importar tips'
    appStore.showSnackbar(detail, 'error')
  } finally {
    importing.value = false
  }
}
</script>

<template>
  <div>
    <v-btn-toggle v-model="importMode" mandatory density="compact" color="primary" class="mb-3">
      <v-btn value="file" prepend-icon="mdi-file-delimited">Archivo</v-btn>
      <v-btn value="paste" prepend-icon="mdi-clipboard-text">Pegar texto</v-btn>
    </v-btn-toggle>

    <div v-if="importMode === 'file'" class="mb-3">
      <v-file-input
        :key="fileInputKey"
        label="Archivo CSV"
        accept=".csv,text/csv"
        :disabled="importing"
        hide-details
        density="compact"
        variant="outlined"
        @change="onFileChange"
      />
      <div class="text-caption text-medium-emphasis mt-1">
        Encabezado requerido: <code>content,topic</code> — el tema debe ser el
        <strong>slug</strong> de un tema existente. Máximo {{ MAX_TIPS_PER_FILE }} tips por archivo.
      </div>
    </div>

    <div v-else class="mb-3">
      <v-textarea
        v-model="csvText"
        label="Pega aquí el CSV"
        :disabled="importing"
        no-auto-grow
        rows="8"
        density="compact"
        variant="outlined"
        class="fixed-textarea"
        placeholder="content,topic&#10;Usa defer para liberar el archivo,go"
      />
    </div>

    <div class="d-flex ga-2 flex-wrap">
      <v-btn
        color="primary"
        :loading="importing"
        :disabled="!canImport"
        prepend-icon="mdi-upload"
        @click="doImport"
      >
        Importar tips
      </v-btn>
      <v-btn variant="tonal" prepend-icon="mdi-download" @click="downloadTemplate">
        Descargar plantilla
      </v-btn>
    </div>

    <div v-if="importResult" class="mt-4">
      <v-alert
        :type="importResult.failed === 0 ? 'success' : 'warning'"
        variant="tonal"
        class="mb-3"
      >
        Importados: {{ importResult.imported }} / {{ importResult.total }} — Fallidos:
        {{ importResult.failed }}
      </v-alert>
      <v-data-table
        v-if="importResult.errors.length"
        :headers="[
          { title: 'Fila', key: 'row' },
          { title: 'Motivo', key: 'reason' },
        ]"
        :items="importResult.errors"
        density="compact"
        hide-default-footer
        no-data-text="Sin errores"
      />
    </div>
  </div>
</template>

<style scoped>
.fixed-textarea :deep(textarea) {
  overflow-y: auto !important;
}
</style>
