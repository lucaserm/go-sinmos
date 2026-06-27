import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { Modal } from '../../components/Modal'
import { useToast } from '../../components/Toast'
import { Button, Field, Input, Select } from '../../components/ui'
import { apiErrorMessage } from '../../lib/api'
import { toRFC3339 } from '../../lib/format'
import { occurrencesApi, occurrenceTypesApi, studentsApi } from '../../lib/resources'

const schema = z.object({
  studentId: z.string().uuid('Selecione um estudante.'),
  occurrenceTypeId: z.string().uuid('Selecione um tipo de ocorrência.'),
  occurredAt: z.string().min(1, 'Informe a data e hora.'),
  status: z.enum(['PENDING', 'APPROVED', 'REPROVED']),
})

type FormValues = z.infer<typeof schema>

export function NewOccurrenceForm({
  open,
  onClose,
}: {
  open: boolean
  onClose: () => void
}) {
  const toast = useToast()
  const qc = useQueryClient()

  const students = useQuery({
    queryKey: ['students', 'all'],
    queryFn: () => studentsApi.list({ limit: 200 }),
    enabled: open,
  })
  const types = useQuery({
    queryKey: ['occurrence-types', 'all'],
    queryFn: () => occurrenceTypesApi.list({ limit: 200 }),
    enabled: open,
  })

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { status: 'PENDING' },
  })

  const mutation = useMutation({
    mutationFn: (values: FormValues) =>
      // userId comes from the bearer token, not the body. occurredAt → RFC3339.
      occurrencesApi.create({
        studentId: values.studentId,
        occurrenceTypeId: values.occurrenceTypeId,
        occurredAt: toRFC3339(values.occurredAt),
        status: values.status,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['occurrences'] })
      toast.success('Ocorrência registrada.')
      reset()
      onClose()
    },
    onError: (err) => toast.error(apiErrorMessage(err)),
  })

  return (
    <Modal open={open} onClose={onClose} title="Registrar Ocorrência">
      <form
        onSubmit={handleSubmit((v) => mutation.mutate(v))}
        className="space-y-4"
        noValidate
      >
        <Field label="Estudante" error={errors.studentId?.message}>
          <Select {...register('studentId')} defaultValue="">
            <option value="" disabled>
              Selecione…
            </option>
            {(students.data ?? []).map((s) => (
              <option key={s.id} value={s.id}>
                {s.name} — {s.ra}
              </option>
            ))}
          </Select>
        </Field>

        <Field label="Tipo de Ocorrência" error={errors.occurrenceTypeId?.message}>
          <Select {...register('occurrenceTypeId')} defaultValue="">
            <option value="" disabled>
              Selecione…
            </option>
            {(types.data ?? []).map((t) => (
              <option key={t.id} value={t.id}>
                {t.code} — {t.description}
              </option>
            ))}
          </Select>
        </Field>

        <Field label="Data e Hora" error={errors.occurredAt?.message}>
          <Input type="datetime-local" {...register('occurredAt')} />
        </Field>

        <Field label="Encaminhamento" error={errors.status?.message}>
          <Select {...register('status')}>
            <option value="PENDING">Ocorrência (pendente)</option>
            <option value="APPROVED">Aprovada</option>
            <option value="REPROVED">Reprovada</option>
          </Select>
        </Field>

        <div className="flex justify-end gap-2 pt-2">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isSubmitting || mutation.isPending}>
            {mutation.isPending ? 'Salvando…' : 'Registrar'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
