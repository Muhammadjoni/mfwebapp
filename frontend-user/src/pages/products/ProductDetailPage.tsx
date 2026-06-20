import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { productsApi } from '../../api/products'
import { useCartStore } from '../../store/slices/cartStore'

export default function ProductDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { t } = useTranslation()
  const addItem = useCartStore(s => s.addItem)

  const { data: product, isLoading } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productsApi.get(id!),
    enabled: !!id,
  })

  if (isLoading) return <p style={{ textAlign: 'center', padding: '3rem', color: '#6b7280' }}>{t('common.loading')}</p>
  if (!product) return <p style={{ textAlign: 'center', padding: '3rem', color: '#ef4444' }}>{t('common.error')}</p>

  const price = product.sale_price ?? product.base_price

  return (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem', alignItems: 'start' }}>
      <div style={{ background: '#f3f4f6', borderRadius: '0.75rem', aspectRatio: '1', overflow: 'hidden', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        {product.images[0]
          ? <img src={product.images[0]} alt={product.name} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
          : <span style={{ fontSize: '5rem' }}>📦</span>
        }
      </div>

      <div>
        <h1 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '0.75rem' }}>{product.name}</h1>
        <div style={{ fontSize: '2rem', fontWeight: 700, color: '#2563eb', marginBottom: '0.25rem' }}>${price}</div>
        {product.sale_price && <div style={{ textDecoration: 'line-through', color: '#9ca3af', marginBottom: '1rem' }}>${product.base_price}</div>}

        <p style={{ color: '#4b5563', lineHeight: 1.6, marginBottom: '1.5rem' }}>{product.short_description}</p>

        <button className="btn-primary" style={{ width: '100%', padding: '0.875rem', fontSize: '1rem' }} onClick={() => addItem(product)}>
          {t('products.addToCart')}
        </button>

        {Object.keys(product.specifications).length > 0 && (
          <div style={{ marginTop: '2rem' }}>
            <h3 style={{ fontWeight: 600, marginBottom: '0.75rem' }}>{t('products.specifications')}</h3>
            <table style={{ width: '100%', fontSize: '0.875rem', borderCollapse: 'collapse' }}>
              <tbody>
                {Object.entries(product.specifications).map(([k, v]) => (
                  <tr key={k} style={{ borderBottom: '1px solid #f3f4f6' }}>
                    <td style={{ padding: '0.5rem 0', color: '#6b7280', width: '50%' }}>{k}</td>
                    <td style={{ padding: '0.5rem 0', color: '#1f2937' }}>{v}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
