// Mirrors the Go API DTOs (internal/*/types.go). Response envelopes are
// single-key wrapped, e.g. {"student": {...}} / {"students": [...]}.

export type Role = 'ADMIN' | 'SUPPORT' | 'RECEPTION'
export type Session = 'MORNING' | 'AFTERNOON' | 'EVENING'
export type PermissionType = 'LEAVE' | 'STAY'
export type OccurrenceStatus = 'PENDING' | 'APPROVED' | 'REPROVED'

export interface User {
  id: string
  name: string
  code: string
  role: Role
}

export interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
}

export interface Student {
  id: string
  cpf: string
  ra: string
  photo_url?: string
  name: string
  email: string
}

export interface Guardian {
  id: string
  name: string
  email: string
  phone: string
}

export interface Course {
  id: string
  name: string
  description?: string
  session: Session
}

export interface Subject {
  id: string
  name: string
  semester: number
  section: string
  courseId: string
}

export interface Schedule {
  id: string
  subjectId: string
  session: Session
  dayOfWeek: string
  startTime: string
  endTime: string
}

export interface Enrollment {
  id: string
  studentId: string
  courseId: string
  year: number
}

export interface OccurrenceType {
  id: string
  code: string
  description: string
  severity: number
}

export interface Occurrence {
  id: string
  userId: string
  occurrenceTypeId: string
  studentId: string
  occurredAt: string
  userRelatedId?: string
  status: OccurrenceStatus
}

export interface Warning {
  id: string
  occurrenceId: string
  report: string
}

export const ROLE_LABELS: Record<Role, string> = {
  ADMIN: 'Coordenação',
  SUPPORT: 'Assistência',
  RECEPTION: 'Portaria',
}

export const SESSION_LABELS: Record<Session, string> = {
  MORNING: 'Matutino',
  AFTERNOON: 'Vespertino',
  EVENING: 'Noturno',
}

export const STATUS_LABELS: Record<OccurrenceStatus, string> = {
  PENDING: 'Pendente',
  APPROVED: 'Aprovada',
  REPROVED: 'Reprovada',
}
