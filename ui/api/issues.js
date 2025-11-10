import { apiClient } from '@/api/index'

export const recentIssuesAPI = (filters) =>
  apiClient.get('/issues/recent', {
    params: filters,
  })

export const issueDetailsAPI = (issueId) => apiClient.get(`/issues/${issueId}`)
