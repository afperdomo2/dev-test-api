<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import {
  sessionsListOptions,
  sessionsResumeOptions,
  createSessionMutation,
  finishSessionMutation,
  deleteSessionMutation,
} from '@/queries/sessions.queries'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import DifficultySelect from '@/components/DifficultySelect.vue'
import TopicAutocomplete from '@/components/TopicAutocomplete.vue'
import SessionRow from '../components/SessionRow.vue'
import type { Session, CreateSessionRequest, SessionStatus } from '@/types/session.types'
import {
  SESSION_STATUS_FILTERS,
  SESSION_SORT_OPTIONS,
  DEFAULT_SESSION_SORT,
} from '@/types/session.types'
import { requiredRule, validateRules } from '@/utils/validators'

const appStore = useAppStore()
const queryClient = useQueryClient()
const router = useRouter()

const { page, perPage, reset: resetPagination } = usePagination()

const selectedStatus = ref<SessionStatus | undefined>(undefined)
const selectedSortTitle = ref<string>(DEFAULT_SESSION_SORT.title)
const sortOption = computed(
  () =>
    SESSION_SORT_OPTIONS.find((option) => option.title === selectedSortTitle.value) ??
    DEFAULT_SESSION_SORT,
)
const sortTitles = SESSION_SORT_OPTIONS.map((option) => option.title)

const { data, isLoading, refetch } = useQuery(
  sessionsListOptions(
    () => page.value,
    () => perPage.value,
    () => selectedStatus.value,
    () => sortOption.value.sortBy,
    () => sortOption.value.sortOrder,
  ),
)

const { data: resumeData } = useQuery(sessionsResumeOptions())

const sessionList = computed<Array<Session>>(() => data.value?.data ?? [])
const totalSessions = computed(() => data.value?.meta?.total ?? 0)
const resumeSession = computed<Session | undefined>(() => resumeData.value?.data?.[0])

function selectStatus(status: SessionStatus | undefined) {
  selectedStatus.value = status
  resetPagination()
}

function selectSort(title: string) {
  selectedSortTitle.value = title
  resetPagination()
}

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['sessions', 'list'] })
  refetch()
}

function openSession(session: Session) {
  if (session.status === 'completed') {
    router.push(`/sessions/${session.id}/review`)
    return
  }
  router.push(`/sessions/${session.id}/study`)
}

const finishMut = useMutation(finishSessionMutation())
const deleteMut = useMutation(deleteSessionMutation())

const actionTarget = ref<Session | null>(null)
const actionType = ref<'finish' | 'delete' | null>(null)

const actionDialog = computed({
  get: () => actionType.value !== null,
  set: (value: boolean) => {
    if (!value) {
      actionType.value = null
      actionTarget.value = null
    }
  },
})

const actionPending = computed(() => finishMut.isPending.value || deleteMut.isPending.value)

const actionTitle = computed(() =>
  actionType.value === 'delete' ? '¿Eliminar sesión?' : '¿Finalizar sesión?',
)

const actionMessage = computed(() =>
  actionType.value === 'delete'
    ? 'Esta acción no se puede deshacer. La sesión se borrará permanentemente.'
    : 'Al finalizar se cerrará la sesión y podrás ver tu resumen de respuestas.',
)

function askFinish(session: Session) {
  actionTarget.value = session
  actionType.value = 'finish'
}

function askDelete(session: Session) {
  actionTarget.value = session
  actionType.value = 'delete'
}

function getErrorDetail(err: unknown, fallback: string) {
  return err && typeof err === 'object' && 'detail' in err
    ? (err as { detail: string }).detail
    : fallback
}

async function confirmAction() {
  if (!actionTarget.value || !actionType.value) return
  const action = actionType.value
  const target = actionTarget.value
  try {
    if (action === 'finish') {
      await finishMut.mutateAsync(target.id)
    } else {
      await deleteMut.mutateAsync(target.id)
    }
    await queryClient.invalidateQueries({ queryKey: ['sessions', 'list'] })
    appStore.showSnackbar(action === 'finish' ? 'Sesión finalizada' : 'Sesión eliminada')
    actionDialog.value = false
  } catch (err: unknown) {
    appStore.showSnackbar(
      getErrorDetail(
        err,
        action === 'finish' ? 'Error al finalizar sesión' : 'Error al eliminar sesión',
      ),
      'error',
    )
  }
}

const createDialog = ref(false)
function getDefaultCreateForm(): CreateSessionRequest {
  return {
    name: '',
    mode: 'review',
    difficulty: 'random',
    topicIds: [],
    questionLimit: undefined,
  }
}
const createForm = ref<CreateSessionRequest>(getDefaultCreateForm())
const createErrors = ref<Record<string, Array<string>>>({})
const creating = ref(false)
const createMut = useMutation(createSessionMutation())

function resetCreateForm() {
  createForm.value = getDefaultCreateForm()
  createErrors.value = {}
}

function openCreateDialog() {
  resetCreateForm()
  createDialog.value = true
}

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
    await queryClient.invalidateQueries({ queryKey: ['sessions', 'list'] })
    createDialog.value = false
    resetCreateForm()
    appStore.showSnackbar('Sesión creada')
  } catch (err: unknown) {
    appStore.showSnackbar(getErrorDetail(err, 'Error al crear sesión'), 'error')
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
      @create="openCreateDialog"
    />

    <v-card v-if="resumeSession" variant="flat" border class="mb-4">
      <div class="d-flex align-center flex-wrap ga-3 pa-3">
        <v-icon color="warning" size="24">mdi-progress-clock</v-icon>
        <div class="flex-grow-1" style="min-width: 200px">
          <div class="text-body-2 font-weight-medium">Continuar: {{ resumeSession.name }}</div>
          <div class="text-caption text-medium-emphasis">
            {{ resumeSession.answerCount }}
            <template v-if="resumeSession.questionLimit">
              / {{ resumeSession.questionLimit }}
            </template>
            preguntas respondidas
          </div>
        </div>
        <v-btn color="primary" size="small" class="text-none" @click="openSession(resumeSession)">
          Continuar
        </v-btn>
      </div>
    </v-card>

    <div class="d-flex flex-wrap align-center ga-2 mb-4">
      <v-chip
        v-for="filter in SESSION_STATUS_FILTERS"
        :key="filter.value ?? 'all'"
        :variant="selectedStatus === filter.value ? 'flat' : 'outlined'"
        :color="selectedStatus === filter.value ? 'primary' : undefined"
        filter
        @click="selectStatus(filter.value)"
      >
        {{ filter.title }}
      </v-chip>

      <v-spacer />

      <v-select
        :model-value="selectedSortTitle"
        :items="sortTitles"
        label="Ordenar por"
        density="compact"
        variant="outlined"
        hide-details
        style="max-width: 220px"
        @update:model-value="selectSort($event as string)"
      />
    </div>

    <v-card variant="flat" border>
      <div v-if="isLoading" class="pa-2">
        <v-skeleton-loader v-for="n in 6" :key="n" type="list-item-two-line" />
      </div>

      <template v-else-if="sessionList.length">
        <SessionRow
          v-for="session in sessionList"
          :key="session.id"
          :session="session"
          :busy="actionPending"
          @open="openSession"
          @finish="askFinish"
          @delete="askDelete"
        />
        <PaginatedFooter
          :page="page"
          :per-page="perPage"
          :total="totalSessions"
          :in-table="true"
          @update:page="page = $event"
          @update:per-page="onPerPageChange"
        />
      </template>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1" class="mb-2"> mdi-play-circle-outline </v-icon>
        <p class="text-body-1 text-medium-emphasis">
          {{ selectedStatus ? 'No hay sesiones en este estado' : 'No hay sesiones aún' }}
        </p>
      </v-card-text>
    </v-card>

    <v-dialog v-model="actionDialog" max-width="420">
      <v-card>
        <v-card-title>{{ actionTitle }}</v-card-title>
        <v-card-text>{{ actionMessage }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" :disabled="actionPending" @click="actionDialog = false">
            Cancelar
          </v-btn>
          <v-btn
            :color="actionType === 'delete' ? 'error' : 'primary'"
            :loading="actionPending"
            @click="confirmAction"
          >
            {{ actionType === 'delete' ? 'Eliminar' : 'Finalizar' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="createDialog" max-width="640">
      <v-card>
        <v-card-title>Nueva sesión de estudio</v-card-title>
        <v-card-text>
          <v-form autocomplete="off" @submit.prevent="handleCreate">
            <v-text-field
              v-model="createForm.name"
              label="Nombre"
              autocomplete="off"
              :error-messages="createErrors.name"
              :disabled="creating"
              required
              placeholder="Ej: Repaso de Go"
            />

            <v-row dense>
              <v-col cols="12" sm="6">
                <DifficultySelect v-model="createForm.difficulty" :disabled="creating" />
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model.number="createForm.questionLimit"
                  label="Límite de preguntas (opcional)"
                  type="number"
                  min="1"
                  max="50"
                  autocomplete="off"
                  density="compact"
                  variant="outlined"
                  :error-messages="createErrors.questionLimit"
                  :disabled="creating"
                  hint="Deja en blanco para ilimitadas"
                  persistent-hint
                />
              </v-col>
            </v-row>

            <TopicAutocomplete
              v-model="createForm.topicIds"
              :active="createDialog"
              label="Temas"
              has-questions
              autocomplete="off"
              :error-messages="createErrors.topicIds"
              :disabled="creating"
              required
              hint="Solo se muestran los temas que tienen preguntas"
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
