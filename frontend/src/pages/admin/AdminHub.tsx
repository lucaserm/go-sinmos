import { Link } from 'react-router-dom'
import { Card, PageTitle } from '../../components/ui'
import { HUB_ORDER, REGISTRY } from './registry'

export function AdminHub() {
  return (
    <div>
      <PageTitle title="Cadastros Gerais" subtitle="Gerencie os dados do sistema." />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {HUB_ORDER.map((key) => {
          const r = REGISTRY[key]
          return (
            <Link key={key} to={`/app/admin/${key}`}>
              <Card className="h-full transition-shadow hover:shadow-md">
                <h2 className="text-lg font-semibold text-brand-700">{r.title}</h2>
                <p className="mt-2 text-sm text-brand-600">{r.subtitle}</p>
              </Card>
            </Link>
          )
        })}
      </div>
    </div>
  )
}
