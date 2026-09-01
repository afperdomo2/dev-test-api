<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { questionStatsOptions } from '@/queries/questions.queries'
import { DIFFICULTY_COLORS, DIFFICULTY_LABELS, TYPE_ICONS } from '@/types/question.types'
import type { TopicCount, CategoryCount, DifficultyCount, TypeCount } from '@/types/question.types'

const router = useRouter()
const queryClient = useQueryClient()

const { data: stats, isLoading, isError } = useQuery(questionStatsOptions())

const total = computed(() => stats.value?.total ?? 0)
const byTopic = computed<Array<TopicCount>>(() => stats.value?.byTopic ?? [])
const byCategory = computed<Array<CategoryCount>>(() => stats.value?.byCategory ?? [])
const byDifficulty = computed<Array<DifficultyCount>>(() => stats.value?.byDifficulty ?? [])
const byType = computed<Array<TypeCount>>(() => stats.value?.byType ?? [])

const maxTopicCount = computed(() => Math.max(1, ...byTopic.value.map((t) => t.count)))
const maxCategoryCount = computed(() => Math.max(1, ...byCategory.value.map((c) => c.count)))
const maxDifficultyCount = computed(() => Math.max(1, ...byDifficulty.value.map((d) => d.count)))
const maxTypeCount = computed(() => Math.max(1, ...byType.value.map((t) => t.count)))

function handleRefresh() {
  queryClient.invalidateQueries({ queryKey: ['questions', 'stats'] })
}

function goToQuestionsByTopic(topicId: string) {
  router.push({ path: '/questions', query: { topicIds: topicId } })
}

function onTopicRowClick(_e: unknown, row: { item: TopicCount }) {
  goToQuestionsByTopic(row.item.topicId)
}

const typeLabels: Record<string, string> = {
  single_choice: 'Selección única',
  multiple_choice: 'Selección múltiple',
  code_completion: 'Completar código',
}

const topicHeaders = [
  { title: 'Tema', key: 'name', sortable: false },
  { title: 'Categoría', key: 'category', sortable: false, align: 'center' as const },
  { title: 'Cantidad', key: 'count', sortable: false, align: 'center' as const },
  { title: '', key: 'action', sortable: false, align: 'center' as const },
]

const categoryHeaders = [
  { title: 'Categoría', key: 'category', sortable: false },
  { title: 'Cantidad', key: 'count', sortable: false, align: 'center' as const },
]
</script>

<template>
  <v-container>
    <div class="d-flex align-center justify-space-between mb-4">
      <h1 class="text-h4">Estadísticas</h1>
      <v-btn variant="tonal" color="secondary" prepend-icon="mdi-refresh" @click="handleRefresh">
        Refrescar
      </v-btn>
    </div>

    <v-alert v-if="isError" type="error" variant="tonal" class="mb-4">
      No se pudieron cargar las estadísticas.
    </v-alert>

    <v-skeleton-loader v-if="isLoading" type="card, card, table" />

    <template v-else>
      <!-- Total + summary chips -->
      <v-row class="mb-4">
        <v-col cols="12" md="4">
          <v-card color="primary" variant="tonal">
            <v-card-text class="d-flex align-center ga-4">
              <v-icon size="48" color="primary">mdi-help-circle</v-icon>
              <div>
                <div class="text-h3 font-weight-bold">{{ total }}</div>
                <div class="text-body-2 text-medium-emphasis">Preguntas visibles</div>
                <div class="text-caption text-medium-emphasis">Tus preguntas + IA compartidas</div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" md="8">
          <v-card variant="flat" border>
            <v-card-text>
              <div class="text-subtitle-2 mb-2 d-flex align-center ga-2">
                <v-icon size="18">mdi-chart-donut</v-icon>
                Resumen
              </div>
              <div class="d-flex flex-wrap ga-2">
                <v-chip
                  v-for="d in byDifficulty"
                  :key="d.difficulty"
                  :color="
                    DIFFICULTY_COLORS[d.difficulty as keyof typeof DIFFICULTY_COLORS] ?? 'default'
                  "
                  variant="tonal"
                  size="small"
                >
                  {{
                    DIFFICULTY_LABELS[d.difficulty as keyof typeof DIFFICULTY_LABELS] ??
                    d.difficulty
                  }}: {{ d.count }}
                </v-chip>
                <v-chip
                  v-for="t in byType"
                  :key="t.type"
                  variant="tonal"
                  size="small"
                  :prepend-icon="TYPE_ICONS[t.type as keyof typeof TYPE_ICONS] ?? 'mdi-help'"
                >
                  {{ typeLabels[t.type] ?? t.type }}: {{ t.count }}
                </v-chip>
              </div>
              <div
                v-if="!byDifficulty.length && !byType.length"
                class="text-body-2 text-medium-emphasis"
              >
                Aún no hay preguntas para mostrar distribución.
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Por tema -->
      <v-card class="mb-4" variant="flat" border>
        <v-card-title class="d-flex align-center ga-2">
          <v-icon>mdi-tag-multiple</v-icon>
          Preguntas por tema
          <v-chip size="small" variant="tonal" class="ml-2">{{ byTopic.length }} temas</v-chip>
        </v-card-title>
        <v-card-subtitle>Clic en una fila para ver las preguntas de ese tema</v-card-subtitle>

        <v-card-text>
          <!-- Barras horizontales -->
          <div v-if="byTopic.length" class="mb-6">
            <div
              v-for="item in byTopic.slice(0, 20)"
              :key="item.topicId"
              class="d-flex align-center ga-3 mb-2"
              style="cursor: pointer"
              @click="goToQuestionsByTopic(item.topicId)"
            >
              <div class="text-body-2 text-no-wrap" style="min-width: 140px; max-width: 180px">
                <span class="font-weight-medium">{{ item.name }}</span>
                <span class="text-caption text-medium-emphasis ml-1">({{ item.category }})</span>
              </div>
              <v-progress-linear
                :model-value="(item.count / maxTopicCount) * 100"
                color="primary"
                height="22"
                rounded
              >
                <span class="text-caption font-weight-bold">{{ item.count }}</span>
              </v-progress-linear>
            </div>
            <div
              v-if="byTopic.length > 20"
              class="text-caption text-medium-emphasis text-center mt-2"
            >
              Mostrando 20 de {{ byTopic.length }} temas — ver tabla completa abajo.
            </div>
          </div>
          <div v-else class="text-body-2 text-medium-emphasis text-center py-4">
            No hay preguntas agrupadas por tema.
          </div>

          <v-data-table
            :headers="topicHeaders"
            :items="byTopic"
            :items-per-page="10"
            no-data-text="Sin datos por tema"
            density="compact"
            class="mt-2"
            @click:row="onTopicRowClick"
          >
            <template #[`item.category`]="{ item }">
              <v-chip size="x-small" variant="tonal">{{ item.category }}</v-chip>
            </template>
            <template #[`item.count`]="{ item }">
              <v-chip size="small" color="primary" variant="tonal">{{ item.count }}</v-chip>
            </template>
            <template #[`item.action`]="{ item }">
              <v-btn
                size="x-small"
                variant="text"
                color="primary"
                prepend-icon="mdi-open-in-new"
                @click.stop="goToQuestionsByTopic(item.topicId)"
              >
                Ver
              </v-btn>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>

      <!-- Por categoría / dificultad / tipo -->
      <v-row>
        <v-col cols="12" md="6">
          <v-card variant="flat" border class="h-100">
            <v-card-title class="d-flex align-center ga-2">
              <v-icon>mdi-shape</v-icon>
              Por categoría
            </v-card-title>
            <v-card-text>
              <div v-if="byCategory.length" class="mb-4">
                <div
                  v-for="c in byCategory"
                  :key="c.category"
                  class="d-flex align-center ga-3 mb-2"
                >
                  <div class="text-body-2 font-weight-medium" style="min-width: 110px">
                    {{ c.category }}
                  </div>
                  <v-progress-linear
                    :model-value="(c.count / maxCategoryCount) * 100"
                    color="#B2DFDB"
                    bg-color="#ECEFF1"
                    bg-opacity="1"
                    height="20"
                    rounded
                  >
                    <span class="text-caption font-weight-bold" style="color: #004d40">{{
                      c.count
                    }}</span>
                  </v-progress-linear>
                </div>
              </div>
              <v-data-table
                :headers="categoryHeaders"
                :items="byCategory"
                :items-per-page="10"
                no-data-text="Sin datos por categoría"
                density="compact"
                hide-default-footer
              >
                <template #[`item.count`]="{ item }">
                  <v-chip size="small" color="teal" variant="tonal">{{ item.count }}</v-chip>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" md="6">
          <v-card variant="flat" border class="mb-4">
            <v-card-title class="d-flex align-center ga-2">
              <v-icon>mdi-signal-cellular-3</v-icon>
              Por dificultad
            </v-card-title>
            <v-card-text>
              <div v-if="byDifficulty.length">
                <div
                  v-for="d in byDifficulty"
                  :key="d.difficulty"
                  class="d-flex align-center ga-3 mb-2"
                >
                  <v-chip
                    :color="
                      DIFFICULTY_COLORS[d.difficulty as keyof typeof DIFFICULTY_COLORS] ?? 'default'
                    "
                    size="small"
                    variant="tonal"
                    style="min-width: 110px; justify-content: center"
                  >
                    {{
                      DIFFICULTY_LABELS[d.difficulty as keyof typeof DIFFICULTY_LABELS] ??
                      d.difficulty
                    }}
                  </v-chip>
                  <v-progress-linear
                    :model-value="(d.count / maxDifficultyCount) * 100"
                    :color="
                      DIFFICULTY_COLORS[d.difficulty as keyof typeof DIFFICULTY_COLORS] ?? 'primary'
                    "
                    height="20"
                    rounded
                  >
                    <span class="text-caption font-weight-bold">{{ d.count }}</span>
                  </v-progress-linear>
                </div>
              </div>
              <div v-else class="text-body-2 text-medium-emphasis">Sin datos por dificultad.</div>
            </v-card-text>
          </v-card>

          <v-card variant="flat" border>
            <v-card-title class="d-flex align-center ga-2">
              <v-icon>mdi-format-list-bulleted-type</v-icon>
              Por tipo
            </v-card-title>
            <v-card-text>
              <div v-if="byType.length">
                <div v-for="t in byType" :key="t.type" class="d-flex align-center ga-3 mb-2">
                  <div class="d-flex align-center ga-2" style="min-width: 160px">
                    <v-icon size="18">{{
                      TYPE_ICONS[t.type as keyof typeof TYPE_ICONS] ?? 'mdi-help'
                    }}</v-icon>
                    <span class="text-body-2">{{ typeLabels[t.type] ?? t.type }}</span>
                  </div>
                  <v-progress-linear
                    :model-value="(t.count / maxTypeCount) * 100"
                    color="info"
                    height="20"
                    rounded
                  >
                    <span class="text-caption font-weight-bold">{{ t.count }}</span>
                  </v-progress-linear>
                </div>
              </div>
              <div v-else class="text-body-2 text-medium-emphasis">Sin datos por tipo.</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>
