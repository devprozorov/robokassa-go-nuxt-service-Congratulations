<template>
  <div class="max-w-2xl">
    
    <NuxtLink to="/dashboard" class="text-sm text-slate-300 hover:text-white">← назад</NuxtLink>

    <div class="mx-auto max-w-3xl px-4 py-6">
        <PaymentBox :order-id="orderId" />
      </div>

    <div class="mt-4 rounded-3xl border border-white/10 bg-white/5 p-6">
      <h1 class="text-2xl font-bold">Заказ</h1>
      <p class="text-sm text-slate-300 mt-1">
        Статус оплаты: <b>{{ order?.status || '...' }}</b>
      </p>
         <PayRobokassa :order-id="orderId" />
      <div v-if="order" class="mt-4 grid md:grid-cols-2 gap-3 text-sm text-slate-200/90">
        <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
          <div class="text-xs text-slate-300">Сумма</div>
          <div class="text-lg font-semibold">{{ order.amountRUB }} {{ order.currency }}</div>
        </div>
        <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
          <div class="text-xs text-slate-300">Услуга</div>
          <div class="text-lg font-semibold">{{ order.productType === 'subdomain' ? 'По поддомену' : 'По коду' }}</div>
        </div>
      </div>

      <div class="mt-5 text-sm text-slate-300">
        <div v-if="order?.status === 'pending'">
        <div class="text-sm text-slate-300 mb-2">Оплатить:</div>
     
        </div>
        <div v-else-if="order?.status === 'paid'">
          Оплата подтверждена. Поздравление активировано на 7 дней.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })

const route = useRoute()
const orderId = computed(() => String(route.params.id || ''))
</script>



