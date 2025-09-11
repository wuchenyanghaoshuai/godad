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
            v-if="profile?.user?.avatar" 
            :src="profile.user.avatar" 
            alt="用户头像"
            class="profile-avatar"
          />
          <div v-else class="profile-avatar default-avatar">
            {{ profile?.user?.username?.charAt(0)?.toUpperCase() || 'U' }}
          </div>
          
          <div class="profile-info">
            <h1 class="profile-name">
              {{ profile?.user?.nickname || profile?.user?.username || '加载中...' }}
            </h1>
            <p class="profile-username">@{{ profile?.user?.username || '' }}</p>
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
                加入于 {{ formatDate(profile?.joined_at) }}
              </div>
            </div>
          </div>
          
          <div class="profile-actions">
            <div class="profile-manage-actions">
              <button
                @click="activeTab = 'settings'"
                class="edit-profile-btn primary"
              >
                编辑资料
              </button>
              <div class="dropdown">
                <button 
                  @click="showManageMenu = !showManageMenu"
                  class="manage-menu-btn"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="1"/>
                    <circle cx="12" cy="5" r="1"/>
                    <circle cx="12" cy="19" r="1"/>
                  </svg>
                </button>
                
                <div v-if="showManageMenu" class="dropdown-menu">
                  <router-link to="/articles/create" class="dropdown-item">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                      <polyline points="14,2 14,8 20,8"/>
                      <line x1="16" y1="13" x2="8" y2="13"/>
                      <line x1="16" y1="17" x2="8" y2="17"/>
                      <polyline points="10,9 9,9 8,9"/>
                    </svg>
                    发布文章
                  </router-link>
                  <router-link to="/my/articles" class="dropdown-item">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                      <polyline points="14,2 14,8 20,8"/>
                    </svg>
                    管理文章
                  </router-link>
                  <div class="dropdown-item" @click="showAvatarUpload = true">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                      <circle cx="12" cy="7" r="4"/>
                    </svg>
                    更换头像
                  </div>
                </div>
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
              <div class="stat-number">{{ stats?.follower_count || 0 }}</div>
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
          <div v-if="articles.length === 0" class="empty-state">
            <Empty message="还没有发布任何文章" />
          </div>
          <div v-else class="articles-grid">
            <!-- TODO: 集成ArticleCard组件 -->
            <div 
              v-for="article in articles"
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

        <!-- 管理标签 - 只对自己可见 -->
        <div v-if="activeTab === 'manage' && isOwnProfile" class="manage-content">
          <div class="mb-6 flex justify-between items-center">
            <h3 class="text-xl font-semibold">我的文章管理</h3>
            <router-link
              to="/articles/create"
              class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
            >
              发布文章
            </router-link>
          </div>
          
          <div v-if="articles.length === 0" class="empty-state">
            <Empty message="还没有发布任何文章">
              <router-link
                to="/articles/create"
                class="text-blue-600 hover:text-blue-500 mt-2 inline-block"
              >
                立即创建第一篇文章
              </router-link>
            </Empty>
          </div>
          <div v-else class="articles-manage-list space-y-4">
            <div 
              v-for="article in articles"
              :key="article.id"
              class="bg-white p-6 rounded-lg shadow-sm border hover:shadow-md transition-shadow"
            >
              <div class="flex justify-between items-start">
                <div class="flex-1">
                  <h4 class="text-lg font-medium text-gray-900 hover:text-blue-600 cursor-pointer"
                      @click="$router.push(`/articles/${article.id}`)">
                    {{ article.title }}
                  </h4>
                  <p v-if="article.summary" class="text-gray-600 mt-2 text-sm">
                    {{ article.summary }}
                  </p>
                  <div class="flex items-center mt-3 space-x-4 text-sm text-gray-500">
                    <span>发布于 {{ new Date(article.created_at).toLocaleDateString('zh-CN') }}</span>
                    <span>阅读 {{ article.view_count || 0 }}</span>
                    <span>点赞 {{ article.like_count || 0 }}</span>
                    <span 
                      class="px-2 py-1 rounded-full text-xs"
                      :class="article.status === 1 ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                    >
                      {{ article.status === 1 ? '已发布' : '草稿' }}
                    </span>
                  </div>
                </div>
                <div class="ml-4 flex space-x-3">
                  <router-link
                    :to="`/articles/${article.id}/edit`"
                    class="text-blue-600 hover:text-blue-700 text-sm px-3 py-1 border border-blue-200 rounded hover:bg-blue-50"
                  >
                    编辑
                  </router-link>
                  <button
                    @click="deleteArticle(article.id)"
                    class="text-red-600 hover:text-red-700 text-sm px-3 py-1 border border-red-200 rounded hover:bg-red-50"
                  >
                    删除
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 设置标签 - 只对自己可见 -->
        <div v-if="activeTab === 'settings' && isOwnProfile" class="settings-content">
          <h3 class="text-xl font-semibold mb-6">账户设置</h3>
          
          <div class="space-y-8">
            <!-- 个人资料编辑 -->
            <div class="bg-gray-50 p-6 rounded-lg">
              <h4 class="text-lg font-medium text-gray-900 mb-4">个人资料</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    昵称
                  </label>
                  <input
                    v-model="profileEditForm.nickname"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    邮箱
                  </label>
                  <input
                    v-model="profileEditForm.email"
                    type="email"
                    disabled
                    class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-500"
                  />
                </div>
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    个人简介
                  </label>
                  <textarea
                    v-model="profileEditForm.bio"
                    rows="4"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="介绍一下自己..."
                  ></textarea>
                </div>
              </div>
              <div class="mt-6">
                <button
                  @click="updateProfileInfo"
                  :disabled="updating"
                  class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ updating ? '保存中...' : '保存资料' }}
                </button>
              </div>
            </div>

            <!-- 修改密码 -->
            <div class="bg-gray-50 p-6 rounded-lg">
              <h4 class="text-lg font-medium text-gray-900 mb-4">修改密码</h4>
              <div class="space-y-4 max-w-md">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    当前密码
                  </label>
                  <input
                    v-model="passwordForm.currentPassword"
                    type="password"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    新密码
                  </label>
                  <input
                    v-model="passwordForm.newPassword"
                    type="password"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    确认新密码
                  </label>
                  <input
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
              <div class="mt-6">
                <button
                  @click="changePassword"
                  :disabled="changingPassword"
                  class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50"
                >
                  {{ changingPassword ? '修改中...' : '修改密码' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import userProfileApi, { type UserProfile, type UserStats, type UserActivity } from '@/api/userProfile'
import { ArticleApi } from '@/api/article'
import likeApi from '@/api/likes'
import { useToast } from '@/composables/useToast'
import Navbar from '@/components/Navbar.vue'
import Empty from '@/components/Empty.vue'
// import ArticleCard from '@/components/ArticleCard.vue'
// import ActivityItem from '@/components/ActivityItem.vue'

const route = useRoute()
const authStore = useAuthStore()
const { toast } = useToast()

const profile = ref<UserProfile | null>(null)
const stats = ref<UserStats | null>(null)
const articles = ref<any[]>([])
const likedArticles = ref<any[]>([])
const activities = ref<UserActivity[]>([])
const loading = ref(false)
const activeTab = ref('articles')
const showFollowingModal = ref(false)
const showFollowersModal = ref(false)
const showManageMenu = ref(false)
const showAvatarUpload = ref(false)

// 用于管理和设置功能的响应式数据
const profileEditForm = ref({
  nickname: '',
  email: '',
  bio: ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const updating = ref(false)
const changingPassword = ref(false)

const userId = computed(() => {
  const id = route.params.id as string
  return id === 'me' ? authStore.user?.id : parseInt(id)
})

const isOwnProfile = computed(() => {
  return userId.value === authStore.user?.id
})

const tabs = computed(() => {
  const baseTabs = [
    { key: 'articles', label: '文章' },
    { key: 'activity', label: '动态' }
  ]
  
  if (isOwnProfile.value) {
    baseTabs.splice(1, 0, { key: 'liked', label: '点赞' })
    baseTabs.push({ key: 'manage', label: '管理' })
    baseTabs.push({ key: 'settings', label: '设置' })
  }
  
  return baseTabs
})

const loadProfile = async () => {
  if (!userId.value) return

  try {
    loading.value = true
    const profileRes = await userProfileApi.getProfile(isOwnProfile.value ? undefined : userId.value)
    
    profile.value = profileRes.data
    
    // 如果是自己的资料，从API获取详细统计；如果是他人，从profile中获取
    if (isOwnProfile.value) {
      const statsRes = await userProfileApi.getUserStats()
      stats.value = statsRes.data
    } else {
      // 从profile中提取统计信息
      stats.value = {
        article_count: profileRes.data.article_count,
        like_count: 0, // 他人资料暂时不显示获赞数
        comment_count: 0, // 他人资料暂时不显示评论数
        follower_count: profileRes.data.follower_count,
        following_count: profileRes.data.following_count,
        view_count: 0 // 他人资料暂时不显示浏览数
      }
    }
  } catch (error: any) {
    toast.error(error.response?.data?.message || '加载用户信息失败')
  } finally {
    loading.value = false
  }
}

const loadArticles = async () => {
  if (!userId.value) return

  try {
    // 使用文章列表API，通过作者ID参数过滤
    const response = await ArticleApi.getArticleList({ 
      author_id: userId.value,
      status: 1, // 只获取已发布的文章
      page: 1,
      size: 20
    })
    articles.value = response.data || []
  } catch (error) {
    console.error('加载文章列表失败:', error)
  }
}

const loadLikedArticles = async () => {
  if (!isOwnProfile.value) return

  try {
    const response = await likeApi.getUserLikes('article')
    likedArticles.value = response.data.data || []
  } catch (error) {
    console.error('加载点赞文章失败:', error)
  }
}

const loadActivities = async () => {
  if (!userId.value) return

  try {
    const response = await userProfileApi.getUserActivity(
      isOwnProfile.value ? undefined : userId.value
    )
    activities.value = response.data.data || []
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

watch(() => route.params.id, () => {
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

// 初始化个人资料编辑表单
const initProfileEditForm = () => {
  if (profile.value?.user) {
    profileEditForm.value = {
      nickname: profile.value.user.nickname || profile.value.user.username || '',
      email: profile.value.user.email || '',
      bio: profile.value.bio || ''
    }
  }
}

// 更新个人资料信息
const updateProfileInfo = async () => {
  if (!profileEditForm.value.nickname.trim()) {
    toast.error('昵称不能为空')
    return
  }

  try {
    updating.value = true
    
    // 更新用户基本信息（昵称）
    if (profileEditForm.value.nickname !== profile.value?.user?.nickname) {
      // TODO: 调用更新用户信息的API
      // await userApi.updateUser({ nickname: profileEditForm.value.nickname })
    }
    
    // 更新用户资料（个人简介等）
    if (profileEditForm.value.bio !== profile.value?.bio) {
      await userProfileApi.updateProfile({
        bio: profileEditForm.value.bio
      })
    }
    
    toast.success('个人资料更新成功')
    
    // 重新加载资料
    await loadProfile()
    initProfileEditForm()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '更新失败')
  } finally {
    updating.value = false
  }
}

// 修改密码
const changePassword = async () => {
  if (!passwordForm.value.currentPassword || !passwordForm.value.newPassword || !passwordForm.value.confirmPassword) {
    toast.error('所有密码字段都必须填写')
    return
  }

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    toast.error('新密码与确认密码不匹配')
    return
  }

  if (passwordForm.value.newPassword.length < 6) {
    toast.error('新密码长度至少6位')
    return
  }

  try {
    changingPassword.value = true
    
    // TODO: 调用修改密码的API
    // await userApi.changePassword({
    //   current_password: passwordForm.value.currentPassword,
    //   new_password: passwordForm.value.newPassword
    // })
    
    toast.success('密码修改成功')
    
    // 清空表单
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error: any) {
    toast.error(error.response?.data?.message || '密码修改失败')
  } finally {
    changingPassword.value = false
  }
}

// 删除文章
const deleteArticle = async (articleId: number) => {
  if (!confirm('确认删除这篇文章吗？此操作不可撤销。')) {
    return
  }

  try {
    await ArticleApi.deleteArticle(articleId)
    toast.success('文章删除成功')
    
    // 从列表中移除该文章
    articles.value = articles.value.filter(article => article.id !== articleId)
    
    // 更新统计数据
    if (stats.value) {
      stats.value.article_count = Math.max(0, stats.value.article_count - 1)
    }
  } catch (error: any) {
    toast.error(error.response?.data?.message || '删除失败')
  }
}

// 点击外部关闭下拉菜单
const handleClickOutside = (event: Event) => {
  const target = event.target as Element
  if (!target.closest('.dropdown')) {
    showManageMenu.value = false
  }
}

onMounted(() => {
  loadProfile().then(() => {
    // 资料加载完成后初始化编辑表单
    if (isOwnProfile.value) {
      initProfileEditForm()
    }
  })
  loadArticles()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
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