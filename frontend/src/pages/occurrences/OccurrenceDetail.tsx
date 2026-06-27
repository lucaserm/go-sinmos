import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useAuth } from '../../auth/useAuth'
import { QueryState } from '../../components/QueryState'
import { useToast } from '../../components/Toast'
import { Badge, Button, Card, Field, PageTitle, Spinner, Textarea } from '../../components/ui'
import { apiErrorMessage } from '../../lib/api'
import { formatDateTime } from '../../lib/format'
import {
  occurrencesApi,
  occurrenceTypesApi,
  studentsApi,
  warningsApi,
} from '../../lib/resources'
import { STATUS_LABELS } from '../../types/api'

export function OccurrenceDetail() {
  const { id = '' } = useParams()
  const { user } = useAuth()
  const toast = useToast()
  const qc = useQueryClient()
  const [report, setReport] = useState('')

  const occurrenceQuery = useQuery({
    queryKey: ['occurrences', id],
    queryFn: () => occurrencesApi.get(id),
  })
  const o = occurrenceQuery.data

  const students = useQuery({
    queryKey: ['students', 'all'],
    queryFn: () => studentsApi.list({ limit: 200 }),
  })
  const types = useQuery({
    queryKey: ['occurrence-types', 'all'],
    queryFn: () => occurrenceTypesApi.list({ limit: 200 }),
  })
  // No per-occurrence warning read — fetch all and match by occurrenceId.
  const warnings = useQuery({
    queryKey: ['warnings', 'all'],
    queryFn: () => warningsApi.list({ limit: 200 }),
  })

  const warning = warnings.data?.find((w) => w.occurrenceId === id)
  const studentName = students.data?.find((s) => s.id === o?.studentId)?.name ?? '—'
  const typeLabel = (() => {
    const t = types.data?.find((x) => x.id === o?.occurrenceTypeId)
    return t ? `${t.code} — ${t.description}` : '—'
  })()

  const approve = useMutation({
    mutationFn: async () => {
      await warningsApi.create({ occurrenceId: id, report })
      await occurrencesApi.update(id, { status: 'APPROVED' })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['occurrences'] })
      qc.invalidateQueries({ queryKey: ['warnings'] })
      toast.success('Advertência aprovada.')
    },
    onError: (err) => toast.error(apiErrorMessage(err)),
  })

  const reject = useMutation({
    mutationFn: () => occurrencesApi.update(id, { status: 'REPROVED' }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['occurrences'] })
      toast.success('Ocorrência reprovada.')
    },
    onError: (err) => toast.error(apiErrorMessage(err)),
  })

  if (occurrenceQuery.isLoading) return <Spinner />

  const canReview = user?.role === 'ADMIN' && o?.status === 'PENDING'

  return (
    <div className="space-y-6">
      <PageTitle title="Ocorrência" subtitle="Detalhes e parecer" />

      <QueryState isLoading={false} error={occurrenceQuery.error}>
        {o && (
          <Card>
            <dl className="grid grid-cols-2 gap-3 text-sm">
              <Info label="Estudante" value={studentName} />
              <Info label="Tipo" value={typeLabel} />
              <Info label="Data" value={formatDateTime(o.occurredAt)} />
              <div>
                <dt className="text-xs uppercase text-brand-500">Status</dt>
                <dd>
                  <Badge tone={o.status}>{STATUS_LABELS[o.status]}</Badge>
                </dd>
              </div>
            </dl>
          </Card>
        )}
      </QueryState>

      {warning && (
        <Card>
          <h2 className="mb-2 text-sm font-semibold uppercase text-brand-500">
            Advertência
          </h2>
          <p className="whitespace-pre-wrap text-sm text-brand-800">{warning.report}</p>
        </Card>
      )}

      {canReview && (
        <Card>
          <h2 className="mb-3 text-lg font-semibold text-brand-800">Parecer</h2>
          <Field label="Medidas que serão tomadas (relatório)">
            <Textarea
              rows={4}
              value={report}
              onChange={(e) => setReport(e.target.value)}
              placeholder="Descreva as medidas da advertência…"
            />
          </Field>
          <div className="mt-4 flex gap-2">
            <Button
              onClick={() => approve.mutate()}
              disabled={approve.isPending || report.trim().length === 0}
            >
              Aprovar Advertência
            </Button>
            <Button
              variant="danger"
              onClick={() => reject.mutate()}
              disabled={reject.isPending}
            >
              Reprovar
            </Button>
          </div>
        </Card>
      )}
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
