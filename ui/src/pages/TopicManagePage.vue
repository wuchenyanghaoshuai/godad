<template>
  <div class="min-h-screen bg-gray-50">
    <BaseHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 页面标题和操作 -->
      <div class="flex justify-between items-center mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">话题管理</h1>
          <p class="text-gray-600 mt-1">管理社区论坛的话题分类</p>
        </div>
        <button
          @click="showCreateModal = true"
          class="bg-pink-600 text-white px-4 py-2 rounded-lg hover:bg-pink-700 transition-colors flex items-center space-x-2"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>新建话题</span>
        </button>
      </div>

      <!-- 话题列表 -->
      <div class="bg-white rounded-lg shadow">
        <!-- 表头 -->
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">话题列表</h3>
            <div class="flex items-center space-x-2">
              <span class="text-sm text-gray-500">共 {{ totalTopics }} 个话题</span>
            </div>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="isLoading" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600"></div>
        </div>

        <!-- 话题表格 -->
        <div v-else-if="topics.length > 0" class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">话题</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">显示名称</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">帖子数量</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">排序</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="topic in topics" :key="topic.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-6 w-6 rounded-full flex items-center justify-center mr-3" :style="{ backgroundColor: topic.color }">
                      <i :class="`fas fa-${topic.icon || 'tag'}`" class="text-white text-xs"></i>
                    </div>
                    <div>
                      <div class="text-sm font-medium text-gray-900">{{ topic.name }}</div>
                      <div class="text-sm text-gray-500" v-if="topic.description">{{ topic.description }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ topic.display_name }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ topic.post_count }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ topic.sort }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="[
                    'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                    topic.is_active
                      ? 'bg-green-100 text-green-800'
                      : 'bg-red-100 text-red-800'
                  ]">
                    {{ topic.is_active ? '启用' : '禁用' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(topic.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex items-center space-x-2">
                    <button
                      @click="editTopic(topic)"
                      class="text-indigo-600 hover:text-indigo-900 transition-colors"
                    >
                      编辑
                    </button>
                    <button
                      @click="confirmDelete(topic)"
                      class="text-red-600 hover:text-red-900 transition-colors"
                      :disabled="topic.post_count > 0"
                      :class="{ 'opacity-50 cursor-not-allowed': topic.post_count > 0 }"
                    >
                      删除
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-12">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无话题</h3>
          <p class="mt-1 text-sm text-gray-500">开始创建第一个话题吧</p>
        </div>

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
          <div class="flex items-center justify-between">
            <div class="flex-1 flex justify-between sm:hidden">
              <button
                @click="changePage(currentPage - 1)"
                :disabled="currentPage <= 1"
                class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                上一页
              </button>
              <button
                @click="changePage(currentPage + 1)"
                :disabled="currentPage >= totalPages"
                class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                下一页
              </button>
            </div>
            <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
              <div>
                <p class="text-sm text-gray-700">
                  显示第 <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span> 到
                  <span class="font-medium">{{ Math.min(currentPage * pageSize, totalTopics) }}</span> 项，
                  共 <span class="font-medium">{{ totalTopics }}</span> 项
                </p>
              </div>
              <div>
                <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                  <button
                    @click="changePage(currentPage - 1)"
                    :disabled="currentPage <= 1"
                    class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">上一页</span>
                    <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </button>

                  <button
                    v-for="page in visiblePages"
                    :key="page"
                    @click="changePage(page)"
                    :class="[
                      'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                      page === currentPage
                        ? 'z-10 bg-pink-50 border-pink-500 text-pink-600'
                        : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                    ]"
                  >
                    {{ page }}
                  </button>

                  <button
                    @click="changePage(currentPage + 1)"
                    :disabled="currentPage >= totalPages"
                    class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">下一页</span>
                    <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                    </svg>
                  </button>
                </nav>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑话题模态框 -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeModals">
      <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">
            {{ showCreateModal ? '新建话题' : '编辑话题' }}
          </h3>

          <form @submit.prevent="submitForm">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- 话题名称 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">话题名称 *</label>
                <input
                  v-model="formData.name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                  placeholder="例如: Baby Care"
                />
              </div>

              <!-- 显示名称 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">显示名称 *</label>
                <input
                  v-model="formData.display_name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                  placeholder="例如: 婴儿护理"
                />
              </div>

              <!-- 颜色 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">主题颜色</label>
                <div class="flex items-center space-x-2">
                  <input
                    v-model="formData.color"
                    type="color"
                    class="h-10 w-16 border border-gray-300 rounded cursor-pointer"
                  />
                  <input
                    v-model="formData.color"
                    type="text"
                    class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                    placeholder="#6366f1"
                  />
                </div>
              </div>

              <!-- 图标 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">图标</label>
                <input
                  v-model="formData.icon"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                  placeholder="例如: baby, heart, book"
                />
              </div>

              <!-- 排序 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                <input
                  v-model.number="formData.sort"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                  placeholder="0"
                />
              </div>

              <!-- 状态 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                <select
                  v-model="formData.is_active"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                >
                  <option :value="true">启用</option>
                  <option :value="false">禁用</option>
                </select>
              </div>
            </div>

            <!-- 描述 -->
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
              <textarea
                v-model="formData.description"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none"
                placeholder="话题描述..."
              ></textarea>
            </div>

            <!-- 按钮 -->
            <div class="flex justify-end space-x-3 mt-6">
              <button
                type="button"
                @click="closeModals"
                class="px-4 py-2 text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
              >
                取消
              </button>
              <button
                type="submit"
                :disabled="isSubmitting"
                class="px-4 py-2 bg-pink-600 text-white rounded-md hover:bg-pink-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ isSubmitting ? '提交中...' : (showCreateModal ? '创建' : '更新') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="showDeleteModal = false">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3 text-center">
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
            <svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-gray-900 mt-4">确认删除</h3>
          <div class="mt-2 px-7 py-3">
            <p class="text-sm text-gray-500">
              确定要删除话题 "<strong>{{ topicToDelete?.display_name }}</strong>" 吗？
              <br />此操作不可撤销。
            </p>
          </div>
          <div class="flex justify-center space-x-3 mt-4">
            <button
              @click="showDeleteModal = false"
              class="px-4 py-2 text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
            >
              取消
            </button>
            <button
              @click="deleteTopic"
              :disabled="isDeleting"
              class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isDeleting ? '删除中...' : '确认删除' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useToast } from '@/composables/useToast'
import BaseHeader from '@/components/BaseHeader.vue'

// 响应式数据
const topics = ref<any[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const isDeleting = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalTopics = ref(0)
const totalPages = ref(0)
const editingTopic = ref<any>(null)
const topicToDelete = ref<any>(null)

// 表单数据
const formData = ref({
  name: '',
  display_name: '',
  description: '',
  color: '#6366f1',
  icon: '',
  sort: 0,
  is_active: true
})

const { showToast } = useToast()

// 计算属性
const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// API 调用函数
const fetchTopics = async (page = 1) => {
  isLoading.value = true
  try {
    const response = await fetch(`http://127.0.0.1:8888/api/admin/topics?page=${page}&size=${pageSize.value}&all=true`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    if (result.code === 200) {
      topics.value = result.data.items || []
      totalTopics.value = result.data.total || 0
      totalPages.value = result.data.total_page || 0
      currentPage.value = page
    } else {
      throw new Error(result.message || '获取话题列表失败')
    }
  } catch (error: any) {
    console.error('获取话题列表失败:', error)
    showToast('获取话题列表失败', 'error')
  } finally {
    isLoading.value = false
  }
}

const createTopic = async () => {
  try {
    const response = await fetch('http://127.0.0.1:8888/api/admin/topics', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData.value)
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    if (result.code === 201) {
      showToast('话题创建成功', 'success')
      closeModals()
      fetchTopics(currentPage.value)
    } else {
      throw new Error(result.message || '创建话题失败')
    }
  } catch (error: any) {
    console.error('创建话题失败:', error)
    showToast(error.message || '创建话题失败', 'error')
  }
}

const updateTopic = async () => {
  try {
    const response = await fetch(`http://127.0.0.1:8888/api/admin/topics/${editingTopic.value.id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData.value)
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    if (result.code === 200) {
      showToast('话题更新成功', 'success')
      closeModals()
      fetchTopics(currentPage.value)
    } else {
      throw new Error(result.message || '更新话题失败')
    }
  } catch (error: any) {
    console.error('更新话题失败:', error)
    showToast(error.message || '更新话题失败', 'error')
  }
}

const deleteTopicAPI = async (id: number) => {
  try {
    const response = await fetch(`http://127.0.0.1:8888/api/admin/topics/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    if (result.code === 200) {
      showToast('话题删除成功', 'success')
      showDeleteModal.value = false
      fetchTopics(currentPage.value)
    } else {
      throw new Error(result.message || '删除话题失败')
    }
  } catch (error: any) {
    console.error('删除话题失败:', error)
    showToast(error.message || '删除话题失败', 'error')
  }
}

// 事件处理函数
const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    fetchTopics(page)
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingTopic.value = null
  formData.value = {
    name: '',
    display_name: '',
    description: '',
    color: '#6366f1',
    icon: '',
    sort: 0,
    is_active: true
  }
}

const editTopic = (topic: any) => {
  editingTopic.value = topic
  formData.value = {
    name: topic.name,
    display_name: topic.display_name,
    description: topic.description || '',
    color: topic.color || '#6366f1',
    icon: topic.icon || '',
    sort: topic.sort || 0,
    is_active: topic.is_active
  }
  showEditModal.value = true
}

const confirmDelete = (topic: any) => {
  topicToDelete.value = topic
  showDeleteModal.value = true
}

const deleteTopic = async () => {
  if (!topicToDelete.value) return

  isDeleting.value = true
  try {
    await deleteTopicAPI(topicToDelete.value.id)
  } finally {
    isDeleting.value = false
    topicToDelete.value = null
  }
}

const submitForm = async () => {
  isSubmitting.value = true
  try {
    if (showCreateModal.value) {
      await createTopic()
    } else if (showEditModal.value) {
      await updateTopic()
    }
  } finally {
    isSubmitting.value = false
  }
}

// 生命周期
onMounted(() => {
  fetchTopics()
})
</script>