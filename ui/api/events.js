import { apiClient } from '@/api/index'

export const eventIssuesAPI = (filters) =>
  apiClient.get('/issues/recent', {
    params: filters,
  })
