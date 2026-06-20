import axios, { type InternalAxiosRequestConfig, type AxiosError } from 'axios'
import { useAuthStore } from '../store/slices/authStore'

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? 'http://localhost:8080',
  timeout: 15_000,
  headers: { 'Content-Type': 'application/json' },
})

apiClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = useAuthStore.getState().accessToken
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let isRefreshing = false
let queue: Array<{ resolve: (t: string) => void; reject: (e: unknown) => void }> = []

const flush = (error: unknown, token: string | null) => {
  queue.forEach(p => (error ? p.reject(error) : p.resolve(token!)))
  queue = []
}

apiClient.interceptors.response.use(
  r => r,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean }
    if (error.response?.status !== 401 || original._retry) return Promise.reject(error)

    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        queue.push({
          resolve: token => { original.headers.Authorization = `Bearer ${token}`; resolve(apiClient(original)) },
          reject,
        })
      })
    }

    original._retry = true
    isRefreshing = true

    try {
      const rt = useAuthStore.getState().refreshToken
      const { data } = await axios.post(`${apiClient.defaults.baseURL}/api/v1/auth/refresh`, { refresh_token: rt })
      const { access_token, refresh_token } = data.data
      useAuthStore.getState().setTokens(access_token, refresh_token)
      flush(null, access_token)
      original.headers.Authorization = `Bearer ${access_token}`
      return apiClient(original)
    } catch (err) {
      flush(err, null)
      useAuthStore.getState().logout()
      return Promise.reject(err)
    } finally {
      isRefreshing = false
    }
  }
)
