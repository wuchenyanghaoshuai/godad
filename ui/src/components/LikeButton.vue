<template>
  <button 
    :class="[
      'like-button',
      { 
        'like-button--liked': isLiked,
        'like-button--loading': loading,
        'like-button--disabled': disabled
      }
    ]"
    :disabled="disabled || loading"
    @click="toggleLike"
  >
    <div class="like-icon">
      <svg 
        v-if="!loading"
        :class="{ 'liked': isLiked }"
        width="20" 
        height="20" 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2"
      >
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
      </svg>
      <div v-else class="loading-spinner">
        <div class="spinner"></div>
      </div>
    </div>
    <span v-if="showCount" class="like-count">
      {{ formatCount(likeCount) }}
    </span>
  </button>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import likeApi from '@/api/likes'
import { useToast } from '@/composables/useToast'
import { useAuthStore } from '@/stores/auth'

interface Props {
  targetType: 'article' | 'comment'
  targetId: number
  initialLiked?: boolean
  initialCount?: number
  showCount?: boolean
  size?: 'small' | 'medium' | 'large'
}

interface Emits {
  (e: 'update:liked', liked: boolean): void
  (e: 'update:count', count: number): void
}

const props = withDefaults(defineProps<Props>(), {
  initialLiked: false,
  initialCount: 0,
  showCount: true,
  size: 'medium'
})

const emit = defineEmits<Emits>()

const authStore = useAuthStore()
const { toast } = useToast()

const isLiked = ref(props.initialLiked)
const likeCount = ref(props.initialCount)
const loading = ref(false)

const disabled = computed(() => !authStore.isAuthenticated)

const toggleLike = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('请先登录')
    return
  }

  if (loading.value) return

  const originalLiked = isLiked.value
  const originalCount = likeCount.value

  // 乐观更新
  isLiked.value = !isLiked.value
  likeCount.value = isLiked.value ? originalCount + 1 : Math.max(0, originalCount - 1)

  emit('update:liked', isLiked.value)
  emit('update:count', likeCount.value)

  try {
    loading.value = true
    
    const response = await likeApi.toggleLike(props.targetType, props.targetId)
    
    // 根据返回的数据判断点赞状态
    // 如果返回数据不为null，说明是点赞操作；如果为null，说明是取消点赞操作
    if (response.data) {
      // 点赞成功
      isLiked.value = true
      likeCount.value = originalCount + 1
    } else {
      // 取消点赞成功
      isLiked.value = false
      likeCount.value = Math.max(0, originalCount - 1)
    }

    emit('update:liked', isLiked.value)
    emit('update:count', likeCount.value)
    
    toast.success(response.message)
  } catch (error: any) {
    
    // 回滚乐观更新
    isLiked.value = originalLiked
    likeCount.value = originalCount

    emit('update:liked', isLiked.value)
    emit('update:count', likeCount.value)

    toast.error(error.response?.data?.message || error.message || '操作失败，请重试')
  } finally {
    loading.value = false
  }
}

const formatCount = (count: number): string => {
  if (count < 1000) return count.toString()
  if (count < 10000) return `${Math.floor(count / 1000)}k`
  return `${Math.floor(count / 10000)}w`
}

const loadLikeStatus = async () => {
  if (!authStore.isAuthenticated) return

  try {
    const response = await likeApi.getLikeStatus(props.targetType, props.targetId)
    isLiked.value = response.data.is_liked
    likeCount.value = response.data.like_count

    emit('update:liked', isLiked.value)
    emit('update:count', likeCount.value)
  } catch (error) {
    console.error('加载点赞状态失败:', error)
  }
}

// 监听认证状态变化
watch(() => authStore.isAuthenticated, (newVal) => {
  if (newVal) {
    loadLikeStatus()
  } else {
    isLiked.value = false
  }
})

// 监听目标ID变化
watch(() => props.targetId, () => {
  isLiked.value = props.initialLiked
  likeCount.value = props.initialCount
  
  if (authStore.isAuthenticated) {
    loadLikeStatus()
  }
})

onMounted(() => {
  if (authStore.isAuthenticated) {
    loadLikeStatus()
  }
})
</script>

<style scoped>
.like-button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 20px;
  font-size: 14px;
  color: #64748b;
  transition: all 0.2s ease;
  user-select: none;
}

.like-button:hover:not(.like-button--disabled):not(.like-button--loading) {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.like-button--liked {
  color: #ef4444;
}

.like-button--liked:hover {
  background: rgba(239, 68, 68, 0.1);
}

.like-button--loading {
  cursor: not-allowed;
  opacity: 0.7;
}

.like-button--disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.like-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.like-icon svg {
  transition: all 0.2s ease;
}

.like-icon svg.liked {
  fill: currentColor;
  color: #ef4444;
  animation: likeAnimation 0.3s ease;
}

@keyframes likeAnimation {
  0% { transform: scale(1); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.loading-spinner {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e2e8f0;
  border-top: 2px solid #ef4444;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.like-count {
  font-weight: 500;
  font-size: 14px;
}

/* 尺寸变体 */
.like-button[data-size="small"] {
  padding: 4px 8px;
  font-size: 12px;
}

.like-button[data-size="small"] .like-icon svg {
  width: 16px;
  height: 16px;
}

.like-button[data-size="large"] {
  padding: 12px 16px;
  font-size: 16px;
}

.like-button[data-size="large"] .like-icon svg {
  width: 24px;
  height: 24px;
}
</style>