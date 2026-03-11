import { Navigate, Route, Routes } from 'react-router-dom'

import { RequireAuth } from './components/RequireAuth'
import { useAuth } from './context/AuthContext'
import { DashboardLayout } from './layout/DashboardLayout'
import { BoardTemplateMVPPage } from './pages/BoardTemplateMVPPage'
import { CRMPage } from './pages/CRMPage'
import { LoginPage } from './pages/LoginPage'
import { PMPage } from './pages/PMPage'
import { TemplateLifecyclePage } from './pages/TemplateLifecyclePage'

export default function App() {
  const { token } = useAuth()

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route element={<RequireAuth />}>
        <Route path="/app" element={<DashboardLayout />}>
          <Route index element={<Navigate to="pm" replace />} />
          <Route path="pm" element={<PMPage />} />
          <Route path="crm" element={<CRMPage />} />
          <Route path="templates" element={<TemplateLifecyclePage />} />
          <Route path="board-template-mvp" element={<BoardTemplateMVPPage />} />
        </Route>
      </Route>
      <Route path="*" element={<Navigate to={token ? '/app/pm' : '/login'} replace />} />
    </Routes>
  )
}
