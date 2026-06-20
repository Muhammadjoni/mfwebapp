import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

export default function HomePage() {
  const { t } = useTranslation()
  return (
    <div>
      <section style={{ background: 'linear-gradient(135deg,#1d4ed8,#3b82f6)', borderRadius: '1rem', padding: '4rem 2rem', textAlign: 'center', color: '#fff', marginBottom: '3rem' }}>
        <h1 style={{ fontSize: '2.5rem', fontWeight: 800, marginBottom: '1rem' }}>MF Shop</h1>
        <p style={{ fontSize: '1.125rem', opacity: 0.85, marginBottom: '2rem' }}>Tech & Electronics Marketplace — RU / TJ</p>
        <Link to="/products" style={{ background: '#fff', color: '#1d4ed8', fontWeight: 700, padding: '0.75rem 2rem', borderRadius: '0.75rem', fontSize: '1rem' }}>
          {t('products.title')}
        </Link>
      </section>
    </div>
  )
}
