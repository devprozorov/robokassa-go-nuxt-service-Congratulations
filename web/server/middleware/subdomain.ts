import { defineEventHandler, getHeader, sendRedirect, setHeader } from 'h3'

export default defineEventHandler((event) => {
  const cfg = useRuntimeConfig(event)

  const host = ((getHeader(event, 'host') || '').toLowerCase().split(':')[0])
  const base = ((cfg.public.baseDomain || '').toLowerCase().split(':')[0])

  const url = event.node.req.url || '/'

  // don't redirect static assets or api
  if (
    url.startsWith('/api') ||
    url.startsWith('/_nuxt') ||
    url.startsWith('/__nuxt') ||
    url.startsWith('/favicon') ||
    url.startsWith('/robots.txt') ||
    url.startsWith('/s/')
  ) return

  if (!host || !base) return
  if (host === base) return
  if (!host.endsWith('.' + base)) return

  const sub = host.slice(0, -(base.length + 1))
  if (!sub || sub.includes('.') || sub === 'www') return

  // helpful for debug
  setHeader(event, 'x-happy-sub', sub)

  const qIndex = url.indexOf('?')
  const query = qIndex >= 0 ? url.slice(qIndex) : ''

  return sendRedirect(event, `/s/${encodeURIComponent(sub)}${query}`, 302)
})
