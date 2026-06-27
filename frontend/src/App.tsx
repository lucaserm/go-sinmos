import { Navigate, Route, Routes } from 'react-router-dom'
import { RequireAuth, RequireRole } from './auth/guards'
import { Layout } from './components/Layout'
import { AdminHub } from './pages/admin/AdminHub'
import { ResourcePage } from './pages/admin/ResourcePage'
import { Dashboard } from './pages/Dashboard'
import { Landing } from './pages/Landing'
import { Login } from './pages/Login'
import { OccurrenceDetail } from './pages/occurrences/OccurrenceDetail'
import { OccurrencesList } from './pages/occurrences/OccurrencesList'
import { StudentDetail } from './pages/students/StudentDetail'
import { StudentsList } from './pages/students/StudentsList'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/login" element={<Login />} />

      <Route element={<RequireAuth />}>
        <Route path="/app" element={<Layout />}>
          <Route index element={<Dashboard />} />

          <Route path="students" element={<StudentsList />} />
          <Route path="students/:id" element={<StudentDetail />} />

          <Route element={<RequireRole roles={['ADMIN', 'SUPPORT']} />}>
            <Route path="occurrences" element={<OccurrencesList />} />
            <Route path="occurrences/:id" element={<OccurrenceDetail />} />
          </Route>

          <Route element={<RequireRole roles={['ADMIN']} />}>
            <Route path="admin" element={<AdminHub />} />
            <Route path="admin/:resource" element={<ResourcePage />} />
          </Route>
        </Route>
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
