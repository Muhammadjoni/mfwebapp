export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  phone: string
  role: 'admin' | 'user' | 'seller'
  status: 'active' | 'inactive' | 'banned' | 'pending'
  language: 'ru' | 'tj'
  avatar_url: string
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface Product {
  id: string
  seller_id: string
  category_id: string
  name: string
  slug: string
  description: string
  short_description: string
  base_price: number
  sale_price: number | null
  currency: string
  sku: string
  stock: number
  status: string
  images: string[]
  tags: string[]
  specifications: Record<string, string>
  rating: number
  review_count: number
  view_count: number
  sold_count: number
  created_at: string
  updated_at: string
}

export interface Category {
  id: string
  parent_id: string | null
  name: string
  slug: string
  description: string
  image_url: string
  is_active: boolean
}

export interface CartItem {
  id: string
  product_id: string
  variant_id: string | null
  quantity: number
  product: Product
}

export type OrderStatus =
  | 'created' | 'paid' | 'processing' | 'shipped'
  | 'delivered' | 'closed' | 'cancelled' | 'refunded'

export interface OrderItem {
  id: string
  product_id: string
  name: string
  sku: string
  image_url: string
  quantity: number
  unit_price: number
  total_price: number
}

export interface Address {
  full_name: string
  phone: string
  country: string
  city: string
  state: string
  street: string
  postal_code: string
  apartment: string
}

export interface Order {
  id: string
  user_id: string
  status: OrderStatus
  items: OrderItem[]
  subtotal: number
  shipping_cost: number
  tax: number
  total: number
  currency: string
  shipping_address: Address
  tracking_number: string
  created_at: string
  updated_at: string
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
