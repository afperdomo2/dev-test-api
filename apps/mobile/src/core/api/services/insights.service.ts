import type { InsightPeriod, InsightSummary } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

/** Resumen de actividad del periodo (ventana móvil de 7 o 30 días). */
export async function getInsightSummary(period: InsightPeriod): Promise<InsightSummary> {
  return apiClient.get<InsightSummary>('/api/v1/insights/summary', { period })
}
