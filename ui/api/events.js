import { apiClient } from '@/api/index'

export const eventIssuesAPI = () => apiClient.get('/events/issues')
