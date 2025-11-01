import { defineStore } from 'pinia'
import { listTickers, type TickerRow } from '@/lib/api'

// 🔒 Store sin Vue Router (no useRoute/useRouter aquí)
export const useStocksStore = defineStore('stocks', {
  state: () => ({
    items: [] as TickerRow[],
    q: '',
    sort: 'score' as string,
    dir: 'DESC' as 'ASC' | 'DESC',
    page: 1,
    limit: 20,
    loading: false,
  }),
  actions: {
    // El componente nos pasa los valores del query ya parseados
    hydrate(params: {
      q?: string
      sort?: string
      dir?: 'ASC' | 'DESC'
      page?: number
      limit?: number
    }) {
      if (params.q !== undefined) this.q = params.q
      if (params.sort) this.sort = params.sort
      if (params.dir) this.dir = params.dir
      if (params.page) this.page = Math.max(params.page, 1)
      if (params.limit) this.limit = Math.min(Math.max(params.limit, 5), 100)
    },

    async fetch() {
      this.loading = true
      const offset = (this.page - 1) * this.limit
      try {
        const res = await listTickers(this.q, this.sort, this.dir, offset, this.limit)
        this.items = res.items
      } finally {
        this.loading = false
      }
    },
  },
})

// Helper que puede usar Router porque lo llamamos desde <script setup>
import { useRoute, useRouter } from 'vue-router'
export function useQueryPusher() {
  const route = useRoute()
  const router = useRouter()
  function push(patch: Record<string, any>) {
    const q: Record<string, any> = { ...route.query, ...patch }
    Object.keys(q).forEach(k => {
      if (q[k] === undefined || q[k] === '') delete q[k]
    })
    router.push({ query: q })
  }
  return { push }
}
