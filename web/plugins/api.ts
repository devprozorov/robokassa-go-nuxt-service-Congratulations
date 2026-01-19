export default defineNuxtPlugin(() => {
  const api = $fetch.create({
    baseURL: '/api',
    credentials: 'include',
    onRequest({ options }) {
      // SSR: прокидываем cookie
      if (import.meta.server) {
        const h = useRequestHeaders(['cookie'])
        options.headers = { ...(options.headers as any), ...h }
      }
    },
  })

  return { provide: { api } }
})
