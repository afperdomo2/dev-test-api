import { queryOptions } from '@tanstack/vue-query'
import * as opencodeService from '@/api/services/opencode.service'

export function openCodeUsageOptions() {
  return queryOptions({
    queryKey: ['opencode', 'usage'],
    queryFn: () => opencodeService.getOpenCodeUsage(),
    staleTime: 60 * 1000,
  })
}
