import { http } from './http'

export interface CreateReportPayload {
  target_type: 'article' | 'forum_post'
  target_id: number
  reason: string
  description?: string
  evidence?: string
}

export class ReportApi {
  static async createReport(data: CreateReportPayload) {
    return http.post('/reports', data)
  }

  static async myReports(params: { page?: number; size?: number } = {}) {
    return http.get('/reports/my', params)
  }
}

export default ReportApi

