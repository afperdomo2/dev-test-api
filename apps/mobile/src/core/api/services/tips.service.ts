import type { Tip } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

/**
 * Tip del día. El backend lo elige de forma determinista por fecha (UTC), así
 * que es estable durante el día y rota al siguiente. 404 si no hay tips.
 */
export async function getTipOfTheDay(): Promise<Tip> {
  return apiClient.get<Tip>('/api/v1/tips/today')
}
