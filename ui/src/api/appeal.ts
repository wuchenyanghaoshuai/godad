import { http } from './http'

export interface CreateAppealPayload {
  report_id?: number
  target_id?: number
  reason: string
  evidence?: string
}

export const AppealApi = {
  create(data: CreateAppealPayload) {
    return http.post('/appeals', data)
  },
}

export default AppealApi

