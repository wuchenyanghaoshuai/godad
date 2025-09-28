<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <h1 class="text-xl font-semibold text-gray-900">举报中心</h1>
          <router-link to="/management-dashboard" class="text-sm text-gray-600 hover:text-blue-600">后台管理</router-link>
          <router-link to="/community-management" class="text-sm text-gray-600 hover:text-blue-600">社区管理</router-link>
          <router-link to="/report-center" class="text-sm text-blue-600">举报中心</router-link>
        </div>
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-600">欢迎，{{ authStore.user?.nickname }}</span>
          <button @click="logout" class="text-sm text-red-600 hover:text-red-700">退出登录</button>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200">
        <!-- 过滤栏 -->
        <div class="p-4 border-b border-gray-200 flex flex-wrap items-center gap-3">
          <input v-model="filters.keyword" placeholder="搜索举报内容/理由" class="px-3 py-2 border rounded w-64" />
          <select v-model="filters.status" class="px-3 py-2 border rounded">
            <option value="">全部状态</option>
            <option value="pending">待处理</option>
            <option value="reviewed">已处理</option>
            <option value="rejected">已驳回</option>
          </select>
          <select v-model="filters.type" class="px-3 py-2 border rounded">
            <option value="">全部类型</option>
            <option value="article">文章</option>
            <option value="forum_post">社区帖子</option>
          </select>
          <button @click="loadReports" class="px-4 py-2 bg-blue-600 text-white rounded">查询</button>
        </div>

        <!-- 卡片列表 -->
        <div class="p-4">
          <div v-if="reports.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="r in reports"
              :key="r.id"
              class="border rounded-lg p-4 shadow-sm hover:shadow-md transition cursor-pointer bg-white"
              @click="openDetail(r)"
            >
              <div class="flex items-center justify-between mb-2">
                <div class="text-sm text-gray-500">
                  <span class="px-2 py-0.5 rounded bg-gray-100 text-gray-700">{{ typeLabel(r.target_type) }}</span>
                  <router-link
                    :to="targetPath(r)"
                    @click.stop
                    class="ml-2 text-blue-600 hover:underline"
                    :title="'前往查看具体内容'"
                  >
                    {{ r.target_title || ('#' + r.target_id) }}
                  </router-link>
                </div>
                <span :class="statusBadgeClass(r.status)">{{ statusLabel(r.status) }}</span>
              </div>
              <div class="mb-2">
                <p class="text-gray-900 font-medium truncate" :title="r.reason">{{ r.reason }}</p>
                <p class="text-gray-600 text-sm line-clamp-2" :title="r.description">{{ r.description || '（无补充说明）' }}</p>
              </div>
              <div class="text-xs text-gray-400 flex items-center justify-between">
                <span>创建时间：{{ formatDate(r.created_at) }}</span>
                <span class="flex items-center gap-2">
                  举报人：
                  <UserAvatar :avatar="r.reporter_avatar" :name="(r.reporter_nickname || r.reporter_username || '#')" :size="20" />
                  <template v-if="r.reporter_username && r.reporter_username.trim() !== ''">
                    <router-link
                      :to="{ name: 'UserProfile', params: { username: r.reporter_username } }"
                      @click.stop
                      class="text-blue-600 hover:underline"
                    >
                      {{ r.reporter_nickname && r.reporter_nickname.trim() !== '' ? r.reporter_nickname : r.reporter_username }}
                    </router-link>
                    <span class="text-gray-400">#{{ r.reporter_id }}</span>
                  </template>
                  <template v-else>
                    <span class="text-gray-700">#{{ r.reporter_id }}</span>
                  </template>
                </span>
              </div>
            </div>
          </div>
          <div v-else class="py-12 text-center text-gray-500">暂无举报</div>
        </div>
      </div>
    </div>

    <!-- 详情弹窗 -->
    <div v-if="showDetail" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div class="bg-white w-full max-w-2xl rounded-lg shadow-lg">
        <div class="px-6 py-4 border-b flex items-center justify-between">
          <h3 class="text-lg font-semibold">举报详情 #{{ selectedReport?.id }}</h3>
          <button class="text-gray-500 hover:text-gray-700" @click="closeDetail">✕</button>
        </div>
        <div class="p-6 space-y-3">
          <div class="flex items-center gap-2 text-sm text-gray-600">
            <span class="px-2 py-0.5 rounded bg-gray-100 text-gray-700">{{ typeLabel(selectedReport?.target_type) }}</span>
            <span class="text-gray-400 flex items-center gap-1">
              目标：
              <router-link
                v-if="selectedReport"
                :to="targetPath(selectedReport)"
                @click.stop
                class="text-blue-600 hover:underline"
              >
                {{ selectedReport?.target_title || ('#' + selectedReport?.target_id) }}
              </router-link>
            </span>
            <span :class="['ml-auto', statusBadgeClass(selectedReport?.status)]">{{ statusLabel(selectedReport?.status) }}</span>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">举报理由</div>
            <div class="text-gray-900">{{ selectedReport?.reason }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">补充说明</div>
            <div class="text-gray-700 whitespace-pre-wrap">{{ selectedReport?.description || '（无）' }}</div>
          </div>
          <div v-if="selectedReport?.evidence">
            <div class="text-xs text-gray-500 mb-1">证据链接</div>
            <a :href="selectedReport?.evidence" target="_blank" class="text-blue-600 hover:underline">{{ selectedReport?.evidence }}</a>
          </div>
          <div class="grid grid-cols-2 gap-4 text-sm text-gray-500">
            <div class="flex items-center gap-2">
              举报人：
              <UserAvatar :avatar="selectedReport?.reporter_avatar" :name="(selectedReport?.reporter_nickname || selectedReport?.reporter_username || '#')" :size="24" />
              <template v-if="selectedReport?.reporter_username && selectedReport?.reporter_username.trim() !== ''">
                <router-link
                  :to="{ name: 'UserProfile', params: { username: selectedReport?.reporter_username } }"
                  @click.stop
                  class="text-blue-600 hover:underline"
                >
                  {{ selectedReport?.reporter_nickname && selectedReport?.reporter_nickname.trim() !== '' ? selectedReport?.reporter_nickname : selectedReport?.reporter_username }}
                </router-link>
                <span class="text-gray-400">#{{ selectedReport?.reporter_id }}</span>
              </template>
              <template v-else>
                <span class="text-gray-700">#{{ selectedReport?.reporter_id }}</span>
              </template>
            </div>
            <div>创建时间：{{ formatDate(selectedReport?.created_at) }}</div>
            <div class="flex items-center gap-2">
              处理人：
              <template v-if="selectedReport?.handled_by">
                <UserAvatar :avatar="selectedReport?.handler_avatar" :name="(selectedReport?.handler_nickname || selectedReport?.handler_username || '#')" :size="24" />
                <template v-if="selectedReport?.handler_username && selectedReport?.handler_username.trim() !== ''">
                  <router-link
                    :to="{ name: 'UserProfile', params: { username: selectedReport?.handler_username } }"
                    @click.stop
                    class="text-blue-600 hover:underline"
                  >
                    {{ selectedReport?.handler_nickname && selectedReport?.handler_nickname.trim() !== '' ? selectedReport?.handler_nickname : selectedReport?.handler_username }}
                  </router-link>
                  <span class="text-gray-400">#{{ selectedReport?.handled_by }}</span>
                </template>
                <template v-else>
                  <span class="text-gray-700">#{{ selectedReport?.handled_by }}</span>
                </template>
              </template>
              <template v-else>
                —
              </template>
            </div>
            <div>更新时间：{{ formatDate(selectedReport?.updated_at) }}</div>
          </div>

          <div class="mt-2">
            <div class="text-xs text-gray-500 mb-1">处理备注</div>
            <textarea v-model="handledNote" class="w-full border rounded px-3 py-2 text-sm" rows="3" placeholder="可填写处理说明（可选）"></textarea>
          </div>
        </div>
        <div class="px-6 py-4 border-t flex items-center justify-end gap-3">
          <button @click="closeDetail" class="px-4 py-2 rounded border">取消</button>
          <button :disabled="actionLoading" @click="submitAction('rejected')" class="px-4 py-2 rounded bg-gray-600 text-white disabled:opacity-50">驳回</button>
          <button :disabled="actionLoading" @click="submitAction('reviewed')" class="px-4 py-2 rounded bg-green-600 text-white disabled:opacity-50">通过</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AdminApi from '@/api/admin'
import { UserApi } from '@/api/user'
import { useToast } from '@/composables/useToast'
import UserAvatar from '@/components/UserAvatar.vue'
import { ArticleApi } from '@/api/article'
import { ForumApi } from '@/api/forum'

const authStore = useAuthStore()
const router = useRouter()
const reports = ref<any[]>([])
const filters = ref({ keyword: '', status: '', type: '' })
const showDetail = ref(false)
const selectedReport = ref<any | null>(null)
const actionLoading = ref(false)
const handledNote = ref('')
const { showToast } = useToast()

const formatDate = (d?: string) => d ? new Date(d).toLocaleString() : '—'
const typeLabel = (t?: string) => t === 'article' ? '文章' : t === 'forum_post' ? '社区帖子' : (t || '未知')
const statusLabel = (s?: string) => s === 'pending' ? '待处理' : s === 'reviewed' ? '已处理' : s === 'rejected' ? '已驳回' : (s || '未知')
const statusBadgeClass = (s?: string) => {
  switch (s) {
    case 'pending': return 'px-2 py-0.5 rounded text-xs bg-yellow-100 text-yellow-800'
    case 'reviewed': return 'px-2 py-0.5 rounded text-xs bg-green-100 text-green-800'
    case 'rejected': return 'px-2 py-0.5 rounded text-xs bg-gray-100 text-gray-700'
    default: return 'px-2 py-0.5 rounded text-xs bg-gray-100 text-gray-700'
  }
}

// 目标详情路径
const targetPath = (r: any) => {
  if (!r) return '/'
  if (r.target_type === 'article') return `/articles/${r.target_id}`
  if (r.target_type === 'forum_post') return `/community/posts/${r.target_id}`
  return '/'
}

// 已不再直接使用 displayUser，改为在模板中分别处理名称与ID，避免重复显示 #ID

const loadReports = async () => {
  try {
    const res = await AdminApi.getReports({
      page: 1,
      size: 20,
      keyword: filters.value.keyword || undefined,
      status: filters.value.status || undefined,
      target_type: filters.value.type || undefined
    })
    reports.value = res.data.items || []
    await Promise.all([
      enrichReportUsers(),
      enrichTargetDetails(),
    ])
    if (import.meta.env.DEV) {
      // 调试输出举报人/处理人信息
      console.log('[ReportCenter] loaded reports:', reports.value.map((r: any) => ({
        id: r.id,
        reporter_id: r.reporter_id,
        reporter_username: r.reporter_username,
        reporter_nickname: r.reporter_nickname,
        reporter_avatar: r.reporter_avatar,
        handled_by: r.handled_by,
        handler_username: r.handler_username,
        handler_nickname: r.handler_nickname,
        target_type: r.target_type,
        target_id: r.target_id,
        target_title: r.target_title,
      })))
    }
  } catch (error) {
    console.error('Failed to load reports:', error)
    reports.value = []
  }
}

const openDetail = (r: any) => {
  selectedReport.value = r
  handledNote.value = r?.handled_note || ''
  showDetail.value = true
  if (import.meta.env.DEV) {
    console.log('[ReportCenter] open detail:', {
      id: r?.id,
      reporter_id: r?.reporter_id,
      reporter_username: r?.reporter_username,
      reporter_nickname: r?.reporter_nickname,
      handled_by: r?.handled_by,
      handler_username: r?.handler_username,
      handler_nickname: r?.handler_nickname,
    })
  }
  ensureReportTargetTitle(selectedReport.value)
}

// 为缺失用户名/昵称的举报项补全用户信息
const enrichReportUsers = async () => {
  try {
    const missingIds = new Set<number>()
    reports.value.forEach((r: any) => {
      if (!r.reporter_username && r.reporter_id) missingIds.add(Number(r.reporter_id))
      if (r.handled_by && !r.handler_username) missingIds.add(Number(r.handled_by))
    })
    if (missingIds.size === 0) return

    const idToUser: Record<number, any> = {}
    await Promise.all(
      Array.from(missingIds).map(async (uid) => {
        try {
          const u = await UserApi.getPublicInfo(uid)
          idToUser[uid] = u.data
        } catch {}
      })
    )

    reports.value = reports.value.map((r: any) => {
      const rr = { ...r }
      if (!rr.reporter_username && rr.reporter_id && idToUser[rr.reporter_id]) {
        const u = idToUser[rr.reporter_id]
        rr.reporter_username = u.username
        rr.reporter_nickname = u.nickname
        rr.reporter_avatar = u.avatar
      }
      if (rr.handled_by && !rr.handler_username && idToUser[rr.handled_by]) {
        const u = idToUser[rr.handled_by]
        rr.handler_username = u.username
        rr.handler_nickname = u.nickname
        rr.handler_avatar = u.avatar
      }
      return rr
    })

    if (import.meta.env.DEV) {
      console.log('[ReportCenter] enriched users:', idToUser)
    }
  } catch {}
}

// 为缺失的目标标题补全（文章/帖子）
const enrichTargetDetails = async () => {
  const articleIds = new Set<number>()
  const postIds = new Set<number>()
  reports.value.forEach((r: any) => {
    if (!r?.target_title) {
      if (r?.target_type === 'article') articleIds.add(Number(r.target_id))
      else if (r?.target_type === 'forum_post') postIds.add(Number(r.target_id))
    }
  })
  const tasks: Promise<any>[] = []
  articleIds.forEach((id) => tasks.push(
    ArticleApi.getArticleDetail(id).then(resp => {
      const title = resp?.data?.title
      if (!title) return
      reports.value = reports.value.map((r: any) => (
        r.target_type === 'article' && Number(r.target_id) === id ? { ...r, target_title: title } : r
      ))
    }).catch(() => {})
  ))
  postIds.forEach((id) => tasks.push(
    ForumApi.getPost(id).then(resp => {
      const title = resp?.data?.title
      if (!title) return
      reports.value = reports.value.map((r: any) => (
        r.target_type === 'forum_post' && Number(r.target_id) === id ? { ...r, target_title: title } : r
      ))
    }).catch(() => {})
  ))
  await Promise.all(tasks)
}

const ensureReportTargetTitle = async (r: any) => {
  try {
    if (!r || r.target_title) return
    const id = Number(r.target_id)
    if (r.target_type === 'article') {
      const resp = await ArticleApi.getArticleDetail(id)
      r.target_title = resp?.data?.title || ''
    } else if (r.target_type === 'forum_post') {
      const resp = await ForumApi.getPost(id)
      r.target_title = resp?.data?.title || ''
    }
  } catch {}
}

const closeDetail = () => {
  showDetail.value = false
  selectedReport.value = null
  handledNote.value = ''
}

const submitAction = async (status: 'reviewed' | 'rejected') => {
  if (!selectedReport.value) return
  try {
    actionLoading.value = true
    await AdminApi.updateReportStatus(selectedReport.value.id, status, handledNote.value || undefined)
    showToast(status === 'reviewed' ? '已标记为通过' : '已驳回', 'success')
    await loadReports()
    // 更新本地选中对象状态
    selectedReport.value.status = status
    selectedReport.value.handled_note = handledNote.value
    closeDetail()
  } catch (e) {
    showToast('操作失败，请稍后重试', 'error')
  } finally {
    actionLoading.value = false
  }
}

// 页面加载时检查权限并加载数据
const initPage = async () => {
  // 检查管理员权限
  if (!authStore.user || authStore.user.role !== 2) {
    router.push('/')
    return
  }

  // 加载举报数据
  await loadReports()
}

initPage()

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
/* 2行截断（已在上方使用 line-clamp-2）若未启用插件，这里兼容实现 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
