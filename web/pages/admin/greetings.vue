<template>
  <div>
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold">Поздравления</h1>
      <NuxtLink to="/admin" class="text-sm text-slate-300 hover:text-white">← назад</NuxtLink>
    </div>

    <div class="mt-6 overflow-x-auto rounded-3xl border border-white/10">
      <table class="min-w-full text-sm">
        <thead class="bg-white/5 text-slate-300">
          <tr>
            <th class="text-left px-4 py-3">ID</th>
            <th class="text-left px-4 py-3">Owner</th>
            <th class="text-left px-4 py-3">Type</th>
            <th class="text-left px-4 py-3">Code/Sub</th>
            <th class="text-left px-4 py-3">Paid</th>
            <th class="text-left px-4 py-3">Expires</th>
            <th class="text-left px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="g in items" :key="g.id" class="border-t border-white/10">
            <td class="px-4 py-3 font-mono text-xs">{{ g.id }}</td>
            <td class="px-4 py-3 font-mono text-xs">{{ g.ownerId }}</td>
            <td class="px-4 py-3">{{ g.type }}</td>
            <td class="px-4 py-3 font-mono text-xs">{{ g.type==='code' ? g.code : g.subdomain }}</td>
            <td class="px-4 py-3">
              <span v-if="g.paid" class="text-emerald-200">yes</span>
              <span v-else class="text-slate-300">no</span>
            </td>
            <td class="px-4 py-3 text-xs text-slate-300">
              {{ g.expiresAt ? new Date(g.expiresAt).toLocaleString() : '-' }}
            </td>
            <td class="px-4 py-3">
              <button class="px-3 py-1.5 rounded-xl bg-rose-500/10 border border-rose-500/20 hover:bg-rose-500/15 text-rose-200" @click="del(g.id)">
                Удалить
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
  const res: any = await api.get('/admin/greetings')
  items.value = res.items || []
}
onMounted(load)

async function del(id: string) {
  if (!confirm('Удалить поздравление и DNS/файлы?')) return
  err.value = ''
  try {
    const res: any = await api.del(`/admin/greetings/${id}`)
    if (res.ok) await load()
    else err.value = res.error || 'DELETE_FAILED'
  } catch {
    err.value = 'DELETE_FAILED'
  }
}
</script>
