<template>
  <div>
    <div ref="toolbarRef" class="editor-toolbar px-3 py-2 sm:px-4 sm:py-2 flex items-center gap-1 sm:gap-1.5 overflow-x-auto">
      <!-- 格式化按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="execCommand('bold')"
          :class="['toolbar-btn', isActive('bold') && 'active']"
          data-tooltip="粗体"
        >
          <BoldIcon class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          @click="execCommand('italic')"
          :class="['toolbar-btn', isActive('italic') && 'active']"
          data-tooltip="斜体"
        >
          <ItalicIcon class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          @click="execCommand('underline')"
          :class="['toolbar-btn', isActive('underline') && 'active']"
          data-tooltip="下划线"
        >
          <UnderlineIcon class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          @click="execCommand('strikeThrough')"
          :class="['toolbar-btn', isActive('strikeThrough') && 'active']"
          data-tooltip="删除线"
        >
          <span class="text-xs font-bold">S</span>
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 颜色按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <div class="relative">
          <input
            type="color"
            @change="changeTextColor($event)"
            class="w-8 h-6 border rounded cursor-pointer"
            data-tooltip="文字颜色"
            value="#000000"
          />
        </div>
        <div class="relative">
          <input
            type="color"
            @change="changeBackgroundColor($event)"
            class="w-8 h-6 border rounded cursor-pointer"
            data-tooltip="背景颜色"
            value="#ffff00"
          />
        </div>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 对齐按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="execCommand('justifyLeft')"
          :class="['toolbar-btn', isActive('justifyLeft') && 'active']"
          data-tooltip="左对齐"
        >
          <span class="text-xs">⬅</span>
        </button>
        <button
          type="button"
          @click="execCommand('justifyCenter')"
          :class="['toolbar-btn', isActive('justifyCenter') && 'active']"
          data-tooltip="居中"
        >
          <span class="text-xs">⬌</span>
        </button>
        <button
          type="button"
          @click="execCommand('justifyRight')"
          :class="['toolbar-btn', isActive('justifyRight') && 'active']"
          data-tooltip="右对齐"
        >
          <span class="text-xs">➡</span>
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 标题按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <select
          @change="formatHeading($event)"
          class="toolbar-select"
          data-tooltip="标题"
        >
          <option value="">正文</option>
          <option value="h1">标题 1</option>
          <option value="h2">标题 2</option>
          <option value="h3">标题 3</option>
        </select>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 列表按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="execCommand('insertUnorderedList')"
          :class="['toolbar-btn', isActive('insertUnorderedList') && 'active']"
          data-tooltip="无序列表"
        >
          <ListIcon class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          @click="execCommand('insertOrderedList')"
          :class="['toolbar-btn', isActive('insertOrderedList') && 'active']"
          data-tooltip="有序列表"
        >
          <ListOrderedIcon class="h-3.5 w-3.5" />
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 链接和图片 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="insertLink"
          class="toolbar-btn"
          data-tooltip="插入链接"
        >
          <LinkIcon class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          @click="showImageUpload = true"
          class="toolbar-btn"
          data-tooltip="插入图片"
        >
          <ImageIcon class="h-3.5 w-3.5" />
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 其他功能 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="execCommand('removeFormat')"
          class="toolbar-btn"
          data-tooltip="清除格式"
        >
          <EraserIcon class="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          @click="toggleSourceMode"
          :class="['toolbar-btn', sourceMode && 'active']"
          data-tooltip="源码模式"
        >
          <CodeIcon class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>
    
    <!-- 编辑器内容区 -->
    <div class="rich-text-editor">
      <div class="editor-content relative">
      <div
        v-if="!sourceMode"
        ref="editorRef"
        contenteditable="true"
        @input="handleInput"
        @keydown="handleKeydown"
        @paste="handlePaste"
        @compositionstart="handleCompositionStart"
        @compositionend="handleCompositionEnd"
        @compositionupdate="handleCompositionUpdate"
        class="editor-area p-3 sm:p-4 focus:outline-none text-sm sm:text-base leading-relaxed"
        :style="{ minHeight: `${minHeight}px` }"
      ></div>
      
      <textarea
        v-else
        v-model="sourceContent"
        @input="handleSourceInput"
        class="source-area p-3 sm:p-4 font-mono text-xs sm:text-sm resize-none focus:outline-none w-full leading-relaxed"
        :style="{ minHeight: `${minHeight}px` }"
      ></textarea>
      
      <!-- 字数统计 -->
      <div class="absolute bottom-2 right-2 text-xs text-gray-500 bg-white bg-opacity-80 px-2 py-1 rounded shadow-sm">
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
            type="button"
            @click="showImageUpload = false"
            class="text-gray-400 hover:text-gray-600"
          >
            <XIcon class="h-5 w-5" />
          </button>
        </div>
        
        <ImageUpload
          v-model="uploadedImageUrls"
          @upload-success="handleImageUpload"
          :multiple="true"
          :max-files="6"
          :max-size="10"
          :article-title="props.articleTitle"
        />
        
        <div class="mt-4 flex justify-end gap-2">
          <button
            type="button"
            @click="showImageUpload = false"
            class="px-4 py-2 text-gray-600 hover:text-gray-800"
          >
            取消
          </button>
          <button
            type="button"
            @click="insertImage"
            :disabled="uploadedImageUrls.length === 0"
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
import {
  BoldIcon,
  ItalicIcon,
  UnderlineIcon,
  ListIcon,
  ListOrderedIcon,
  LinkIcon,
  ImageIcon,
  EraserIcon,
  CodeIcon,
  XIcon
} from 'lucide-vue-next'
import ImageUpload from './ImageUpload.vue'
import type { ImageUploadResponse } from '@/api/types'

// Props
interface Props {
  modelValue?: string
  placeholder?: string
  minHeight?: number
  disabled?: boolean
  articleTitle?: string // 新增：文章标题，用于图片上传时创建文件夹
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入内容...',
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
const editorRef = ref<HTMLElement>()
const toolbarRef = ref<HTMLElement>()
const content = ref('')
const sourceContent = ref('')
const sourceMode = ref(false)
const showImageUpload = ref(false)
const uploadedImageUrls = ref<string[]>([])
const currentSelection = ref<Range | null>(null)
const isComposing = ref(false)

// 计算属性
const wordCount = computed(() => {
  let text = ''
  if (sourceMode.value) {
    // 源码模式：直接使用源码内容，但需要去除HTML标签
    text = sourceContent.value.replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ')
  } else {
    // 普通模式：使用content的纯文本版本
    const tempDiv = document.createElement('div')
    tempDiv.innerHTML = content.value || ''
    text = tempDiv.textContent || tempDiv.innerText || ''
  }
  // 去除空白字符后计算长度
  return text.replace(/\s/g, '').length
})

// 监听modelValue变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue !== content.value && !isComposing.value) {
      content.value = newValue || ''
      sourceContent.value = newValue || ''
      // 同步到编辑器DOM
      if (editorRef.value && !sourceMode.value) {
        nextTick(() => {
          if (editorRef.value && editorRef.value.innerHTML !== newValue) {
            editorRef.value.innerHTML = newValue || ''
          }
        })
      } else if (!editorRef.value) {
        // 如果editor还没有准备好，延迟同步
        let retryCount = 0
        const maxRetries = 10
        
        const retrySync = () => {
          if (editorRef.value && newValue) {
            editorRef.value.innerHTML = newValue || ''
            content.value = newValue || ''
            sourceContent.value = newValue || ''
          } else if (retryCount < maxRetries) {
            retryCount++
            setTimeout(retrySync, 50)
          }
        }
        
        setTimeout(retrySync, 100)
      }
    }
  },
  { immediate: true }
)

// 监听内容变化
watch(content, (newValue) => {
  emit('update:modelValue', newValue)
  emit('change', newValue)
})

// 检测工具栏是否处于sticky状态
const checkStickyState = () => {
  if (toolbarRef.value) {
    const rect = toolbarRef.value.getBoundingClientRect()
    const editorContainer = toolbarRef.value.parentElement
    
    if (editorContainer) {
      const containerRect = editorContainer.getBoundingClientRect()
      // 如果工具栏固定在容器顶部，说明处于sticky状态
      const isSticky = rect.top <= containerRect.top + 5
      
      if (isSticky) {
        toolbarRef.value.classList.add('is-sticky')
      } else {
        toolbarRef.value.classList.remove('is-sticky')
      }
    }
  }
}

// 检查sticky支持并提供降级方案
const checkStickySupport = () => {
  // 检查浏览器是否支持sticky定位
  const testElement = document.createElement('div')
  testElement.style.position = 'sticky'
  const supportsSticky = testElement.style.position === 'sticky'
  
  if (!supportsSticky) {
    console.warn('浏览器不支持sticky定位，工具栏可能无法固定')
  }
  
  return supportsSticky
}

// 生命周期
onMounted(() => {
  nextTick(() => {
    if (editorRef.value) {
      editorRef.value.addEventListener('focus', saveSelection)
      editorRef.value.addEventListener('mouseup', saveSelection)
      editorRef.value.addEventListener('keyup', saveSelection)
      
      // 如果有初始内容，设置到编辑器中
      if (props.modelValue) {
        editorRef.value.innerHTML = props.modelValue
        content.value = props.modelValue
        sourceContent.value = props.modelValue
      }
    }
    
    // 检查sticky支持
    checkStickySupport()
    
    // 监听滚动事件以检测sticky状态
    window.addEventListener('scroll', checkStickyState)
    
    // 延迟检测，确保页面完全加载
    setTimeout(checkStickyState, 100)
  })
})

onBeforeUnmount(() => {
  if (editorRef.value) {
    editorRef.value.removeEventListener('focus', saveSelection)
    editorRef.value.removeEventListener('mouseup', saveSelection)
    editorRef.value.removeEventListener('keyup', saveSelection)
  }
  // 清理事件监听器
  window.removeEventListener('scroll', checkStickyState)
})

// 方法
const handleInput = () => {
  if (editorRef.value && !isComposing.value) {
    content.value = editorRef.value.innerHTML
    // 触发字数统计更新（通过改变content.value会自动触发wordCount计算属性更新）
  }
}

const handleCompositionStart = () => {
  isComposing.value = true
}

const handleCompositionUpdate = () => {
  // 在输入过程中不更新内容，避免干扰输入
}

const handleCompositionEnd = () => {
  isComposing.value = false
  if (editorRef.value) {
    content.value = editorRef.value.innerHTML
  }
}

const handleSourceInput = () => {
  content.value = sourceContent.value
  // 源码模式下也会触发字数统计更新
}

const handleKeydown = (e: KeyboardEvent) => {
  // 处理Tab键
  if (e.key === 'Tab') {
    e.preventDefault()
    execCommand('insertHTML', '&nbsp;&nbsp;&nbsp;&nbsp;')
  }
  
  // 处理Enter键
  if (e.key === 'Enter' && e.ctrlKey) {
    e.preventDefault()
    execCommand('insertHTML', '<br><br>')
  }
}

const handlePaste = (e: ClipboardEvent) => {
  e.preventDefault()
  const text = e.clipboardData?.getData('text/plain') || ''
  execCommand('insertText', text)
}

const saveSelection = () => {
  const selection = window.getSelection()
  if (selection && selection.rangeCount > 0) {
    currentSelection.value = selection.getRangeAt(0)
  }
}

const restoreSelection = () => {
  if (currentSelection.value) {
    const selection = window.getSelection()
    selection?.removeAllRanges()
    selection?.addRange(currentSelection.value)
  }
}

const execCommand = (command: string, value?: string) => {
  if (props.disabled) return
  
  restoreSelection()
  document.execCommand(command, false, value)
  
  nextTick(() => {
    if (editorRef.value) {
      editorRef.value.focus()
      handleInput()
    }
  })
}

const isActive = (command: string): boolean => {
  try {
    return document.queryCommandState(command)
  } catch {
    return false
  }
}

const formatHeading = (e: Event) => {
  const target = e.target as HTMLSelectElement
  const value = target.value
  
  if (value) {
    execCommand('formatBlock', `<${value}>`)
  } else {
    execCommand('formatBlock', '<p>')
  }
  
  target.value = ''
}

const insertLink = () => {
  // 先检查是否有选中的文字
  const selectedText = window.getSelection()?.toString()
  
  const inputUrl = prompt('请输入链接地址:')
  if (inputUrl && inputUrl.trim()) {
    let url = inputUrl.trim()
    
    // 如果没有协议，自动添加 https://
    if (!url.match(/^https?:\/\//)) {
      // 检查是否是邮箱地址
      if (url.includes('@')) {
        url = `mailto:${url}`
      } else {
        url = `https://${url}`
      }
    }
    
    let linkText = selectedText
    
    // 如果没有选中文字，询问用户要显示的文字
    if (!selectedText || selectedText.trim() === '') {
      const displayText = prompt('请输入要显示的链接文字:', getDomainFromUrl(inputUrl.trim()))
      linkText = displayText && displayText.trim() ? displayText.trim() : getDomainFromUrl(inputUrl.trim())
    }
    
    execCommand('insertHTML', `<a href="${url}" target="_blank" rel="noopener noreferrer" class="editor-link">${linkText}</a>`)
  }
}

// 从URL提取域名作为默认显示文字
const getDomainFromUrl = (url: string): string => {
  try {
    // 处理不带协议的URL
    let processedUrl = url
    if (!url.match(/^https?:\/\//)) {
      processedUrl = `https://${url}`
    }
    
    const urlObj = new URL(processedUrl)
    let domain = urlObj.hostname
    
    // 移除 www. 前缀
    domain = domain.replace(/^www\./, '')
    
    // 对于常见网站，返回更友好的名称
    const friendlyNames: { [key: string]: string } = {
      'baidu.com': '百度',
      'google.com': 'Google',
      'github.com': 'GitHub',
      'stackoverflow.com': 'Stack Overflow',
      'zhihu.com': '知乎',
      'juejin.cn': '掘金',
      'csdn.net': 'CSDN',
      'bilibili.com': '哔哩哔哩',
      'youtube.com': 'YouTube',
      'twitter.com': 'Twitter',
      'weibo.com': '微博'
    }
    
    return friendlyNames[domain] || domain
  } catch {
    // 如果URL解析失败，返回原始文本
    return url
  }
}

const changeTextColor = (e: Event) => {
  const target = e.target as HTMLInputElement
  execCommand('foreColor', target.value)
}

const changeBackgroundColor = (e: Event) => {
  const target = e.target as HTMLInputElement
  execCommand('hiliteColor', target.value)
}

const handleImageUpload = (response: ImageUploadResponse | ImageUploadResponse[]) => {
  if (Array.isArray(response)) {
    uploadedImageUrls.value = response.map(img => img.public_url)
  } else {
    uploadedImageUrls.value = [response.public_url]
  }
}

const insertImage = () => {
  if (uploadedImageUrls.value.length > 0) {
    const imagesHTML = uploadedImageUrls.value.map(url => 
      `<img src="${url}" alt="插入的图片" style="max-width: 100%; height: auto; margin: 5px;" />`
    ).join('')
    execCommand('insertHTML', imagesHTML)
    showImageUpload.value = false
    uploadedImageUrls.value = []
  }
}

const toggleSourceMode = () => {
  if (sourceMode.value) {
    // 从源码模式切换到可视化模式
    content.value = sourceContent.value
    sourceMode.value = false
    nextTick(() => {
      if (editorRef.value) {
        editorRef.value.innerHTML = content.value
      }
    })
  } else {
    // 从可视化模式切换到源码模式
    sourceContent.value = content.value
    sourceMode.value = true
  }
}
</script>

<style scoped>
.rich-text-editor {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  max-height: 80vh; /* 限制编辑器最大高度 */
  overflow: hidden; /* 让内容区域可滚动 */
  display: flex;
  flex-direction: column;
}

.editor-content {
  flex: 1;
  overflow-y: auto; /* 内容区域可滚动 */
  max-height: calc(80vh - 60px); /* 减去工具栏高度 */
}

.editor-toolbar {
  @apply bg-white/95 backdrop-blur-sm z-50;
  position: sticky;
  top: 0;
  border-bottom: 1px solid #e5e7eb;
  border-radius: 8px 8px 0 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.2s ease;
  margin: 0;
  padding: 6px 12px;
  width: 100%;
  overflow: visible;
  min-height: 40px;
}

/* 工具栏处于sticky状态时的增强效果 */
.editor-toolbar.is-sticky {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .editor-toolbar {
    padding: 4px 8px;
    min-height: 36px;
  }
}

/* 确保工具栏内容不被遮挡 */
.editor-toolbar::before {
  content: '';
  position: absolute;
  top: -1px;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(0, 0, 0, 0.1), transparent);
  opacity: 0;
  transition: opacity 0.2s ease-in-out;
}

.editor-toolbar.is-sticky::before {
  opacity: 1;
}

.toolbar-btn {
  @apply px-1.5 py-1.5 sm:px-2 sm:py-1.5 rounded hover:bg-white/70 transition-all text-xs;
  border: 1px solid transparent;
  color: #6b7280;
  min-height: 28px;
  height: 28px;
  min-width: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  white-space: nowrap;
}

.toolbar-btn:hover {
  @apply bg-white/80 shadow-sm;
  border-color: rgba(0, 0, 0, 0.1);
}

.toolbar-btn.active {
  @apply bg-blue-100 border-blue-300 text-blue-700;
}

.toolbar-select {
  @apply px-2 py-1 rounded text-xs focus:outline-none focus:ring-1 focus:ring-blue-400;
  border: 1px solid transparent;
  background: rgba(255, 255, 255, 0.7);
  color: #6b7280;
  min-height: 28px;
  height: 28px;
  line-height: 1;
  white-space: nowrap;
  font-size: 11px;
}

.toolbar-select:hover {
  @apply bg-white/80 shadow-sm;
  border-color: rgba(0, 0, 0, 0.1);
}

.editor-area {
  min-height: 200px;
}

.source-area {
  min-height: 200px;
}

.editor-area:empty:before {
  content: attr(placeholder);
  color: #9ca3af;
  pointer-events: none;
}

/* 富文本内容样式 */
.editor-area h1 {
  @apply text-xl sm:text-2xl font-bold mb-3 sm:mb-4;
}

.editor-area h2 {
  @apply text-lg sm:text-xl font-bold mb-2 sm:mb-3;
}

.editor-area h3 {
  @apply text-base sm:text-lg font-bold mb-2;
}

.editor-area ul {
  @apply list-disc list-inside mb-3 sm:mb-4;
}

.editor-area ol {
  @apply list-decimal list-inside mb-3 sm:mb-4;
}

.editor-area p {
  @apply mb-2;
}

/* 使用:deep()来穿透scoped样式，作用于动态插入的链接 */
.editor-area :deep(a) {
  @apply text-blue-600;
  text-decoration: underline !important;
  text-decoration-color: rgba(37, 99, 235, 0.6) !important;
  text-underline-offset: 2px !important;
  color: #2563eb !important;
  cursor: pointer !important;
  transition: all 0.2s ease !important;
  font-weight: 500 !important;
  position: relative !important;
}

.editor-area :deep(a:hover) {
  @apply text-blue-800;
  background-color: rgba(37, 99, 235, 0.1) !important;
  text-decoration-color: #1e40af !important;
  border-radius: 3px !important;
  padding: 2px 4px !important;
  margin: 0 -2px !important;
}

/* 为链接添加小图标 */
.editor-area :deep(a.editor-link::after) {
  content: '🔗' !important;
  display: inline !important;
  margin-left: 2px !important;
  font-size: 0.8em !important;
  opacity: 0.7 !important;
  transition: opacity 0.2s ease !important;
}

.editor-area :deep(a.editor-link:hover::after) {
  opacity: 1 !important;
}

/* 外部链接图标的替代方案，使用CSS实现的小箭头 */
.editor-area :deep(a[target="_blank"]::before) {
  content: '' !important;
  display: inline-block !important;
  width: 0 !important;
  height: 0 !important;
  border-left: 3px solid currentColor !important;
  border-top: 3px solid transparent !important;
  border-bottom: 3px solid transparent !important;
  margin-right: 3px !important;
  opacity: 0.6 !important;
  vertical-align: middle !important;
  transition: opacity 0.2s ease !important;
}

.editor-area :deep(a[target="_blank"]:hover::before) {
  opacity: 1 !important;
}

.editor-area img {
  @apply max-w-full h-auto rounded;
}

/* 自定义tooltip样式 - 立即显示 */
[data-tooltip] {
  position: relative;
}

[data-tooltip]:hover::before {
  content: attr(data-tooltip);
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.9);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 11px;
  white-space: nowrap;
  z-index: 1000;
  opacity: 1;
  visibility: visible;
  transition: opacity 0.1s ease-in-out;
  margin-top: 4px;
  pointer-events: none;
}

[data-tooltip]:hover::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 4px solid transparent;
  border-bottom-color: rgba(0, 0, 0, 0.9);
  z-index: 1000;
  margin-top: -4px;
  pointer-events: none;
}

/* 隐藏默认tooltip */
[data-tooltip]::before,
[data-tooltip]::after {
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.1s ease-in-out, visibility 0.1s ease-in-out;
}

/* 移动端优化 */
@media (max-width: 640px) {
  .toolbar-btn {
    min-width: 24px;
    min-height: 24px;
    height: 24px;
    padding: 2px;
  }
  
  .toolbar-select {
    min-height: 24px;
    height: 24px;
    padding: 2px 4px;
    font-size: 10px;
  }
  
  /* 移动端不显示tooltip，避免干扰触摸操作 */
  [data-tooltip]:hover::before,
  [data-tooltip]:hover::after {
    display: none;
  }
}
</style>