<template>
  <footer class="border-t border-white/10 mt-10">
    <div class="max-w-6xl mx-auto px-4 py-8 text-sm text-slate-400">
      <div class="flex flex-col md:flex-row gap-4 md:items-center md:justify-between">
        <div>© {{ new Date().getFullYear() }} ООО "Унистар" </div>

        <div class="flex flex-wrap gap-x-4 gap-y-2">
          <NuxtLink to="/" class="hover:text-slate-200">Главная</NuxtLink>
          <NuxtLink to="/docs" class="hover:text-slate-200">Документация</NuxtLink>
          <NuxtLink to="/legal/offer" class="hover:text-slate-200">Оферта</NuxtLink>
          <NuxtLink to="/legal/company" class="hover:text-slate-200">Юр. лицо</NuxtLink>
          <NuxtLink to="/legal/privacy" class="hover:text-slate-200">Политика ПДн</NuxtLink>
          <NuxtLink to="/legal/consent" class="hover:text-slate-200">Согласие ПДн</NuxtLink>
        </div>
      </div>
    </div>
  </footer>
</template>
