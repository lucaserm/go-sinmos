import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { QueryState } from '../../components/QueryState'
import { Table, Td, Th } from '../../components/Table'
import { Badge, Button, PageTitle } from '../../components/ui'
import { formatDateTime } from '../../lib/format'
import { occurrencesApi, studentsApi } from '../../lib/resources'
import { STATUS_LABELS, type OccurrenceStatus } from '../../types/api'
import { NewOccurrenceForm } from './NewOccurrenceForm'

const FILTERS: { value: 'ALL' | OccurrenceStatus; label: string }[] = [
  { value: 'ALL', label: 'Todas' },
  { value: 'PENDING', label: 'Pendentes' },
  { value: 'APPROVED', label: 'Aprovadas' },
  { value: 'REPROVED', label: 'Reprovadas' },
]

export function OccurrencesList() {
  const [filter, setFilter] = useState<'ALL' | OccurrenceStatus>('ALL')
  const [showForm, setShowForm] = useState(false)

  const { data, isLoading, error } = useQuery({
    queryKey: ['occurrences', 'all'],
    queryFn: () => occurrencesApi.list({ limit: 200 }),
  })

  const students = useQuery({
    queryKey: ['students', 'all'],
    queryFn: () => studentsApi.list({ limit: 200 }),
  })

  const studentName = useMemo(() => {
    const map = new Map(students.data?.map((s) => [s.id, s.name]))
    return (id: string) => map.get(id) ?? '—'
  }, [students.data])

  const rows = (data ?? []).filter((o) => filter === 'ALL' || o.status === filter)

  return (
    <div>
      <PageTitle
        title="Ocorrências"
        subtitle="Registre e acompanhe ocorrências disciplinares."
        action={<Button onClick={() => setShowForm(true)}>Nova Ocorrência</Button>}
      />

      <div className="mb-4 flex gap-2">
        {FILTERS.map((f) => (
          <button
            key={f.value}
            onClick={() => setFilter(f.value)}
            className={`rounded-full px-3 py-1 text-sm font-medium transition-colors ${
              filter === f.value
                ? 'bg-brand-600 text-white'
                : 'bg-white text-brand-600 hover:bg-brand-100'
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      <QueryState
        isLoading={isLoading}
        error={error}
        isEmpty={rows.length === 0}
        emptyText="Nenhuma ocorrência encontrada."
      >
        <Table
          head={
            <>
              <Th>Estudante</Th>
              <Th>Data</Th>
              <Th>Status</Th>
              <Th>Ação</Th>
            </>
          }
        >
          {rows.map((o) => (
            <tr key={o.id} className="hover:bg-brand-50">
              <Td>{studentName(o.studentId)}</Td>
              <Td>{formatDateTime(o.occurredAt)}</Td>
              <Td>
                <Badge tone={o.status}>{STATUS_LABELS[o.status]}</Badge>
              </Td>
              <Td>
                <Link
                  to={`/app/occurrences/${o.id}`}
                  className="text-brand-600 hover:underline"
                >
                  Ver
                </Link>
              </Td>
            </tr>
          ))}
        </Table>
      </QueryState>

      <NewOccurrenceForm open={showForm} onClose={() => setShowForm(false)} />
    </div>
  )
}
