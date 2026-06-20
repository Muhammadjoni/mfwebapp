import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAdminStore } from '../../store/adminStore'
import { authApi } from '../../api/auth'

const NAV_ITEMS = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/users', label: 'Users' },
  { to: '/sellers', label: 'Sellers' },
  { to: '/products', label: 'Products' },
  { to: '/orders', label: 'Orders' },
  { to: '/payments', label: 'Payments' },
]

export default function AdminLayout() {
  const { user, logout } = useAdminStore()
  const navigate = useNavigate()

  const handleLogout = async () => {
    try { await authApi.logout() } catch {}
    logout()
    navigate('/login')
  }

  return (
    <div style={{ display: 'flex', minHeight: '100vh', fontFamily: 'system-ui, sans-serif' }}>
      <aside style={{ width: 240, background: '#1e293b', color: '#cbd5e1', display: 'flex', flexDirection: 'column', flexShrink: 0 }}>
        <div style={{ padding: '1.5rem 1.25rem', borderBottom: '1px solid #334155' }}>
          <Link to="/dashboard" style={{ color: '#f1f5f9', fontWeight: 700, fontSize: '1.125rem' }}>MF Admin</Link>
        </div>

        <nav style={{ flex: 1, padding: '1rem 0' }}>
          {NAV_ITEMS.map(item => (
            <NavLink key={item.to} to={item.to} style={({ isActive }) => ({
              display: 'block', padding: '0.625rem 1.25rem', fontSize: '0.875rem', fontWeight: isActive ? 600 : 400,
              color: isActive ? '#f1f5f9' : '#94a3b8',
              background: isActive ? '#334155' : 'transparent',
              textDecoration: 'none',
              borderLeft: isActive ? '3px solid #3b82f6' : '3px solid transparent',
            })}>
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div style={{ padding: '1rem 1.25rem', borderTop: '1px solid #334155', fontSize: '0.875rem' }}>
          <p style={{ color: '#94a3b8', marginBottom: '0.5rem' }}>{user?.first_name} {user?.last_name}</p>
          <button onClick={handleLogout} style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer', fontSize: '0.875rem', padding: 0 }}>Logout</button>
        </div>
      </aside>

      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        <header style={{ background: '#fff', borderBottom: '1px solid #e2e8f0', padding: '0 1.5rem', height: 56, display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }}>
          <span style={{ fontSize: '0.875rem', color: '#64748b' }}>{user?.email}</span>
        </header>
        <main style={{ flex: 1, padding: '1.5rem', overflow: 'auto', background: '#f8fafc' }}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}
