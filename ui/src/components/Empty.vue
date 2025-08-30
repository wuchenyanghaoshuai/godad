<template>
  <div class="empty-state flex flex-col items-center justify-center py-8 sm:py-12 px-4 text-center">
    <!-- 图标 -->
    <div class="mb-4 sm:mb-6">
      <component 
        :is="iconComponent" 
        class="h-12 w-12 sm:h-16 sm:w-16 text-gray-300 mx-auto"
      />
    </div>
    
    <!-- 标题 -->
    <h3 class="text-base sm:text-lg font-medium text-gray-900 mb-2">
      {{ title || '暂无数据' }}
    </h3>
    
    <!-- 描述 -->
    <p class="text-sm sm:text-base text-gray-500 mb-4 sm:mb-6 max-w-sm">
      {{ description || '当前没有任何内容，您可以尝试其他操作' }}
    </p>
    
    <!-- 操作按钮 -->
    <div v-if="showAction" class="flex flex-col sm:flex-row gap-2 sm:gap-3">
      <button
        v-if="actionText"
        @click="$emit('action')"
        class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
      >
        {{ actionText }}
      </button>
      
      <button
        v-if="secondaryActionText"
        @click="$emit('secondary-action')"
        class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors"
      >
        {{ secondaryActionText }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { 
  InboxIcon, 
  DocumentTextIcon, 
  UserGroupIcon, 
  PhotoIcon,
  ChatBubbleLeftRightIcon,
  FolderIcon
} from '@heroicons/vue/24/outline'

// 定义组件属性
interface Props {
  type?: 'default' | 'article' | 'comment' | 'user' | 'image' | 'folder'
  title?: string
  description?: string
  actionText?: string
  secondaryActionText?: string
  showAction?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  showAction: true
})

// 定义事件
defineEmits<{
  action: []
  'secondary-action': []
}>()

// 根据类型选择图标
const iconComponent = computed(() => {
  const iconMap = {
    default: InboxIcon,
    article: DocumentTextIcon,
    comment: ChatBubbleLeftRightIcon,
    user: UserGroupIcon,
    image: PhotoIcon,
    folder: FolderIcon
  }
  return iconMap[props.type] || InboxIcon
})
</script>
