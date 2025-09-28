<template>
  <AppLayout background="gray" :show-footer="false">
    <PageContainer background="transparent" padding="md" max-width="7xl">
      <div class="flex gap-6">
        <!-- Left Sidebar - Topics -->
        <div class="w-80 bg-white rounded-lg shadow-sm border border-gray-200">
          <!-- Search Bar -->
          <div class="p-4 border-b border-gray-200">
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <SearchIcon class="h-5 w-5 text-gray-400" />
              </div>
              <input
                v-model="searchQuery"
                @keyup.enter="handleSearch"
                type="text"
                placeholder="搜索帖子..."
                class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent bg-gray-50 focus:bg-white transition-colors"
              />
            </div>
          </div>

          <!-- Topics List -->
          <div class="p-4">
            <h3 class="text-lg font-bold text-gray-900 mb-4">话题分类</h3>
            <div class="space-y-1">
              <div
                v-for="topic in topics"
                :key="topic.key"
                @click="setSelectedTopic(topic.key)"
                :class="[
                  'flex items-center px-4 py-3 rounded-lg cursor-pointer transition-all duration-200',
                  selectedTopic === topic.key
                    ? 'bg-pink-100 text-pink-700 border-l-4 border-pink-500'
                    : 'text-gray-700 hover:bg-gray-100'
                ]"
              >
                <span class="text-sm font-medium">{{ topic.label }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Main Content - Forum Posts -->
        <div class="flex-1 bg-white rounded-lg shadow-sm border border-gray-200">
          <!-- Page Header -->
          <div class="flex justify-between items-center p-6 border-b border-gray-200">
            <h1 class="text-2xl font-bold text-gray-900">育儿社区</h1>
            <button
              @click="handleNewPost"
              class="bg-gradient-to-r from-pink-500 to-rose-400 text-white px-6 py-2 rounded-lg hover:from-pink-600 hover:to-rose-500 transition-all duration-200 font-medium shadow-md hover:shadow-lg transform hover:scale-105"
            >
              发布新帖
            </button>
          </div>

          <!-- Forum Posts List -->
          <div class="divide-y divide-gray-200">
            <!-- 加载状态 -->
            <div v-if="isLoading && forumPosts.length === 0" class="flex flex-col items-center justify-center py-16">
              <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-pink-600 mb-4"></div>
              <p class="text-gray-500">加载中...</p>
            </div>

            <!-- 无帖子状态 -->
            <div v-else-if="!isLoading && filteredPosts.length === 0" class="flex flex-col items-center justify-center py-16">
              <div class="text-gray-400 mb-4">
                <UsersIcon class="h-16 w-16 mx-auto" />
              </div>
              <h3 class="text-lg font-medium text-gray-900 mb-2">暂无相关帖子</h3>
              <p class="text-gray-500 text-center">
                {{ forumPosts.length === 0 ? '还没有人发帖，来发第一个帖子吧！' : '尝试搜索其他关键词或选择不同的话题分类' }}
              </p>
            </div>

            <!-- 帖子列表 -->
            <ForumPost
              v-for="post in filteredPosts"
              :key="post.id"
              :post="post"
              @click="handlePostClick(post)"
            />
          </div>

          <!-- Load More -->
          <div v-if="filteredPosts.length > 0 && currentPage < totalPages" class="p-6 border-t border-gray-200">
            <div class="flex justify-center">
              <button
                @click="loadMore"
                :disabled="isLoading"
                class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ isLoading ? '加载中...' : '加载更多' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </PageContainer>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { SearchIcon, UsersIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { ForumApi } from '@/api/forum'
import { useToast } from '@/composables/useToast'
import { AppLayout, PageContainer } from '@/components/layout'
import ForumPost from '@/components/ForumPost.vue'
import type { ForumPost as ForumPostType, TopicConfig } from '@/api/types'

const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const selectedTopic = ref('All')
const searchQuery = ref('')
const isLoading = ref(false)
const forumPosts = ref<ForumPostType[]>([])
const currentPage = ref(1)
const totalPages = ref(1)
const totalPosts = ref(0)

// 话题分类
const topics = ref<TopicConfig[]>([
  { key: 'All', label: '全部' },
  { key: 'Baby Care', label: '婴儿护理' },
  { key: 'Feeding', label: '喂养' },
  { key: 'Sleep', label: '睡眠' },
  { key: 'Health', label: '健康' },
  { key: 'Development', label: '发育' },
  { key: 'Activities', label: '活动' },
  { key: 'Gear', label: '用品' },
  { key: 'Parenting', label: '育儿' },
  { key: 'Family Life', label: '家庭生活' },
  { key: 'Work & Life Balance', label: '工作生活平衡' },
  { key: 'Relationships', label: '人际关系' },
  { key: 'Mental Health', label: '心理健康' },
  { key: 'Finances', label: '财务' },
  { key: 'Legal', label: '法律' },
  { key: 'Other', label: '其他' }
])

// 计算属性：过滤后的帖子
const filteredPosts = computed(() => {
  if (!forumPosts.value || !Array.isArray(forumPosts.value)) {
    return []
  }
  const filtered = forumPosts.value.filter(post => {
    const matchesTopic = selectedTopic.value === 'All' || post.topic === selectedTopic.value
    const matchesSearch = post.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesTopic && matchesSearch
  })

  // 前端兜底排序：置顶优先、精华其次、再按创建时间倒序
  return filtered.slice().sort((a: any, b: any) => {
    const topA = a.is_top ? 1 : 0
    const topB = b.is_top ? 1 : 0
    if (topA !== topB) return topB - topA
    const hotA = a.is_hot ? 1 : 0
    const hotB = b.is_hot ? 1 : 0
    if (hotA !== hotB) return hotB - hotA
    const ta = new Date(a.created_at || a.createdAt || 0).getTime()
    const tb = new Date(b.created_at || b.createdAt || 0).getTime()
    return tb - ta
  })
})

// 方法
const handleNewPost = () => {
  if (authStore.isAuthenticated) {
    router.push('/community/posts/create')
  } else {
    router.push('/login')
  }
}

const handlePostClick = (post: ForumPostType) => {
  console.log('查看帖子:', post.title)
  router.push(`/community/posts/${post.id}`)
}

// 加载帖子列表
const loadPosts = async (page: number = 1) => {
  if (isLoading.value) return

  isLoading.value = true
  try {
    const params = {
      page,
      size: 10,
      topic: selectedTopic.value === 'All' ? undefined : selectedTopic.value,
      keyword: searchQuery.value.trim() || undefined,
      sort: 'created_at desc'
    }

    const response = await ForumApi.getPostList(params)
    const data = response.data

    forumPosts.value = data.items || []
    currentPage.value = data.page
    totalPages.value = data.total_pages || Math.ceil((data.total || 0) / (data.size || 1))
    totalPosts.value = data.total
  } catch (error: any) {
    console.error('加载帖子列表失败:', error)
    showToast('加载帖子列表失败', 'error')
  } finally {
    isLoading.value = false
  }
}

// 刷新帖子列表
const refreshPosts = () => {
  currentPage.value = 1
  loadPosts(1)
}

// 加载更多帖子
const loadMore = () => {
  if (currentPage.value < totalPages.value) {
    loadPosts(currentPage.value + 1)
  }
}

// 监听话题和搜索变化
const setSelectedTopic = (topic: string) => {
  selectedTopic.value = topic
  refreshPosts()
}

// 搜索处理
const handleSearch = () => {
  refreshPosts()
}

// 组件挂载时加载数据
onMounted(() => {
  loadPosts()
})
</script>
