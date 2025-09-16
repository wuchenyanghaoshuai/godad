// 积分系统API接口
import { http } from './http'

// 积分相关接口类型定义
export interface UserPoints {
  id: number
  user_id: number
  total_points: number
  current_level: number
  next_level_points: number
  created_at: string
  updated_at: string
  user: {
    id: number
    username: string
    nickname?: string
    avatar?: string
  }
  level: UserLevel
}

export interface UserLevel {
  id: number
  name: string
  level: number
  min_points: number
  max_points: number
  color: string
  icon: string
  badge?: string
  description: string
  privileges?: string
  status: number
  created_at: string
  updated_at: string
}

export interface PointsTransaction {
  id: number
  user_id: number
  action: string
  points: number
  description: string
  source_type: string
  source_id: number
  created_at: string
  user: {
    id: number
    username: string
    nickname?: string
  }
}

export interface PointsRule {
  id: number
  action: string
  name: string
  points: number
  daily_limit: number
  description: string
  status: number
  created_at: string
  updated_at: string
}

export interface PointsStats {
  total_points: number
  current_level: number
  next_level_points: number
  today_points: number
  rank: number
  level_info: UserLevel
}

export interface PointsHistoryResponse {
  transactions: PointsTransaction[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

// 积分系统API类
export class PointsAPI {
  // 获取用户积分信息
  static async getUserPoints() {
    return http.get<UserPoints>('/points/user')
  }

  // 获取积分历史记录
  static async getPointsHistory(params: { page?: number; limit?: number } = {}) {
    return http.get<PointsHistoryResponse>('/points/history', {
      page: params.page || 1,
      limit: params.limit || 20
    })
  }

  // 获取积分统计信息
  static async getPointsStats() {
    return http.get<PointsStats>('/points/stats')
  }

  // 获取等级配置列表
  static async getLevels() {
    return http.get<UserLevel[]>('/points/levels')
  }

  // 获取积分规则列表
  static async getPointsRules() {
    return http.get<PointsRule[]>('/points/rules')
  }
}

// 导出积分相关API方法
export const {
  getUserPoints,
  getPointsHistory,
  getPointsStats,
  getLevels,
  getPointsRules
} = PointsAPI