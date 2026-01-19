<template>
  <div class="grid lg:grid-cols-[260px,1fr] gap-6">
    <!-- Notion-like sidebar -->
    <aside class="hidden lg:block">
      <div class="sticky top-6 rounded-2xl border border-white/10 bg-white/5 p-4">
        <div class="text-xs text-slate-300 mb-2">Документация</div>
        <nav class="space-y-1 text-sm">
          <a class="block px-3 py-2 rounded-xl hover:bg-white/5" href="#start">Как начать</a>
          <a class="block px-3 py-2 rounded-xl hover:bg-white/5" href="#types">2 услуги: код и поддомен</a>
          <a class="block px-3 py-2 rounded-xl hover:bg-white/5" href="#builder">Редактор</a>
          <a class="block px-3 py-2 rounded-xl hover:bg-white/5" href="#photos">Фото и просмотр</a>
          <a class="block px-3 py-2 rounded-xl hover:bg-white/5" href="#publish">Оплата и 7 дней</a>
          <a class="block px-3 py-2 rounded-xl hover:bg-white/5" href="#security">Вход и восстановление</a>          <div class="mt-3 pt-3 border-t border-white/10">
            <div class="text-xs text-slate-400 mb-2">Юридическая информация</div>
            <NuxtLink class="block px-3 py-2 rounded-xl hover:bg-white/5" to="/legal/offer">Оферта</NuxtLink>
            <NuxtLink class="block px-3 py-2 rounded-xl hover:bg-white/5" to="/legal/company">Юр. лицо</NuxtLink>
            <NuxtLink class="block px-3 py-2 rounded-xl hover:bg-white/5" to="/legal/privacy">Политика ПДн</NuxtLink>
            <NuxtLink class="block px-3 py-2 rounded-xl hover:bg-white/5" to="/legal/consent">Согласие ПДн</NuxtLink>
          </div>
        </nav>
      </div>
    </aside>

    <section class="space-y-6">
      <header class="rounded-3xl border border-white/10 bg-gradient-to-br from-slate-950 via-slate-950 to-indigo-950 p-7 md:p-10">
        <div class="text-xs text-slate-300">Happy</div>
        <h1 class="mt-2 text-3xl md:text-4xl font-extrabold tracking-tight">Документация</h1>
        <p class="mt-3 text-slate-200/90 max-w-2xl">
          Коротко и по делу: как создать поздравление, опубликовать его на 7 дней и поделиться ссылкой.
        </p>
      </header>

      <article id="start" class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
        <h2 class="text-2xl font-bold">Как начать</h2>
        <div class="mt-4 grid md:grid-cols-3 gap-3 text-sm text-slate-300">
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="text-xs text-slate-400">Шаг 1</div>
            <div class="font-semibold mt-1 text-slate-100">Создай</div>
            <div class="mt-2">Выбери тему, напиши заголовок и основной текст.</div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="text-xs text-slate-400">Шаг 2</div>
            <div class="font-semibold mt-1 text-slate-100">Проверь</div>
            <div class="mt-2">Предпросмотр доступен всем — можно довести до идеала.</div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="text-xs text-slate-400">Шаг 3</div>
            <div class="font-semibold mt-1 text-slate-100">Опубликуй</div>
            <div class="mt-2">Оплата включает 7 дней доступа, потом страница удаляется.</div>
          </div>
        </div>

        <div class="mt-5 flex flex-col sm:flex-row gap-3">
          <NuxtLink to="/auth/register" class="px-5 py-3 rounded-2xl bg-white text-slate-950 font-semibold hover:bg-slate-100">
            Начать
          </NuxtLink>
          <NuxtLink to="/" class="px-5 py-3 rounded-2xl bg-white/10 border border-white/10 hover:bg-white/15">
            На главную
          </NuxtLink>
        </div>
      </article>

      <article id="types" class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
        <h2 class="text-2xl font-bold">2 услуги: по коду и по поддомену</h2>
        <div class="mt-4 grid md:grid-cols-2 gap-4 text-sm text-slate-300">
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-5">
            <div class="text-xs text-slate-400">Услуга 1(Разовая оплата страницы на 7 дней - 59 рублей)</div>
            <div class="text-lg font-semibold text-slate-100 mt-1">По коду</div>
            <ul class="mt-3 space-y-2">
              <li>• Открывается по вводу на главном экране в поле "Введите код и откройте страницу"</li>
              <li>• Быстрее и дешевле</li>
            </ul>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-5">
            <div class="text-xs text-slate-400">Услуга 2(Разовая оплата страницы на 7 дней - 99 рублей)</div>
            <div class="text-lg font-semibold text-slate-100 mt-1">По поддомену</div>
            <ul class="mt-3 space-y-2">
              <li>• <span class="font-mono">имя.{{ baseDomain }}</span></li>
              <li>• Максимальный “вау” эффект</li>
            </ul>
          </div>
        </div>
      </article>

      <article id="builder" class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
        <h2 class="text-2xl font-bold">Редактор</h2>
        <div class="mt-4 text-sm text-slate-300 space-y-3">
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="font-semibold text-slate-100">Заголовок</div>
            <div class="mt-1">Показывается крупно и жирно.</div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="font-semibold text-slate-100">Основной текст</div>
            <div class="mt-1">Поддерживает переносы строк, текст читаемый на мобильных.</div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="font-semibold text-slate-100">Подарок</div>
            <div class="mt-1">Нажатие сопровождается эффектом конфетти. Доступные варианты: промокод / текст / переход на сайт.</div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
            <div class="font-semibold text-slate-100">Тема и эффекты</div>
            <div class="mt-1">Новый год — снег и гирлянды, День рождения — шарики и конфетти и т.д.</div>
          </div>
        </div>
      </article>

      <article id="photos" class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
        <h2 class="text-2xl font-bold">Фото и просмотр</h2>
        <div class="mt-4 text-sm text-slate-300 space-y-2">
          <div>• Загружай несколько фото (jpg/png/webp).</div>
          <div>• Фото собираются в “карточки” и анимируются.</div>
          <div>• По клику фото раскрывается на весь экран прямо на странице.</div>
        </div>
      </article>

      <article id="publish" class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
        <h2 class="text-2xl font-bold">Оплата и 7 дней</h2>
        <div class="mt-4 text-sm text-slate-300 space-y-2">
          <div>• Создание и предпросмотр доступны бесплатно.</div>
          <div>• Публикация активна 7 дней после оплаты.</div>
          <div>• После окончания срока контент удаляется автоматически.</div>
        </div>
      </article>

      <article id="security" class="rounded-3xl border border-white/10 bg-white/5 p-6 md:p-8">
        <h2 class="text-2xl font-bold">Вход и восстановление</h2>
        <div class="mt-4 text-sm text-slate-300 space-y-2">
          <div>• Регистрация по email + username.</div>
          <div>• При регистрации выдаётся одноразовый резервный код для восстановления.</div>
        </div>
      </article>

      <div class="lg:hidden rounded-2xl border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
        <div class="font-semibold text-slate-100">Юридические документы</div>
        <div class="mt-2 flex flex-col gap-2">
          <NuxtLink class="hover:underline" to="/legal/offer">Оферта</NuxtLink>
          <NuxtLink class="hover:underline" to="/legal/company">Юр. лицо</NuxtLink>
          <NuxtLink class="hover:underline" to="/legal/privacy">Политика ПДн</NuxtLink>
          <NuxtLink class="hover:underline" to="/legal/consent">Согласие ПДн</NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const baseDomain = useRuntimeConfig().public.baseDomain
</script>
