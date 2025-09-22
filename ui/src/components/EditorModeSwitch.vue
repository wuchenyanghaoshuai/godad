<template>
  <div class="editor-mode-switch">
    <!-- 模式选择器 -->
    <div class="mode-selector mb-4">
      <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
        <div class="flex items-center gap-2">
          <EditIcon class="h-5 w-5 text-gray-600" />
          <span class="text-sm font-medium text-gray-700">编辑器模式</span>
        </div>
        
        <div class="flex items-center bg-white rounded-lg p-1 shadow-sm">
          <button
            @click="switchMode('rich')"
            :class="[
              'px-3 py-2 text-sm font-medium rounded-md transition-all duration-200',
              currentMode === 'rich'
                ? 'bg-blue-100 text-blue-700 shadow-sm'
                : 'text-gray-600 hover:text-gray-800 hover:bg-gray-50'
            ]"
          >
            <div class="flex items-center gap-2">
              <TypeIcon class="h-4 w-4" />
              <span>富文本</span>
            </div>
          </button>
          
          <button
            @click="switchMode('markdown')"
            :class="[
              'px-3 py-2 text-sm font-medium rounded-md transition-all duration-200',
              currentMode === 'markdown'
                ? 'bg-blue-100 text-blue-700 shadow-sm'
                : 'text-gray-600 hover:text-gray-800 hover:bg-gray-50'
            ]"
          >
            <div class="flex items-center gap-2">
              <CodeIcon class="h-4 w-4" />
              <span>Markdown</span>
            </div>
          </button>
        </div>
      </div>
      
      <!-- 模式说明 -->
      <div class="mt-2 text-xs text-gray-500">
        <div v-if="currentMode === 'rich'" class="flex items-center gap-1">
          <InfoIcon class="h-3 w-3" />
          <span>使用所见即所得编辑器，支持实时预览、颜色格式化和拖拽操作</span>
        </div>
        <div v-else class="flex items-center gap-1">
          <InfoIcon class="h-3 w-3" />
          <span>使用Markdown语法编写，支持实时预览和代码高亮</span>
        </div>
      </div>
    </div>
    
    <!-- 编辑器内容 -->
    <div class="editor-container">
      <!-- 富文本编辑器 -->
      <TipTapEditor
        v-if="currentMode === 'rich'"
        v-model="editorContent"
        @change="handleEditorChange"
        :placeholder="richTextPlaceholder"
        :min-height="minHeight"
        :disabled="disabled"
      />
      
      <!-- Markdown编辑器 -->
      <MarkdownEditor
        v-else-if="currentMode === 'markdown'"
        v-model="editorContent"
        @change="handleEditorChange"
        :placeholder="markdownPlaceholder"
        :min-height="minHeight"
        :disabled="disabled"
      />
    </div>
    
    <!-- 模式切换确认对话框 -->
    <div
      v-if="showSwitchConfirm"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="showSwitchConfirm = false"
    >
      <div
        class="bg-white rounded-lg p-6 max-w-md w-full mx-4"
        @click.stop
      >
        <div class="flex items-start gap-3">
          <AlertTriangleIcon class="h-6 w-6 text-yellow-500 mt-0.5" />
          <div class="flex-1">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">切换编辑器模式</h3>
            <p class="text-sm text-gray-600 mb-4">
              切换编辑器模式可能会改变内容格式。确定要从
              <span class="font-medium">{{ getModeDisplayName(currentMode) }}</span>
              切换到
              <span class="font-medium">{{ getModeDisplayName(pendingMode) }}</span>
              吗？
            </p>
            
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mb-4">
              <p class="text-xs text-yellow-700">
                <strong>注意：</strong>
                <span v-if="pendingMode === 'markdown'">
                  富文本内容将转换为HTML格式存储，但编辑时会显示为Markdown语法。
                </span>
                <span v-else>
                  Markdown内容将被转换为HTML格式，格式化信息可能会丢失。
                </span>
              </p>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end gap-3 mt-4">
          <button
            @click="showSwitchConfirm = false"
            class="px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-50 rounded-lg transition-colors"
          >
            取消
          </button>
          <button
            @click="confirmModeSwitch"
            class="px-4 py-2 text-sm font-medium bg-blue-600 text-white hover:bg-blue-700 rounded-lg transition-colors"
          >
            确认切换
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  EditIcon,
  TypeIcon,
  CodeIcon,
  InfoIcon,
  AlertTriangleIcon
} from 'lucide-vue-next'
import MarkdownEditor from './MarkdownEditor.vue'

// Props
interface Props {
  modelValue?: string
  contentType?: 'rich' | 'markdown'
  richTextPlaceholder?: string
  markdownPlaceholder?: string
  minHeight?: number
  disabled?: boolean
  showModeSwitch?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  contentType: 'rich', // 默认使用富文本模式
  richTextPlaceholder: '请输入内容...',
  markdownPlaceholder: '请输入Markdown内容...',
  minHeight: 300,
  disabled: false,
  showModeSwitch: true
})

// Emits
interface Emits {
  'update:modelValue': [value: string]
  'update:contentType': [contentType: 'rich' | 'markdown']
  'change': [content: string, contentType: 'rich' | 'markdown']
  'mode-switch': [fromMode: 'rich' | 'markdown', toMode: 'rich' | 'markdown']
}

const emit = defineEmits<Emits>()

// 响应式数据
const currentMode = ref<'rich' | 'markdown'>(props.contentType)
const editorContent = ref(props.modelValue)
const showSwitchConfirm = ref(false)
const pendingMode = ref<'rich' | 'markdown'>('rich')

// 计算属性
const getModeDisplayName = (mode: 'rich' | 'markdown'): string => {
  return mode === 'rich' ? '富文本编辑器' : 'Markdown编辑器'
}

// 监听外部props变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue !== editorContent.value) {
      editorContent.value = newValue
    }
  }
)

watch(
  () => props.contentType,
  (newType) => {
    if (newType !== currentMode.value) {
      currentMode.value = newType
    }
  }
)

// 监听编辑器内容变化
watch(editorContent, (newContent) => {
  emit('update:modelValue', newContent)
  emit('change', newContent, currentMode.value)
})

// 监听模式变化
watch(currentMode, (newMode) => {
  emit('update:contentType', newMode)
})

// 方法
const handleEditorChange = (content: string) => {
  editorContent.value = content
}

const switchMode = (targetMode: 'rich' | 'markdown') => {
  if (targetMode === currentMode.value) return
  
  // 如果有内容且模式不同，显示确认对话框
  if (editorContent.value.trim() && targetMode !== currentMode.value) {
    pendingMode.value = targetMode
    showSwitchConfirm.value = true
  } else {
    // 没有内容时直接切换
    performModeSwitch(targetMode)
  }
}

const confirmModeSwitch = () => {
  performModeSwitch(pendingMode.value)
  showSwitchConfirm.value = false
}

const performModeSwitch = (targetMode: 'rich' | 'markdown') => {
  const previousMode = currentMode.value
  currentMode.value = targetMode
  
  // 发射模式切换事件
  emit('mode-switch', previousMode, targetMode)
  
  // 内容格式转换（如果需要）
  // 注意：实际的内容转换会在后端处理，这里主要是界面切换
  // 富文本 -> Markdown: 内容保持HTML格式，但用Markdown编辑器编辑
  // Markdown -> 富文本: 内容转换为HTML，用富文本编辑器编辑
}

// 暴露给父组件的方法
defineExpose({
  switchMode,
  getCurrentMode: () => currentMode.value,
  getContent: () => editorContent.value,
  setContent: (content: string) => {
    editorContent.value = content
  }
})
</script>

<style scoped>
.editor-mode-switch {
  @apply w-full;
}

.mode-selector {
  @apply transition-all duration-200;
}

.editor-container {
  @apply transition-all duration-200;
}

/* 平滑过渡效果 */
.editor-container > * {
  @apply transition-opacity duration-200;
}

/* 移动端优化 */
@media (max-width: 640px) {
  .mode-selector .flex {
    @apply flex-col gap-3;
  }
  
  .mode-selector .flex:first-child {
    @apply flex-row;
  }
  
  .mode-selector button {
    @apply text-xs px-2 py-1.5;
  }
}
</style>
