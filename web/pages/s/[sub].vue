<template>
  <div v-if="loading" class="min-h-[100svh] flex items-center justify-center text-slate-400 px-4">Загрузка...</div>
  <div v-else-if="!item" class="min-h-[100svh] flex items-center justify-center text-slate-400 px-4 text-center">Поздравление не найдено или срок истёк</div>

  <div v-else>
    <GreetingView
      :title="item.title"
      :body="item.body"
      :gift="giftView"
      :photos="photoUrls"
      :bgClass="bgClass"
      :themeKind="item.theme?.kind"
    />
  </div>
</template>

<script setup lang="ts">
import GreetingView from '~/components/GreetingView.vue'
import { THEMES, type ThemeKey } from '~/utils/themes'

const api = useApi()
const route = useRoute()
const sub = route.params.sub as string

const loading = ref(true)
const item = ref<any>(null)

const photoUrls = computed(() => {
  const photos: string[] = item.value?.photos || []
  return photos.map((fn) => `/api/media/${item.value.id}/${fn}`)
})

const bgClass = computed(() => {
  const kind = (item.value?.theme?.kind || 'birthday') as ThemeKey
  return item.value?.theme?.background || (THEMES[kind]?.bg ?? THEMES.birthday.bg)
})

const giftView = computed(() => {
  const g = item.value?.gift || { kind: 'text', value: '' }
  if (g.kind === 'image') return { kind: 'text', value: '' }
  return g
})

async function load() {
  loading.value = true
  const res: any = await api.get(`/public/subdomain/${sub}`)
  item.value = res.item || null
  loading.value = false
}
onMounted(load)
</script>