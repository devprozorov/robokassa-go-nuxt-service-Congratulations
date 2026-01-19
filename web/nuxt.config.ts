export default defineNuxtConfig({
  ssr: true,
  css: ['~/assets/css/tailwind.css'],
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  app: {
    head: {
      title: 'Happy',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Happy — персональные поздравления по коду или поддомену' },
      ],
       link: [
        { rel: 'icon', type: 'image/png', href: '/happy4u.png' }, // иконка
      ]
    },
  },
  runtimeConfig: {
    public: {
      baseDomain: process.env.NUXT_PUBLIC_BASE_DOMAIN || 'happy4u.online',
      appName: process.env.NUXT_PUBLIC_APP_NAME || 'Happy',
    },
  },
})
