<template>
  <div class="max-w-4xl mx-auto p-6">
    <div class="bg-white rounded-lg shadow-sm border border-gray-200">
      <!-- 表单标题 -->
      <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-xl font-semibold text-gray-900">
          {{ isEditing ? '编辑帖子' : '发布新帖' }}
        </h2>
        <p class="text-sm text-gray-600 mt-1">
          分享您的育儿经验，与其他家长交流讨论
        </p>
      </div>

      <!-- 表单内容 -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
        <!-- 帖子标题 -->
        <div>
          <label for="title" class="block text-sm font-medium text-gray-700 mb-2">
            标题 <span class="text-red-500">*</span>
          </label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            required
            maxlength="200"
            placeholder="请输入帖子标题"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors"
            :class="{ 'border-red-300': errors.title }"
          />
          <div class="flex justify-between items-center mt-1">
            <p v-if="errors.title" class="text-sm text-red-600">{{ errors.title }}</p>
            <span class="text-sm text-gray-500">{{ form.title.length }}/200</span>
          </div>
        </div>

        <!-- 话题分类 -->
        <div>
          <label for="topic" class="block text-sm font-medium text-gray-700 mb-2">
            话题分类 <span class="text-red-500">*</span>
          </label>
          <select
            id="topic"
            v-model="form.topic"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors"
            :class="{ 'border-red-300': errors.topic }"
          >
            <option value="">请选择话题分类</option>
            <option v-for="topic in topics" :key="topic.key" :value="topic.key">
              {{ topic.label }}
            </option>
          </select>
          <p v-if="errors.topic" class="text-sm text-red-600 mt-1">{{ errors.topic }}</p>
        </div>

        <!-- 帖子内容 -->
        <div>
          <label for="content" class="block text-sm font-medium text-gray-700 mb-2">
            内容 <span class="text-red-500">*</span>
          </label>
          <textarea
            id="content"
            v-model="form.content"
            rows="12"
            required
            maxlength="10000"
            placeholder="请输入帖子内容，可以详细描述您的问题或分享经验..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-vertical transition-colors"
            :class="{ 'border-red-300': errors.content }"
          />
          <div class="flex justify-between items-center mt-1">
            <p v-if="errors.content" class="text-sm text-red-600">{{ errors.content }}</p>
            <span class="text-sm text-gray-500">{{ form.content.length }}/10000</span>
          </div>
        </div>

        <!-- 发布选项 -->
        <div class="bg-gray-50 rounded-lg p-4 border border-gray-200">
          <h3 class="text-sm font-medium text-gray-700 mb-3">发布设置</h3>
          <div class="space-y-2">
            <label class="flex items-center">
              <input
                v-model="form.publishImmediately"
                type="checkbox"
                class="rounded border-gray-300 text-pink-600 focus:ring-pink-500"
              />
              <span class="ml-2 text-sm text-gray-700">立即发布</span>
            </label>
            <p class="text-xs text-gray-500">
              取消勾选将保存为草稿，您可以稍后再发布
            </p>
          </div>
        </div>

        <!-- 表单按钮 -->
        <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
          <button
            type="button"
            @click="handleCancel"
            class="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors font-medium"
            :disabled="isSubmitting"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            class="px-6 py-2 bg-gradient-to-r from-pink-500 to-rose-400 text-white rounded-lg hover:from-pink-600 hover:to-rose-500 transition-all duration-200 font-medium shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isSubmitting" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isEditing ? '更新中...' : '发布中...' }}
            </span>
            <span v-else>
              {{ isEditing ? '更新帖子' : (form.publishImmediately ? '发布帖子' : '保存草稿') }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ForumApi } from '@/api/forum'
import { useToast } from '@/composables/useToast'
import type { ForumPost, ForumPostCreateRequest, ForumPostUpdateRequest, TopicConfig } from '@/api/types'

// Props
interface Props {
  post?: ForumPost // 编辑时传入的帖子数据
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  success: [post: ForumPost]
  cancel: []
}>()

// 路由和工具
const router = useRouter()
const { showToast } = useToast()

// 表单数据
const form = ref({
  title: '',
  content: '',
  topic: '',
  publishImmediately: true
})

// 表单状态
const isSubmitting = ref(false)
const errors = ref<Record<string, string>>({})

// 话题选项
const topics = ref<TopicConfig[]>([
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

// 计算属性
const isEditing = computed(() => !!props.post)

const isFormValid = computed(() => {
  return form.value.title.trim() !== '' &&
         form.value.content.trim() !== '' &&
         form.value.topic !== '' &&
         Object.keys(errors.value).length === 0
})

// 验证表单
const validateForm = () => {
  errors.value = {}

  if (!form.value.title.trim()) {
    errors.value.title = '请输入帖子标题'
  } else if (form.value.title.length > 200) {
    errors.value.title = '标题不能超过200个字符'
  }

  if (!form.value.content.trim()) {
    errors.value.content = '请输入帖子内容'
  } else if (form.value.content.length > 10000) {
    errors.value.content = '内容不能超过10000个字符'
  }

  if (!form.value.topic) {
    errors.value.topic = '请选择话题分类'
  }

  return Object.keys(errors.value).length === 0
}

// 处理表单提交
const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }

  isSubmitting.value = true

  try {
    if (isEditing.value && props.post) {
      // 更新帖子
      const updateData: ForumPostUpdateRequest = {
        title: form.value.title.trim(),
        content: form.value.content.trim(),
        topic: form.value.topic,
        status: form.value.publishImmediately ? 1 : 0
      }

      const response = await ForumApi.updatePost(props.post.id, updateData)
      showToast('帖子更新成功', 'success')
      emit('success', response.data)
    } else {
      // 创建新帖子
      const createData: ForumPostCreateRequest = {
        title: form.value.title.trim(),
        content: form.value.content.trim(),
        topic: form.value.topic
      }

      const response = await ForumApi.createPost(createData)
      showToast('帖子发布成功', 'success')
      emit('success', response.data)
    }
  } catch (error: any) {
    console.error('提交帖子失败:', error)
    showToast(error.message || '操作失败，请重试', 'error')
  } finally {
    isSubmitting.value = false
  }
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}

// 初始化表单数据
const initializeForm = () => {
  if (props.post) {
    form.value = {
      title: props.post.title,
      content: props.post.content,
      topic: props.post.topic,
      publishImmediately: props.post.status === 1
    }
  }
}

// 监听props变化
watch(() => props.post, () => {
  initializeForm()
}, { immediate: true })

// 组件挂载时加载话题列表
onMounted(async () => {
  try {
    const response = await ForumApi.getTopics()
    // 将API返回的话题字符串转换为配置对象
    const apiTopics = response.data
    topics.value = topics.value.filter(topic => apiTopics.includes(topic.key))
  } catch (error) {
    console.error('加载话题列表失败:', error)
    // 使用默认话题列表
  }
})
</script>