import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/slices/authStore'
import MainLayout from './components/layout/MainLayout'
import HomePage from './pages/HomePage'
import ProductListPage from './pages/products/ProductListPage'
import ProductDetailPage from './pages/products/ProductDetailPage'
import CartPage from './pages/cart/CartPage'
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import OrdersPage from './pages/orders/OrdersPage'
import AccountPage from './pages/account/AccountPage'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore(s => s.isAuthenticated)
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<MainLayout />}>
        <Route index element={<HomePage />} />
        <Route path="products" element={<ProductListPage />} />
        <Route path="products/:id" element={<ProductDetailPage />} />
        <Route path="cart" element={<CartPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="orders" element={<PrivateRoute><OrdersPage /></PrivateRoute>} />
        <Route path="account" element={<PrivateRoute><AccountPage /></PrivateRoute>} />
      </Route>
    </Routes>
  )
}
