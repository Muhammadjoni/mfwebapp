import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { ordersApi } from '../../api/orders'
import type { OrderStatus } from '../../types'

const STATUS_COLOR: Record<OrderStatus, string> = {
  created: '#6b7280',
  paid: '#2563eb',
  processing: '#d97706',
  shipped: '#7c3aed',
  delivered: '#059669',
  closed: '#374151',
  cancelled: '#ef4444',
  refunded: '#f59e0b',
}

export default function OrdersPage() {
  const { t } = useTranslation()
  const { data, isLoading } = useQuery({ queryKey: ['orders'], queryFn: ordersApi.list })

  if (isLoading) return <p style={{ textAlign: 'center', padding: '3rem', color: '#6b7280' }}>{t('common.loading')}</p>

  if (!data?.data?.length) {
    return <p style={{ textAlign: 'center', padding: '3rem', color: '#9ca3af' }}>{t('orders.empty')}</p>
  }

  return (
    <div>
      <h2 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.25rem' }}>{t('nav.orders')}</h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.875rem' }}>
        {data.data.map(order => (
          <div key={order.id} className="card" style={{ display: 'grid', gridTemplateColumns: '1fr auto auto', gap: '1rem', alignItems: 'center' }}>
            <div>
              <p style={{ fontWeight: 600, fontSize: '0.875rem', marginBottom: '0.25rem' }}>#{order.id.slice(0, 8).toUpperCase()}</p>
              <p style={{ color: '#6b7280', fontSize: '0.75rem' }}>{new Date(order.created_at).toLocaleDateString()}</p>
              <p style={{ color: '#4b5563', fontSize: '0.75rem', marginTop: '0.25rem' }}>
                {order.items.length} {t('orders.items')}
              </p>
            </div>
            <span style={{ fontSize: '0.75rem', fontWeight: 600, color: STATUS_COLOR[order.status], background: STATUS_COLOR[order.status] + '1a', padding: '0.25rem 0.625rem', borderRadius: '999px' }}>
              {t(`orders.status.${order.status}`)}
            </span>
            <span style={{ fontWeight: 700, color: '#1f2937' }}>${order.total.toFixed(2)}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
