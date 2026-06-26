<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { sessionsInfiniteOptions, createSessionMutation } from '@/queries/sessions.queries'
import { useAppStore } from '@/stores/app.store'
import ListPageHeader from '@/components/ListPageHeader.vue'
import SessionCard from '../components/SessionCard.vue'
import type { Session, CreateSessionRequest, SessionStatus } from '@/types/session.types'
import { SESSION_MODES, SESSION_DIFFICULTIES, SESSION_STATUS_FILTERS } from '@/types/session.types'
import { requiredRule, validateRules } from '@/utils/validators'
import { listTopics } from '@/api/services/topics.service'

const appStore = useAppStore()
const queryClient = useQueryClient()

const selectedStatus = ref<SessionStatus | undefined>(undefined)

const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useInfiniteQuery(
  sessionsInfiniteOptions(() => selectedStatus.value),
)

const sessionList = computed<Array<Session>>(() => data.value?.pages.flatMap((p) => p.data) ?? [])

function onIntersect(isIntersecting: boolean) {
  if (isIntersecting && hasNextPage.value && !isFetchingNextPage.value) {
    fetchNextPage()
  }
}

function handleRefresh() {
  queryClient.invalidateQueries({ queryKey: ['sessions', 'list', 'infinite'] })
}

const createDialog = ref(false)
const createForm = ref<CreateSessionRequest>({
  name: '',
  mode: 'generate',
  difficulty: 'beginner',
  topicIds: [],
  questionLimit: undefined,
})
const createErrors = ref<Record<string, Array<string>>>({})
const creating = ref(false)
const createMut = useMutation(createSessionMutation())

const { data: topicsData, isLoading: topicsLoading } = useQuery({
  queryKey: ['topics', 'list', 1, 100, 'name', 'asc'],
  queryFn: () => listTopics(1, 100, 'name', 'asc'),
  staleTime: 60 * 1000,
  enabled: computed(() => createDialog.value),
})
const topicItems = computed(() =>
  (topicsData.value?.data ?? []).map((t) => ({
    title: t.name,
    value: t.id,
    props: { subtitle: t.category },
  })),
)

function validateCreate(): boolean {
  const newErrors: Record<string, Array<string>> = {}
  newErrors.name = validateRules([requiredRule()], createForm.value.name)
  newErrors.topicIds = validateRules(
    [
      {
        validate: () => createForm.value.topicIds.length > 0,
        message: 'Selecciona al menos un tema',
      },
    ],
    '',
  )
  if (createForm.value.questionLimit !== undefined) {
    newErrors.questionLimit = validateRules(
      [
        {
          validate: () => {
            const v = createForm.value.questionLimit
            return v !== undefined && v >= 1
          },
          message: 'Minimo 1 pregunta',
        },
        {
          validate: () => {
            const v = createForm.value.questionLimit
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
    queryClient.invalidateQueries({ queryKey: ['sessions', 'list', 'infinite'] })
    createDialog.value = false
    createForm.value = {
      name: '',
      mode: 'generate',
      difficulty: 'beginner',
      topicIds: [],
      questionLimit: undefined,
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

    <div class="d-flex ga-2 mb-4 flex-wrap">
      <v-chip
        v-for="filter in SESSION_STATUS_FILTERS"
        :key="filter.value ?? 'all'"
        :variant="selectedStatus === filter.value ? 'flat' : 'outlined'"
        :color="selectedStatus === filter.value ? 'primary' : undefined"
        filter
        @click="selectedStatus = filter.value"
      >
        {{ filter.title }}
      </v-chip>
    </div>

    <v-row v-if="isLoading">
      <v-col v-for="n in 4" :key="n" cols="12" sm="6">
        <v-skeleton-loader type="card" />
      </v-col>
    </v-row>

    <v-row v-else-if="sessionList.length">
      <v-col v-for="session in sessionList" :key="session.id" cols="12" sm="6" lg="4">
        <SessionCard :session="session" />
      </v-col>
    </v-row>

    <div v-else class="text-center py-8">
      <v-icon size="48" color="grey-lighten-1" class="mb-2"> mdi-play-circle-outline </v-icon>
      <p class="text-body-1 text-medium-emphasis">No hay sesiones aún</p>
    </div>

    <div v-intersect="onIntersect" class="d-flex justify-center py-4">
      <v-progress-circular v-if="isFetchingNextPage" indeterminate color="primary" size="32" />
      <span
        v-else-if="!hasNextPage && sessionList.length > 0"
        class="text-caption text-medium-emphasis"
      >
        No hay más sesiones
      </span>
    </div>

    <v-dialog v-model="createDialog" max-width="640">
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

            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="createForm.mode"
                  label="Modo"
                  :items="SESSION_MODES"
                  :disabled="creating"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="createForm.difficulty"
                  label="Dificultad"
                  :items="SESSION_DIFFICULTIES"
                  :disabled="creating"
                  required
                />
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="createForm.questionLimit"
                  label="Limite de preguntas (opcional)"
                  type="number"
                  min="1"
                  max="50"
                  :error-messages="createErrors.questionLimit"
                  :disabled="creating"
                  hint="Deja en blanco para preguntas ilimitadas"
                  persistent-hint
                />
              </v-col>
            </v-row>

            <v-autocomplete
              v-model="createForm.topicIds"
              label="Temas"
              :items="topicItems"
              item-title="title"
              item-props="props"
              :error-messages="createErrors.topicIds"
              :disabled="creating"
              :loading="topicsLoading"
              multiple
              chips
              clearable
              required
              hint="Busca y selecciona los temas a incluir"
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
