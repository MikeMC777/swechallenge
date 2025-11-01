<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
const props = defineProps<{ field: string; label: string }>();
const route = useRoute();
const router = useRouter();

function toggle(){
  const active = route.query.sort === props.field;
  const dir = (route.query.dir as string) || 'DESC';
  router.push({ query: { ...route.query, sort: props.field, dir: active && dir==='DESC' ? 'ASC' : 'DESC' } });
}
</script>
<template>
  <button @click="toggle" class="text-left" :class="{ 'underline': $route.query.sort===field }">
    {{ label }}<span v-if="$route.query.sort===field"> {{ ($route.query.dir||'DESC')==='DESC' ? '↓' : '↑' }}</span>
  </button>
</template>