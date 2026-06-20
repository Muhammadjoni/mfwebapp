import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, TokenPair } from '../types'

interface AdminState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  setUser: (user: User) => void
  setTokens: (access: string, refresh: string) => void
  logout: () => void
}

export const useAdminStore = create<AdminState>()(
  persist(
    set => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,

      setUser: user => set({ user }),

      setTokens: (accessToken, refreshToken) => {
        localStorage.setItem('mf-admin-access', accessToken)
        localStorage.setItem('mf-admin-refresh', refreshToken)
        set({ accessToken, refreshToken, isAuthenticated: true })
      },

      logout: () => {
        localStorage.removeItem('mf-admin-access')
        localStorage.removeItem('mf-admin-refresh')
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false })
      },
    }),
    { name: 'mf-admin', partialize: s => ({ user: s.user, accessToken: s.accessToken, refreshToken: s.refreshToken, isAuthenticated: s.isAuthenticated }) }
  )
)
