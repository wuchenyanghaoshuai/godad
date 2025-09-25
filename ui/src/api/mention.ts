import { http } from './http'

export interface MentionUser {
  id: number
  username: string
  nickname: string
  avatar?: string
}

export const MentionApi = {
  suggest: (query?: string) => {
    const params: any = {}
    if (query && query.trim()) params.query = query.trim()
    return http.get<MentionUser[]>(`/mentions/suggestions`, params)
  }
}

