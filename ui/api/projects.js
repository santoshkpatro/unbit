import { apiClient } from '@/api/index'

export const projectListAPI = () => apiClient.get('/projects')
