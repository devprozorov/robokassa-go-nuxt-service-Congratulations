<template>
  <div class="relative rounded-3xl overflow-hidden border border-white/10 bg-white/5 p-3 sm:p-4">
    <div v-if="photos.length === 0" class="text-sm text-slate-400 p-6 text-center">
      Фото появятся тут ✨
    </div>

    <div v-else class="flex gap-3 overflow-x-auto pb-2 snap-x snap-mandatory">
      <button
        v-for="(src, i) in photos"
        :key="src + i"
        type="button"
        class="group relative shrink-0 w-24 h-24 sm:w-28 sm:h-28 md:w-36 md:h-36 snap-start rounded-2xl overflow-hidden border border-white/10 bg-slate-900 hover:border-white/30 transition"
        :style="{ transform: `rotate(${((i%2)*2-1)*2}deg)` }"
        @click="open(i)"
      >
        <img :src="src" class="w-full h-full object-cover collage-float group-hover:scale-[1.03] transition" />
        <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition"></div>
      </button>
    </div>

    <div class="pointer-events-none absolute -bottom-10 left-0 right-0 h-20 bg-gradient-to-t from-white/10 to-transparent blur-2xl"></div>

    <!-- Lightbox -->
    <Teleport to="body">
      <div v-if="openIndex !== null" class="fixed inset-0 z-50">
        <div class="absolute inset-0 bg-black/70" @click="close"></div>

        <div class="relative h-full w-full flex items-center justify-center p-4">
          <div class="relative max-w-4xl w-full">
            <img
              :src="photos[openIndex]"
              class="w-full max-h-[85vh] object-contain rounded-2xl border border-white/10 bg-slate-950"
            />

            <button
              class="absolute -top-3 -right-3 w-10 h-10 rounded-full bg-white/10 border border-white/10 hover:bg-white/15"
              @click="close"
              aria-label="Закрыть"
              type="button"
            >
              ✕
            </button>

            <button
              v-if="photos.length > 1"
              class="absolute left-3 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/10 border border-white/10 hover:bg-white/15"
              @click="prev"
              type="button"
              aria-label="Предыдущее"
            >
              ‹
            </button>

            <button
              v-if="photos.length > 1"
              class="absolute right-3 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/10 border border-white/10 hover:bg-white/15"
              @click="next"
              type="button"
              aria-label="Следующее"
            >
              ›
            </button>

            <div class="mt-3 text-center text-xs text-slate-300">
              {{ openIndex + 1 }} / {{ photos.length }} • ESC — закрыть
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ photos: string[] }>()

const openIndex = ref<number | null>(null)

function open(i: number) {
  openIndex.value = i
}
function close() {
  openIndex.value = null
}
function next() {
  if (openIndex.value === null) return
  openIndex.value = (openIndex.value + 1) % props.photos.length
}
function prev() {
  if (openIndex.value === null) return
  openIndex.value = (openIndex.value - 1 + props.photos.length) % props.photos.length
}

function onKey(e: KeyboardEvent) {
  if (openIndex.value === null) return
  if (e.key === 'Escape') close()
  if (e.key === 'ArrowRight') next()
  if (e.key === 'ArrowLeft') prev()
}

onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
</script>

<style scoped>
.collage-float {
  animation: floaty 4.8s ease-in-out infinite;
}
@keyframes floaty {
  0%, 100% { transform: translateY(0); filter: saturate(1); }
  50% { transform: translateY(-4px); filter: saturate(1.12); }
}
</style>
