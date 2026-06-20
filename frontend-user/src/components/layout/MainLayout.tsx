import { Outlet } from 'react-router-dom'
import Navbar from './Navbar'

export default function MainLayout() {
  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Navbar />
      <main style={{ flex: 1, maxWidth: 1200, margin: '0 auto', width: '100%', padding: '1.5rem 1rem' }}>
        <Outlet />
      </main>
      <footer style={{ background: '#fff', borderTop: '1px solid #e5e7eb', padding: '1.5rem', textAlign: 'center', fontSize: '0.875rem', color: '#9ca3af' }}>
        © {new Date().getFullYear()} MF Shop. All rights reserved.
      </footer>
    </div>
  )
}
