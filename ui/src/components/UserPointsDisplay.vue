<template>
  <div class="user-points-display">
    <!-- 简洁模式 -->
    <div v-if="mode === 'simple'" class="flex items-center space-x-2">
      <UserLevelBadge
        :level="userPoints?.current_level"
        :points="userPoints?.total_points"
        :size="size"
        :show-text="showLevel"
        :show-points="false"
      />
      <span
        v-if="showPoints"
        class="points-text font-medium text-gray-700"
        :class="textSizeClasses"
      >
        {{ formatPoints(userPoints?.total_points || 0) }}
      </span>
    </div>

    <!-- 详细模式 -->
    <div v-else-if="mode === 'detailed'" class="space-y-4">
      <!-- 当前积分和等级卡片 -->
      <!-- 采用与用户中心统计卡一致的浅色渐变样式 -->
      <div class="bg-gradient-to-r from-primary-50 to-secondary-50 rounded-lg p-4 border border-primary-100">
        <div class="flex items-center justify-between mb-3">
          <div>
            <div class="text-2xl font-bold text-gray-900">
              {{ userPoints?.total_points || 0 }} 积分
            </div>
            <div class="text-sm text-gray-600">当前积分总数</div>
          </div>
          <UserLevelBadge
            :level="userPoints?.current_level"
            :points="userPoints?.total_points"
            size="lg"
            :show-text="true"
            :show-points="false"
          />
        </div>

        <!-- 等级进度信息 -->
        <div v-if="progressData.showProgress" class="space-y-2">
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-700">
              {{ progressData.currentLevelName }} → {{ progressData.nextLevelName }}
            </span>
            <span class="font-medium text-primary-600">
              还需 {{ progressData.needPoints }} 积分升级
            </span>
          </div>
          <!-- 进度条 -->
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div
              class="bg-gradient-to-r from-primary-500 to-secondary-500 h-2 rounded-full transition-all duration-300"
              :style="{ width: `${progressData.progressPercent}%` }"
            />
          </div>
          <div class="flex justify-between text-xs text-gray-500">
            <span>{{ formatPoints(progressData.currentLevelMin) }}</span>
            <span>{{ formatPoints(progressData.nextLevelMin) }}</span>
          </div>
        </div>

        <!-- 已达到最高等级 -->
        <div v-else class="text-center py-2">
          <span class="text-sm font-medium text-primary-600">🎉 已达到最高等级！</span>
        </div>
      </div>

      <!-- 今日积分 -->
      <div v-if="showTodayPoints && stats?.today_points" class="flex items-center justify-between p-2 bg-secondary-50 rounded-lg">
        <span class="text-sm text-secondary-700">今日获得</span>
        <span class="font-medium text-secondary-800">+{{ stats.today_points }}</span>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center space-x-2">
      <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-500"></div>
      <span class="text-sm text-gray-500">加载中...</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { PointsAPI, type UserPoints, type PointsStats, type UserLevel } from '@/api'
import UserLevelBadge from './UserLevelBadge.vue'

// Props
interface Props {
  mode?: 'simple' | 'detailed'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  showLevel?: boolean
  showPoints?: boolean
  showProgress?: boolean
  showTodayPoints?: boolean
  autoRefresh?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'simple',
  size: 'md',
  showLevel: true,
  showPoints: true,
  showProgress: true,
  showTodayPoints: true,
  autoRefresh: false
})

// 响应式数据
const userPoints = ref<UserPoints | null>(null)
const stats = ref<PointsStats | null>(null)
const levels = ref<UserLevel[]>([])
const loading = ref(false)
const error = ref<string>('')

// 计算属性
const textSizeClasses = computed(() => {
  const sizeMap = {
    xs: 'text-xs',
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-lg'
  }
  return sizeMap[props.size]
})

const progressData = computed(() => {
  if (!userPoints.value || !levels.value.length) {
    return {
      showProgress: false,
      progressPercent: 0,
      needPoints: 0,
      currentLevelMin: 0,
      nextLevelMin: 0,
      currentLevelName: '',
      nextLevelName: ''
    }
  }

  const currentLevel = levels.value.find(l => l.level === userPoints.value!.current_level)
  const nextLevel = levels.value.find(l => l.level === userPoints.value!.current_level + 1)

  if (!currentLevel || !nextLevel) {
    return {
      showProgress: false,
      progressPercent: 100,
      needPoints: 0,
      currentLevelMin: currentLevel?.min_points || 0,
      nextLevelMin: currentLevel?.max_points || 0,
      currentLevelName: currentLevel?.name || '',
      nextLevelName: ''
    }
  }

  const currentPoints = userPoints.value.total_points
  const currentLevelPoints = currentPoints - currentLevel.min_points
  const totalNeedPoints = nextLevel.min_points - currentLevel.min_points
  const progressPercent = Math.min((currentLevelPoints / totalNeedPoints) * 100, 100)
  const needPoints = nextLevel.min_points - currentPoints

  return {
    showProgress: true,
    progressPercent,
    needPoints: Math.max(needPoints, 0),
    currentLevelMin: currentLevel.min_points,
    nextLevelMin: nextLevel.min_points,
    currentLevelName: currentLevel.name,
    nextLevelName: nextLevel.name
  }
})

// 方法
const formatPoints = (points: number): string => {
  if (points >= 10000) {
    return `${(points / 10000).toFixed(1)}万`
  }
  if (points >= 1000) {
    return `${(points / 1000).toFixed(1)}k`
  }
  return points.toString()
}

const loadUserPoints = async () => {
  try {
    loading.value = true
    error.value = ''

    const [pointsResponse, levelsResponse] = await Promise.all([
      PointsAPI.getUserPoints(),
      PointsAPI.getLevels()
    ])

    userPoints.value = pointsResponse.data
    levels.value = levelsResponse.data

    // 如果是详细模式，也加载统计信息
    if (props.mode === 'detailed') {
      const statsResponse = await PointsAPI.getPointsStats()
      stats.value = statsResponse.data
    }
  } catch (err: any) {
    error.value = err.message || '加载用户积分失败'
    console.error('加载用户积分失败:', err)
  } finally {
    loading.value = false
  }
}

const refresh = () => {
  loadUserPoints()
}

// 暴露方法给父组件
defineExpose({
  refresh,
  userPoints,
  stats
})

// 生命周期
onMounted(() => {
  loadUserPoints()

  // 自动刷新
  if (props.autoRefresh) {
    const interval = setInterval(() => {
      loadUserPoints()
    }, 30000) // 30秒刷新一次

    // 组件卸载时清理定时器
    onBeforeUnmount(() => {
      clearInterval(interval)
    })
  }
})
</script>

<style scoped>
.user-points-display {
  @apply select-none;
}

.points-text {
  @apply bg-gradient-to-r from-primary-600 to-secondary-600 bg-clip-text text-transparent;
}
</style>
