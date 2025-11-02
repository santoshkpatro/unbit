import { apiClient } from '@/api/index'

export const authStatusAPI = () => apiClient.get('/auth/status')

export const loginAPI = (credentials) => apiClient.post('/auth/login', credentials)
