<script setup lang="ts">
import { ref, watch } from 'vue'

defineOptions({ name: 'PasswordField' })

const props = withDefaults(
  defineProps<{
    modelValue: string
    label?: string
    name?: string
    autocomplete?: string
    errorMessages?: string | Array<string>
    disabled?: boolean
    dataTest?: string
  }>(),
  {
    label: 'Contraseña',
    name: 'password',
    autocomplete: 'current-password',
    errorMessages: undefined,
    disabled: false,
    dataTest: undefined,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const visible = ref(false)

// Al enviar (campo deshabilitado) se vuelve a ocultar: la contraseña no debe
// quedarse en texto plano en pantalla.
watch(
  () => props.disabled,
  (disabled) => {
    if (disabled) visible.value = false
  },
)

function toggleVisibility() {
  visible.value = !visible.value
}
</script>

<template>
  <v-text-field
    :model-value="modelValue"
    :label="label"
    :name="name"
    :type="visible ? 'text' : 'password'"
    :autocomplete="autocomplete"
    :error-messages="errorMessages"
    :disabled="disabled"
    :data-test="dataTest"
    required
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #append-inner>
      <v-btn
        icon
        type="button"
        variant="text"
        density="comfortable"
        :aria-label="visible ? 'Ocultar contraseña' : 'Mostrar contraseña'"
        :disabled="disabled"
        :data-test="dataTest ? `${dataTest}-toggle` : undefined"
        @click="toggleVisibility"
      >
        <v-icon :icon="visible ? 'mdi-eye-off-outline' : 'mdi-eye-outline'" />
      </v-btn>
    </template>
  </v-text-field>
</template>
