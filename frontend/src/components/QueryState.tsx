import type { ReactNode } from 'react'
import { apiErrorMessage } from '../lib/api'
import { Spinner } from './ui'

/** Renders loading/error/empty boilerplate around a React Query result. */
export function QueryState({
  isLoading,
  error,
  isEmpty,
  emptyText = 'Nenhum registro encontrado.',
  children,
}: {
  isLoading: boolean
  error: unknown
  isEmpty?: boolean
  emptyText?: string
  children: ReactNode
}) {
  if (isLoading) return <Spinner />
  if (error)
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
        {apiErrorMessage(error)}
      </div>
    )
  if (isEmpty)
    return <div className="py-8 text-center text-sm text-brand-600">{emptyText}</div>
  return <>{children}</>
}
