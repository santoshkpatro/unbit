import axios from 'axios'
import { push } from 'notivue'

export const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

apiClient.interceptors.response.use(
  function (response) {
    const payload = response.data

    if (payload && typeof payload === 'object' && 'success' in payload) {
      if (payload.success === true) {
        // if has message → show info
        if (payload.message) {
          push.success(payload.message)
        }
        return payload.data
      }

      // fail path
      const msg = payload.message || 'Request failed'
      push.error(msg)
      return Promise.reject(new Error(msg))
    }

    return response.data
  },
  function (error) {
    const msg = error?.response?.data?.message || error?.message || 'Network error'

    push.error(msg)
    return Promise.reject(error)
  },
)
