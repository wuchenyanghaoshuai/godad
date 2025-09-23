<template>
  <AppLayout
    :header-config="{
      showSearch: false,
      showCreateButton: false,
      showNotifications: true,
      showUserPoints: true,
      showNavigation: false,
      showUserMenu: true
    }"
    :show-footer="false"
    background-class="bg-gray-50"
  >
    <PageContainer background="gray" padding="lg">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 侧边栏 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow p-6">
            <!-- 用户头像和基本信息 -->
            <div class="text-center mb-6">
              <div class="w-24 h-24 mx-auto mb-4 relative group cursor-pointer" @click="showAvatarModal = true">
                <UserAvatar
                  :avatar="user?.avatar || ''"
                  :name="user?.nickname || user?.username || 'U'"
                  :size="96"
                />
                <!-- 悬停提示 -->
                <div class="absolute inset-0 rounded-full bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                  <CameraIcon class="h-6 w-6 text-white" />
                </div>
              </div>
              <p class="text-xs text-gray-500 mb-2">点击头像更换</p>
              <h2 class="text-xl font-semibold text-gray-900">{{ user?.username }}</h2>
              <p class="text-gray-600">{{ user?.email }}</p>
              <span class="inline-block px-2 py-1 text-xs font-medium rounded-full mt-2"
                    :class="roleClasses">
                {{ roleText }}
              </span>

              <!-- 用户等级和积分 -->
              <div class="mt-4 pt-4 border-t border-gray-200">
                <UserPointsDisplay
                  mode="detailed"
                  :auto-refresh="true"
                  class="mb-4"
                />
              </div>

              <!-- 关注统计信息 -->
              <div class="flex justify-center space-x-6 mt-4 pt-4 border-t border-gray-200">
                <div class="text-center cursor-pointer hover:text-pink-600 transition-colors" @click="activeTab = 'following'">
                  <div class="text-lg font-semibold">{{ followingCount }}</div>
                  <div class="text-xs text-gray-500">关注</div>
                </div>
                <div class="text-center cursor-pointer hover:text-pink-600 transition-colors" @click="activeTab = 'followers'">
                  <div class="text-lg font-semibold">{{ followersCount }}</div>
                  <div class="text-xs text-gray-500">粉丝</div>
                </div>
                <div class="text-center">
                  <div class="text-lg font-semibold">{{ articlesCount }}</div>
                  <div class="text-xs text-gray-500">文章</div>
                </div>
              </div>
            </div>

            <!-- 导航菜单 -->
            <nav class="space-y-2">
              <button
                v-for="item in menuItems"
                :key="item.key"
                @click="activeTab = item.key"
                :class="[
                  'w-full flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                  activeTab === item.key
                    ? 'bg-pink-50 text-pink-700 border-r-2 border-pink-500'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                ]"
              >
                <component :is="item.icon" class="h-5 w-5 mr-3" />
                {{ item.label }}
              </button>
            </nav>
          </div>
        </div>

        <!-- 主内容区 -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-lg shadow">
            <!-- 个人信息 -->
            <div v-if="activeTab === 'profile'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">个人信息</h3>
              <form @submit.prevent="updateProfile" class="space-y-6">

                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      用户名
                    </label>
                    <input
                      v-model="profileForm.username"
                      type="text"
                      disabled
                      class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-500"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      邮箱
                    </label>
                    <input
                      v-model="profileForm.email"
                      type="email"
                      disabled
                      class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-500"
                    />
                  </div>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    昵称
                  </label>
                  <input
                    v-model="profileForm.nickname"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                    placeholder="请输入昵称"
                  />
                </div>

                <!-- 手机号和性别 -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      手机号
                    </label>
                    <input
                      v-model="profileForm.phone"
                      type="tel"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      placeholder="请输入手机号"
                      pattern="[0-9]{11}"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      性别
                    </label>
                    <select
                      v-model="profileForm.gender"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                    >
                      <option value="">请选择性别</option>
                      <option value="1">男</option>
                      <option value="2">女</option>
                      <option value="0">保密</option>
                    </select>
                  </div>
                </div>

                <!-- 生日 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    生日
                  </label>
                  <input
                    v-model="profileForm.birthday"
                    type="date"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    个人简介
                  </label>
                  <textarea
                    v-model="profileForm.bio"
                    rows="4"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                    placeholder="介绍一下自己..."
                  ></textarea>
                </div>
                <div class="flex justify-end">
                  <button
                    type="submit"
                    :disabled="isUpdating"
                    class="bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700 disabled:opacity-50"
                  >
                    {{ isUpdating ? '保存中...' : '保存更改' }}
                  </button>
                </div>
              </form>
            </div>

            <!-- 我的文章 -->
            <div v-else-if="activeTab === 'articles'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的文章</h3>
                <router-link
                  to="/articles/create"
                  class="bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700"
                >
                  写文章
                </router-link>
              </div>

              <!-- 加载状态 -->
              <div v-if="isLoadingArticles" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="mt-4 text-gray-500">加载中...</p>
              </div>

              <!-- 错误状态 -->
              <div v-else-if="articlesError" class="text-center py-12 text-red-500">
                <p>{{ articlesError }}</p>
                <button
                  @click="loadMyArticles"
                  class="mt-4 bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700"
                >
                  重试
                </button>
              </div>

              <!-- 文章列表 -->
              <div v-else-if="myArticles.length > 0" class="space-y-4">
                <div
                  v-for="article in myArticles"
                  :key="article.id"
                  class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
                  @click="router.push(`/articles/${article.id}`)"
                >
                  <div class="flex justify-between items-start">
                    <div class="flex-1">
                      <h4 class="text-lg font-medium text-gray-900 hover:text-pink-600">
                        {{ article.title }}
                      </h4>
                      <p v-if="article.summary" class="text-gray-600 mt-2 text-sm line-clamp-2">
                        {{ article.summary }}
                      </p>
                      <div class="flex items-center mt-3 space-x-4 text-sm text-gray-500">
                        <span>{{ formatDate(article.created_at) }}</span>
                        <span>阅读 {{ article.view_count || 0 }}</span>
                        <span>点赞 {{ article.like_count || 0 }}</span>
                        <span
                          class="px-2 py-1 rounded-full text-xs"
                          :class="getStatusClass(article.status)"
                        >
                          {{ getStatusText(article.status) }}
                        </span>
                      </div>
                    </div>
                    <div class="ml-4 flex space-x-2">
                      <router-link
                        :to="`/articles/${article.id}/edit`"
                        class="text-blue-600 hover:text-blue-700 text-sm"
                        @click.stop
                      >
                        编辑
                      </router-link>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <FileTextIcon class="h-12 w-12 mx-auto mb-4" />
                <p>您还没有发布任何文章</p>
                <router-link
                  to="/articles/create"
                  class="text-pink-600 hover:text-pink-500 mt-2 inline-block"
                >
                  立即创建第一篇文章
                </router-link>
              </div>
            </div>

            <!-- 其他标签页内容会在这里添加... -->
            <div v-else class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">{{ activeTab }}</h3>
              <p class="text-gray-500">功能开发中...</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 头像上传弹窗 -->
      <AvatarModal
        :is-visible="showAvatarModal"
        @close="closeAvatarModal"
        @success="handleAvatarUpload"
        @error="handleUploadError"
      />
    </PageContainer>
  </AppLayout>
</template>

<script setup lang="ts">
/* eslint-disable @typescript-eslint/no-unused-vars */
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  UserIcon,
  FileTextIcon,
  StarIcon,
  SettingsIcon,
  CameraIcon,
  UsersIcon,
  HeartIcon,
  BellIcon
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { AppLayout, PageContainer } from '@/components/layout'
import AvatarModal from '../components/AvatarModal.vue'
import UserPointsDisplay from '../components/UserPointsDisplay.vue'
import { useToast } from '../composables/useToast'
import { useUserDataSync } from '../composables/useUserDataSync'
import type { ImageUploadResponse } from '../api/types'
import UserAvatar from '@/components/UserAvatar.vue'

// 路由
const router = useRouter()

// 认证store
const authStore = useAuthStore()
const { toast } = useToast()

// 用户数据同步
const userDataSync = useUserDataSync()

// 响应式数据
const showAvatarModal = ref(false)
const activeTab = ref('profile')
const isUpdating = ref(false)

// 使用同步的数据 - 直接使用userDataSync返回的computed属性
const myArticles = userDataSync.articles
const isLoadingArticles = userDataSync.isLoading
const followingCount = computed(() => userDataSync.stats.value.following_count)
const followersCount = computed(() => userDataSync.stats.value.followers_count)
const articlesCount = computed(() => userDataSync.stats.value.article_count)
const articlesError = ref('')

// 用户信息
const user = computed(() => userDataSync.profile.value || authStore.user)

// 角色显示
const roleText = computed(() => {
  switch (user.value?.role) {
    case 'admin':
      return '管理员'
    case 'content_manager':
      return '内容管理员'
    default:
      return '普通用户'
  }
})

const roleClasses = computed(() => {
  switch (user.value?.role) {
    case 'admin':
      return 'bg-red-100 text-red-800'
    case 'content_manager':
      return 'bg-blue-100 text-blue-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})

// 菜单项
const menuItems = [
  { key: 'profile', label: '个人信息', icon: UserIcon },
  { key: 'articles', label: '我的文章', icon: FileTextIcon },
  { key: 'favorites', label: '我的收藏', icon: StarIcon },
  { key: 'settings', label: '设置', icon: SettingsIcon }
]

// 个人信息表单
const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  phone: '',
  gender: '',
  birthday: '',
  bio: '',
  avatar: ''
})

// 初始化用户信息
const initUserInfo = () => {
  if (user.value) {
    profileForm.username = user.value.username
    profileForm.email = user.value.email
    profileForm.nickname = user.value.nickname || ''
    profileForm.phone = user.value.phone || ''
    profileForm.gender = user.value.gender?.toString() || ''
    profileForm.birthday = user.value.birthday ? user.value.birthday.split('T')[0] : ''
    profileForm.bio = user.value.bio || ''
    profileForm.avatar = user.value.avatar || ''
  }
}

// 头像上传成功处理
const handleAvatarUpload = async (response: ImageUploadResponse) => {
  const avatarUrl = response.url

  if (avatarUrl) {
    try {
      // 更新表单数据
      profileForm.avatar = avatarUrl

      // 使用数据同步函数更新头像
      await userDataSync.updateUserProfile({
        avatar: avatarUrl
      })

      // 确保DOM更新
      await nextTick()

      // 关闭头像上传弹窗
      closeAvatarModal()
    } catch (error) {
      console.error('头像更新失败:', error)
      toast.error('头像更新失败，请重试')
    }
  }
}

// 关闭头像弹窗
const closeAvatarModal = () => {
  showAvatarModal.value = false
}

// 上传错误处理
const handleUploadError = (error: string) => {
  toast.error(`上传失败: ${error}`)
}

// 更新个人信息
const updateProfile = async () => {
  try {
    isUpdating.value = true

    // 调用更新用户信息API
    const updateData: any = {
      nickname: profileForm.nickname,
      phone: profileForm.phone,
      gender: profileForm.gender ? parseInt(profileForm.gender) : undefined,
      birthday: profileForm.birthday ? `${profileForm.birthday}T00:00:00Z` : undefined,
      bio: profileForm.bio,
      avatar: profileForm.avatar
    }

    // 移除空值
    Object.keys(updateData).forEach(key => {
      if (updateData[key] === '' || updateData[key] === undefined) {
        delete updateData[key]
      }
    })

    // 使用数据同步函数更新
    await userDataSync.updateUserProfile(updateData)

    // 重新初始化表单
    initUserInfo()
  } catch (error) {
    console.error('更新失败:', error)
    toast.error('更新失败，请重试')
  } finally {
    isUpdating.value = false
  }
}

// 加载我的文章（使用数据同步）
const loadMyArticles = async () => {
  try {
    articlesError.value = ''
    await userDataSync.loadUserArticles({ page: 1, size: 20 })
  } catch (error: any) {
    articlesError.value = error.message || '加载文章失败'
    console.error('加载我的文章失败:', error)
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 获取状态文本
const getStatusText = (status: number) => {
  switch (status) {
    case 0: return '草稿'
    case 1: return '已发布'
    case 2: return '已下架'
    default: return '未知'
  }
}

// 获取状态样式
const getStatusClass = (status: number) => {
  switch (status) {
    case 0: return 'bg-gray-100 text-gray-600'
    case 1: return 'bg-green-100 text-green-600'
    case 2: return 'bg-red-100 text-red-600'
    default: return 'bg-gray-100 text-gray-600'
  }
}

// 监听activeTab变化，根据标签加载不同数据
watch(activeTab, (newTab) => {
  if (newTab === 'articles') {
    loadMyArticles()
  }
})

// 组件挂载时初始化
onMounted(() => {
  // 检查登录状态
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }

  // 处理URL参数中的tab
  const urlParams = new URLSearchParams(window.location.search)
  const tab = urlParams.get('tab')
  if (tab && ['profile', 'articles', 'favorites', 'settings'].includes(tab)) {
    activeTab.value = tab
  }

  initUserInfo()
  // 初始化用户数据同步
  userDataSync.initUserData()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
