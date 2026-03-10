import { NavLink, Outlet, useNavigate } from 'react-router-dom'

import { useAuth } from '../context/AuthContext'

export function DashboardLayout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  function onLogout() {
    logout()
    navigate('/login')
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div>
          <p className="brand">通用管理系统</p>
          <p className="subtitle">Core + PM</p>
        </div>
        <nav className="menu">
          <NavLink to="/app/pm" className={({ isActive }) => (isActive ? 'menu-link active' : 'menu-link')}>
            Project Management
          </NavLink>
        </nav>
        <div className="account-box">
          <p>{user?.name}</p>
          <p className="muted">{user?.role}</p>
          <button type="button" onClick={onLogout} className="btn-secondary full-width">
            Logout
          </button>
        </div>
      </aside>
      <main className="content">
        <Outlet />
      </main>
    </div>
  )
}
