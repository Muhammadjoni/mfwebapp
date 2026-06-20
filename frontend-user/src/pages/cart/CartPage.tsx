import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useCartStore } from '../../store/slices/cartStore'
import { useAuthStore } from '../../store/slices/authStore'

export default function CartPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { items, removeItem, updateQuantity, subtotal } = useCartStore()
  const isAuthenticated = useAuthStore(s => s.isAuthenticated)

  const handleCheckout = () => {
    if (!isAuthenticated) { navigate('/login'); return }
    navigate('/checkout')
  }

  if (items.length === 0) {
    return (
      <div style={{ textAlign: 'center', padding: '5rem 2rem' }}>
        <div style={{ fontSize: '4rem', marginBottom: '1rem' }}>🛒</div>
        <p style={{ color: '#6b7280', marginBottom: '1.5rem' }}>{t('cart.empty')}</p>
        <Link to="/products" className="btn-primary" style={{ display: 'inline-block', padding: '0.75rem 2rem' }}>{t('products.title')}</Link>
      </div>
    )
  }

  return (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 320px', gap: '2rem', alignItems: 'start' }}>
      <div>
        <h2 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1rem' }}>{t('nav.cart')}</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
          {items.map(item => (
            <div key={item.product.id} className="card" style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
              <div style={{ width: 80, height: 80, background: '#f3f4f6', borderRadius: '0.5rem', flexShrink: 0, display: 'flex', alignItems: 'center', justifyContent: 'center', overflow: 'hidden' }}>
                {item.product.images[0]
                  ? <img src={item.product.images[0]} alt={item.product.name} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                  : <span style={{ fontSize: '1.75rem' }}>📦</span>
                }
              </div>

              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{ fontWeight: 500, fontSize: '0.875rem', marginBottom: '0.25rem', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{item.product.name}</p>
                <p style={{ color: '#2563eb', fontWeight: 700, fontSize: '0.875rem' }}>
                  ${((item.product.sale_price ?? item.product.base_price) * item.quantity).toFixed(2)}
                </p>
              </div>

              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', flexShrink: 0 }}>
                <button onClick={() => updateQuantity(item.product.id, item.quantity - 1)} disabled={item.quantity <= 1}
                  style={{ width: 28, height: 28, border: '1px solid #d1d5db', borderRadius: '0.375rem', background: '#fff', cursor: 'pointer', fontSize: '1rem' }}>−</button>
                <span style={{ minWidth: 28, textAlign: 'center', fontWeight: 500 }}>{item.quantity}</span>
                <button onClick={() => updateQuantity(item.product.id, item.quantity + 1)}
                  style={{ width: 28, height: 28, border: '1px solid #d1d5db', borderRadius: '0.375rem', background: '#fff', cursor: 'pointer', fontSize: '1rem' }}>+</button>
              </div>

              <button onClick={() => removeItem(item.product.id)}
                style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer', flexShrink: 0, fontSize: '1.25rem' }}>×</button>
            </div>
          ))}
        </div>
      </div>

      <div className="card" style={{ position: 'sticky', top: '5rem' }}>
        <h3 style={{ fontWeight: 700, marginBottom: '1rem' }}>{t('cart.summary')}</h3>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem', fontSize: '0.875rem', color: '#4b5563' }}>
          <span>{t('cart.items', { count: items.reduce((s, i) => s + i.quantity, 0) })}</span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.25rem', fontWeight: 700, fontSize: '1.1rem', borderTop: '1px solid #e5e7eb', paddingTop: '0.75rem' }}>
          <span>{t('cart.total')}</span>
          <span>${subtotal().toFixed(2)}</span>
        </div>
        <button onClick={handleCheckout} className="btn-primary" style={{ width: '100%', padding: '0.875rem' }}>
          {t('cart.checkout')}
        </button>
      </div>
    </div>
  )
}
