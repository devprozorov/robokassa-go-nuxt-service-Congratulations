<template>
  <div class="grid grid-cols-2 md:grid-cols-5 gap-3">
    <button
      v-for="(t, key) in themes"
      :key="key"
      class="rounded-xl border px-3 py-3 text-left hover:border-white/30 transition"
      :class="selected === key ? 'border-white/40 bg-white/10' : 'border-white/10 bg-white/5'"
      @click="$emit('update:selected', key)"
      type="button"
    >
      <div class="text-xs text-slate-300">Тема</div>
      <div class="font-medium">{{ t.label }}</div>
    </button>
  </div>
</template>

<script setup lang="ts">
import { THEMES, type ThemeKey } from '~/utils/themes'
defineProps<{ selected: ThemeKey }>()
defineEmits<{ (e: 'update:selected', v: ThemeKey): void }>()

const themes = THEMES
</script>
