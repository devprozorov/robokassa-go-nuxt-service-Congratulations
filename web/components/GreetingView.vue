<template>
  <section class="relative min-h-[100svh] w-full overflow-hidden" :class="bgClass">
    <ThemeEffects :themeKind="themeKind" />

    <!-- soft bottom glow like the reference -->
    <div class="pointer-events-none absolute inset-0">
      <div class="absolute inset-x-0 bottom-0 h-44 bg-gradient-to-t from-black/45 to-transparent"></div>
      <div class="absolute -bottom-16 left-1/2 -translate-x-1/2 w-[1100px] h-[240px] rounded-full bg-amber-200/10 blur-3xl"></div>
      <div class="absolute -bottom-20 left-1/3 -translate-x-1/2 w-[700px] h-[220px] rounded-full bg-fuchsia-300/10 blur-3xl"></div>
    </div>

    <div class="relative min-h-[100svh] flex flex-col px-4 sm:px-6 md:px-10 py-10 md:py-14">
      <div class="text-xs text-slate-200/80">Happy</div>

      <div class="flex-1 flex flex-col justify-center">
        <div class="mt-3 text-center">
          <h1 class="title-font text-3xl sm:text-4xl md:text-6xl font-extrabold tracking-tight leading-tight">
            {{ title }}
          </h1>
          <p class="mt-3 md:mt-4 text-slate-100/90 text-base sm:text-lg whitespace-pre-wrap max-w-3xl mx-auto">
            {{ body }}
          </p>
        </div>

        <div class="mt-8 md:mt-10 w-full max-w-5xl mx-auto">
          <PhotoCollage :photos="photoUrls" />
        </div>

        <div class="mt-7 flex justify-center">
          <button
            type="button"
            class="px-5 py-3 rounded-2xl bg-white/10 border border-white/10 hover:bg-white/15 active:scale-[0.99] transition font-semibold"
            @click="openGift"
          >
            {{ giftOpen ? 'Скрыть подарок' : 'Открыть подарок' }}
          </button>
        </div>

        <div
          v-if="giftOpen"
          class="mt-4 w-full max-w-2xl mx-auto rounded-2xl border border-white/10 bg-black/20 p-5"
        >
          <div class="text-xs text-slate-300 mb-2">🎁 {{ giftLabel }}</div>

          <template v-if="gift.kind === 'promo'">
            <div class="text-sm text-slate-300">Промокод</div>
            <div class="mt-1 inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-white/10 border border-white/10 font-mono">
              {{ gift.value }}
            </div>
          </template>

          <template v-else-if="gift.kind === 'text'">
            <div class="whitespace-pre-wrap">{{ gift.value }}</div>
          </template>

          <template v-else-if="gift.kind === 'redirect'">
            <div class="text-sm text-slate-300">Ссылка</div>
            <a :href="gift.value" target="_blank" class="mt-1 inline-block underline break-all">
              {{ gift.value }}
            </a>
          </template>

          <template v-else>
            <div class="text-slate-300 text-sm">Подарок не задан</div>
          </template>
        </div>
      </div>

      <div class="mt-6 text-center text-xs text-slate-200/70">
        Публикация активна ограниченное время (7 дней после оплаты).
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import PhotoCollage from '~/components/PhotoCollage.vue'
import ThemeEffects from '~/components/ThemeEffects.vue'

const props = defineProps<{
  title: string
  body: string
  gift: { kind: string; value: string }
  photos: string[]
  bgClass: string
  themeKind?: string
}>()

const giftOpen = ref(false)

const giftLabel = computed(() => {
  const k = props.gift.kind
  if (k === 'promo') return 'Промокод'
  if (k === 'text') return 'Сообщение'
  if (k === 'redirect') return 'Переход на сайт'
  return 'Подарок'
})

const photoUrls = computed(() => props.photos)

async function fireConfetti() {
  if (process.server) return
  try {
    const mod: any = await import('canvas-confetti')
    const confetti = mod.default || mod
    confetti({ particleCount: 110, spread: 70, origin: { y: 0.65 } })
    setTimeout(() => confetti({ particleCount: 60, spread: 100, origin: { y: 0.65 } }), 120)
  } catch {
    // ignore
  }
}

const openGift = async () => {
  const next = !giftOpen.value
  giftOpen.value = next
  if (next) await fireConfetti()
}
</script>

<style scoped>
.title-font{
  font-family: ui-serif, Georgia, 'Times New Roman', Times, serif;
  font-style: italic;
  text-shadow: 0 10px 30px rgba(0,0,0,.35);
}
</style>
