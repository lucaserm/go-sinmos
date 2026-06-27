import { createResource } from './resource'
import type {
  Course,
  Enrollment,
  Guardian,
  Occurrence,
  OccurrenceType,
  Schedule,
  Student,
  Subject,
  Warning,
} from '../types/api'

export const studentsApi = createResource<Student>('/students', 'student', 'students')
export const guardiansApi = createResource<Guardian>('/guardians', 'guardian', 'guardians')
export const coursesApi = createResource<Course>('/courses', 'course', 'courses')
export const subjectsApi = createResource<Subject>('/subjects', 'subject', 'subjects')
export const schedulesApi = createResource<Schedule>('/schedules', 'schedule', 'schedules')
export const enrollmentsApi = createResource<Enrollment>(
  '/enrollments',
  'enrollment',
  'enrollments',
)
// occurrence-types uses camelCase envelope keys, unlike the others.
export const occurrenceTypesApi = createResource<OccurrenceType>(
  '/occurrence-types',
  'occurrenceType',
  'occurrenceTypes',
)
export const occurrencesApi = createResource<Occurrence>(
  '/occurrences',
  'occurrence',
  'occurrences',
)
export const warningsApi = createResource<Warning>('/warnings', 'warning', 'warnings')
