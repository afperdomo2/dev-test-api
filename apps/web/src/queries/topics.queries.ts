import { queryOptions, infiniteQueryOptions } from '@tanstack/vue-query'
import * as topicsService from '@/api/services/topics.service'
import type { PaginatedResponse } from '@/types/api.types'
import type { Topic } from '@/types/topic.types'

const TOPICS_PER_PAGE = 20

export function topicsInfiniteOptions(search: () => string, enabled: () => boolean) {
  return infiniteQueryOptions({
    queryKey: ['topics', 'list', 'infinite', search],
    queryFn: ({ pageParam = 1 }) =>
      topicsService.listTopics(pageParam as number, TOPICS_PER_PAGE, 'name', 'asc', search()),
    getNextPageParam: (
      lastPage: PaginatedResponse<Topic>,
      allPages: Array<PaginatedResponse<Topic>>,
    ) => {
      const totalFetched = allPages.reduce((sum, p) => sum + p.data.length, 0)
      return totalFetched < lastPage.meta.total ? allPages.length + 1 : undefined
    },
    initialPageParam: 1,
    staleTime: 60 * 1000,
    enabled,
  })
}

export function topicsListOptions(
  page: () => number,
  perPage: () => number,
  sortBy?: () => string,
  sortOrder?: () => string,
  search?: () => string,
  category?: () => string,
  myOnly?: () => boolean,
) {
  return queryOptions({
    queryKey: ['topics', 'list', page, perPage, sortBy, sortOrder, search, category, myOnly],
    queryFn: () =>
      topicsService.listTopics(
        page(),
        perPage(),
        sortBy?.(),
        sortOrder?.(),
        search?.(),
        category?.(),
        myOnly?.(),
      ),
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
    mutationFn: ({
      id,
      data,
    }: {
      id: string
      data: Parameters<typeof topicsService.updateTopic>[1]
    }) => topicsService.updateTopic(id, data),
  }
}

export function deleteTopicMutation() {
  return {
    mutationKey: ['topics', 'delete'],
    mutationFn: topicsService.deleteTopic,
  }
}
