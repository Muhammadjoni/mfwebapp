import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { authApi } from '../../api/auth'
import { useAuthStore } from '../../store/slices/authStore'

export default function RegisterPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { setUser, setTokens } = useAuthStore()
  const [form, setForm] = useState({ first_name: '', last_name: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const set = (key: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm(f => ({ ...f, [key]: e.target.value }))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (form.password.length < 8) { setError(t('auth.passwordMin')); return }
    setLoading(true)
    try {
      const res = await authApi.register(form)
      setUser(res.user)
      setTokens(res.tokens.access_token, res.tokens.refresh_token)
      navigate('/')
    } catch (err: any) {
      const code = err.response?.data?.error?.code
      if (code === 'EMAIL_TAKEN') setError(t('auth.emailTaken'))
      else setError(t('common.error'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 400, margin: '4rem auto' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '1.5rem', textAlign: 'center' }}>{t('auth.register')}</h2>

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }} className="card">
        {error && <p style={{ color: '#ef4444', fontSize: '0.875rem', textAlign: 'center' }}>{error}</p>}

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem' }}>
          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, marginBottom: '0.375rem' }}>{t('auth.firstName')}</label>
            <input value={form.first_name} onChange={set('first_name')} required style={{ width: '100%' }} />
          </div>
          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, marginBottom: '0.375rem' }}>{t('auth.lastName')}</label>
            <input value={form.last_name} onChange={set('last_name')} style={{ width: '100%' }} />
          </div>
        </div>

        <div>
          <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, marginBottom: '0.375rem' }}>{t('auth.email')}</label>
          <input type="email" value={form.email} onChange={set('email')} required style={{ width: '100%' }} />
        </div>

        <div>
          <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, marginBottom: '0.375rem' }}>{t('auth.password')}</label>
          <input type="password" value={form.password} onChange={set('password')} required minLength={8} style={{ width: '100%' }} />
          <p style={{ fontSize: '0.75rem', color: '#9ca3af', marginTop: '0.25rem' }}>{t('auth.passwordMin')}</p>
        </div>

        <button type="submit" className="btn-primary" disabled={loading} style={{ width: '100%', padding: '0.75rem' }}>
          {loading ? t('common.loading') : t('auth.register')}
        </button>

        <p style={{ textAlign: 'center', fontSize: '0.875rem', color: '#6b7280' }}>
          {t('auth.hasAccount')}{' '}
          <Link to="/login" style={{ color: '#2563eb', fontWeight: 500 }}>{t('auth.login')}</Link>
        </p>
      </form>
    </div>
  )
}
