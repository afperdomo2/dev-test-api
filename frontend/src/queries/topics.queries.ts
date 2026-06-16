import { queryOptions } from '@tanstack/vue-query'
import * as topicsService from '@/api/services/topics.service'

export function topicsListOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['topics', 'list', page, perPage],
    queryFn: () => topicsService.listTopics(page(), perPage()),
    staleTime: 60 * 1000,
  })
}

export function topicDetailOptions(id: () => string) {
  return queryOptions({
    queryKey: ['topics', 'detail', id],
    queryFn: () => topicsService.getTopicById(id()),
  })
}

export function createTopicMutation() {
  return {
    mutationKey: ['topics', 'create'],
    mutationFn: topicsService.createTopic,
  }
}

export function updateTopicMutation() {
  return {
    mutationKey: ['topics', 'update'],
    mutationFn: ({ id, data }: { id: string; data: Parameters<typeof topicsService.updateTopic>[1] }) =>
      topicsService.updateTopic(id, data),
  }
}

export function deleteTopicMutation() {
  return {
    mutationKey: ['topics', 'delete'],
    mutationFn: topicsService.deleteTopic,
  }
}
