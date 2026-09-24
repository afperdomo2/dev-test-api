import { queryOptions } from '@tanstack/vue-query'
import * as tipsService from '@/api/services/tips.service'
import type { TipsFilters } from '@/api/services/tips.service'

export function tipsListOptions(
  page: () => number,
  perPage: () => number,
  filters: () => TipsFilters,
) {
  return queryOptions({
    queryKey: ['tips', 'list', page, perPage, filters],
    queryFn: () => tipsService.listTips(page(), perPage(), filters()),
    staleTime: 30 * 1000,
  })
}

export function tipDetailOptions(id: () => string) {
  return queryOptions({
    queryKey: ['tips', 'detail', id],
    queryFn: () => tipsService.getTipById(id()),
  })
}

export function createTipMutation() {
  return {
    mutationKey: ['tips', 'create'],
    mutationFn: tipsService.createTip,
  }
}

export function updateTipMutation() {
  return {
    mutationKey: ['tips', 'update'],
    mutationFn: ({ id, data }: { id: string; data: Parameters<typeof tipsService.updateTip>[1] }) =>
      tipsService.updateTip(id, data),
  }
}

export function deleteTipMutation() {
  return {
    mutationKey: ['tips', 'delete'],
    mutationFn: tipsService.deleteTip,
  }
}

export function importTipsMutation() {
  return {
    mutationKey: ['tips', 'import'],
    mutationFn: (formData: FormData) => tipsService.importTips(formData),
  }
}

export function importTipsTextMutation() {
  return {
    mutationKey: ['tips', 'import-text'],
    mutationFn: (content: string) => tipsService.importTipsFromText(content),
  }
}
