<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 头部导航 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <router-link to="/" class="text-xl font-bold text-gray-900">
            GoDad
          </router-link>
          <div class="flex items-center space-x-4">
            <span class="text-gray-700">{{ user?.username }}</span>
            <button
              @click="handleLogout"
              class="text-gray-500 hover:text-gray-700"
            >
              <LogOutIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 侧边栏 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow p-6">
            <!-- 用户头像和基本信息 -->
            <div class="text-center mb-6">
              <div class="w-24 h-24 mx-auto mb-4 relative group cursor-pointer" @click="showAvatarModal = true">
                <img
                  v-if="user?.avatar"
                  :src="user.avatar"
                  :alt="user.username"
                  class="w-24 h-24 rounded-full object-cover transition-opacity group-hover:opacity-75"
                />
                <div
                  v-else
                  class="w-24 h-24 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center transition-opacity group-hover:opacity-75"
                >
                  <span class="text-2xl font-bold text-white">
                    {{ user?.username?.charAt(0).toUpperCase() }}
                  </span>
                </div>
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

            <!-- 我的关注 -->
            <div v-else-if="activeTab === 'following'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的关注 ({{ followingCount }})</h3>
              </div>
              
              <div v-if="isLoadingFollows" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="text-gray-500 mt-2">加载中...</p>
              </div>
              
              <!-- 关注列表 -->
              <div v-else-if="followingList.length > 0" class="space-y-4">
                <div v-for="user in followingList" :key="user.id" class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-r from-pink-400 to-orange-400 flex items-center justify-center">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.username" class="w-full h-full object-cover">
                      <span v-else class="text-white font-semibold">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900">{{ user.nickname || user.username }}</h4>
                      <p class="text-sm text-gray-500">@{{ user.username }}</p>
                      <p v-if="user.bio" class="text-sm text-gray-600 mt-1">{{ user.bio }}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span v-if="user.is_mutual_follow" class="text-xs bg-pink-100 text-pink-600 px-2 py-1 rounded-full">互关</span>
                    <button 
                      @click="unfollowUser(user.id)"
                      class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm"
                    >
                      取消关注
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <HeartIcon class="h-12 w-12 mx-auto mb-4" />
                <p>您还没有关注任何人</p>
                <p class="text-sm mt-2">去发现更多有趣的用户吧</p>
              </div>
            </div>

            <!-- 我的粉丝 -->
            <div v-else-if="activeTab === 'followers'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的粉丝 ({{ followersCount }})</h3>
              </div>
              
              <div v-if="isLoadingFollows" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="text-gray-500 mt-2">加载中...</p>
              </div>
              
              <!-- 粉丝列表 -->
              <div v-else-if="followersList.length > 0" class="space-y-4">
                <div v-for="user in followersList" :key="user.id" class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-r from-pink-400 to-orange-400 flex items-center justify-center">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.username" class="w-full h-full object-cover">
                      <span v-else class="text-white font-semibold">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900">{{ user.nickname || user.username }}</h4>
                      <p class="text-sm text-gray-500">@{{ user.username }}</p>
                      <p v-if="user.bio" class="text-sm text-gray-600 mt-1">{{ user.bio }}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span v-if="user.is_mutual_follow" class="text-xs bg-pink-100 text-pink-600 px-2 py-1 rounded-full">互关</span>
                    <button 
                      v-if="!user.is_following" 
                      @click="followUser(user.id)"
                      class="px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 transition-colors text-sm"
                    >
                      回关
                    </button>
                    <button 
                      v-else 
                      @click="unfollowUser(user.id)"
                      class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm"
                    >
                      取消关注
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <UsersIcon class="h-12 w-12 mx-auto mb-4" />
                <p>还没有人关注您</p>
                <p class="text-sm mt-2">发布有趣的内容来吸引更多关注吧</p>
              </div>
            </div>

            <!-- 互相关注 -->
            <div v-else-if="activeTab === 'mutual'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">互相关注 ({{ mutualCount }})</h3>
              </div>
              
              <div v-if="isLoadingFollows" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="text-gray-500 mt-2">加载中...</p>
              </div>
              
              <!-- 互关列表 -->
              <div v-else-if="mutualFollowsList.length > 0" class="space-y-4">
                <div v-for="user in mutualFollowsList" :key="user.id" class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-r from-pink-400 to-orange-400 flex items-center justify-center">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.username" class="w-full h-full object-cover">
                      <span v-else class="text-white font-semibold">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900">{{ user.nickname || user.username }}</h4>
                      <p class="text-sm text-gray-500">@{{ user.username }}</p>
                      <p v-if="user.bio" class="text-sm text-gray-600 mt-1">{{ user.bio }}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span class="text-xs bg-pink-100 text-pink-600 px-3 py-2 rounded-full font-medium">💕 互相关注</span>
                    <button 
                      @click="unfollowUser(user.id)"
                      class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm"
                    >
                      取消关注
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <HeartIcon class="h-12 w-12 mx-auto mb-4 text-pink-300" />
                <p>还没有互相关注的用户</p>
                <p class="text-sm mt-2">关注其他用户，等待他们回关吧</p>
              </div>
            </div>

            <!-- 设置 -->
            <div v-else-if="activeTab === 'settings'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">账户设置</h3>
              <div class="space-y-6">
                <!-- 修改密码 -->
                <div class="border-b border-gray-200 pb-6">
                  <h4 class="text-md font-medium text-gray-900 mb-4">修改密码</h4>
                  <form @submit.prevent="changePassword" class="space-y-4">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        当前密码
                      </label>
                      <input
                        v-model="passwordForm.currentPassword"
                        type="password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        新密码
                      </label>
                      <input
                        v-model="passwordForm.newPassword"
                        type="password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        确认新密码
                      </label>
                      <input
                        v-model="passwordForm.confirmPassword"
                        type="password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      />
                    </div>
                    <button
                      type="submit"
                      :disabled="isChangingPassword"
                      class="bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700 disabled:opacity-50"
                    >
                      {{ isChangingPassword ? '修改中...' : '修改密码' }}
                    </button>
                  </form>
                </div>

                <!-- 危险操作 -->
                <div>
                  <h4 class="text-md font-medium text-red-600 mb-4">危险操作</h4>
                  <button
                    @click="handleLogout"
                    class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
                  >
                    退出登录
                  </button>
                </div>
              </div>
            </div>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  UserIcon,
  FileTextIcon,
  SettingsIcon,
  LogOutIcon,
  CameraIcon,
  UsersIcon,
  HeartIcon
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import AvatarModal from '../components/AvatarModal.vue'
import { useToast } from '../composables/useToast'
import { UserApi } from '../api/user'
import { ArticleApi } from '../api/article'
import { FollowApi } from '../api/follow'
import type { ImageUploadResponse, Article } from '../api/types'

// 路由
const router = useRouter()

// 认证store
const authStore = useAuthStore()
const { toast } = useToast()

// 响应式数据
const showAvatarModal = ref(false)
const activeTab = ref('profile')
const isUpdating = ref(false)
const isChangingPassword = ref(false)

// 文章相关
const myArticles = ref<Article[]>([])
const isLoadingArticles = ref(false)

// 关注相关数据
const followingCount = ref(0)
const followersCount = ref(0)
const mutualCount = ref(0)
const articlesCount = ref(0)
const followingList = ref<any[]>([])
const followersList = ref<any[]>([])
const mutualFollowsList = ref<any[]>([])
const isLoadingFollows = ref(false)
const articlesError = ref('')

// 用户信息
const user = computed(() => authStore.user)

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
  { key: 'following', label: '我的关注', icon: HeartIcon },
  { key: 'followers', label: '我的粉丝', icon: UsersIcon },
  { key: 'mutual', label: '互相关注', icon: HeartIcon },
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

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
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
      
      // 调用API更新头像到后端
      const updateResponse = await UserApi.updateProfile({
        avatar: avatarUrl
      })
      
      // 更新本地用户信息
      authStore.updateUserInfo(updateResponse.data)
      
      // 确保DOM更新
      await nextTick()
      
      // 关闭头像上传弹窗
      closeAvatarModal()
      
      toast.success('头像更换成功')
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
      birthday: profileForm.birthday || undefined,
      bio: profileForm.bio,
      avatar: profileForm.avatar
    }
    
    // 移除空值
    Object.keys(updateData).forEach(key => {
      if (updateData[key] === '' || updateData[key] === undefined) {
        delete updateData[key]
      }
    })
    
    const response = await UserApi.updateProfile(updateData)
    
    // 更新本地用户信息
    authStore.updateUserInfo(response.data)
    
    toast.success('个人信息更新成功')
  } catch (error) {
    console.error('更新失败:', error)
    toast.error('更新失败，请重试')
  } finally {
    isUpdating.value = false
  }
}

// 修改密码
const changePassword = async () => {
  try {
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      alert('新密码和确认密码不一致')
      return
    }

    isChangingPassword.value = true
    // TODO: 调用修改密码API
    console.log('修改密码')
    
    // 重置表单
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error) {
    console.error('修改密码失败:', error)
  } finally {
    isChangingPassword.value = false
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}

// 加载我的文章
const loadMyArticles = async () => {
  try {
    isLoadingArticles.value = true
    articlesError.value = ''
    
    const response = await ArticleApi.getMyArticles({ page: 1, size: 20 })
    myArticles.value = (response.data && Array.isArray(response.data)) ? response.data : []
  } catch (error: any) {
    articlesError.value = error.message || '加载文章失败'
    console.error('加载我的文章失败:', error)
  } finally {
    isLoadingArticles.value = false
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

// 加载关注统计信息
const loadFollowStats = async () => {
  try {
    const response = await FollowApi.getFollowStats()
    console.log('Follow stats response:', response)
    // 检查响应结构
    if (response && response.data) {
      followingCount.value = response.data.following_count || 0
      followersCount.value = response.data.followers_count || 0
    } else {
      // 如果直接是数据对象
      followingCount.value = response.following_count || 0
      followersCount.value = response.followers_count || 0
    }
  } catch (error) {
    console.error('加载关注统计失败:', error)
  }
}

// 加载关注列表
const loadFollowing = async () => {
  try {
    isLoadingFollows.value = true
    const response = await FollowApi.getFollowing({ page: 1, limit: 50 })
    console.log('Following response:', response)
    // 检查响应结构
    if (response && response.data && response.data.users) {
      followingList.value = response.data.users
    } else if (response && response.users) {
      followingList.value = response.users
    } else {
      followingList.value = []
    }
  } catch (error) {
    console.error('加载关注列表失败:', error)
    toast.error('加载关注列表失败')
  } finally {
    isLoadingFollows.value = false
  }
}

// 加载粉丝列表
const loadFollowers = async () => {
  try {
    isLoadingFollows.value = true
    const response = await FollowApi.getFollowers({ page: 1, limit: 50 })
    console.log('Followers response:', response)
    // 检查响应结构
    if (response && response.data && response.data.users) {
      followersList.value = response.data.users
    } else if (response && response.users) {
      followersList.value = response.users
    } else {
      followersList.value = []
    }
  } catch (error) {
    console.error('加载粉丝列表失败:', error)
    toast.error('加载粉丝列表失败')
  } finally {
    isLoadingFollows.value = false
  }
}

// 取消关注
const unfollowUser = async (userId) => {
  try {
    await FollowApi.unfollowUser(userId)
    toast.success('取消关注成功')
    // 重新加载数据
    await loadFollowStats()
    if (activeTab.value === 'following') {
      await loadFollowing()
    }
  } catch (error) {
    console.error('取消关注失败:', error)
    toast.error('取消关注失败')
  }
}

// 关注用户
const followUser = async (userId) => {
  try {
    await FollowApi.followUser(userId)
    toast.success('关注成功')
    // 重新加载数据
    await loadFollowStats()
    if (activeTab.value === 'followers') {
      await loadFollowers()
    }
  } catch (error) {
    console.error('关注失败:', error)
    toast.error('关注失败')
  }
}

// 加载互相关注列表
const loadMutualFollows = async () => {
  try {
    isLoadingFollows.value = true
    const response = await FollowApi.getMutualFollows({ page: 1, limit: 50 })
    // 检查响应结构
    if (response && response.data && response.data.users) {
      mutualFollowsList.value = response.data.users
      mutualCount.value = response.data.total || response.data.users.length
    } else if (response && response.users) {
      mutualFollowsList.value = response.users
      mutualCount.value = response.total || response.users.length
    } else {
      mutualFollowsList.value = []
      mutualCount.value = 0
    }
  } catch (error) {
    console.error('加载互关列表失败:', error)
    toast.error('加载互关列表失败')
    mutualFollowsList.value = []
    mutualCount.value = 0
  } finally {
    isLoadingFollows.value = false
  }
}

// 监听activeTab变化，根据标签加载不同数据
watch(activeTab, (newTab) => {
  if (newTab === 'articles') {
    loadMyArticles()
  } else if (newTab === 'following') {
    loadFollowing()
  } else if (newTab === 'followers') {
    loadFollowers()
  } else if (newTab === 'mutual') {
    loadMutualFollows()
  }
})

// 组件挂载时初始化
onMounted(() => {
  // 检查登录状态
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  
  initUserInfo()
  loadFollowStats()
  
  // 如果默认是文章标签，立即加载文章
  if (activeTab.value === 'articles') {
    loadMyArticles()
  }
})
</script>