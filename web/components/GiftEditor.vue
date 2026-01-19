<template>
  <div class="rounded-2xl border border-white/10 bg-white/5 p-4">
    <div class="font-semibold mb-3">🎁 Подарок</div>

    <div class="grid md:grid-cols-3 gap-2 mb-3">
      <button
        v-for="k in kinds"
        :key="k"
        class="px-3 py-2 rounded-xl text-sm border transition"
        :class="gift.kind === k ? 'border-white/40 bg-white/10' : 'border-white/10 hover:bg-white/5'"
        type="button"
        @click="setKind(k)"
      >
        {{ labels[k] }}
      </button>
    </div>

    <div v-if="gift.kind === 'promo'" class="space-y-2">
      <div class="text-xs text-slate-300">Промокод</div>
      <input v-model="gift.value" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="PROMO2026" />
      <div class="text-xs text-slate-500">Покажется после нажатия «Открыть подарок».</div>
    </div>

    <div v-else-if="gift.kind === 'text'" class="space-y-2">
      <div class="text-xs text-slate-300">Текст</div>
      <textarea v-model="gift.value" rows="4" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="Сюрприз внутри 🙂"></textarea>
    </div>

    <div v-else class="space-y-2">
      <div class="text-xs text-slate-300">Ссылка</div>
      <input v-model="gift.value" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="https://example.com" />
      <div class="text-xs text-slate-500">Откроется в новой вкладке.</div>
    </div>
  </div>
</template>

<script setup lang="ts">
type GiftKind = 'promo' | 'text' | 'redirect'
type Gift = { kind: GiftKind; value: string }

const props = defineProps<{ modelValue: Gift }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: Gift): void }>()

const kinds: GiftKind[] = ['promo', 'text', 'redirect']
const labels: Record<GiftKind, string> = {
  promo: 'Промокод',
  text: 'Текст',
  redirect: 'Переход',
}

const gift = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

function setKind(k: GiftKind) {
  gift.value = { kind: k, value: '' }
}
</script>
