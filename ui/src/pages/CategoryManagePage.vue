<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- 页面标题和操作按钮 -->
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">分类管理</h1>
          <p class="mt-2 text-gray-600">管理文章分类，创建和编辑分类信息</p>
        </div>
        <button
          @click="showCreateModal = true"
          class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors duration-200 flex items-center gap-2"
        >
          <Plus class="w-5 h-5" />
          新建分类
        </button>
      </div>

      <!-- 搜索和筛选 -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <div class="relative">
              <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                v-model="searchKeyword"
                type="text"
                placeholder="搜索分类名称..."
                class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                @input="handleSearch"
              />
            </div>
          </div>
          <button
            @click="loadCategories(true)"
            class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-lg transition-colors duration-200 flex items-center gap-2"
          >
            <RefreshCw class="w-4 h-4" />
            刷新
          </button>
        </div>
      </div>

      <!-- 分类列表 -->
      <div class="bg-white rounded-lg shadow-sm overflow-hidden">
        <!-- 加载状态 -->
        <div v-if="loading" class="p-8 text-center">
          <div class="inline-flex items-center gap-2 text-gray-600">
            <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
            加载中...
          </div>
        </div>

        <!-- 错误状态 -->
        <div v-else-if="error" class="p-8 text-center">
          <div class="text-red-600 mb-4">
            <AlertCircle class="w-12 h-12 mx-auto mb-2" />
            <p class="font-medium">加载失败</p>
            <p class="text-sm text-gray-600 mt-1">{{ error }}</p>
          </div>
          <button
            @click="loadCategories(true)"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors duration-200"
          >
            重试
          </button>
        </div>

        <!-- 分类表格 -->
        <div v-else-if="categories.length > 0">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    分类名称
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    描述
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    创建时间
                  </th>
                  <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="category in categories" :key="category.id" class="hover:bg-gray-50">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="font-medium text-gray-900">{{ category.name }}</div>
                  </td>
                  <td class="px-6 py-4">
                    <div class="text-gray-600 max-w-xs truncate">
                      {{ category.description || '暂无描述' }}
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-gray-600">
                    {{ formatDate(category.created_at) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right">
                    <div class="flex justify-end gap-2">
                      <button
                        @click="editCategory(category)"
                        class="text-blue-600 hover:text-blue-800 p-2 rounded-lg hover:bg-blue-50 transition-colors duration-200"
                        title="编辑"
                      >
                        <Edit class="w-4 h-4" />
                      </button>
                      <button
                        @click="confirmDelete(category)"
                        class="text-red-600 hover:text-red-800 p-2 rounded-lg hover:bg-red-50 transition-colors duration-200"
                        title="删除"
                      >
                        <Trash2 class="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 分页 -->
          <div v-if="totalPages > 1" class="px-6 py-4 border-t border-gray-200">
            <div class="flex items-center justify-between">
              <div class="text-sm text-gray-600">
                共 {{ totalCategories }} 个分类
              </div>
              <div class="flex gap-2">
                <button
                  @click="changePage(currentPage - 1)"
                  :disabled="currentPage <= 1"
                  class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  上一页
                </button>
                <span class="px-3 py-1 text-sm text-gray-600">
                  {{ currentPage }} / {{ totalPages }}
                </span>
                <button
                  @click="changePage(currentPage + 1)"
                  :disabled="currentPage >= totalPages"
                  class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  下一页
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="p-12 text-center">
          <Folder class="w-16 h-16 mx-auto text-gray-300 mb-4" />
          <h3 class="text-lg font-medium text-gray-900 mb-2">暂无分类</h3>
          <p class="text-gray-600 mb-6">还没有创建任何分类，点击上方按钮创建第一个分类</p>
          <button
            @click="showCreateModal = true"
            class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors duration-200"
          >
            创建分类
          </button>
        </div>
      </div>
    </div>

    <!-- 创建/编辑分类模态框 -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">
            {{ showCreateModal ? '创建分类' : '编辑分类' }}
          </h3>
        </div>
        
        <form @submit.prevent="submitForm" class="p-6">
          <div class="mb-4">
            <label for="categoryName" class="block text-sm font-medium text-gray-700 mb-2">
              分类名称 *
            </label>
            <input
              id="categoryName"
              v-model="formData.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="请输入分类名称"
            />
          </div>
          
          <div class="mb-6">
            <label for="categoryDescription" class="block text-sm font-medium text-gray-700 mb-2">
              分类描述
            </label>
            <textarea
              id="categoryDescription"
              v-model="formData.description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="请输入分类描述（可选）"
            ></textarea>
          </div>
          
          <div class="flex justify-end gap-3">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors duration-200"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <div v-if="submitting" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              {{ submitting ? '保存中...' : (showCreateModal ? '创建' : '保存') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="flex-shrink-0">
              <AlertTriangle class="w-6 h-6 text-red-600" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">确认删除</h3>
              <p class="text-gray-600 mt-1">
                确定要删除分类 "{{ categoryToDelete?.name }}" 吗？此操作不可撤销。
              </p>
            </div>
          </div>
          
          <div class="flex justify-end gap-3">
            <button
              @click="showDeleteModal = false"
              class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors duration-200"
            >
              取消
            </button>
            <button
              @click="deleteCategory"
              :disabled="deleting"
              class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <div v-if="deleting" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              {{ deleting ? '删除中...' : '确认删除' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'
import {
  Plus,
  Search,
  RefreshCw,
  Edit,
  Trash2,
  AlertCircle,
  AlertTriangle,
  Folder
} from 'lucide-vue-next'
import { CategoryApi } from '@/api/category'
import { useAuthStore } from '@/stores/auth'
import type { Category, CategoryRequest } from '@/api/types'

const router = useRouter()
const authStore = useAuthStore()
const { toast } = useToast()

// 检查管理员权限
if (!authStore.isAdmin) {
  router.push('/')
}

// 响应式数据
const categories = ref<Category[]>([])
const loading = ref(false)
const error = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const totalCategories = ref(0)

// 模态框状态
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const submitting = ref(false)
const deleting = ref(false)

// 表单数据
const formData = ref<CategoryRequest>({
  name: '',
  description: ''
})

// 编辑和删除的分类
const editingCategory = ref<Category | null>(null)
const categoryToDelete = ref<Category | null>(null)

// 计算属性
const totalPages = computed(() => Math.ceil(totalCategories.value / pageSize.value))

// 加载分类列表
const loadCategories = async (reset = false) => {
  if (reset) {
    currentPage.value = 1
  }
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await CategoryApi.getAdminCategoryPage({
      page: currentPage.value,
      size: pageSize.value,
      keyword: searchKeyword.value || undefined
    })
    const page = response.data
    categories.value = page.items
    totalCategories.value = page.total
    currentPage.value = page.page
    pageSize.value = page.size
  } catch (err: any) {
    error.value = err.message || '获取分类列表失败'
    console.error('获取分类列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 搜索处理
let searchTimer: number
const handleSearch = () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    loadCategories(true)
  }, 500)
}

// 分页处理
const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadCategories()
  }
}

// 编辑分类
const editCategory = (category: Category) => {
  editingCategory.value = category
  formData.value = {
    name: category.name,
    description: category.description || ''
  }
  showEditModal.value = true
}

// 确认删除
const confirmDelete = (category: Category) => {
  categoryToDelete.value = category
  showDeleteModal.value = true
}

// 关闭模态框
const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingCategory.value = null
  formData.value = {
    name: '',
    description: ''
  }
}

// 提交表单
const submitForm = async () => {
  if (!formData.value.name.trim()) {
    toast.error('请输入分类名称')
    return
  }
  
  submitting.value = true
  
  try {
    if (showCreateModal.value) {
      // 创建分类
      const response = await CategoryApi.createCategory(formData.value)
      if (response.data) {
        toast.success('分类创建成功')
        closeModal()
        loadCategories(true)
      } else {
        throw new Error(response.message || '创建分类失败')
      }
    } else if (showEditModal.value && editingCategory.value) {
      // 更新分类
      const response = await CategoryApi.updateCategory(editingCategory.value.id, formData.value)
      if (response.data) {
        toast.success('分类更新成功')
        closeModal()
        loadCategories()
      } else {
        throw new Error(response.message || '更新分类失败')
      }
    }
  } catch (err: any) {
    toast.error(err.message || '操作失败')
    console.error('分类操作失败:', err)
  } finally {
    submitting.value = false
  }
}

// 删除分类
const deleteCategory = async () => {
  if (!categoryToDelete.value) return
  
  deleting.value = true
  
  try {
    const response = await CategoryApi.deleteCategory(categoryToDelete.value.id)
    if (response.data !== undefined) {
      toast.success('分类删除成功')
      showDeleteModal.value = false
      categoryToDelete.value = null
      loadCategories()
    } else {
      throw new Error(response.message || '删除分类失败')
    }
  } catch (err: any) {
    toast.error(err.message || '删除分类失败')
    console.error('删除分类失败:', err)
  } finally {
    deleting.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 组件挂载时加载数据
onMounted(() => {
  loadCategories()
})
</script>
