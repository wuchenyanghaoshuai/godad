<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center items-center min-h-screen">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-pink-600"></div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="flex flex-col items-center justify-center min-h-screen">
      <div class="text-red-600 mb-4 text-lg">{{ error }}</div>
      <button
        @click="loadArticle"
        class="bg-pink-600 text-white px-6 py-2 rounded-lg hover:bg-pink-700 transition-colors"
      >
        重试
      </button>
    </div>

    <!-- 编辑表单 -->
    <div v-else-if="originalArticle" class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
      <!-- 页面标题 -->
      <div class="mb-6 sm:mb-8">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">编辑文章</h1>
            <p class="mt-1 sm:mt-2 text-sm sm:text-base text-gray-600">修改您的文章内容</p>
          </div>
          <button
            @click="goBack"
            class="flex items-center justify-center sm:justify-start text-gray-600 hover:text-gray-900 transition-colors bg-white/80 backdrop-blur-sm px-4 py-2 rounded-full shadow-sm hover:shadow-md border border-gray-200 hover:border-gray-300 text-sm sm:text-base"
          >
            <ArrowLeftIcon class="h-4 w-4 sm:h-5 sm:w-5 mr-2" />
            返回
          </button>
        </div>
      </div>

      <!-- 文章表单 -->
      <form @submit.prevent="submitArticle" class="space-y-6">
        <!-- 基本信息卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
          <h2 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">基本信息</h2>
          
          <!-- 文章标题 -->
          <div class="mb-4 sm:mb-6">
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
              class="w-full px-3 py-2.5 sm:py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent text-sm sm:text-base"
            />
            <div class="mt-1 text-sm text-gray-500 text-right">
              {{ form.title.length }}/100
            </div>
          </div>

          <!-- 文章摘要 -->
          <div class="mb-4 sm:mb-6">
            <label for="summary" class="block text-sm font-medium text-gray-700 mb-2">
              文章摘要
            </label>
            <textarea
              id="summary"
              v-model="form.summary"
              rows="3"
              maxlength="200"
              placeholder="请输入文章摘要（可选）"
              class="w-full px-3 py-2.5 sm:py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none text-sm sm:text-base"
            ></textarea>
            <div class="mt-1 text-sm text-gray-500 text-right">
              {{ form.summary.length }}/200
            </div>
          </div>

          <!-- 分类选择 -->
          <div class="mb-4 sm:mb-6">
            <label for="category" class="block text-sm font-medium text-gray-700 mb-2">
              文章分类 <span class="text-red-500">*</span>
            </label>
            <select
              id="category"
              v-model="form.category_id"
              required
              class="w-full px-3 py-2.5 sm:py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent text-sm sm:text-base"
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
          <div class="mb-4 sm:mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              文章标签
            </label>
            <div class="flex flex-wrap gap-1.5 sm:gap-2 mb-2">
              <span
                v-for="(tag, index) in form.tags"
                :key="index"
                class="inline-flex items-center bg-pink-100 text-pink-700 px-2.5 sm:px-3 py-1 rounded-full text-xs sm:text-sm"
              >
                {{ tag }}
                <button
                  type="button"
                  @click="removeTag(index)"
                  class="ml-1.5 sm:ml-2 text-pink-500 hover:text-pink-700"
                >
                  <XIcon class="h-3 w-3 sm:h-4 sm:w-4" />
                </button>
              </span>
            </div>
            <div class="flex flex-col sm:flex-row gap-2">
              <input
                v-model="newTag"
                type="text"
                placeholder="输入标签后按回车添加"
                @keydown.enter.prevent="addTag"
                class="flex-1 px-3 py-2.5 sm:py-2 border border-gray-300 rounded-lg sm:rounded-l-lg sm:rounded-r-none focus:ring-2 focus:ring-pink-500 focus:border-transparent text-sm sm:text-base"
              />
              <button
                type="button"
                @click="addTag"
                class="px-4 py-2.5 sm:py-2 bg-pink-600 text-white rounded-lg sm:rounded-l-none sm:rounded-r-lg hover:bg-pink-700 transition-colors text-sm sm:text-base font-medium"
              >
                添加
              </button>
            </div>
          </div>
        </div>

        <!-- 封面图片卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
          <h2 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">封面图片</h2>
          
          <!-- 图片上传组件 -->
          <ImageUpload
            v-model="form.cover_image"
            :multiple="false"
            :max-size="5"
            :article-title="originalArticle?.title || form.title || `article-${originalArticle?.id || 'editing'}`"
            @upload-success="handleCoverUpload"
            @upload-error="handleUploadError"
          />
        </div>

        <!-- 文章内容卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
          <h2 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">文章内容</h2>
          
          <!-- 富文本编辑器 -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              正文内容 <span class="text-red-500">*</span>
            </label>
            <RichTextEditor
              v-model="form.content"
              placeholder="请输入文章内容..."
              :min-height="300"
              class="sm:min-h-[400px]"
              :article-title="originalArticle?.title || form.title || `article-${originalArticle?.id || 'editing'}`"
              @change="handleContentChange"
            />
          </div>
        </div>

        <!-- 发布设置卡片 -->
        <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
          <h2 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">发布设置</h2>
          
          <!-- 发布状态 -->
          <div class="mb-4 sm:mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2 sm:mb-3">
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
                <span class="ml-2 text-sm sm:text-base text-gray-700">保存为草稿</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="form.status"
                  type="radio"
                  value="published"
                  class="text-pink-600 focus:ring-pink-500"
                />
                <span class="ml-2 text-sm sm:text-base text-gray-700">发布文章</span>
              </label>
            </div>
          </div>

          <!-- 精选功能已移除 -->

          <!-- 文章信息 -->
          <div v-if="originalArticle" class="bg-gray-50 rounded-lg p-3 sm:p-4 mb-3 sm:mb-4">
            <h3 class="text-sm font-medium text-gray-700 mb-2">文章信息</h3>
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-2 sm:gap-4 text-xs sm:text-sm">
              <div class="flex flex-col sm:block">
                <span class="text-gray-500">创建时间：</span>
                <span class="text-gray-900 font-medium">{{ formatDate(originalArticle.created_at) }}</span>
              </div>
              <div class="flex flex-col sm:block">
                <span class="text-gray-500">浏览量：</span>
                <span class="text-gray-900 font-medium">{{ originalArticle.view_count || 0 }}</span>
              </div>
              <div class="flex flex-col sm:block">
                <span class="text-gray-500">点赞数：</span>
                <span class="text-gray-900 font-medium">{{ originalArticle.like_count || 0 }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 提交按钮 -->
        <div class="flex flex-col sm:flex-row justify-end gap-3 sm:gap-4">
          <button
            type="button"
            @click="goBack"
            class="w-full sm:w-auto px-6 py-2.5 sm:py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors text-sm sm:text-base font-medium"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="!isFormValid || isSubmitting"
            class="w-full sm:w-auto px-6 py-2.5 sm:py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm sm:text-base font-medium"
          >
            {{ isSubmitting ? '保存中...' : '保存文章' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
import { useToast } from '@/composables/useToast'
import type { Article, Category, ArticleUpdateRequest, ArticleFormData, ImageUploadResponse } from '@/api/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { toast } = useToast()

// 响应式数据
const originalArticle = ref<Article | null>(null)
const categories = ref<Category[]>([])
const isLoading = ref(false)
const error = ref('')
const isSubmitting = ref(false)
const newTag = ref('')


// 表单数据
const form = ref<ArticleFormData>({
  title: '',
  content: '',
  summary: '',
  category_id: 0,
  cover_image: '',
  tags: [],
  status: 'draft'
})

// 计算属性
const isFormValid = computed(() => {
  return (form.value.title || '').trim() && 
         (form.value.content || '').trim() && 
         form.value.category_id > 0
})

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

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
  if (!originalArticle.value) return false
  
  return form.value.title !== originalArticle.value.title ||
         form.value.content !== originalArticle.value.content ||
         form.value.summary !== (originalArticle.value.summary || '') ||
         form.value.category_id !== originalArticle.value.category_id ||
         form.value.cover_image !== (originalArticle.value.cover_image || '') ||
         form.value.status !== originalArticle.value.status ||
         JSON.stringify(form.value.tags) !== JSON.stringify([])
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
  const url = Array.isArray(response) ? response[0]?.public_url : response.public_url
  if (url) {
    form.value.cover_image = url
  }
}

// 处理上传错误
const handleUploadError = (error: string) => {
  console.error('上传错误:', error)
  toast.error('上传失败：' + error)
}

// 生成 slug
const generateSlug = (title: string): string => {
  return title
    .toLowerCase()
    .replace(/[^\w\s\u4e00-\u9fa5-]/g, '') // 保留字母、数字、空格和中文
    .replace(/\s+/g, '-') // 空格替换为连字符
    .trim()
    .substring(0, 100) // 限制长度
    || `article-${Date.now()}` // fallback
}

// 处理内容变化
const handleContentChange = (content: string) => {
  form.value.content = content
}

// 提交文章
const submitArticle = async () => {
  if (!isFormValid.value || isSubmitting.value || !originalArticle.value) return
  
  try {
    isSubmitting.value = true
    
    // 准备提交数据，转换数据类型并包含所有必需字段
    const submitData: ArticleUpdateRequest = {
      title: form.value.title.trim(),
      slug: generateSlug(form.value.title.trim()),
      content: form.value.content.trim(),
      summary: form.value.summary?.trim() || '',
      cover_image: form.value.cover_image || '',
      category_id: form.value.category_id,
      tags: form.value.tags.join(','),
      status: form.value.status === 'published' ? 1 : 0
    }
    
    await ArticleApi.updateArticle(originalArticle.value.id, submitData)
    
    // 跳转到文章详情页
    router.push(`/articles/${originalArticle.value.id}`)
  } catch (err: any) {
    console.error('更新文章失败:', err)
    alert('更新失败：' + (err.message || '未知错误'))
  } finally {
    isSubmitting.value = false
  }
}

// 加载文章详情
const loadArticle = async () => {
  const articleId = Number(route.params.id)
  if (!articleId) {
    error.value = '无效的文章ID'
    return
  }
  
  try {
    isLoading.value = true
    error.value = ''
    
    const response = await ArticleApi.getArticleDetail(articleId)
    originalArticle.value = response.data
    
    // 检查权限
    if (!authStore.isAuthenticated || 
        (authStore.user?.id !== originalArticle.value.author_id && !authStore.isAdmin)) {
      error.value = '您没有权限编辑此文章'
      return
    }
    
    // 填充表单数据
    form.value = {
      title: originalArticle.value.title || '',
      content: originalArticle.value.content || '',
      summary: originalArticle.value.summary || '',
      category_id: originalArticle.value.category_id,
      cover_image: originalArticle.value.cover_image || '',
      tags: originalArticle.value.tags ? originalArticle.value.tags.split(',').map(tag => tag.trim()).filter(tag => tag) : [],
      status: originalArticle.value.status === 1 ? 'published' : 'draft',
      // is_featured 属性已移除
    }
    
  } catch (err: any) {
    // 根据错误状态码或错误信息提供更精确的错误提示
    if (err.response?.status === 404) {
      error.value = '文章不存在或已被删除'
    } else if (err.response?.status === 403) {
      error.value = '您没有权限访问此文章'
    } else if (err.response?.data?.message) {
      error.value = err.response.data.message
    } else {
      error.value = err.message || '加载文章失败'
    }
    console.error('加载文章失败:', err)
  } finally {
    isLoading.value = false
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
  loadArticle()
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