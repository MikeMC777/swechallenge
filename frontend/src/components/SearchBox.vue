<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
const router = useRouter();
const route = useRoute();
const q = ref<string>((route.query.q as string) || '');
let t: number | undefined;
watch(q, (val)=>{
  clearTimeout(t);
  t = window.setTimeout(()=>{
    const query = { ...route.query } as Record<string, any>;
    if (val) query.q = val; else delete query.q;
    query.page = 1; // reset page on new search
    router.push({ query });
  }, 300);
});
</script>
<template>
  <input v-model="q" placeholder="Search symbol, name, sector…" class="w-full rounded-xl border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-zinc-600" />
</template>