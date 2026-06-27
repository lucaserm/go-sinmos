import { Link } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { Card, PageTitle } from '../components/ui'
import { ROLE_LABELS, type Role } from '../types/api'

interface Tile {
  to: string
  title: string
  desc: string
  roles: Role[]
}

const TILES: Tile[] = [
  {
    to: '/app/students',
    title: 'Estudantes',
    desc: 'Buscar estudantes, ver horários, responsáveis e ocorrências.',
    roles: ['ADMIN', 'SUPPORT', 'RECEPTION'],
  },
  {
    to: '/app/occurrences',
    title: 'Ocorrências',
    desc: 'Registrar e acompanhar ocorrências e advertências.',
    roles: ['ADMIN', 'SUPPORT'],
  },
  {
    to: '/app/admin',
    title: 'Cadastros Gerais',
    desc: 'Cursos, disciplinas, horários, matrículas, usuários e mais.',
    roles: ['ADMIN'],
  },
]

export function Dashboard() {
  const { user } = useAuth()
  if (!user) return null
  const tiles = TILES.filter((t) => t.roles.includes(user.role))

  return (
    <div>
      <PageTitle
        title={`Olá, ${user.name.split(' ')[0]}`}
        subtitle={`Painel — ${ROLE_LABELS[user.role]}`}
      />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {tiles.map((t) => (
          <Link key={t.to} to={t.to}>
            <Card className="h-full transition-shadow hover:shadow-md">
              <h2 className="text-lg font-semibold text-brand-700">{t.title}</h2>
              <p className="mt-2 text-sm text-brand-600">{t.desc}</p>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  )
}
