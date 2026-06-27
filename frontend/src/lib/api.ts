import axios, {
  AxiosError,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from 'axios'
import { authStorage } from './auth-storage'

const BASE = import.meta.env.VITE_API_BASE ?? '/api/v1'

export const api = axios.create({ baseURL: BASE })

// Attach the bearer token to every request.
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = authStorage.getAccess()
  if (token) config.headers.set('Authorization', `Bearer ${token}`)
  return config
})

// Called when refresh is impossible/failed — wired up by AuthProvider so it can
// clear React state and redirect, not just nuke localStorage.
let onAuthFailure: () => void = () => authStorage.clear()
export function setAuthFailureHandler(fn: () => void) {
  onAuthFailure = fn
}

type RetriableConfig = AxiosRequestConfig & { _retried?: boolean }

// Single in-flight refresh shared by concurrent 401s.
let refreshing: Promise<string | null> | null = null

async function refreshAccessToken(): Promise<string | null> {
  if (!authStorage.hasValidRefresh()) return null
  try {
    // Bare axios (not `api`) to skip the interceptors and avoid a loop.
    const { data } = await axios.post(`${BASE}/auth/refresh`, {
      refreshToken: authStorage.getRefresh(),
    })
    authStorage.set(data.accessToken, data.refreshToken)
    return data.accessToken as string
  } catch {
    return null
  }
}

api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as RetriableConfig | undefined
    const isAuthCall = original?.url?.includes('/auth/')

    if (error.response?.status === 401 && original && !original._retried && !isAuthCall) {
      original._retried = true
      refreshing ??= refreshAccessToken().finally(() => {
        refreshing = null
      })
      const newToken = await refreshing
      if (newToken) {
        original.headers = { ...original.headers, Authorization: `Bearer ${newToken}` }
        return api(original)
      }
      onAuthFailure()
    }
    return Promise.reject(error)
  },
)

/** Normalizes the API's `{error, timestamp}` envelope into a message string. */
export function apiErrorMessage(err: unknown, fallback = 'Algo deu errado.'): string {
  if (err instanceof AxiosError) {
    const data = err.response?.data as { error?: string } | undefined
    if (data?.error) return data.error
    if (err.message) return err.message
  }
  return fallback
}
