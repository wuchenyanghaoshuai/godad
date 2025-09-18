<template>
  <div class="min-h-screen bg-gray-50">
    <UniversalHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
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
            <div v-if="filteredPosts.length === 0" class="flex flex-col items-center justify-center py-16">
              <div class="text-gray-400 mb-4">
                <UsersIcon class="h-16 w-16 mx-auto" />
              </div>
              <h3 class="text-lg font-medium text-gray-900 mb-2">暂无相关帖子</h3>
              <p class="text-gray-500 text-center">
                尝试搜索其他关键词或选择不同的话题分类
              </p>
            </div>

            <ForumPost
              v-for="post in filteredPosts"
              :key="post.id"
              :post="post"
              @click="handlePostClick(post)"
            />
          </div>

          <!-- Load More -->
          <div v-if="filteredPosts.length > 0" class="p-6 border-t border-gray-200">
            <div class="flex justify-center">
              <button
                class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
              >
                加载更多
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { SearchIcon, UsersIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import UniversalHeader from '@/components/UniversalHeader.vue'
import ForumPost from '@/components/ForumPost.vue'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const selectedTopic = ref('All')
const searchQuery = ref('')

// 话题分类
const topics = [
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
]

// 论坛帖子数据
const forumPosts = [
  {
    id: 1,
    title: '4个月宝宝睡眠倒退怎么办？',
    author: {
      name: '新手妈妈123',
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b5bc?auto=format&fit=crop&q=80&w=200'
    },
    replies: 12,
    views: 234,
    timeAgo: '2天前',
    topic: 'Sleep'
  },
  {
    id: 2,
    title: '大房子用什么婴儿监视器比较好？',
    author: {
      name: '爸爸小王',
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=200'
    },
    replies: 8,
    views: 187,
    timeAgo: '3天前',
    topic: 'Gear'
  },
  {
    id: 3,
    title: '6个月宝宝添加辅食的技巧分享',
    author: {
      name: '营养师小李',
      avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?auto=format&fit=crop&q=80&w=200'
    },
    replies: 15,
    views: 312,
    timeAgo: '4天前',
    topic: 'Feeding'
  },
  {
    id: 4,
    title: '产后焦虑症如何应对？',
    author: {
      name: '心理咨询师',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=crop&q=80&w=200'
    },
    replies: 10,
    views: 205,
    timeAgo: '5天前',
    topic: 'Mental Health'
  },
  {
    id: 5,
    title: '城市生活用什么婴儿车最好？',
    author: {
      name: '城市妈妈',
      avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&q=80&w=200'
    },
    replies: 7,
    views: 156,
    timeAgo: '6天前',
    topic: 'Gear'
  },
  {
    id: 6,
    title: '3个月宝宝的刺激活动推荐',
    author: {
      name: '早教老师',
      avatar: 'https://images.unsplash.com/photo-1489424731084-a5d8b219a5bb?auto=format&fit=crop&q=80&w=200'
    },
    replies: 11,
    views: 220,
    timeAgo: '7天前',
    topic: 'Activities'
  },
  {
    id: 7,
    title: '新生儿来了如何处理兄弟姐妹争宠？',
    author: {
      name: '二胎妈妈',
      avatar: 'https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?auto=format&fit=crop&q=80&w=200'
    },
    replies: 9,
    views: 198,
    timeAgo: '8天前',
    topic: 'Family Life'
  },
  {
    id: 8,
    title: '徒步旅行用什么婴儿背带？',
    author: {
      name: '户外爸爸',
      avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&q=80&w=200'
    },
    replies: 13,
    views: 250,
    timeAgo: '9天前',
    topic: 'Gear'
  },
  {
    id: 9,
    title: '带宝宝旅行的实用技巧',
    author: {
      name: '旅行达人',
      avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=200'
    },
    replies: 6,
    views: 145,
    timeAgo: '10天前',
    topic: 'Family Life'
  },
  {
    id: 10,
    title: '母乳宝宝如何接受奶瓶？',
    author: {
      name: '哺乳顾问',
      avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&q=80&w=200'
    },
    replies: 14,
    views: 280,
    timeAgo: '11天前',
    topic: 'Feeding'
  }
]

// 计算属性：过滤后的帖子
const filteredPosts = computed(() => {
  return forumPosts.filter(post => {
    const matchesTopic = selectedTopic.value === 'All' || post.topic === selectedTopic.value
    const matchesSearch = post.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesTopic && matchesSearch
  })
})

// 方法
const setSelectedTopic = (topic: string) => {
  selectedTopic.value = topic
}

const handleNewPost = () => {
  if (authStore.isAuthenticated) {
    console.log('创建新帖子')
    // 这里可以跳转到创建帖子页面
  } else {
    router.push('/login')
  }
}

const handlePostClick = (post: any) => {
  console.log('查看帖子:', post.title)
  // 这里可以跳转到帖子详情页面
}
</script>