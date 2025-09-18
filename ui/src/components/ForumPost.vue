<template>
  <div
    @click="$emit('click', post)"
    class="p-6 hover:bg-gray-50 transition-colors cursor-pointer"
  >
    <div class="flex items-start space-x-4">
      <!-- 用户头像 -->
      <div class="flex-shrink-0">
        <img
          :src="post.author.avatar"
          :alt="post.author.name"
          class="w-12 h-12 rounded-full object-cover border-2 border-gray-200"
        />
      </div>

      <!-- 帖子内容 -->
      <div class="flex-1 min-w-0">
        <!-- 帖子标题 -->
        <h3 class="text-lg font-semibold text-gray-900 hover:text-pink-600 transition-colors mb-2 line-clamp-2">
          {{ post.title }}
        </h3>

        <!-- 作者和时间 -->
        <div class="flex items-center text-sm text-gray-500 mb-3">
          <span class="font-medium text-gray-700">{{ post.author.name }}</span>
          <span class="mx-2">·</span>
          <span>{{ post.timeAgo }}</span>
          <span class="mx-2">·</span>
          <span class="bg-gray-100 text-gray-600 px-2 py-1 rounded-full text-xs">
            {{ getTopicLabel(post.topic) }}
          </span>
        </div>

        <!-- 统计信息 -->
        <div class="flex items-center space-x-6 text-sm text-gray-500">
          <div class="flex items-center space-x-1">
            <MessageSquareIcon class="h-4 w-4" />
            <span>{{ post.replies }} 回复</span>
          </div>
          <div class="flex items-center space-x-1">
            <EyeIcon class="h-4 w-4" />
            <span>{{ post.views }} 浏览</span>
          </div>
        </div>
      </div>

      <!-- 右侧箭头 -->
      <div class="flex-shrink-0">
        <ChevronRightIcon class="h-5 w-5 text-gray-400" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MessageSquareIcon, EyeIcon, ChevronRightIcon } from 'lucide-vue-next'

// Props
interface Author {
  name: string
  avatar: string
}

interface Post {
  id: number
  title: string
  author: Author
  replies: number
  views: number
  timeAgo: string
  topic: string
}

defineProps<{
  post: Post
}>()

// Emits
defineEmits<{
  click: [post: Post]
}>()

// 话题标签映射
const topicLabels: Record<string, string> = {
  'All': '全部',
  'Baby Care': '婴儿护理',
  'Feeding': '喂养',
  'Sleep': '睡眠',
  'Health': '健康',
  'Development': '发育',
  'Activities': '活动',
  'Gear': '用品',
  'Parenting': '育儿',
  'Family Life': '家庭生活',
  'Work & Life Balance': '工作生活平衡',
  'Relationships': '人际关系',
  'Mental Health': '心理健康',
  'Finances': '财务',
  'Legal': '法律',
  'Other': '其他'
}

const getTopicLabel = (topic: string): string => {
  return topicLabels[topic] || topic
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>