import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { QueryState } from '../../components/QueryState'
import { Table, Td, Th } from '../../components/Table'
import { Input, PageTitle } from '../../components/ui'
import { studentsApi } from '../../lib/resources'

export function StudentsList() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')

  const { data, isLoading, error } = useQuery({
    queryKey: ['students', 'all'],
    // No server-side search param exists yet — fetch a wide page and filter
    // client-side. Flagged as a backend follow-up (search endpoint).
    queryFn: () => studentsApi.list({ limit: 200 }),
  })

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase()
    if (!q) return data ?? []
    return (data ?? []).filter((s) =>
      [s.name, s.cpf, s.ra, s.email].some((f) => f?.toLowerCase().includes(q)),
    )
  }, [data, search])

  return (
    <div>
      <PageTitle title="Estudantes" subtitle="Busque por nome, CPF, RA ou e-mail." />

      <div className="mb-4 max-w-md">
        <Input
          placeholder="Buscar estudante…"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      <QueryState
        isLoading={isLoading}
        error={error}
        isEmpty={filtered.length === 0}
        emptyText="Nenhum estudante encontrado."
      >
        <Table
          head={
            <>
              <Th>Nome</Th>
              <Th>CPF</Th>
              <Th>RA</Th>
              <Th>E-mail</Th>
            </>
          }
        >
          {filtered.map((s) => (
            <tr
              key={s.id}
              onClick={() => navigate(`/app/students/${s.id}`)}
              className="cursor-pointer hover:bg-brand-50"
            >
              <Td>{s.name}</Td>
              <Td>{s.cpf}</Td>
              <Td>{s.ra}</Td>
              <Td>{s.email}</Td>
            </tr>
          ))}
        </Table>
      </QueryState>
    </div>
  )
}
