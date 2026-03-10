import { Navigate, Outlet, useLocation } from 'react-router-dom'

import { useAuth } from '../context/AuthContext'

export function RequireAuth() {
  const { token, isReady } = useAuth()
  const location = useLocation()

  if (!isReady) {
    return <div className="screen-center">Loading session...</div>
  }

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return <Outlet />
}
