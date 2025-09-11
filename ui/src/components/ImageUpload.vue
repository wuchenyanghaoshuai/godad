<template>
  <div class="image-upload">
    <!-- 上传区域 -->
    <div
      ref="uploadAreaRef"
      @click="triggerFileInput"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent="handleDrop"
      :class="[
        'border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors',
        isDragOver
          ? 'border-pink-500 bg-pink-50'
          : 'border-gray-300 hover:border-pink-400 hover:bg-gray-50',
        disabled && 'opacity-50 cursor-not-allowed'
      ]"
    >
      <input
        ref="fileInputRef"
        type="file"
        :accept="accept"
        :multiple="multiple"
        @change="handleFileSelect"
        class="hidden"
        :disabled="disabled"
      />
      
      <div v-if="!previewUrls.length">
        <UploadIcon class="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <p class="text-gray-600 mb-2">点击上传或拖拽图片到此处</p>
        <p class="text-sm text-gray-400">
          支持 {{ acceptText }}，单个文件不超过 {{ maxSizeText }}
        </p>
      </div>
      
      <!-- 预览区域 -->
      <div v-else class="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div
          v-for="(url, index) in previewUrls"
          :key="index"
          class="relative group"
        >
          <img
            :src="url"
            :alt="`预览图 ${index + 1}`"
            class="w-full h-24 object-cover rounded-lg"
          />
          <button
            @click.stop="removeImage(index)"
            class="absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <XIcon class="h-4 w-4" />
          </button>
        </div>
        
        <!-- 添加更多按钮 -->
        <div
          v-if="multiple && previewUrls.length < maxFiles"
          class="border-2 border-dashed border-gray-300 rounded-lg h-24 flex items-center justify-center cursor-pointer hover:border-pink-400 transition-colors"
          @click.stop="triggerFileInput"
        >
          <PlusIcon class="h-8 w-8 text-gray-400" />
        </div>
      </div>
    </div>
    
    <!-- 上传进度 -->
    <div v-if="uploading" class="mt-4">
      <div class="flex items-center justify-between text-sm text-gray-600 mb-2">
        <span>上传中...</span>
        <span>{{ uploadProgress }}%</span>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-2">
        <div
          class="bg-pink-500 h-2 rounded-full transition-all duration-300"
          :style="{ width: `${uploadProgress}%` }"
        ></div>
      </div>
    </div>
    
    <!-- 错误信息 -->
    <div v-if="errorMessage" class="mt-2 text-sm text-red-600">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { UploadIcon, XIcon, PlusIcon } from 'lucide-vue-next'
import { UploadApi, UploadUtils } from '@/api/upload'
import type { ImageUploadResponse } from '@/api/types'

// Props
interface Props {
  modelValue?: string | string[]
  multiple?: boolean
  maxFiles?: number
  maxSize?: number // MB
  accept?: string
  disabled?: boolean
  compress?: boolean
  quality?: number
  uploadType?: 'image' | 'avatar' // 新增：上传类型
  articleTitle?: string // 新增：文章标题
}

const props = withDefaults(defineProps<Props>(), {
  multiple: false,
  maxFiles: 5,
  maxSize: 5,
  accept: 'image/*',
  disabled: false,
  compress: true,
  quality: 0.8,
  uploadType: 'image'
})

// Emits
interface Emits {
  'update:modelValue': [value: string | string[]]
  'upload-success': [response: ImageUploadResponse | ImageUploadResponse[]]
  'upload-error': [error: string]
}

const emit = defineEmits<Emits>()

// 响应式数据
const fileInputRef = ref<HTMLInputElement>()
const uploadAreaRef = ref<HTMLElement>()
const isDragOver = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const errorMessage = ref('')
const previewUrls = ref<string[]>([])
const selectedFiles = ref<File[]>([])

// 计算属性
const acceptText = computed(() => {
  const types = props.accept.split(',').map(type => {
    if (type.includes('jpeg') || type.includes('jpg')) return 'JPEG'
    if (type.includes('png')) return 'PNG'
    if (type.includes('gif')) return 'GIF'
    if (type.includes('webp')) return 'WebP'
    return type.replace('image/', '').toUpperCase()
  })
  return types.join('、')
})

const maxSizeText = computed(() => {
  return props.maxSize >= 1 ? `${props.maxSize}MB` : `${props.maxSize * 1024}KB`
})

// 监听modelValue变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue) {
      if (Array.isArray(newValue)) {
        previewUrls.value = [...newValue]
      } else {
        previewUrls.value = [newValue]
      }
    } else {
      previewUrls.value = []
    }
  },
  { immediate: true }
)

// 方法
const triggerFileInput = () => {
  if (!props.disabled) {
    fileInputRef.value?.click()
  }
}

const handleDragOver = (e: DragEvent) => {
  if (!props.disabled) {
    isDragOver.value = true
  }
}

const handleDragLeave = (e: DragEvent) => {
  isDragOver.value = false
}

const handleDrop = (e: DragEvent) => {
  isDragOver.value = false
  if (props.disabled) return
  
  const files = Array.from(e.dataTransfer?.files || [])
  processFiles(files)
}

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  const files = Array.from(target.files || [])
  processFiles(files)
}

const processFiles = async (files: File[]) => {
  errorMessage.value = ''
  
  // 过滤图片文件
  const imageFiles = files.filter(file => file.type.startsWith('image/'))
  
  if (imageFiles.length === 0) {
    errorMessage.value = '请选择有效的图片文件'
    return
  }
  
  // 检查文件数量限制
  const totalFiles = props.multiple 
    ? Math.min(previewUrls.value.length + imageFiles.length, props.maxFiles)
    : 1
  
  const filesToProcess = imageFiles.slice(0, totalFiles - previewUrls.value.length)
  
  // 验证文件
  for (const file of filesToProcess) {
    const validation = UploadUtils.validateImageFile(file)
    if (!validation.valid) {
      errorMessage.value = validation.error || '文件验证失败'
      return
    }
  }
  
  // 处理文件
  const processedFiles: File[] = []
  
  for (const file of filesToProcess) {
    try {
      const processedFile = props.compress 
        ? await UploadUtils.compressImage(file, props.quality)
        : file
      processedFiles.push(processedFile)
    } catch (error) {
      console.error('图片处理失败:', error)
      processedFiles.push(file)
    }
  }
  
  // 生成预览
  const newPreviewUrls = processedFiles.map(file => UploadUtils.createPreviewUrl(file))
  
  if (props.multiple) {
    previewUrls.value = [...previewUrls.value, ...newPreviewUrls]
    selectedFiles.value = [...selectedFiles.value, ...processedFiles]
  } else {
    // 清理旧的预览URL
    previewUrls.value.forEach(url => {
      if (url.startsWith('blob:')) {
        UploadUtils.revokePreviewUrl(url)
      }
    })
    previewUrls.value = newPreviewUrls
    selectedFiles.value = processedFiles
  }
  
  // 自动上传
  await uploadFiles()
}

const uploadFiles = async () => {
  if (selectedFiles.value.length === 0) return
  
  uploading.value = true
  uploadProgress.value = 0
  errorMessage.value = ''
  
  try {
    if (props.multiple && selectedFiles.value.length > 1) {
      // 批量上传
      const response = await UploadApi.uploadImages(selectedFiles.value, props.articleTitle)
      if (response.data && Array.isArray(response.data)) {
        const urls = response.data.map(item => item.public_url)
        const currentUrls = Array.isArray(props.modelValue) ? props.modelValue : []
        const newUrls = [...currentUrls, ...urls]
        
        emit('update:modelValue', newUrls)
        emit('upload-success', response.data)
        
        // 更新预览URL - 清理blob URL并替换为服务器URL
        const oldBlobUrls = previewUrls.value.slice(-selectedFiles.value.length)
        oldBlobUrls.forEach(url => {
          if (url.startsWith('blob:')) {
            UploadUtils.revokePreviewUrl(url)
          }
        })
        previewUrls.value = [...previewUrls.value.slice(0, -selectedFiles.value.length), ...urls]
      }
    } else {
      // 单个上传
      const file = selectedFiles.value[0]
      // 根据上传类型调用不同的API方法
      const response = props.uploadType === 'avatar' 
        ? await UploadApi.uploadAvatar(file)
        : await UploadApi.uploadImage(file, props.articleTitle)
      
      if (response.data) {
        const url = response.data.public_url
        
        if (props.multiple) {
          const currentUrls = Array.isArray(props.modelValue) ? props.modelValue : []
          const newUrls = [...currentUrls, url]
          emit('update:modelValue', newUrls)
        } else {
          emit('update:modelValue', url)
        }
        
        emit('upload-success', response.data)
        
        // 更新预览URL - 清理blob URL并替换为服务器URL
        const lastIndex = previewUrls.value.length - 1
        const oldBlobUrl = previewUrls.value[lastIndex]
        if (oldBlobUrl && oldBlobUrl.startsWith('blob:')) {
          UploadUtils.revokePreviewUrl(oldBlobUrl)
        }
        previewUrls.value.splice(lastIndex, 1, url)
      }
    }
    
    uploadProgress.value = 100
    selectedFiles.value = []
  } catch (error: any) {
    console.error('上传失败:', error)
    const errorMsg = error.response?.data?.message || error.message || '上传失败，请重试'
    errorMessage.value = errorMsg
    emit('upload-error', errorMsg)
    
    // 移除失败的预览
    const failedCount = selectedFiles.value.length
    previewUrls.value = previewUrls.value.slice(0, -failedCount)
    selectedFiles.value = []
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

const removeImage = (index: number) => {
  const url = previewUrls.value[index]
  
  // 清理blob URL
  if (url.startsWith('blob:')) {
    UploadUtils.revokePreviewUrl(url)
  }
  
  previewUrls.value.splice(index, 1)
  
  if (props.multiple) {
    const currentUrls = Array.isArray(props.modelValue) ? [...props.modelValue] : []
    currentUrls.splice(index, 1)
    emit('update:modelValue', currentUrls)
  } else {
    emit('update:modelValue', '')
  }
}
</script>

<style scoped>
.image-upload {
  @apply w-full;
}
</style>