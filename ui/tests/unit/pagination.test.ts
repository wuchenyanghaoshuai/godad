import { describe, it, expect } from 'vitest'
import { normalizePageResponse } from '@/api/pagination'

describe('normalizePageResponse', () => {
  it('handles flat array with total/page/size', () => {
    const raw = { code: 200, message: 'success', data: [1, 2], total: 100, page: 2, size: 10 }
    const res = normalizePageResponse<number>(raw)
    expect(res.items).toEqual([1, 2])
    expect(res.total).toBe(100)
    expect(res.page).toBe(2)
    expect(res.size).toBe(10)
    expect(res.total_pages).toBe(10)
  })

  it('handles nested items object', () => {
    const raw = { code: 200, data: { items: [{ a: 1 }], total: 1, page: 1, size: 20, total_pages: 1 } }
    const res = normalizePageResponse<any>(raw)
    expect(res.items.length).toBe(1)
    expect(res.total).toBe(1)
    expect(res.page).toBe(1)
    expect(res.size).toBe(20)
    expect(res.total_pages).toBe(1)
  })

  it('handles notifications shape', () => {
    const raw = { code: 200, data: { notifications: [{ id: 1 }], pagination: { total: 5, current_page: 1, per_page: 10, total_pages: 1 } } }
    const res = normalizePageResponse<any>(raw)
    expect(res.items.length).toBe(1)
    expect(res.total).toBe(5)
    expect(res.page).toBe(1)
    expect(res.size).toBe(10)
    expect(res.total_pages).toBe(1)
  })
})

