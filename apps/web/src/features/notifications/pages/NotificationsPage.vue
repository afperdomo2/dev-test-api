<script setup lang="ts">
import { computed } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  markAllReadMutation,
  markReadMutation,
  notificationsOptions,
} from '@/queries/notifications.queries'
import { usePagination } from '@/composables/usePagination'
import { useAppStore } from '@/stores/app.store'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import type { Notification } from '@/types/notification.types'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset } = usePagination()

const { data, isLoading, error, isFetching } = useQuery(
  notificationsOptions(
    () => page.value,
    () => perPage.value,
  ),
)

const notifications = computed<Array<Notification>>(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)
const hasUnread = computed(() => notifications.value.some((n) => !n.isRead))

const markRead = useMutation(markReadMutation())
const markAll = useMutation(markAllReadMutation())

const relativeFormatter = new Intl.RelativeTimeFormat('es', { numeric: 'auto' })

function timeAgo(dateString: string): string {
  const diffMinutes = Math.round((new Date(dateString).getTime() - Date.now()) / 60_000)
  if (Math.abs(diffMinutes) < 60) return relativeFormatter.format(diffMinutes, 'minute')

  const diffHours = Math.round(diffMinutes / 60)
  if (Math.abs(diffHours) < 24) return relativeFormatter.format(diffHours, 'hour')

  return relativeFormatter.format(Math.round(diffHours / 24), 'day')
}

function handleRefresh() {
  reset()
  queryClient.invalidateQueries({ queryKey: ['notifications'] })
}

function handlePerPageChange(value: number) {
  perPage.value = value
  page.value = 1
}

async function handleMarkRead(id: string) {
  try {
    await markRead.mutateAsync({ id })
    await queryClient.invalidateQueries({ queryKey: ['notifications'] })
  } catch {
    appStore.showSnackbar('No se pudo marcar la notificación como leída', 'error')
  }
}

async function handleMarkAll() {
  try {
    await markAll.mutateAsync()
    await queryClient.invalidateQueries({ queryKey: ['notifications'] })
    appStore.showSnackbar('Notificaciones marcadas como leídas', 'success')
  } catch {
    appStore.showSnackbar('No se pudieron marcar las notificaciones como leídas', 'error')
  }
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Notificaciones"
      create-label=""
      :show-create="false"
      @refresh="handleRefresh"
    >
      <template #actions>
        <v-btn
          variant="tonal"
          color="primary"
          :disabled="!hasUnread"
          :loading="markAll.isPending.value"
          @click="handleMarkAll"
        >
          Marcar todas como leídas
        </v-btn>
      </template>
    </ListPageHeader>

    <v-skeleton-loader v-if="isLoading" type="list-item-two-line@4" />

    <v-alert v-else-if="error" type="error" variant="tonal">
      No se pudieron cargar las notificaciones.
      <template #append>
        <v-btn variant="text" @click="handleRefresh">Reintentar</v-btn>
      </template>
    </v-alert>

    <template v-else>
      <v-card variant="flat" border>
        <div v-if="notifications.length === 0" class="text-center py-8">
          <v-icon size="48" color="grey-lighten-1" class="mb-2">mdi-bell-off-outline</v-icon>
          <p class="text-body-1 text-medium-emphasis">No tienes notificaciones</p>
        </div>

        <v-list v-else lines="two">
          <v-list-item
            v-for="notification in notifications"
            :key="notification.id"
            :class="{ 'bg-surface-light': !notification.isRead }"
          >
            <template #prepend>
              <v-icon
                :icon="notification.isRead ? 'mdi-bell-outline' : 'mdi-bell-badge-outline'"
                :color="notification.isRead ? 'grey' : 'primary'"
                size="28"
              />
            </template>

            <v-list-item-title class="font-weight-medium">
              {{ notification.title }}
            </v-list-item-title>
            <v-list-item-subtitle>{{ notification.body }}</v-list-item-subtitle>

            <template #append>
              <div class="d-flex align-center ga-2 flex-wrap justify-end">
                <v-chip
                  v-if="!notification.isRead"
                  color="primary"
                  size="small"
                  variant="tonal"
                  prepend-icon="mdi-circle"
                >
                  Sin leer
                </v-chip>
                <span class="text-caption text-medium-emphasis">
                  {{ timeAgo(notification.createdAt) }}
                </span>
                <v-btn
                  v-if="!notification.isRead"
                  variant="text"
                  size="small"
                  :loading="markRead.isPending.value"
                  @click="handleMarkRead(notification.id)"
                >
                  Marcar leída
                </v-btn>
              </div>
            </template>
          </v-list-item>
        </v-list>
      </v-card>

      <PaginatedFooter
        :page="page"
        :per-page="perPage"
        :total="total"
        :in-table="false"
        @update:page="page = $event"
        @update:per-page="handlePerPageChange"
      />
    </template>

    <span v-if="isFetching && !isLoading" class="sr-only">Actualizando notificaciones</span>
  </v-container>
</template>
