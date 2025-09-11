<template>
  <div>
    <div ref="toolbarRef" class="editor-toolbar px-3 py-3 sm:px-4 sm:py-3 flex flex-wrap gap-2 sm:gap-2">
      <!-- 格式化按钮 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="execCommand('bold')"
          :class="['toolbar-btn', isActive('bold') && 'active']"
          title="粗体"
        >
          <BoldIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="execCommand('italic')"
          :class="['toolbar-btn', isActive('italic') && 'active']"
          title="斜体"
        >
          <ItalicIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="execCommand('underline')"
          :class="['toolbar-btn', isActive('underline') && 'active']"
          title="下划线"
        >
          <UnderlineIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="execCommand('strikeThrough')"
          :class="['toolbar-btn', isActive('strikeThrough') && 'active']"
          title="删除线"
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
            title="文字颜色"
            value="#000000"
          />
        </div>
        <div class="relative">
          <input
            type="color"
            @change="changeBackgroundColor($event)"
            class="w-8 h-6 border rounded cursor-pointer"
            title="背景颜色"
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
          title="左对齐"
        >
          <span class="text-xs">⬅</span>
        </button>
        <button
          type="button"
          @click="execCommand('justifyCenter')"
          :class="['toolbar-btn', isActive('justifyCenter') && 'active']"
          title="居中"
        >
          <span class="text-xs">⬌</span>
        </button>
        <button
          type="button"
          @click="execCommand('justifyRight')"
          :class="['toolbar-btn', isActive('justifyRight') && 'active']"
          title="右对齐"
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
          title="标题"
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
          title="无序列表"
        >
          <ListIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="execCommand('insertOrderedList')"
          :class="['toolbar-btn', isActive('insertOrderedList') && 'active']"
          title="有序列表"
        >
          <ListOrderedIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
      </div>
      
      <div class="border-l border-gray-300 mx-0.5 sm:mx-1 hidden sm:block"></div>
      
      <!-- 链接和图片 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="insertLink"
          class="toolbar-btn"
          title="插入链接"
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
      
      <!-- 其他功能 -->
      <div class="flex gap-0.5 sm:gap-1">
        <button
          type="button"
          @click="execCommand('removeFormat')"
          class="toolbar-btn"
          title="清除格式"
        >
          <EraserIcon class="h-3 w-3 sm:h-4 sm:w-4" />
        </button>
        <button
          type="button"
          @click="toggleSourceMode"
          :class="['toolbar-btn', sourceMode && 'active']"
          title="源码模式"
        >
          <CodeIcon class="h-3 w-3 sm:h-4 sm:w-4" />
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

// 动态调整工具栏位置
const adjustToolbarPosition = () => {
  if (toolbarRef.value) {
    // 获取导航栏的实际高度
    const navbar = document.querySelector('nav')
    const navbarHeight = navbar ? navbar.offsetHeight : 64
    
    // 设置CSS变量，工具栏紧贴导航栏底部
    toolbarRef.value.style.setProperty('--navbar-height', `${navbarHeight}px`)
  }
}

// 检测工具栏是否处于sticky状态
const checkStickyState = () => {
  if (toolbarRef.value) {
    const rect = toolbarRef.value.getBoundingClientRect()
    const navbar = document.querySelector('nav')
    const navbarHeight = navbar ? navbar.offsetHeight : 64
    
    // 如果工具栏距离顶部的距离等于导航栏高度，说明处于sticky状态
    const isSticky = Math.abs(rect.top - navbarHeight) < 2
    
    if (isSticky) {
      toolbarRef.value.classList.add('is-sticky')
    } else {
      toolbarRef.value.classList.remove('is-sticky')
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
    
    // 动态调整工具栏位置
    adjustToolbarPosition()
    
    // 监听窗口大小变化和滚动事件
    window.addEventListener('resize', adjustToolbarPosition)
    window.addEventListener('scroll', adjustToolbarPosition)
    
    // 延迟调整，确保页面完全加载
    setTimeout(adjustToolbarPosition, 100)
    setTimeout(adjustToolbarPosition, 500)
  })
})

onBeforeUnmount(() => {
  if (editorRef.value) {
    editorRef.value.removeEventListener('focus', saveSelection)
    editorRef.value.removeEventListener('mouseup', saveSelection)
    editorRef.value.removeEventListener('keyup', saveSelection)
  }
  // 清理事件监听器
  window.removeEventListener('resize', adjustToolbarPosition)
  window.removeEventListener('scroll', adjustToolbarPosition)
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
    
    const text = window.getSelection()?.toString() || inputUrl.trim()
    execCommand('insertHTML', `<a href="${url}" target="_blank" rel="noopener noreferrer">${text}</a>`)
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
  /* overflow: hidden; 移除这个，它阻止了sticky定位 */
}

.editor-toolbar {
  @apply bg-gray-50/90 backdrop-blur-sm z-40;
  position: sticky;
  top: var(--navbar-height, 64px);
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
  margin: 0 -24px; /* 负边距突破父容器的padding限制 */
  padding: 6px 24px;
  width: calc(100% + 48px); /* 补偿负边距 */
  /* 确保内容不被裁剪 */
  overflow: visible;
  min-height: 48px;
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
    /* 移动端使用相同的navbar-height变量 */
    top: var(--navbar-height, 56px);
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
  @apply px-2 py-2 sm:px-3 sm:py-2 rounded hover:bg-white/70 transition-all text-xs sm:text-sm;
  border: 1px solid transparent;
  color: #6b7280;
  min-height: 36px; /* 增加高度确保文本不被截断 */
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1.2; /* 确保文本行高正常 */
  white-space: nowrap; /* 防止文本换行 */
}

.toolbar-btn:hover {
  @apply bg-white/80 shadow-sm;
  border-color: rgba(0, 0, 0, 0.1);
}

.toolbar-btn.active {
  @apply bg-blue-100 border-blue-300 text-blue-700;
}

.toolbar-select {
  @apply px-2 py-2 sm:px-3 sm:py-2 rounded text-xs sm:text-sm focus:outline-none focus:ring-1 focus:ring-blue-400;
  border: 1px solid transparent;
  background: rgba(255, 255, 255, 0.7);
  color: #6b7280;
  min-height: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  line-height: 1.2;
  white-space: nowrap;
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

.editor-area a {
  @apply text-blue-600 underline;
}

.editor-area img {
  @apply max-w-full h-auto rounded;
}

/* 移动端优化 */
@media (max-width: 640px) {
  .toolbar-btn {
    @apply min-w-[36px] min-h-[36px] flex items-center justify-center;
    height: 36px;
  }
  
  .toolbar-select {
    @apply min-h-[36px];
    height: 36px;
  }
}
</style>