// web/middleware/panel-auth.ts
//
// Задача middleware:
// 1) На SSR (и после F5) Pinia-store пустой, поэтому мы обязаны спросить бэкенд: "кто я?"
// 2) Если пользователь не авторизован — отправляем на страницу входа и сохраняем куда он хотел попасть.
// 3) Если авторизован — пропускаем дальше.
//
// ВАЖНО про SSR:
// - В SSR запросы $fetch НЕ прокидывают cookie сами, поэтому fetchMe() должен прокидывать cookie через headers
//   (ты это уже поправил в auth store).
// - Здесь обязательно использовать await, иначе редирект случится раньше, чем fetchMe успеет заполнить store.

import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()

  if (!auth.user) {
    await auth.fetchMe()
  }

  if (!auth.user) {
    return navigateTo(`/auth?next=${encodeURIComponent(to.fullPath)}`)
  }
})
