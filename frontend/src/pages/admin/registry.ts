import type { ReactNode } from 'react'
import { api } from '../../lib/api'
import {
  coursesApi,
  enrollmentsApi,
  guardiansApi,
  occurrenceTypesApi,
  schedulesApi,
  studentsApi,
  subjectsApi,
} from '../../lib/resources'
import { SESSION_LABELS, type Session } from '../../types/api'

type Row = Record<string, unknown>
export type Option = { value: string; label: string }

export interface FieldConfig {
  name: string
  label: string
  type?: 'text' | 'number' | 'textarea' | 'password'
  required?: boolean
  /** Static <select> options (enums). */
  staticOptions?: Option[]
  /** Dynamic <select> options loaded from another resource. */
  selectKey?: string
  selectOptions?: () => Promise<Option[]>
}

export interface ColumnConfig {
  key: string
  label: string
  render?: (row: Row) => ReactNode
}

export interface ResourceConfig {
  key: string
  title: string
  subtitle: string
  list?: () => Promise<unknown[]>
  get?: (id: string) => Promise<unknown>
  create: (body: unknown) => Promise<unknown>
  update?: (id: string, body: unknown) => Promise<unknown>
  remove?: (id: string) => Promise<void>
  columns: ColumnConfig[]
  fields: FieldConfig[]
}

const sessionOptions: Option[] = (Object.keys(SESSION_LABELS) as Session[]).map((s) => ({
  value: s,
  label: SESSION_LABELS[s],
}))

const courseOptions = async () =>
  (await coursesApi.list({ limit: 200 })).map((c) => ({ value: c.id, label: c.name }))
const subjectOptions = async () =>
  (await subjectsApi.list({ limit: 200 })).map((s) => ({ value: s.id, label: s.name }))
const studentOptions = async () =>
  (await studentsApi.list({ limit: 200 })).map((s) => ({
    value: s.id,
    label: `${s.name} — ${s.ra}`,
  }))

export const REGISTRY: Record<string, ResourceConfig> = {
  students: {
    key: 'students',
    title: 'Estudantes',
    subtitle: 'Cadastro de estudantes.',
    list: () => studentsApi.list({ limit: 200 }),
    create: studentsApi.create,
    update: studentsApi.update,
    remove: studentsApi.remove,
    columns: [
      { key: 'name', label: 'Nome' },
      { key: 'cpf', label: 'CPF' },
      { key: 'ra', label: 'RA' },
      { key: 'email', label: 'E-mail' },
    ],
    fields: [
      { name: 'name', label: 'Nome', required: true },
      { name: 'email', label: 'E-mail', type: 'text', required: true },
      { name: 'cpf', label: 'CPF', required: true },
      { name: 'ra', label: 'RA', required: true },
      { name: 'photo_url', label: 'URL da Foto' },
    ],
  },

  courses: {
    key: 'courses',
    title: 'Cursos',
    subtitle: 'Cadastro de cursos.',
    list: () => coursesApi.list({ limit: 200 }),
    create: coursesApi.create,
    update: coursesApi.update,
    remove: coursesApi.remove,
    columns: [
      { key: 'name', label: 'Nome' },
      { key: 'description', label: 'Descrição' },
      {
        key: 'session',
        label: 'Período',
        render: (r) => SESSION_LABELS[r.session as Session] ?? String(r.session),
      },
    ],
    fields: [
      { name: 'name', label: 'Nome', required: true },
      { name: 'description', label: 'Descrição', type: 'textarea' },
      { name: 'session', label: 'Período', required: true, staticOptions: sessionOptions },
    ],
  },

  subjects: {
    key: 'subjects',
    title: 'Disciplinas',
    subtitle: 'Cadastro de disciplinas.',
    list: () => subjectsApi.list({ limit: 200 }),
    create: subjectsApi.create,
    update: subjectsApi.update,
    remove: subjectsApi.remove,
    columns: [
      { key: 'name', label: 'Nome' },
      { key: 'semester', label: 'Semestre' },
      { key: 'section', label: 'Turma' },
    ],
    fields: [
      { name: 'name', label: 'Nome', required: true },
      { name: 'semester', label: 'Semestre', type: 'number', required: true },
      { name: 'section', label: 'Turma', required: true },
      {
        name: 'courseId',
        label: 'Curso',
        required: true,
        selectKey: 'courses',
        selectOptions: courseOptions,
      },
    ],
  },

  schedules: {
    key: 'schedules',
    title: 'Horários',
    subtitle: 'Horários das disciplinas. Use o formato HH:MM (ex.: 08:00).',
    list: () => schedulesApi.list({ limit: 200 }),
    create: schedulesApi.create,
    update: schedulesApi.update,
    remove: schedulesApi.remove,
    columns: [
      { key: 'dayOfWeek', label: 'Dia' },
      {
        key: 'session',
        label: 'Período',
        render: (r) => SESSION_LABELS[r.session as Session] ?? String(r.session),
      },
      { key: 'startTime', label: 'Início' },
      { key: 'endTime', label: 'Fim' },
    ],
    fields: [
      {
        name: 'subjectId',
        label: 'Disciplina',
        required: true,
        selectKey: 'subjects',
        selectOptions: subjectOptions,
      },
      { name: 'session', label: 'Período', required: true, staticOptions: sessionOptions },
      { name: 'dayOfWeek', label: 'Dia da Semana', required: true },
      { name: 'startTime', label: 'Início (HH:MM)', required: true },
      { name: 'endTime', label: 'Fim (HH:MM)', required: true },
    ],
  },

  enrollments: {
    key: 'enrollments',
    title: 'Matrículas',
    subtitle: 'Matrícula de estudantes em cursos por ano.',
    list: () => enrollmentsApi.list({ limit: 200 }),
    create: enrollmentsApi.create,
    update: enrollmentsApi.update,
    remove: enrollmentsApi.remove,
    columns: [
      { key: 'studentId', label: 'Estudante (id)' },
      { key: 'courseId', label: 'Curso (id)' },
      { key: 'year', label: 'Ano' },
    ],
    fields: [
      {
        name: 'studentId',
        label: 'Estudante',
        required: true,
        selectKey: 'students',
        selectOptions: studentOptions,
      },
      {
        name: 'courseId',
        label: 'Curso',
        required: true,
        selectKey: 'courses',
        selectOptions: courseOptions,
      },
      { name: 'year', label: 'Ano', type: 'number', required: true },
    ],
  },

  'occurrence-types': {
    key: 'occurrence-types',
    title: 'Tipos de Ocorrência',
    subtitle: 'Catálogo de tipos de ocorrência.',
    list: () => occurrenceTypesApi.list({ limit: 200 }),
    create: occurrenceTypesApi.create,
    update: occurrenceTypesApi.update,
    remove: occurrenceTypesApi.remove,
    columns: [
      { key: 'code', label: 'Código' },
      { key: 'description', label: 'Descrição' },
      { key: 'severity', label: 'Gravidade' },
    ],
    fields: [
      { name: 'code', label: 'Código', required: true },
      { name: 'description', label: 'Descrição', required: true },
      { name: 'severity', label: 'Gravidade', type: 'number', required: true },
    ],
  },

  guardians: {
    key: 'guardians',
    title: 'Responsáveis',
    subtitle: 'Cadastro de responsáveis.',
    list: () => guardiansApi.list({ limit: 200 }),
    create: guardiansApi.create,
    update: guardiansApi.update,
    remove: guardiansApi.remove,
    columns: [
      { key: 'name', label: 'Nome' },
      { key: 'email', label: 'E-mail' },
      { key: 'phone', label: 'Telefone' },
    ],
    fields: [
      { name: 'name', label: 'Nome', required: true },
      { name: 'email', label: 'E-mail', required: true },
      { name: 'phone', label: 'Telefone', required: true },
    ],
  },

  // Users are create-only: the API exposes registration but no list/update/delete.
  users: {
    key: 'users',
    title: 'Usuários',
    subtitle: 'A API permite apenas o cadastro de usuários (sem listagem).',
    create: (body) => api.post('/auth/register', body).then((r) => r.data),
    columns: [],
    fields: [
      { name: 'name', label: 'Nome', required: true },
      { name: 'code', label: 'Código do Servidor', required: true },
      { name: 'password', label: 'Senha', type: 'password', required: true },
      {
        name: 'role',
        label: 'Cargo',
        required: true,
        staticOptions: [
          { value: 'ADMIN', label: 'Coordenação (ADMIN)' },
          { value: 'SUPPORT', label: 'Assistência (SUPPORT)' },
          { value: 'RECEPTION', label: 'Portaria (RECEPTION)' },
        ],
      },
    ],
  },
}

export const HUB_ORDER = [
  'students',
  'courses',
  'subjects',
  'schedules',
  'enrollments',
  'occurrence-types',
  'guardians',
  'users',
]
