import { createContext, useContext, useEffect, useMemo, useState } from 'react'

import { apiClient } from '../api/client'

export type Role = 'owner' | 'admin' | 'member' | 'viewer'

export type User = {
  id: number
  name: string
  email: string
  role: Role
}

type LoginResponse = {
  token: string
  user: User
}

type MeResponse = {
  user: User
}

type AuthContextValue = {
  token: string | null
  user: User | null
  isReady: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

const tokenKey = 'gms-token'

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(tokenKey))
  const [user, setUser] = useState<User | null>(null)
  const [isReady, setIsReady] = useState(false)

  useEffect(() => {
    let mounted = true

    async function loadMe() {
      if (!token) {
        if (mounted) {
          setUser(null)
          setIsReady(true)
        }
        return
      }

      try {
        const data = await apiClient.get<MeResponse>('/api/auth/me', token)
        if (mounted) {
          setUser(data.user)
        }
      } catch {
        if (mounted) {
          setToken(null)
          setUser(null)
          localStorage.removeItem(tokenKey)
        }
      } finally {
        if (mounted) {
          setIsReady(true)
        }
      }
    }

    setIsReady(false)
    void loadMe()

    return () => {
      mounted = false
    }
  }, [token])

  async function login(email: string, password: string) {
    const payload = await apiClient.post<LoginResponse>('/api/auth/login', { email, password })
    localStorage.setItem(tokenKey, payload.token)
    setToken(payload.token)
    setUser(payload.user)
  }

  function logout() {
    localStorage.removeItem(tokenKey)
    setToken(null)
    setUser(null)
  }

  const value = useMemo(
    () => ({ token, user, isReady, login, logout }),
    [token, user, isReady],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
