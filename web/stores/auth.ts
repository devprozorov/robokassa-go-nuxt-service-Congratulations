import { defineStore } from 'pinia'

type User = { id: string; email: string; username: string; role: 'user' | 'admin' }

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    loading: false,
    error: '' as string,
  }),
  getters: {
    isAuthed: (s) => !!s.user,
    isAdmin: (s) => s.user?.role === 'admin',
  },
  actions: {
    async fetchMe() {
      try {
        const headers = import.meta.server ? useRequestHeaders(['cookie']) : undefined
        const res: any = await $fetch('/api/auth/me', { credentials: 'include', headers })
        this.user = res?.user ?? null
      } catch {
        this.user = null
      }
    },

    async login(loginOrEmail: string, password: string) {
      this.error = ''
      this.loading = true
      try {
        const res: any = await $fetch('/api/auth/login', {
          method: 'POST',
          body: { loginOrEmail, password },
          credentials: 'include',
        })
        if (res?.ok) {
          this.user = res.user
          return { ok: true }
        }
        this.error = res?.error || 'LOGIN_FAILED'
        return { ok: false }
      } catch {
        this.error = 'LOGIN_FAILED'
        return { ok: false }
      } finally {
        this.loading = false
      }
    },

    async register(email: string, username: string, password: string) {
      this.error = ''
      this.loading = true
      try {
        const res: any = await $fetch('/api/auth/register', {
          method: 'POST',
          body: { email, username, password },
          credentials: 'include',
        })
        if (res?.ok) {
          this.user = res.user
          return { ok: true, recoveryCode: res.recoveryCode as string }
        }
        this.error = res?.error || 'REGISTER_FAILED'
        return { ok: false }
      } catch {
        this.error = 'REGISTER_FAILED'
        return { ok: false }
      } finally {
        this.loading = false
      }
    },

    async logout() {
      await $fetch('/api/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => {})
      this.user = null
    },
  },
})
