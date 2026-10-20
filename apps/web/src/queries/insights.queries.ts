import { queryOptions } from '@tanstack/vue-query'
import * as insightsService from '@/api/services/insights.service'
import type { InsightPeriod } from '@/types/insight.types'

export function insightSummaryOptions(period: () => InsightPeriod) {
  return queryOptions({
    queryKey: ['insights', 'summary', period],
    queryFn: () => insightsService.getInsightSummary(period()),
    // El informe cambia con cada sesión, pero no hace falta refrescarlo al instante.
    staleTime: 60 * 1000,
  })
}
