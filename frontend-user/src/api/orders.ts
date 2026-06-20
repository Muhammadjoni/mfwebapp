import { apiClient } from './client'
import type { Order, PaginatedResponse } from '../types'

export const ordersApi = {
  list: () => apiClient.get<PaginatedResponse<Order>>('/orders').then(r => r.data),
  get: (id: string) => apiClient.get<{ data: Order }>(`/orders/${id}`).then(r => r.data.data),
  create: (payload: { items: { product_id: string; quantity: number }[]; address: object }) =>
    apiClient.post<{ data: Order }>('/orders', payload).then(r => r.data.data),
}
