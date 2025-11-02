import { apiClient } from '@/api/index'

export const settingMetaAPI = () => apiClient.get('/setting/meta')
