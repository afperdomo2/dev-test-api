<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { sessionsListOptions, createSessionMutation } from '@/queries/sessions.queries'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import SessionCard from '../components/SessionCard.vue'
import type { Session, CreateSessionRequest } from '@/types/session.types'
import { SESSION_MODES, SESSION_DIFFICULTIES } from '@/types/session.types'
import { requiredRule, validateRules } from '@/utils/validators'
import { topicsListOptions } from '@/queries/topics.queries'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const { data, isLoading } = useQuery(
  sessionsListOptions(
    () => page.value,
    () => perPage.value,
  ),
)

const sessionList = computed<Array<Session>>(() => {
  return data.value?.data ?? []
})

const totalSessions = computed(() => data.value?.meta?.total ?? 0)

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['sessions', 'list'] })
}

// Create dialog
const createDialog = ref(false)
const createForm = ref<CreateSessionRequest>({
  name: '',
  mode: 'generate',
  difficulty: 'beginner',
  topic_ids: [],
  question_limit: undefined,
})
const createErrors = ref<Record<string, Array<string>>>({})
const creating = ref(false)
const createMut = useMutation(createSessionMutation())

const { data: topicsData } = useQuery(
  topicsListOptions(
    () => 1,
    () => 100,
  ),
)
const topicItems = computed(() =>
  (topicsData.value?.data ?? []).map((t) => ({ title: t.name, value: t.id })),
)

function validateCreate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.name = validateRules([requiredRule()], createForm.value.name)
  newErrors.topic_ids = validateRules(
    [
      {
        validate: () => createForm.value.topic_ids.length > 0,
        message: 'Selecciona al menos un tema',
      },
    ],
    '',
  )
  if (createForm.value.question_limit !== undefined) {
    newErrors.question_limit = validateRules(
      [
        {
          validate: () => {
            const v = createForm.value.question_limit
            return v !== undefined && v >= 1
          },
          message: 'Minimo 1 pregunta',
        },
        {
          validate: () => {
            const v = createForm.value.question_limit
            return v !== undefined && v <= 50
          },
          message: 'Maximo 50 preguntas',
        },
      ],
      '',
    )
  }
  createErrors.value = newErrors
  return Object.values(newErrors).every((e) => e.length === 0)
}

async function handleCreate() {
  if (!validateCreate()) return
  creating.value = true
  try {
    await createMut.mutateAsync(createForm.value)
    queryClient.invalidateQueries({ queryKey: ['sessions', 'list'] })
    createDialog.value = false
    createForm.value = {
      name: '',
      mode: 'generate',
      difficulty: 'beginner',
      topic_ids: [],
      question_limit: undefined,
    }
    appStore.showSnackbar('Sesión creada')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al crear sesión'
    appStore.showSnackbar(detail, 'error')
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <v-container>
    <ListPageHeader
      title="Sesiones"
      create-label="Nueva sesión"
      @refresh="handleRefresh"
      @create="createDialog = true"
    />

    <v-card variant="flat" border class="mt-0">
      <v-row v-if="isLoading" class="ma-0">
        <v-col v-for="n in 4" :key="n" cols="12" sm="6">
          <v-skeleton-loader type="card" />
        </v-col>
      </v-row>

      <v-row v-else-if="sessionList.length" class="ma-0">
        <v-col v-for="session in sessionList" :key="session.id" cols="12" sm="6" lg="4">
          <SessionCard :session="session" />
        </v-col>
      </v-row>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1" class="mb-2"> mdi-play-circle-outline </v-icon>
        <p class="text-body-1 text-medium-emphasis">No hay sesiones aún</p>
      </v-card-text>

      <template v-if="!isLoading && sessionList.length">
        <PaginatedFooter
          :page="page"
          :per-page="perPage"
          :total="totalSessions"
          :in-table="true"
          @update:page="page = $event"
          @update:per-page="perPage = $event"
        />
      </template>
    </v-card>

    <!-- Create Session Dialog -->
    <v-dialog v-model="createDialog" max-width="480">
      <v-card>
        <v-card-title>Nueva sesión de estudio</v-card-title>
        <v-card-text>
          <v-form @submit.prevent="handleCreate">
            <v-text-field
              v-model="createForm.name"
              label="Nombre"
              :error-messages="createErrors.name"
              :disabled="creating"
              required
              placeholder="Ej: Repaso de Go"
            />

            <v-select
              v-model="createForm.mode"
              label="Modo"
              :items="SESSION_MODES"
              :disabled="creating"
              required
            />

            <v-select
              v-model="createForm.difficulty"
              label="Dificultad"
              :items="SESSION_DIFFICULTIES"
              :disabled="creating"
              required
            />

            <v-text-field
              v-model.number="createForm.question_limit"
              label="Limite de preguntas (opcional)"
              type="number"
              min="1"
              max="50"
              :error-messages="createErrors.question_limit"
              :disabled="creating"
              hint="Deja en blanco para preguntas ilimitadas"
              persistent-hint
            />

            <v-select
              v-model="createForm.topic_ids"
              label="Temas"
              :items="topicItems"
              :error-messages="createErrors.topic_ids"
              :disabled="creating"
              multiple
              chips
              required
              hint="Selecciona los temas a incluir"
              persistent-hint
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" :disabled="creating" @click="createDialog = false">
            Cancelar
          </v-btn>
          <v-btn color="primary" :loading="creating" @click="handleCreate"> Crear sesión </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
