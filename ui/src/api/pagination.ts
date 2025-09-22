import type { PaginatedResponse } from './types'

// 统一分页响应适配器，将后端不同结构统一为 { items, total, page, size, total_pages }
export function normalizePageResponse<T = any>(raw: any): PaginatedResponse<T> {
  // 可能的形态：
  // 1) { code, message, data: [...], total, page, size }
  if (raw && Array.isArray(raw.data) && typeof raw.total !== 'undefined') {
    const totalPages = Math.ceil((raw.total || 0) / (raw.size || 1))
    return {
      items: raw.data as T[],
      total: Number(raw.total) || 0,
      page: Number(raw.page) || 1,
      size: Number(raw.size) || (raw.data?.length || 0),
      total_pages: Number(raw.pages) || totalPages,
    }
  }

  // 2) { code, message, data: { items: [...], total, page, size, total_pages } } 或含 page_size/total_page
  if (raw && raw.data && Array.isArray(raw.data.items)) {
    const d = raw.data
    return {
      items: d.items as T[],
      total: Number(d.total) || 0,
      page: Number(d.page) || 1,
      size: Number(d.size) || Number(d.page_size) || (d.items?.length || 0),
      total_pages: Number(d.total_pages) || Number(d.pages) || Number(d.total_page) || 1,
    }
  }

  // 3) 通知：{ code, data: { notifications: [...], pagination: {...} } }
  if (raw && raw.data && Array.isArray(raw.data.notifications)) {
    const list = raw.data.notifications as T[]
    const p = raw.data.pagination || {}
    return {
      items: list,
      total: Number(p.total) || list.length || 0,
      page: Number(p.current_page) || Number(p.page) || 1,
      size: Number(p.per_page) || Number(p.size) || list.length || 0,
      total_pages: Number(p.total_pages) || 1,
    }
  }

  // 回退：尽力从常见字段推断
  const items: T[] = Array.isArray(raw?.data?.data) ? raw.data.data
                    : Array.isArray(raw?.data) ? raw.data
                    : Array.isArray(raw) ? raw
                    : []
  return {
    items,
    total: Number(raw?.total) || Number(raw?.data?.total) || items.length,
    page: Number(raw?.page) || Number(raw?.data?.page) || 1,
    size: Number(raw?.size) || Number(raw?.data?.size) || items.length,
    total_pages: Number(raw?.total_pages) || Number(raw?.data?.total_pages) || 1,
  }
}
