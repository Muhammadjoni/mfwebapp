import { apiClient } from './client'
import type { User, Seller, Product, Order, Payment, PaginatedResponse } from '../types'

export const adminApi = {
  users: {
    list: (params?: { page?: number; limit?: number; role?: string; status?: string }) =>
      apiClient.get<PaginatedResponse<User>>('/admin/users', { params }).then(r => r.data),
    updateStatus: (id: string, status: string) =>
      apiClient.patch<{ data: User }>(`/admin/users/${id}/status`, { status }).then(r => r.data.data),
    updateRole: (id: string, role: string) =>
      apiClient.patch<{ data: User }>(`/admin/users/${id}/role`, { role }).then(r => r.data.data),
  },
  sellers: {
    list: (params?: { page?: number; limit?: number; status?: string }) =>
      apiClient.get<PaginatedResponse<Seller>>('/admin/sellers', { params }).then(r => r.data),
    approve: (id: string) =>
      apiClient.patch(`/admin/sellers/${id}/approve`).then(r => r.data),
    reject: (id: string) =>
      apiClient.patch(`/admin/sellers/${id}/reject`).then(r => r.data),
  },
  products: {
    list: (params?: { page?: number; limit?: number; status?: string }) =>
      apiClient.get<PaginatedResponse<Product>>('/admin/products', { params }).then(r => r.data),
    approve: (id: string) =>
      apiClient.patch(`/seller/products/${id}/approve`).then(r => r.data),
    reject: (id: string, reason?: string) =>
      apiClient.patch(`/seller/products/${id}/reject`, { reason }).then(r => r.data),
  },
  orders: {
    list: (params?: { page?: number; limit?: number; status?: string }) =>
      apiClient.get<PaginatedResponse<Order>>('/admin/orders', { params }).then(r => r.data),
    updateStatus: (id: string, status: string) =>
      apiClient.patch(`/orders/${id}/status`, { status, admin_override: true }).then(r => r.data),
  },
  payments: {
    list: (params?: { page?: number; limit?: number; status?: string }) =>
      apiClient.get<PaginatedResponse<Payment>>('/admin/payments', { params }).then(r => r.data),
  },
}
