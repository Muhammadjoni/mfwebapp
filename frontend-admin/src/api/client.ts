import axios from 'axios'

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
})

let isRefreshing = false
let queue: Array<(token: string) => void> = []

apiClient.interceptors.request.use(cfg => {
  const token = localStorage.getItem('mf-admin-access')
  if (token && cfg.headers) cfg.headers['Authorization'] = `Bearer ${token}`
  return cfg
})

apiClient.interceptors.response.use(
  r => r,
  async err => {
    const orig = err.config
    if (err.response?.status !== 401 || orig._retry) return Promise.reject(err)
    if (isRefreshing) {
      return new Promise(resolve => queue.push(token => { orig.headers['Authorization'] = `Bearer ${token}`; resolve(apiClient(orig)) }))
    }
    orig._retry = true
    isRefreshing = true
    try {
      const refresh = localStorage.getItem('mf-admin-refresh')
      const { data } = await axios.post(`${apiClient.defaults.baseURL}/auth/refresh`, { refresh_token: refresh })
      const newToken = data.data.access_token
      localStorage.setItem('mf-admin-access', newToken)
      localStorage.setItem('mf-admin-refresh', data.data.refresh_token)
      queue.forEach(cb => cb(newToken))
      queue = []
      orig.headers['Authorization'] = `Bearer ${newToken}`
      return apiClient(orig)
    } catch {
      localStorage.removeItem('mf-admin-access')
      localStorage.removeItem('mf-admin-refresh')
      window.location.href = '/login'
      return Promise.reject(err)
    } finally {
      isRefreshing = false
    }
  }
)
