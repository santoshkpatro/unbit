import { apiClient } from '@/api/index'

export const eventListAPI = () => apiClient.get('/events')
