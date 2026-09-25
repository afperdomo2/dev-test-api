<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { reviewReportMutation } from '@/queries/moderation.queries'
import { useAppStore } from '@/stores/app.store'
import type { Report } from '@/types/moderation.types'

const props = defineProps<{
  modelValue: boolean
  report: Report | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  reviewed: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()

const unpublish = ref(false)
const saving = ref(false)

const reviewMut = useMutation(reviewReportMutation())

/** El backend solo des-publica preguntas; en tips el flag se ignora. */
const canUnpublish = computed(() => props.report?.targetType === 'question')

watch(
  () => props.modelValue,
  (open) => {
    if (open) unpublish.value = false
  },
)

async function submit() {
  if (!props.report) return
  saving.value = true
  try {
    await reviewMut.mutateAsync({
      id: props.report.id,
      data: { unpublish: canUnpublish.value ? unpublish.value : undefined },
    })
    appStore.showSnackbar(
      unpublish.value ? 'Reporte revisado y contenido des-publicado' : 'Reporte revisado',
    )
    queryClient.invalidateQueries({ queryKey: ['reports'] })
    emit('reviewed')
    close()
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al revisar el reporte'
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
    max-width="520"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>Revisar reporte</v-card-title>

      <v-card-text v-if="report">
        <v-list density="compact" class="py-0 mb-2">
          <v-list-item>
            <v-list-item-subheader>Motivo</v-list-item-subheader>
            <v-list-item-title class="text-body-2">{{ report.reason }}</v-list-item-title>
          </v-list-item>
          <v-list-item>
            <v-list-item-subheader>Contenido</v-list-item-subheader>
            <v-list-item-title class="text-body-2">
              {{ report.targetType === 'question' ? 'Pregunta' : 'Tip' }} ·
              <span class="text-medium-emphasis">{{ report.targetId }}</span>
            </v-list-item-title>
          </v-list-item>
        </v-list>

        <v-divider class="my-3" />

        <v-switch
          v-if="canUnpublish"
          v-model="unpublish"
          color="error"
          hide-details
          inset
          :disabled="saving"
          label="Des-publicar el contenido reportado"
        />
        <v-alert v-else type="info" variant="tonal" density="compact">
          Los tips son siempre públicos: la revisión solo marcará el reporte como atendido.
        </v-alert>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" :disabled="saving" @click="close"> Cancelar </v-btn>
        <v-btn color="primary" :loading="saving" @click="submit"> Marcar revisado </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
