<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { finishSessionMutation, deleteSessionMutation } from '@/queries/sessions.queries'
import { useAppStore } from '@/stores/app.store'
import type { Session } from '@/types/session.types'
import {
  SESSION_STATUS_LABELS,
  SESSION_STATUS_COLORS,
  SESSION_MODE_LABELS,
  SESSION_DIFFICULTY_LABELS,
} from '@/types/session.types'
import { formatDate, formatScore } from '@/utils/format'
import { DIFFICULTY_COLORS } from '@/types/question.types'

interface Props {
  session: Session
}

const props = defineProps<Props>()
const router = useRouter()
const queryClient = useQueryClient()
const appStore = useAppStore()

const finishMut = useMutation(finishSessionMutation())
const deleteMut = useMutation(deleteSessionMutation())

const finishDialog = ref(false)
const deleteDialog = ref(false)

const canDelete = computed(() => {
  if (props.session.answerCount > 0) return false
  const created = new Date(props.session.createdAt)
  const limit = Date.now() - 24 * 60 * 60 * 1000
  return created.getTime() > limit
})

const isGenerating = computed(
  () =>
    props.session.mode === 'generate' &&
    props.session.status === 'in_progress' &&
    props.session.questionsGenerated < 2,
)

const displayedTopics = computed(() => props.session.topics.slice(0, 2))
const hiddenCount = computed(() => Math.max(0, props.session.topics.length - 2))

const STATUS_ICONS: Record<string, string> = {
  in_progress: 'mdi-progress-clock',
  completed: 'mdi-check-circle',
  cancelled: 'mdi-cancel',
}

const MODE_ICONS: Record<string, string> = {
  generate: 'mdi-auto-fix',
  review: 'mdi-book-clock-outline',
}

const MODE_COLORS: Record<string, string> = {
  generate: 'deep-purple',
  review: 'blue',
}

const DIFFICULTY_ICONS: Record<string, string> = {
  beginner: 'mdi-signal-cellular-1',
  intermediate: 'mdi-signal-cellular-2',
  advanced: 'mdi-signal-cellular-3',
}

async function handleFinish() {
  finishDialog.value = false
  try {
    await finishMut.mutateAsync(props.session.id)
    queryClient.invalidateQueries({ queryKey: ['sessions', 'list', 'infinite'] })
    appStore.showSnackbar('Sesión finalizada')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al finalizar sesión'
    appStore.showSnackbar(detail, 'error')
  }
}

async function handleDelete() {
  deleteDialog.value = false
  try {
    await deleteMut.mutateAsync(props.session.id)
    queryClient.invalidateQueries({ queryKey: ['sessions', 'list', 'infinite'] })
    appStore.showSnackbar('Sesión eliminada')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al eliminar sesión'
    appStore.showSnackbar(detail, 'error')
  }
}

function handleClick() {
  if (isGenerating.value) {
    appStore.showSnackbar('Las preguntas aún se están generando. Espera unos segundos.', 'info')
    return
  }
  if (props.session.status === 'completed') {
    router.push(`/sessions/${props.session.id}/review`)
    return
  }
  router.push(`/sessions/${props.session.id}/study`)
}

const cardClasses = computed(() => ({
  'session-card': true,
  'session-card--blocked': isGenerating.value,
}))
</script>

<template>
  <v-card :class="cardClasses" hover :ripple="!isGenerating" @click="handleClick">
    <v-card-item>
      <template #prepend>
        <v-icon :color="MODE_COLORS[session.mode]" size="32">
          {{ MODE_ICONS[session.mode] }}
        </v-icon>
      </template>
      <v-card-title>{{ session.name }}</v-card-title>
      <template #append>
        <v-menu>
          <template #activator="{ props: menuProps }">
            <v-btn
              v-bind="menuProps"
              icon="mdi-dots-vertical"
              variant="text"
              size="small"
              @click.stop
            />
          </template>
          <v-list density="compact">
            <v-list-item
              v-if="session.status === 'in_progress'"
              prepend-icon="mdi-check-circle"
              title="Finalizar sesión"
              :disabled="finishMut.isPending.value"
              @click="finishDialog = true"
            />
            <v-list-item
              v-if="canDelete"
              prepend-icon="mdi-delete"
              title="Eliminar sesión"
              :disabled="deleteMut.isPending.value"
              @click="deleteDialog = true"
            />
          </v-list>
        </v-menu>
      </template>
      <v-card-subtitle class="d-flex flex-wrap ga-1">
        <v-chip
          :color="SESSION_STATUS_COLORS[session.status]"
          size="x-small"
          variant="flat"
          class="mr-1"
        >
          <v-icon start size="14">{{ STATUS_ICONS[session.status] }}</v-icon>
          {{ SESSION_STATUS_LABELS[session.status] }}
        </v-chip>
        <v-chip size="x-small" variant="tonal" class="mr-1">
          <v-icon start size="14">{{ MODE_ICONS[session.mode] }}</v-icon>
          {{ SESSION_MODE_LABELS[session.mode] }}
        </v-chip>
        <v-tooltip location="bottom" :text="SESSION_DIFFICULTY_LABELS[session.difficulty]">
          <template #activator="{ props: tooltipProps }">
            <v-icon v-bind="tooltipProps" :color="DIFFICULTY_COLORS[session.difficulty]" size="16">
              {{ DIFFICULTY_ICONS[session.difficulty] }}
            </v-icon>
          </template>
        </v-tooltip>
        <v-chip v-if="isGenerating" size="x-small" variant="flat" color="warning" class="ml-1">
          <v-icon start size="14">mdi-cog</v-icon>
          Generando...
        </v-chip>
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      <div class="d-flex align-center ga-4">
        <div>
          <span class="text-caption text-medium-emphasis">Puntuación</span>
          <div class="text-body-2 font-weight-medium">
            {{ session.status === 'completed' ? formatScore(session.score) : '—' }}
          </div>
        </div>
        <div>
          <span class="text-caption text-medium-emphasis">Respuestas</span>
          <div class="text-body-2 font-weight-medium">{{ session.answerCount }}</div>
        </div>
        <div v-if="session.topics.length">
          <span class="text-caption text-medium-emphasis">Temas</span>
          <div class="text-body-2 d-flex flex-wrap ga-1 mt-1">
            <v-chip
              v-for="(topic, idx) in displayedTopics"
              :key="idx"
              size="x-small"
              variant="outlined"
            >
              {{ topic }}
            </v-chip>
            <v-tooltip v-if="hiddenCount > 0" location="bottom" :text="session.topics.join(', ')">
              <template #activator="{ props: tooltipProps }">
                <v-chip v-bind="tooltipProps" size="x-small" variant="tonal" color="grey">
                  +{{ hiddenCount }} más
                </v-chip>
              </template>
            </v-tooltip>
          </div>
        </div>
      </div>
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <span class="text-caption text-medium-emphasis">
        {{ formatDate(session.createdAt) }}
      </span>
    </v-card-actions>

    <v-dialog v-model="finishDialog" max-width="400">
      <v-card>
        <v-card-title>¿Finalizar sesión?</v-card-title>
        <v-card-text>
          Al finalizar se cerrará la sesión y podrás ver tu resumen de respuestas.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="finishDialog = false">Cancelar</v-btn>
          <v-btn color="primary" :loading="finishMut.isPending.value" @click="handleFinish">
            Finalizar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title>¿Eliminar sesión?</v-card-title>
        <v-card-text>
          Esta acción no se puede deshacer. La sesión se borrará permanentemente.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false">Cancelar</v-btn>
          <v-btn color="error" :loading="deleteMut.isPending.value" @click="handleDelete">
            Eliminar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<style scoped>
.session-card--blocked {
  opacity: 0.6;
  cursor: default;
  pointer-events: auto;
}

.session-card--blocked :deep(.v-card-item__prepend .v-icon) {
  color: rgb(var(--v-theme-warning)) !important;
}
</style>
