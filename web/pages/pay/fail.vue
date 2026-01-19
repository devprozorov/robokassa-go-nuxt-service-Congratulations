<template>
  <div class="min-h-screen flex items-center justify-center p-6">
    <div class="max-w-md w-full rounded-3xl border border-white/10 bg-white/5 p-6 text-center">
      <div class="text-2xl font-bold text-white">Оплата не завершена ❌</div>
      <div class="mt-2 text-white/70">Можно попробовать ещё раз из кабинета.</div>
      <NuxtLink class="mt-6 inline-block rounded-xl bg-white/90 px-4 py-3 font-semibold text-slate-900" to="/dashboard">
        В кабинет
      </NuxtLink>
    </div>
  </div>
</template>
