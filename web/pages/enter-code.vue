<template>
  <div class="max-w-md mx-auto">
    <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
      <h1 class="text-2xl font-bold">Открыть поздравление по коду</h1>
      <p class="text-sm text-slate-300 mt-1">Введи код и открой страницу</p>

      <form class="mt-5 space-y-3" @submit.prevent="go">
        <input v-model="code" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10 font-mono" placeholder="ABCDEF" />
        <button class="w-full px-4 py-2 rounded-xl bg-white text-slate-950 font-semibold hover:bg-slate-100">
          Открыть
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const code = ref('')
async function go() {
  const c = code.value.trim().toUpperCase()
  if (!c) return
  await navigateTo(`/c/${c}`)
}
</script>
