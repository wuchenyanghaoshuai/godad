<template>
  <AppLayout
    :header-config="{
      showSearch: false,
      showCreateButton: false,
      showNotifications: true,
      showUserPoints: true,
      showNavigation: false,
      showUserMenu: true
    }"
    :show-footer="false"
    background-class="bg-gray-50"
  >
    <PageContainer background="gray" padding="lg">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 侧边栏 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow p-6">
            <!-- 用户头像和基本信息 -->
            <div class="text-center mb-6">
              <div class="w-24 h-24 mx-auto mb-4 relative group cursor-pointer" @click="showAvatarModal = true">
                <UserAvatar
                  :avatar="user?.avatar || ''"
                  :name="user?.nickname || user?.username || 'U'"
                  :size="96"
                />
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

              <!-- 用户等级和积分 -->
              <div class="mt-4 pt-4 border-t border-gray-200">
                <UserPointsDisplay
                  mode="detailed"
                  :auto-refresh="true"
                  class="mb-4"
                />
              </div>

              <!-- 关注统计信息 -->
              <div class="flex justify-center space-x-6 mt-4 pt-4 border-t border-gray-200">
                <div class="text-center cursor-pointer hover:text-primary-600 transition-colors" @click="activeTab = 'following'">
                  <div class="text-lg font-semibold">{{ followingCount }}</div>
                  <div class="text-xs text-gray-500">关注</div>
                </div>
                <div class="text-center cursor-pointer hover:text-primary-600 transition-colors" @click="activeTab = 'followers'">
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
                    ? 'bg-primary-50 text-primary-700 border-r-2 border-primary-500'
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
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
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
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
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
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    个人简介
                  </label>
                  <textarea
                    v-model="profileForm.bio"
                    rows="4"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                    placeholder="介绍一下自己..."
                  ></textarea>
                </div>
                <div class="flex justify-end">
                  <button
                    type="submit"
                    :disabled="isUpdating"
                    class="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 disabled:opacity-50"
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
                  class="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700"
                >
                  写文章
                </router-link>
              </div>

              <!-- 加载状态 -->
              <div v-if="isLoadingArticles" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                <p class="mt-4 text-gray-500">加载中...</p>
              </div>

              <!-- 错误状态 -->
              <div v-else-if="articlesError" class="text-center py-12 text-red-500">
                <p>{{ articlesError }}</p>
                <button
                  @click="loadMyArticles"
                  class="mt-4 bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700"
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
                      <h4 class="text-lg font-medium text-gray-900 hover:text-primary-600">
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
                        class="text-primary-600 hover:text-primary-700 text-sm"
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
                  class="text-primary-600 hover:text-primary-500 mt-2 inline-block"
                >
                  立即创建第一篇文章
                </router-link>
              </div>
            </div>

            <!-- 我的收藏 -->
            <div v-else-if="activeTab === 'favorites'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">我的收藏</h3>

              <!-- 加载状态 -->
              <div v-if="isLoadingFavorites" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                <p class="mt-4 text-gray-500">加载中...</p>
              </div>

              <!-- 错误状态 -->
              <div v-else-if="favoritesError" class="text-center py-12 text-red-500">
                <p>{{ favoritesError }}</p>
                <button
                  @click="loadMyFavorites"
                  class="mt-4 bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700"
                >
                  重试
                </button>
              </div>

              <!-- 收藏列表 -->
              <div v-else-if="myFavorites && myFavorites.length > 0" class="space-y-4">
                <div
                  v-for="favorite in myFavorites"
                  :key="favorite.id"
                  class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
                  @click="router.push(`/articles/${favorite.article?.id}`)"
                >
                  <div class="flex justify-between items-start">
                    <div class="flex-1">
                      <h4 class="text-lg font-medium text-gray-900 hover:text-primary-600">
                        {{ favorite.article?.title }}
                      </h4>
                      <p v-if="favorite.article?.excerpt" class="text-gray-600 mt-2 text-sm line-clamp-2">
                        {{ favorite.article.excerpt }}
                      </p>
                      <div class="flex items-center mt-3 space-x-4 text-sm text-gray-500">
                        <span>{{ formatFavoriteTime(favorite.created_at) }}</span>
                        <span>阅读 {{ favorite.article?.view_count || 0 }}</span>
                        <span>点赞 {{ favorite.article?.like_count || 0 }}</span>
                        <span>收藏 {{ favorite.article?.favorite_count || 0 }}</span>
                      </div>
                    </div>
                    <div class="ml-4 flex space-x-2">
                      <button
                        @click.stop="removeFavorite(favorite)"
                        class="text-red-600 hover:text-red-700 text-sm"
                      >
                        取消收藏
                      </button>
                    </div>
                  </div>
                </div>

                <!-- 分页 -->
                <div v-if="favoritesPagination.total_pages > 1" class="flex justify-center mt-8">
                  <div class="flex space-x-2">
                    <button
                      v-for="page in getPaginationRange(favoritesPagination.current_page, favoritesPagination.total_pages)"
                      :key="page"
                      @click="loadMyFavorites(page)"
                      :class="[
                        'px-3 py-2 text-sm font-medium rounded-md',
                        page === favoritesPagination.current_page
                          ? 'bg-primary-600 text-white'
                          : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                      ]"
                    >
                      {{ page }}
                    </button>
                  </div>
                </div>
              </div>

              <!-- 空状态 -->
              <div v-else class="text-center py-12 text-gray-500">
                <StarIcon class="h-12 w-12 mx-auto mb-4" />
                <p>您还没有收藏任何文章</p>
                <router-link
                  to="/articles"
                  class="text-primary-600 hover:text-primary-500 mt-2 inline-block"
                >
                  去发现更多优质文章
                </router-link>
              </div>
            </div>

            <!-- 我的积分 -->
            <div v-else-if="activeTab === 'points'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">我的积分</h3>

              <!-- 积分统计卡片（与下方积分卡一致的浅色渐变） -->
              <div v-if="pointsStats" class="bg-gradient-to-r from-primary-50 to-secondary-50 rounded-lg p-6 border border-primary-100 mb-6">
                <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
                  <div class="flex flex-col items-center space-y-1">
                    <div class="text-3xl md:text-2xl font-bold text-gray-900 tabular-nums tracking-tight">{{ pointsStats.total_points }}</div>
                    <div class="text-xs md:text-sm text-gray-600">总积分</div>
                  </div>
                  <div class="flex flex-col items-center space-y-1">
                    <div class="text-3xl md:text-2xl font-bold text-gray-900 tabular-nums tracking-tight">{{ pointsStats.current_level }}</div>
                    <div class="text-xs md:text-sm text-gray-600">当前等级</div>
                  </div>
                  <div class="flex flex-col items-center space-y-1">
                    <div class="text-3xl md:text-2xl font-bold text-gray-900 tabular-nums tracking-tight">{{ pointsStats.today_points }}</div>
                    <div class="text-xs md:text-sm text-gray-600">今日获得</div>
                  </div>
                  <div class="flex flex-col items-center space-y-1">
                    <div class="text-3xl md:text-2xl font-bold text-gray-900 tabular-nums tracking-tight">#{{ pointsStats.rank || '--' }}</div>
                    <div class="text-xs md:text-sm text-gray-600">积分排名</div>
                  </div>
                </div>

                <!-- 等级进度条 -->
                <div v-if="pointsStats.level_info" class="mt-4 pt-4 border-t border-gray-200">
                  <div class="flex items-center justify-between text-sm mb-2">
                    <span class="text-gray-700">{{ pointsStats.level_info.name }}</span>
                    <span class="font-medium text-primary-600">还需 {{ pointsStats.next_level_points }} 分升级</span>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2">
                    <div
                      class="bg-gradient-to-r from-primary-500 to-secondary-500 h-2 rounded-full transition-all duration-300"
                      :style="{
                        width: `${Math.min(100, (pointsStats.total_points - pointsStats.level_info.min_points) / (pointsStats.level_info.max_points - pointsStats.level_info.min_points) * 100)}%`
                      }"
                    ></div>
                  </div>
                </div>

                <!-- 规则快捷入口 -->
                <div class="mt-4 text-right">
                  <button
                    class="text-primary-600 hover:text-primary-700 text-sm underline"
                    @click="showRulesDrawer = true"
                  >
                    查看全部积分规则
                  </button>
                </div>
              </div>

              <!-- 加载状态 -->
              <div v-else-if="isLoadingPoints" class="text-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                <p class="mt-4 text-gray-500">加载中...</p>
              </div>

              <!-- 选项卡 -->
              <div class="border-b border-gray-200 mb-6">
                <nav class="-mb-px flex space-x-8">
                  <button
                    @click="pointsTab = 'history'"
                    :class="[
                      'py-2 px-1 border-b-2 font-medium text-sm',
                      pointsTab === 'history'
                        ? 'border-primary-500 text-primary-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    ]"
                  >
                    积分记录
                  </button>
                  <button
                    @click="pointsTab = 'rules'"
                    :class="[
                      'py-2 px-1 border-b-2 font-medium text-sm',
                      pointsTab === 'rules'
                        ? 'border-primary-500 text-primary-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    ]"
                  >
                    积分规则
                  </button>
                </nav>
              </div>

              <!-- 积分记录 -->
              <div v-if="pointsTab === 'history'">
                <!-- 筛选条 -->
                <div class="mb-4 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                  <!-- 类型筛选 -->
                  <div class="flex items-center gap-2">
                    <button
                      v-for="t in pointsTypeOptions"
                      :key="t.value"
                      @click="pointsTypeFilter = t.value"
                      :class="[
                        'px-3 py-1.5 text-sm rounded-md border',
                        pointsTypeFilter === t.value
                          ? 'bg-primary-50 text-primary-700 border-primary-200'
                          : 'bg-white text-gray-700 border-gray-200 hover:bg-gray-50'
                      ]"
                    >{{ t.label }}</button>
                  </div>

                  <!-- 时间筛选 -->
                  <div class="flex items-center gap-2">
                    <button
                      v-for="o in pointsTimeOptions"
                      :key="o.value"
                      @click="pointsTimeFilter = o.value"
                      :class="[
                        'px-3 py-1.5 text-sm rounded-md border',
                        pointsTimeFilter === o.value
                          ? 'bg-primary-50 text-primary-700 border-primary-200'
                          : 'bg-white text-gray-700 border-gray-200 hover:bg-gray-50'
                      ]"
                    >{{ o.label }}</button>
                    <div v-if="pointsTimeFilter === 'custom'" class="flex items-center gap-2">
                      <input type="date" v-model="customDateRange.start" class="px-2 py-1 text-sm border border-gray-300 rounded" />
                      <span class="text-gray-400">-</span>
                      <input type="date" v-model="customDateRange.end" class="px-2 py-1 text-sm border border-gray-300 rounded" />
                    </div>
                  </div>

                  <!-- 搜索 -->
                  <div class="relative w-full md:w-64">
                    <input
                      v-model="pointsSearchQuery"
                      type="text"
                      placeholder="搜索描述/来源"
                      class="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-200"
                    />
                    <SearchIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                  </div>
                </div>
                <!-- 加载状态 -->
                <div v-if="isLoadingPointsHistory" class="text-center py-12">
                  <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                  <p class="mt-4 text-gray-500">加载中...</p>
                </div>

                <!-- 积分历史记录（分组） -->
                <div v-else-if="groupedPointsHistory.length > 0">
                  <div class="bg-white border border-gray-200 rounded-lg">
                    <div v-for="group in groupedPointsHistory" :key="group.label" class="">
                      <div class="px-4 py-2 text-xs font-medium text-soft-500 bg-gray-50">{{ group.label }}</div>
                      <div class="divide-y divide-gray-100">
                        <div
                          v-for="record in group.items"
                          :key="record.id"
                          class="p-4 flex items-center justify-between"
                        >
                          <div class="flex-1">
                            <h4 class="text-sm font-medium text-gray-900">
                              {{ record.description }}
                            </h4>
                            <p class="text-xs text-gray-500 mt-1">
                              {{ formatTimeOnly(record.created_at) }}
                            </p>
                          </div>
                          <div class="ml-4">
                            <span
                              :class="[
                                'text-sm font-semibold',
                                record.points > 0 ? 'text-secondary-600' : 'text-red-600'
                              ]"
                            >
                              {{ formatPointsChange(record.points) }}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 分页：按产品要求不再展示，默认仅显示最近 N 条 -->
                </div>

                <!-- 空状态 -->
                <div v-else class="text-center py-12 text-gray-500">
                  <TrophyIcon class="h-12 w-12 mx-auto mb-4" />
                  <p>暂无积分记录</p>
                  <p class="text-sm mt-2">发布文章、点赞、评论等行为可以获得积分</p>
                </div>
              </div>

              <!-- 积分规则 -->
              <div v-else-if="pointsTab === 'rules'">
                <div v-if="pointsRules.length > 0">
                  <!-- Top 3 + 查看全部按钮 -->
                  <div class="bg-white border border-gray-200 rounded-lg divide-y divide-gray-100">
                    <div
                      v-for="rule in topRules"
                      :key="rule.id"
                      class="p-4 flex items-start justify-between"
                    >
                      <div class="pr-4">
                        <h4 class="text-sm font-medium text-gray-900">{{ rule.name }}</h4>
                        <p class="text-xs text-gray-500 mt-1">{{ rule.description }}</p>
                        <p v-if="rule.daily_limit > 0" class="text-xs text-secondary-600 mt-1">每日限制: {{ rule.daily_limit }} 次</p>
                      </div>
                      <div class="ml-4">
                        <span class="text-sm font-semibold text-secondary-600">+{{ rule.points }}</span>
                      </div>
                    </div>
                  </div>
                  <div class="mt-4 flex justify-end">
                    <button class="text-sm text-primary-600 hover:text-primary-700 underline" @click="showRulesDrawer = true">查看全部规则</button>
                  </div>
                </div>

                <!-- 空状态 -->
                <div v-else class="text-center py-12 text-gray-500">
                  <TrophyIcon class="h-12 w-12 mx-auto mb-4" />
                  <p>暂无积分规则</p>
                </div>
              </div>

              <!-- 规则抽屉 -->
              <div v-if="showRulesDrawer" class="fixed inset-0 z-50">
                <div class="absolute inset-0 bg-black/40" @click="showRulesDrawer = false"></div>
                <div class="absolute right-0 top-0 h-full w-full max-w-md bg-white shadow-xl flex flex-col">
                  <div class="px-4 py-3 border-b border-gray-200 flex items-center justify-between">
                    <h4 class="text-base font-medium text-gray-900">积分规则</h4>
                    <button @click="showRulesDrawer = false" class="p-1 rounded hover:bg-gray-100">
                      <XIcon class="h-5 w-5" />
                    </button>
                  </div>
                  <div class="p-4 border-b border-gray-100">
                    <div class="relative">
                      <input
                        v-model="rulesQuery"
                        type="text"
                        placeholder="搜索规则..."
                        class="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-200"
                      />
                      <SearchIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                    </div>
                  </div>
                  <div class="flex-1 overflow-y-auto">
                    <div v-for="group in groupedRules" :key="group.key" class="border-b border-gray-100">
                      <button
                        class="w-full px-4 py-3 flex items-center justify-between hover:bg-gray-50"
                        @click="toggleRuleGroup(group.key)"
                      >
                        <div class="flex items-center gap-2">
                          <ChevronDownIcon :class="['h-4 w-4 transition-transform', expandedRuleGroups.has(group.key) ? 'rotate-180' : '']" />
                          <span class="text-sm font-medium text-gray-900">{{ group.label }}</span>
                          <span class="text-xs text-gray-500">({{ group.items.length }})</span>
                        </div>
                      </button>
                      <div v-if="expandedRuleGroups.has(group.key)" class="divide-y divide-gray-100">
                        <div v-for="rule in group.items" :key="rule.id" class="px-4 py-3 flex items-start justify-between">
                          <div class="pr-4">
                            <h5 class="text-sm font-medium text-gray-900">{{ rule.name }}</h5>
                            <p class="text-xs text-gray-500 mt-1">{{ rule.description }}</p>
                            <p v-if="rule.daily_limit > 0" class="text-xs text-secondary-600 mt-1">每日限制: {{ rule.daily_limit }} 次</p>
                          </div>
                          <div class="ml-4">
                            <span class="text-sm font-semibold text-secondary-600">+{{ rule.points }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 关注（我关注的人） -->
            <div v-else-if="activeTab === 'following'" class="p-6">
              <div class="flex items-center justify-between mb-6">
                <h3 class="text-lg font-medium text-gray-900">关注（{{ followingCount }}）</h3>
              </div>

              <div v-if="followingList.length > 0" class="bg-white border border-gray-200 rounded-lg divide-y divide-gray-100">
                <div
                  v-for="u in followingList"
                  :key="u.id"
                  class="p-4 flex items-center justify-between"
                >
                  <div class="flex items-center">
                    <router-link :to="`/users/${u.username}`" class="flex items-center hover:opacity-90">
                      <UserAvatar :avatar="u.avatar || ''" :name="u.nickname || u.username" :size="40" />
                    </router-link>
                    <div class="ml-3">
                      <router-link :to="`/users/${u.username}`" class="text-sm font-medium text-gray-900 hover:text-primary-600">{{ u.nickname || u.username }}</router-link>
                      <div class="text-xs text-gray-500">@{{ u.username }}</div>
                    </div>
                  </div>
                  <div class="flex items-center gap-3">
                    <span v-if="u.is_mutual_follow" class="px-2 py-0.5 rounded-full text-xs bg-secondary-50 text-secondary-700 border border-secondary-200">互相关注</span>
                    <button
                      class="px-3 py-1.5 text-sm rounded-md border border-gray-200 text-gray-700 hover:bg-gray-50"
                      @click="onUnfollow(u.id)"
                    >取消关注</button>
                  </div>
                </div>
              </div>

              <div v-else class="text-center py-12 text-gray-500">
                <UsersIcon class="h-12 w-12 mx-auto mb-4" />
                <p>还没有关注任何人</p>
              </div>
            </div>

            <!-- 粉丝（关注我的人） -->
            <div v-else-if="activeTab === 'followers'" class="p-6">
              <div class="flex items-center justify-between mb-6">
                <h3 class="text-lg font-medium text-gray-900">粉丝（{{ followersCount }}）</h3>
              </div>

              <div v-if="followersList.length > 0" class="bg-white border border-gray-200 rounded-lg divide-y divide-gray-100">
                <div
                  v-for="u in followersList"
                  :key="u.id"
                  class="p-4 flex items-center justify-between"
                >
                  <div class="flex items-center">
                    <router-link :to="`/users/${u.username}`" class="flex items-center hover:opacity-90">
                      <UserAvatar :avatar="u.avatar || ''" :name="u.nickname || u.username" :size="40" />
                    </router-link>
                    <div class="ml-3">
                      <router-link :to="`/users/${u.username}`" class="text-sm font-medium text-gray-900 hover:text-primary-600">{{ u.nickname || u.username }}</router-link>
                      <div class="text-xs text-gray-500">@{{ u.username }}</div>
                    </div>
                  </div>
                  <div class="flex items-center gap-3">
                    <span v-if="u.is_mutual_follow" class="px-2 py-0.5 rounded-full text-xs bg-secondary-50 text-secondary-700 border border-secondary-200">互相关注</span>
                    <template v-else>
                      <button
                        class="px-3 py-1.5 text-sm rounded-md bg-primary-600 text-white hover:bg-primary-700"
                        @click="onFollow(u.id)"
                      >回关</button>
                    </template>
                  </div>
                </div>
              </div>

              <div v-else class="text-center py-12 text-gray-500">
                <UsersIcon class="h-12 w-12 mx-auto mb-4" />
                <p>还没有粉丝</p>
              </div>
            </div>

            <!-- 互相关注 -->
            <div v-else-if="activeTab === 'mutual'" class="p-6">
              <div class="flex items-center justify-between mb-6">
                <h3 class="text-lg font-medium text-gray-900">互相关注（{{ mutualFollowsList.length }}）</h3>
              </div>

              <div v-if="mutualFollowsList.length > 0" class="bg-white border border-gray-200 rounded-lg divide-y divide-gray-100">
                <div
                  v-for="u in mutualFollowsList"
                  :key="u.id"
                  class="p-4 flex items-center justify-between"
                >
                  <div class="flex items-center">
                    <router-link :to="`/users/${u.username}`" class="flex items-center hover:opacity-90">
                      <UserAvatar :avatar="u.avatar || ''" :name="u.nickname || u.username" :size="40" />
                    </router-link>
                    <div class="ml-3">
                      <router-link :to="`/users/${u.username}`" class="text-sm font-medium text-gray-900 hover:text-primary-600">{{ u.nickname || u.username }}</router-link>
                      <div class="text-xs text-gray-500">@{{ u.username }}</div>
                    </div>
                  </div>
                  <div class="flex items-center gap-3">
                    <span class="px-2 py-0.5 rounded-full text-xs bg-secondary-50 text-secondary-700 border border-secondary-200">互相关注</span>
                    <button
                      class="px-3 py-1.5 text-sm rounded-md border border-gray-200 text-gray-700 hover:bg-gray-50"
                      @click="onUnfollow(u.id)"
                    >取消关注</button>
                  </div>
                </div>
              </div>

              <div v-else class="text-center py-12 text-gray-500">
                <UsersIcon class="h-12 w-12 mx-auto mb-4" />
                <p>还没有互相关注</p>
              </div>
            </div>

            <!-- 设置 -->
            <div v-else-if="activeTab === 'settings'" class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">设置</h3>

              <!-- 修改密码 -->
              <div class="bg-gray-50 rounded-lg p-6 mb-6">
                <h4 class="text-md font-medium text-gray-900 mb-4">修改密码</h4>
                <form @submit.prevent="changePassword" class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      当前密码
                    </label>
                    <input
                      v-model="passwordForm.old_password"
                      type="password"
                      required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                      placeholder="请输入当前密码"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      新密码
                    </label>
                    <input
                      v-model="passwordForm.new_password"
                      type="password"
                      required
                      minlength="6"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                      placeholder="请输入新密码（至少6位）"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      确认新密码
                    </label>
                    <input
                      v-model="passwordForm.confirm_password"
                      type="password"
                      required
                      minlength="6"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
                      placeholder="请再次输入新密码"
                    />
                  </div>

                  <!-- 密码强度提示 -->
                  <div v-if="passwordForm.new_password" class="text-xs text-gray-500">
                    <p>密码要求：</p>
                    <ul class="list-disc list-inside mt-1 space-y-1">
                      <li :class="passwordForm.new_password.length >= 6 ? 'text-green-600' : 'text-red-500'">
                        至少6个字符
                      </li>
                      <li :class="passwordForm.new_password !== passwordForm.old_password ? 'text-green-600' : 'text-red-500'">
                        与当前密码不同
                      </li>
                      <li :class="passwordForm.new_password === passwordForm.confirm_password ? 'text-green-600' : 'text-red-500'">
                        两次密码输入一致
                      </li>
                    </ul>
                  </div>

                  <div class="flex justify-end">
                    <button
                      type="submit"
                      :disabled="isChangingPassword || !isPasswordFormValid"
                      class="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {{ isChangingPassword ? '修改中...' : '修改密码' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- 账户安全提示 -->
              <div class="bg-secondary-50 rounded-lg p-4">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg class="h-5 w-5 text-secondary-400" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm text-secondary-700">
                      <strong>安全提示：</strong>
                    </p>
                    <div class="mt-2 text-sm text-secondary-600">
                      <ul class="list-disc list-inside space-y-1">
                        <li>定期更换密码可以提高账户安全性</li>
                        <li>请使用包含字母、数字和特殊字符的强密码</li>
                        <li>不要在其他网站使用相同的密码</li>
                        <li>如发现账户异常活动，请立即修改密码</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 其他标签页内容会在这里添加... -->
            <div v-else class="p-6">
              <h3 class="text-lg font-medium text-gray-900 mb-6">{{ activeTab }}</h3>
              <p class="text-gray-500">功能开发中...</p>
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
    </PageContainer>
  </AppLayout>
</template>

<script setup lang="ts">
/* eslint-disable @typescript-eslint/no-unused-vars */
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  UserIcon,
  FileTextIcon,
  StarIcon,
  SettingsIcon,
  CameraIcon,
  UsersIcon,
  HeartIcon,
  BellIcon,
  TrophyIcon
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { AppLayout, PageContainer } from '@/components/layout'
import AvatarModal from '../components/AvatarModal.vue'
import UserPointsDisplay from '../components/UserPointsDisplay.vue'
import { useToast } from '../composables/useToast'
import { useUserDataSync } from '../composables/useUserDataSync'
import type { ImageUploadResponse } from '../api/types'
import UserAvatar from '@/components/UserAvatar.vue'
import { FavoriteApi, FavoriteUtils, type Favorite } from '@/api/favorite'
import { UserApi } from '@/api/user'
import type { ChangePasswordRequest } from '@/api/types'
import { PointsAPI, type PointsTransaction, type PointsRule, type PointsStats } from '@/api/points'
import { SearchIcon, XIcon, ChevronDownIcon } from 'lucide-vue-next'

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

// 使用同步的数据 - 直接使用userDataSync返回的computed属性
const myArticles = userDataSync.articles
const isLoadingArticles = userDataSync.isLoading
const followingList = userDataSync.followingList
const followersList = userDataSync.followersList
const mutualFollowsList = userDataSync.mutualFollowsList
const followingCount = computed(() => userDataSync.stats.value.following_count)
const followersCount = computed(() => userDataSync.stats.value.followers_count)
const articlesCount = computed(() => userDataSync.stats.value.article_count)
const articlesError = ref('')

// 收藏相关状态
const myFavorites = ref<Favorite[]>([])
const isLoadingFavorites = ref(false)
const favoritesError = ref('')
const favoritesPagination = ref({
  total: 0,
  current_page: 1,
  per_page: 20,
  total_pages: 1
})

// 密码修改相关状态
const isChangingPassword = ref(false)
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// 积分相关状态
const pointsStats = ref<PointsStats | null>(null)
const pointsHistory = ref<PointsTransaction[]>([])
const pointsRules = ref<PointsRule[]>([])
const isLoadingPoints = ref(false)
const isLoadingPointsHistory = ref(false)
const pointsError = ref('')
const pointsPagination = ref({
  page: 1,
  limit: 8, // 仅展示最近 8 条
  total: 0,
  total_pages: 1
})
const pointsTab = ref('history')

// 关注操作
const onFollow = async (userId: number) => {
  try {
    await userDataSync.followUser(userId)
  } catch (e) {
    console.error('关注失败', e)
  }
}

const onUnfollow = async (userId: number) => {
  try {
    await userDataSync.unfollowUser(userId)
  } catch (e) {
    console.error('取消关注失败', e)
  }
}

// 记录筛选状态
const pointsTypeOptions = [
  { label: '全部', value: 'all' },
  { label: '获得', value: 'income' },
  { label: '扣减', value: 'expense' }
]
const pointsTypeFilter = ref<'all' | 'income' | 'expense'>('all')

const pointsTimeOptions = [
  { label: '7天', value: '7d' },
  { label: '30天', value: '30d' },
  { label: '全部', value: 'all' },
  { label: '自定义', value: 'custom' }
]
const pointsTimeFilter = ref<'7d' | '30d' | 'all' | 'custom'>('all')
const customDateRange = reactive<{ start: string; end: string }>({ start: '', end: '' })
const pointsSearchQuery = ref('')

// 记录过滤与分组
const filteredPointsHistory = computed(() => {
  let list = [...pointsHistory.value]
  // 类型
  if (pointsTypeFilter.value === 'income') list = list.filter(i => i.points > 0)
  if (pointsTypeFilter.value === 'expense') list = list.filter(i => i.points < 0)
  // 时间
  const now = new Date()
  let startDate: Date | null = null
  if (pointsTimeFilter.value === '7d') {
    startDate = new Date(now)
    startDate.setDate(now.getDate() - 7)
  } else if (pointsTimeFilter.value === '30d') {
    startDate = new Date(now)
    startDate.setDate(now.getDate() - 30)
  } else if (pointsTimeFilter.value === 'custom' && customDateRange.start) {
    startDate = new Date(customDateRange.start + 'T00:00:00')
  }
  const endDate = pointsTimeFilter.value === 'custom' && customDateRange.end
    ? new Date(customDateRange.end + 'T23:59:59')
    : null

  if (startDate) list = list.filter(i => new Date(i.created_at) >= startDate!)
  if (endDate) list = list.filter(i => new Date(i.created_at) <= endDate!)

  // 搜索
  const q = pointsSearchQuery.value.trim().toLowerCase()
  if (q) {
    list = list.filter(i =>
      (i.description || '').toLowerCase().includes(q) ||
      (i.source_type || '').toLowerCase().includes(q)
    )
  }
  return list
})

const groupedPointsHistory = computed(() => {
  const groups: { label: string; items: PointsTransaction[]; key: string }[] = []
  const map = new Map<string, PointsTransaction[]>()
  for (const item of filteredPointsHistory.value) {
    const key = getDateKey(item.created_at)
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(item)
  }
  // 日期 key 倒序
  const keys = Array.from(map.keys()).sort((a, b) => +new Date(b) - +new Date(a))
  for (const key of keys) {
    const items = map.get(key) || []
    items.sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at))
    groups.push({ key, label: getDateGroupLabel(key), items })
  }
  return groups
})

const getDateKey = (iso: string) => {
  const d = new Date(iso)
  const y = d.getFullYear()
  const m = (d.getMonth() + 1).toString().padStart(2, '0')
  const day = d.getDate().toString().padStart(2, '0')
  return `${y}-${m}-${day}`
}

const getDateGroupLabel = (dateKey: string) => {
  const d = new Date(dateKey)
  const today = new Date()
  const ytd = new Date()
  ytd.setDate(today.getDate() - 1)
  const isSameDay = (d1: Date, d2: Date) => d1.toDateString() === d2.toDateString()
  if (isSameDay(d, today)) return '今天'
  if (isSameDay(d, ytd)) return '昨天'
  const y = d.getFullYear()
  const m = (d.getMonth() + 1).toString().padStart(2, '0')
  const day = d.getDate().toString().padStart(2, '0')
  const curY = today.getFullYear()
  return y === curY ? `${m}-${day}` : `${y}-${m}-${day}`
}

const formatTimeOnly = (iso: string) => {
  const d = new Date(iso)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 规则抽屉与分组
const showRulesDrawer = ref(false)
const rulesQuery = ref('')

const filteredRules = computed(() => {
  const q = rulesQuery.value.trim().toLowerCase()
  let list = pointsRules.value.filter(r => r.status === 1)
  if (q) {
    list = list.filter(r =>
      r.name.toLowerCase().includes(q) ||
      (r.description || '').toLowerCase().includes(q) ||
      (r.action || '').toLowerCase().includes(q)
    )
  }
  return list
})

const groupedRules = computed(() => {
  const groups: Record<string, PointsRule[]> = {}
  for (const r of filteredRules.value) {
    const key = detectRuleGroupKey(r)
    if (!groups[key]) groups[key] = []
    groups[key].push(r)
  }
  const order = ['creation', 'interaction', 'system', 'other']
  const labelMap: Record<string, string> = {
    creation: '创作相关',
    interaction: '互动相关',
    system: '系统奖励',
    other: '其他规则'
  }
  return order
    .filter(k => groups[k]?.length)
    .map(k => ({ key: k, label: labelMap[k], items: groups[k].sort((a, b) => b.points - a.points) }))
})

const detectRuleGroupKey = (r: PointsRule): 'creation' | 'interaction' | 'system' | 'other' => {
  const a = (r.action || '').toLowerCase()
  if (/(article|post|publish|create|write|upload)/.test(a)) return 'creation'
  if (/(like|comment|follow|share|favorite|collect)/.test(a)) return 'interaction'
  if (/(system|daily|login|sign|register)/.test(a)) return 'system'
  return 'other'
}

const expandedRuleGroups = reactive<Set<string>>(new Set(['creation', 'interaction', 'system', 'other']))
const toggleRuleGroup = (key: string) => {
  if (expandedRuleGroups.has(key)) expandedRuleGroups.delete(key)
  else expandedRuleGroups.add(key)
}

// Top 3 规则（按积分倒序）
const topRules = computed(() => filteredRules.value.slice().sort((a, b) => b.points - a.points).slice(0, 3))

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

// 密码表单验证
const isPasswordFormValid = computed(() => {
  return passwordForm.old_password.length > 0 &&
         passwordForm.new_password.length >= 6 &&
         passwordForm.new_password === passwordForm.confirm_password &&
         passwordForm.new_password !== passwordForm.old_password
})

// 菜单项
const menuItems = [
  { key: 'profile', label: '个人信息', icon: UserIcon },
  { key: 'articles', label: '我的文章', icon: FileTextIcon },
  { key: 'favorites', label: '我的收藏', icon: StarIcon },
  { key: 'following', label: '关注', icon: UsersIcon },
  { key: 'followers', label: '粉丝', icon: UsersIcon },
  { key: 'mutual', label: '互相关注', icon: UsersIcon },
  { key: 'points', label: '我的积分', icon: TrophyIcon },
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
      birthday: profileForm.birthday ? `${profileForm.birthday}T00:00:00Z` : undefined,
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

// 加载我的收藏
const loadMyFavorites = async (page: number = 1) => {
  try {
    isLoadingFavorites.value = true
    favoritesError.value = ''

    const response = await FavoriteApi.getUserFavorites({
      page,
      size: favoritesPagination.value.per_page
    })

    myFavorites.value = response.data.favorites || []
    favoritesPagination.value = response.data.pagination
  } catch (error: any) {
    favoritesError.value = error.message || '加载收藏失败'
    console.error('加载收藏失败:', error)
  } finally {
    isLoadingFavorites.value = false
  }
}

// 取消收藏
const removeFavorite = async (favorite: Favorite) => {
  if (!confirm('确定要取消收藏这篇文章吗？')) {
    return
  }

  try {
    await FavoriteApi.deleteFavorite(favorite.id)
    toast.success('已取消收藏')

    // 从列表中移除
    myFavorites.value = myFavorites.value.filter(f => f.id !== favorite.id)

    // 更新分页信息
    favoritesPagination.value.total--
    if (myFavorites.value.length === 0 && favoritesPagination.value.current_page > 1) {
      // 如果当前页没有数据了，返回上一页
      loadMyFavorites(favoritesPagination.value.current_page - 1)
    }
  } catch (error: any) {
    toast.error(error.message || '取消收藏失败')
    console.error('取消收藏失败:', error)
  }
}

// 格式化收藏时间
const formatFavoriteTime = (timeString: string) => {
  return FavoriteUtils.formatFavoriteTime(timeString)
}

// 分页范围计算
const getPaginationRange = (currentPage: number, totalPages: number, maxVisible: number = 5) => {
  if (totalPages <= maxVisible) {
    return Array.from({ length: totalPages }, (_, i) => i + 1)
  }

  const half = Math.floor(maxVisible / 2)
  let start = Math.max(1, currentPage - half)
  let end = Math.min(totalPages, start + maxVisible - 1)

  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1)
  }

  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
}

// 修改密码
const changePassword = async () => {
  try {
    isChangingPassword.value = true

    const data: ChangePasswordRequest = {
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    }

    await UserApi.changePassword(data)
    toast.success('密码修改成功')

    // 清空表单
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (error: any) {
    toast.error(error.message || '密码修改失败')
    console.error('密码修改失败:', error)
  } finally {
    isChangingPassword.value = false
  }
}

// 加载积分统计
const loadPointsStats = async () => {
  try {
    isLoadingPoints.value = true
    pointsError.value = ''

    const response = await PointsAPI.getPointsStats()
    pointsStats.value = response.data
  } catch (error: any) {
    pointsError.value = error.message || '加载积分信息失败'
    console.error('加载积分统计失败:', error)
  } finally {
    isLoadingPoints.value = false
  }
}

// 加载积分历史记录
const loadPointsHistory = async (page: number = 1) => {
  try {
    isLoadingPointsHistory.value = true
    pointsError.value = ''

    const response = await PointsAPI.getPointsHistory({
      page,
      limit: pointsPagination.value.limit
    })

    pointsHistory.value = response.data.transactions
    pointsPagination.value = response.data.pagination
  } catch (error: any) {
    pointsError.value = error.message || '加载积分历史失败'
    console.error('加载积分历史失败:', error)
  } finally {
    isLoadingPointsHistory.value = false
  }
}

// 加载积分规则
const loadPointsRules = async () => {
  try {
    const response = await PointsAPI.getPointsRules()
    pointsRules.value = response.data
  } catch (error: any) {
    console.error('加载积分规则失败:', error)
  }
}

// 格式化积分变化
const formatPointsChange = (points: number) => {
  return points > 0 ? `+${points}` : points.toString()
}

// 格式化积分时间
const formatPointsTime = (timeString: string) => {
  const time = new Date(timeString)
  return time.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 监听activeTab变化，根据标签加载不同数据
watch(activeTab, (newTab) => {
  if (newTab === 'articles') {
    loadMyArticles()
  } else if (newTab === 'favorites') {
    loadMyFavorites()
  } else if (newTab === 'points') {
    loadPointsStats()
    loadPointsHistory()
    loadPointsRules()
  } else if (newTab === 'following') {
    userDataSync.loadFollowingList()
  } else if (newTab === 'followers') {
    userDataSync.loadFollowersList()
  } else if (newTab === 'mutual') {
    userDataSync.loadMutualFollowsList()
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
  if (tab && ['profile', 'articles', 'favorites', 'points', 'settings'].includes(tab)) {
    activeTab.value = tab
  }

  initUserInfo()
  // 初始化用户数据同步
  userDataSync.initUserData()

  // 避免首次进入时用户信息尚未注入导致关注列表不加载
  watch(
    () => authStore.user?.username,
    (val) => {
      if (val) {
        userDataSync.initUserData()
        if (activeTab.value === 'following') userDataSync.loadFollowingList()
        if (activeTab.value === 'followers') userDataSync.loadFollowersList()
        if (activeTab.value === 'mutual') userDataSync.loadMutualFollowsList()
      }
    },
    { immediate: true }
  )
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
