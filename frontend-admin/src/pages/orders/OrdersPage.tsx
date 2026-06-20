import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { adminApi } from '../../api/admin'
import type { OrderStatus } from '../../types'

const STATUS_COLOR: Record<OrderStatus, string> = {
  created: '#6b7280', paid: '#2563eb', processing: '#d97706',
  shipped: '#7c3aed', delivered: '#059669', closed: '#374151',
  cancelled: '#ef4444', refunded: '#f59e0b',
}

const NEXT_STATUS: Partial<Record<OrderStatus, OrderStatus>> = {
  paid: 'processing', processing: 'shipped', shipped: 'delivered', delivered: 'closed',
}

export default function AdminOrdersPage() {
  const qc = useQueryClient()
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['admin-orders', page, statusFilter],
    queryFn: () => adminApi.orders.list({ page, limit: 20, status: statusFilter || undefined }),
  })

  const updateStatus = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) => adminApi.orders.updateStatus(id, status),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-orders'] }),
  })

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
        <h1 style={{ fontSize: '1.25rem', fontWeight: 700, color: '#1e293b' }}>Orders</h1>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
          style={{ padding: '0.375rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '0.5rem', fontSize: '0.875rem' }}>
          <option value="">All statuses</option>
          {['created','paid','processing','shipped','delivered','closed','cancelled','refunded'].map(s => (
            <option key={s} value={s}>{s}</option>
          ))}
        </select>
      </div>

      {isLoading ? <p style={{ color: '#64748b' }}>Loading...</p> : (
        <div style={{ background: '#fff', borderRadius: '0.75rem', boxShadow: '0 1px 2px rgba(0,0,0,0.05)', overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #e2e8f0', background: '#f8fafc' }}>
                {['Order ID', 'Customer', 'Total', 'Status', 'Date', 'Actions'].map(h => (
                  <th key={h} style={{ textAlign: 'left', padding: '0.75rem 1rem', fontWeight: 600, color: '#475569' }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data?.data.map(order => (
                <tr key={order.id} style={{ borderBottom: '1px solid #f1f5f9' }}>
                  <td style={{ padding: '0.75rem 1rem', fontFamily: 'monospace', color: '#475569' }}>#{order.id.slice(0, 8).toUpperCase()}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#475569' }}>{order.user_id.slice(0, 8)}</td>
                  <td style={{ padding: '0.75rem 1rem', fontWeight: 600, color: '#1e293b' }}>${order.total.toFixed(2)}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{ color: STATUS_COLOR[order.status], fontWeight: 600, fontSize: '0.75rem', background: STATUS_COLOR[order.status] + '1a', padding: '0.125rem 0.5rem', borderRadius: '999px', textTransform: 'capitalize' }}>{order.status}</span>
                  </td>
                  <td style={{ padding: '0.75rem 1rem', color: '#94a3b8' }}>{new Date(order.created_at).toLocaleDateString()}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    {NEXT_STATUS[order.status] && (
                      <button onClick={() => updateStatus.mutate({ id: order.id, status: NEXT_STATUS[order.status]! })}
                        style={{ background: 'none', border: '1px solid #2563eb', color: '#2563eb', borderRadius: '0.375rem', padding: '0.25rem 0.625rem', fontSize: '0.75rem', cursor: 'pointer' }}>
                        → {NEXT_STATUS[order.status]}
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {data && data.meta.total_pages > 1 && (
            <div style={{ padding: '0.75rem 1rem', display: 'flex', gap: '0.5rem', justifyContent: 'flex-end', fontSize: '0.875rem' }}>
              <button onClick={() => setPage(p => p - 1)} disabled={page <= 1}
                style={{ padding: '0.375rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '0.375rem', cursor: 'pointer', background: '#fff' }}>←</button>
              <span style={{ padding: '0.375rem 0.5rem', color: '#64748b' }}>{page} / {data.meta.total_pages}</span>
              <button onClick={() => setPage(p => p + 1)} disabled={page >= data.meta.total_pages}
                style={{ padding: '0.375rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '0.375rem', cursor: 'pointer', background: '#fff' }}>→</button>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
