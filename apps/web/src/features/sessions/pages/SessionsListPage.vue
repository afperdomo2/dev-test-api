<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { sessionsInfiniteOptions, createSessionMutation } from '@/queries/sessions.queries'
import { getAiQuota } from '@/api/services/questions.service'
import { useAppStore } from '@/stores/app.store'
import { useAuthStore } from '@/stores/auth.store'
import ListPageHeader from '@/components/ListPageHeader.vue'
import SessionCard from '../components/SessionCard.vue'
import type { Session, CreateSessionRequest, SessionStatus } from '@/types/session.types'
import {
  SESSION_MODE_DESCRIPTIONS,
  SESSION_MODE_ICONS,
  SESSION_MODES,
  SESSION_DIFFICULTIES,
  SESSION_STATUS_FILTERS,
} from '@/types/session.types'
import { requiredRule, validateRules } from '@/utils/validators'
import TopicAutocomplete from '@/components/TopicAutocomplete.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
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

const { data: aiQuota } = useQuery({
  queryKey: ['questions', 'ai-quota'],
  queryFn: () => getAiQuota(),
  staleTime: 30 * 1000,
  enabled: computed(
    () => createDialog.value && !authStore.isAdmin && createForm.value.mode === 'generate',
  ),
})

const aiQuotaPercent = computed(() => {
  if (!aiQuota.value || aiQuota.value.dailyLimit === 0) return 0
  return Math.min(100, Math.round((aiQuota.value.usedToday / aiQuota.value.dailyLimit) * 100))
})

const aiQuotaColor = computed(() => {
  const p = aiQuotaPercent.value
  if (p >= 90) return 'error'
  if (p >= 70) return 'warning'
  return 'success'
})

const aiQuotaExceeds = computed(() => {
  if (
    createForm.value.mode !== 'generate' ||
    !aiQuota.value ||
    createForm.value.questionLimit === undefined
  )
    return false
  return createForm.value.questionLimit > aiQuota.value.remaining
})

watch(
  () => createForm.value.mode,
  (mode) => {
    if (mode === 'generate' && !authStore.isAdmin) {
      queryClient.invalidateQueries({ queryKey: ['questions', 'ai-quota'] })
    }
  },
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
    queryClient.invalidateQueries({ queryKey: ['questions', 'ai-quota'] })
    createDialog.value = false
    resetCreateForm()
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
      @create="openCreateDialog"
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

            <div class="mb-4">
              <p class="text-subtitle-2 mb-2">Modo</p>
              <v-radio-group
                v-model="createForm.mode"
                :disabled="creating"
                hide-details
                class="mt-0"
              >
                <v-card
                  v-for="mode in SESSION_MODES"
                  :key="mode.value"
                  :variant="createForm.mode === mode.value ? 'tonal' : 'outlined'"
                  :color="createForm.mode === mode.value ? 'primary' : undefined"
                  class="mb-2 cursor-pointer"
                  @click="createForm.mode = mode.value"
                >
                  <v-card-text class="d-flex align-center ga-3 py-3">
                    <v-radio
                      :value="mode.value"
                      :disabled="creating"
                      hide-details
                      density="compact"
                      color="primary"
                    />
                    <v-icon :icon="SESSION_MODE_ICONS[mode.value]" size="28" color="primary" />
                    <div>
                      <div class="text-body-2 font-weight-medium">{{ mode.title }}</div>
                      <div class="text-caption text-medium-emphasis">
                        {{ SESSION_MODE_DESCRIPTIONS[mode.value] }}
                      </div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-radio-group>
            </div>

            <v-alert
              v-if="createForm.mode === 'generate' && aiQuota"
              :type="aiQuota.remaining === 0 ? 'error' : aiQuotaExceeds ? 'warning' : 'info'"
              variant="tonal"
              density="compact"
              class="mb-4"
            >
              <div class="d-flex align-center justify-space-between flex-wrap ga-2">
                <span class="text-body-2">
                  <v-icon start size="18">mdi-robot</v-icon>
                  Cupo IA hoy: {{ aiQuota.usedToday }} / {{ aiQuota.dailyLimit }} usadas
                  <span class="font-weight-bold">— Restan {{ aiQuota.remaining }}</span>
                </span>
              </div>
              <v-progress-linear
                :model-value="aiQuotaPercent"
                :color="aiQuotaColor"
                height="16"
                rounded
                class="mt-2"
              >
                <template #default>
                  <span class="text-caption font-weight-medium">{{ aiQuotaPercent }}%</span>
                </template>
              </v-progress-linear>
              <div v-if="aiQuota.remaining === 0" class="text-caption mt-2">
                Has agotado tu cupo de IA hoy. Podrás generar nuevas preguntas mañana. Las preguntas
                existentes seguirán disponibles.
              </div>
              <div v-else-if="aiQuotaExceeds" class="text-caption mt-2">
                Pediste {{ createForm.questionLimit }} preguntas pero solo te quedan
                {{ aiQuota.remaining }} nuevas con IA hoy. El resto se completará con preguntas
                existentes o la sesión terminará antes.
              </div>
              <div v-else class="text-caption mt-2">
                La IA generará preguntas hasta agotar tu cupo diario.
              </div>
            </v-alert>

            <v-select
              v-model="createForm.difficulty"
              label="Dificultad"
              :items="SESSION_DIFFICULTIES"
              item-props="props"
              :disabled="creating"
              required
            />

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

            <TopicAutocomplete
              v-model="createForm.topicIds"
              :active="createDialog"
              label="Temas"
              :error-messages="createErrors.topicIds"
              :disabled="creating"
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
