import { Navigate, Outlet, useLocation } from 'react-router-dom'

import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'

export function RequireAuth() {
  const { token, isReady } = useAuth()
  const { t } = useI18n()
  const location = useLocation()

  if (!isReady) {
    return <div className="screen-center">{t('common.loadingSession')}</div>
  }

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return <Outlet />
}
