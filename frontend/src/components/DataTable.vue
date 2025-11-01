<script setup lang="ts">
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { TickerRow } from '@/lib/api'
const props = defineProps<{ items: TickerRow[] }>()
const route = useRoute()
const router = useRouter()
function setSort(field: string) {
  const active = route.query.sort === field
  const dir = (route.query.dir as string) || 'DESC'
  router.push({ query: { ...route.query, sort: field, dir: active && dir === 'DESC' ? 'ASC' : 'DESC', page: 1 } })
}
</script>

<template>
  <div class="overflow-hidden rounded-2xl border border-zinc-800">
    <table class="w-full text-sm">
      <thead class="bg-zinc-900">
        <tr class="grid grid-cols-12 gap-2 px-3 py-2">
          <th class="col-span-2"><button @click="setSort('ticker')" class="text-left" :class="{ underline: $route.query.sort==='ticker' }">Ticker <span v-if="$route.query.sort==='ticker'">{{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span></button></th>
          <th class="col-span-3"><button @click="setSort('company')" class="text-left" :class="{ underline: $route.query.sort==='company' }">Company <span v-if="$route.query.sort==='company'">{{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span></button></th>
          <th class="col-span-2"><button @click="setSort('last_action')" class="text-left" :class="{ underline: $route.query.sort==='last_action' }">Last Action <span v-if="$route.query.sort==='last_action'">{{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span></button></th>
          <th class="col-span-2"><button @click="setSort('last_rating')" class="text-left" :class="{ underline: $route.query.sort==='last_rating' }">Rating <span v-if="$route.query.sort==='last_rating'">{{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span></button></th>
          <th class="col-span-1 text-right"><button @click="setSort('last_target')" class="text-right w-full" :class="{ underline: $route.query.sort==='last_target' }">Target <span v-if="$route.query.sort==='last_target'">{{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span></button></th>
          <th class="col-span-2"><button @click="setSort('score')" class="text-left" :class="{ underline: $route.query.sort==='score' }">Score <span v-if="$route.query.sort==='score'">{{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span></button></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in props.items" :key="s.ticker" class="grid grid-cols-12 items-center gap-2 border-t border-zinc-800 px-3 py-2 hover:bg-zinc-900/50">
          <td class="col-span-2 font-medium"><RouterLink :to="`/stocks/${s.ticker}`">{{ s.ticker }}</RouterLink></td>
          <td class="col-span-3 opacity-80 truncate">{{ s.company }}</td>
          <td class="col-span-2">{{ s.last_action }}<span v-if="s.last_brokerage" class="opacity-60"> • {{ s.last_brokerage }}</span></td>
          <td class="col-span-2">{{ s.last_rating }}</td>
          <td class="col-span-1 text-right">{{ s.last_target?.toLocaleString?.('en-US',{ style:'currency', currency:'USD' }) }}</td>
          <td class="col-span-2">{{ s.score.toFixed(2) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
