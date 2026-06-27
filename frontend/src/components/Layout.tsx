import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { ROLE_LABELS, type Role } from '../types/api'
import { Button } from './ui'

interface NavItem {
  to: string
  label: string
  roles: Role[]
}

const NAV: NavItem[] = [
  { to: '/app/students', label: 'Estudantes', roles: ['ADMIN', 'SUPPORT', 'RECEPTION'] },
  { to: '/app/occurrences', label: 'Ocorrências', roles: ['ADMIN', 'SUPPORT'] },
  { to: '/app/admin', label: 'Cadastros', roles: ['ADMIN'] },
]

export function Layout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login', { replace: true })
  }

  const items = NAV.filter((i) => user && i.roles.includes(user.role))

  return (
    <div className="flex min-h-full flex-col">
      <header className="border-b border-brand-100 bg-white">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-4">
          <div className="flex items-center gap-6">
            <NavLink to="/app" className="text-xl font-bold text-brand-600">
              SINMOS
            </NavLink>
            <nav className="flex gap-1">
              {items.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  className={({ isActive }) =>
                    `rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                      isActive
                        ? 'bg-brand-100 text-brand-700'
                        : 'text-brand-600 hover:bg-brand-50'
                    }`
                  }
                >
                  {item.label}
                </NavLink>
              ))}
            </nav>
          </div>

          <div className="flex items-center gap-3">
            {user && (
              <span className="hidden text-right text-sm sm:block">
                <span className="block font-medium text-brand-800">{user.name}</span>
                <span className="block text-xs text-brand-500">{ROLE_LABELS[user.role]}</span>
              </span>
            )}
            <Button variant="ghost" onClick={handleLogout}>
              Sair
            </Button>
          </div>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
