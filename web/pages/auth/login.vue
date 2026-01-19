<template>
  <div class="max-w-md mx-auto">
    <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
      <h1 class="text-2xl font-bold">Вход</h1>
      <p class="text-sm text-slate-300 mt-1">Email или логин + пароль</p>

      <form class="mt-5 space-y-3" @submit.prevent="submit">
        <input v-model="loginOrEmail" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="email или username" />
        <input v-model="password" type="password" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="пароль" />
        <button class="w-full px-4 py-2 rounded-xl bg-white text-slate-950 font-semibold hover:bg-slate-100" :disabled="auth.loading">
          {{ auth.loading ? '...' : 'Войти' }}
        </button>
        <div v-if="auth.error" class="text-sm text-rose-300">{{ auth.error }}</div>
      </form>

      <div class="mt-4 text-sm text-slate-300 flex justify-between">
        <NuxtLink to="/auth/register" class="underline">Регистрация</NuxtLink>
        <NuxtLink to="/auth/recover" class="underline">Забыл пароль</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })
const auth = useAuthStore()

const route = useRoute()
const loginOrEmail = ref('')
const password = ref('')

async function submit() {
  const res = await auth.login(loginOrEmail.value, password.value)
  if (res.ok) {
    const next = (route.query.next as string) || '/dashboard'
    await navigateTo(next)
  }
}
</script>
