import { NavLink, Outlet, useNavigate } from 'react-router-dom'

import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'
import { Locale } from '../i18n/messages'

export function DashboardLayout() {
  const { user, logout } = useAuth()
  const { locale, setLocale, t } = useI18n()
  const navigate = useNavigate()

  function onLogout() {
    logout()
    navigate('/login')
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="stack-md">
          <div>
            <p className="brand">{t('app.name')}</p>
            <p className="subtitle">{t('app.subtitle')}</p>
          </div>

          <div className="locale-picker row-between">
            <span className="muted">{t('common.language')}</span>
            <select value={locale} onChange={(event) => setLocale(event.target.value as Locale)}>
              <option value="zh-CN">中文</option>
              <option value="en-US">English</option>
            </select>
          </div>

          <nav className="menu">
            <NavLink to="/app/pm" className={({ isActive }) => (isActive ? 'menu-link active' : 'menu-link')}>
              {t('nav.projectManagement')}
            </NavLink>
            <NavLink to="/app/crm" className={({ isActive }) => (isActive ? 'menu-link active' : 'menu-link')}>
              {t('nav.crm')}
            </NavLink>
            <NavLink to="/app/templates" className={({ isActive }) => (isActive ? 'menu-link active' : 'menu-link')}>
              {t('nav.templates')}
            </NavLink>
          </nav>
        </div>

        <div className="account-box">
          <p>{user?.name}</p>
          <p className="muted">{user?.role ? t(`role.${user.role}`) : ''}</p>
          <button type="button" onClick={onLogout} className="btn-secondary full-width">
            {t('common.logout')}
          </button>
        </div>
      </aside>
      <main className="content">
        <Outlet />
      </main>
    </div>
  )
}
