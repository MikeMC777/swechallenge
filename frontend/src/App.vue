<template>
  <div class="min-h-dvh max-w-6xl mx-auto p-6 space-y-6">
    <header class="flex items-center gap-3">
      <h1 class="text-2xl font-semibold">Stocks</h1>
      <nav class="ml-auto flex gap-4 text-sm">
        <RouterLink to="/">Home</RouterLink>
      </nav>
    </header>
    <RouterView />
  </div>
</template>

<script setup lang="ts"></script>
