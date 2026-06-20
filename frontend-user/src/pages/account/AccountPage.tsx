import { useTranslation } from 'react-i18next'
import { useAuthStore } from '../../store/slices/authStore'
import { authApi } from '../../api/auth'
import { useNavigate } from 'react-router-dom'

export default function AccountPage() {
  const { t } = useTranslation()
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = async () => {
    try { await authApi.logout() } catch {}
    logout()
    navigate('/login')
  }

  return (
    <div style={{ maxWidth: 600, margin: '0 auto' }}>
      <h2 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem' }}>{t('nav.account')}</h2>

      <div className="card" style={{ marginBottom: '1rem' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', fontSize: '0.875rem' }}>
          <div>
            <p style={{ color: '#9ca3af', marginBottom: '0.25rem' }}>{t('auth.firstName')}</p>
            <p style={{ fontWeight: 500 }}>{user?.first_name}</p>
          </div>
          <div>
            <p style={{ color: '#9ca3af', marginBottom: '0.25rem' }}>{t('auth.lastName')}</p>
            <p style={{ fontWeight: 500 }}>{user?.last_name || '—'}</p>
          </div>
          <div>
            <p style={{ color: '#9ca3af', marginBottom: '0.25rem' }}>{t('auth.email')}</p>
            <p style={{ fontWeight: 500 }}>{user?.email}</p>
          </div>
          <div>
            <p style={{ color: '#9ca3af', marginBottom: '0.25rem' }}>Role</p>
            <p style={{ fontWeight: 500, textTransform: 'capitalize' }}>{user?.role}</p>
          </div>
        </div>
      </div>

      <button onClick={handleLogout} style={{ background: '#ef4444', color: '#fff', border: 'none', borderRadius: '0.5rem', padding: '0.625rem 1.5rem', fontWeight: 600, cursor: 'pointer', fontSize: '0.875rem' }}>
        {t('auth.logout')}
      </button>
    </div>
  )
}
