import { apiClient } from './client'
import type { Product, PaginatedResponse } from '../types'

export interface ProductFilter {
  page?: number
  limit?: number
  search?: string
  category_id?: string
  min_price?: number
  max_price?: number
  sort_by?: string
  sort_dir?: 'asc' | 'desc'
}

export const productsApi = {
  list: (filter: ProductFilter = {}) =>
    apiClient.get<PaginatedResponse<Product>>('/api/v1/products', { params: filter }).then(r => r.data),

  get: (id: string) =>
    apiClient.get<{ data: Product }>(`/api/v1/products/${id}`).then(r => r.data.data),
}
