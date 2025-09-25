<template>
  <AppLayout
    :header-config="{
      showSearch: true,
      showCreateButton: false,
      showNotifications: true,
      showUserPoints: true,
      showNavigation: true,
      showUserMenu: true
    }"
    :show-footer="true"
    background-class="bg-gray-50"
  >
    <PageContainer background="gray" padding="lg">
      <div class="max-w-5xl mx-auto">
        <!-- 标题与操作栏 -->
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <MegaphoneIcon class="h-6 w-6 text-primary-600" />
            <h1 class="text-xl font-semibold text-gray-900">公告中心</h1>
          </div>
        </div>

        <!-- 筛选与搜索 -->
        <div class="bg-white rounded-lg border border-gray-200 p-4 mb-6">
          <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
            <!-- 类型筛选 -->
            <div class="flex items-center gap-2">
              <button
                v-for="opt in typeOptions"
                :key="opt.value"
                @click="activeType = opt.value"
                :class="[
                  'px-3 py-1.5 text-sm rounded-md border',
                  activeType === opt.value
                    ? 'bg-primary-50 text-primary-700 border-primary-200'
                    : 'bg-white text-gray-700 border-gray-200 hover:bg-gray-50'
                ]"
              >
                {{ opt.label }}
              </button>
            </div>

            <!-- 搜索与仅进行中 -->
            <div class="flex items-center gap-3 w-full lg:w-auto">
              <div class="relative w-full lg:w-72">
                <input
                  v-model="q"
                  type="text"
                  placeholder="搜索标题/内容"
                  class="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-200"
                />
                <SearchIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              </div>
              <label class="inline-flex items-center gap-2 text-sm text-gray-700 select-none">
                <input type="checkbox" v-model="onlyActive" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                仅进行中
              </label>
            </div>
          </div>
        </div>

        <!-- 公告列表 -->
        <div v-if="filteredAnnouncements.length > 0" class="space-y-4">
          <div
            v-for="a in filteredAnnouncements"
            :key="a.id"
            class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-sm transition-shadow"
          >
            <div class="flex items-start justify-between">
              <div class="pr-4">
                <div class="flex items-center gap-2 mb-1">
                  <span
                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="typeBadgeClass(a.type)"
                  >
                    {{ typeLabel(a.type) }}
                  </span>
                  <span v-if="a.pinned" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-secondary-50 text-secondary-700 border border-secondary-200">
                    <PinIcon class="h-3.5 w-3.5 mr-1" /> 置顶
                  </span>
                  <span v-if="a.require_ack" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-red-50 text-red-700 border border-red-200">
                    <AlertTriangleIcon class="h-3.5 w-3.5 mr-1" /> 重要
                  </span>
                </div>
                <h2 class="text-base font-semibold text-gray-900 mb-1">
                  {{ a.title }}
                </h2>
                <p class="text-sm text-gray-600 line-clamp-2">{{ a.summary }}</p>
                <div class="mt-3 flex items-center gap-3 text-xs text-gray-500">
                  <span class="inline-flex items-center">
                    <CalendarIcon class="h-4 w-4 mr-1" />
                    {{ formatDateRange(a.start_at, a.end_at) }}
                  </span>
                  <span v-if="isActive(a)" class="inline-flex items-center px-1.5 py-0.5 rounded bg-green-50 text-green-700 border border-green-200">进行中</span>
                </div>
              </div>
              <div class="flex-shrink-0 flex items-start gap-2">
                <a
                  v-if="a.cta_url"
                  class="inline-flex items-center px-3 py-1.5 text-sm rounded-md border border-gray-200 text-gray-700 hover:bg-gray-50"
                  :href="a.cta_url"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <ExternalLinkIcon class="h-4 w-4 mr-1" /> 前往
                </a>
                <button
                  class="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-primary-600 text-white hover:bg-primary-700"
                  @click="openDetail(a)"
                >
                  查看详情
                </button>
              </div>
            </div>
          </div>

          <!-- 加载更多（示例） -->
          <div class="flex justify-center mt-4" v-if="canLoadMore">
            <button class="px-4 py-2 text-sm rounded-md border border-gray-200 text-gray-700 hover:bg-gray-50" @click="loadMore">
              加载更多
            </button>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-16 text-gray-500">
          <MegaphoneIcon class="h-12 w-12 mx-auto mb-4 text-gray-300" />
          <p>暂无公告</p>
        </div>
      </div>
    </PageContainer>

    <!-- 详情抽屉 -->
    <div v-if="showDrawer" class="fixed inset-0 z-50">
      <div class="absolute inset-0 bg-black/40" @click="closeDetail"></div>
      <div class="absolute right-0 top-0 h-full w-full max-w-xl bg-white shadow-xl flex flex-col">
        <div class="px-4 py-3 border-b border-gray-200 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <MegaphoneIcon class="h-5 w-5 text-primary-600" />
            <h4 class="text-base font-medium text-gray-900">{{ current?.title }}</h4>
          </div>
          <button @click="closeDetail" class="p-1 rounded hover:bg-gray-100">
            <XIcon class="h-5 w-5" />
          </button>
        </div>
        <div class="p-4 overflow-y-auto">
          <div class="flex items-center gap-2 mb-3">
            <span
              class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
              :class="typeBadgeClass(current?.type || 'system')"
            >
              {{ typeLabel(current?.type || 'system') }}
            </span>
            <span v-if="current?.pinned" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-secondary-50 text-secondary-700 border border-secondary-200">
              <PinIcon class="h-3.5 w-3.5 mr-1" /> 置顶
            </span>
            <span v-if="current?.require_ack" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-red-50 text-red-700 border border-red-200">
              <AlertTriangleIcon class="h-3.5 w-3.5 mr-1" /> 重要
            </span>
          </div>
          <div class="text-xs text-gray-500 mb-4 flex items-center">
            <CalendarIcon class="h-4 w-4 mr-1" /> {{ formatDateRange(current?.start_at, current?.end_at) }}
            <span v-if="current && isActive(current)" class="ml-2 inline-flex items-center px-1.5 py-0.5 rounded bg-green-50 text-green-700 border border-green-200">进行中</span>
          </div>
          <div class="prose prose-sm max-w-none">
            <p class="whitespace-pre-line text-gray-800">{{ current?.content }}</p>
          </div>
        </div>
        <div class="px-4 py-3 border-t border-gray-200 flex items-center justify-end gap-2">
          <a
            v-if="current?.cta_url"
            class="inline-flex items-center px-3 py-1.5 text-sm rounded-md border border-gray-200 text-gray-700 hover:bg-gray-50"
            :href="current.cta_url"
            target="_blank"
            rel="noopener noreferrer"
          >
            <ExternalLinkIcon class="h-4 w-4 mr-1" /> 前往
          </a>
          <button class="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-primary-600 text-white hover:bg-primary-700" @click="closeDetail">
            知道了
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
  
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { AppLayout, PageContainer } from '@/components/layout'
import { MegaphoneIcon, PinIcon, AlertTriangleIcon, CalendarIcon, ExternalLinkIcon, SearchIcon, XIcon } from 'lucide-vue-next'

type AnnouncementType = 'activity' | 'system' | 'release'

interface Announcement {
  id: number
  title: string
  type: AnnouncementType
  summary: string
  content: string
  start_at: string
  end_at?: string
  pinned?: boolean
  require_ack?: boolean
  cta_url?: string
}

const typeOptions = [
  { label: '全部', value: 'all' },
  { label: '活动', value: 'activity' },
  { label: '系统', value: 'system' },
  { label: '版本', value: 'release' }
]

const activeType = ref<'all' | AnnouncementType>('all')
const q = ref('')
const onlyActive = ref(true)

// 示例数据（后续可替换为接口数据）
const announcements = ref<Announcement[]>([
  {
    id: 1,
    title: '中秋亲子活动报名开启',
    type: 'activity',
    summary: '限时报名，参与即可获得积分奖励与精美礼品，详情点击查看。',
    content: '活动时间：9/25 - 9/30\n地点：线上+线下联合\n报名方式：点击前往活动页，填写报名信息。\n奖励：参与即送积分，优胜家庭可获定制礼品。',
    start_at: new Date(Date.now() - 86400000).toISOString(),
    end_at: new Date(Date.now() + 5 * 86400000).toISOString(),
    pinned: true,
    cta_url: 'https://example.com/activity'
  },
  {
    id: 2,
    title: '系统维护通知',
    type: 'system',
    summary: '本周六凌晨进行例行维护，期间可能出现短暂不可用。',
    content: '维护窗口：周六 02:00 - 03:00\n影响范围：站点访问、图片上传\n建议：如遇异常，请稍后重试。',
    start_at: new Date(Date.now() + 2 * 86400000).toISOString(),
    end_at: new Date(Date.now() + 2 * 86400000 + 3600000).toISOString(),
    require_ack: true
  },
  {
    id: 3,
    title: 'v2.1.0 发布说明',
    type: 'release',
    summary: '新增“公告中心”，优化积分展示与筛选体验，修复若干问题。',
    content: '亮点：\n- 新增公告中心页面\n- 积分记录仅展示最近8条\n- 优化配色与可读性\nBugfix：若干已知问题修复。',
    start_at: new Date(Date.now() - 2 * 86400000).toISOString()
  },
  {
    id: 4,
    title: '线下沙龙报名火热进行中',
    type: 'activity',
    summary: '主题：高效亲子沟通技巧；名额有限，先到先得。',
    content: '时间：下周日 14:00\n地点：XX中心 3F 多功能厅\n报名：扫描活动页二维码报名。',
    start_at: new Date(Date.now() - 3 * 86400000).toISOString(),
    end_at: new Date(Date.now() + 7 * 86400000).toISOString(),
    cta_url: 'https://example.com/salon'
  },
  {
    id: 5,
    title: '安全提醒：请勿泄露账号信息',
    type: 'system',
    summary: '近期出现冒充官方的诈骗行为，请提高警惕，注意账号安全。',
    content: '官方不会以任何名义索要您的密码或验证码。若遇可疑信息，请通过站内客服反馈并保存证据。',
    start_at: new Date(Date.now() - 10 * 86400000).toISOString()
  },
  {
    id: 6,
    title: 'v2.0.0 大版本升级',
    type: 'release',
    summary: '全新界面风格，支持积分等级与成长任务，欢迎体验。',
    content: '新增：积分系统、成长任务；\n优化：首页加载与渲染性能；\n兼容性：建议升级至最新浏览器版本。',
    start_at: new Date(Date.now() - 30 * 86400000).toISOString()
  }
])

const isActive = (a: Announcement): boolean => {
  const now = Date.now()
  const start = a.start_at ? new Date(a.start_at).getTime() : 0
  const end = a.end_at ? new Date(a.end_at).getTime() : Number.POSITIVE_INFINITY
  return now >= start && now <= end
}

const filteredAnnouncements = computed(() => {
  let list = announcements.value.slice()
  // 类型
  if (activeType.value !== 'all') list = list.filter(a => a.type === activeType.value)
  // 搜索
  const query = q.value.trim().toLowerCase()
  if (query) list = list.filter(a => a.title.toLowerCase().includes(query) || a.summary.toLowerCase().includes(query) || a.content.toLowerCase().includes(query))
  // 仅进行中
  if (onlyActive.value) list = list.filter(a => isActive(a))
  // 排序：置顶优先，其次按开始时间倒序
  list.sort((a, b) => {
    if ((b.pinned ? 1 : 0) - (a.pinned ? 1 : 0) !== 0) {
      return (b.pinned ? 1 : 0) - (a.pinned ? 1 : 0)
    }
    return (new Date(b.start_at).getTime()) - (new Date(a.start_at).getTime())
  })
  return list.slice(0, pageSize.value * pageIndex.value)
})

const pageIndex = ref(1)
const pageSize = ref(5)
const canLoadMore = computed(() => (announcements.value.length > pageSize.value * pageIndex.value))
const loadMore = () => { pageIndex.value += 1 }

const typeLabel = (t: AnnouncementType | 'all'): string => {
  switch (t) {
    case 'activity': return '活动'
    case 'system': return '系统'
    case 'release': return '版本'
    default: return '全部'
  }
}

const typeBadgeClass = (t: AnnouncementType) => {
  switch (t) {
    case 'activity':
      return 'bg-primary-50 text-primary-700 border border-primary-200'
    case 'system':
      return 'bg-secondary-50 text-secondary-700 border border-secondary-200'
    case 'release':
      return 'bg-gray-100 text-gray-700 border border-gray-200'
    default:
      return 'bg-gray-100 text-gray-700 border border-gray-200'
  }
}

const formatDateRange = (start?: string, end?: string) => {
  if (!start && !end) return '时间未定'
  const fmt = (s: string) => new Date(s).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
  if (start && end) return `${fmt(start)} - ${fmt(end)}`
  if (start) return `${fmt(start)} 开始`
  return `${fmt(end!)} 结束`
}

// 详情抽屉
const showDrawer = ref(false)
const current = ref<Announcement | null>(null)
const openDetail = (a: Announcement) => { current.value = a; showDrawer.value = true }
const closeDetail = () => { showDrawer.value = false; current.value = null }
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

