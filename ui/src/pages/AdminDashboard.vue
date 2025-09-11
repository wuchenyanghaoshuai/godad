<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-900">管理员后台</h1>
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
      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-blue-100 text-blue-600">
              <ArticleIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">文章总数</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.articleCount }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-green-100 text-green-600">
              <UserIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">用户总数</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.userCount }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-yellow-100 text-yellow-600">
              <CategoryIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">分类总数</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.categoryCount }}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-purple-100 text-purple-600">
              <CommentIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">评论总数</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.commentCount }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 标签页 -->
      <div class="bg-white rounded-lg shadow">
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                'whitespace-nowrap py-4 px-6 border-b-2 font-medium text-sm',
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              {{ tab.name }}
            </button>
          </nav>
        </div>

        <div class="p-6">
          <!-- 文章管理 -->
          <div v-if="activeTab === 'articles'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">文章管理</h2>
              <div class="flex space-x-2">
                <input
                  v-model="articleFilter.keyword"
                  @input="searchArticles"
                  type="text"
                  placeholder="搜索文章..."
                  class="px-3 py-2 border border-gray-300 rounded-md text-sm"
                >
                <select
                  v-model="articleFilter.status"
                  @change="loadArticles"
                  class="px-3 py-2 border border-gray-300 rounded-md text-sm"
                >
                  <option value="">全部状态</option>
                  <option value="0">草稿</option>
                  <option value="1">已发布</option>
                  <option value="2">已下架</option>
                </select>
              </div>
            </div>

            <div class="overflow-hidden">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      文章标题
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      作者
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      分类
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      状态
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      浏览量
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      创建时间
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      操作
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="article in articles" :key="article.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm font-medium text-gray-900 max-w-xs truncate">
                        {{ article.title }}
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ article.author?.nickname }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ article.category?.name }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          article.status === 1
                            ? 'bg-green-100 text-green-800'
                            : article.status === 0
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                        ]"
                      >
                        {{ getArticleStatusText(article.status) }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ article.view_count }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(article.created_at) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        @click="toggleArticleStatus(article)"
                        :class="[
                          'mr-2 inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md',
                          article.status === 1
                            ? 'text-red-700 bg-red-100 hover:bg-red-200'
                            : 'text-green-700 bg-green-100 hover:bg-green-200'
                        ]"
                      >
                        {{ article.status === 1 ? '下架' : '发布' }}
                      </button>
                      <button
                        @click="deleteArticle(article.id)"
                        class="text-red-600 hover:text-red-900"
                      >
                        删除
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- 分页 -->
            <div class="flex justify-between items-center mt-4">
              <span class="text-sm text-gray-700">
                显示 {{ articlePagination.offset + 1 }} 到 {{ Math.min(articlePagination.offset + articlePagination.limit, articlePagination.total) }} 条，共 {{ articlePagination.total }} 条
              </span>
              <div class="flex space-x-2">
                <button
                  @click="prevPage('article')"
                  :disabled="articlePagination.page === 1"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  上一页
                </button>
                <button
                  @click="nextPage('article')"
                  :disabled="articlePagination.page >= Math.ceil(articlePagination.total / articlePagination.limit)"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  下一页
                </button>
              </div>
            </div>
          </div>

          <!-- 用户管理 -->
          <div v-if="activeTab === 'users'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">用户管理</h2>
              <input
                v-model="userFilter.keyword"
                @input="searchUsers"
                type="text"
                placeholder="搜索用户..."
                class="px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
            </div>

            <div class="overflow-hidden">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      用户信息
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      邮箱
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      角色
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      状态
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      注册时间
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      操作
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="user in users" :key="user.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <img
                          class="h-10 w-10 rounded-full"
                          :src="user.avatar || '/default-avatar.png'"
                          :alt="user.nickname"
                        >
                        <div class="ml-4">
                          <div class="text-sm font-medium text-gray-900">
                            {{ user.nickname }}
                          </div>
                          <div class="text-sm text-gray-500">
                            @{{ user.username }}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ user.email }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          user.role === 2
                            ? 'bg-purple-100 text-purple-800'
                            : 'bg-gray-100 text-gray-800'
                        ]"
                      >
                        {{ user.role === 2 ? '管理员' : '普通用户' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          user.status === 1
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-800'
                        ]"
                      >
                        {{ user.status === 1 ? '正常' : '禁用' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(user.created_at) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        @click="toggleUserStatus(user)"
                        :class="[
                          'mr-2 inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md',
                          user.status === 1
                            ? 'text-red-700 bg-red-100 hover:bg-red-200'
                            : 'text-green-700 bg-green-100 hover:bg-green-200'
                        ]"
                      >
                        {{ user.status === 1 ? '禁用' : '启用' }}
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- 分页 -->
            <div class="flex justify-between items-center mt-4">
              <span class="text-sm text-gray-700">
                显示 {{ userPagination.offset + 1 }} 到 {{ Math.min(userPagination.offset + userPagination.limit, userPagination.total) }} 条，共 {{ userPagination.total }} 条
              </span>
              <div class="flex space-x-2">
                <button
                  @click="prevPage('user')"
                  :disabled="userPagination.page === 1"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  上一页
                </button>
                <button
                  @click="nextPage('user')"
                  :disabled="userPagination.page >= Math.ceil(userPagination.total / userPagination.limit)"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  下一页
                </button>
              </div>
            </div>
          </div>

          <!-- 分类管理 -->
          <div v-if="activeTab === 'categories'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">分类管理</h2>
              <button
                @click="showCreateCategoryModal = true"
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
              >
                新增分类
              </button>
            </div>

            <div class="overflow-hidden">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      分类名称
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      别名
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      描述
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      排序
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      状态
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      操作
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="category in categories" :key="category.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div
                          v-if="category.color"
                          class="w-4 h-4 rounded-full mr-3"
                          :style="{ backgroundColor: category.color }"
                        ></div>
                        <div class="text-sm font-medium text-gray-900">
                          {{ category.name }}
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ category.slug }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 max-w-xs truncate">
                      {{ category.description }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ category.sort }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          category.status === 1
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-800'
                        ]"
                      >
                        {{ category.status === 1 ? '启用' : '禁用' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        @click="editCategory(category)"
                        class="text-blue-600 hover:text-blue-900 mr-2"
                      >
                        编辑
                      </button>
                      <button
                        @click="toggleCategoryStatus(category)"
                        :class="[
                          'mr-2',
                          category.status === 1
                            ? 'text-red-600 hover:text-red-900'
                            : 'text-green-600 hover:text-green-900'
                        ]"
                      >
                        {{ category.status === 1 ? '禁用' : '启用' }}
                      </button>
                      <button
                        @click="deleteCategory(category.id)"
                        class="text-red-600 hover:text-red-900"
                      >
                        删除
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑分类模态框 -->
    <div v-if="showCreateCategoryModal || showEditCategoryModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">
          {{ showCreateCategoryModal ? '新增分类' : '编辑分类' }}
        </h3>
        <form @submit.prevent="submitCategory">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">分类名称</label>
            <input
              v-model="categoryForm.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">别名</label>
            <input
              v-model="categoryForm.slug"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="留空将自动生成"
            >
            <p class="text-xs text-gray-500 mt-1">用于URL的英文标识，留空时将根据分类名称自动生成</p>
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">描述</label>
            <textarea
              v-model="categoryForm.description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            ></textarea>
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">颜色</label>
            <input
              v-model="categoryForm.color"
              type="color"
              class="w-16 h-10 border border-gray-300 rounded-md"
            >
          </div>
          <div class="flex justify-end space-x-2">
            <button
              type="button"
              @click="closeCategoryModal"
              class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              取消
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              {{ showCreateCategoryModal ? '创建' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { http } from '../api/http'

// Icons
import {
  FileTextIcon as ArticleIcon,
  UserIcon,
  TagIcon as CategoryIcon,
  MessageCircleIcon as CommentIcon
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const activeTab = ref('articles')
const stats = reactive({
  articleCount: 0,
  userCount: 0,
  categoryCount: 0,
  commentCount: 0
})

const tabs = [
  { id: 'articles', name: '文章管理' },
  { id: 'users', name: '用户管理' },
  { id: 'categories', name: '分类管理' }
]

// 文章管理
const articles = ref([])
const articleFilter = reactive({
  keyword: '',
  status: ''
})
const articlePagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  offset: 0
})

// 用户管理
const users = ref([])
const userFilter = reactive({
  keyword: ''
})
const userPagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  offset: 0
})

// 分类管理
const categories = ref([])
const showCreateCategoryModal = ref(false)
const showEditCategoryModal = ref(false)
const categoryForm = reactive({
  id: null,
  name: '',
  slug: '',
  description: '',
  color: '#3B82F6'
})

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await http.get('/admin/stats')
    stats.articleCount = response.data.articleCount
    stats.userCount = response.data.userCount
    stats.categoryCount = response.data.categoryCount
    stats.commentCount = response.data.commentCount
  } catch (error) {
    console.error('加载统计数据失败:', error)
    // 如果API调用失败，使用默认值
    stats.articleCount = 0
    stats.userCount = 0
    stats.categoryCount = 0
    stats.commentCount = 0
  }
}

// 加载文章列表
const loadArticles = async () => {
  try {
    const params = {
      page: articlePagination.page,
      size: articlePagination.limit,
      status: articleFilter.status,
      keyword: articleFilter.keyword
    }
    const response = await http.get('/admin/articles', params)
    articles.value = response.data.articles || []
    articlePagination.total = response.data.total || 0
    articlePagination.offset = (articlePagination.page - 1) * articlePagination.limit
  } catch (error) {
    console.error('加载文章列表失败:', error)
  }
}

// 搜索文章
let articleSearchTimeout = null
const searchArticles = () => {
  clearTimeout(articleSearchTimeout)
  articleSearchTimeout = setTimeout(() => {
    articlePagination.page = 1
    loadArticles()
  }, 500)
}

// 切换文章状态
const toggleArticleStatus = async (article) => {
  try {
    const newStatus = article.status === 1 ? 2 : 1
    await http.put(`/admin/articles/${article.id}/status`, { status: newStatus })
    article.status = newStatus
  } catch (error) {
    console.error('切换文章状态失败:', error)
    alert('操作失败，请重试')
  }
}

// 删除文章
const deleteArticle = async (articleId) => {
  if (!confirm('确定要删除这篇文章吗？')) return
  try {
    await http.delete(`/admin/articles/${articleId}`)
    loadArticles()
  } catch (error) {
    console.error('删除文章失败:', error)
    alert('删除失败，请重试')
  }
}

// 加载用户列表
const loadUsers = async () => {
  try {
    const params = {
      page: userPagination.page,
      size: userPagination.limit,
      keyword: userFilter.keyword
    }
    const response = await http.get('/admin/users', params)
    users.value = response.data.users || []
    userPagination.total = response.data.total || 0
    userPagination.offset = (userPagination.page - 1) * userPagination.limit
  } catch (error) {
    console.error('加载用户列表失败:', error)
  }
}

// 搜索用户
let userSearchTimeout = null
const searchUsers = () => {
  clearTimeout(userSearchTimeout)
  userSearchTimeout = setTimeout(() => {
    userPagination.page = 1
    loadUsers()
  }, 500)
}

// 切换用户状态
const toggleUserStatus = async (user) => {
  try {
    const newStatus = user.status === 1 ? 0 : 1
    await http.put(`/admin/users/${user.id}/status`, { status: newStatus })
    user.status = newStatus
  } catch (error) {
    console.error('切换用户状态失败:', error)
    alert('操作失败，请重试')
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    console.log('AdminDashboard: 开始加载分类数据...')
    const response = await http.get('/admin/categories')
    console.log('AdminDashboard: API 响应:', response)
    
    // 处理分页响应格式
    if (response.data && Array.isArray(response.data)) {
      categories.value = response.data
      console.log('AdminDashboard: 成功加载分类数据 (格式1):', categories.value.length, '个分类')
    } else if (response && response.data && Array.isArray(response.data.data)) {
      // 处理嵌套的分页响应格式
      categories.value = response.data.data
      console.log('AdminDashboard: 成功加载分类数据 (格式2):', categories.value.length, '个分类')
    } else {
      categories.value = []
      console.log('AdminDashboard: 未找到分类数据，响应格式:', typeof response.data, response.data)
    }
  } catch (error) {
    console.error('AdminDashboard: 加载分类列表失败:', error)
    categories.value = []
  }
}

// 编辑分类
const editCategory = (category) => {
  categoryForm.id = category.id
  categoryForm.name = category.name
  categoryForm.slug = category.slug
  categoryForm.description = category.description
  categoryForm.color = category.color
  showEditCategoryModal.value = true
}

// 提交分类表单
const submitCategory = async () => {
  try {
    // 如果没有填写 slug，自动从 name 生成
    if (!categoryForm.slug && categoryForm.name) {
      categoryForm.slug = categoryForm.name
        .toLowerCase()
        .replace(/\s+/g, '-')  // 空格替换为连字符
        .replace(/[^\w\-\u4e00-\u9fa5]/g, '') // 只保留字母、数字、连字符和中文
        .replace(/--+/g, '-')  // 多个连字符合并为一个
        .replace(/^-|-$/g, '') // 去除首尾连字符
    }
    
    if (showCreateCategoryModal.value) {
      await http.post('/admin/categories', categoryForm)
    } else {
      await http.put(`/admin/categories/${categoryForm.id}`, categoryForm)
    }
    closeCategoryModal()
    loadCategories()
  } catch (error) {
    console.error('保存分类失败:', error)
    alert('保存失败，请重试')
  }
}

// 关闭分类模态框
const closeCategoryModal = () => {
  showCreateCategoryModal.value = false
  showEditCategoryModal.value = false
  Object.assign(categoryForm, {
    id: null,
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6'
  })
}

// 切换分类状态
const toggleCategoryStatus = async (category) => {
  try {
    const newStatus = category.status === 1 ? 0 : 1
    await http.put(`/admin/categories/${category.id}/status`, { status: newStatus })
    category.status = newStatus
  } catch (error) {
    console.error('切换分类状态失败:', error)
    alert('操作失败，请重试')
  }
}

// 删除分类
const deleteCategory = async (categoryId) => {
  if (!confirm('确定要删除这个分类吗？')) return
  try {
    await http.delete(`/admin/categories/${categoryId}`)
    loadCategories()
  } catch (error) {
    console.error('删除分类失败:', error)
    alert('删除失败，请重试')
  }
}

// 分页控制
const prevPage = (type) => {
  if (type === 'article' && articlePagination.page > 1) {
    articlePagination.page--
    loadArticles()
  } else if (type === 'user' && userPagination.page > 1) {
    userPagination.page--
    loadUsers()
  }
}

const nextPage = (type) => {
  if (type === 'article') {
    const maxPage = Math.ceil(articlePagination.total / articlePagination.limit)
    if (articlePagination.page < maxPage) {
      articlePagination.page++
      loadArticles()
    }
  } else if (type === 'user') {
    const maxPage = Math.ceil(userPagination.total / userPagination.limit)
    if (userPagination.page < maxPage) {
      userPagination.page++
      loadUsers()
    }
  }
}

// 工具函数
const getArticleStatusText = (status) => {
  switch (status) {
    case 0: return '草稿'
    case 1: return '已发布'
    case 2: return '已下架'
    default: return '未知'
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
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
  
  loadStats()
  loadArticles()
  loadUsers()
  loadCategories()
})
</script>