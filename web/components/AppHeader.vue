<template>
  <header class="border-b border-white/10 bg-slate-950/70 backdrop-blur sticky top-0 z-50">
    <div class="max-w-6xl mx-auto px-3 sm:px-4 md:px-6 py-3 flex items-center justify-between">
      <NuxtLink to="/" class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-xl bg-gradient-to-br from-fuchsia-500 to-cyan-400"></div>
        <div class="font-semibold tracking-tight">{{ appName }}</div>
      </NuxtLink>

      <!-- Desktop nav -->
      <nav class="hidden md:flex items-center gap-3 text-sm">
        <NuxtLink to="/docs" class="text-slate-300 hover:text-white">Документация</NuxtLink>
        <NuxtLink to="/enter-code" class="text-slate-300 hover:text-white">Открыть по коду</NuxtLink>
        <NuxtLink v-if="auth.isAuthed" to="/dashboard" class="text-slate-300 hover:text-white">Кабинет</NuxtLink>
        <NuxtLink v-if="auth.isAdmin" to="/admin" class="text-slate-300 hover:text-white">Админ</NuxtLink>

        <div class="w-px h-5 bg-white/10 mx-1"></div>

        <NuxtLink
          v-if="!auth.isAuthed"
          to="/auth/login"
          class="px-3 py-2 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition"
        >
          Войти
        </NuxtLink>

        <NuxtLink
          v-if="!auth.isAuthed"
          to="/auth/register"
          class="px-3 py-2 rounded-xl bg-white text-slate-900 font-semibold hover:bg-white/90 transition"
        >
          Регистрация
        </NuxtLink>

        <button
          v-if="auth.isAuthed"
          type="button"
          class="px-3 py-2 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition"
          @click="logout"
        >
          Выйти
        </button>
      </nav>

      <!-- Mobile button -->
      <button
        type="button"
        class="md:hidden inline-flex items-center justify-center w-10 h-10 rounded-xl border border-white/10 bg-white/5 hover:bg-white/10 transition"
        aria-label="Меню"
        @click="menuOpen = !menuOpen"
      >
        <svg v-if="!menuOpen" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" d="M6 6l12 12M18 6L6 18" />
        </svg>
      </button>
    </div>

    <!-- Mobile menu -->
    <div v-if="menuOpen" class="md:hidden border-t border-white/10 bg-slate-950/95">
      <div class="max-w-6xl mx-auto px-3 sm:px-4 md:px-6 py-3 flex flex-col gap-2 text-sm">
        <NuxtLink to="/docs" class="px-3 py-2 rounded-xl hover:bg-white/5" @click="menuOpen=false">Документация</NuxtLink>
        <NuxtLink to="/enter-code" class="px-3 py-2 rounded-xl hover:bg-white/5" @click="menuOpen=false">Открыть по коду</NuxtLink>
        <NuxtLink v-if="auth.isAuthed" to="/dashboard" class="px-3 py-2 rounded-xl hover:bg-white/5" @click="menuOpen=false">Кабинет</NuxtLink>
        <NuxtLink v-if="auth.isAdmin" to="/admin" class="px-3 py-2 rounded-xl hover:bg-white/5" @click="menuOpen=false">Админ</NuxtLink>

        <div class="h-px bg-white/10 my-2"></div>

        <NuxtLink
          v-if="!auth.isAuthed"
          to="/auth/login"
          class="px-3 py-2 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition"
          @click="menuOpen=false"
        >
          Войти
        </NuxtLink>

        <NuxtLink
          v-if="!auth.isAuthed"
          to="/auth/register"
          class="px-3 py-2 rounded-xl bg-white text-slate-900 font-semibold hover:bg-white/90 transition"
          @click="menuOpen=false"
        >
          Регистрация
        </NuxtLink>

        <button
          v-if="auth.isAuthed"
          type="button"
          class="text-left px-3 py-2 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition"
          @click="logout"
        >
          Выйти
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '~/stores/auth'

const appName = 'Happy'
const auth = useAuthStore()

const menuOpen = ref(false)
const route = useRoute()
watch(() => route.fullPath, () => { menuOpen.value = false })

async function logout() {
  await auth.logout()
  menuOpen.value = false
}
</script>
