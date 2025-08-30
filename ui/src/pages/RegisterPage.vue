<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 to-orange-50 flex items-center justify-center px-4">
    <div class="max-w-md w-full space-y-8">
      <!-- 头部 -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">加入GoDad</h1>
        <p class="text-gray-600">创建您的育儿知识分享账户</p>
      </div>

      <!-- 注册表单 -->
      <div class="bg-white rounded-2xl shadow-xl p-8">
        <form @submit.prevent="handleRegister" class="space-y-6">
          <!-- 用户名输入 -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
              用户名
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors"
              placeholder="请输入用户名"
              :disabled="isLoading"
            />
          </div>

          <!-- 邮箱输入 -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              邮箱地址
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors"
              placeholder="请输入您的邮箱"
              :disabled="isLoading"
            />
          </div>

          <!-- 昵称输入 -->
          <div>
            <label for="nickname" class="block text-sm font-medium text-gray-700 mb-2">
              昵称
            </label>
            <div class="flex gap-2">
              <input
                id="nickname"
                v-model="form.nickname"
                type="text"
                required
                class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors"
                placeholder="请输入您的昵称"
                :disabled="isLoading"
                @input="debouncedValidateNickname(form.nickname)"
              />
              <button
                type="button"
                @click="generateRandomNickname"
                :disabled="isLoading || isGeneratingNickname"
                class="px-4 py-3 bg-gradient-to-r from-purple-500 to-pink-500 text-white rounded-lg hover:from-purple-600 hover:to-pink-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all whitespace-nowrap"
                title="随机生成一个可爱的昵称"
              >
                <span v-if="isGeneratingNickname" class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  生成中
                </span>
                <span v-else>🎲 随机昵称</span>
              </button>
            </div>
            <!-- 昵称验证状态显示 -->
            <div v-if="isCheckingNickname || nicknameValidationMessage" class="mt-2">
              <div v-if="isCheckingNickname" class="flex items-center text-sm text-gray-500">
                <span class="animate-spin mr-2">⟳</span>
                检查昵称可用性...
              </div>
              <div v-else-if="nicknameValidationMessage" class="text-sm" :class="{
                'text-green-600': nicknameValidationMessage === '昵称可用',
                'text-red-600': nicknameValidationMessage !== '昵称可用'
              }">
                {{ nicknameValidationMessage }}
              </div>
            </div>
            <p class="text-xs text-gray-500 mt-1">点击随机昵称按钮生成一个可爱的昵称，或手动输入您喜欢的昵称</p>
          </div>

          <!-- 密码输入 -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              密码
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors pr-12"
                placeholder="请输入密码（至少8位，包含大小写字母）"
                :disabled="isLoading"
                minlength="8"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                :disabled="isLoading"
              >
                <EyeIcon v-if="!showPassword" class="h-5 w-5" />
                <EyeOffIcon v-else class="h-5 w-5" />
              </button>
            </div>
            <!-- 密码要求提示 -->
            <div class="mt-2 text-xs space-y-1" v-if="form.password">
              <div class="flex items-center space-x-2">
                <span :class="passwordStrength.hasMinLength ? 'text-green-600' : 'text-gray-400'">
                  {{ passwordStrength.hasMinLength ? '✓' : '○' }}
                </span>
                <span :class="passwordStrength.hasMinLength ? 'text-green-600' : 'text-gray-400'">
                  至少8位字符
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <span :class="passwordStrength.hasLowercase ? 'text-green-600' : 'text-gray-400'">
                  {{ passwordStrength.hasLowercase ? '✓' : '○' }}
                </span>
                <span :class="passwordStrength.hasLowercase ? 'text-green-600' : 'text-gray-400'">
                  包含至少一个小写字母
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <span :class="passwordStrength.hasUppercase ? 'text-green-600' : 'text-gray-400'">
                  {{ passwordStrength.hasUppercase ? '✓' : '○' }}
                </span>
                <span :class="passwordStrength.hasUppercase ? 'text-green-600' : 'text-gray-400'">
                  包含至少一个大写字母
                </span>
              </div>
            </div>
            <div v-else class="mt-2 text-xs text-gray-500">
              <ul class="list-disc list-inside space-y-1">
                <li>至少8位字符</li>
                <li>包含至少一个小写字母</li>
                <li>包含至少一个大写字母</li>
              </ul>
            </div>
          </div>

          <!-- 确认密码输入 -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              确认密码
            </label>
            <div class="relative">
              <input
                id="confirmPassword"
                v-model="form.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors pr-12"
                placeholder="请再次输入密码"
                :disabled="isLoading"
              />
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                :disabled="isLoading"
              >
                <EyeIcon v-if="!showConfirmPassword" class="h-5 w-5" />
                <EyeOffIcon v-else class="h-5 w-5" />
              </button>
            </div>
          </div>

          <!-- 服务条款 -->
          <div class="flex items-start">
            <input
              id="agree"
              v-model="form.agree"
              type="checkbox"
              required
              class="h-4 w-4 text-pink-600 focus:ring-pink-500 border-gray-300 rounded mt-1"
            />
            <label for="agree" class="ml-2 block text-sm text-gray-700">
              我已阅读并同意
              <a href="#" class="text-pink-600 hover:text-pink-500">服务条款</a>
              和
              <a href="#" class="text-pink-600 hover:text-pink-500">隐私政策</a>
            </label>
          </div>

          <!-- 错误信息 -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-3">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>

          <!-- 注册按钮 -->
          <button
            type="submit"
            :disabled="isLoading || !isFormValid"
            class="w-full bg-gradient-to-r from-pink-500 to-orange-500 text-white py-3 px-4 rounded-lg font-medium hover:from-pink-600 hover:to-orange-600 focus:outline-none focus:ring-2 focus:ring-pink-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <span v-if="isLoading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              注册中...
            </span>
            <span v-else>注册</span>
          </button>
        </form>

        <!-- 登录链接 -->
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600">
            已有账户？
            <router-link
              to="/login"
              class="text-pink-600 hover:text-pink-500 font-medium"
            >
              立即登录
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { EyeIcon, EyeOffIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { UserApi } from '@/api/user'
import type { UserRegisterRequest } from '@/api/types'

// 路由
const router = useRouter()

// 认证store
const authStore = useAuthStore()

// 响应式数据
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const isLoading = ref(false)
const isGeneratingNickname = ref(false)
const isCheckingNickname = ref(false)
const nicknameValidationMessage = ref('')
const error = ref<string | null>(null)

// 表单数据
const form = reactive<UserRegisterRequest & { confirmPassword: string; agree: boolean }>({
  username: '',
  email: '',
  password: '',
  nickname: '',
  confirmPassword: '',
  agree: false
})

// 密码强度检查
const passwordStrength = computed(() => {
  const password = form.password
  return {
    hasMinLength: password.length >= 8,
    hasLowercase: /[a-z]/.test(password),
    hasUppercase: /[A-Z]/.test(password),
  }
})

const isPasswordValid = computed(() => {
  const strength = passwordStrength.value
  return strength.hasMinLength && strength.hasLowercase && strength.hasUppercase
})

// 表单验证
const isFormValid = computed(() => {
  return (
    form.username.trim() &&
    form.email.trim() &&
    form.nickname.trim() &&
    isPasswordValid.value &&
    form.password === form.confirmPassword &&
    form.agree
  )
})

// 处理注册
const handleRegister = async () => {
  try {
    isLoading.value = true
    error.value = null

    // 验证密码强度
    if (form.password.length < 8) {
      error.value = '密码长度至少需要8位'
      return
    }
    
    if (!/[a-z]/.test(form.password)) {
      error.value = '密码必须包含至少一个小写字母'
      return
    }
    
    if (!/[A-Z]/.test(form.password)) {
      error.value = '密码必须包含至少一个大写字母'
      return
    }
    
    // 验证密码一致性
    if (form.password !== form.confirmPassword) {
      error.value = '两次输入的密码不一致'
      return
    }

    // 调用注册API
    await authStore.register({
      username: form.username,
      email: form.email,
      password: form.password,
      nickname: form.nickname
    })

    // 注册成功，跳转到登录页面
    router.push({
      path: '/login',
      query: {
        message: '注册成功！请使用您的账号密码登录'
      }
    })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '注册失败，请重试'
  } finally {
    isLoading.value = false
  }
}

// 随机昵称生成词库 - 确保生成的昵称在2-12个字符范围内
const adjectives = [
  '可爱', '温柔', '聪明', '活泼', '甜美', '阳光', '快乐', '善良',
  '机智', '勇敢', '优雅', '淘气', '开朗', '温暖', '纯真', '灵动',
  '清新', '梦幻', '闪亮', '柔软', '迷人', '俏皮', '安静', '热情'
]

const nouns = [
  '小猫', '小兔', '小熊', '小鸟', '小鱼', '小狐', '小鹿', '小羊',
  '星星', '月亮', '彩虹', '花朵', '蝴蝶', '雪花', '云朵', '露珠',
  '糖果', '棉花', '泡泡', '樱花', '薄荷', '柠檬', '草莓', '蜜桃',
  '天使', '精灵', '公主', '王子', '宝石', '珍珠', '水晶', '钻石'
]

// 检查昵称是否已存在
const checkNicknameExists = async (nickname: string): Promise<boolean> => {
  try {
    const response = await UserApi.checkNickname(nickname)
    return response.data.exists
  } catch (err) {
    console.error('检查昵称失败:', err)
    // 如果API调用失败，返回false（假设昵称可用）
    return false
  }
}

// 防抖计时器
let nicknameValidationTimer: number | null = null

// 验证昵称可用性
const validateNickname = async (nickname: string) => {
  if (!nickname || nickname.length < 2) {
    nicknameValidationMessage.value = ''
    return
  }
  
  if (nickname.length > 12) {
    nicknameValidationMessage.value = '昵称长度不能超过12个字符'
    return
  }
  
  isCheckingNickname.value = true
  nicknameValidationMessage.value = ''
  
  try {
    const exists = await checkNicknameExists(nickname)
    if (exists) {
      nicknameValidationMessage.value = '该昵称已被使用，请选择其他昵称'
    } else {
      nicknameValidationMessage.value = '昵称可用'
    }
  } catch (err) {
    nicknameValidationMessage.value = '检查昵称时出错，请稍后重试'
  } finally {
    isCheckingNickname.value = false
  }
}

// 防抖的昵称验证
const debouncedValidateNickname = (nickname: string) => {
  if (nicknameValidationTimer) {
    clearTimeout(nicknameValidationTimer)
  }
  
  nicknameValidationTimer = setTimeout(() => {
    validateNickname(nickname)
  }, 500) // 500ms防抖
}

// 生成符合长度限制的随机昵称（2-12个字符）
const generateRandomNickname = async () => {
  if (isGeneratingNickname.value) return
  
  try {
    isGeneratingNickname.value = true
    
    let attempts = 0
    const maxAttempts = 10
    
    while (attempts < maxAttempts) {
      // 随机选择形容词和名词
      const randomAdjective = adjectives[Math.floor(Math.random() * adjectives.length)]
      const randomNoun = nouns[Math.floor(Math.random() * nouns.length)]
      let generatedNickname = randomAdjective + randomNoun
      
      // 确保昵称长度在12个字符以内
      if (generatedNickname.length > 12) {
        // 如果太长，尝试只使用名词
        generatedNickname = randomNoun
        // 如果名词还是太长，截取前12个字符
        if (generatedNickname.length > 12) {
          generatedNickname = generatedNickname.substring(0, 12)
        }
      }
      
      // 确保昵称长度至少2个字符
      if (generatedNickname.length < 2) {
        generatedNickname = randomNoun
      }
      
      // 检查昵称是否已存在
      const exists = await checkNicknameExists(generatedNickname)
      
      if (!exists) {
        form.nickname = generatedNickname
        // 生成昵称后立即验证
        await validateNickname(generatedNickname)
        break
      }
      
      attempts++
    }
    
    // 如果尝试多次仍然重复，添加1-2位随机数字
    if (attempts >= maxAttempts) {
      const randomAdjective = adjectives[Math.floor(Math.random() * adjectives.length)]
      const randomNoun = nouns[Math.floor(Math.random() * nouns.length)]
      let baseNickname = randomAdjective + randomNoun
      
      // 确保基础昵称 + 数字不超过12个字符
      const maxBaseLength = 10 // 留2个字符给数字
      if (baseNickname.length > maxBaseLength) {
        baseNickname = baseNickname.substring(0, maxBaseLength)
      }
      
      const randomNumber = Math.floor(Math.random() * 99) + 1
      const finalNickname = baseNickname + randomNumber
      form.nickname = finalNickname
      // 验证生成的昵称
      await validateNickname(finalNickname)
    }
    
  } catch (err) {
    console.error('生成随机昵称失败:', err)
    // 如果生成失败，使用备用方案
    const randomNoun = nouns[Math.floor(Math.random() * nouns.length)]
    const randomNumber = Math.floor(Math.random() * 9) + 1
    let backupNickname = randomNoun + randomNumber
    
    // 确保备用昵称也符合长度限制
    if (backupNickname.length > 12) {
      backupNickname = randomNoun.substring(0, 11) + randomNumber
    }
    
    form.nickname = backupNickname
    // 验证生成的昵称
    await validateNickname(backupNickname)
  } finally {
    isGeneratingNickname.value = false
  }
}
</script>