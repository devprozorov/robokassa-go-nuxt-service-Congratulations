<template>
  <div>
    <div class="flex items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold">Кабинет</h1>
        <p class="text-sm text-slate-300 mt-1">Создавай, редактируй и публикуй поздравления</p>
      </div>

      <button class="px-4 py-2 rounded-xl bg-white text-slate-950 font-semibold hover:bg-slate-100" @click="openCreate = true">
        + Создать
      </button>
    </div>

    <div v-if="openCreate" class="mt-6 rounded-3xl border border-white/10 bg-white/5 p-6">
      <div class="flex items-center justify-between">
        <div class="font-semibold">Новое поздравление</div>
        <button class="text-sm text-slate-300 hover:text-white" @click="openCreate=false">Закрыть</button>
      </div>

      <div class="mt-4 grid md:grid-cols-2 gap-4">
        <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
          <div class="font-semibold">По коду</div>
          <div class="text-sm text-slate-300 mt-1">Откроется по ссылке <code>/c/CODE</code></div>
          <button class="mt-4 px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15" @click="create('code')">
            Создать
          </button>
        </div>

        <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
          <div class="font-semibold">По поддомену</div>
          <div class="text-sm text-slate-300 mt-1">Откроется на <code>имя.{{ baseDomain }}</code></div>
          <div class="mt-3 flex gap-2">
            <input v-model="subdomain" class="flex-1 px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="имя" />
            <button class="px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15" @click="create('subdomain')">
              Создать
            </button>
          </div>
          <div class="text-xs text-slate-400 mt-2">Поддомен можно поменять до оплаты.</div>
        </div>
      </div>

      <div v-if="createError" class="mt-3 text-sm text-rose-300">{{ createError }}</div>
    </div>

    <div class="mt-8 grid md:grid-cols-2 gap-4">
      <div v-for="g in items" :key="g.id" class="rounded-3xl border border-white/10 bg-white/5 p-6">
        <div class="flex items-start justify-between gap-3">
          <div>
            <div class="text-xs text-slate-300">{{ g.type === 'code' ? 'По коду' : 'По поддомену' }}</div>
            <div class="text-lg font-semibold">{{ g.title }}</div>
            <div class="text-sm text-slate-300 line-clamp-2 mt-1">{{ g.body }}</div>
          </div>
          <div class="text-xs">
            <span v-if="g.paid" class="px-2 py-1 rounded-full bg-emerald-400/10 border border-emerald-400/20 text-emerald-200">Оплачено</span>
            <span v-else class="px-2 py-1 rounded-full bg-white/10 border border-white/10 text-slate-200">Черновик</span>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap gap-2">
          <NuxtLink :to="`/builder/${g.id}`" class="px-3 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15 text-sm">
            Редактировать
          </NuxtLink>
          <NuxtLink :to="`/preview/${g.id}`" class="px-3 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15 text-sm">
            Предпросмотр
          </NuxtLink>

          <NuxtLink v-if="g.paid && g.type==='code'" :to="`/c/${g.code}`" class="px-3 py-2 rounded-xl bg-emerald-400/10 border border-emerald-400/20 hover:bg-emerald-400/15 text-sm">
            Открыть
          </NuxtLink>
        </div>

        <div v-if="g.paid && g.expiresAt" class="mt-3 text-xs text-slate-400">
          Активно до: {{ new Date(g.expiresAt).toLocaleString() }}
        </div>
      </div>
    </div>

    <div v-if="items.length === 0" class="mt-8 text-slate-400">
      Пока нет поздравлений. Нажми «Создать».
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const api = useApi()
const baseDomain = useRuntimeConfig().public.baseDomain

const items = ref<any[]>([])
const openCreate = ref(false)
const subdomain = ref('')
const createError = ref('')

async function load() {
  const res: any = await api.get('/greetings')
  items.value = res.items || []
}
onMounted(load)

async function create(type: 'code' | 'subdomain') {
  createError.value = ''
  try {
    const body: any = { type }
    if (type === 'subdomain') body.subdomain = subdomain.value
    const res: any = await api.post('/greetings', body)
    if (res.ok) {
      openCreate.value = false
      await navigateTo(`/builder/${res.item.id}`)
      return
    }
    createError.value = res.error || 'CREATE_FAILED'
  } catch {
    createError.value = 'CREATE_FAILED'
  }
}
</script>
