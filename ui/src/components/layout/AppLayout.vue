<template>
  <div class="min-h-screen flex flex-col">
    <!-- 头部导航 -->
    <AppHeader
      v-if="showHeader"
      :nav-items="headerConfig.navItems"
      :show-search="headerConfig.showSearch"
      :show-create-button="headerConfig.showCreateButton"
      :show-notifications="headerConfig.showNotifications"
      :show-user-points="headerConfig.showUserPoints"
      :show-navigation="headerConfig.showNavigation"
      :show-user-menu="headerConfig.showUserMenu"
    />

    <!-- 主内容区域 -->
    <main class="flex-1" :class="mainClasses">
      <slot />
    </main>

    <!-- 页脚 -->
    <AppFooter
      v-if="showFooter"
      :quick-links="footerConfig.quickLinks"
      :hot-categories="footerConfig.hotCategories"
      :show-social-media="footerConfig.showSocialMedia"
      :show-contact="footerConfig.showContact"
      :theme="footerTheme"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppHeader from './AppHeader.vue'
import AppFooter from './AppFooter.vue'

// 导航项接口
interface NavItem {
  path: string
  label: string
}

// 页脚链接接口
interface FooterLink {
  path: string
  label: string
}

// 头部配置接口
interface HeaderConfig {
  navItems?: NavItem[]
  showSearch?: boolean
  showCreateButton?: boolean
  showNotifications?: boolean
  showUserPoints?: boolean
  showNavigation?: boolean
  showUserMenu?: boolean
}

// 页脚配置接口
interface FooterConfig {
  quickLinks?: FooterLink[]
  hotCategories?: FooterLink[]
  showSocialMedia?: boolean
  showContact?: boolean
}

// Props
const props = withDefaults(defineProps<{
  // 布局控制
  showHeader?: boolean
  showFooter?: boolean

  // 主内容区样式
  backgroundClass?: string
  containerClass?: string

  // 头部配置
  headerConfig?: HeaderConfig

  // 页脚配置
  footerConfig?: FooterConfig
  // 页脚主题
  footerTheme?: 'dark' | 'tint' | 'brand' | 'light'
}>(), {
  showHeader: true,
  showFooter: true,
  backgroundClass: 'bg-gray-50',
  containerClass: '',
  headerConfig: () => ({
    navItems: [
      { path: '/articles', label: '文章' },
      { path: '/community', label: '社区' },
      { path: '/resources', label: '资源' }
    ],
    showSearch: true,
    showCreateButton: true,
    showNotifications: true,
    showUserPoints: true,
    showNavigation: true,
    showUserMenu: true
  }),
  footerConfig: () => ({
    quickLinks: [
      { path: '/about', label: '关于我们' },
      { path: '/contact', label: '联系我们' },
      { path: '/privacy', label: '隐私政策' },
      { path: '/terms', label: '服务条款' },
      { path: '/help', label: '帮助中心' }
    ],
    hotCategories: [
      { path: '/categories?category=newborn', label: '新生儿护理' },
      { path: '/categories?category=nutrition', label: '营养健康' },
      { path: '/categories?category=education', label: '早期教育' },
      { path: '/categories?category=psychology', label: '心理发展' },
      { path: '/categories?category=safety', label: '安全防护' }
    ],
    showSocialMedia: true,
    showContact: true
  }),
  footerTheme: 'dark'
})

// 计算主内容区样式
const mainClasses = computed(() => {
  const classes = [props.backgroundClass]
  if (props.containerClass) {
    classes.push(props.containerClass)
  }
  return classes.join(' ')
})
</script>
