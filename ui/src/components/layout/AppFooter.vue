<template>
  <footer :class="rootClass">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8">
        <!-- 品牌信息 -->
        <div class="sm:col-span-2 md:col-span-1">
          <div class="flex items-center mb-3 sm:mb-4">
            <!-- Logo -->
            <img
              src="@/assets/images/logo/GoDad_logo.png"
              alt="GoDad Logo"
              class="w-12 h-12 sm:w-14 sm:h-14 mr-2 sm:mr-3 rounded-lg shadow-lg"
            />
            <span class="text-lg sm:text-xl font-bold">GoDad</span>
          </div>
          <p :class="[mutedTextClass, 'mb-3 sm:mb-6 text-sm sm:text-base leading-relaxed']">
            专业的育儿知识分享平台，陪伴每一个家庭的成长之路。
          </p>
          <div class="flex space-x-2 sm:space-x-3">
            <div :class="['w-8 h-8 sm:w-10 sm:h-10 rounded-full flex items-center justify-center transition-all duration-300 cursor-pointer group shadow-lg hover:bg-gradient-to-r hover:from-pink-600 hover:to-orange-600', socialIconClass]">
              <span class="text-xs sm:text-sm font-medium group-hover:text-white">微</span>
            </div>
            <div :class="['w-8 h-8 sm:w-10 sm:h-10 rounded-full flex items-center justify-center transition-all duration-300 cursor-pointer group shadow-lg hover:bg-gradient-to-r hover:from-pink-600 hover:to-orange-600', socialIconClass]">
              <span class="text-xs sm:text-sm font-medium group-hover:text-white">博</span>
            </div>
            <div :class="['w-8 h-8 sm:w-10 sm:h-10 rounded-full flex items-center justify-center transition-all duration-300 cursor-pointer group shadow-lg hover:bg-gradient-to-r hover:from-pink-600 hover:to-orange-600', socialIconClass]">
              <span class="text-xs sm:text-sm font-medium group-hover:text-white">抖</span>
            </div>
          </div>
        </div>

        <!-- 快速链接 -->
        <div>
          <h3 :class="['text-base sm:text-lg font-semibold mb-3 sm:mb-4', headingTextClass]">快速链接</h3>
          <ul class="space-y-2 sm:space-y-3">
            <li v-for="link in quickLinks" :key="link.path">
              <router-link
                :to="link.path"
                :class="linkClass"
              >
                {{ link.label }}
              </router-link>
            </li>
          </ul>
        </div>

        <!-- 热门分类 -->
        <div>
          <h3 :class="['text-base sm:text-lg font-semibold mb-3 sm:mb-4', headingTextClass]">热门分类</h3>
          <ul class="space-y-2 sm:space-y-3">
            <li v-for="category in hotCategories" :key="category.path">
              <router-link
                :to="category.path"
                :class="linkClass"
              >
                {{ category.label }}
              </router-link>
            </li>
          </ul>
        </div>

        <!-- 联系我们 -->
        <div>
          <h3 :class="['text-base sm:text-lg font-semibold mb-3 sm:mb-4', headingTextClass]">联系我们</h3>
          <div class="space-y-3 sm:space-y-4">
            <div class="flex items-center group">
              <MailIcon :class="['h-4 w-4 sm:h-5 sm:w-5 mr-2 sm:mr-3 transition-colors duration-300', iconMutedClass, 'group-hover:text-pink-400']" />
              <span :class="[mutedTextClass, 'text-sm sm:text-base transition-colors duration-300 group-hover:text-current']">support@godad.com</span>
            </div>
            <div class="flex items-center group">
              <PhoneIcon :class="['h-4 w-4 sm:h-5 sm:w-5 mr-2 sm:mr-3 transition-colors duration-300', iconMutedClass, 'group-hover:text-pink-400']" />
              <span :class="[mutedTextClass, 'text-sm sm:text-base transition-colors duration-300 group-hover:text-current']">400-123-4567</span>
            </div>
            <div class="flex items-start group">
              <MapPinIcon :class="['h-4 w-4 sm:h-5 sm:w-5 mr-2 sm:mr-3 mt-0.5 transition-colors duration-300', iconMutedClass, 'group-hover:text-pink-400']" />
              <span :class="[mutedTextClass, 'text-sm sm:text-base transition-colors duration-300 group-hover:text-current']">北京市朝阳区科技园区</span>
            </div>
          </div>
        </div>
      </div>

      <div :class="['pt-6 sm:pt-8 border-t', borderLineClass]">
        <div class="flex flex-col sm:flex-row justify-between items-center space-y-3 sm:space-y-0">
          <p :class="[mutedTextClass, 'text-xs sm:text-sm text-center sm:text-left']">
            © {{ currentYear }} GoDad. 保留所有权利。
          </p>
          <div class="flex flex-col sm:flex-row space-y-1 sm:space-y-0 sm:space-x-4 text-xs sm:text-sm text-center sm:text-right">
            <a href="#" :class="linkSecondaryClass">备案号：京ICP备12345678号</a>
            <a href="#" :class="linkSecondaryClass">京公网安备 11010502012345号</a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { MailIcon, PhoneIcon, MapPinIcon } from 'lucide-vue-next'

// Props 支持自定义配置
interface FooterLink {
  path: string
  label: string
}

const props = withDefaults(defineProps<{
  quickLinks?: FooterLink[]
  hotCategories?: FooterLink[]
  showSocialMedia?: boolean
  showContact?: boolean
  theme?: 'dark' | 'tint' | 'brand' | 'light'
}>(), {
  quickLinks: () => [
    { path: '/about', label: '关于我们' },
    { path: '/contact', label: '联系我们' },
    { path: '/privacy', label: '隐私政策' },
    { path: '/terms', label: '服务条款' },
    { path: '/help', label: '帮助中心' }
  ],
  hotCategories: () => [
    { path: '/categories?category=newborn', label: '新生儿护理' },
    { path: '/categories?category=nutrition', label: '营养健康' },
    { path: '/categories?category=education', label: '早期教育' },
    { path: '/categories?category=psychology', label: '心理发展' },
    { path: '/categories?category=safety', label: '安全防护' }
  ],
  showSocialMedia: true,
  showContact: true,
  theme: 'dark'
})

// 当前年份
const currentYear = computed(() => new Date().getFullYear())

// 主题样式映射
const rootClass = computed(() => {
  switch (props.theme) {
    case 'tint':
      return 'bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50 text-slate-800'
    case 'brand':
      return 'bg-gradient-to-br from-pink-600 via-orange-600 to-amber-500 text-white'
    case 'light':
      return 'bg-white text-slate-800'
    case 'dark':
    default:
      return 'bg-slate-900 text-white'
  }
})

const headingTextClass = computed(() => (props.theme === 'dark' || props.theme === 'brand') ? 'text-white' : 'text-slate-900')

const mutedTextClass = computed(() => (props.theme === 'dark' || props.theme === 'brand') ? 'text-white/70' : 'text-slate-600')

const iconMutedClass = mutedTextClass

const linkBase = 'transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block'
const linkClass = computed(() => {
  switch (props.theme) {
    case 'tint':
    case 'light':
      return `${linkBase} text-pink-600 hover:text-pink-700`
    case 'brand':
      return `${linkBase} text-white/80 hover:text-white`
    case 'dark':
    default:
      return `${linkBase} text-white/70 hover:text-pink-300`
  }
})

const linkSecondaryClass = computed(() => {
  switch (props.theme) {
    case 'tint':
    case 'light':
      return 'text-slate-600 hover:text-slate-900 transition-colors duration-300'
    case 'brand':
      return 'text-white/80 hover:text-white transition-colors duration-300'
    case 'dark':
    default:
      return 'text-white/70 hover:text-white transition-colors duration-300'
  }
})

const socialIconClass = computed(() => {
  switch (props.theme) {
    case 'tint':
    case 'light':
      return 'bg-white border border-gray-200 text-slate-700'
    case 'brand':
      return 'bg-white/10 text-white'
    case 'dark':
    default:
      return 'bg-gray-800'
  }
})

const borderLineClass = computed(() => {
  switch (props.theme) {
    case 'tint':
    case 'light':
      return 'border-gray-200'
    case 'brand':
      return 'border-white/20'
    case 'dark':
    default:
      return 'border-white/10'
  }
})
</script>
