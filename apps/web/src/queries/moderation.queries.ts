import { queryOptions } from '@tanstack/vue-query'
import * as moderationService from '@/api/services/moderation.service'
import type { ReportStatus } from '@/types/moderation.types'

export function reportsListOptions(
  page: () => number,
  perPage: () => number,
  status: () => ReportStatus | undefined,
) {
  return queryOptions({
    queryKey: ['reports', 'list', page, perPage, status],
    queryFn: () => moderationService.listReports(page(), perPage(), status()),
    staleTime: 30 * 1000,
  })
}

export function createReportMutation() {
  return {
    mutationKey: ['reports', 'create'],
    mutationFn: moderationService.createReport,
  }
}

export function reviewReportMutation() {
  return {
    mutationKey: ['reports', 'review'],
    mutationFn: ({
      id,
      data,
    }: {
      id: string
      data: Parameters<typeof moderationService.reviewReport>[1]
    }) => moderationService.reviewReport(id, data),
  }
}
