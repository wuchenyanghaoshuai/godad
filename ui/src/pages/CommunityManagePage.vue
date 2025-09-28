<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-6">
            <h1 class="text-xl font-semibold text-gray-900 flex items-center">
              <MessageSquareIcon class="h-6 w-6 text-pink-600 mr-2" />
              社区管理
            </h1>
            <nav class="flex space-x-4">
              <AdminNavLink to="/management-dashboard" label="系统管理" />
              <AdminNavLink to="/community-management" label="社区管理" />
              <AdminNavLink to="/report-center" label="举报中心" :badge="pendingCount" />
            </nav>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-600">欢迎，{{ authStore.user?.nickname }}</span>
            <button
              @click="logout"
              class="text-sm text-red-600 hover:text-red-700"
            >
              退出登录
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 社区统计卡片 -->
      <div class="mb-8">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- 论坛帖子总数 -->
        <div class="bg-white rounded-lg shadow p-6 border-l-4 border-pink-500">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-pink-100 text-pink-600">
              <FileTextIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">论坛帖子</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.forumPostCount }}</p>
              <p class="text-xs text-green-600 flex items-center mt-1">
                <TrendingUpIcon class="h-3 w-3 mr-1" />
                +{{ stats.todayNewPosts }} 今日新增
              </p>
            </div>
          </div>
        </div>

        <!-- 论坛回复总数 -->
        <div class="bg-white rounded-lg shadow p-6 border-l-4 border-blue-500">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-blue-100 text-blue-600">
              <MessageCircleIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">论坛回复</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.forumReplyCount }}</p>
              <p class="text-xs text-green-600 flex items-center mt-1">
                <TrendingUpIcon class="h-3 w-3 mr-1" />
                +{{ stats.todayNewReplies }} 今日新增
              </p>
            </div>
          </div>
        </div>

        <!-- 活跃用户数 -->
        <div class="bg-white rounded-lg shadow p-6 border-l-4 border-green-500">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-green-100 text-green-600">
              <UsersIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">活跃用户</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.activeUserCount }}</p>
              <p class="text-xs text-gray-500">最近7天发帖用户</p>
            </div>
          </div>
        </div>

        <!-- 精华帖子数 -->
        <div class="bg-white rounded-lg shadow p-6 border-l-4 border-yellow-500">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-yellow-100 text-yellow-600">
              <StarIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">精华帖子</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.featuredPostCount }}</p>
              <p class="text-xs text-gray-500">置顶: {{ stats.pinnedPostCount }}</p>
            </div>
          </div>
        </div>
          </div>
        </div>
      </div>

      <!-- 主内容区域 -->
      <div class="flex flex-col gap-6">
        <div class="w-full bg-white rounded-lg shadow-sm border border-gray-200">
          <!-- 内容标题栏 -->
          <div class="flex justify-between items-center p-6 border-b border-gray-200">
            <div>
              <h2 class="text-lg font-semibold text-gray-900">{{ getCurrentTabName() }}</h2>
              <p class="text-sm text-gray-600 mt-1">{{ getCurrentTabDescription() }}</p>
            </div>
            <div class="flex items-center space-x-3"></div>
          </div>

          <!-- 顶部标签切换 -->
          <div class="border-b border-gray-200">
            <nav class="-mb-px flex px-6">
              <button
                @click="activeTab = 'posts'"
                :class="[
                  'whitespace-nowrap py-4 px-4 border-b-2 font-medium text-sm',
                  activeTab === 'posts' ? 'border-pink-500 text-pink-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                ]"
              >
                社区帖子管理
              </button>
              
            </nav>
          </div>

          <!-- 标签页内容 -->
          <div class="p-6">
            <!-- 论坛帖子管理 -->
            <div v-if="activeTab === 'posts'">
              <div class="flex justify-between items-center mb-4">
                <div class="flex space-x-2">
                  <input v-model="postFilter.keyword" @input="searchForumPosts" type="text" placeholder="搜索帖子标题或内容..." class="px-3 py-2 border border-gray-300 rounded-md text-sm" >
                  <select v-model="postFilter.topic" @change="loadForumPosts" class="px-3 py-2 border border-gray-300 rounded-md text-sm">
                    <option value="">全部话题</option>
                    <option value="Baby Care">婴儿护理</option>
                    <option value="Feeding">喂养</option>
                    <option value="Sleep">睡眠</option>
                    <option value="Health">健康</option>
                    <option value="Development">发育</option>
                    <option value="Other">其他</option>
                  </select>
                  <select v-model="postFilter.status" @change="loadForumPosts" class="px-3 py-2 border border-gray-300 rounded-md text-sm">
                    <option value="">全部状态</option>
                    <option value="normal">普通</option>
                    <option value="pinned">置顶</option>
                    <option value="featured">精华</option>
                  </select>
                </div>
              </div>
              <div class="overflow-x-auto border border-gray-200 rounded-lg">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-50">
                    <tr>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">帖子标题</th>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">作者</th>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">话题</th>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">回复数</th>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">浏览数</th>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
                      <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                    </tr>
                  </thead>
                  <tbody class="bg-white divide-y divide-gray-200">
                    <tr v-for="post in forumPosts" :key="post.id" class="hover:bg-gray-50">
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="flex items-center">
                          <span v-if="post.is_top" class="px-2 py-1 text-xs bg-yellow-100 text-yellow-800 rounded mr-2">置顶</span>
                          <span v-if="post.is_hot" class="px-2 py-1 text-xs bg-red-100 text-red-800 rounded mr-2">精华</span>
                          <div class="text-sm font-medium text-gray-900">{{ post.title }}</div>
                        </div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ post.author?.nickname || post.author?.username || '匿名' }}</td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <span class="px-2 py-1 text-xs bg-pink-100 text-pink-800 rounded">{{ getTopicLabel(post.topic) }}</span>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ post.reply_count || 0 }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ post.view_count || 0 }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDateTime(post.created_at) }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <button class="text-blue-600 hover:text-blue-900 mr-2" @click="viewForumPost(post.id)">查看</button>
                        <button class="text-yellow-600 hover:text-yellow-900 mr-2" @click="toggleForumTop(post)">{{ post.is_top ? '取消置顶' : '置顶' }}</button>
                        <button class="text-amber-600 hover:text-amber-900 mr-2" @click="toggleForumHot(post)">{{ post.is_hot ? '取消精华' : '精华' }}</button>
                        <button class="text-gray-600 hover:text-gray-900 mr-2" @click="toggleForumLock(post)">{{ post.is_locked ? '解锁' : '锁定' }}</button>
                        <button class="text-red-600 hover:text-red-900" @click="deleteForumPost(post)">删除</button>
                      </td>
                    </tr>
                    <tr v-if="!forumPosts.length">
                      <td colspan="7" class="px-6 py-8 text-center text-gray-500">暂无帖子</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="flex justify-between items-center mt-4">
                <span class="text-sm text-gray-700">显示 {{ forumPostPagination.offset + 1 }} 到 {{ Math.min(forumPostPagination.offset + forumPostPagination.limit, forumPostPagination.total) }} 条，共 {{ forumPostPagination.total }} 条</span>
                <div class="flex space-x-2">
                  <button @click="prevPage('forum')" :disabled="forumPostPagination.page === 1" class="px-3 py-1 border rounded-md text-sm disabled:opacity-50">上一页</button>
                  <button @click="nextPage('forum')" :disabled="forumPostPagination.page >= Math.ceil(forumPostPagination.total / forumPostPagination.limit)" class="px-3 py-1 border rounded-md text-sm disabled:opacity-50">下一页</button>
                </div>
              </div>
            </div>

            

            <!-- 社区统计 -->
            <div v-if="activeTab === 'stats'">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <!-- 发帖趋势图占位 -->
                <div class="bg-gray-50 rounded-lg p-6 h-64 flex items-center justify-center">
                  <div class="text-center">
                    <BarChartIcon class="h-12 w-12 mx-auto text-gray-400 mb-3" />
                    <p class="text-gray-500">发帖趋势图</p>
                    <p class="text-sm text-gray-400 mt-1">功能开发中...</p>
                  </div>
                </div>

                <!-- 热门话题排行占位 -->
                <div class="bg-gray-50 rounded-lg p-6 h-64 flex items-center justify-center">
                  <div class="text-center">
                    <PieChartIcon class="h-12 w-12 mx-auto text-gray-400 mb-3" />
                    <p class="text-gray-500">热门话题排行</p>
                    <p class="text-sm text-gray-400 mt-1">功能开发中...</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- 内容审核 -->
            <div v-if="activeTab === 'moderation'">
              <div class="text-center py-12">
                <ShieldCheckIcon class="h-12 w-12 mx-auto text-gray-400 mb-4" />
                <h3 class="text-lg font-medium text-gray-900 mb-2">内容审核中心</h3>
                <p class="text-gray-500">暂无需要审核的内容</p>
                <p class="text-sm text-gray-400 mt-1">功能开发中...</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 批量操作模态框 -->
    <div v-if="showBatchActionsModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">批量操作 ({{ selectedPosts.length }} 个帖子)</h3>
        <div class="space-y-3">
          <button
            @click="batchTogglePin()"
            class="w-full flex items-center px-4 py-3 text-left text-sm bg-gray-50 hover:bg-gray-100 rounded-lg"
          >
            <PinIcon class="h-4 w-4 mr-3 text-gray-500" />
            批量置顶/取消置顶
          </button>
          <button
            @click="batchToggleFeature()"
            class="w-full flex items-center px-4 py-3 text-left text-sm bg-gray-50 hover:bg-gray-100 rounded-lg"
          >
            <StarIcon class="h-4 w-4 mr-3 text-gray-500" />
            批量设为精华/取消精华
          </button>
          <button
            @click="batchDelete()"
            class="w-full flex items-center px-4 py-3 text-left text-sm bg-red-50 hover:bg-red-100 rounded-lg text-red-700"
          >
            <TrashIcon class="h-4 w-4 mr-3" />
            批量删除
          </button>
        </div>
        <div class="flex justify-end space-x-2 mt-6">
          <button
            @click="showBatchActionsModal = false"
            class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
          >
            取消
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { usePendingReports } from '@/composables/usePendingReports'
import AdminNavLink from '@/components/AdminNavLink.vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import AdminApi from '@/api/admin'
import { ForumApi } from '@/api/forum'

// Icons
import {
  MessageSquareIcon,
  FileTextIcon,
  MessageCircleIcon,
  UsersIcon,
  StarIcon,
  TrendingUpIcon,
  SearchIcon,
  EyeIcon,
  PinIcon,
  TrashIcon,
  DownloadIcon,
  BarChartIcon,
  PieChartIcon,
  ShieldCheckIcon
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()
const { pendingCount, refresh: refreshPendingReports, startPolling: pollPendingReports } = usePendingReports()

// 响应式数据
const activeTab = ref('posts')
const selectedPosts = ref<number[]>([])
const showBatchActionsModal = ref(false)

// 社区统计数据（从后端拉取）
const stats = reactive({
  forumPostCount: 0,
  forumReplyCount: 0,
  activeUserCount: 0,
  featuredPostCount: 0,
  pinnedPostCount: 0,
  todayNewPosts: 0,
  todayNewReplies: 0
})

const loadForumStats = async () => {
  try {
    const res = await AdminApi.getForumStats()
    const s = res?.data || {}
    stats.forumPostCount = s.forumPostCount ?? 0
    stats.forumReplyCount = s.forumReplyCount ?? 0
    stats.activeUserCount = s.activeUserCount ?? 0
    stats.featuredPostCount = s.featuredPostCount ?? 0
    stats.pinnedPostCount = s.pinnedPostCount ?? 0
    stats.todayNewPosts = s.todayNewPosts ?? 0
    stats.todayNewReplies = s.todayNewReplies ?? 0
  } catch (e) {
    console.error('Failed to load forum stats:', e)
    showToast('获取论坛统计数据失败，请检查权限或重新登录', 'error')
  }
}

// 顶部标签
const tabs = [
  { id: 'posts', name: '社区帖子管理', icon: FileTextIcon, badge: null },
  { id: 'replies', name: '回复管理', icon: MessageCircleIcon, badge: null },
]

// 帖子管理
const posts = ref<any[]>([
  {
    id: 1,
    title: 'DNF装备求助 - 新手玩家装备选择指导',
    content: '大家好，我是个DNF新手，现在角色等级35了，不知道该选择什么装备比较好...',
    topic: 'Other',
    author: {
      id: 3,
      nickname: '淘气宝石',
      username: 'chenyangwu',
      avatar: 'https://godad.kaikeba.com.cn/UserImage/41f84f30-a506-4017-b716-23c214504987.jpg'
    },
    view_count: 156,
    reply_count: 23,
    is_pinned: true,
    is_featured: false,
    created_at: '2025-09-26T10:30:00Z'
  },
  {
    id: 2,
    title: '经验交流：婴儿护理的实用技巧',
    content: '作为新手妈妈，想分享一些我在婴儿护理方面的心得体会...',
    topic: 'Baby Care',
    author: {
      id: 5,
      nickname: '温柔妈妈',
      username: 'dongxueli',
      avatar: '/default-avatar.png'
    },
    view_count: 89,
    reply_count: 12,
    is_pinned: false,
    is_featured: true,
    created_at: '2025-09-26T09:15:00Z'
  }
])

const postFilters = reactive({
  keyword: '',
  topic: '',
  status: ''
})

const postPagination = reactive({
  page: 1,
  limit: 10,
  total: 2
})

// 回复管理
/* const replies = ref<any[]>([
  {
    id: 1,
    content: '这个建议很实用，我试试看！感谢分享~',
    post_id: 1,
    post: { title: 'DNF装备求助 - 新手玩家装备选择指导' },
    author: {
      id: 16,
      nickname: '新手小白',
      username: 'zhangzhi',
      avatar: '/default-avatar.png'
    },
    created_at: '2025-09-26T11:00:00Z'
  },
  {
    id: 2,
    content: '我也是新手妈妈，这些技巧对我很有帮助，收藏了！',
    post_id: 2,
    post: { title: '经验交流：婴儿护理的实用技巧' },
    author: {
      id: 18,
      nickname: '学习妈妈',
      username: 'mama123',
      avatar: '/default-avatar.png'
    },
    created_at: '2025-09-26T10:45:00Z'
  }
]) */

// const replyFilters = reactive({ keyword: '', postId: '' })

// const replyPagination = reactive({ page: 1, limit: 15, total: 2 })

// 与后台一致的社区帖子管理数据/方法
const forumPosts = ref<any[]>([])
const postFilter = reactive({ keyword: '', topic: '', status: '' })
const forumPostPagination = reactive({ page: 1, limit: 10, total: 0, offset: 0 })

const formatDateTime = (str: string | Date) => {
  const d = typeof str === 'string' ? new Date(str) : str
  return d.toLocaleString()
}

const loadForumPosts = async () => {
  try {
    const params: any = { page: forumPostPagination.page, size: forumPostPagination.limit }
    if (postFilter.keyword?.trim()) params.keyword = postFilter.keyword.trim()
    if (postFilter.topic) params.topic = postFilter.topic
    if (postFilter.status === 'pinned') params.is_top = true
    else if (postFilter.status === 'featured') params.is_hot = true
    const res = await AdminApi.getForumPostsPage(params)
    const page = res.data
    forumPosts.value = page.items || []
    forumPostPagination.total = page.total || 0
    forumPostPagination.offset = (page.page - 1) * page.size
  } catch (error) {
    // 后台端点不可用时回退到公开列表
    try {
      const params: any = { page: forumPostPagination.page, size: forumPostPagination.limit }
      if (postFilter.keyword?.trim()) params.keyword = postFilter.keyword.trim()
      if (postFilter.topic) params.topic = postFilter.topic
      if (postFilter.status === 'pinned') params.is_top = true
      else if (postFilter.status === 'featured') params.is_hot = true
      const res = await ForumApi.getPostList(params)
      const page = res.data
      forumPosts.value = page.items || []
      forumPostPagination.total = page.total || 0
      forumPostPagination.offset = (page.page - 1) * page.size
    } catch (e) {
      forumPosts.value = []
      forumPostPagination.total = 0
    }
  }
}

let forumPostSearchTimeout: any = null
const searchForumPosts = () => {
  clearTimeout(forumPostSearchTimeout)
  forumPostSearchTimeout = setTimeout(() => {
    forumPostPagination.page = 1
    loadForumPosts()
  }, 500)
}

const viewForumPost = (id: number) => {
  router.push(`/community/posts/${id}`)
}

const toggleForumTop = async (post: any) => {
  try {
    await AdminApi.toggleForumPostTop(post.id, !post.is_top)
    post.is_top = !post.is_top
    showToast(`已${post.is_top ? '置顶' : '取消置顶'}`, 'success')
  } catch (e) { showToast('操作失败，请重试', 'error') }
}

const toggleForumHot = async (post: any) => {
  try {
    await AdminApi.toggleForumPostHot(post.id, !post.is_hot)
    post.is_hot = !post.is_hot
    showToast(`已${post.is_hot ? '设为精华' : '取消精华'}`, 'success')
  } catch (e) { showToast('操作失败，请重试', 'error') }
}

const toggleForumLock = async (post: any) => {
  try {
    const locked = !!post.is_locked
    await AdminApi.toggleForumPostLock(post.id, !locked)
    post.is_locked = !locked
    showToast(`帖子已${post.is_locked ? '锁定' : '解锁'}`, 'success')
  } catch (e) { showToast('操作失败，请重试', 'error') }
}

const deleteForumPost = async (post: any) => {
  if (!confirm('确定要删除该帖子吗？')) return
  try {
    await AdminApi.deleteForumPost(post.id)
    const newTotal = forumPostPagination.total - 1
    const maxPage = Math.ceil(newTotal / forumPostPagination.limit)
    if (forumPostPagination.page > maxPage && maxPage > 0) forumPostPagination.page = maxPage
    else if (maxPage === 0) forumPostPagination.page = 1
    await loadForumPosts()
    showToast('删除成功', 'success')
  } catch (e) { showToast('删除失败，请重试', 'error') }
}

// 计算属性和方法
const getCurrentTabName = () => {
  const tab = tabs.find(t => t.id === activeTab.value)
  return tab?.name || '未知'
}

const getCurrentTabDescription = () => {
  const descriptions: Record<string, string> = {
    posts: '管理论坛帖子，包括置顶、精华、删除等操作',
    replies: '管理论坛回复，查看和删除不当回复',
    stats: '查看社区数据统计和趋势分析',
    moderation: '审核举报内容，维护社区秩序'
  }
  return descriptions[activeTab.value] || ''
}

const getTopicLabel = (topic: string) => {
  const labels: Record<string, string> = {
    'Baby Care': '婴儿护理',
    'Feeding': '喂养',
    'Sleep': '睡眠',
    'Health': '健康',
    'Development': '发育',
    'Other': '其他'
  }
  return labels[topic] || topic
}

const truncateText = (text: string, maxLength: number) => {
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 帖子操作
const toggleAllPosts = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.checked) {
    selectedPosts.value = posts.value.map(p => p.id)
  } else {
    selectedPosts.value = []
  }
}

const viewPost = (post: any) => {
  window.open(`/community/posts/${post.id}`, '_blank')
}

const togglePin = (post: any) => {
  post.is_pinned = !post.is_pinned
  showToast(`帖子已${post.is_pinned ? '置顶' : '取消置顶'}`, 'success')
}

const toggleFeature = (post: any) => {
  post.is_featured = !post.is_featured
  showToast(`帖子已${post.is_featured ? '设为精华' : '取消精华'}`, 'success')
}

const deletePost = (postId: number) => {
  if (!confirm('确定要删除这个帖子吗？')) return
  posts.value = posts.value.filter(p => p.id !== postId)
  selectedPosts.value = selectedPosts.value.filter(id => id !== postId)
  showToast('帖子删除成功', 'success')
}

const deleteReply = (replyId: number) => {
  if (!confirm('确定要删除这个回复吗？')) return
  replies.value = replies.value.filter(r => r.id !== replyId)
  showToast('回复删除成功', 'success')
}

// 批量操作
const batchTogglePin = () => {
  showToast(`已对 ${selectedPosts.value.length} 个帖子执行置顶操作`, 'success')
  selectedPosts.value = []
  showBatchActionsModal.value = false
}

const batchToggleFeature = () => {
  showToast(`已对 ${selectedPosts.value.length} 个帖子执行精华操作`, 'success')
  selectedPosts.value = []
  showBatchActionsModal.value = false
}

const batchDelete = () => {
  if (!confirm(`确定要删除选中的 ${selectedPosts.value.length} 个帖子吗？`)) return
  posts.value = posts.value.filter(p => !selectedPosts.value.includes(p.id))
  showToast(`已删除 ${selectedPosts.value.length} 个帖子`, 'success')
  selectedPosts.value = []
  showBatchActionsModal.value = false
}

// 搜索和分页
const searchPosts = () => {
  // TODO: 实现搜索逻辑
}

// const searchReplies = () => {}

const loadPosts = () => {
  // TODO: 实现加载逻辑
}

// const loadReplies = () => {}

const prevPage = (type: string) => {
  if (type === 'forum' && forumPostPagination.page > 1) {
    forumPostPagination.page--
    loadForumPosts()
  }
}

const nextPage = (type: string) => {
  if (type === 'forum') {
    const maxPage = Math.ceil(forumPostPagination.total / forumPostPagination.limit)
    if (forumPostPagination.page < maxPage) {
      forumPostPagination.page++
      loadForumPosts()
    }
  }
}

const exportPosts = () => {
  showToast('数据导出功能开发中...', 'info')
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}

// 初始化
onMounted(() => {
  // 检查管理员权限
  if (!authStore.user || authStore.user.role !== 2) {
    router.push('/')
    return
  }

  // 加载帖子列表
  loadForumPosts()
  // 加载论坛统计
  loadForumStats()
  // 刷新未处理举报数量并开始轮询
  refreshPendingReports(); pollPendingReports(60000)
  // loadReplies()
})

watch(activeTab, (val) => {
  if (val === 'posts') {
    forumPostPagination.page = 1
    loadForumPosts()
  }
})
</script>

<style scoped>
/* 自定义样式 */
.custom-shadow {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}
</style>
