import { apiClient } from './client'
import type { User, TokenPair } from '../types'

interface AuthResponse { user: User; tokens: TokenPair }

export const authApi = {
  register: (payload: { email: string; password: string; first_name: string; last_name?: string; phone?: string; language?: 'ru' | 'tj' }) =>
    apiClient.post<{ data: AuthResponse }>('/api/v1/auth/register', payload).then(r => r.data.data),

  login: (payload: { email: string; password: string }) =>
    apiClient.post<{ data: AuthResponse }>('/api/v1/auth/login', payload).then(r => r.data.data),

  refresh: (refreshToken: string) =>
    apiClient.post<{ data: TokenPair }>('/api/v1/auth/refresh', { refresh_token: refreshToken }).then(r => r.data.data),

  logout: (refreshToken: string) =>
    apiClient.post('/api/v1/auth/logout', { refresh_token: refreshToken }),
}
