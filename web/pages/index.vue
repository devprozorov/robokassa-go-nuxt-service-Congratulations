<template>
  <div class="space-y-10">

     <section class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
      <div class="text-xs text-slate-300">Открыть поздравление по коду</div>
      <div class="mt-2 text-lg font-semibold">Введите код и откройте страницу</div>

      <div class="mt-4 flex flex-col sm:flex-row gap-3">
        <input
          v-model="code"
          class="w-full px-4 py-3 rounded-2xl bg-slate-900 border border-white/10"
          placeholder="Например: A1B2C3"
        />
        <button
          class="px-5 py-3 rounded-2xl bg-white text-slate-950 font-semibold hover:bg-slate-100 disabled:opacity-60"
          :disabled="!code.trim()"
          @click="openByCode"
        >
          Открыть
        </button>
      </div>

      <div class="mt-2 text-xs text-slate-400">
        Код можно получить у отправителя или в кабинете после публикации.
      </div>
    </section>


    <section class="rounded-3xl border border-white/10 bg-gradient-to-br from-slate-950 via-slate-950 to-indigo-950 overflow-hidden">
      <div class="p-8 md:p-12">
        <div class="inline-flex items-center gap-2 text-xs px-3 py-1 rounded-full bg-white/10 border border-white/10">
          <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
          Создай поздравление за 2 минуты
        </div>

        <h1 class="mt-4 text-4xl md:text-6xl font-extrabold tracking-tight leading-tight">
          Happy — персональные поздравления<br class="hidden md:block" />
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-fuchsia-400 to-cyan-300">по коду</span>
          или
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-300 to-rose-300">по поддомену</span>
        </h1>

        <p class="mt-4 text-slate-200/90 max-w-2xl">
          Выбирай тематику, добавляй фото-коллаж, подарок и получай красивую страницу.
          Публикация активна 7 дней после оплаты.
        </p>

        <div class="mt-6 flex flex-col sm:flex-row gap-3">
          <NuxtLink to="/auth/register" class="px-5 py-3 rounded-2xl bg-white text-slate-950 font-semibold hover:bg-slate-100">
            Начать бесплатно
          </NuxtLink>
          <NuxtLink to="/docs" class="px-5 py-3 rounded-2xl bg-white/10 border border-white/10 hover:bg-white/15">
            Как это работает
          </NuxtLink>
        </div>
      </div>
    </section>

    <section class="grid md:grid-cols-3 gap-4">
      <div class="rounded-2xl border border-white/10 bg-white/5 p-5">
        <div class="text-sm text-slate-300">Шаг 1</div>
        <div class="font-semibold mt-1">Создай и настрой</div>
        <div class="text-sm text-slate-300 mt-2">Тема, текст, подарок, фото-коллаж, предпросмотр.</div>
      </div>
      <div class="rounded-2xl border border-white/10 bg-white/5 p-5">
        <div class="text-sm text-slate-300">Шаг 2</div>
        <div class="font-semibold mt-1">Оплати публикацию</div>
        <div class="text-sm text-slate-300 mt-2">По коду дешевле, по поддомену — дороже но оригинальней.</div>
      </div>
      <div class="rounded-2xl border border-white/10 bg-white/5 p-5">
        <div class="text-sm text-slate-300">Шаг 3</div>
        <div class="font-semibold mt-1">Поделись</div>
        <div class="text-sm text-slate-300 mt-2">Код / ссылка на поддомен. Активно 7 дней.</div>
      </div>
    </section>

    <section class="grid md:grid-cols-2 gap-4">
      <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
        <div class="text-xs text-slate-300">Разовая оплата страницы на 7 дней - 59 рублей</div>
        <div class="text-2xl font-bold mt-1">Поздравление по коду</div>
        <ul class="mt-3 text-sm text-slate-300 space-y-2">
          <li>• Открывается по вводу на главном экране в поле "Введите код и откройте страницу"</li>
          <li>• Код генерируется случайным образом и состоит из букв и чисел</li>
          <li>• Быстро и удобно</li>
        </ul>
        <div class="mt-5">
          <NuxtLink to="/auth/register" class="inline-flex px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15">
            Создать
          </NuxtLink>
        </div>
      </div>

      <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
        <div class="text-xs text-slate-300">(Разовая оплата страницы на 7 дней - 99 рублей)</div>
        <div class="text-2xl font-bold mt-1">Поздравление по поддомену</div>
        <ul class="mt-3 text-sm text-slate-300 space-y-2">
          <li>• Например `ForRoma.{{ baseDomain }}`</li>
          <li>• Красивая ссылка которой можно сразу поделиться и при открытии будет поздравление</li>
          <li>• Лучший способ удивить родного человека</li>
        </ul>
        <div class="mt-5">
          <NuxtLink to="/auth/register" class="inline-flex px-4 py-2 rounded-xl bg-white/10 border border-white/10 hover:bg-white/15">
            Создать
          </NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const baseDomain = useRuntimeConfig().public.baseDomain

const code = ref('')
function openByCode() {
  const v = code.value.trim()
  if (!v) return
  return navigateTo(`/c/${encodeURIComponent(v)}`)
}

</script>
