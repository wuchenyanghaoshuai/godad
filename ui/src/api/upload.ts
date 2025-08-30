// 上传相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type {
  ImageUploadResponse,
  ApiResponse
} from './types'

// 上传API服务类
export class UploadApi {
  // 上传图片
  static async uploadImage(file: File): Promise<ApiResponse<ImageUploadResponse>> {
    const formData = new FormData()
    formData.append('file', file)  // 修复：后端期望的字段名是 'file'
    
    return http.post<ImageUploadResponse>(API_CONFIG.ENDPOINTS.UPLOAD.IMAGE, formData)
  }

  // 批量上传图片
  static async uploadImages(files: File[]): Promise<ApiResponse<ImageUploadResponse[]>> {
    // 由于后端只有单个图片上传接口，这里使用循环上传多个文件
    const uploadPromises = files.map(file => this.uploadImage(file))
    const results = await Promise.all(uploadPromises)
    
    // 提取所有上传结果的data部分
    const uploadData = results.map(result => result.data).filter(Boolean) as ImageUploadResponse[]
    
    return {
      code: 200,
      message: '批量上传成功',
      data: uploadData
    }
  }

  // 上传头像
  static async uploadAvatar(file: File): Promise<ApiResponse<ImageUploadResponse>> {
    const formData = new FormData()
    formData.append('file', file)  // 修复：后端期望的字段名是 'file' 而不是 'avatar'
    
    return http.post<ImageUploadResponse>('/upload/avatar', formData)
  }
}

// 上传工具函数
export class UploadUtils {
  // 验证图片文件
  static validateImageFile(file: File): { valid: boolean; error?: string } {
    // 检查文件类型
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
    if (!allowedTypes.includes(file.type)) {
      return {
        valid: false,
        error: '只支持 JPEG、PNG、GIF、WebP 格式的图片'
      }
    }

    // 检查文件大小（5MB）
    const maxSize = 5 * 1024 * 1024
    if (file.size > maxSize) {
      return {
        valid: false,
        error: '图片大小不能超过 5MB'
      }
    }

    return { valid: true }
  }

  // 压缩图片
  static compressImage(file: File, quality: number = 0.8): Promise<File> {
    return new Promise((resolve) => {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')!
      const img = new Image()

      img.onload = () => {
        // 计算压缩后的尺寸
        const maxWidth = 1920
        const maxHeight = 1080
        let { width, height } = img

        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width *= ratio
          height *= ratio
        }

        canvas.width = width
        canvas.height = height

        // 绘制压缩后的图片
        ctx.drawImage(img, 0, 0, width, height)

        // 转换为Blob
        canvas.toBlob(
          (blob) => {
            if (blob) {
              const compressedFile = new File([blob], file.name, {
                type: file.type,
                lastModified: Date.now()
              })
              resolve(compressedFile)
            } else {
              resolve(file)
            }
          },
          file.type,
          quality
        )
      }

      img.src = URL.createObjectURL(file)
    })
  }

  // 生成图片预览URL
  static createPreviewUrl(file: File): string {
    return URL.createObjectURL(file)
  }

  // 释放预览URL
  static revokePreviewUrl(url: string): void {
    URL.revokeObjectURL(url)
  }
}