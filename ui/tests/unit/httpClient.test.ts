import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { HttpClient } from '@/api/http'

// Mock errorHandler to avoid side effects
vi.mock('@/utils/errorHandler', () => {
  return {
    errorHandler: {
      handleApiError: (err: any) => ({ isAppError: true, code: 'API_ERROR', message: err?.message || 'error', details: err, timestamp: Date.now() }),
      handleBusinessError: (code: string, message?: string, details?: any) => ({ isAppError: true, code, message: message || 'biz', details, timestamp: Date.now() })
    }
  }
})

// Helper to create Response
const jsonResponse = (body: any, init?: Partial<Response>) =>
  new Response(JSON.stringify(body), { status: init?.status ?? 200, headers: { 'Content-Type': 'application/json' } })

describe('HttpClient', () => {
  const client = new HttpClient()
  const originalFetch = globalThis.fetch

  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
    globalThis.fetch = originalFetch
  })

  it('sends credentials and parses success code', async () => {
    const fetchMock = vi.fn().mockImplementation((_url: any, init?: RequestInit) => {
      expect(init?.credentials).toBe('include')
      return Promise.resolve(jsonResponse({ code: 200, message: 'ok', data: { ok: true } }))
    })
    // @ts-expect-error override
    globalThis.fetch = fetchMock

    const res = await client.get<any>('/test')
    expect(res.data.ok).toBe(true)
    expect(fetchMock).toHaveBeenCalled()
  })

  it('throws business error when code is non-success', async () => {
    // @ts-expect-error override
    globalThis.fetch = vi.fn().mockResolvedValue(jsonResponse({ code: 400, message: 'bad', data: null }))
    await expect(client.get<any>('/bad')).rejects.toMatchObject({ isAppError: true })
  })

  it('refreshes on 401 then retries original request', async () => {
    const calls: string[] = []
    // @ts-expect-error override
    globalThis.fetch = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
      const u = url.toString()
      calls.push(u)
      if (u.includes('/auth/refresh-token')) {
        return Promise.resolve(jsonResponse({ code: 200, message: 'ok' }))
      }
      // First attempt returns 401
      if (calls.filter(x => !x.includes('/auth/refresh-token')).length === 1) {
        return Promise.resolve(new Response('', { status: 401, statusText: 'Unauthorized' }))
      }
      // Retry returns success
      return Promise.resolve(jsonResponse({ code: 200, message: 'ok', data: { retried: true } }))
    })

    const res = await client.get<any>('/needs-auth')
    expect(res.data.retried).toBe(true)
    expect(calls.some(u => u.includes('/auth/refresh-token'))).toBe(true)
  })
})

