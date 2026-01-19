<template>
  <div class="rounded-2xl border border-white/10 bg-white/5 p-4">
    <div class="text-sm text-slate-200 mb-3">Оплата</div>

    <button
      class="px-4 py-2 rounded-xl bg-white/10 hover:bg-white/15 text-white"
      :disabled="loading"
      @click="pay"
    >
      {{ loading ? '...' : 'Оплатить' }}
    </button>

    <div v-if="err" class="mt-3 text-sm text-red-300">{{ err }}</div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ orderId: string }>()
const loading = ref(false)
const err = ref('')

async function pay() {
  err.value = ''
  loading.value = true
  try {
    const res: any = await $fetch('/api/payments/robokassa/init', {
      method: 'POST',
      credentials: 'include',
      body: { orderId: props.orderId },
    })

    if (!res?.ok || !res?.payUrl) {
      err.value = res?.error || 'ROBO_INIT_FAILED'
      return
    }

    window.location.href = res.payUrl
  } catch {
    err.value = 'ROBO_INIT_FAILED'
  } finally {
    loading.value = false
  }
}
</script>
