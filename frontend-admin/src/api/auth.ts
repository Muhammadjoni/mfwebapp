import { apiClient } from './client'
import type { User, TokenPair } from '../types'

interface LoginResponse { user: User; tokens: TokenPair }

export const authApi = {
  login: (payload: { email: string; password: string }) =>
    apiClient.post<{ data: LoginResponse }>('/auth/login', payload).then(r => r.data.data),
  logout: () =>
    apiClient.post('/auth/logout'),
}
