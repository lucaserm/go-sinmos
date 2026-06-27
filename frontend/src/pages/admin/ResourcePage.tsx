import { useMutation, useQueries, useQuery, useQueryClient } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { Navigate, useParams } from 'react-router-dom'
import { Modal } from '../../components/Modal'
import { QueryState } from '../../components/QueryState'
import { Table, Td, Th } from '../../components/Table'
import { useToast } from '../../components/Toast'
import { Button, PageTitle } from '../../components/ui'
import { apiErrorMessage } from '../../lib/api'
import { REGISTRY, type Option } from './registry'
import { ResourceForm } from './ResourceForm'

type Row = Record<string, unknown> & { id: string }

export function ResourcePage() {
  const { resource = '' } = useParams()
  const config = REGISTRY[resource]
  const toast = useToast()
  const qc = useQueryClient()

  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Row | null>(null)

  // Load dynamic <select> option sources for this resource's fields.
  const selectFields = useMemo(
    () => (config?.fields ?? []).filter((f) => f.selectOptions),
    [config],
  )
  const optionQueries = useQueries({
    queries: selectFields.map((f) => ({
      queryKey: ['options', f.selectKey],
      queryFn: f.selectOptions!,
    })),
  })
  const options: Record<string, Option[]> = {}
  selectFields.forEach((f, i) => {
    options[f.selectKey!] = (optionQueries[i].data as Option[]) ?? []
  })

  const listQuery = useQuery({
    queryKey: ['resource', resource],
    queryFn: () => config.list!() as Promise<Row[]>,
    enabled: !!config?.list,
  })

  const invalidate = () => {
    qc.invalidateQueries({ queryKey: ['resource', resource] })
    // shared caches the rest of the app reads from
    qc.invalidateQueries({ queryKey: [resource] })
    qc.invalidateQueries({ queryKey: ['options', resource] })
  }

  const createMut = useMutation({
    mutationFn: (body: Record<string, unknown>) => config.create(body),
    onSuccess: () => {
      invalidate()
      toast.success('Registro salvo.')
      setModalOpen(false)
    },
    onError: (err) => toast.error(apiErrorMessage(err)),
  })

  const updateMut = useMutation({
    mutationFn: ({ id, body }: { id: string; body: Record<string, unknown> }) =>
      config.update!(id, body),
    onSuccess: () => {
      invalidate()
      toast.success('Registro atualizado.')
      setModalOpen(false)
      setEditing(null)
    },
    onError: (err) => toast.error(apiErrorMessage(err)),
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => config.remove!(id),
    onSuccess: () => {
      invalidate()
      toast.success('Registro removido.')
    },
    onError: (err) => toast.error(apiErrorMessage(err)),
  })

  if (!config) return <Navigate to="/app/admin" replace />

  const openCreate = () => {
    setEditing(null)
    setModalOpen(true)
  }
  const openEdit = (row: Row) => {
    setEditing(row)
    setModalOpen(true)
  }

  const handleSubmit = (body: Record<string, unknown>) => {
    if (editing) updateMut.mutate({ id: editing.id, body })
    else createMut.mutate(body)
  }

  const initialValues = editing
    ? Object.fromEntries(
        config.fields.map((f) => [
          f.name,
          editing[f.name] != null ? String(editing[f.name]) : '',
        ]),
      )
    : undefined

  const showActions = !!config.update || !!config.remove

  return (
    <div>
      <PageTitle
        title={config.title}
        subtitle={config.subtitle}
        action={<Button onClick={openCreate}>Novo</Button>}
      />

      {!config.list ? (
        <div className="rounded-lg border border-brand-100 bg-white px-4 py-6 text-sm text-brand-600">
          Este recurso não possui listagem na API. Use o botão “Novo” para cadastrar.
        </div>
      ) : (
        <QueryState
          isLoading={listQuery.isLoading}
          error={listQuery.error}
          isEmpty={(listQuery.data?.length ?? 0) === 0}
        >
          <Table
            head={
              <>
                {config.columns.map((c) => (
                  <Th key={c.key}>{c.label}</Th>
                ))}
                {showActions && <Th>Ações</Th>}
              </>
            }
          >
            {(listQuery.data ?? []).map((row) => (
              <tr key={row.id} className="hover:bg-brand-50">
                {config.columns.map((c) => (
                  <Td key={c.key}>
                    {c.render ? c.render(row) : String(row[c.key] ?? '—')}
                  </Td>
                ))}
                {showActions && (
                  <Td>
                    <div className="flex gap-3">
                      {config.update && (
                        <button
                          onClick={() => openEdit(row)}
                          className="text-brand-600 hover:underline"
                        >
                          Editar
                        </button>
                      )}
                      {config.remove && (
                        <button
                          onClick={() => {
                            if (confirm('Remover este registro?')) deleteMut.mutate(row.id)
                          }}
                          className="text-red-600 hover:underline"
                        >
                          Excluir
                        </button>
                      )}
                    </div>
                  </Td>
                )}
              </tr>
            ))}
          </Table>
        </QueryState>
      )}

      <Modal
        open={modalOpen}
        onClose={() => {
          setModalOpen(false)
          setEditing(null)
        }}
        title={editing ? `Editar ${config.title}` : `Novo — ${config.title}`}
      >
        <ResourceForm
          fields={config.fields}
          options={options}
          initial={initialValues}
          submitting={createMut.isPending || updateMut.isPending}
          onSubmit={handleSubmit}
          onCancel={() => {
            setModalOpen(false)
            setEditing(null)
          }}
        />
      </Modal>
    </div>
  )
}
