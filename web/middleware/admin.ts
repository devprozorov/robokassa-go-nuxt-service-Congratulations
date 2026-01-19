export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()
  if (!auth.user) await auth.fetchMe()
  if (!auth.user || auth.user.role !== 'admin') {
    return navigateTo('/')
  }
})
