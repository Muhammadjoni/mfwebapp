import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { adminApi } from '../../api/admin'
import type { PaymentStatus } from '../../types'

const STATUS_COLOR: Record<PaymentStatus, string> = {
  pending: '#d97706', completed: '#059669', failed: '#ef4444', refunded: '#f59e0b',
}

export default function PaymentsPage() {
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['admin-payments', page, statusFilter],
    queryFn: () => adminApi.payments.list({ page, limit: 20, status: statusFilter || undefined }),
  })

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
        <h1 style={{ fontSize: '1.25rem', fontWeight: 700, color: '#1e293b' }}>Payments</h1>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
          style={{ padding: '0.375rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '0.5rem', fontSize: '0.875rem' }}>
          <option value="">All statuses</option>
          <option value="pending">Pending</option>
          <option value="completed">Completed</option>
          <option value="failed">Failed</option>
          <option value="refunded">Refunded</option>
        </select>
      </div>

      {isLoading ? <p style={{ color: '#64748b' }}>Loading...</p> : (
        <div style={{ background: '#fff', borderRadius: '0.75rem', boxShadow: '0 1px 2px rgba(0,0,0,0.05)', overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #e2e8f0', background: '#f8fafc' }}>
                {['Payment ID', 'Order ID', 'Amount', 'Provider', 'Status', 'Date'].map(h => (
                  <th key={h} style={{ textAlign: 'left', padding: '0.75rem 1rem', fontWeight: 600, color: '#475569' }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data?.data.map(payment => (
                <tr key={payment.id} style={{ borderBottom: '1px solid #f1f5f9' }}>
                  <td style={{ padding: '0.75rem 1rem', fontFamily: 'monospace', color: '#475569' }}>{payment.id.slice(0, 8)}</td>
                  <td style={{ padding: '0.75rem 1rem', fontFamily: 'monospace', color: '#475569' }}>{payment.order_id.slice(0, 8)}</td>
                  <td style={{ padding: '0.75rem 1rem', fontWeight: 600, color: '#1e293b' }}>{payment.currency} {payment.amount.toFixed(2)}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{ background: '#f1f5f9', color: '#475569', padding: '0.125rem 0.5rem', borderRadius: '999px', fontSize: '0.75rem', fontWeight: 500, textTransform: 'capitalize' }}>{payment.provider}</span>
                  </td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{ color: STATUS_COLOR[payment.status], fontWeight: 600, fontSize: '0.75rem', textTransform: 'capitalize' }}>{payment.status}</span>
                  </td>
                  <td style={{ padding: '0.75rem 1rem', color: '#94a3b8' }}>{new Date(payment.created_at).toLocaleDateString()}</td>
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
