import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { QueryState } from '../../components/QueryState'
import { Table, Td, Th } from '../../components/Table'
import { Badge, Card, PageTitle, Spinner } from '../../components/ui'
import { occurrencesApi, studentsApi } from '../../lib/resources'
import { STATUS_LABELS } from '../../types/api'
import { formatDateTime } from '../../lib/format'

export function StudentDetail() {
  const { id = '' } = useParams()

  const studentQuery = useQuery({
    queryKey: ['students', id],
    queryFn: () => studentsApi.get(id),
  })

  // No per-student occurrences endpoint — fetch and filter client-side.
  const occurrencesQuery = useQuery({
    queryKey: ['occurrences', 'all'],
    queryFn: () => occurrencesApi.list({ limit: 200 }),
  })

  const student = studentQuery.data
  const occurrences = (occurrencesQuery.data ?? []).filter((o) => o.studentId === id)

  if (studentQuery.isLoading) return <Spinner />

  return (
    <div className="space-y-6">
      <PageTitle title={student?.name ?? 'Estudante'} subtitle="Dados do estudante" />

      <QueryState isLoading={false} error={studentQuery.error}>
        {student && (
          <Card>
            <div className="flex items-start gap-6">
              {student.photo_url ? (
                <img
                  src={student.photo_url}
                  alt={student.name}
                  className="h-24 w-24 rounded-lg object-cover"
                />
              ) : (
                <div className="flex h-24 w-24 items-center justify-center rounded-lg bg-brand-100 text-2xl font-semibold text-brand-500">
                  {student.name.charAt(0)}
                </div>
              )}
              <dl className="grid flex-1 grid-cols-2 gap-3 text-sm">
                <Info label="Nome" value={student.name} />
                <Info label="E-mail" value={student.email} />
                <Info label="CPF" value={student.cpf} />
                <Info label="RA" value={student.ra} />
              </dl>
            </div>
          </Card>
        )}
      </QueryState>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-brand-800">Ocorrências</h2>
        <QueryState
          isLoading={occurrencesQuery.isLoading}
          error={occurrencesQuery.error}
          isEmpty={occurrences.length === 0}
          emptyText="Nenhuma ocorrência para este estudante."
        >
          <Table
            head={
              <>
                <Th>Data</Th>
                <Th>Status</Th>
                <Th>Ação</Th>
              </>
            }
          >
            {occurrences.map((o) => (
              <tr key={o.id} className="hover:bg-brand-50">
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
      </section>

      <div className="rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
        Horários e responsáveis exigem endpoints de leitura das tabelas de
        associação (ainda não disponíveis na API). Pendência de backend.
      </div>
    </div>
  )
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs uppercase text-brand-500">{label}</dt>
      <dd className="font-medium text-brand-800">{value}</dd>
    </div>
  )
}
