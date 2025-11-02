import { apiClient } from '@/api/index'

export const issueListAPI = () => apiClient.get('/issues')
