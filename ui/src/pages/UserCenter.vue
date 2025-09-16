<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 头部导航 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <router-link to="/" class="text-xl font-bold text-gray-900">
            GoDad
          </router-link>
          <div class="flex items-center space-x-4">
            <span class="text-gray-700">{{ user?.username }}</span>
            <button
              @click="handleLogout"
              class="text-gray-500 hover:text-gray-700"
            >
              <LogOutIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 侧边栏 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow p-6">
            <!-- 用户头像和基本信息 -->
            <div class="text-center mb-6">
              <div class="w-24 h-24 mx-auto mb-4 relative group cursor-pointer" @click="showAvatarModal = true">
                <img
                  v-if="user?.avatar"
                  :src="user.avatar"
                  :alt="user.username"
                  class="w-24 h-24 rounded-full object-cover transition-opacity group-hover:opacity-75"
                />
                <div
                  v-else
                  class="w-24 h-24 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center transition-opacity group-hover:opacity-75"
                >
                  <span class="text-2xl font-bold text-white">
                    {{ user?.username?.charAt(0).toUpperCase() }}
                  </span>
                </div>
                <!-- 悬停提示 -->
                <div class="absolute inset-0 rounded-full bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                  <CameraIcon class="h-6 w-6 text-white" />
                </div>
              </div>
              <p class="text-xs text-gray-500 mb-2">点击头像更换</p>
              <h2 class="text-xl font-semibold text-gray-900">{{ user?.username }}</h2>
              <p class="text-gray-600">{{ user?.email }}</p>
              <span class="inline-block px-2 py-1 text-xs font-medium rounded-full mt-2"
                    :class="roleClasses">
                {{ roleText }}
              </span>
              
              <!-- 关注统计信息 -->
              <div class="flex justify-center space-x-6 mt-4 pt-4 border-t border-gray-200">
                <div class="text-center cursor-pointer hover:text-pink-600 transition-colors" @click="activeTab = 'following'">
                  <div class="text-lg font-semibold">{{ followingCount }}</div>
                  <div class="text-xs text-gray-500">关注</div>
                </div>
                <div class="text-center cursor-pointer hover:text-pink-600 transition-colors" @click="activeTab = 'followers'">
                  <div class="text-lg font-semibold">{{ followersCount }}</div>
                  <div class="text-xs text-gray-500">粉丝</div>
                </div>
                <div class="text-center">
                  <div class="text-lg font-semibold">{{ articlesCount }}</div>
                  <div class="text-xs text-gray-500">文章</div>
                </div>
              </div>
            </div>

            <!-- 导航菜单 -->
            <nav class="space-y-2">
              <button
                v-for="item in menuItems"
                :key="item.key"
                @click="activeTab = item.key"
                :class="[
                  'w-full flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                  activeTab === item.key
                    ? 'bg-pink-50 text-pink-700 border-r-2 border-pink-500'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                ]"
              >
                <component :is="item.icon" class="h-5 w-5 mr-3" />
                {{ item.label }}
              </button>
            </nav>
          </div>
        </div>

        <!-- 主内容区 -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-lg shadow">
            <!-- 个人信息 -->
            <div v-if="activeTab === 'profile'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">个人信息</h3>
              <form @submit.prevent="updateProfile" class="space-y-6">

                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      用户名
                    </label>
                    <input
                      v-model="profileForm.username"
                      type="text"
                      disabled
                      class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-500"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      邮箱
                    </label>
                    <input
                      v-model="profileForm.email"
                      type="email"
                      disabled
                      class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-500"
                    />
                  </div>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    昵称
                  </label>
                  <input
                    v-model="profileForm.nickname"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                    placeholder="请输入昵称"
                  />
                </div>
                
                <!-- 手机号和性别 -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      手机号
                    </label>
                    <input
                      v-model="profileForm.phone"
                      type="tel"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      placeholder="请输入手机号"
                      pattern="[0-9]{11}"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      性别
                    </label>
                    <select
                      v-model="profileForm.gender"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                    >
                      <option value="">请选择性别</option>
                      <option value="1">男</option>
                      <option value="2">女</option>
                      <option value="0">保密</option>
                    </select>
                  </div>
                </div>

                <!-- 生日 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    生日
                  </label>
                  <input
                    v-model="profileForm.birthday"
                    type="date"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    个人简介
                  </label>
                  <textarea
                    v-model="profileForm.bio"
                    rows="4"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                    placeholder="介绍一下自己..."
                  ></textarea>
                </div>
                <div class="flex justify-end">
                  <button
                    type="submit"
                    :disabled="isUpdating"
                    class="bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700 disabled:opacity-50"
                  >
                    {{ isUpdating ? '保存中...' : '保存更改' }}
                  </button>
                </div>
              </form>
            </div>

            <!-- 我的文章 -->
            <div v-else-if="activeTab === 'articles'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的文章</h3>
                <router-link
                  to="/articles/create"
                  class="bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700"
                >
                  写文章
                </router-link>
              </div>


              <!-- 加载状态 -->
              <div v-if="isLoadingArticles" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="mt-4 text-gray-500">加载中...</p>
              </div>

              <!-- 错误状态 -->
              <div v-else-if="articlesError" class="text-center py-12 text-red-500">
                <p>{{ articlesError }}</p>
                <button 
                  @click="loadMyArticles"
                  class="mt-4 bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700"
                >
                  重试
                </button>
              </div>

              <!-- 文章列表 -->
              <div v-else-if="myArticles.length > 0" class="space-y-4">
                <div 
                  v-for="article in myArticles" 
                  :key="article.id"
                  class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
                  @click="router.push(`/articles/${article.id}`)"
                >
                  <div class="flex justify-between items-start">
                    <div class="flex-1">
                      <h4 class="text-lg font-medium text-gray-900 hover:text-pink-600">
                        {{ article.title }}
                      </h4>
                      <p v-if="article.summary" class="text-gray-600 mt-2 text-sm line-clamp-2">
                        {{ article.summary }}
                      </p>
                      <div class="flex items-center mt-3 space-x-4 text-sm text-gray-500">
                        <span>{{ formatDate(article.created_at) }}</span>
                        <span>阅读 {{ article.view_count || 0 }}</span>
                        <span>点赞 {{ article.like_count || 0 }}</span>
                        <span 
                          class="px-2 py-1 rounded-full text-xs"
                          :class="getStatusClass(article.status)"
                        >
                          {{ getStatusText(article.status) }}
                        </span>
                      </div>
                    </div>
                    <div class="ml-4 flex space-x-2">
                      <router-link
                        :to="`/articles/${article.id}/edit`"
                        class="text-blue-600 hover:text-blue-700 text-sm"
                        @click.stop
                      >
                        编辑
                      </router-link>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <FileTextIcon class="h-12 w-12 mx-auto mb-4" />
                <p>您还没有发布任何文章</p>
                <router-link
                  to="/articles/create"
                  class="text-pink-600 hover:text-pink-500 mt-2 inline-block"
                >
                  立即创建第一篇文章
                </router-link>
              </div>
            </div>


            <!-- 消息通知 -->
            <div v-else-if="activeTab === 'notifications'" class="p-6">
              <NotificationList :auto-refresh="true" />
            </div>

            <!-- 我的关注 -->
            <div v-else-if="activeTab === 'following'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的关注 ({{ followingCount }})</h3>
              </div>
              
              <div v-if="isLoadingFollows" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="text-gray-500 mt-2">加载中...</p>
              </div>
              
              <!-- 关注列表 -->
              <div v-else-if="followingList.length > 0" class="space-y-4">
                <div v-for="user in followingList" :key="user.id" class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-r from-pink-400 to-orange-400 flex items-center justify-center">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.username" class="w-full h-full object-cover">
                      <span v-else class="text-white font-semibold">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <div class="flex-1">
                      <h4 class="font-medium text-gray-900">{{ user.nickname || user.username }}</h4>
                      <p class="text-sm text-gray-500">@{{ user.username }}</p>
                      <p v-if="user.bio" class="text-sm text-gray-600 mt-1">{{ user.bio }}</p>
                      <p v-if="user.followed_at" class="text-xs text-gray-400 mt-1">{{ formatFollowDate(user.followed_at) }}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span v-if="user.is_mutual_follow" class="text-xs bg-pink-100 text-pink-600 px-2 py-1 rounded-full">互关</span>
                    <button 
                      @click="unfollowUser(user.id)"
                      class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm"
                    >
                      取消关注
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <HeartIcon class="h-12 w-12 mx-auto mb-4" />
                <p>您还没有关注任何人</p>
                <p class="text-sm mt-2">去发现更多有趣的用户吧</p>
              </div>
            </div>

            <!-- 我的粉丝 -->
            <div v-else-if="activeTab === 'followers'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的粉丝 ({{ followersCount }})</h3>
              </div>
              
              <div v-if="isLoadingFollows" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="text-gray-500 mt-2">加载中...</p>
              </div>
              
              <!-- 粉丝列表 -->
              <div v-else-if="followersList.length > 0" class="space-y-4">
                <div v-for="user in followersList" :key="user.id" class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-r from-pink-400 to-orange-400 flex items-center justify-center">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.username" class="w-full h-full object-cover">
                      <span v-else class="text-white font-semibold">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <div class="flex-1">
                      <h4 class="font-medium text-gray-900">{{ user.nickname || user.username }}</h4>
                      <p class="text-sm text-gray-500">@{{ user.username }}</p>
                      <p v-if="user.bio" class="text-sm text-gray-600 mt-1">{{ user.bio }}</p>
                      <p v-if="user.followed_at" class="text-xs text-gray-400 mt-1">{{ formatFollowDate(user.followed_at) }}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span
                      v-if="user.is_mutual_follow"
                      class="text-xs bg-pink-100 text-pink-600 px-3 py-2 rounded-full"
                    >
                      互相关注
                    </span>
                    <button
                      v-else-if="!user.is_following"
                      @click="followUser(user.id)"
                      class="px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 transition-colors text-sm"
                    >
                      回关
                    </button>
                    <button
                      v-else
                      @click="unfollowUser(user.id)"
                      class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm"
                    >
                      取消关注
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <UsersIcon class="h-12 w-12 mx-auto mb-4" />
                <p>还没有人关注您</p>
                <p class="text-sm mt-2">发布有趣的内容来吸引更多关注吧</p>
              </div>
            </div>

            <!-- 互相关注 -->
            <div v-else-if="activeTab === 'mutual'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">互相关注 ({{ mutualCount }})</h3>
              </div>
              
              <div v-if="isLoadingFollows" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600 mx-auto"></div>
                <p class="text-gray-500 mt-2">加载中...</p>
              </div>
              
              <!-- 互关列表 -->
              <div v-else-if="mutualFollowsList.length > 0" class="space-y-4">
                <div v-for="user in mutualFollowsList" :key="user.id" class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-r from-pink-400 to-orange-400 flex items-center justify-center">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.username" class="w-full h-full object-cover">
                      <span v-else class="text-white font-semibold">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <div class="flex-1">
                      <h4 class="font-medium text-gray-900">{{ user.nickname || user.username }}</h4>
                      <p class="text-sm text-gray-500">@{{ user.username }}</p>
                      <p v-if="user.bio" class="text-sm text-gray-600 mt-1">{{ user.bio }}</p>
                      <p v-if="user.followed_at" class="text-xs text-gray-400 mt-1">{{ formatFollowDate(user.followed_at) }}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span class="text-xs bg-pink-100 text-pink-600 px-3 py-2 rounded-full font-medium">💕 互相关注</span>
                    <button 
                      @click="unfollowUser(user.id)"
                      class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm"
                    >
                      取消关注
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <HeartIcon class="h-12 w-12 mx-auto mb-4 text-pink-300" />
                <p>还没有互相关注的用户</p>
                <p class="text-sm mt-2">关注其他用户，等待他们回关吧</p>
              </div>
            </div>

            <!-- 我的收藏 -->
            <div v-else-if="activeTab === 'favorites'" class="p-6">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-lg font-medium text-gray-900">我的收藏</h3>
              </div>
              
              <!-- 收藏列表 -->
              <div v-if="favoritesLoading" class="text-center py-12">
                <LoaderIcon class="h-8 w-8 mx-auto mb-4 text-gray-400 animate-spin" />
                <p class="text-gray-500">加载中...</p>
              </div>

              <div v-else-if="favorites.length === 0" class="text-center py-12 text-gray-500">
                <StarIcon class="h-12 w-12 mx-auto mb-4 text-yellow-400" />
                <p>暂无收藏的文章</p>
                <p class="text-sm mt-2">去收藏一些有趣的文章吧</p>
              </div>

              <div v-else class="space-y-4">
                <div
                  v-for="favorite in favorites"
                  :key="favorite.id"
                  class="bg-gray-50 rounded-lg p-4 hover:bg-gray-100 transition-colors cursor-pointer"
                  @click="goToArticle(favorite.article_id)"
                >
                  <div class="flex items-start space-x-4">
                    <img
                      v-if="favorite.article?.cover_image"
                      :src="favorite.article.cover_image"
                      :alt="favorite.article.title"
                      class="w-16 h-16 object-cover rounded-lg flex-shrink-0"
                    />
                    <div v-else class="w-16 h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-lg flex-shrink-0 flex items-center justify-center">
                      <FileTextIcon class="h-6 w-6 text-white" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <h4 class="text-sm font-medium text-gray-900 truncate mb-1">
                        {{ favorite.article?.title || '文章标题' }}
                      </h4>
                      <p class="text-xs text-gray-500 mb-2">
                        {{ favorite.article?.excerpt || '暂无摘要' }}
                      </p>
                      <div class="flex items-center space-x-4 text-xs text-gray-400">
                        <span>收藏时间: {{ formatDate(favorite.created_at) }}</span>
                        <div class="flex items-center space-x-1">
                          <EyeIcon class="h-3 w-3" />
                          <span>{{ favorite.article?.view_count || 0 }}</span>
                        </div>
                        <div class="flex items-center space-x-1">
                          <HeartIcon class="h-3 w-3" />
                          <span>{{ favorite.article?.like_count || 0 }}</span>
                        </div>
                      </div>
                    </div>
                    <button
                      @click.stop="removeFavorite(favorite.id)"
                      class="text-gray-400 hover:text-red-500 transition-colors"
                      title="取消收藏"
                    >
                      <TrashIcon class="h-4 w-4" />
                    </button>
                  </div>
                </div>

                <!-- 分页 -->
                <div v-if="favoritesPagination.total_pages > 1" class="flex justify-center mt-6">
                  <nav class="flex items-center space-x-2">
                    <button
                      @click="loadFavorites(favoritesPagination.current_page - 1)"
                      :disabled="favoritesPagination.current_page <= 1 || favoritesLoading"
                      class="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      上一页
                    </button>
                    <span class="text-sm text-gray-600">
                      第 {{ favoritesPagination.current_page }} 页，共 {{ favoritesPagination.total_pages }} 页
                    </span>
                    <button
                      @click="loadFavorites(favoritesPagination.current_page + 1)"
                      :disabled="favoritesPagination.current_page >= favoritesPagination.total_pages || favoritesLoading"
                      class="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      下一页
                    </button>
                  </nav>
                </div>
              </div>
            </div>

            <!-- 设置 -->
            <div v-else-if="activeTab === 'settings'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">账户设置</h3>
              <div class="space-y-6">
                <!-- 修改密码 -->
                <div class="border-b border-gray-200 pb-6">
                  <h4 class="text-md font-medium text-gray-900 mb-4">修改密码</h4>
                  <form @submit.prevent="changePassword" class="space-y-4">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        当前密码
                      </label>
                      <input
                        v-model="passwordForm.currentPassword"
                        type="password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        新密码
                      </label>
                      <input
                        v-model="passwordForm.newPassword"
                        type="password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        确认新密码
                      </label>
                      <input
                        v-model="passwordForm.confirmPassword"
                        type="password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-pink-500 focus:border-pink-500"
                      />
                    </div>
                    <button
                      type="submit"
                      :disabled="isChangingPassword"
                      class="bg-pink-600 text-white px-4 py-2 rounded-md hover:bg-pink-700 disabled:opacity-50"
                    >
                      {{ isChangingPassword ? '修改中...' : '修改密码' }}
                    </button>
                  </form>
                </div>

                <!-- 危险操作 -->
                <div>
                  <h4 class="text-md font-medium text-red-600 mb-4">危险操作</h4>
                  <button
                    @click="handleLogout"
                    class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
                  >
                    退出登录
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 头像上传弹窗 -->
    <AvatarModal
      :is-visible="showAvatarModal"
      @close="closeAvatarModal"
      @success="handleAvatarUpload"
      @error="handleUploadError"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  UserIcon,
  FileTextIcon,
  StarIcon,
  SettingsIcon,
  LogOutIcon,
  CameraIcon,
  UsersIcon,
  HeartIcon,
  BellIcon,
  LoaderIcon,
  EyeIcon,
  TrashIcon
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import AvatarModal from '../components/AvatarModal.vue'
import NotificationList from '../components/NotificationList.vue'
import { useToast } from '../composables/useToast'
import { useUserDataSync } from '../composables/useUserDataSync'
import { UserApi } from '../api/user'
import { ArticleApi } from '../api/article'
import { FollowApi } from '../api/follow'
import { FavoriteApi } from '../api/favorite'
import type { ImageUploadResponse, Article, Favorite } from '../api/types'

// 路由
const router = useRouter()

// 认证store
const authStore = useAuthStore()
const { toast } = useToast()

// 用户数据同步
const userDataSync = useUserDataSync()

// 响应式数据
const showAvatarModal = ref(false)
const activeTab = ref('profile')
const isUpdating = ref(false)
const isChangingPassword = ref(false)

// 收藏相关
const favorites = ref<Favorite[]>([])
const favoritesLoading = ref(false)
const favoritesPagination = ref({
  total: 0,
  current_page: 1,
  per_page: 10,
  total_pages: 0
})

// 使用同步的数据 - 直接使用userDataSync返回的computed属性
const myArticles = userDataSync.articles
const isLoadingArticles = userDataSync.isLoading
const followingCount = computed(() => userDataSync.stats.value.following_count)
const followersCount = computed(() => userDataSync.stats.value.followers_count)
const mutualCount = computed(() => userDataSync.mutualFollowsList.value.length)
const articlesCount = computed(() => userDataSync.stats.value.article_count)
const followingList = userDataSync.followingList
const followersList = userDataSync.followersList
const mutualFollowsList = userDataSync.mutualFollowsList
const isLoadingFollows = userDataSync.isLoading
const articlesError = ref('')

// 用户信息
const user = computed(() => userDataSync.profile.value || authStore.user)

// 角色显示
const roleText = computed(() => {
  switch (user.value?.role) {
    case 'admin':
      return '管理员'
    case 'content_manager':
      return '内容管理员'
    default:
      return '普通用户'
  }
})

const roleClasses = computed(() => {
  switch (user.value?.role) {
    case 'admin':
      return 'bg-red-100 text-red-800'
    case 'content_manager':
      return 'bg-blue-100 text-blue-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})

// 菜单项
const menuItems = [
  { key: 'profile', label: '个人信息', icon: UserIcon },
  { key: 'articles', label: '我的文章', icon: FileTextIcon },
  { key: 'favorites', label: '我的收藏', icon: StarIcon },
  { key: 'notifications', label: '消息通知', icon: BellIcon },
  { key: 'following', label: '我的关注', icon: HeartIcon },
  { key: 'followers', label: '我的粉丝', icon: UsersIcon },
  { key: 'mutual', label: '互相关注', icon: HeartIcon },
  { key: 'settings', label: '设置', icon: SettingsIcon }
]

// 个人信息表单
const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  phone: '',
  gender: '',
  birthday: '',
  bio: '',
  avatar: ''
})

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 初始化用户信息
const initUserInfo = () => {
  if (user.value) {
    profileForm.username = user.value.username
    profileForm.email = user.value.email
    profileForm.nickname = user.value.nickname || ''
    profileForm.phone = user.value.phone || ''
    profileForm.gender = user.value.gender?.toString() || ''
    profileForm.birthday = user.value.birthday ? user.value.birthday.split('T')[0] : ''
    profileForm.bio = user.value.bio || ''
    profileForm.avatar = user.value.avatar || ''
  }
}

// 头像上传成功处理
const handleAvatarUpload = async (response: ImageUploadResponse) => {
  const avatarUrl = response.url
  
  if (avatarUrl) {
    try {
      // 更新表单数据
      profileForm.avatar = avatarUrl
      
      // 使用数据同步函数更新头像
      await userDataSync.updateUserProfile({
        avatar: avatarUrl
      })
      
      // 确保DOM更新
      await nextTick()
      
      // 关闭头像上传弹窗
      closeAvatarModal()
    } catch (error) {
      console.error('头像更新失败:', error)
      toast.error('头像更新失败，请重试')
    }
  }
}

// 关闭头像弹窗
const closeAvatarModal = () => {
  showAvatarModal.value = false
}

// 上传错误处理
const handleUploadError = (error: string) => {
  toast.error(`上传失败: ${error}`)
}

// 更新个人信息
const updateProfile = async () => {
  try {
    isUpdating.value = true
    
    // 调用更新用户信息API
    const updateData: any = {
      nickname: profileForm.nickname,
      phone: profileForm.phone,
      gender: profileForm.gender ? parseInt(profileForm.gender) : undefined,
      birthday: profileForm.birthday || undefined,
      bio: profileForm.bio,
      avatar: profileForm.avatar
    }
    
    // 移除空值
    Object.keys(updateData).forEach(key => {
      if (updateData[key] === '' || updateData[key] === undefined) {
        delete updateData[key]
      }
    })
    
    // 使用数据同步函数更新
    await userDataSync.updateUserProfile(updateData)
    
    // 重新初始化表单
    initUserInfo()
  } catch (error) {
    console.error('更新失败:', error)
    toast.error('更新失败，请重试')
  } finally {
    isUpdating.value = false
  }
}

// 修改密码
const changePassword = async () => {
  try {
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      alert('新密码和确认密码不一致')
      return
    }

    isChangingPassword.value = true
    // TODO: 调用修改密码API
    
    // 重置表单
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error) {
    console.error('修改密码失败:', error)
  } finally {
    isChangingPassword.value = false
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}

// 加载我的文章（使用数据同步）
const loadMyArticles = async () => {
  try {
    articlesError.value = ''
    await userDataSync.loadUserArticles({ page: 1, size: 20 })
  } catch (error: any) {
    articlesError.value = error.message || '加载文章失败'
    console.error('加载我的文章失败:', error)
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 格式化关注时间
const formatFollowDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)
  
  if (diffInSeconds < 60) {
    return '刚刚关注'
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60)
    return `${minutes}分钟前关注`
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600)
    return `${hours}小时前关注`
  } else if (diffInSeconds < 2592000) {
    const days = Math.floor(diffInSeconds / 86400)
    return `${days}天前关注`
  } else {
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }) + ' 关注'
  }
}

// 获取状态文本
const getStatusText = (status: number) => {
  switch (status) {
    case 0: return '草稿'
    case 1: return '已发布'
    case 2: return '已下架'
    default: return '未知'
  }
}

// 获取状态样式
const getStatusClass = (status: number) => {
  switch (status) {
    case 0: return 'bg-gray-100 text-gray-600'
    case 1: return 'bg-green-100 text-green-600'
    case 2: return 'bg-red-100 text-red-600'
    default: return 'bg-gray-100 text-gray-600'
  }
}

// 加载关注统计信息（使用数据同步）
const loadFollowStats = async () => {
  try {
    await userDataSync.loadUserStats()
  } catch (error) {
    console.error('加载关注统计失败:', error)
  }
}

// 加载关注列表（使用数据同步）
const loadFollowing = async () => {
  try {
    await userDataSync.loadFollowingList({ page: 1, limit: 50 })
  } catch (error) {
    console.error('加载关注列表失败:', error)
    toast.error('加载关注列表失败')
  }
}

// 加载粉丝列表（使用数据同步）
const loadFollowers = async () => {
  try {
    await userDataSync.loadFollowersList({ page: 1, limit: 50 })
  } catch (error) {
    console.error('加载粉丝列表失败:', error)
    toast.error('加载粉丝列表失败')
  }
}

// 取消关注（使用数据同步）
const unfollowUser = async (userId) => {
  try {
    await userDataSync.unfollowUser(userId)
    // 数据同步函数会自动更新相关统计和列表
  } catch (error) {
    console.error('取消关注失败:', error)
  }
}

// 关注用户（使用数据同步）
const followUser = async (userId) => {
  try {
    await userDataSync.followUser(userId)
    // 数据同步函数会自动更新相关统计和列表
  } catch (error) {
    console.error('关注失败:', error)
  }
}

// 加载互相关注列表（使用数据同步）
const loadMutualFollows = async () => {
  try {
    await userDataSync.loadMutualFollowsList({ page: 1, limit: 50 })
  } catch (error) {
    console.error('加载互关列表失败:', error)
    toast.error('加载互关列表失败')
  }
}

// 收藏相关函数
const loadFavorites = async (page = 1) => {
  try {
    favoritesLoading.value = true
    const response = await FavoriteApi.getUserFavorites({
      page,
      size: favoritesPagination.value.per_page
    })

    favorites.value = response.data.favorites
    favoritesPagination.value = response.data.pagination
  } catch (error) {
    console.error('加载收藏列表失败:', error)
    toast.error('加载收藏列表失败')
  } finally {
    favoritesLoading.value = false
  }
}

const removeFavorite = async (favoriteId: number) => {
  if (!confirm('确定要取消收藏这篇文章吗？')) {
    return
  }

  try {
    await FavoriteApi.deleteFavorite(favoriteId)
    toast.success('取消收藏成功')
    // 重新加载当前页的收藏列表
    await loadFavorites(favoritesPagination.value.current_page)
  } catch (error) {
    console.error('取消收藏失败:', error)
    toast.error('取消收藏失败')
  }
}

const goToArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

// 监听activeTab变化，根据标签加载不同数据
watch(activeTab, (newTab) => {
  if (newTab === 'articles') {
    loadMyArticles()
  } else if (newTab === 'following') {
    loadFollowing()
  } else if (newTab === 'followers') {
    loadFollowers()
  } else if (newTab === 'mutual') {
    loadMutualFollows()
  } else if (newTab === 'favorites') {
    loadFavorites()
  }
})

// 组件挂载时初始化
onMounted(() => {
  // 检查登录状态
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  
  // 处理URL参数中的tab
  const urlParams = new URLSearchParams(window.location.search)
  const tab = urlParams.get('tab')
  if (tab && ['profile', 'articles', 'notifications', 'following', 'followers', 'mutual', 'settings'].includes(tab)) {
    activeTab.value = tab
  }
  
  initUserInfo()
  // 初始化用户数据同步
  userDataSync.initUserData()
})
</script>