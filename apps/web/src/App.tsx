import { Navigate, Route, Routes } from 'react-router-dom'

import { RequireAuth } from './components/RequireAuth'
import { useAuth } from './context/AuthContext'
import { DashboardLayout } from './layout/DashboardLayout'
import { LoginPage } from './pages/LoginPage'
import { PMPage } from './pages/PMPage'

export default function App() {
  const { token } = useAuth()

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route element={<RequireAuth />}>
        <Route path="/app" element={<DashboardLayout />}>
          <Route index element={<Navigate to="pm" replace />} />
          <Route path="pm" element={<PMPage />} />
        </Route>
      </Route>
      <Route path="*" element={<Navigate to={token ? '/app/pm' : '/login'} replace />} />
    </Routes>
  )
}
