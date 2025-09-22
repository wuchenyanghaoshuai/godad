<template>
  <div class="avatar-upload">
    <!-- 当前头像显示 -->
    <div class="flex items-center space-x-6 mb-4">
      <div class="w-20 h-20">
        <img
          v-if="currentAvatar"
          :src="currentAvatar"
          alt="当前头像"
          class="w-20 h-20 rounded-full object-cover border-2 border-gray-200"
        />
        <div
          v-else
          class="w-20 h-20 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center"
        >
          <span class="text-xl font-bold text-white">
            {{ username?.charAt(0).toUpperCase() || 'U' }}
          </span>
        </div>
      </div>
      <div class="flex-1">
        <h3 class="text-sm font-medium text-gray-700 mb-1">当前头像</h3>
        <p class="text-xs text-gray-500">点击下方上传新头像</p>
      </div>
    </div>

    <!-- 上传区域 -->
    <div
      ref="dropZone"
      class="upload-zone border-2 border-dashed border-gray-300 rounded-lg p-6 text-center transition-colors"
      :class="{
        'border-pink-500 bg-pink-50': isDragOver,
        'border-gray-300': !isDragOver
      }"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent="handleDrop"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        accept="image/jpeg,image/png"
        class="hidden"
        @change="handleFileSelect"
      />
      
      <div v-if="!selectedFile" class="cursor-pointer">
        <div class="w-12 h-12 mx-auto mb-4 text-gray-400">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
        </div>
        <p class="text-sm text-gray-600 mb-2">拖拽图片到此处或点击上传</p>
        <p class="text-xs text-gray-500">支持 JPG、PNG 格式，文件大小不超过 5MB</p>
      </div>

      <!-- 图片预览和裁剪 -->
      <div v-else class="space-y-4">
        <!-- 原图预览 -->
        <div class="relative">
          <img
            ref="imagePreview"
            :src="previewUrl"
            alt="预览图片"
            class="max-w-full max-h-64 mx-auto rounded"
            @load="initCrop"
          />
          
          <!-- 裁剪框 -->
          <div
            v-if="showCropBox"
            ref="cropBox"
            class="absolute border-2 border-pink-500 bg-pink-500 bg-opacity-20 cursor-move"
            :style="cropBoxStyle"
            @mousedown="startDrag"
          >
            <!-- 裁剪框调整手柄 -->
            <div class="absolute -top-1 -left-1 w-3 h-3 bg-pink-500 cursor-nw-resize" @mousedown.stop="startResize('nw')"></div>
            <div class="absolute -top-1 -right-1 w-3 h-3 bg-pink-500 cursor-ne-resize" @mousedown.stop="startResize('ne')"></div>
            <div class="absolute -bottom-1 -left-1 w-3 h-3 bg-pink-500 cursor-sw-resize" @mousedown.stop="startResize('sw')"></div>
            <div class="absolute -bottom-1 -right-1 w-3 h-3 bg-pink-500 cursor-se-resize" @mousedown.stop="startResize('se')"></div>
          </div>
        </div>

        <!-- 裁剪后预览 -->
        <div v-if="croppedPreview" class="text-center">
          <h4 class="text-sm font-medium text-gray-700 mb-2">裁剪预览</h4>
          <div class="inline-block">
            <img
              :src="croppedPreview"
              alt="裁剪预览"
              class="w-20 h-20 rounded-full object-cover border-2 border-gray-200"
            />
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex justify-center space-x-3">
          <button
            type="button"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            @click="resetUpload"
          >
            重新选择
          </button>
          <button
            type="button"
            class="px-4 py-2 text-sm font-medium text-white bg-pink-500 rounded-md hover:bg-pink-600 disabled:opacity-50"
            :disabled="isUploading"
            @click="uploadAvatar"
          >
            {{ isUploading ? '上传中...' : '保存头像' }}
          </button>
        </div>

        <!-- 上传进度 -->
        <div v-if="isUploading" class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="bg-pink-500 h-2 rounded-full transition-all duration-300"
            :style="{ width: uploadProgress + '%' }"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, nextTick } from 'vue'
import { UploadApi } from '../api/upload'
// type import removed (unused)

export default defineComponent({
  name: 'AvatarUpload',
  props: {
    modelValue: {
      type: String,
      default: ''
    },
    username: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue', 'upload-success', 'upload-error'],
  setup(props, { emit }) {
    // 响应式数据
    const fileInput = ref<HTMLInputElement>()
    const dropZone = ref<HTMLDivElement>()
    const imagePreview = ref<HTMLImageElement>()
    const cropBox = ref<HTMLDivElement>()
    
    const selectedFile = ref<File | null>(null)
    const previewUrl = ref('')
    const isDragOver = ref(false)
    const isUploading = ref(false)
    const uploadProgress = ref(0)
    const showCropBox = ref(false)
    const croppedPreview = ref('')
    
    // 裁剪相关数据
    const cropData = ref({
      x: 0,
      y: 0,
      width: 100,
      height: 100
    })
    
    const isDragging = ref(false)
    const isResizing = ref(false)
    const resizeDirection = ref('')
    const dragStart = ref({ x: 0, y: 0 })
    
    // 计算属性
    const currentAvatar = computed(() => props.modelValue)
    
    const cropBoxStyle = computed(() => ({
      left: cropData.value.x + 'px',
      top: cropData.value.y + 'px',
      width: cropData.value.width + 'px',
      height: cropData.value.height + 'px'
    }))
    
    // 文件处理方法
    const validateFile = (file: File): boolean => {
      if (!file.type.match(/^image\/(jpeg|png)$/)) {
        emit('upload-error', '只支持 JPG 和 PNG 格式的图片')
        return false
      }
      
      if (file.size > 5 * 1024 * 1024) {
        emit('upload-error', '文件大小不能超过 5MB')
        return false
      }
      
      return true
    }
    
    const processFile = (file: File) => {
      if (!validateFile(file)) return
      
      selectedFile.value = file
      const reader = new FileReader()
      reader.onload = (e) => {
        previewUrl.value = e.target?.result as string
      }
      reader.readAsDataURL(file)
    }
    
    // 拖拽处理
    const handleDragOver = (e: DragEvent) => {
      e.preventDefault()
      isDragOver.value = true
    }
    
    const handleDragLeave = (e: DragEvent) => {
      e.preventDefault()
      isDragOver.value = false
    }
    
    const handleDrop = (e: DragEvent) => {
      e.preventDefault()
      isDragOver.value = false
      
      const files = e.dataTransfer?.files
      if (files && files.length > 0) {
        processFile(files[0])
      }
    }
    
    // 文件选择
    const triggerFileInput = () => {
      if (!selectedFile.value) {
        fileInput.value?.click()
      }
    }
    
    const handleFileSelect = (e: Event) => {
      const target = e.target as HTMLInputElement
      const files = target.files
      if (files && files.length > 0) {
        processFile(files[0])
      }
    }
    
    // 裁剪功能
    const initCrop = async () => {
      await nextTick()
      if (imagePreview.value) {
        const rect = imagePreview.value.getBoundingClientRect()
        const size = Math.min(rect.width, rect.height) * 0.6
        
        cropData.value = {
          x: (rect.width - size) / 2,
          y: (rect.height - size) / 2,
          width: size,
          height: size
        }
        
        showCropBox.value = true
        updateCroppedPreview()
      }
    }
    
    const updateCroppedPreview = () => {
      if (!imagePreview.value || !selectedFile.value) return
      
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d', { 
        colorSpace: 'srgb',
        alpha: false
      })
      if (!ctx) return
      
      // 设置图像渲染选项以保持颜色准确性
      ctx.imageSmoothingEnabled = true
      ctx.imageSmoothingQuality = 'high'
      
      const img = imagePreview.value
      const rect = img.getBoundingClientRect()
      
      // 计算缩放比例
      const scaleX = img.naturalWidth / rect.width
      const scaleY = img.naturalHeight / rect.height
      
      // 设置画布大小
      canvas.width = cropData.value.width * scaleX
      canvas.height = cropData.value.height * scaleY
      
      // 绘制裁剪后的图片
      ctx.drawImage(
        img,
        cropData.value.x * scaleX,
        cropData.value.y * scaleY,
        cropData.value.width * scaleX,
        cropData.value.height * scaleY,
        0,
        0,
        canvas.width,
        canvas.height
      )
      
      // 使用更高的JPEG质量以减少颜色失真
      croppedPreview.value = canvas.toDataURL('image/jpeg', 0.95)
    }
    
    // 拖拽和调整大小
    const startDrag = (e: MouseEvent) => {
      isDragging.value = true
      dragStart.value = { x: e.clientX - cropData.value.x, y: e.clientY - cropData.value.y }
      
      const handleMouseMove = (e: MouseEvent) => {
        if (isDragging.value && imagePreview.value) {
          const rect = imagePreview.value.getBoundingClientRect()
          const newX = Math.max(0, Math.min(e.clientX - dragStart.value.x, rect.width - cropData.value.width))
          const newY = Math.max(0, Math.min(e.clientY - dragStart.value.y, rect.height - cropData.value.height))
          
          cropData.value.x = newX
          cropData.value.y = newY
          updateCroppedPreview()
        }
      }
      
      const handleMouseUp = () => {
        isDragging.value = false
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
      }
      
      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
    }
    
    const startResize = (direction: string) => {
      isResizing.value = true
      resizeDirection.value = direction
      
      const handleMouseMove = (e: MouseEvent) => {
        if (isResizing.value && imagePreview.value) {
          const rect = imagePreview.value.getBoundingClientRect()
          // 简化的调整大小逻辑
          const minSize = 50
          const maxSize = Math.min(rect.width, rect.height)
          
          if (direction.includes('e')) {
            const newWidth = Math.max(minSize, Math.min(e.clientX - cropData.value.x, maxSize - cropData.value.x))
            cropData.value.width = newWidth
            cropData.value.height = newWidth // 保持正方形
          }
          
          updateCroppedPreview()
        }
      }
      
      const handleMouseUp = () => {
        isResizing.value = false
        resizeDirection.value = ''
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
      }
      
      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
    }
    
    // 上传功能
    const uploadAvatar = async () => {
      if (!selectedFile.value || !croppedPreview.value) {
        emit('upload-error', '请先选择并裁剪图片')
        return
      }
      
      try {
        isUploading.value = true
        uploadProgress.value = 0
        
        // 将裁剪后的图片转换为文件
        const canvas = document.createElement('canvas')
        const ctx = canvas.getContext('2d', { 
          colorSpace: 'srgb',
          alpha: false
        })
        if (!ctx) {
          throw new Error('无法创建画布')
        }
        
        // 设置图像渲染选项以保持颜色准确性
        ctx.imageSmoothingEnabled = true
        ctx.imageSmoothingQuality = 'high'
        
        // 使用Promise包装图片加载，确保异步处理正确
        const processImage = () => {
          return new Promise<void>((resolve, reject) => {
            const img = new Image()
            img.onload = async () => {
              try {
                canvas.width = 200
                canvas.height = 200
                ctx.drawImage(img, 0, 0, 200, 200)
                
                // 使用Promise包装toBlob，确保异步处理正确
                const blob = await new Promise<Blob | null>((blobResolve) => {
                  canvas.toBlob(blobResolve, 'image/jpeg', 0.95)
                })
                
                if (!blob) {
                  reject(new Error('图片处理失败'))
                  return
                }
                
                const croppedFile = new File([blob], selectedFile.value!.name, {
                  type: 'image/jpeg',
                  lastModified: Date.now()
                })
                
                // 模拟上传进度
                const progressInterval = setInterval(() => {
                  uploadProgress.value += 10
                  if (uploadProgress.value >= 90) {
                    clearInterval(progressInterval)
                  }
                }, 100)
                
                try {
                  const response = await UploadApi.uploadAvatar(croppedFile)
                  uploadProgress.value = 100
                  
                  // 确保事件正确触发
                  emit('update:modelValue', response.data.url)
                  emit('upload-success', response.data)
                  
                  // 延迟重置，确保用户能看到成功状态
                  setTimeout(() => {
                    resetUpload()
                  }, 500)
                  
                  resolve()
                } catch (error) {
                  console.error('上传失败:', error)
                  emit('upload-error', '上传失败，请重试')
                  reject(error)
                } finally {
                  clearInterval(progressInterval)
                }
              } catch (error) {
                reject(error)
              }
            }
            
            img.onerror = () => {
              reject(new Error('图片加载失败'))
            }
            
            img.src = croppedPreview.value
          })
        }
        
        await processImage()
        
      } catch (error) {
        console.error('处理图片失败:', error)
        emit('upload-error', '处理图片失败，请重试')
      } finally {
        isUploading.value = false
        uploadProgress.value = 0
      }
    }
    
    // 重置上传
    const resetUpload = () => {
      selectedFile.value = null
      previewUrl.value = ''
      croppedPreview.value = ''
      showCropBox.value = false
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }
    
    return {
      // 模板引用
      fileInput,
      dropZone,
      imagePreview,
      cropBox,
      
      // 响应式数据
      selectedFile,
      previewUrl,
      isDragOver,
      isUploading,
      uploadProgress,
      showCropBox,
      croppedPreview,
      
      // 计算属性
      currentAvatar,
      cropBoxStyle,
      
      // 方法
      handleDragOver,
      handleDragLeave,
      handleDrop,
      triggerFileInput,
      handleFileSelect,
      initCrop,
      startDrag,
      startResize,
      uploadAvatar,
      resetUpload
    }
  }
})
</script>

<style scoped>
.avatar-upload {
  max-width: 500px;
}

.upload-zone {
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-zone:hover {
  border-color: #ec4899;
  background-color: #fdf2f8;
}
</style>
