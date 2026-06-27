// Date helpers for the pt-BR UI. The API returns occurredAt as RFC3339.
const dateTimeFmt = new Intl.DateTimeFormat('pt-BR', {
  dateStyle: 'short',
  timeStyle: 'short',
})

export function formatDateTime(iso?: string): string {
  if (!iso) return '—'
  const d = new Date(iso)
  return Number.isNaN(d.getTime()) ? iso : dateTimeFmt.format(d)
}

/** Converts a datetime-local input value to the RFC3339 the API expects. */
export function toRFC3339(localValue: string): string {
  if (!localValue) return ''
  const d = new Date(localValue)
  return Number.isNaN(d.getTime()) ? localValue : d.toISOString()
}
