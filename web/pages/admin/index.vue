<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold">Админ-панель</h1>

    <div class="grid md:grid-cols-3 gap-4">
      <NuxtLink to="/admin/users" class="rounded-3xl border border-white/10 bg-white/5 p-6 hover:bg-white/10 transition">
        <div class="text-sm text-slate-300">Пользователи</div>
        <div class="text-xl font-semibold mt-1">Список / удалить</div>
      </NuxtLink>
      <NuxtLink to="/admin/orders" class="rounded-3xl border border-white/10 bg-white/5 p-6 hover:bg-white/10 transition">
        <div class="text-sm text-slate-300">Продажи</div>
        <div class="text-xl font-semibold mt-1">Заказы / оплатить</div>
      </NuxtLink>
      <NuxtLink to="/admin/greetings" class="rounded-3xl border border-white/10 bg-white/5 p-6 hover:bg-white/10 transition">
        <div class="text-sm text-slate-300">Страницы</div>
        <div class="text-xl font-semibold mt-1">Поздравления</div>
      </NuxtLink>
    </div>

    <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
      <div class="font-semibold">Статистика</div>
      <div class="mt-3 grid md:grid-cols-2 gap-3">
        <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
          <div class="text-xs text-slate-300">Оплаченных заказов</div>
          <div class="text-2xl font-bold">{{ stats.paidOrders ?? '...' }}</div>
        </div>
        <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
          <div class="text-xs text-slate-300">Сумма RUB</div>
          <div class="text-2xl font-bold">{{ stats.paidSumRUB ?? '...' }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['admin'] })
const api = useApi()
const stats = ref<any>({})

onMounted(async () => {
  const res: any = await api.get('/admin/stats')
  stats.value = res
})
</script>
