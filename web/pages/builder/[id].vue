<template>
  <div v-if="loading" class="text-slate-400">Загрузка...</div>

  <div v-else-if="!greeting" class="text-slate-400">Не найдено</div>

  <div v-else class="space-y-6">
    <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
      <div>
        <div class="text-xs text-slate-300">Редактор</div>
        <h1 class="text-xl sm:text-2xl font-bold break-words">{{ greeting.title }}</h1>
        <div class="text-sm text-slate-300 mt-1">
          Тип: <b>{{ greeting.type === 'code' ? 'по коду' : 'по поддомену' }}</b>
          <span v-if="greeting.type==='code' && greeting.code"> • код: <span class="font-mono">{{ greeting.code }}</span></span>
          <span v-if="greeting.type==='subdomain' && greeting.subdomain"> • поддомен: <span class="font-mono">{{ greeting.subdomain }}.{{ baseDomain }}</span></span>
        </div>
      </div>

      <div class="flex flex-wrap gap-2">
        <NuxtLink :to="`/preview/${greeting.id}`" class="px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15">
          Предпросмотр
        </NuxtLink>

        <button
          class="px-4 py-2 rounded-xl font-semibold"
          :class="greeting.paid ? 'bg-emerald-400/10 border border-emerald-400/20 text-emerald-200' : 'bg-white text-slate-950 hover:bg-slate-100'"
          @click="publish"
          :disabled="publishing || greeting.paid"
        >
          {{ greeting.paid ? 'Оплачено' : (publishing ? '...' : 'Оплатить и опубликовать') }}
        </button>
      </div>
    </div>

    <ContentPolicyBanner class="mb-4" />

    <div class="grid lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-4">
        <div class="rounded-3xl border border-white/10 bg-white/5 p-6 space-y-3">
          <div class="font-semibold">Текст</div>
          <input v-model="form.title" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="Заголовок" />
          <textarea v-model="form.body" rows="6" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="Основной текст" />
          <div v-if="greeting.type==='subdomain'" class="pt-2">
            <div class="text-sm text-slate-300 mb-2">Поддомен</div>
            <input v-model="form.subdomain" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="имя" />
            <div class="text-xs text-slate-400 mt-1">Например: <span class="font-mono">love</span> → love.{{ baseDomain }}</div>
          </div>

          <div class="pt-2">
            <div class="text-sm text-slate-300 mb-2">Тема</div>
            <ThemePicker :selected="themeKey" @update:selected="(v)=>themeKey=v" />
          </div>

          <div v-if="themeKey==='custom'" class="grid md:grid-cols-2 gap-3 pt-2">
            <div>
              <div class="text-xs text-slate-300 mb-1">Primary</div>
              <input v-model="customPrimary" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="#a78bfa" />
            </div>
            <div>
              <div class="text-xs text-slate-300 mb-1">Secondary</div>
              <input v-model="customSecondary" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="#22d3ee" />
            </div>
            <div class="md:col-span-2">
              <div class="text-xs text-slate-300 mb-1">Background (CSS classes)</div>
              <input v-model="customBg" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950" />
            </div>
          </div>

          <div class="flex gap-2 pt-2">
            <button class="px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15" @click="save">
              Сохранить
            </button>
            <button class="px-4 py-2 rounded-xl bg-rose-500/10 border border-rose-500/20 hover:bg-rose-500/15 text-rose-200" @click="remove">
              Удалить
            </button>
          </div>
          <div v-if="msg" class="text-sm text-emerald-200">{{ msg }}</div>
          <div v-if="err" class="text-sm text-rose-300">{{ err }}</div>
        </div>

        <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div class="flex items-center justify-between">
            <div class="font-semibold">Фото-коллаж</div>
            <div class="text-xs text-slate-400">jpg/png/webp</div>
          </div>

          <div class="mt-3 flex flex-col md:flex-row gap-3">
            <input ref="fileInput" type="file" class="w-full" multiple accept="image/*" />
            <button class="px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15" @click="upload">
              Загрузить
            </button>
          </div>

          <div class="mt-4">
            <PhotoCollage :photos="photoUrls" />
          </div>
        </div>
      </div>

      <div class="space-y-4">
        <GiftEditor v-model="gift" />

        <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div class="font-semibold">Статус</div>
          <div class="mt-2 text-sm text-slate-300">
            <div>Оплата: <b>{{ greeting.paid ? 'Да' : 'Нет' }}</b></div>
            <div v-if="greeting.paid && greeting.expiresAt">Активно до: <b>{{ new Date(greeting.expiresAt).toLocaleString() }}</b></div>
          </div>
          <div class="mt-4 text-xs text-slate-400">
            После оплаты публикация активна 7 дней, затем автоматически удаляется worker-ом.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
import ThemePicker from '~/components/ThemePicker.vue'
import GiftEditor from '~/components/GiftEditor.vue'
import PhotoCollage from '~/components/PhotoCollage.vue'
import ContentPolicyBanner from '~/components/ContentPolicyBanner.vue'
import { THEMES, type ThemeKey } from '~/utils/themes'

const api = useApi()
const baseDomain = useRuntimeConfig().public.baseDomain

const route = useRoute()
const id = route.params.id as string

const loading = ref(true)
const greeting = ref<any>(null)

const form = reactive({ title: '', body: '', subdomain: '' })
const gift = ref<{kind: any, value: string}>({ kind: 'text', value: '' })

const themeKey = ref<ThemeKey>('birthday')
const customPrimary = ref('#a78bfa')
const customSecondary = ref('#22d3ee')
const customBg = ref('bg-slate-950')

const msg = ref('')
const err = ref('')
const publishing = ref(false)

const fileInput = ref<HTMLInputElement | null>(null)

const photoUrls = computed(() => {
  const photos: string[] = greeting.value?.photos || []
  return photos.map((fn) => `/api/media/${greeting.value.id}/${fn}`)
})

function applyFromGreeting() {
  form.title = greeting.value.title || ''
  form.body = greeting.value.body || ''
  form.subdomain = greeting.value.subdomain || ''
  gift.value = greeting.value.gift || { kind: 'text', value: '' }

  const kind = (greeting.value.theme?.kind || 'birthday') as ThemeKey
  if (THEMES[kind]) {
    themeKey.value = kind
  } else {
    themeKey.value = 'birthday'
  }

  if (themeKey.value === 'custom') {
    customPrimary.value = greeting.value.theme?.primary || '#a78bfa'
    customSecondary.value = greeting.value.theme?.secondary || '#22d3ee'
    customBg.value = greeting.value.theme?.background || 'bg-slate-950'
  }
}

async function load() {
  loading.value = true
  const res: any = await api.get(`/greetings/${id}`)
  greeting.value = res.item || null
  if (greeting.value) applyFromGreeting()
  loading.value = false
}
onMounted(load)

async function save() {
  msg.value = ''
  err.value = ''
  try {
    const theme = {
      kind: themeKey.value,
      primary: themeKey.value === 'custom' ? customPrimary.value : THEMES[themeKey.value].primary,
      secondary: themeKey.value === 'custom' ? customSecondary.value : THEMES[themeKey.value].secondary,
      background: themeKey.value === 'custom' ? customBg.value : THEMES[themeKey.value].bg,
    }
    const res: any = await api.put(`/greetings/${id}`, {
      title: form.title,
      body: form.body,
      // subdomain: form.subdomain,
      theme,
      // gift: gift.value,
    })
    if (res.ok) {
      await load()
      msg.value = 'Сохранено'
      return
    }
    err.value = res.error || 'SAVE_FAILED'
  } catch {
    err.value = 'SAVE_FAILED'
  }
}

async function upload() {
  msg.value = ''
  err.value = ''
  const el = fileInput.value
  if (!el || !el.files || el.files.length === 0) return

  const fd = new FormData()
  for (const f of Array.from(el.files)) fd.append('files', f)

  try {
    const res: any = await $fetch(`/api/greetings/${id}/photos`, { method: 'POST', body: fd, credentials: 'include' })
    if (res.ok) {
      await load()
      msg.value = 'Фото загружены'
      el.value = ''
      return
    }
    err.value = res.error || 'UPLOAD_FAILED'
  } catch {
    err.value = 'UPLOAD_FAILED'
  }
}

async function publish() {
  publishing.value = true
  err.value = ''
  try {
    await save()
    const res: any = await api.post(`/greetings/${id}/publish`)
    if (res.ok) {
      // redirect to checkout (stub by default)
      await navigateTo(res.checkoutUrl, { external: true })
      return
    }
    err.value = res.error || 'PUBLISH_FAILED'
  } catch {
    err.value = 'PUBLISH_FAILED'
  } finally {
    publishing.value = false
  }
}

async function remove() {
  if (!confirm('Удалить поздравление?')) return
  try {
    const res: any = await api.del(`/greetings/${id}`)
    if (res.ok) {
      await navigateTo('/dashboard')
    }
  } catch {}
}
</script>
