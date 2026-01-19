<template>
  <div class="max-w-md mx-auto">
    <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
      <h1 class="text-2xl font-bold">Регистрация</h1>
      <p class="text-sm text-slate-300 mt-1">Сразу получите одноразовый код восстановления</p>

      <form class="mt-5 space-y-3" @submit.prevent="submit">
        <input v-model="email" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="email" />
        <input v-model="username" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="username (3-32 символа)" />
        <input v-model="password" type="password" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="пароль (мин 8)" />
        <button class="w-full px-4 py-2 rounded-xl bg-white text-slate-950 font-semibold hover:bg-slate-100" :disabled="auth.loading">
          {{ auth.loading ? '...' : 'Создать аккаунт' }}
        </button>
        <div v-if="auth.error" class="text-sm text-rose-300">{{ auth.error }}</div>
      </form>

      <div v-if="recoveryCode" class="mt-4 rounded-2xl border border-emerald-400/20 bg-emerald-400/10 p-4">
        <div class="font-semibold">Ваш одноразовый код восстановления</div>
        <div class="mt-2 font-mono text-lg break-all">{{ recoveryCode }}</div>
        <div class="text-xs text-emerald-100/80 mt-2">Сохраните его. Показан один раз.</div>
        <button class="mt-3 px-3 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15" @click="goDashboard">
          В кабинет
        </button>
      </div>

      <div class="mt-4 text-sm text-slate-300">
        Уже есть аккаунт? <NuxtLink to="/auth/login" class="underline">Войти</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })
const auth = useAuthStore()

const email = ref('')
const username = ref('')
const password = ref('')
const recoveryCode = ref('')

async function submit() {
  const res: any = await auth.register(email.value, username.value, password.value)
  if (res.ok) {
    recoveryCode.value = res.recoveryCode || ''
  }
}
async function goDashboard() {
  await navigateTo('/dashboard')
}
</script>
