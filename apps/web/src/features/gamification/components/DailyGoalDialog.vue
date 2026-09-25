<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { setDailyGoalMutation } from '@/queries/gamification.queries'
import { useAppStore } from '@/stores/app.store'
import { useFormErrors } from '@/composables/useFormErrors'
import { DAILY_GOAL_MAX_XP, DAILY_GOAL_MIN_XP, DAILY_GOAL_OPTIONS } from '../constants'

const props = defineProps<{
  modelValue: boolean
  currentTarget: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const appStore = useAppStore()
const queryClient = useQueryClient()
const { extractFieldErrors } = useFormErrors()

const goalXp = ref<number>(DAILY_GOAL_MIN_XP)
const saving = ref(false)
const serverErrors = ref<Record<string, string>>({})
const clientError = ref('')

const goalMut = useMutation(setDailyGoalMutation())

const goalError = computed(() => {
  if (clientError.value) return clientError.value
  return serverErrors.value.dailyGoalXp ?? serverErrors.value.dailyGoalp ?? ''
})

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return
    goalXp.value = props.currentTarget || DAILY_GOAL_MIN_XP
    serverErrors.value = {}
    clientError.value = ''
  },
  { immediate: true },
)

function validate(): boolean {
  const value = Number(goalXp.value)
  if (!Number.isInteger(value)) {
    clientError.value = 'Valor inválido'
    return false
  }
  if (value < DAILY_GOAL_MIN_XP || value > DAILY_GOAL_MAX_XP) {
    clientError.value = `Entre ${DAILY_GOAL_MIN_XP} y ${DAILY_GOAL_MAX_XP} XP`
    return false
  }
  clientError.value = ''
  return true
}

async function submit() {
  if (!validate()) return
  saving.value = true
  serverErrors.value = {}
  try {
    await goalMut.mutateAsync({ dailyGoalXp: Number(goalXp.value) })
    appStore.showSnackbar('Meta diaria actualizada')
    queryClient.invalidateQueries({ queryKey: ['gamification'] })
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
        : 'Error al guardar la meta diaria'
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
    max-width="420"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>Meta diaria de XP</v-card-title>

      <v-card-text>
        <p class="text-body-2 text-medium-emphasis mb-4">
          Objetivo de XP que quieres alcanzar cada día. Se reinicia a medianoche (UTC).
        </p>

        <v-text-field
          v-model.number="goalXp"
          label="XP por día"
          type="number"
          :min="DAILY_GOAL_MIN_XP"
          :max="DAILY_GOAL_MAX_XP"
          :error-messages="goalError || undefined"
          :disabled="saving"
          hint="Entre 10 y 200 XP"
          persistent-hint
        />

        <div class="text-caption text-medium-emphasis mt-4 mb-1">Valores frecuentes</div>
        <div class="d-flex flex-wrap ga-2">
          <v-chip
            v-for="option in DAILY_GOAL_OPTIONS"
            :key="option"
            :color="goalXp === option ? 'primary' : undefined"
            :variant="goalXp === option ? 'flat' : 'tonal'"
            size="small"
            @click="goalXp = option"
          >
            {{ option }} XP
          </v-chip>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" :disabled="saving" @click="close"> Cancelar </v-btn>
        <v-btn color="primary" :loading="saving" @click="submit"> Guardar </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
