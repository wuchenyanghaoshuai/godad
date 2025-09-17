<template>
  <div class="markdown-editor">
    <div class="editor-toolbar border-b border-gray-200 p-2 sm:p-3 flex flex-wrap gap-1 sm:gap-2">
      <!-- 格式化按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="insertMarkdown('bold')"
          class="toolbar-btn"
          title="粗体 (Ctrl+B)"
        >
          <BoldIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="insertMarkdown('italic')"
          class="toolbar-btn"
          title="斜体 (Ctrl+I)"
        >
          <ItalicIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="insertMarkdown('code')"
          class="toolbar-btn"
          title="内联代码"
        >
          <CodeIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="insertMarkdown('strikethrough')"
          class="toolbar-btn"
          title="删除线"
        >
          <StrikethroughIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 标题按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="insertMarkdown('h1')"
          class="toolbar-btn"
          title="标题 1"
        >
          H1
        </button>
        <button
          type="button"
          @click="insertMarkdown('h2')"
          class="toolbar-btn"
          title="标题 2"
        >
          H2
        </button>
        <button
          type="button"
          @click="insertMarkdown('h3')"
          class="toolbar-btn"
          title="标题 3"
        >
          H3
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 列表按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="insertMarkdown('ul')"
          class="toolbar-btn"
          title="无序列表"
        >
          <ListIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="insertMarkdown('ol')"
          class="toolbar-btn"
          title="有序列表"
        >
          <ListOrderedIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="insertMarkdown('quote')"
          class="toolbar-btn"
          title="引用"
        >
          <QuoteIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 链接和图片 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="insertMarkdown('link')"
          class="toolbar-btn"
          title="插入链接 (Ctrl+K)"
        >
          <LinkIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="showImageUpload = true"
          class="toolbar-btn"
          title="插入图片"
        >
          <ImageIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 模式切换 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="togglePreviewMode"
          :class="['toolbar-btn', previewMode && 'active']"
          title="预览模式"
        >
          <EyeIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
      </div>
    </div>
    
    <!-- 编辑器内容区 -->
    <div class="editor-content relative">
      <!-- 拖拽提示层 -->
      <div
        v-if="isDragOver"
        class="absolute inset-0 bg-blue-50 border-2 border-dashed border-blue-300 flex items-center justify-center z-10"
      >
        <div class="text-center text-blue-600">
          <ImageIcon class="h-12 w-12 mx-auto mb-2" />
          <p class="text-lg font-medium">拖拽图片到这里上传</p>
          <p class="text-sm">支持 JPG、PNG、GIF 格式，最大 10MB</p>
        </div>
      </div>
      
      <div class="flex" :style="{ minHeight: `${minHeight}px` }">
        <!-- Markdown编辑区 -->
        <div :class="[previewMode ? 'w-1/2 border-r border-gray-200' : 'w-full']">
          <textarea
            ref="editorRef"
            v-model="content"
            @input="handleInput"
            @keydown="handleKeydown"
            @paste="handlePaste"
            @dragover.prevent="handleDragOver"
            @dragleave.prevent="handleDragLeave"
            @drop.prevent="handleDrop"
            :placeholder="placeholder"
            class="markdown-area p-3 sm:p-4 font-mono text-sm resize-none focus:outline-none w-full h-full leading-relaxed border-none"
            :style="{ minHeight: `${minHeight}px` }"
            :disabled="props.disabled"
          ></textarea>
        </div>
        
        <!-- 预览区 -->
        <div v-if="previewMode" class="w-1/2 p-3 sm:p-4 overflow-y-auto">
          <div 
            class="preview-area prose prose-sm max-w-none"
            v-html="htmlContent"
          ></div>
        </div>
      </div>
      
      <!-- 状态栏 -->
      <div class="absolute bottom-2 left-2 right-2 flex justify-between items-center text-xs text-gray-500">
        <!-- 自动保存状态 -->
        <div class="bg-white bg-opacity-80 px-2 py-1 rounded shadow-sm">
          <span v-if="autoSaveStatus === 'saving'" class="text-blue-600">
            <svg class="inline-block w-3 h-3 mr-1 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            正在保存草稿...
          </span>
          <span v-else-if="autoSaveStatus === 'saved'" class="text-green-600">
            ✓ 草稿已自动保存
          </span>
          <span v-else-if="autoSaveStatus === 'error'" class="text-red-600">
            ✗ 草稿保存失败
          </span>
          <span v-else class="text-gray-400">草稿自动保存</span>
        </div>
        
        <!-- 字数统计 -->
        <div class="bg-white bg-opacity-80 px-2 py-1 rounded shadow-sm">
          <span class="font-medium">{{ wordCount }}</span> 字
        </div>
      </div>
    </div>
    
    <!-- 图片上传弹窗 -->
    <div
      v-if="showImageUpload"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="showImageUpload = false"
    >
      <div
        class="bg-white rounded-lg p-6 max-w-md w-full mx-4"
        @click.stop
      >
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">插入图片</h3>
          <button
            @click="showImageUpload = false"
            class="text-gray-400 hover:text-gray-600"
          >
            <XIcon class="h-5 w-5" />
          </button>
        </div>
        
        <ImageUpload
          v-model="uploadedImageUrl"
          @upload-success="handleImageUpload"
          :multiple="false"
          :max-size="10"
        />
        
        <div class="mt-4 flex justify-end gap-2">
          <button
            @click="showImageUpload = false"
            class="px-4 py-2 text-gray-600 hover:text-gray-800"
          >
            取消
          </button>
          <button
            @click="insertImage"
            :disabled="!uploadedImageUrl"
            class="px-4 py-2 bg-pink-500 text-white rounded hover:bg-pink-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            插入
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { marked } from 'marked'
import {
  BoldIcon,
  ItalicIcon,
  StrikethroughIcon,
  ListIcon,
  ListOrderedIcon,
  LinkIcon,
  ImageIcon,
  CodeIcon,
  XIcon,
  EyeIcon,
  QuoteIcon
} from 'lucide-vue-next'
import ImageUpload from './ImageUpload.vue'
import type { ImageUploadResponse } from '@/api/types'
import { useToast } from '@/composables/useToast'

// Props
interface Props {
  modelValue?: string
  placeholder?: string
  minHeight?: number
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入Markdown内容...',
  minHeight: 300,
  disabled: false
})

// Emits
interface Emits {
  'update:modelValue': [value: string]
  'change': [value: string]
}

const emit = defineEmits<Emits>()

// 响应式数据
const editorRef = ref<HTMLTextAreaElement>()
const content = ref('')
const previewMode = ref(false)
const showImageUpload = ref(false)
const uploadedImageUrl = ref('')
const isDragOver = ref(false)
const autoSaveStatus = ref<'idle' | 'saving' | 'saved' | 'error'>('idle')

// Toast
const { toast } = useToast()
const dragCounter = ref(0)

// 计算属性
const wordCount = computed(() => {
  // Markdown模式：计算纯文本长度
  const text = content.value.replace(/[#*`>\-_\[\]()!]/g, '').replace(/\s/g, '')
  return text.length
})

// Markdown 转 HTML
const htmlContent = computed(() => {
  if (!content.value) return ''
  try {
    return marked(content.value)
  } catch (error) {
    console.error('Markdown 解析错误:', error)
    return '<p>Markdown 解析出错</p>'
  }
})

// 监听modelValue变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue !== content.value) {
      content.value = newValue || ''
    }
  },
  { immediate: true }
)

// 监听内容变化
watch(content, (newValue) => {
  emit('update:modelValue', newValue)
  emit('change', newValue)
})

// 生命周期
onMounted(() => {
  nextTick(() => {
    if (props.modelValue) {
      content.value = props.modelValue
    } else {
      // 尝试加载草稿
      loadDraft()
    }
  })
})

onBeforeUnmount(() => {
  // 清理自动保存定时器
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
    autoSaveTimer = null
  }
})

// 方法
const handleInput = () => {
  emit('update:modelValue', content.value)
  emit('change', content.value)
  // 触发自动保存
  startAutoSave()
}

const handleKeydown = (e: KeyboardEvent) => {
  const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0
  const cmdKey = isMac ? e.metaKey : e.ctrlKey
  
  // 快捷键映射
  if (cmdKey && !e.shiftKey && !e.altKey) {
    switch (e.key.toLowerCase()) {
      case 'b': // Ctrl/Cmd + B: 粗体
        e.preventDefault()
        insertMarkdown('bold')
        break
      case 'i': // Ctrl/Cmd + I: 斜体
        e.preventDefault()
        insertMarkdown('italic')
        break
      case 'k': // Ctrl/Cmd + K: 插入链接
        e.preventDefault()
        insertMarkdown('link')
        break
      case 'p': // Ctrl/Cmd + P: 预览模式
        e.preventDefault()
        togglePreviewMode()
        break
    }
  }
  
  // 处理Tab键
  if (e.key === 'Tab') {
    e.preventDefault()
    insertAtCursor('    ') // 4个空格缩进
  }
}

const handlePaste = (e: ClipboardEvent) => {
  // 让浏览器默认处理粘贴，保持Markdown格式
  const text = e.clipboardData?.getData('text/plain') || ''
  // 可以在这里添加粘贴内容的处理逻辑
}

// Markdown 操作方法
const insertMarkdown = (type: string) => {
  if (props.disabled || !editorRef.value) return
  
  const textarea = editorRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end)
  
  let markdownText = ''
  let newCursorPos = start
  
  switch (type) {
    case 'bold':
      markdownText = `**${selectedText || '粗体文字'}**`
      newCursorPos = selectedText ? end + 4 : start + 2
      break
    case 'italic':
      markdownText = `*${selectedText || '斜体文字'}*`
      newCursorPos = selectedText ? end + 2 : start + 1
      break
    case 'code':
      markdownText = `\`${selectedText || '内联代码'}\``
      newCursorPos = selectedText ? end + 2 : start + 1
      break
    case 'strikethrough':
      markdownText = `~~${selectedText || '删除线文字'}~~`
      newCursorPos = selectedText ? end + 4 : start + 2
      break
    case 'h1':
      markdownText = `# ${selectedText || '一级标题'}`
      newCursorPos = selectedText ? end + 2 : start + 2
      break
    case 'h2':
      markdownText = `## ${selectedText || '二级标题'}`
      newCursorPos = selectedText ? end + 3 : start + 3
      break
    case 'h3':
      markdownText = `### ${selectedText || '三级标题'}`
      newCursorPos = selectedText ? end + 4 : start + 4
      break
    case 'ul':
      markdownText = `- ${selectedText || '列表项'}`
      newCursorPos = selectedText ? end + 2 : start + 2
      break
    case 'ol':
      markdownText = `1. ${selectedText || '列表项'}`
      newCursorPos = selectedText ? end + 3 : start + 3
      break
    case 'quote':
      markdownText = `> ${selectedText || '引用内容'}`
      newCursorPos = selectedText ? end + 2 : start + 2
      break
    case 'link':
      const url = prompt('请输入链接地址:')
      if (url) {
        markdownText = `[${selectedText || '链接文字'}](${url})`
        newCursorPos = selectedText ? end + url.length + 4 : start + 4
      } else {
        return
      }
      break
    default:
      return
  }
  
  const newContent = content.value.substring(0, start) + markdownText + content.value.substring(end)
  content.value = newContent
  
  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(newCursorPos, newCursorPos)
    handleInput()
  })
}

const insertAtCursor = (text: string) => {
  if (!editorRef.value) return
  
  const textarea = editorRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  
  const newContent = content.value.substring(0, start) + text + content.value.substring(end)
  content.value = newContent
  
  nextTick(() => {
    const newPos = start + text.length
    textarea.setSelectionRange(newPos, newPos)
    textarea.focus()
    handleInput()
  })
}

const handleImageUpload = (response: ImageUploadResponse) => {
  uploadedImageUrl.value = response.public_url
}

const insertImage = () => {
  if (uploadedImageUrl.value) {
    const markdownImage = `![插入的图片](${uploadedImageUrl.value})`
    insertAtCursor(markdownImage)
    showImageUpload.value = false
    uploadedImageUrl.value = ''
  }
}

// 拖拽上传功能
const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
  dragCounter.value++
  if (!isDragOver.value) {
    isDragOver.value = true
  }
}

const handleDragLeave = (e: DragEvent) => {
  e.preventDefault()
  dragCounter.value--
  if (dragCounter.value === 0) {
    isDragOver.value = false
  }
}

const handleDrop = async (e: DragEvent) => {
  e.preventDefault()
  isDragOver.value = false
  dragCounter.value = 0
  
  const files = e.dataTransfer?.files
  if (!files || files.length === 0) return
  
  const file = files[0]
  if (!file.type.startsWith('image/')) {
    toast.error('只能上传图片文件')
    return
  }
  
  if (file.size > 10 * 1024 * 1024) { // 10MB
    toast.error('图片大小不能超过 10MB')
    return
  }
  
  try {
    // 这里应该上传到服务器，暂时使用本地预览
    const reader = new FileReader()
    reader.onload = (e) => {
      const imageUrl = e.target?.result as string
      const markdownImage = `![拖拽上传的图片](${imageUrl})`
      insertAtCursor(markdownImage)
    }
    reader.readAsDataURL(file)
    
  } catch (error) {
    console.error('图片上传失败:', error)
    toast.error('图片上传失败，请重试')
  }
}

// 自动保存草稿功能
let autoSaveTimer: NodeJS.Timeout | null = null

const startAutoSave = () => {
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
  }
  
  autoSaveTimer = setTimeout(() => {
    saveDraft()
  }, 5000) // 5秒后自动保存
}

const saveDraft = async () => {
  if (!content.value.trim()) return
  
  try {
    autoSaveStatus.value = 'saving'
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 保存到本地存储作为备份
    localStorage.setItem('markdown_editor_draft', content.value)
    localStorage.setItem('markdown_editor_draft_time', new Date().toISOString())
    
    autoSaveStatus.value = 'saved'
    
    // 3秒后隐藏保存成功状态
    setTimeout(() => {
      if (autoSaveStatus.value === 'saved') {
        autoSaveStatus.value = 'idle'
      }
    }, 3000)
    
  } catch (error) {
    console.error('保存草稿失败:', error)
    autoSaveStatus.value = 'error'
    
    // 3秒后隐藏错误状态
    setTimeout(() => {
      if (autoSaveStatus.value === 'error') {
        autoSaveStatus.value = 'idle'
      }
    }, 3000)
  }
}

// 加载草稿
const loadDraft = () => {
  const draft = localStorage.getItem('markdown_editor_draft')
  const draftTime = localStorage.getItem('markdown_editor_draft_time')
  
  if (draft && draftTime && !props.modelValue) {
    const timeDiff = Date.now() - new Date(draftTime).getTime()
    // 如果草稿时间在24小时内，则提示用户是否恢复
    if (timeDiff < 24 * 60 * 60 * 1000) {
      if (confirm('检测到未保存的Markdown草稿，是否恢复？')) {
        content.value = draft
        emit('update:modelValue', draft)
      }
    }
  }
}

const togglePreviewMode = () => {
  previewMode.value = !previewMode.value
  if (!previewMode.value && editorRef.value) {
    nextTick(() => {
      editorRef.value?.focus()
    })
  }
}
</script>

<style scoped>
.markdown-editor {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.toolbar-btn {
  @apply px-1.5 py-1 sm:px-2 sm:py-1.5 border border-gray-300 rounded hover:bg-gray-50 transition-colors text-xs sm:text-sm;
}

.toolbar-btn.active {
  @apply bg-blue-100 border-blue-300 text-blue-700;
}

.markdown-area {
  min-height: 200px;
  border: none;
  outline: none;
  resize: none;
}

.markdown-area:focus {
  outline: none;
}

.preview-area {
  background-color: #fafafa;
}

/* 预览区域样式 */
.preview-area h1 {
  @apply text-2xl font-bold mb-4 mt-6 pb-2 border-b border-gray-200;
}

.preview-area h2 {
  @apply text-xl font-bold mb-3 mt-5;
}

.preview-area h3 {
  @apply text-lg font-bold mb-2 mt-4;
}

.preview-area ul {
  @apply list-disc list-inside mb-4 pl-4;
}

.preview-area ol {
  @apply list-decimal list-inside mb-4 pl-4;
}

.preview-area p {
  @apply mb-3 leading-relaxed;
}

.preview-area a {
  @apply text-blue-600 hover:text-blue-800 underline;
}

.preview-area img {
  @apply max-w-full h-auto rounded shadow-sm my-4;
}

.preview-area blockquote {
  @apply border-l-4 border-gray-300 pl-4 py-2 my-4 bg-gray-50 italic;
}

.preview-area code {
  @apply bg-gray-100 px-1 py-0.5 rounded text-sm font-mono;
}

.preview-area pre {
  @apply bg-gray-100 p-3 rounded overflow-x-auto my-4;
}

.preview-area pre code {
  @apply bg-transparent p-0;
}

/* 移动端优化 */
@media (max-width: 640px) {
  .editor-toolbar {
    @apply sticky top-0 bg-white z-10;
  }
  
  .toolbar-btn {
    @apply min-w-[32px] min-h-[32px] flex items-center justify-center;
  }
}
</style>