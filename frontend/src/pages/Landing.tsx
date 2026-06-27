import { Link } from 'react-router-dom'
import { Button } from '../components/ui'

export function Landing() {
  return (
    <div className="flex min-h-full flex-col items-center justify-center bg-brand-800 p-6 text-center text-white">
      <h1 className="text-5xl font-bold tracking-tight">SINMOS</h1>
      <p className="mt-3 max-w-md text-brand-100">
        Sistema Integrado de Monitoramento e Segurança
      </p>
      <Link to="/login" className="mt-8">
        <Button variant="secondary">Login</Button>
      </Link>
    </div>
  )
}
