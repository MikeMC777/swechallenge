<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useStocksStore, useQueryPusher } from '@/stores/stocks'
import DataTable from '@/components/DataTable.vue'
import TopPicks from '@/components/TopPicks.vue'
import RecomputeButton from '@/components/RecomputeButton.vue'

const store = useStocksStore()
const route = useRoute()
const { push } = useQueryPusher()

function parseQuery() {
  const q = (route.query.q as string) || ''
  const sort = (route.query.sort as string) || 'score'
  const dir = ((route.query.dir as string) === 'ASC' ? 'ASC' : 'DESC') as 'ASC' | 'DESC'
  const page = Math.max(parseInt((route.query.page as string) || '1', 10), 1)
  const limit = Math.min(Math.max(parseInt((route.query.limit as string) || '20', 10), 5), 100)
  return { q, sort, dir, page, limit }
}

onMounted(async () => {
  store.hydrate(parseQuery())
  await store.fetch()
})

watch(() => route.query, async () => {
  store.hydrate(parseQuery())
  await store.fetch()
})

// Debounce búsqueda
let t: number | undefined
function setSearchDebounced(v: string) {
  clearTimeout(t)
  // @ts-ignore
  t = setTimeout(() => push({ q: v || undefined, page: 1 }), 300)
}

function nextPage(){ push({ page: store.page + 1 }) }
function prevPage(){ push({ page: Math.max(store.page - 1, 1) }) }
function setLimit(l:number){ push({ limit: l, page: 1 }) }
</script>

<template>
  <div class="space-y-6">
    <!-- Top picks -->
    <TopPicks />

    <!-- Toolbar: búsqueda, page size y acciones -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
      <div class="flex items-center gap-3">
        <input
          class="px-3 py-2 rounded-xl bg-zinc-900 border border-zinc-800 w-80"
          type="text" :value="store.q" placeholder="Search ticker/company/brokerage"
          @input="setSearchDebounced(($event.target as HTMLInputElement).value)"
        />
        <select class="px-3 py-2 rounded-xl bg-zinc-900 border border-zinc-800"
                :value="store.limit" @change="setLimit(parseInt(($event.target as HTMLSelectElement).value))">
          <option :value="10">10</option>
          <option :value="20">20</option>
          <option :value="50">50</option>
        </select>
      </div>

      <div class="sm:ml-auto flex items-center gap-2">
        <button class="px-3 py-2 rounded-lg border border-zinc-700 disabled:opacity-60"
                :disabled="store.page<=1" @click="prevPage">Prev</button>
        <span>Page {{ store.page }}</span>
        <button class="px-3 py-2 rounded-lg border border-zinc-700" @click="nextPage">Next</button>
        <RecomputeButton />
      </div>
    </div>

    <!-- Tabla -->
    <DataTable :items="store.items" />

    <!-- Loading / Empty state -->
    <div v-if="store.loading" class="text-sm opacity-70">Loading…</div>
    <div v-else-if="store.items.length===0" class="text-sm opacity-70">No results. Try a different search.</div>
  </div>
</template>
