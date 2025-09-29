<template>
  <div
    v-if="showNotifications"
    class="relative"
    ref="notifMenuWrapper"
    @mouseenter="openNotifMenuHover"
    @mouseleave="closeNotifMenuHover"
  >
    <button
      class="relative p-2 rounded-lg text-gray-600 hover:text-primary-600 hover:bg-primary-50 transition-all duration-200"
      title="消息通知"
      @click="toggleNotifMenu"
      @mouseenter="cancelNotifMenuHide"
    >
      <BellIcon class="h-5 w-5" />
      <!-- 未读通知红点 -->
      <span
        v-if="totalUnreadCount > 0"
        class="absolute -top-1 -right-1 inline-flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-red-500 rounded-full"
      >
        {{ totalUnreadCount > 99 ? '99+' : totalUnreadCount }}
      </span>
    </button>

    <!-- hover 桥接层，避免按钮与菜单之间空隙导致抖动 -->
    <div
      v-if="showNotifMenu"
      class="absolute right-0 top-full w-56 h-3"
      @mouseenter="cancelNotifMenuHide"
    ></div>

    <!-- 下拉列表 -->
    <div
      v-if="showNotifMenu"
      class="absolute right-0 top-full mt-1 w-56 bg-white border border-gray-200 rounded-lg shadow-lg z-50 py-1"
      role="menu"
      @mouseenter="cancelNotifMenuHide"
      @mouseleave="closeNotifMenuHover"
    >
      <!-- 回复我的（评论/回复） -->
      <button
        class="w-full flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
        role="menuitem"
        @click="goNotifications('comment')"
      >
        <MessageCircleIcon class="h-4 w-4 mr-2 text-gray-400" />
        <span>回复我的</span>
        <span v-if="unreadByType.comment > 0" class="ml-auto inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full text-xs font-bold bg-red-500 text-white">{{ unreadByType.comment > 99 ? '99+' : unreadByType.comment }}</span>
      </button>

      <!-- @我的 -->
      <button
        class="w-full flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
        role="menuitem"
        @click="goNotifications('mention')"
      >
        <AtSignIcon class="h-4 w-4 mr-2 text-gray-400" />
        <span>@我的</span>
        <span v-if="unreadByType.mention > 0" class="ml-auto inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full text-xs font-bold bg-red-500 text-white">{{ unreadByType.mention > 99 ? '99+' : unreadByType.mention }}</span>
      </button>

      <!-- 收到的赞 -->
      <button
        class="w-full flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
        role="menuitem"
        @click="goNotifications('like')"
      >
        <ThumbsUpIcon class="h-4 w-4 mr-2 text-gray-400" />
        <span>收到的赞</span>
        <span v-if="unreadByType.like > 0" class="ml-auto inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full text-xs font-bold bg-red-500 text-white">{{ unreadByType.like > 99 ? '99+' : unreadByType.like }}</span>
      </button>

      <!-- 系统消息 -->
      <button
        class="w-full flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
        role="menuitem"
        @click="goNotifications('system')"
      >
        <AlertCircleIcon class="h-4 w-4 mr-2 text-gray-400" />
        <span>系统消息</span>
        <span v-if="unreadByType.system && unreadByType.system > 0" class="ml-auto inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full text-xs font-bold bg-red-500 text-white">{{ unreadByType.system > 99 ? '99+' : unreadByType.system }}</span>
      </button>

      <!-- 我的消息（私信） -->
      <button
        class="w-full flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
        role="menuitem"
        @click="goNotifications('message')"
      >
        <MessageSquareIcon class="h-4 w-4 mr-2 text-gray-400" />
        <span>我的消息</span>
        <span v-if="unreadByType.message > 0" class="ml-auto inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full text-xs font-bold bg-red-500 text-white">{{ unreadByType.message > 99 ? '99+' : unreadByType.message }}</span>
      </button>

      <!-- 其他通知（收藏、关注等） -->
      <button
        class="w-full flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
        role="menuitem"
        @click="goNotifications('other')"
      >
        <AlertCircleIcon class="h-4 w-4 mr-2 text-gray-400" />
        <span>其他通知</span>
        <span v-if="otherUnreadCount > 0" class="ml-auto inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full text-xs font-bold bg-red-500 text-white">{{ otherUnreadCount > 99 ? '99+' : otherUnreadCount }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  BellIcon,
  MessageCircleIcon,
  AtSignIcon,
  ThumbsUpIcon,
  AlertCircleIcon,
  MessageSquareIcon
} from 'lucide-vue-next'

interface Props {
  showNotifications?: boolean
  totalUnreadCount?: number
  unreadByType?: { [k: string]: number }
}

const props = withDefaults(defineProps<Props>(), {
  showNotifications: true,
  totalUnreadCount: 0,
  unreadByType: () => ({ message: 0, like: 0, comment: 0, follow: 0, bookmark: 0, system: 0, mention: 0, total_unread: 0 })
})

const emit = defineEmits<{
  fetchUnreadByType: []
}>()

const router = useRouter()

const notifMenuWrapper = ref<HTMLElement | null>(null)
const showNotifMenu = ref(false)
let notifMenuHideTimer: number | null = null

// 计算属性：其他通知数量（收藏、关注等）
const otherUnreadCount = computed(() => {
  return (props.unreadByType.bookmark || 0) + (props.unreadByType.follow || 0)
})

// 通知下拉控制与路由跳转
const cancelNotifMenuHide = () => {
  if (notifMenuHideTimer) {
    clearTimeout(notifMenuHideTimer)
    notifMenuHideTimer = null
  }
}

const openNotifMenuHover = () => {
  cancelNotifMenuHide()
  showNotifMenu.value = true
  emit('fetchUnreadByType')
}

const closeNotifMenuHover = () => {
  cancelNotifMenuHide()
  notifMenuHideTimer = window.setTimeout(() => {
    showNotifMenu.value = false
    notifMenuHideTimer = null
  }, 220)
}

const toggleNotifMenu = () => {
  showNotifMenu.value = !showNotifMenu.value
}

const goNotifications = (category: 'message' | 'like' | 'comment' | 'system' | 'mention' | 'other') => {
  showNotifMenu.value = false
  if (category === 'message') {
    router.push({ path: '/notifications', query: { tab: 'message' } })
    return
  }
  // 其他通知类型，包含收藏和关注
  if (category === 'other') {
    router.push({ path: '/notifications', query: { tab: 'notify', category: 'bookmark,follow' } })
    return
  }
  // 其余归类到通知
  const query: any = { tab: 'notify' }
  if (['like', 'comment', 'follow', 'bookmark', 'system', 'mention'].includes(category)) {
    query.category = category
  }
  router.push({ path: '/notifications', query })
}
</script>