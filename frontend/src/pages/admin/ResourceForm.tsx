import { useState } from 'react'
import { Button, Field, Input, Select, Textarea } from '../../components/ui'
import type { FieldConfig, Option } from './registry'

type Values = Record<string, string>

/**
 * Builds a form from a resource's field config. Validation is intentionally
 * light (required-only) — the API validates and returns `{error}` messages that
 * surface as toasts. Number fields are coerced on submit.
 */
export function ResourceForm({
  fields,
  options,
  initial,
  submitting,
  onSubmit,
  onCancel,
}: {
  fields: FieldConfig[]
  options: Record<string, Option[]>
  initial?: Values
  submitting: boolean
  onSubmit: (body: Record<string, unknown>) => void
  onCancel: () => void
}) {
  const [values, setValues] = useState<Values>(() => {
    const v: Values = {}
    for (const f of fields) v[f.name] = initial?.[f.name] ?? ''
    return v
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const set = (name: string, value: string) =>
    setValues((prev) => ({ ...prev, [name]: value }))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const errs: Record<string, string> = {}
    const body: Record<string, unknown> = {}

    for (const f of fields) {
      const raw = values[f.name]?.trim() ?? ''
      if (f.required && !raw) {
        errs[f.name] = 'Campo obrigatório.'
        continue
      }
      if (!raw) continue // omit empty optional fields (API disallows unknown, allows missing)
      body[f.name] = f.type === 'number' ? Number(raw) : raw
    }

    setErrors(errs)
    if (Object.keys(errs).length === 0) onSubmit(body)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4" noValidate>
      {fields.map((f) => {
        const opts = f.staticOptions ?? options[f.selectKey ?? ''] ?? null
        return (
          <Field key={f.name} label={f.label} error={errors[f.name]}>
            {opts ? (
              <Select value={values[f.name]} onChange={(e) => set(f.name, e.target.value)}>
                <option value="" disabled>
                  Selecione…
                </option>
                {opts.map((o) => (
                  <option key={o.value} value={o.value}>
                    {o.label}
                  </option>
                ))}
              </Select>
            ) : f.type === 'textarea' ? (
              <Textarea
                rows={3}
                value={values[f.name]}
                onChange={(e) => set(f.name, e.target.value)}
              />
            ) : (
              <Input
                type={f.type === 'number' ? 'number' : f.type === 'password' ? 'password' : 'text'}
                value={values[f.name]}
                onChange={(e) => set(f.name, e.target.value)}
              />
            )}
          </Field>
        )
      })}

      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="secondary" onClick={onCancel}>
          Cancelar
        </Button>
        <Button type="submit" disabled={submitting}>
          {submitting ? 'Salvando…' : 'Salvar'}
        </Button>
      </div>
    </form>
  )
}
