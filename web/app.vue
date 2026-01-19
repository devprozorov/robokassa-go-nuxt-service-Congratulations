<template>
  <div class="min-h-screen bg-slate-950 text-slate-100">
    <AppHeader v-if="!isGreetingRoute" />

    <main :class="mainClass">
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
    </main>

    <AppFooter v-if="!isGreetingRoute" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from '~/components/AppHeader.vue'
import AppFooter from '~/components/AppFooter.vue'

const route = useRoute()

const isGreetingRoute = computed(() => {
  const p = route.path || '/'
  return p.startsWith('/s/') || p.startsWith('/c/') || p.startsWith('/preview/')
})

const mainClass = computed(() => {
  if (isGreetingRoute.value) return 'p-0'
  return 'max-w-6xl mx-auto px-3 sm:px-4 md:px-6 py-5 sm:py-6'
})
</script>
