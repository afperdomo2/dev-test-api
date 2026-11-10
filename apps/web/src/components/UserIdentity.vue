<script setup lang="ts">
import { computed } from 'vue'
import type { User } from '@/types/user.types'
import { roleLabel, userInitial } from '@/utils/identity'

const props = withDefaults(
  defineProps<{
    user: User | null
    variant?: 'icon' | 'sidebar'
  }>(),
  { variant: 'sidebar' },
)

const isAdmin = computed(() => props.user?.isAdmin ?? false)
const role = computed(() => roleLabel(isAdmin.value))
const initial = computed(() => userInitial(props.user?.email))

/**
 * En la app bar (`APP_COLOR` configurable) el avatar va blanco con la inicial
 * en `primary` para contrastar con cualquier preset. En el sidebar (fondo
 * `surface`) usa el color del rol, como `UserSummaryCard`.
 */
const avatarColor = computed(() => {
  if (props.variant === 'icon') return 'white'
  return isAdmin.value ? 'error' : 'info'
})

const initialColor = computed(() => (props.variant === 'icon' ? 'text-primary' : 'text-white'))
</script>

<template>
  <component :is="variant === 'icon' ? 'span' : 'div'" class="identity d-flex align-center ga-2">
    <v-avatar
      :color="avatarColor"
      :size="variant === 'icon' ? 28 : 32"
      :tag="variant === 'icon' ? 'span' : 'div'"
    >
      <span class="text-caption font-weight-bold" :class="initialColor">{{ initial }}</span>
    </v-avatar>

    <span v-if="variant === 'sidebar'" class="identity__text d-block">
      <span class="d-block text-body-2 text-truncate" :title="user?.email">
        {{ user?.email }}
      </span>
      <v-chip
        v-if="user"
        size="x-small"
        :color="isAdmin ? 'error' : 'info'"
        variant="tonal"
        :title="`Rol: ${role}`"
        class="mt-1"
      >
        <v-icon start size="14">{{ isAdmin ? 'mdi-shield-account' : 'mdi-account' }}</v-icon>
        {{ role }}
      </v-chip>
    </span>
  </component>
</template>

<style scoped>
.identity__text {
  min-width: 0;
}
</style>
