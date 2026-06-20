import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { adminApi } from '../../api/admin'
import type { User } from '../../types'

const STATUS_COLOR: Record<string, string> = { active: '#059669', inactive: '#6b7280', banned: '#ef4444', pending: '#d97706' }

export default function UsersPage() {
  const qc = useQueryClient()
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['admin-users', page, statusFilter],
    queryFn: () => adminApi.users.list({ page, limit: 20, status: statusFilter || undefined }),
  })

  const updateStatus = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) => adminApi.users.updateStatus(id, status),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-users'] }),
  })

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1.25rem' }}>
        <h1 style={{ fontSize: '1.25rem', fontWeight: 700, color: '#1e293b' }}>Users</h1>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
          style={{ padding: '0.375rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '0.5rem', fontSize: '0.875rem' }}>
          <option value="">All statuses</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
          <option value="banned">Banned</option>
          <option value="pending">Pending</option>
        </select>
      </div>

      {isLoading ? <p style={{ color: '#64748b' }}>Loading...</p> : (
        <div style={{ background: '#fff', borderRadius: '0.75rem', boxShadow: '0 1px 2px rgba(0,0,0,0.05)', overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #e2e8f0', background: '#f8fafc' }}>
                {['Name', 'Email', 'Role', 'Status', 'Joined', 'Actions'].map(h => (
                  <th key={h} style={{ textAlign: 'left', padding: '0.75rem 1rem', fontWeight: 600, color: '#475569' }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data?.data.map((user: User) => (
                <tr key={user.id} style={{ borderBottom: '1px solid #f1f5f9' }}>
                  <td style={{ padding: '0.75rem 1rem', color: '#1e293b', fontWeight: 500 }}>{user.first_name} {user.last_name}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#475569' }}>{user.email}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{ background: '#dbeafe', color: '#1d4ed8', padding: '0.125rem 0.5rem', borderRadius: '999px', fontSize: '0.75rem', fontWeight: 500, textTransform: 'capitalize' }}>{user.role}</span>
                  </td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{ color: STATUS_COLOR[user.status], fontWeight: 600, fontSize: '0.75rem', textTransform: 'capitalize' }}>{user.status}</span>
                  </td>
                  <td style={{ padding: '0.75rem 1rem', color: '#94a3b8' }}>{new Date(user.created_at).toLocaleDateString()}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    {user.status !== 'banned'
                      ? <button onClick={() => updateStatus.mutate({ id: user.id, status: 'banned' })}
                          style={{ background: 'none', border: '1px solid #ef4444', color: '#ef4444', borderRadius: '0.375rem', padding: '0.25rem 0.625rem', fontSize: '0.75rem', cursor: 'pointer' }}>Ban</button>
                      : <button onClick={() => updateStatus.mutate({ id: user.id, status: 'active' })}
                          style={{ background: 'none', border: '1px solid #059669', color: '#059669', borderRadius: '0.375rem', padding: '0.25rem 0.625rem', fontSize: '0.75rem', cursor: 'pointer' }}>Unban</button>
                    }
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
