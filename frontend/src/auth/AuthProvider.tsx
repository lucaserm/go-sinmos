import { useCallback, useEffect, useState, type ReactNode } from 'react'
import { api, setAuthFailureHandler } from '../lib/api'
import { authStorage } from '../lib/auth-storage'
import { queryClient } from '../lib/queryClient'
import type { AuthResponse, User } from '../types/api'
import { AuthContext } from './AuthContext'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const logout = useCallback(() => {
    authStorage.clear()
    setUser(null)
    queryClient.clear()
  }, [])

  // Let the axios interceptor force a logout when refresh is impossible.
  useEffect(() => {
    setAuthFailureHandler(logout)
  }, [logout])

  // On boot, resolve the session from a stored token via /auth/me.
  useEffect(() => {
    if (!authStorage.getAccess()) {
      setLoading(false)
      return
    }
    api
      .get<{ user: User }>('/auth/me')
      .then((res) => setUser(res.data.user))
      .catch(() => authStorage.clear())
      .finally(() => setLoading(false))
  }, [])

  const login = useCallback(async (code: string, password: string) => {
    const { data } = await api.post<AuthResponse>('/auth/login', { code, password })
    authStorage.set(data.accessToken, data.refreshToken)
    setUser(data.user)
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
