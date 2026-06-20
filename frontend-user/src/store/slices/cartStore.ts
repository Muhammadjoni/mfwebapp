import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { CartItem, Product } from '../../types'

interface CartState {
  items: CartItem[]
  addItem: (product: Product, quantity?: number) => void
  removeItem: (productId: string) => void
  updateQuantity: (productId: string, quantity: number) => void
  clear: () => void
  totalItems: () => number
  subtotal: () => number
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],

      addItem: (product, quantity = 1) =>
        set(state => {
          const existing = state.items.find(i => i.product_id === product.id)
          if (existing) {
            return { items: state.items.map(i => i.product_id === product.id ? { ...i, quantity: i.quantity + quantity } : i) }
          }
          return { items: [...state.items, { id: crypto.randomUUID(), product_id: product.id, variant_id: null, quantity, product }] }
        }),

      removeItem: productId =>
        set(state => ({ items: state.items.filter(i => i.product_id !== productId) })),

      updateQuantity: (productId, quantity) => {
        if (quantity <= 0) { get().removeItem(productId); return }
        set(state => ({ items: state.items.map(i => i.product_id === productId ? { ...i, quantity } : i) }))
      },

      clear: () => set({ items: [] }),

      totalItems: () => get().items.reduce((sum, i) => sum + i.quantity, 0),

      subtotal: () =>
        get().items.reduce((sum, i) => sum + (i.product?.sale_price ?? i.product?.base_price ?? 0) * i.quantity, 0),
    }),
    { name: 'mf-cart' }
  )
)
