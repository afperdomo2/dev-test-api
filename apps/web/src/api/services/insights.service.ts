import apiClient from '@/api/client'
import type { InsightPeriod, InsightSummary } from '@/types/insight.types'

/**
 * Resumen de actividad del periodo. El backend usa una ventana móvil
 * (7 o 30 días) y devuelve el informe completo en una sola llamada.
 */
export async function getInsightSummary(period: InsightPeriod): Promise<InsightSummary> {
  const res = await apiClient.get<InsightSummary>('/api/v1/insights/summary', {
    params: { period },
  })
  return res.data
}
