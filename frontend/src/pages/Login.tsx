import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { z } from 'zod'
import { useAuth } from '../auth/useAuth'
import { Button, Card, Field, Input } from '../components/ui'
import { apiErrorMessage } from '../lib/api'

const schema = z.object({
  code: z.string().min(1, 'Informe o código do servidor.'),
  password: z.string().min(6, 'A senha deve ter ao menos 6 caracteres.'),
})

type FormValues = z.infer<typeof schema>

export function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation() as { state?: { from?: { pathname?: string } } }
  const [serverError, setServerError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema) })

  const onSubmit = handleSubmit(async (values) => {
    setServerError(null)
    try {
      await login(values.code, values.password)
      navigate(location.state?.from?.pathname ?? '/app', { replace: true })
    } catch (err) {
      setServerError(apiErrorMessage(err, 'Usuário ou senha inválidos!'))
    }
  })

  return (
    <div className="flex min-h-full items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <div className="mb-6 text-center">
          <Link to="/" className="text-3xl font-bold text-brand-600">
            SINMOS
          </Link>
          <p className="mt-1 text-sm text-brand-600">Acesso restrito a servidores</p>
        </div>

        <form onSubmit={onSubmit} className="space-y-4" noValidate>
          <Field label="Código do Servidor" error={errors.code?.message}>
            <Input autoFocus {...register('code')} />
          </Field>
          <Field label="Senha" error={errors.password?.message}>
            <Input type="password" {...register('password')} />
          </Field>

          {serverError && (
            <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
              {serverError}
            </div>
          )}

          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Entrando…' : 'Entrar'}
          </Button>
        </form>
      </Card>
    </div>
  )
}
