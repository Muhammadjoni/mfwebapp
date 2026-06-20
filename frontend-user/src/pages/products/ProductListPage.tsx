import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { productsApi, type ProductFilter } from '../../api/products'
import { useCartStore } from '../../store/slices/cartStore'

export default function ProductListPage() {
  const { t } = useTranslation()
  const addItem = useCartStore(s => s.addItem)
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState<ProductFilter>({ page: 1, limit: 20 })

  const { data, isLoading, isError } = useQuery({
    queryKey: ['products', filter],
    queryFn: () => productsApi.list(filter),
  })

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setFilter(f => ({ ...f, search, page: 1 }))
  }

  if (isLoading) return <p style={{ textAlign: 'center', padding: '3rem', color: '#6b7280' }}>{t('common.loading')}</p>
  if (isError) return <p style={{ textAlign: 'center', padding: '3rem', color: '#ef4444' }}>{t('common.error')}</p>

  return (
    <div>
      <form onSubmit={handleSearch} style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.5rem' }}>
        <input
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder={t('products.search')}
          style={{ flex: 1 }}
        />
        <button type="submit" className="btn-primary">{t('common.search')}</button>
      </form>

      {data?.data.length === 0 && (
        <p style={{ textAlign: 'center', padding: '3rem', color: '#9ca3af' }}>{t('products.noProducts')}</p>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))', gap: '1rem' }}>
        {data?.data.map(product => (
          <div key={product.id} className="card" style={{ overflow: 'hidden' }}>
            <Link to={`/products/${product.id}`}>
              <div style={{ aspectRatio: '1', background: '#f3f4f6', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                {product.images[0]
                  ? <img src={product.images[0]} alt={product.name} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                  : <span style={{ fontSize: '3rem' }}>📦</span>
                }
              </div>
            </Link>
            <div style={{ padding: '0.875rem' }}>
              <p style={{ fontWeight: 500, fontSize: '0.875rem', marginBottom: '0.5rem', overflow: 'hidden', display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical' }}>
                {product.name}
              </p>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <span style={{ fontWeight: 700, color: '#2563eb' }}>${product.sale_price ?? product.base_price}</span>
                <button className="btn-primary" style={{ padding: '0.25rem 0.625rem', fontSize: '0.75rem' }} onClick={() => addItem(product)}>
                  {t('products.addToCart')}
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {data && data.meta.total_pages > 1 && (
        <div style={{ display: 'flex', justifyContent: 'center', gap: '0.5rem', marginTop: '2rem' }}>
          <button onClick={() => setFilter(f => ({ ...f, page: (f.page ?? 1) - 1 }))} disabled={(filter.page ?? 1) <= 1} className="btn-primary">←</button>
          <span style={{ padding: '0.5rem 1rem', color: '#6b7280', fontSize: '0.875rem' }}>{filter.page} / {data.meta.total_pages}</span>
          <button onClick={() => setFilter(f => ({ ...f, page: (f.page ?? 1) + 1 }))} disabled={(filter.page ?? 1) >= data.meta.total_pages} className="btn-primary">→</button>
        </div>
      )}
    </div>
  )
}
