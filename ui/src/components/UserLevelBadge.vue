<template>
  <div
    class="user-level-badge inline-flex items-center"
    :class="sizeClasses"
    :title="tooltip"
  >
    <!-- 等级图标 -->
    <span
      class="level-icon flex-shrink-0"
      :class="iconSizeClasses"
      :style="{ color: levelData?.color || '#9CA3AF' }"
    >
      {{ levelData?.icon || '🌱' }}
    </span>

    <!-- 等级信息 -->
    <div v-if="showText" class="level-info ml-1 flex flex-col">
      <span
        class="level-name font-medium leading-tight"
        :class="textSizeClasses"
        :style="{ color: levelData?.color || '#9CA3AF' }"
      >
        {{ levelData?.name || 'LV1' }}
      </span>
      <span
        v-if="showPoints && size !== 'xs'"
        class="level-points text-xs text-gray-500 leading-tight"
      >
        {{ points }}积分
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { PointsAPI, type UserLevel } from '@/api'

// Props
interface Props {
  level?: number
  points?: number
  size?: 'xs' | 'sm' | 'md' | 'lg'
  showText?: boolean
  showPoints?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  level: 1,
  points: 0,
  size: 'md',
  showText: true,
  showPoints: false
})

// 响应式数据
const levels = ref<UserLevel[]>([])
const loading = ref(false)

// 计算属性
const levelData = computed(() => {
  return levels.value.find(l => l.level === props.level) || {
    id: 1,
    name: 'LV1',
    level: 1,
    min_points: 0,
    max_points: 99,
    color: '#9CA3AF',
    icon: '🌱',
    description: '新手用户',
    status: 1,
    created_at: '',
    updated_at: ''
  }
})

const sizeClasses = computed(() => {
  const sizeMap = {
    xs: 'text-xs',
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-lg'
  }
  return sizeMap[props.size]
})

const iconSizeClasses = computed(() => {
  const sizeMap = {
    xs: 'text-sm',
    sm: 'text-base',
    md: 'text-lg',
    lg: 'text-xl'
  }
  return sizeMap[props.size]
})

const textSizeClasses = computed(() => {
  const sizeMap = {
    xs: 'text-xs',
    sm: 'text-sm',
    md: 'text-sm',
    lg: 'text-base'
  }
  return sizeMap[props.size]
})

const tooltip = computed(() => {
  const level = levelData.value
  if (!level) return ''

  let tooltipText = `${level.name} (LV${level.level})`
  if (level.description) {
    tooltipText += `\n${level.description}`
  }
  if (props.showPoints && props.points !== undefined) {
    tooltipText += `\n当前积分：${props.points}`
    if (level.max_points < 999999999) {
      tooltipText += `\n升级需要：${level.max_points + 1}积分`
    }
  }
  return tooltipText
})

// 方法
const loadLevels = async () => {
  try {
    loading.value = true
    const response = await PointsAPI.getLevels()
    levels.value = response.data
  } catch (error) {
    console.error('加载等级配置失败:', error)
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  loadLevels()
})
</script>

<style scoped>
.user-level-badge {
  @apply rounded-md px-1 py-0.5;
  transition: all 0.2s ease-in-out;
}

.user-level-badge:hover {
  @apply shadow-sm;
  transform: translateY(-1px);
}

.level-icon {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.level-name {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}
</style>