<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { recommend } from '@/lib/api'

type Rec = { symbol: string; name?: string; score: number; rationale: string }
const items = ref<Rec[]>([])
const loading = ref(true)
const error = ref<string|null>(null)

onMounted(async () => {
  try {
    const res = await recommend()
    items.value = res.items
  } catch (e:any) {
    error.value = e?.message || 'error'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="space-y-2">
    <h2 class="text-lg font-semibold">Top picks</h2>
    <div v-if="loading" class="text-sm opacity-70">Loading…</div>
    <div v-else-if="error" class="text-sm text-red-400">Error: {{ error }}</div>
    <div v-else class="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
      <div v-for="r in items" :key="r.symbol"
           class="rounded-2xl border border-zinc-800 p-4 bg-zinc-900/40">
        <div class="flex items-center justify-between">
          <div class="text-base font-semibold">{{ r.symbol }}</div>
          <div class="text-sm px-2 py-0.5 rounded-lg border border-zinc-700">
            {{ r.score.toFixed(2) }}
          </div>
        </div>
        <div class="text-sm opacity-80 truncate">{{ r.name }}</div>
        <div class="text-xs opacity-70 mt-2 line-clamp-3">{{ r.rationale }}</div>
      </div>
    </div>
  </div>
</template>
