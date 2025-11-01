<script setup lang="ts">
import { ref } from 'vue'

const API = import.meta.env.VITE_API_URL || 'http://localhost:8081'

const busy = ref(false)
const msg = ref<string>('')

async function run() {
  busy.value = true
  msg.value = ''
  try {
    const res = await fetch(`${API}/api/score/recompute`, { method: 'POST' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    msg.value = 'Scores recomputed'
  } catch (e: any) {
    msg.value = e?.message || 'Error'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="flex items-center gap-2">
    <button
      class="px-3 py-2 rounded-lg border border-zinc-700 bg-zinc-900 hover:bg-zinc-800 disabled:opacity-60"
      :disabled="busy"
      @click="run"
    >
      {{ busy ? 'Working…' : 'Recompute scores' }}
    </button>
    <span class="text-sm opacity-70">{{ msg }}</span>
  </div>
</template>
