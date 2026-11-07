import { queryOptions } from '@tanstack/vue-query'
import * as generationService from '@/api/services/generation.service'
import type {
  CreateGenerationJobRequest,
  GenerationJobStatus,
  GenerationKind,
  UpdateSettingsRequest,
} from '@/types/generation.types'

/** Encolar un trabajo de generación. */
export function createGenerationJobMutation() {
  return {
    mutationKey: ['generation', 'create'],
    mutationFn: (data: CreateGenerationJobRequest) => generationService.createGenerationJob(data),
  }
}

/** Historial paginado. Se refresca solo mientras haya trabajo en curso. */
export function generationJobsListOptions(
  page: () => number,
  perPage: () => number,
  status: () => GenerationJobStatus | undefined,
  kind: () => GenerationKind,
) {
  return queryOptions({
    queryKey: ['generation', 'list', page, perPage, status, kind],
    queryFn: () => generationService.listGenerationJobs(page(), perPage(), status(), kind()),
    staleTime: 5 * 1000,
  })
}

export function generationStatsOptions(kind: () => GenerationKind | undefined = () => undefined) {
  return queryOptions({
    queryKey: ['generation', 'stats', kind],
    queryFn: () => generationService.getGenerationStats(kind()),
    staleTime: 5 * 1000,
  })
}

export function settingsListOptions() {
  return queryOptions({
    queryKey: ['settings', 'list'],
    queryFn: () => generationService.listSettings(),
    staleTime: 30 * 1000,
  })
}

export function updateSettingsMutation() {
  return {
    mutationKey: ['settings', 'update'],
    mutationFn: (data: UpdateSettingsRequest) => generationService.updateSettings(data),
  }
}
