import { useQuery } from '@tanstack/react-query'
import { adminApi } from '../../api/admin'

interface StatCardProps { label: string; value: number | string; color: string }
function StatCard({ label, value, color }: StatCardProps) {
  return (
    <div style={{ background: '#fff', borderRadius: '0.75rem', padding: '1.25rem 1.5rem', boxShadow: '0 1px 2px rgba(0,0,0,0.05)', borderLeft: `4px solid ${color}` }}>
      <p style={{ color: '#64748b', fontSize: '0.875rem', marginBottom: '0.5rem' }}>{label}</p>
      <p style={{ fontSize: '1.75rem', fontWeight: 700, color: '#1e293b' }}>{value}</p>
    </div>
  )
}

export default function DashboardPage() {
  const { data: users } = useQuery({ queryKey: ['admin-users-count'], queryFn: () => adminApi.users.list({ limit: 1 }) })
  const { data: orders } = useQuery({ queryKey: ['admin-orders-count'], queryFn: () => adminApi.orders.list({ limit: 1 }) })
  const { data: products } = useQuery({ queryKey: ['admin-products-count'], queryFn: () => adminApi.products.list({ limit: 1 }) })
  const { data: pendingProducts } = useQuery({ queryKey: ['admin-pending-products'], queryFn: () => adminApi.products.list({ status: 'pending', limit: 1 }) })

  return (
    <div>
      <h1 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem', color: '#1e293b' }}>Dashboard</h1>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '1rem', marginBottom: '2rem' }}>
        <StatCard label="Total Users" value={users?.meta.total_items ?? '—'} color="#3b82f6" />
        <StatCard label="Total Orders" value={orders?.meta.total_items ?? '—'} color="#10b981" />
        <StatCard label="Total Products" value={products?.meta.total_items ?? '—'} color="#f59e0b" />
        <StatCard label="Pending Approval" value={pendingProducts?.meta.total_items ?? '—'} color="#ef4444" />
      </div>

      {(pendingProducts?.meta.total_items ?? 0) > 0 && (
        <div style={{ background: '#fffbeb', border: '1px solid #fcd34d', borderRadius: '0.75rem', padding: '1rem 1.25rem', fontSize: '0.875rem', color: '#92400e' }}>
          <strong>{pendingProducts?.meta.total_items}</strong> product(s) are awaiting approval.{' '}
          <a href="/products?status=pending" style={{ color: '#2563eb', fontWeight: 500 }}>Review now →</a>
        </div>
      )}
    </div>
  )
}
