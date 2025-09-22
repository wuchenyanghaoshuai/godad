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
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
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
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-orange-100 text-orange-600">
              <ResourceIcon class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">资源总数</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.resourceCount }}</p>
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
                      <div
                        class="text-sm font-medium text-gray-900 max-w-xs truncate cursor-pointer hover:text-blue-600 hover:underline"
                        @click="openArticle(article.id)"
                        :title="article.title"
                      >
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

          <!-- 文章分类管理 -->
          <div v-if="activeTab === 'categories'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">文章分类管理</h2>
              <button
                @click="showCreateCategoryModal = true"
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
              >
                新增分类
              </button>
            </div>

            <!-- 加载状态 -->
            <div v-if="categoryLoading" class="text-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
              <p class="mt-4 text-gray-500">加载中...</p>
            </div>

            <div v-else class="overflow-hidden">
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

            <!-- 分页 -->
            <div class="flex justify-between items-center mt-4">
              <span class="text-sm text-gray-700">
                显示 {{ categoryPagination.offset + 1 }} 到 {{ Math.min(categoryPagination.offset + categoryPagination.limit, categoryPagination.total) }} 条，共 {{ categoryPagination.total }} 条
              </span>
              <div class="flex space-x-2">
                <button
                  @click="prevPage('category')"
                  :disabled="categoryPagination.page === 1"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  上一页
                </button>
                <button
                  @click="nextPage('category')"
                  :disabled="categoryPagination.page >= Math.ceil(categoryPagination.total / categoryPagination.limit)"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  下一页
                </button>
              </div>
            </div>
          </div>

          <!-- 话题管理 -->
          <div v-if="activeTab === 'topics'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">话题管理</h2>
              <button
                @click="showCreateTopicModal = true"
                class="px-4 py-2 bg-pink-600 text-white rounded-md hover:bg-pink-700 text-sm"
              >
                新增话题
              </button>
            </div>

            <div class="bg-white rounded-lg shadow overflow-hidden">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">话题</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">显示名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">帖子数量</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">排序</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
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
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        @click="editTopic(topic)"
                        class="text-indigo-600 hover:text-indigo-900 mr-3"
                      >
                        编辑
                      </button>
                      <button
                        @click="toggleTopicStatus(topic)"
                        :class="topic.is_active ? 'text-red-600 hover:text-red-900' : 'text-green-600 hover:text-green-900'"
                      >
                        {{ topic.is_active ? '禁用' : '启用' }}
                      </button>
                      <button
                        @click="deleteTopic(topic.id)"
                        class="text-red-600 hover:text-red-900 ml-3"
                        :disabled="topic.post_count > 0"
                        :class="{ 'opacity-50 cursor-not-allowed': topic.post_count > 0 }"
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
                显示 {{ topicPagination.offset + 1 }} 到 {{ Math.min(topicPagination.offset + topicPagination.limit, topicPagination.total) }} 条，共 {{ topicPagination.total }} 条
              </span>
              <div class="flex space-x-2">
                <button
                  @click="loadTopics(topicPagination.page - 1)"
                  :disabled="topicPagination.page <= 1"
                  class="px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50"
                >
                  上一页
                </button>
                <button
                  @click="loadTopics(topicPagination.page + 1)"
                  :disabled="topicPagination.page >= Math.ceil(topicPagination.total / topicPagination.limit)"
                  class="px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50"
                >
                  下一页
                </button>
              </div>
            </div>
          </div>

          <!-- 资源管理 -->
          <div v-if="activeTab === 'resources'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">资源管理</h2>
              <div class="flex space-x-2">
                <select
                  v-model="resourceFilter.status"
                  @change="loadResources"
                  class="px-3 py-2 border border-gray-300 rounded-md text-sm"
                >
                  <option value="">全部状态</option>
                  <option value="0">待审核</option>
                  <option value="1">已发布</option>
                  <option value="2">已拒绝</option>
                </select>
                <select
                  v-model="resourceFilter.category"
                  @change="loadResources"
                  class="px-3 py-2 border border-gray-300 rounded-md text-sm"
                >
                  <option value="">全部分类</option>
                  <option value="E-books">电子书</option>
                  <option value="Videos">视频</option>
                  <option value="Tools">工具</option>
                </select>
                <input
                  v-model="resourceFilter.keyword"
                  @input="searchResources"
                  type="text"
                  placeholder="搜索资源..."
                  class="px-3 py-2 border border-gray-300 rounded-md text-sm"
                >
                <button
                  @click="showCreateResourceModal = true"
                  class="px-4 py-2 bg-orange-600 text-white rounded-md hover:bg-orange-700 text-sm"
                >
                  新增资源
                </button>
              </div>
            </div>

            <!-- 加载状态 -->
            <div v-if="resourceLoading" class="text-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500 mx-auto"></div>
              <p class="mt-4 text-gray-500">加载中...</p>
            </div>

            <div v-else class="overflow-hidden">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      资源信息
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      类型
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      分类
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      上传者
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      状态
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      下载量
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
                  <tr v-for="resource in resources" :key="resource.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <img
                          class="h-10 w-10 rounded object-cover"
                          :src="resource.image || '/default-resource.png'"
                          :alt="resource.title"
                        >
                        <div class="ml-4 max-w-xs">
                          <div class="text-sm font-medium text-gray-900 truncate" :title="resource.title">
                            {{ resource.title }}
                          </div>
                          <div class="text-sm text-gray-500 truncate" :title="resource.description">
                            {{ resource.description }}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          resource.type === 'e-book'
                            ? 'bg-blue-100 text-blue-800'
                            : resource.type === 'video'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-purple-100 text-purple-800'
                        ]"
                      >
                        {{ getResourceTypeText(resource.type) }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ getCategoryText(resource.category) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ resource.uploader?.nickname || '系统' }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          resource.status === 1
                            ? 'bg-green-100 text-green-800'
                            : resource.status === 0
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                        ]"
                      >
                        {{ getResourceStatusText(resource.status) }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ resource.download_count || 0 }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(resource.created_at) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        v-if="resource.status === 0"
                        @click="approveResource(resource)"
                        class="text-green-600 hover:text-green-900 mr-2"
                      >
                        通过
                      </button>
                      <button
                        v-if="resource.status === 0"
                        @click="rejectResource(resource)"
                        class="text-red-600 hover:text-red-900 mr-2"
                      >
                        拒绝
                      </button>
                      <button
                        @click="editResource(resource)"
                        class="text-blue-600 hover:text-blue-900 mr-2"
                      >
                        编辑
                      </button>
                      <button
                        @click="toggleResourceStatus(resource)"
                        :class="[
                          'mr-2',
                          resource.status === 1
                            ? 'text-yellow-600 hover:text-yellow-900'
                            : 'text-green-600 hover:text-green-900'
                        ]"
                      >
                        {{ resource.status === 1 ? '下架' : '发布' }}
                      </button>
                      <button
                        @click="deleteResource(resource.id)"
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
                显示 {{ resourcePagination.offset + 1 }} 到 {{ Math.min(resourcePagination.offset + resourcePagination.limit, resourcePagination.total) }} 条，共 {{ resourcePagination.total }} 条
              </span>
              <div class="flex space-x-2">
                <button
                  @click="prevPage('resource')"
                  :disabled="resourcePagination.page === 1"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  上一页
                </button>
                <button
                  @click="nextPage('resource')"
                  :disabled="resourcePagination.page >= Math.ceil(resourcePagination.total / resourcePagination.limit)"
                  class="px-3 py-1 border rounded-md text-sm disabled:opacity-50"
                >
                  下一页
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑资源模态框 -->
    <div v-if="showCreateResourceModal || showEditResourceModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">
          {{ showCreateResourceModal ? '新增资源' : '编辑资源' }}
        </h3>
        <form @submit.prevent="submitResource">
          <div class="grid grid-cols-2 gap-4">
            <div class="col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">资源标题</label>
              <input
                v-model="resourceForm.title"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                placeholder="请输入资源标题"
              >
            </div>
            <div class="col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">资源描述</label>
              <textarea
                v-model="resourceForm.description"
                rows="3"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                placeholder="请输入资源描述"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">资源类型</label>
              <select
                v-model="resourceForm.type"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              >
                <option value="">请选择类型</option>
                <option value="e-book">电子书</option>
                <option value="video">视频</option>
                <option value="tool">工具</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">分类</label>
              <select
                v-model="resourceForm.category"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              >
                <option value="">请选择分类</option>
                <option value="E-books">电子书</option>
                <option value="Videos">视频</option>
                <option value="Tools">工具</option>
              </select>
            </div>
            <div class="col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">封面图片URL</label>
              <input
                v-model="resourceForm.image"
                type="url"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                placeholder="https://example.com/image.jpg"
              >
            </div>
            <div class="col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">资源文件URL</label>
              <input
                v-model="resourceForm.file_url"
                type="url"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                placeholder="https://example.com/resource.pdf"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">按钮文本</label>
              <input
                v-model="resourceForm.buttonText"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                placeholder="立即下载"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">状态</label>
              <select
                v-model="resourceForm.status"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              >
                <option :value="0">待审核</option>
                <option :value="1">已发布</option>
                <option :value="2">已拒绝</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end space-x-2 mt-6">
            <button
              type="button"
              @click="closeResourceModal"
              class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              取消
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-orange-600 text-white rounded-md hover:bg-orange-700"
            >
              {{ showCreateResourceModal ? '创建' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 创建/编辑话题模态框 -->
    <div v-if="showCreateTopicModal || showEditTopicModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">
          {{ showCreateTopicModal ? '新增话题' : '编辑话题' }}
        </h3>
        <form @submit.prevent="submitTopic">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">话题名称</label>
            <input
              v-model="topicForm.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500"
              placeholder="例如: Baby Care"
            >
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">显示名称</label>
            <input
              v-model="topicForm.display_name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500"
              placeholder="例如: 婴儿护理"
            >
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">描述</label>
            <textarea
              v-model="topicForm.description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500"
              placeholder="话题描述..."
            ></textarea>
          </div>
          <div class="mb-4 flex space-x-4">
            <div class="flex-1">
              <label class="block text-sm font-medium text-gray-700 mb-2">颜色</label>
              <input
                v-model="topicForm.color"
                type="color"
                class="w-full h-10 border border-gray-300 rounded-md"
              >
            </div>
            <div class="flex-1">
              <label class="block text-sm font-medium text-gray-700 mb-2">图标</label>
              <input
                v-model="topicForm.icon"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500"
                placeholder="例如: baby, heart"
              >
            </div>
          </div>
          <div class="mb-4 flex space-x-4">
            <div class="flex-1">
              <label class="block text-sm font-medium text-gray-700 mb-2">排序</label>
              <input
                v-model.number="topicForm.sort"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500"
              >
            </div>
            <div class="flex-1">
              <label class="block text-sm font-medium text-gray-700 mb-2">状态</label>
              <select
                v-model="topicForm.is_active"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pink-500"
              >
                <option :value="true">启用</option>
                <option :value="false">禁用</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end space-x-2">
            <button
              type="button"
              @click="closeTopicModal"
              class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              取消
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-pink-600 text-white rounded-md hover:bg-pink-700"
            >
              {{ showCreateTopicModal ? '创建' : '保存' }}
            </button>
          </div>
        </form>
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
// import { http } from '../api/http' // unused
import AdminApi from '@/api/admin'
import { ResourceApi } from '@/api/resource'
import { TopicApi } from '@/api/topic'
import { useToast } from '../composables/useToast'
import { CategoryApi } from '@/api/category'

// Icons
import {
  FileTextIcon as ArticleIcon,
  UserIcon,
  TagIcon as CategoryIcon,
  MessageCircleIcon as CommentIcon,
  BookOpenIcon as ResourceIcon
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const activeTab = ref('articles')
const stats = reactive({
  articleCount: 0,
  userCount: 0,
  categoryCount: 0,
  commentCount: 0,
  resourceCount: 0
})

const tabs = [
  { id: 'articles', name: '文章管理' },
  { id: 'users', name: '用户管理' },
  { id: 'categories', name: '文章分类管理' },
  { id: 'topics', name: '话题管理' },
  { id: 'resources', name: '资源管理' }
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

// 文章分类管理
const categories = ref([])
const showCreateCategoryModal = ref(false)
const showEditCategoryModal = ref(false)
const categoryLoading = ref(false)
const categoryForm = reactive({
  id: null,
  name: '',
  slug: '',
  description: '',
  color: '#3B82F6'
})
const categoryPagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  offset: 0
})

// 话题管理
const topics = ref([])
const topicLoading = ref(false)
const showCreateTopicModal = ref(false)
const showEditTopicModal = ref(false)
const topicForm = reactive({
  id: null,
  name: '',
  display_name: '',
  description: '',
  color: '#EC4899',
  icon: '',
  sort: 0,
  is_active: true
})
const topicPagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  offset: 0
})

// 资源管理
const resources = ref([])
const resourceLoading = ref(false)
const showCreateResourceModal = ref(false)
const showEditResourceModal = ref(false)
const resourceFilter = reactive({
  keyword: '',
  status: '',
  category: ''
})
const resourceForm = reactive({
  id: null,
  title: '',
  description: '',
  type: '',
  category: '',
  image: '',
  file_url: '',
  buttonText: '',
  status: 1
})
const resourcePagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  offset: 0
})

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await AdminApi.getStats()
    stats.articleCount = response.data.articleCount
    stats.userCount = response.data.userCount
    stats.categoryCount = response.data.categoryCount
    stats.commentCount = response.data.commentCount
    stats.resourceCount = response.data.resourceStats?.total || 0
  } catch (error) {
    console.error('加载统计数据失败:', error)
    // 如果API调用失败，使用默认值
    stats.articleCount = 0
    stats.userCount = 0
    stats.categoryCount = 0
    stats.commentCount = 0
    stats.resourceCount = 0
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
    const response = await AdminApi.getArticlesPage(params)
    articles.value = response.data.items || []
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
    await AdminApi.updateArticleStatus(article.id, newStatus)
    article.status = newStatus
  } catch (error) {
    console.error('切换文章状态失败:', error)
    toast.error('操作失败，请重试')
  }
}

// 删除文章
const deleteArticle = async (articleId) => {
  if (!confirm('确定要删除这篇文章吗？')) return
  try {
    await AdminApi.deleteArticle(articleId)
    loadArticles()
  } catch (error) {
    console.error('删除文章失败:', error)
    toast.error('删除失败，请重试')
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
    const response = await AdminApi.getUsersPage(params)
    users.value = response.data.items || []
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
    await AdminApi.updateUserStatus(user.id, newStatus)
    user.status = newStatus
  } catch (error) {
    console.error('切换用户状态失败:', error)
    toast.error('操作失败，请重试')
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    categoryLoading.value = true
    const params = {
      page: categoryPagination.page,
      size: categoryPagination.limit
    }
    const response = await CategoryApi.getAdminCategoryPage(params)
    const page = response.data
    categories.value = page.items
    categoryPagination.total = page.total

    categoryPagination.offset = (categoryPagination.page - 1) * categoryPagination.limit

    // 处理边界情况：当前页没有数据且不是第一页时，自动跳转到上一页
    if (categories.value.length === 0 && categoryPagination.page > 1) {
      categoryPagination.page = Math.max(1, Math.ceil(categoryPagination.total / categoryPagination.limit))
      if (categoryPagination.page > 0) {
        await loadCategories() // 递归加载正确的页面
        return // 递归调用时不重复设置loading状态
      }
    }
  } catch (error) {
    console.error('AdminDashboard: 加载分类列表失败:', error)
    categories.value = []
    categoryPagination.total = 0
  } finally {
    categoryLoading.value = false
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
      await CategoryApi.createCategory(categoryForm as any)
      // 创建分类后重置到第一页，因为新分类可能在任何位置
      categoryPagination.page = 1
    } else {
      await CategoryApi.updateCategory(categoryForm.id as any, categoryForm as any)
      // 编辑分类不需要重置页码
    }
    closeCategoryModal()
    loadCategories()
  } catch (error) {
    console.error('保存分类失败:', error)
    toast.error('保存失败，请重试')
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
    await CategoryApi.toggleStatus(category.id, newStatus)
    category.status = newStatus
  } catch (error) {
    console.error('切换分类状态失败:', error)
    toast.error('操作失败，请重试')
  }
}

// 删除分类
const deleteCategory = async (categoryId) => {
  if (!confirm('确定要删除这个分类吗？')) return
  try {
    await CategoryApi.deleteCategory(categoryId)

    // 计算删除后的总数和当前页是否还有效
    const newTotal = categoryPagination.total - 1
    const maxPage = Math.ceil(newTotal / categoryPagination.limit)

    // 如果当前页超出范围，自动跳转到最后一页
    if (categoryPagination.page > maxPage && maxPage > 0) {
      categoryPagination.page = maxPage
    } else if (maxPage === 0) {
      // 如果没有数据了，回到第一页
      categoryPagination.page = 1
    }

    loadCategories()
  } catch (error) {
    console.error('删除分类失败:', error)
    toast.error('删除失败，请重试')
  }
}

// 话题管理函数
// 加载话题列表
const loadTopics = async (page = topicPagination.page) => {
  topicLoading.value = true
  try {
    const res = await TopicApi.getAdminTopics({ page, size: topicPagination.limit, all: true })
    const data = res.data
    topics.value = data.items as any
    topicPagination.total = data.total
    topicPagination.page = data.page
    topicPagination.offset = (topicPagination.page - 1) * topicPagination.limit
  } catch (error) {
    console.error('加载话题列表失败:', error)
    topics.value = []
    topicPagination.total = 0
    showToast('加载话题列表失败', 'error')
  } finally {
    topicLoading.value = false
  }
}

// 编辑话题
const editTopic = (topic) => {
  topicForm.id = topic.id
  topicForm.name = topic.name
  topicForm.display_name = topic.display_name
  topicForm.description = topic.description || ''
  topicForm.color = topic.color || '#EC4899'
  topicForm.icon = topic.icon || ''
  topicForm.sort = topic.sort || 0
  topicForm.is_active = topic.is_active
  showEditTopicModal.value = true
}

// 提交话题表单
const submitTopic = async () => {
  try {
    if (showCreateTopicModal.value) {
      await TopicApi.createTopic(topicForm as any)
      showToast('话题创建成功', 'success')
    } else {
      await TopicApi.updateTopic(topicForm.id as any, topicForm as any)
      showToast('话题更新成功', 'success')
    }
    closeTopicModal()
    loadTopics()
  } catch (error) {
    console.error('话题操作失败:', error)
    showToast(showCreateTopicModal.value ? '创建失败，请重试' : '更新失败，请重试', 'error')
  }
}

// 关闭话题模态框
const closeTopicModal = () => {
  showCreateTopicModal.value = false
  showEditTopicModal.value = false
  topicForm.id = null
  topicForm.name = ''
  topicForm.display_name = ''
  topicForm.description = ''
  topicForm.color = '#EC4899'
  topicForm.icon = ''
  topicForm.sort = 0
  topicForm.is_active = true
}

// 切换话题状态
const toggleTopicStatus = async (topic) => {
  try {
    const newStatus = !topic.is_active
    await TopicApi.updateTopic(topic.id, { is_active: newStatus } as any)
    topic.is_active = newStatus
    showToast(`话题已${topic.is_active ? '启用' : '禁用'}`, 'success')
  } catch (error) {
    console.error('切换话题状态失败:', error)
    showToast('操作失败，请重试', 'error')
  }
}

// 删除话题
const deleteTopic = async (topicId) => {
  if (!confirm('确定要删除这个话题吗？')) return
  try {
    await TopicApi.deleteTopic(topicId)

    // 计算删除后的总数和当前页是否还有效
    const newTotal = topicPagination.total - 1
    const maxPage = Math.ceil(newTotal / topicPagination.limit)

    // 如果当前页超出范围，自动跳转到最后一页
    if (topicPagination.page > maxPage && maxPage > 0) {
      topicPagination.page = maxPage
    } else if (maxPage === 0) {
      // 如果没有数据了，回到第一页
      topicPagination.page = 1
    }

    loadTopics()
    showToast('话题删除成功', 'success')
  } catch (error) {
    console.error('删除话题失败:', error)
    showToast('删除失败，请重试', 'error')
  }
}

// 资源管理函数
// 加载资源列表
const loadResources = async () => {
  resourceLoading.value = true
  try {
    const params = {
      page: resourcePagination.page,
      size: resourcePagination.limit,
      status: resourceFilter.status,
      category: resourceFilter.category,
      keyword: resourceFilter.keyword
    }
    const response = await ResourceApi.getAllResources(params)
    if (response.data && response.data.items) {
      resources.value = response.data.items
      resourcePagination.total = response.data.total || 0
    } else {
      resources.value = []
      resourcePagination.total = 0
    }

    resourcePagination.offset = (resourcePagination.page - 1) * resourcePagination.limit
  } catch (error) {
    console.error('加载资源列表失败:', error)
    resources.value = []
    resourcePagination.total = 0
    showToast('加载资源列表失败', 'error')
  } finally {
    resourceLoading.value = false
  }
}

// 搜索资源
let resourceSearchTimeout = null
const searchResources = () => {
  clearTimeout(resourceSearchTimeout)
  resourceSearchTimeout = setTimeout(() => {
    resourcePagination.page = 1
    loadResources()
  }, 500)
}

// 编辑资源
const editResource = (resource) => {
  resourceForm.id = resource.id
  resourceForm.title = resource.title
  resourceForm.description = resource.description
  resourceForm.type = resource.type
  resourceForm.category = resource.category
  resourceForm.image = resource.image
  resourceForm.file_url = resource.file_url
  resourceForm.buttonText = resource.buttonText
  resourceForm.status = resource.status
  showEditResourceModal.value = true
}

// 提交资源表单
const submitResource = async () => {
  try {
    if (showCreateResourceModal.value) {
      await ResourceApi.adminCreateResource(resourceForm as any)
      showToast('资源创建成功', 'success')
    } else {
      await ResourceApi.updateResource(resourceForm.id as any, resourceForm as any)
      showToast('资源更新成功', 'success')
    }
    closeResourceModal()
    loadResources()
  } catch (error) {
    console.error('资源操作失败:', error)
    showToast(showCreateResourceModal.value ? '创建失败，请重试' : '更新失败，请重试', 'error')
  }
}

// 关闭资源模态框
const closeResourceModal = () => {
  showCreateResourceModal.value = false
  showEditResourceModal.value = false
  Object.assign(resourceForm, {
    id: null,
    title: '',
    description: '',
    type: '',
    category: '',
    image: '',
    file_url: '',
    buttonText: '',
    status: 1
  })
}

// 审核资源
const approveResource = async (resource) => {
  try {
    await ResourceApi.updateResourceStatus(resource.id, 1)
    resource.status = 1
    showToast('资源审核通过', 'success')
  } catch (error) {
    console.error('审核资源失败:', error)
    showToast('操作失败，请重试', 'error')
  }
}

// 拒绝资源
const rejectResource = async (resource) => {
  if (!confirm('确定要拒绝这个资源吗？')) return
  try {
    await ResourceApi.updateResourceStatus(resource.id, 2)
    resource.status = 2
    showToast('资源已拒绝', 'success')
  } catch (error) {
    console.error('拒绝资源失败:', error)
    showToast('操作失败，请重试', 'error')
  }
}

// 切换资源状态
const toggleResourceStatus = async (resource) => {
  try {
    const newStatus = resource.status === 1 ? 0 : 1
    await ResourceApi.updateResourceStatus(resource.id, newStatus)
    resource.status = newStatus
    showToast(`资源已${newStatus === 1 ? '发布' : '下架'}`, 'success')
  } catch (error) {
    console.error('切换资源状态失败:', error)
    showToast('操作失败，请重试', 'error')
  }
}

// 删除资源
const deleteResource = async (resourceId) => {
  if (!confirm('确定要删除这个资源吗？')) return
  try {
    await ResourceApi.deleteResource(resourceId)

    const newTotal = resourcePagination.total - 1
    const maxPage = Math.ceil(newTotal / resourcePagination.limit)

    if (resourcePagination.page > maxPage && maxPage > 0) {
      resourcePagination.page = maxPage
    } else if (maxPage === 0) {
      resourcePagination.page = 1
    }

    loadResources()
    showToast('资源删除成功', 'success')
  } catch (error) {
    console.error('删除资源失败:', error)
    showToast('删除失败，请重试', 'error')
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
  } else if (type === 'category' && categoryPagination.page > 1) {
    categoryPagination.page--
    loadCategories()
  } else if (type === 'topic' && topicPagination.page > 1) {
    topicPagination.page--
    loadTopics()
  } else if (type === 'resource' && resourcePagination.page > 1) {
    resourcePagination.page--
    loadResources()
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
  } else if (type === 'category') {
    const maxPage = Math.ceil(categoryPagination.total / categoryPagination.limit)
    if (categoryPagination.page < maxPage) {
      categoryPagination.page++
      loadCategories()
    }
  } else if (type === 'topic') {
    const maxPage = Math.ceil(topicPagination.total / topicPagination.limit)
    if (topicPagination.page < maxPage) {
      topicPagination.page++
      loadTopics()
    }
  } else if (type === 'resource') {
    const maxPage = Math.ceil(resourcePagination.total / resourcePagination.limit)
    if (resourcePagination.page < maxPage) {
      resourcePagination.page++
      loadResources()
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

const getResourceStatusText = (status) => {
  switch (status) {
    case 0: return '待审核'
    case 1: return '已发布'
    case 2: return '已拒绝'
    default: return '未知'
  }
}

const getResourceTypeText = (type) => {
  switch (type) {
    case 'e-book': return '电子书'
    case 'video': return '视频'
    case 'tool': return '工具'
    default: return type
  }
}

const getCategoryText = (category) => {
  switch (category) {
    case 'E-books': return '电子书'
    case 'Videos': return '视频'
    case 'Tools': return '工具'
    default: return category
  }
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}

// 打开文章页面
const openArticle = (articleId) => {
  const url = `/articles/${articleId}`
  window.open(url, '_blank')
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
  loadTopics()
  loadResources()
})
</script>
