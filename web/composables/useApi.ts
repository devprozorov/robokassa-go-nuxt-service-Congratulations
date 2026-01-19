export function useApi() {
  const base = '/api'
  return {
    async get<T>(path: string): Promise<T> {
      return await $fetch<T>(base + path, { credentials: 'include' })
    },
    async post<T>(path: string, body?: any): Promise<T> {
      return await $fetch<T>(base + path, { method: 'POST', body, credentials: 'include' })
    },
    async put<T>(path: string, body?: any): Promise<T> {
      return await $fetch<T>(base + path, { method: 'PUT', body, credentials: 'include' })
    },
    async del<T>(path: string): Promise<T> {
      return await $fetch<T>(base + path, { method: 'DELETE', credentials: 'include' })
    },
  }
}
