<template>
  <!-- 头像上传弹窗 -->
  <div v-if="isVisible" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" @click="closeModal">
    <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4" @click.stop>
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-medium text-gray-900">更换头像</h3>
        <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
          <XIcon class="h-6 w-6" />
        </button>
      </div>

      <!-- 文件选择区域 -->
      <div 
        v-if="!selectedFile" 
        class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-gray-400 transition-colors"
        :class="{ 'border-pink-400 bg-pink-50': isDragOver }"
        @dragover.prevent="isDragOver = true"
        @dragleave.prevent="isDragOver = false"
        @drop.prevent="handleDrop"
      >
        <input
          ref="fileInput"
          type="file"
          accept="image/jpeg,image/jpg,image/png"
          @change="handleFileSelect"
          class="hidden"
        />
        <div class="space-y-2">
          <ImageIcon class="h-12 w-12 text-gray-400 mx-auto" />
          <div>
            <button
              @click="fileInput?.click()"
              class="text-pink-600 hover:text-pink-700 font-medium transition-colors"
            >
              选择图片
            </button>
            <span class="text-gray-500 mx-2">或拖拽到此处</span>
            <p class="text-sm text-gray-500 mt-1">支持 JPG、PNG 格式，最大 5MB</p>
            <p class="text-xs text-gray-400 mt-1">建议尺寸：200x200 像素</p>
          </div>
        </div>
      </div>

      <!-- 图片裁剪区域 -->
      <div v-if="selectedFile && !isUploading" class="space-y-4">
        <!-- 裁剪画布 -->
        <div class="relative bg-gray-100 rounded-lg overflow-hidden" style="height: 300px;">
          <canvas
            ref="cropCanvas"
            class="absolute inset-0 cursor-move"
            @mousedown="startCrop"
            @mousemove="updateCrop"
            @mouseup="endCrop"
            @mouseleave="endCrop"
          ></canvas>
        </div>

        <!-- 预览区域 -->
        <div class="flex items-center space-x-4">
          <div class="text-sm text-gray-600">预览:</div>
          <div class="w-16 h-16 rounded-full overflow-hidden border-2 border-gray-200">
            <canvas ref="previewCanvas" width="64" height="64" class="w-full h-full"></canvas>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex space-x-3">
          <button
            @click="resetSelection"
            class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50"
          >
            重新选择
          </button>
          <button
            @click="saveAvatar"
            class="flex-1 px-4 py-2 bg-pink-600 text-white rounded-md hover:bg-pink-700"
          >
            保存头像
          </button>
        </div>
      </div>

      <!-- 上传进度 -->
      <div v-if="isUploading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600"></div>
        <p class="mt-2 text-sm text-gray-600">正在上传...</p>
        <div class="w-full bg-gray-200 rounded-full h-2 mt-4">
          <div class="bg-pink-600 h-2 rounded-full transition-all duration-300" :style="{ width: uploadProgress + '%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, nextTick } from 'vue'
import { XIcon, ImageIcon } from 'lucide-vue-next'
import { UploadApi } from '../api/upload'
import type { ImageUploadResponse } from '../api/types'

export default defineComponent({
  name: 'AvatarModal',
  components: {
    XIcon,
    ImageIcon
  },
  props: {
    isVisible: {
      type: Boolean,
      required: true
    }
  },
  emits: ['close', 'success', 'error'],
  setup(props, { emit }) {
    // 响应式数据
    const fileInput = ref<HTMLInputElement>()
    const cropCanvas = ref<HTMLCanvasElement>()
    const previewCanvas = ref<HTMLCanvasElement>()
    const selectedFile = ref<File | null>(null)
    const isUploading = ref(false)
    const uploadProgress = ref(0)
    const isDragOver = ref(false)

    // 裁剪相关状态
    const cropState = ref({
      isDragging: false,
      startX: 0,
      startY: 0,
      cropX: 0,
      cropY: 0,
      cropSize: 200
    })

    const imageData = ref({
      img: null as HTMLImageElement | null,
      scale: 1,
      offsetX: 0,
      offsetY: 0
    })

    // 关闭弹窗
    const closeModal = () => {
      resetSelection()
      emit('close')
    }

    // 重置选择
    const resetSelection = () => {
      selectedFile.value = null
      isUploading.value = false
      uploadProgress.value = 0
      cropState.value = {
        isDragging: false,
        startX: 0,
        startY: 0,
        cropX: 0,
        cropY: 0,
        cropSize: 200
      }
      imageData.value = {
        img: null,
        scale: 1,
        offsetX: 0,
        offsetY: 0
      }
    }

    // 文件验证
    const validateFile = (file: File): boolean => {
      // 验证文件类型
      if (!['image/jpeg', 'image/jpg', 'image/png'].includes(file.type)) {
        emit('error', '请选择 JPG 或 PNG 格式的图片')
        return false
      }
      
      // 验证文件大小 (5MB)
      if (file.size > 5 * 1024 * 1024) {
        emit('error', '图片大小不能超过 5MB')
        return false
      }
      
      return true
    }

    // 处理文件
    const processFile = (file: File) => {
      if (validateFile(file)) {
        selectedFile.value = file
        loadImage(file)
      }
    }

    // 文件选择处理
    const handleFileSelect = (event: Event) => {
      const target = event.target as HTMLInputElement
      const file = target.files?.[0]
      
      if (!file) return
      
      processFile(file)
    }

    // 拖拽处理
    const handleDrop = (event: DragEvent) => {
      isDragOver.value = false
      
      const files = event.dataTransfer?.files
      if (!files || files.length === 0) return
      
      const file = files[0]
      processFile(file)
    }

    // 加载图片
    const loadImage = (file: File) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        const img = new Image()
        img.onload = () => {
          imageData.value.img = img
          nextTick(() => {
            initializeCrop()
          })
        }
        img.src = e.target?.result as string
      }
      reader.readAsDataURL(file)
    }

    // 初始化裁剪
    const initializeCrop = () => {
      if (!cropCanvas.value || !imageData.value.img) return
      
      const canvas = cropCanvas.value
      const img = imageData.value.img
      
      // 设置画布尺寸
      canvas.width = 400
      canvas.height = 300
      
      // 计算图片缩放比例
      const scaleX = canvas.width / img.width
      const scaleY = canvas.height / img.height
      imageData.value.scale = Math.min(scaleX, scaleY)
      
      // 计算图片居中位置
      const scaledWidth = img.width * imageData.value.scale
      const scaledHeight = img.height * imageData.value.scale
      imageData.value.offsetX = (canvas.width - scaledWidth) / 2
      imageData.value.offsetY = (canvas.height - scaledHeight) / 2
      
      // 初始化裁剪区域（居中）
      cropState.value.cropX = (canvas.width - cropState.value.cropSize) / 2
      cropState.value.cropY = (canvas.height - cropState.value.cropSize) / 2
      
      drawCanvas()
    }

    // 绘制画布
    const drawCanvas = () => {
      if (!cropCanvas.value || !imageData.value.img) return
      
      const canvas = cropCanvas.value
      const ctx = canvas.getContext('2d')!
      const img = imageData.value.img
      
      // 清空画布
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      
      // 绘制图片
      const scaledWidth = img.width * imageData.value.scale
      const scaledHeight = img.height * imageData.value.scale
      ctx.drawImage(
        img,
        imageData.value.offsetX,
        imageData.value.offsetY,
        scaledWidth,
        scaledHeight
      )
      
      // 绘制遮罩
      ctx.fillStyle = 'rgba(0, 0, 0, 0.5)'
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      
      // 清除裁剪区域的遮罩
      ctx.globalCompositeOperation = 'destination-out'
      ctx.fillRect(
        cropState.value.cropX,
        cropState.value.cropY,
        cropState.value.cropSize,
        cropState.value.cropSize
      )
      
      // 重置合成模式
      ctx.globalCompositeOperation = 'source-over'
      
      // 绘制裁剪框边框
      ctx.strokeStyle = '#ec4899'
      ctx.lineWidth = 2
      ctx.strokeRect(
        cropState.value.cropX,
        cropState.value.cropY,
        cropState.value.cropSize,
        cropState.value.cropSize
      )
      
      // 更新预览
      updatePreview()
    }

    // 更新预览
    const updatePreview = () => {
      if (!previewCanvas.value || !cropCanvas.value || !imageData.value.img) return
      
      const previewCtx = previewCanvas.value.getContext('2d')!
      const mainCanvas = cropCanvas.value
      
      // 清空预览画布
      previewCtx.clearRect(0, 0, 64, 64)
      
      // 从主画布裁剪区域复制到预览画布
      previewCtx.drawImage(
        mainCanvas,
        cropState.value.cropX,
        cropState.value.cropY,
        cropState.value.cropSize,
        cropState.value.cropSize,
        0,
        0,
        64,
        64
      )
    }

    // 开始裁剪拖拽
    const startCrop = (event: MouseEvent) => {
      const rect = cropCanvas.value!.getBoundingClientRect()
      const x = event.clientX - rect.left
      const y = event.clientY - rect.top
      
      // 检查是否点击在裁剪区域内
      if (
        x >= cropState.value.cropX &&
        x <= cropState.value.cropX + cropState.value.cropSize &&
        y >= cropState.value.cropY &&
        y <= cropState.value.cropY + cropState.value.cropSize
      ) {
        cropState.value.isDragging = true
        cropState.value.startX = x - cropState.value.cropX
        cropState.value.startY = y - cropState.value.cropY
      }
    }

    // 更新裁剪位置
    const updateCrop = (event: MouseEvent) => {
      if (!cropState.value.isDragging || !cropCanvas.value) return
      
      const rect = cropCanvas.value.getBoundingClientRect()
      const x = event.clientX - rect.left
      const y = event.clientY - rect.top
      
      // 计算新位置
      let newX = x - cropState.value.startX
      let newY = y - cropState.value.startY
      
      // 边界限制
      newX = Math.max(0, Math.min(newX, cropCanvas.value.width - cropState.value.cropSize))
      newY = Math.max(0, Math.min(newY, cropCanvas.value.height - cropState.value.cropSize))
      
      cropState.value.cropX = newX
      cropState.value.cropY = newY
      
      drawCanvas()
    }

    // 结束裁剪拖拽
    const endCrop = () => {
      cropState.value.isDragging = false
    }

    // 保存头像
    const saveAvatar = async () => {
      if (!selectedFile.value || !cropCanvas.value || !imageData.value.img) return
      
      try {
        isUploading.value = true
        uploadProgress.value = 0
        
        // 创建裁剪后的图片
        const croppedCanvas = document.createElement('canvas')
        croppedCanvas.width = 200
        croppedCanvas.height = 200
        const croppedCtx = croppedCanvas.getContext('2d')!
        
        // 从主画布复制裁剪区域
        croppedCtx.drawImage(
          cropCanvas.value,
          cropState.value.cropX,
          cropState.value.cropY,
          cropState.value.cropSize,
          cropState.value.cropSize,
          0,
          0,
          200,
          200
        )
        
        // 转换为Blob (使用PNG格式保持最佳质量)
        const blob = await new Promise<Blob>((resolve) => {
          croppedCanvas.toBlob((blob) => {
            resolve(blob!)
          }, 'image/png')
        })
        
        // 创建File对象
        const croppedFile = new File([blob], 'avatar.png', { type: 'image/png' })
        
        // 模拟上传进度
        const progressInterval = setInterval(() => {
          uploadProgress.value += 10
          if (uploadProgress.value >= 90) {
            clearInterval(progressInterval)
          }
        }, 100)
        
        // 上传头像
        const response = await UploadApi.uploadAvatar(croppedFile)
        uploadProgress.value = 100
        
        clearInterval(progressInterval)
        
        emit('success', response.data)
        closeModal()
        
      } catch (error) {
        console.error('上传失败:', error)
        emit('error', '上传失败，请重试')
      } finally {
        isUploading.value = false
        uploadProgress.value = 0
      }
    }

    return {
      fileInput,
      cropCanvas,
      previewCanvas,
      selectedFile,
      isUploading,
      uploadProgress,
      isDragOver,
      cropState,
      imageData,
      closeModal,
      resetSelection,
      validateFile,
      processFile,
      handleFileSelect,
      handleDrop,
      loadImage,
      initializeCrop,
      drawCanvas,
      updatePreview,
      startCrop,
      updateCrop,
      endCrop,
      saveAvatar
    }
  }
})
</script>

<style scoped>
.cursor-move {
  cursor: move;
}
</style>