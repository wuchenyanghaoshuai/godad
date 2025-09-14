<template>
  <div class="min-h-screen">
    <!-- 导航栏 -->
    <Navbar />
    
    <div class="user-profile-page">
      <div class="container">
      <!-- 用户信息头部 -->
      <div class="profile-header">
        <div class="profile-cover">
          <img 
            v-if="profile?.avatar" 
            :src="profile.avatar" 
            alt="用户头像"
            class="profile-avatar"
          />
          <div v-else class="profile-avatar default-avatar">
            {{ profile?.username?.charAt(0)?.toUpperCase() || 'U' }}
          </div>
          
          <div class="profile-info">
            <h1 class="profile-name">
              {{ profile?.nickname || profile?.username || '加载中...' }}
            </h1>
            <p class="profile-username">@{{ profile?.username || '' }}</p>
            <p v-if="profile?.bio" class="profile-bio">{{ profile.bio }}</p>
            
            <div class="profile-meta">
              <div v-if="profile?.location" class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                  <circle cx="12" cy="10" r="3"/>
                </svg>
                {{ profile.location }}
              </div>
              <div v-if="profile?.website" class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
                </svg>
                <a :href="profile.website" target="_blank" rel="noopener noreferrer">
                  {{ profile.website }}
                </a>
              </div>
              <div class="meta-item">
                <svg class="meta-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                  <line x1="16" y1="2" x2="16" y2="6"/>
                  <line x1="8" y1="2" x2="8" y2="6"/>
                  <line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
                加入于 {{ formatDate(profile?.created_at) }}
              </div>
            </div>
          </div>
          
          <div class="profile-actions">
            <div class="profile-manage-actions">
              <!-- 他人个人中心显示关注和私信按钮 -->
              <div v-if="!isOwnProfile" class="flex items-center space-x-3">
                <button
                  @click="handleFollowAction"
                  class="follow-btn"
                  :class="getFollowButtonClass()"
                  :disabled="followActionLoading || followStatusLoading"
                >
                  <svg v-if="followActionLoading" class="w-4 h-4 animate-spin mr-2" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
                  </svg>
                  <svg v-else-if="isFollowing" class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                  </svg>
                  {{ followActionLoading ? '处理中...' : getFollowButtonText() }}
                </button>
                
                <button
                  @click="startChat"
                  class="chat-btn"
                  :disabled="chatLoading"
                >
                  <svg v-if="chatLoading" class="w-4 h-4 animate-spin mr-2" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
                  </svg>
                  <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
                  </svg>
                  {{ chatLoading ? '处理中...' : '私信' }}
                </button>
              </div>
              
            </div>
          </div>
        </div>

        <!-- 统计信息 -->
        <div class="stats-container">
          <div class="stats-grid">
            <div class="stat-item" @click="showFollowingModal = true">
              <div class="stat-number">{{ stats?.following_count || 0 }}</div>
              <div class="stat-label">关注</div>
            </div>
            <div class="stat-item" @click="showFollowersModal = true">
              <div class="stat-number">{{ stats?.followers_count || 0 }}</div>
              <div class="stat-label">粉丝</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ stats?.article_count || 0 }}</div>
              <div class="stat-label">文章</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ stats?.like_count || 0 }}</div>
              <div class="stat-label">获赞</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 内容标签栏 -->
      <div class="content-tabs">
        <button 
          v-for="tab in tabs"
          :key="tab.key"
          :class="['tab-button', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- 内容区域 -->
      <div class="content-area">
        <div v-if="activeTab === 'articles'" class="articles-content">
          <div v-if="!Array.isArray(articles) || articles.length === 0" class="empty-state">
            <Empty message="还没有发布任何文章" />
          </div>
          <div v-else class="articles-grid">
            <!-- TODO: 集成ArticleCard组件 -->
            <div
              v-for="article in (Array.isArray(articles) ? articles.filter(a => a && a.id) : [])"
              :key="article.id"
              class="article-item bg-white rounded-lg shadow-sm border cursor-pointer hover:shadow-md transition-shadow overflow-hidden"
              @click="$router.push(`/articles/${article.id}`)"
            >
              <!-- 封面图片 -->
              <div v-if="article.cover_image" class="relative h-48 bg-gray-100">
                <img 
                  :src="article.cover_image" 
                  :alt="article.title"
                  class="w-full h-full object-cover"
                  @error="$event.target.style.display='none'"
                />
              </div>
              
              <!-- 文章内容 -->
              <div class="p-4">
                <h3 class="font-medium text-gray-900 hover:text-blue-600 transition-colors line-clamp-2">{{ article.title }}</h3>
                <p class="text-sm text-gray-600 mt-2 line-clamp-2">{{ article.summary }}</p>
                <div class="flex items-center justify-between mt-3 text-xs text-gray-500">
                  <span>{{ new Date(article.created_at).toLocaleDateString('zh-CN') }}</span>
                  <div class="flex items-center space-x-3">
                    <span>阅读 {{ article.view_count || 0 }}</span>
                    <span>点赞 {{ article.like_count || 0 }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'liked'" class="liked-content">
          <div v-if="likedArticles.length === 0" class="empty-state">
            <Empty message="还没有点赞任何文章" />
          </div>
          <div v-else class="articles-grid">
            <!-- TODO: 集成ArticleCard组件 -->
            <div 
              v-for="article in likedArticles"
              :key="article.id"
              class="article-item bg-white rounded-lg shadow-sm border cursor-pointer hover:shadow-md transition-shadow overflow-hidden"
              @click="$router.push(`/articles/${article.id}`)"
            >
              <!-- 封面图片 -->
              <div v-if="article.cover_image" class="relative h-48 bg-gray-100">
                <img 
                  :src="article.cover_image" 
                  :alt="article.title"
                  class="w-full h-full object-cover"
                  @error="$event.target.style.display='none'"
                />
              </div>
              
              <!-- 文章内容 -->
              <div class="p-4">
                <h3 class="font-medium text-gray-900 hover:text-blue-600 transition-colors line-clamp-2">{{ article.title }}</h3>
                <p class="text-sm text-gray-600 mt-2 line-clamp-2">{{ article.summary }}</p>
                <div class="flex items-center justify-between mt-3 text-xs text-gray-500">
                  <span>{{ new Date(article.created_at).toLocaleDateString('zh-CN') }}</span>
                  <div class="flex items-center space-x-3">
                    <span>阅读 {{ article.view_count || 0 }}</span>
                    <span>点赞 {{ article.like_count || 0 }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'activity'" class="activity-content">
          <div v-if="activities.length === 0" class="empty-state">
            <Empty message="暂无动态" />
          </div>
          <div v-else class="activity-list">
            <!-- TODO: 集成ActivityItem组件 -->
            <div 
              v-for="activity in activities"
              :key="activity.id"
              class="activity-item bg-white p-4 rounded-lg shadow-sm border mb-3"
            >
              <div class="font-medium text-gray-900">{{ activity.title }}</div>
              <div class="text-sm text-gray-600 mt-1">{{ activity.description }}</div>
            </div>
          </div>
        </div>

      </div>
    </div>

    </div>
  </div>

</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserDataSync } from '@/composables/useUserDataSync'
import { UserApi } from '@/api/user'
import { useToast } from '@/composables/useToast'
import { FollowApi } from '@/api/follow'
import Navbar from '@/components/Navbar.vue'
import Empty from '@/components/Empty.vue'
import type { User } from '@/api/types'
// import ArticleCard from '@/components/ArticleCard.vue'
// import ActivityItem from '@/components/ActivityItem.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

// 用户名参数
const username = computed(() => route.params.username as string)

// 使用数据同步组合函数
const userDataSync = useUserDataSync(username.value)

// 使用同步的数据 - 直接使用userDataSync返回的computed属性
const profile = userDataSync.profile
const stats = userDataSync.stats
const articles = userDataSync.articles
const likedArticles = ref<any[]>([])
const activities = ref<any[]>([])
const loading = computed(() => userDataSync.isLoading)
const activeTab = ref('articles')
const showFollowingModal = ref(false)
const showFollowersModal = ref(false)

const isFollowing = ref<boolean | null>(null) // null表示未加载
const isMutualFollow = ref<boolean | null>(null) // null表示未加载
const followActionLoading = ref(false)
const chatLoading = ref(false)

// 关注状态加载状态
const followStatusLoading = ref(false)

// 在用户资料加载后检查关注状态
const checkFollowStatus = async () => {
  if (!userId.value || isOwnProfile.value || !authStore.isAuthenticated) {
    isFollowing.value = false
    isMutualFollow.value = false
    return
  }

  // 如果已经在加载，避免重复请求
  if (followStatusLoading.value) return

  followStatusLoading.value = true

  try {
    // 检查我是否关注了对方
    const myFollowResponse = await FollowApi.checkFollowStatus(userId.value)
    console.log('关注状态API响应:', myFollowResponse)
    isFollowing.value = myFollowResponse.data?.is_following || myFollowResponse.is_following || false
    console.log('是否关注:', isFollowing.value)

    // 如果我关注了对方，再检查是否为互关
    if (isFollowing.value) {
      try {
        // 检查互关列表，看目标用户是否在我的互关列表中
        const mutualResponse = await FollowApi.getMutualFollows({ page: 1, limit: 100 })
        console.log('互关列表API响应:', mutualResponse)
        const mutualUsers = mutualResponse.data?.users || mutualResponse.users || []
        console.log('互关用户列表:', mutualUsers)

        // 检查目标用户是否在互关列表中
        isMutualFollow.value = mutualUsers.some((user: any) => user.id === userId.value)
        console.log('是否互关:', isMutualFollow.value)
      } catch (error) {
        console.error('检查互关状态失败:', error)
        isMutualFollow.value = false
      }
    } else {
      isMutualFollow.value = false
    }

  } catch (error) {
    console.error('检查关注状态失败:', error)
    isFollowing.value = false
    isMutualFollow.value = false
  } finally {
    followStatusLoading.value = false
  }
}

// username 已在上面定义

const userId = computed(() => {
  return profile.value?.id
})

const isOwnProfile = computed(() => {
  const currentUsername = authStore.user?.username
  const profileUsername = username.value
  return currentUsername === profileUsername
})

const tabs = computed(() => {
  return [
    { key: 'articles', label: '文章' },
    { key: 'activity', label: '动态' }
  ]
})

// 获取关注按钮文字
const getFollowButtonText = () => {
  if (followStatusLoading.value || isFollowing.value === null) {
    return '加载中...'
  }
  if (!isFollowing.value) {
    return '关注'
  } else if (isMutualFollow.value) {
    return '互相关注'
  } else {
    return '已关注'
  }
}

// 获取关注按钮样式类
const getFollowButtonClass = () => {
  if (followStatusLoading.value || isFollowing.value === null) {
    return 'loading'
  }
  if (!isFollowing.value) {
    return 'not-following'
  } else if (isMutualFollow.value) {
    return 'mutual-following'
  } else {
    return 'following'
  }
}

const loadProfile = async () => {
  if (!username.value) return

  try {
    await userDataSync.loadUserProfile()
    await userDataSync.loadUserStats()
    await checkFollowStatus()
  } catch (error: any) {
    showToast(error.response?.data?.message || '加载用户信息失败', 'error')
  }
}

const loadArticles = async () => {
  if (!username.value) return

  try {
    await userDataSync.loadUserArticles({ page: 1, size: 20 })
  } catch (error) {
    console.error('加载文章列表失败:', error)
  }
}

const loadLikedArticles = async () => {
  if (!isOwnProfile.value) return

  try {
    // TODO: 实现获取用户点赞文章的API
    likedArticles.value = []
  } catch (error) {
    console.error('加载点赞文章失败:', error)
  }
}

const loadActivities = async () => {
  if (!username.value) return

  try {
    // TODO: 实现获取用户动态的API
    activities.value = []
  } catch (error) {
    console.error('加载用户动态失败:', error)
  }
}

const formatDate = (dateString?: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long'
  })
}

watch(() => route.params.username, () => {
  loadProfile()
  loadArticles()
  if (activeTab.value === 'liked') {
    loadLikedArticles()
  } else if (activeTab.value === 'activity') {
    loadActivities()
  }
})

watch(activeTab, (newTab) => {
  if (newTab === 'liked' && likedArticles.value.length === 0) {
    loadLikedArticles()
  } else if (newTab === 'activity' && activities.value.length === 0) {
    loadActivities()
  }
})



// 关注/取消关注用户
const handleFollowAction = async () => {
  if (!userId.value || isOwnProfile.value) return

  try {
    followActionLoading.value = true

    if (isFollowing.value) {
      await FollowApi.unfollowUser(userId.value)
      isFollowing.value = false
      isMutualFollow.value = false
      showToast('已取消关注', 'success')
    } else {
      await FollowApi.followUser(userId.value)
      showToast('关注成功', 'success')
      // 重新检查关注状态，因为可能建立了互关
      await checkFollowStatus()
    }
  } catch (error: any) {
    showToast(error.response?.data?.message || '操作失败', 'error')
  } finally {
    followActionLoading.value = false
  }
}

// 开始聊天
const startChat = async () => {
  if (!userId.value || isOwnProfile.value) return

  try {
    chatLoading.value = true
    // 跳转到私信页面，传递用户ID参数
    await router.push(`/messages?user=${userId.value}`)
  } catch (error) {
    showToast('无法发送私信', 'error')
  } finally {
    chatLoading.value = false
  }
}


// 监听用户ID变化，当用户数据加载完成后立即检查关注状态
watch(userId, async (newUserId) => {
  if (newUserId && !isOwnProfile.value && authStore.isAuthenticated) {
    console.log('=== userId变化，立即检查关注状态 ===')
    console.log('新的userId:', newUserId)
    await checkFollowStatus()
  }
}, { immediate: true })

onMounted(async () => {
  console.log('=== UserProfilePage onMounted 开始 ===')
  console.log('当前用户名参数:', username.value)
  console.log('当前认证状态:', authStore.isAuthenticated)
  console.log('当前用户:', authStore.user)

  // 初始化用户数据
  await userDataSync.initUserData()
})
</script>

<style scoped>
.user-profile-page {
  min-height: 100vh;
  background: #f8fafc;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.profile-header {
  background: white;
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.profile-cover {
  display: flex;
  gap: 24px;
  align-items: flex-start;
  margin-bottom: 32px;
}

.profile-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  flex-shrink: 0;
}

.default-avatar {
  background: linear-gradient(135deg, #3B82F6, #8B5CF6);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 600;
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 4px;
}

.profile-username {
  font-size: 18px;
  color: #64748b;
  margin-bottom: 12px;
}

.profile-bio {
  font-size: 16px;
  color: #475569;
  line-height: 1.6;
  margin-bottom: 16px;
  max-width: 500px;
}

.profile-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #64748b;
}

.meta-icon {
  flex-shrink: 0;
}

.meta-item a {
  color: #3B82F6;
  text-decoration: none;
}

.meta-item a:hover {
  text-decoration: underline;
}

.profile-actions {
  display: flex;
  align-items: flex-start;
}

.profile-manage-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.edit-profile-btn {
  background: #3B82F6;
  color: white;
  text-decoration: none;
  padding: 10px 20px;
  border-radius: 20px;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.edit-profile-btn:hover {
  background: #2563EB;
  transform: translateY(-1px);
}

.dropdown {
  position: relative;
}

.manage-menu-btn {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #64748b;
  padding: 10px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.manage-menu-btn:hover {
  background: #f1f5f9;
  color: #475569;
  transform: translateY(-1px);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  min-width: 160px;
  z-index: 50;
  overflow: hidden;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: #374151;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s;
  cursor: pointer;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
}

.dropdown-item:hover {
  background: #f8fafc;
  color: #3B82F6;
}

.dropdown-item svg {
  flex-shrink: 0;
}

/* 关注按钮样式 */
.follow-btn {
  background: #3B82F6;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 20px;
  font-weight: 500;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 100px;
}

.follow-btn:hover {
  background: #2563EB;
  transform: translateY(-1px);
}

.follow-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.follow-btn.following {
  background: #059669;
}

.follow-btn.following:hover {
  background: #047857;
}

/* 互相关注按钮样式 */
.follow-btn.mutual-following {
  background: #7C3AED;
}

.follow-btn.mutual-following:hover {
  background: #6D28D9;
}

/* 加载状态按钮样式 */
.follow-btn.loading {
  background: #F1F5F9;
  color: #94A3B8;
  cursor: not-allowed;
  opacity: 0.8;
}

.follow-btn.loading:hover {
  background: #F1F5F9;
  transform: none;
}

/* 私信按钮样式 */
.chat-btn {
  background: #F1F5F9;
  color: #475569;
  border: 1px solid #E2E8F0;
  padding: 10px 20px;
  border-radius: 20px;
  font-weight: 500;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 100px;
}

.chat-btn:hover {
  background: #E2E8F0;
  color: #334155;
  transform: translateY(-1px);
}

.chat-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.stats-container {
  border-top: 1px solid #e2e8f0;
  padding-top: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 32px;
}

.stat-item {
  text-align: center;
  cursor: pointer;
  transition: opacity 0.2s;
}

.stat-item:hover {
  opacity: 0.8;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #64748b;
}

.content-tabs {
  display: flex;
  gap: 32px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e2e8f0;
}

.tab-button {
  background: none;
  border: none;
  font-size: 16px;
  font-weight: 500;
  color: #64748b;
  padding: 16px 0;
  cursor: pointer;
  position: relative;
  transition: color 0.2s;
}

.tab-button:hover {
  color: #3B82F6;
}

.tab-button.active {
  color: #3B82F6;
}

.tab-button.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #3B82F6;
}

.content-area {
  background: white;
  border-radius: 16px;
  padding: 32px;
  min-height: 400px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 文本截断样式 */
.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

@media (max-width: 768px) {
  .profile-cover {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .profile-avatar {
    width: 100px;
    height: 100px;
  }
  
  .profile-name {
    font-size: 24px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
  
  .content-tabs {
    overflow-x: auto;
    white-space: nowrap;
  }
  
  .articles-grid {
    grid-template-columns: 1fr;
  }
}
</style>