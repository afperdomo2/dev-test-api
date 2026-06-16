<script setup lang="ts">
import type { User } from '@/types/user.types'
import { formatDateTime } from '@/utils/format'

interface Props {
  users: Array<User>
  loading: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  delete: [user: User]
}>()

const headers = [
  { title: 'Email', key: 'email', sortable: false },
  { title: 'Rol', key: 'is_admin', sortable: false },
  { title: 'Creado', key: 'created_at', sortable: false },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'end' as const },
]
</script>

<template>
  <v-data-table
    :headers="headers"
    :items="users"
    :loading="loading"
    no-data-text="No hay usuarios"
    loading-text="Cargando usuarios..."
    hover
  >
    <template #[`item.is_admin`]="{ item }">
      <v-chip :color="item.is_admin ? 'primary' : ''" size="small" variant="tonal">
        {{ item.is_admin ? 'Admin' : 'Usuario' }}
      </v-chip>
    </template>

    <template #[`item.created_at`]="{ item }">
      {{ formatDateTime(item.created_at) }}
    </template>

    <template #[`item.actions`]="{ item }">
      <v-btn
        v-if="!item.is_admin"
        icon="mdi-delete"
        variant="text"
        size="small"
        color="error"
        @click="emit('delete', item)"
      />
    </template>

    <template #bottom>
      <slot name="footer" />
    </template>
  </v-data-table>
</template>
