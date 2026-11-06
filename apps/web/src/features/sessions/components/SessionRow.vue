<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import type { Session } from '@/types/session.types'
import {
  SESSION_STATUS_LABELS,
  SESSION_STATUS_COLORS,
  SESSION_MODE_LABELS,
  SESSION_MODE_ICONS,
  SESSION_DIFFICULTY_LABELS,
  SESSION_DIFFICULTY_ICONS,
  SESSION_DIFFICULTY_COLORS,
} from '@/types/session.types'
import { formatRelativeTime, formatDateTime, formatScore, formatDuration } from '@/utils/format'

interface Props {
  session: Session
  busy?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  busy: false,
})

const emit = defineEmits<{
  open: [session: Session]
  finish: [session: Session]
  delete: [session: Session]
}>()

const STATUS_ICONS: Record<string, string> = {
  in_progress: 'mdi-progress-clock',
  completed: 'mdi-check-circle',
  cancelled: 'mdi-cancel',
}

const isInProgress = computed(() => props.session.status === 'in_progress')
const isCompleted = computed(() => props.session.status === 'completed')

const statusColor = computed(() => SESSION_STATUS_COLORS[props.session.status])
const statusIcon = computed(() => STATUS_ICONS[props.session.status] ?? 'mdi-help-circle')
const modeIcon = computed(() => SESSION_MODE_ICONS[props.session.mode])

const targetRoute = computed(() =>
  isCompleted.value
    ? `/sessions/${props.session.id}/review`
    : `/sessions/${props.session.id}/study`,
)

const ctaLabel = computed(() => {
  if (isInProgress.value) return 'Continuar'
  if (isCompleted.value) return 'Ver resumen'
  return 'Ver detalle'
})

const canDelete = computed(() => {
  if (props.session.answerCount > 0) return false
  const created = new Date(props.session.createdAt).getTime()
  return created > Date.now() - 24 * 60 * 60 * 1000
})

const hasActions = computed(() => isInProgress.value || canDelete.value)

const limit = computed(() => props.session.questionLimit)
const progressPercent = computed(() => {
  const total = limit.value
  if (!total || total <= 0) return 0
  return Math.min(100, Math.round((props.session.answerCount / total) * 100))
})

const progressLabel = computed(() => {
  if (isInProgress.value && limit.value) {
    return `${props.session.answerCount}/${limit.value}`
  }
  return `${props.session.answerCount} respuestas`
})

const duration = computed(() => formatDuration(props.session.startedAt, props.session.finishedAt))

const showsProgressBar = computed(() => isInProgress.value && !!limit.value)

const secondaryLine = computed(() => {
  if (isCompleted.value) {
    const answered = `${props.session.answerCount} respuestas`
    return duration.value ? `${answered} · ${duration.value}` : answered
  }
  if (isInProgress.value && !limit.value) return 'Sin límite de preguntas'
  return duration.value
})

const displayedTopics = computed(() => props.session.topics.slice(0, 2))
const hiddenTopics = computed(() => Math.max(0, props.session.topics.length - 2))

const rowAriaLabel = computed(
  () =>
    `Sesión ${props.session.name}, ${SESSION_STATUS_LABELS[props.session.status]}, ` +
    `${props.session.answerCount} respuestas`,
)
</script>

<template>
  <div class="session-row" role="group" :aria-label="rowAriaLabel" @click="emit('open', session)">
    <div class="session-row__main">
      <RouterLink
        :to="targetRoute"
        class="session-row__name text-body-1 font-weight-medium"
        @click.stop
      >
        {{ session.name }}
      </RouterLink>

      <div class="d-flex flex-wrap align-center ga-2 mt-1">
        <v-chip :color="statusColor" size="x-small" variant="flat">
          <v-icon start size="14">{{ statusIcon }}</v-icon>
          {{ SESSION_STATUS_LABELS[session.status] }}
        </v-chip>

        <v-tooltip location="bottom" :text="formatDateTime(session.createdAt)">
          <template #activator="{ props: tooltipProps }">
            <span v-bind="tooltipProps" class="text-caption text-medium-emphasis">
              {{ formatRelativeTime(session.createdAt) }}
            </span>
          </template>
        </v-tooltip>

        <span class="text-caption text-medium-emphasis d-inline-flex align-center">
          <v-icon size="14" class="mr-1">{{ modeIcon }}</v-icon>
          {{ SESSION_MODE_LABELS[session.mode] }}
        </span>

        <span
          class="text-caption d-inline-flex align-center"
          :class="`text-${SESSION_DIFFICULTY_COLORS[session.difficulty]}`"
        >
          <v-icon size="14" class="mr-1">{{ SESSION_DIFFICULTY_ICONS[session.difficulty] }}</v-icon>
          {{ SESSION_DIFFICULTY_LABELS[session.difficulty] }}
        </span>

        <span class="text-caption text-medium-emphasis session-row__progress-inline">
          · {{ progressLabel }}
        </span>
      </div>
    </div>

    <div class="session-row__progress">
      <div class="text-body-2 font-weight-medium">
        <template v-if="isCompleted">{{ formatScore(session.score) }}</template>
        <template v-else>{{ progressLabel }}</template>
      </div>
      <v-progress-linear
        v-if="showsProgressBar"
        :model-value="progressPercent"
        :color="statusColor"
        height="6"
        rounded
        class="mt-1"
      />
      <div v-if="secondaryLine" class="text-caption text-medium-emphasis mt-1">
        {{ secondaryLine }}
      </div>
    </div>

    <div class="session-row__topics">
      <div v-if="session.topics.length" class="d-flex flex-wrap ga-1">
        <v-chip
          v-for="(topic, idx) in displayedTopics"
          :key="`${topic}-${idx}`"
          size="x-small"
          variant="outlined"
        >
          {{ topic }}
        </v-chip>
        <v-tooltip v-if="hiddenTopics > 0" location="bottom" :text="session.topics.join(', ')">
          <template #activator="{ props: tooltipProps }">
            <v-chip v-bind="tooltipProps" size="x-small" variant="tonal" color="grey">
              +{{ hiddenTopics }}
            </v-chip>
          </template>
        </v-tooltip>
      </div>
      <span v-else class="text-caption text-medium-emphasis">Sin temas</span>
    </div>

    <div class="session-row__actions">
      <v-btn
        :to="targetRoute"
        :variant="isInProgress ? 'flat' : 'tonal'"
        :color="isInProgress ? 'primary' : undefined"
        size="small"
        class="text-none"
        :disabled="busy"
        @click.stop
      >
        {{ ctaLabel }}
      </v-btn>

      <v-menu v-if="hasActions" location="bottom end">
        <template #activator="{ props: menuProps }">
          <v-btn
            v-bind="menuProps"
            icon="mdi-dots-vertical"
            variant="text"
            size="small"
            :disabled="busy"
            :aria-label="`Acciones de ${session.name}`"
            @click.stop
          />
        </template>
        <v-list density="compact">
          <v-list-item
            v-if="isInProgress"
            prepend-icon="mdi-check-circle"
            title="Finalizar sesión"
            @click="emit('finish', session)"
          />
          <v-list-item
            v-if="canDelete"
            prepend-icon="mdi-delete"
            title="Eliminar sesión"
            @click="emit('delete', session)"
          />
        </v-list>
      </v-menu>
    </div>
  </div>
</template>

<style scoped>
.session-row {
  display: grid;
  grid-template-columns: minmax(0, 2.2fr) minmax(150px, 1fr) minmax(0, 1.2fr) auto;
  align-items: center;
  gap: 16px;
  min-height: 68px;
  padding: 10px 12px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.session-row:last-child {
  border-bottom: none;
}

.session-row:hover {
  background-color: rgba(var(--v-theme-on-surface), 0.04);
}

.session-row__main {
  min-width: 0;
}

.session-row__name {
  display: block;
  color: inherit;
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-row__name:hover {
  text-decoration: underline;
}

.session-row__main,
.session-row__progress,
.session-row__topics {
  min-width: 0;
}

.session-row__actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.session-row__progress-inline {
  display: none;
}

@media (max-width: 1279px) {
  .session-row {
    grid-template-columns: minmax(0, 2fr) minmax(130px, 1fr) auto;
  }

  .session-row__topics {
    display: none;
  }
}

@media (max-width: 959px) {
  .session-row {
    grid-template-columns: minmax(0, 1fr) auto;
    gap: 12px;
  }

  .session-row__progress {
    display: none;
  }

  .session-row__progress-inline {
    display: inline;
  }
}
</style>
