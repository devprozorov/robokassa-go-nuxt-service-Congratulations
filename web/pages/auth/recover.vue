<template>
  <div class="max-w-md mx-auto">
    <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
      <h1 class="text-2xl font-bold">Восстановление пароля</h1>
      <p class="text-sm text-slate-300 mt-1">Используй одноразовый резервный код</p>

      <form class="mt-5 space-y-3" @submit.prevent="submit">
        <input v-model="loginOrEmail" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="email или username" />
        <input v-model="recoveryCode" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="одноразовый код" />
        <input v-model="newPassword" type="password" class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-white/10" placeholder="новый пароль" />
        <button class="w-full px-4 py-2 rounded-xl bg-white text-slate-950 font-semibold hover:bg-slate-100">
          Сбросить пароль
        </button>
        <div v-if="error" class="text-sm text-rose-300">{{ error }}</div>
      </form>

      <div v-if="newRecovery" class="mt-4 rounded-2xl border border-emerald-400/20 bg-emerald-400/10 p-4">
        <div class="font-semibold">Новый код восстановления</div>
        <div class="mt-2 font-mono text-lg break-all">{{ newRecovery }}</div>
        <div class="text-xs text-emerald-100/80 mt-2">Сохраните его. Этот код тоже одноразовый.</div>
        <button class="mt-3 px-3 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15" @click="goDashboard">
          В кабинет
        </button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })
const error = ref('')
const newRecovery = ref('')

const loginOrEmail = ref('')
const recoveryCode = ref('')
const newPassword = ref('')

async function submit() {
  error.value = ''
  newRecovery.value = ''
  try {
    const res: any = await $fetch('/api/auth/recover/reset', { method: 'POST', body: { loginOrEmail: loginOrEmail.value, recoveryCode: recoveryCode.value, newPassword: newPassword.value }, credentials: 'include' })
    if (res?.ok) {
      newRecovery.value = res.newRecoveryCode
      await useAuthStore().fetchMe()
      return
    }
    error.value = res?.error || 'RECOVER_FAILED'
  } catch {
    error.value = 'RECOVER_FAILED'
  }
}

async function goDashboard() {
  await navigateTo('/dashboard')
}
</script>
