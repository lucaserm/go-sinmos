import { createContext } from 'react'
import type { User } from '../types/api'

export interface AuthContextValue {
  user: User | null
  loading: boolean
  login: (code: string, password: string) => Promise<void>
  logout: () => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)
