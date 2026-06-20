import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { adminApi } from '../../api/admin'

const STATUS_COLOR: Record<string, string> = { active: '#059669', pending: '#d97706', suspended: '#ef4444', rejected: '#6b7280' }

export default function SellersPage() {
  const qc = useQueryClient()
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['admin-sellers', page, statusFilter],
    queryFn: () => adminApi.sellers.list({ page, limit: 20, status: statusFilter || undefined }),
  })

  const approve = useMutation({
    mutationFn: (id: string) => adminApi.sellers.approve(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-sellers'] }),
  })

  const reject = useMutation({
    mutationFn: (id: string) => adminApi.sellers.reject(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-sellers'] }),
  })

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
        <h1 style={{ fontSize: '1.25rem', fontWeight: 700, color: '#1e293b' }}>Sellers</h1>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
          style={{ padding: '0.375rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '0.5rem', fontSize: '0.875rem' }}>
          <option value="">All statuses</option>
          <option value="pending">Pending</option>
          <option value="active">Active</option>
          <option value="suspended">Suspended</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>

      {isLoading ? <p style={{ color: '#64748b' }}>Loading...</p> : (
        <div style={{ background: '#fff', borderRadius: '0.75rem', boxShadow: '0 1px 2px rgba(0,0,0,0.05)', overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #e2e8f0', background: '#f8fafc' }}>
                {['Business', 'Email', 'Commission', 'Status', 'Joined', 'Actions'].map(h => (
                  <th key={h} style={{ textAlign: 'left', padding: '0.75rem 1rem', fontWeight: 600, color: '#475569' }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data?.data.map(seller => (
                <tr key={seller.id} style={{ borderBottom: '1px solid #f1f5f9' }}>
                  <td style={{ padding: '0.75rem 1rem', fontWeight: 500, color: '#1e293b' }}>{seller.business_name}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#475569' }}>{seller.business_email}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#475569' }}>{(seller.commission_rate * 100).toFixed(1)}%</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{ color: STATUS_COLOR[seller.status], fontWeight: 600, fontSize: '0.75rem', textTransform: 'capitalize' }}>{seller.status}</span>
                  </td>
                  <td style={{ padding: '0.75rem 1rem', color: '#94a3b8' }}>{new Date(seller.created_at).toLocaleDateString()}</td>
                  <td style={{ padding: '0.75rem 1rem', display: 'flex', gap: '0.5rem' }}>
                    {seller.status === 'pending' && (
                      <>
                        <button onClick={() => approve.mutate(seller.id)}
                          style={{ background: 'none', border: '1px solid #059669', color: '#059669', borderRadius: '0.375rem', padding: '0.25rem 0.625rem', fontSize: '0.75rem', cursor: 'pointer' }}>Approve</button>
                        <button onClick={() => reject.mutate(seller.id)}
                          style={{ background: 'none', border: '1px solid #ef4444', color: '#ef4444', borderRadius: '0.375rem', padding: '0.25rem 0.625rem', fontSize: '0.75rem', cursor: 'pointer' }}>Reject</button>
                      </>
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
