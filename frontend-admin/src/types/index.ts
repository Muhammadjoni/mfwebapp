export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  phone: string
  role: 'admin' | 'user' | 'seller'
  status: 'active' | 'inactive' | 'banned' | 'pending'
  language: 'ru' | 'tj'
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface Seller {
  id: string
  user_id: string
  business_name: string
  business_email: string
  status: 'pending' | 'active' | 'suspended' | 'rejected'
  commission_rate: number
  total_sales: number
  verified_at: string | null
  created_at: string
}

export type ProductStatus = 'pending' | 'active' | 'rejected' | 'archived'

export interface Product {
  id: string
  seller_id: string
  name: string
  slug: string
  base_price: number
  sale_price: number | null
  stock: number
  status: ProductStatus
  images: string[]
  created_at: string
}

export type OrderStatus =
  | 'created' | 'paid' | 'processing' | 'shipped'
  | 'delivered' | 'closed' | 'cancelled' | 'refunded'

export interface Order {
  id: string
  user_id: string
  status: OrderStatus
  total: number
  currency: string
  created_at: string
  updated_at: string
}

export type PaymentStatus = 'pending' | 'completed' | 'failed' | 'refunded'

export interface Payment {
  id: string
  order_id: string
  amount: number
  currency: string
  provider: 'stripe' | 'visa' | 'alif'
  status: PaymentStatus
  idempotency_key: string
  created_at: string
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: { code: string; message: string }
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  meta: {
    page: number
    limit: number
    total_items: number
    total_pages: number
  }
}
