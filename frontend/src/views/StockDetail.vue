<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getTicker } from '@/lib/api'
import { useRoute } from 'vue-router'

const route = useRoute()
const data = ref<Awaited<ReturnType<typeof getTicker>> | null>(null)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    data.value = await getTicker(route.params.symbol as string)
  } catch (e:any) {
    error.value = e?.message || 'error'
  }
})
</script>

<template>
  <div class="space-y-4">
    <RouterLink to="/">← Back</RouterLink>
    <div v-if="error" class="text-red-400">Error: {{ error }}</div>
    <div v-else-if="!data">Loading…</div>
    <div v-else>
      <h2 class="text-xl font-semibold">{{ data.ticker }}</h2>
      <div class="overflow-hidden rounded-2xl border border-zinc-800">
        <table class="w-full text-sm">
          <thead class="bg-zinc-900">
            <tr>
              <th class="text-left p-2">Time</th>
              <th class="text-left p-2">Action</th>
              <th class="text-left p-2">Brokerage</th>
              <th class="text-left p-2">Rating</th>
              <th class="text-left p-2">Target</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="e in data.events" :key="e.time" class="border-t border-zinc-800">
              <td class="p-2">{{ new Date(e.time).toLocaleString() }}</td>
              <td class="p-2">{{ e.action }}</td>
              <td class="p-2">{{ e.brokerage }}</td>
              <td class="p-2">{{ e.rating_from }} → {{ e.rating_to }}</td>
              <td class="p-2">
                <span v-if="e.target_from != null">{{ e.target_from.toLocaleString('en-US',{style:'currency', currency:'USD'}) }} → </span>
                <span v-if="e.target_to != null">{{ e.target_to.toLocaleString('en-US',{style:'currency', currency:'USD'}) }}</span>
                <span v-if="e.target_delta != null" class="opacity-60"> (Δ {{ e.target_delta.toFixed(2) }})</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
