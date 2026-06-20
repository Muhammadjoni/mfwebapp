import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuthStore } from '../../store/slices/authStore'
import { useCartStore } from '../../store/slices/cartStore'

export default function Navbar() {
  const { t, i18n } = useTranslation()
  const { isAuthenticated, user, logout } = useAuthStore()
  const totalItems = useCartStore(s => s.totalItems())

  const toggleLang = () => {
    const next = i18n.language === 'ru' ? 'tj' : 'ru'
    i18n.changeLanguage(next)
    localStorage.setItem('mf-lang', next)
  }

  return (
    <nav style={{ background: '#fff', borderBottom: '1px solid #e5e7eb', position: 'sticky', top: 0, zIndex: 50 }}>
      <div style={{ maxWidth: 1200, margin: '0 auto', padding: '0 1rem', height: 64, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <Link to="/" style={{ fontWeight: 700, fontSize: '1.25rem', color: '#2563eb' }}>MF Shop</Link>

        <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
          <Link to="/products" style={{ color: '#4b5563', fontSize: '0.875rem' }}>{t('nav.products')}</Link>

          <Link to="/cart" style={{ position: 'relative', color: '#4b5563', fontSize: '0.875rem' }}>
            {t('nav.cart')}
            {totalItems > 0 && (
              <span style={{ position: 'absolute', top: -8, right: -10, background: '#ef4444', color: '#fff', borderRadius: '50%', width: 18, height: 18, fontSize: 11, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                {totalItems}
              </span>
            )}
          </Link>

          {isAuthenticated ? (
            <>
              <Link to="/orders" style={{ color: '#4b5563', fontSize: '0.875rem' }}>{t('nav.orders')}</Link>
              <Link to="/account" style={{ color: '#4b5563', fontSize: '0.875rem' }}>{user?.first_name}</Link>
              <button onClick={logout} style={{ background: 'none', border: 'none', color: '#ef4444', fontSize: '0.875rem', cursor: 'pointer' }}>
                {t('auth.logout')}
              </button>
            </>
          ) : (
            <>
              <Link to="/login" style={{ color: '#4b5563', fontSize: '0.875rem' }}>{t('auth.login')}</Link>
              <Link to="/register" style={{ background: '#2563eb', color: '#fff', padding: '0.375rem 0.875rem', borderRadius: '0.5rem', fontSize: '0.875rem', fontWeight: 600 }}>
                {t('auth.register')}
              </Link>
            </>
          )}

          <button onClick={toggleLang} style={{ background: 'none', border: '1px solid #d1d5db', borderRadius: '0.375rem', padding: '0.25rem 0.5rem', fontSize: '0.75rem', cursor: 'pointer', color: '#4b5563' }}>
            {i18n.language.toUpperCase()}
          </button>
        </div>
      </div>
    </nav>
  )
}
