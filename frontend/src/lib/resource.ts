import { api } from './api'

export interface ListParams {
  offset?: number
  limit?: number
}

/**
 * Factory for the API's uniform CRUD resources. Responses are single-key
 * envelope-wrapped ({"student": {...}} / {"students": [...]}); this unwraps them.
 * Note `singular`/`plural` are the envelope keys, which are camelCase for some
 * resources (e.g. occurrenceType / occurrenceTypes).
 */
export function createResource<T>(path: string, singular: string, plural: string) {
  return {
    path,
    list: async (params?: ListParams): Promise<T[]> => {
      const res = await api.get(path, { params })
      return res.data[plural] as T[]
    },
    get: async (id: string): Promise<T> => {
      const res = await api.get(`${path}/${id}`)
      return res.data[singular] as T
    },
    create: async (body: unknown): Promise<T> => {
      const res = await api.post(path, body)
      return res.data[singular] as T
    },
    update: async (id: string, body: unknown): Promise<T> => {
      const res = await api.put(`${path}/${id}`, body)
      return res.data[singular] as T
    },
    remove: async (id: string): Promise<void> => {
      await api.delete(`${path}/${id}`)
    },
  }
}
