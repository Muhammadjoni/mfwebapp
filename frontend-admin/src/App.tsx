import { Routes, Route, Navigate } from 'react-router-dom'
import { useAdminStore } from './store/adminStore'
import AdminLayout from './components/layout/AdminLayout'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/dashboard/DashboardPage'
import UsersPage from './pages/users/UsersPage'
import SellersPage from './pages/sellers/SellersPage'
import AdminProductsPage from './pages/products/ProductsPage'
import AdminOrdersPage from './pages/orders/OrdersPage'
import PaymentsPage from './pages/payments/PaymentsPage'

function AdminRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, user } = useAdminStore()
  if (!isAuthenticated || user?.role !== 'admin') return <Navigate to="/login" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<AdminRoute><AdminLayout /></AdminRoute>}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="users" element={<UsersPage />} />
        <Route path="sellers" element={<SellersPage />} />
        <Route path="products" element={<AdminProductsPage />} />
        <Route path="orders" element={<AdminOrdersPage />} />
        <Route path="payments" element={<PaymentsPage />} />
      </Route>
    </Routes>
  )
}
