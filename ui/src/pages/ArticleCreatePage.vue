<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <BaseHeader :showCreateButton="true" />
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 页面标题 -->
      <div class="mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">发布文章</h1>
            <p class="mt-2 text-gray-600">分享您的育儿经验和知识</p>
          </div>
          <button
            @click="goBack"
            class="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
          >
            <ArrowLeftIcon class="h-5 w-5 mr-2" />
            返回
          </button>
        </div>
      </div>

      <!-- 文章表单 -->
      <form @submit.prevent="submitArticle" class="space-y-6">
        <!-- 基本信息卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">基本信息</h2>
          
          <!-- 文章标题 -->
          <div class="mb-6">
            <label for="title" class="block text-sm font-medium text-gray-700 mb-2">
              文章标题 <span class="text-red-500">*</span>
            </label>
            <input
              id="title"
              v-model="form.title"
              type="text"
              required
              maxlength="100"
              placeholder="请输入文章标题"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
            />
            <div class="mt-1 text-sm text-gray-500 text-right">
              {{ form.title.length }}/100
            </div>
          </div>

          <!-- 文章摘要 -->
          <div class="mb-6">
            <label for="summary" class="block text-sm font-medium text-gray-700 mb-2">
              文章摘要
            </label>
            <textarea
              id="summary"
              v-model="form.summary"
              rows="3"
              maxlength="200"
              placeholder="请输入文章摘要（可选）"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none"
            ></textarea>
            <div class="mt-1 text-sm text-gray-500 text-right">
              {{ form.summary.length }}/200
            </div>
          </div>

          <!-- 分类选择 -->
          <div class="mb-6">
            <label for="category" class="block text-sm font-medium text-gray-700 mb-2">
              文章分类 <span class="text-red-500">*</span>
            </label>
            <select
              id="category"
              v-model="form.category_id"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
            >
              <option value="">请选择分类</option>
              <option
                v-for="category in categories"
                :key="category.id"
                :value="category.id"
              >
                {{ category.name }}
              </option>
            </select>
          </div>

          <!-- 文章标签 -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              文章标签
            </label>
            <div class="flex flex-wrap gap-2 mb-2">
              <span
                v-for="(tag, index) in form.tags"
                :key="index"
                class="inline-flex items-center bg-pink-100 text-pink-700 px-3 py-1 rounded-full text-sm"
              >
                {{ tag }}
                <button
                  type="button"
                  @click="removeTag(index)"
                  class="ml-2 text-pink-500 hover:text-pink-700"
                >
                  <XIcon class="h-4 w-4" />
                </button>
              </span>
            </div>
            <div class="flex">
              <input
                v-model="newTag"
                type="text"
                placeholder="输入标签后按回车添加"
                @keydown.enter.prevent="addTag"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent"
              />
              <button
                type="button"
                @click="addTag"
                class="px-4 py-2 bg-pink-600 text-white rounded-r-lg hover:bg-pink-700 transition-colors"
              >
                添加
              </button>
            </div>
          </div>
        </div>

        <!-- 封面图片卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">封面图片</h2>
          
          <!-- 图片上传组件 -->
          <ImageUpload
            v-model="form.cover_image"
            :multiple="false"
            :max-size="5"
            :article-title="form.title || 'new-article'"
            @upload-success="handleCoverUpload"
            @upload-error="handleUploadError"
          />
        </div>

        <!-- 文章内容卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">文章内容</h2>
          
          <!-- 富文本编辑器 -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              正文内容 <span class="text-red-500">*</span>
            </label>
            <RichTextEditor
              v-model="form.content"
              placeholder="请输入文章内容..."
              :min-height="400"
              :article-title="form.title || 'new-article'"
              @change="handleContentChange"
            />
          </div>
        </div>

        <!-- 发布设置卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">发布设置</h2>
          
          <!-- 发布状态 -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-3">
              发布状态
            </label>
            <div class="space-y-2">
              <label class="flex items-center">
                <input
                  v-model="form.status"
                  type="radio"
                  value="draft"
                  class="text-pink-600 focus:ring-pink-500"
                />
                <span class="ml-2 text-gray-700">保存为草稿</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="form.status"
                  type="radio"
                  value="published"
                  class="text-pink-600 focus:ring-pink-500"
                />
                <span class="ml-2 text-gray-700">立即发布</span>
              </label>
            </div>
          </div>

          <!-- 是否精选 -->
          <div v-if="authStore.isAdmin" class="mb-6">
            <label class="flex items-center">
              <input
                v-model="form.is_featured"
                type="checkbox"
                class="text-pink-600 focus:ring-pink-500 rounded"
              />
              <span class="ml-2 text-gray-700">设为精选文章</span>
            </label>
          </div>
        </div>

        <!-- 提交按钮 -->
        <div class="flex justify-end space-x-4">
          <button
            type="button"
            @click="goBack"
            class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            :class="[
              'px-6 py-2 rounded-lg transition-colors flex items-center',
              isFormValid && !isSubmitting
                ? 'bg-pink-600 text-white hover:bg-pink-700'
                : 'bg-gray-300 text-gray-500 cursor-not-allowed'
            ]"
          >
            <LoaderIcon v-if="isSubmitting" class="h-4 w-4 mr-2 animate-spin" />
            {{ form.status === 'published' ? '发布文章' : '保存草稿' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeftIcon,
  ImageIcon,
  XIcon,
  LoaderIcon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { ArticleApi } from '@/api/article'
import { CategoryApi } from '@/api/category'
import { UploadApi, UploadUtils } from '@/api/upload'
import RichTextEditor from '@/components/RichTextEditor.vue'
import ImageUpload from '@/components/ImageUpload.vue'
import BaseHeader from '@/components/BaseHeader.vue'
import { useToast } from '@/composables/useToast'
import type { Category, ArticleCreateRequest, ImageUploadResponse } from '@/api/types'

const router = useRouter()
const authStore = useAuthStore()
const { toast } = useToast()

// 响应式数据
const categories = ref<Category[]>([])
const isSubmitting = ref(false)
const newTag = ref('')

// 前端表单数据类型（用户友好格式）
interface FormData {
  title: string
  content: string
  summary: string
  category_id: number
  cover_image: string
  tags: string[] // 前端使用数组格式
  status: 'draft' | 'published' // 前端使用字符串格式
  is_featured: boolean // 前端使用is_featured字段名
}

// 表单数据
const form = ref<FormData>({
  title: '',
  content: '',
  summary: '',
  category_id: 0,
  cover_image: '',
  tags: [],
  status: 'draft',
  is_featured: false
})

// 计算属性
const isFormValid = computed(() => {
  return (form.value.title || '').trim() && 
         (form.value.content || '').trim() && 
         form.value.category_id > 0
})

// 返回上一页
const goBack = () => {
  if (hasUnsavedChanges()) {
    if (confirm('您有未保存的更改，确定要离开吗？')) {
      router.back()
    }
  } else {
    router.back()
  }
}

// 检查是否有未保存的更改
const hasUnsavedChanges = () => {
  return (form.value.title || '').trim() || 
         (form.value.content || '').trim() || 
         (form.value.summary || '').trim() ||
         form.value.cover_image ||
         form.value.tags.length > 0
}

// 添加标签
const addTag = () => {
  const tag = newTag.value.trim()
  if (tag && !form.value.tags.includes(tag) && form.value.tags.length < 10) {
    form.value.tags.push(tag)
    newTag.value = ''
  }
}

// 移除标签
const removeTag = (index: number) => {
  form.value.tags.splice(index, 1)
}

// 处理封面图片上传成功
const handleCoverUpload = (response: ImageUploadResponse | ImageUploadResponse[]) => {
  // v-model已经自动更新了form.cover_image，不需要显示成功消息
}

// 处理上传错误
const handleUploadError = (error: string) => {
  toast.error(error)
}

// 处理内容变化
const handleContentChange = (content: string) => {
  // Content change handling can be added here if needed
}

// 提交文章
const submitArticle = async () => {
  if (!isFormValid.value || isSubmitting.value) return
  
  try {
    isSubmitting.value = true
    
    // 转换数据格式以匹配后端API要求
    const submitData: ArticleCreateRequest = {
      title: form.value.title,
      slug: form.value.title.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/^-+|-+$/g, ''), // 生成slug
      content: form.value.content,
      summary: form.value.summary,
      category_id: form.value.category_id,
      cover_image: form.value.cover_image,
      tags: form.value.tags.join(','), // 将数组转换为逗号分隔的字符串
      status: form.value.status === 'published' ? 1 : 0, // 将字符串转换为数字
      is_recommend: form.value.is_featured // 字段名映射
    }
    
    const response = await ArticleApi.createArticle(submitData)
    
    if (response.data) {
      // 跳转到文章详情页
      router.push(`/articles/${response.data.id}`)
    } else {
      throw new Error('服务器返回数据异常')
    }
  } catch (err: any) {
    console.error('发布文章失败:', err)
    const errorMessage = err.response?.data?.message || err.message || '发布失败，请重试'
    toast.error(errorMessage)
  } finally {
    isSubmitting.value = false
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await CategoryApi.getCategoryList()
    categories.value = (response.data && Array.isArray(response.data)) ? response.data : []
  } catch (err) {
    console.error('加载分类失败:', err)
  }
}

// 组件挂载时加载数据
onMounted(() => {
  // 检查用户是否已登录
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  
  loadCategories()
})

// 页面离开前确认
window.addEventListener('beforeunload', (event) => {
  if (hasUnsavedChanges()) {
    event.preventDefault()
    event.returnValue = ''
  }
})
</script>

<style scoped>
/* 自定义样式 */
.font-mono {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}
</style>