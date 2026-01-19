<template>
  <div>
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold">Заказы</h1>
      <NuxtLink to="/admin" class="text-sm text-slate-300 hover:text-white">← назад</NuxtLink>
    </div>

    <div class="mt-6 overflow-x-auto rounded-3xl border border-white/10">
      <table class="min-w-full text-sm">
        <thead class="bg-white/5 text-slate-300">
          <tr>
            <th class="text-left px-4 py-3">ID</th>
            <th class="text-left px-4 py-3">Owner</th>
            <th class="text-left px-4 py-3">Greeting</th>
            <th class="text-left px-4 py-3">Type</th>
            <th class="text-left px-4 py-3">Amount</th>
            <th class="text-left px-4 py-3">Status</th>
            <th class="text-left px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="o in items" :key="o.id" class="border-t border-white/10">
            <td class="px-4 py-3 font-mono text-xs">{{ o.id }}</td>
            <td class="px-4 py-3 font-mono text-xs">{{ o.ownerId }}</td>
            <td class="px-4 py-3 font-mono text-xs">{{ o.greetingId }}</td>
            <td class="px-4 py-3">{{ o.productType }}</td>
            <td class="px-4 py-3">{{ o.amountRUB }} {{ o.currency }}</td>
            <td class="px-4 py-3">
              <span v-if="o.status==='paid'" class="text-emerald-200">paid</span>
              <span v-else class="text-slate-300">{{ o.status }}</span>
            </td>
            <td class="px-4 py-3">
              <button
                v-if="o.status !== 'paid'"
                class="px-3 py-1.5 rounded-xl bg-emerald-400/10 border border-emerald-400/20 hover:bg-emerald-400/15 text-emerald-200"
                @click="markPaid(o.id)"
              >
                Пометить оплачено
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="err" class="mt-4 text-sm text-rose-300">{{ err }}</div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['admin'] })
const api = useApi()
const items = ref<any[]>([])
const err = ref('')

async function load() {
  const res: any = await api.get('/admin/orders')
  items.value = res.items || []
}
onMounted(load)

async function markPaid(id: string) {
  if (!confirm('Подтвердить оплату?')) return
  err.value = ''
  try {
    const res: any = await api.post(`/admin/orders/${id}/mark-paid`)
    if (res.ok) await load()
    else err.value = res.error || 'MARK_FAILED'
  } catch {
    err.value = 'MARK_FAILED'
  }
}
</script>
