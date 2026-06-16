<script setup lang="ts">
import type { User } from '@/types/user.types'
import { formatDateTime } from '@/utils/format'

interface Props {
  users: Array<User>
  loading: boolean
  itemsPerPage: number
}

defineProps<Props>()

const emit = defineEmits<{
  delete: [user: User]
}>()

const headers = [
  { title: 'Email', key: 'email', sortable: false },
  { title: 'Rol', key: 'isAdmin', sortable: false },
  { title: 'Creado', key: 'createdAt', sortable: false },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'end' as const },
]
</script>

<template>
  <v-data-table
    :headers="headers"
    :items="users"
    :items-per-page="itemsPerPage"
    :loading="loading"
    no-data-text="No hay usuarios"
    loading-text="Cargando usuarios..."
    hover
  >
    <template #[`item.isAdmin`]="{ item }">
      <v-chip :color="item.isAdmin ? 'primary' : ''" size="small" variant="tonal">
        {{ item.isAdmin ? 'Admin' : 'Usuario' }}
      </v-chip>
    </template>

    <template #[`item.createdAt`]="{ item }">
      {{ formatDateTime(item.createdAt) }}
    </template>

    <template #[`item.actions`]="{ item }">
      <v-btn
        v-if="!item.isAdmin"
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
