<template>
  <div v-if="loading" class="min-h-[100svh] flex items-center justify-center text-slate-400 px-4">Загрузка...</div>
  <div v-else-if="!greeting" class="min-h-[100svh] flex items-center justify-center text-slate-400 px-4">Не найдено</div>

  <div v-else class="min-h-[100svh] flex flex-col">
    <div class="px-4 sm:px-6 py-4 border-b border-white/10 bg-slate-950/70 backdrop-blur sticky top-0 z-40 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
      <div class="text-sm text-slate-300">Предпросмотр</div>
      <NuxtLink :to="`/builder/${greeting.id}`" class="px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15">
        Назад в редактор
      </NuxtLink>
    </div>

  <div class="flex-1">
      <GreetingView
        :title="greeting.title"
        :body="greeting.body"
        :gift="giftView"
        :photos="photoUrls"
        :bgClass="bgClass"
        :themeKind="greeting.theme?.kind"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
import GreetingView from '~/components/GreetingView.vue'
import { THEMES, type ThemeKey } from '~/utils/themes'

const api = useApi()
const route = useRoute()
const id = route.params.id as string

const loading = ref(true)
const greeting = ref<any>(null)

const photoUrls = computed(() => {
  const photos: string[] = greeting.value?.photos || []
  return photos.map((fn) => `/api/media/${greeting.value.id}/${fn}`)
})

const bgClass = computed(() => {
  const kind = (greeting.value?.theme?.kind || 'birthday') as ThemeKey
  return greeting.value?.theme?.background || (THEMES[kind]?.bg ?? THEMES.birthday.bg)
})

const giftView = computed(() => {
  const g = greeting.value?.gift || { kind: 'text', value: '' }
  if (g.kind === 'image') return { kind: 'text', value: '' }
  return g
})

async function load() {
  loading.value = true
  const res: any = await api.get(`/greetings/${id}`)
  greeting.value = res.item || null
  loading.value = false
}
onMounted(load)
</script>