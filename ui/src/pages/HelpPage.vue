<script setup lang="ts">
import { ref, computed } from 'vue'
import { AppLayout, PageContainer } from '@/components/layout'
import {
  HelpCircleIcon,
  BookIcon,
  UserIcon,
  MessageSquareIcon,
  SettingsIcon,
  SearchIcon,
  ChevronDownIcon,
  ChevronUpIcon
} from 'lucide-vue-next'

// 搜索关键字
const searchQuery = ref('')

// 帮助分类
const helpCategories = ref([
  {
    id: 'getting-started',
    title: '新手入门',
    icon: BookIcon,
    description: '快速了解如何使用GoDad平台',
    articles: [
      { title: '如何注册账户', link: '#register' },
      { title: '完善个人资料', link: '#profile' },
      { title: '发布第一篇文章', link: '#first-post' },
      { title: '关注其他用户', link: '#follow' }
    ]
  },
  {
    id: 'account',
    title: '账户管理',
    icon: UserIcon,
    description: '管理您的账户设置和个人信息',
    articles: [
      { title: '修改密码', link: '#change-password' },
      { title: '更新个人信息', link: '#update-profile' },
      { title: '隐私设置', link: '#privacy' },
      { title: '注销账户', link: '#delete-account' }
    ]
  },
  {
    id: 'content',
    title: '内容发布',
    icon: MessageSquareIcon,
    description: '了解如何创建和管理内容',
    articles: [
      { title: '文章编写指南', link: '#writing-guide' },
      { title: '图片上传规范', link: '#image-upload' },
      { title: '使用Markdown格式', link: '#markdown' },
      { title: '内容审核规则', link: '#content-rules' }
    ]
  },
  {
    id: 'settings',
    title: '平台设置',
    icon: SettingsIcon,
    description: '自定义您的平台使用体验',
    articles: [
      { title: '通知设置', link: '#notifications' },
      { title: '主题切换', link: '#themes' },
      { title: '语言设置', link: '#language' },
      { title: '数据导出', link: '#export' }
    ]
  }
])

// 常见问题
const faqItems = ref([
  {
    id: 1,
    question: '如何重置密码？',
    answer: '在登录页面点击"忘记密码"链接，输入您的注册邮箱，我们会发送重置密码的链接到您的邮箱。请检查您的收件箱和垃圾邮件文件夹。',
    expanded: false
  },
  {
    id: 2,
    question: '文章发布后多久会显示？',
    answer: '文章发布后通常会立即显示在您的个人页面。如果文章包含需要审核的内容，可能需要1-24小时的审核时间。审核完成后会通过站内信通知您。',
    expanded: false
  },
  {
    id: 3,
    question: '如何删除已发布的文章？',
    answer: '进入您的文章详情页，点击右上角的"..."按钮，选择"删除文章"。请注意，删除操作无法撤销，建议您谨慎操作。',
    expanded: false
  },
  {
    id: 4,
    question: '为什么我的账户被限制了？',
    answer: '账户可能因为违反社区规则而被限制，包括但不限于发布不当内容、恶意举报他人、使用虚假信息等。如有疑问，请联系客服申请复审。',
    expanded: false
  },
  {
    id: 5,
    question: '如何举报不当内容？',
    answer: '在相关内容下方点击"举报"按钮，选择举报原因并提供详细说明。我们会在24小时内处理您的举报，并通过站内信通知处理结果。',
    expanded: false
  },
  {
    id: 6,
    question: '平台支持哪些图片格式？',
    answer: '平台支持JPG、PNG、GIF格式的图片，单张图片大小不超过10MB。建议使用JPG格式以获得最佳加载速度。',
    expanded: false
  }
])

// 快捷操作
const quickActions = ref([
  {
    title: '重置密码',
    description: '忘记密码？快速重置',
    icon: UserIcon,
    link: '#reset-password'
  },
  {
    title: '联系客服',
    description: '在线客服随时为您服务',
    icon: MessageSquareIcon,
    link: '/contact'
  },
  {
    title: '反馈建议',
    description: '帮助我们改进平台',
    icon: HelpCircleIcon,
    link: '#feedback'
  },
  {
    title: '使用教程',
    description: '观看视频教程',
    icon: BookIcon,
    link: '#tutorials'
  }
])

// 切换FAQ展开状态
const toggleFaq = (id: number) => {
  const faq = faqItems.value.find(item => item.id === id)
  if (faq) {
    faq.expanded = !faq.expanded
  }
}

// 搜索过滤
const filteredFaq = computed(() => {
  if (!searchQuery.value) return faqItems.value
  return faqItems.value.filter(item =>
    item.question.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    item.answer.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})
</script>

<template>
  <AppLayout footer-theme="tint">
    <PageContainer background="tint" padding="xl">
      <div class="text-center max-w-3xl mx-auto">
        <HelpCircleIcon class="h-12 w-12 sm:h-16 sm:w-16 mx-auto mb-4 text-pink-600" />
        <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">帮助中心</h1>
        <p class="text-gray-600">快速找到您需要的帮助和指导</p>
      </div>
    </PageContainer>

    <PageContainer background="white" padding="xl" max-width="7xl">

      <!-- 搜索栏 -->
      <section class="mb-16">
        <div class="max-w-2xl mx-auto">
          <div class="relative">
            <SearchIcon class="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索帮助内容..."
              class="w-full pl-12 pr-4 py-4 border border-gray-300 rounded-2xl focus:ring-2 focus:ring-pink-500 focus:border-transparent bg-white shadow-lg text-lg"
            >
          </div>
        </div>
      </section>

      <!-- 快捷操作 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">快捷操作</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto"></div>
        </div>

        <div class="grid md:grid-cols-4 gap-6">
          <div
            v-for="action in quickActions"
            :key="action.title"
            class="bg-white rounded-2xl shadow-lg p-6 text-center hover:shadow-xl transition-shadow duration-300 cursor-pointer"
          >
            <div class="w-12 h-12 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mx-auto mb-4">
              <component :is="action.icon" class="h-6 w-6 text-white" />
            </div>
            <h3 class="text-lg font-bold text-gray-900 mb-2">{{ action.title }}</h3>
            <p class="text-gray-600 text-sm">{{ action.description }}</p>
          </div>
        </div>
      </section>

      <!-- 帮助分类 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">帮助分类</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto">
            按照功能分类查找您需要的帮助信息
          </p>
        </div>

        <div class="grid md:grid-cols-2 gap-8">
          <div
            v-for="category in helpCategories"
            :key="category.id"
            class="bg-white rounded-2xl shadow-lg p-8 hover:shadow-xl transition-shadow duration-300"
          >
            <div class="flex items-center mb-6">
              <div class="w-12 h-12 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mr-4">
                <component :is="category.icon" class="h-6 w-6 text-white" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900">{{ category.title }}</h3>
                <p class="text-gray-600 text-sm">{{ category.description }}</p>
              </div>
            </div>
            <ul class="space-y-3">
              <li
                v-for="article in category.articles"
                :key="article.title"
                class="flex items-center"
              >
                <div class="w-1.5 h-1.5 bg-pink-500 rounded-full mr-3"></div>
                <a
                  :href="article.link"
                  class="text-gray-600 hover:text-pink-600 transition-colors duration-300"
                >
                  {{ article.title }}
                </a>
              </li>
            </ul>
          </div>
        </div>
      </section>

      <!-- 常见问题 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">常见问题</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto">
            查看用户最常询问的问题和解答
          </p>
        </div>

        <div class="max-w-4xl mx-auto">
          <div class="space-y-4">
            <div
              v-for="faq in filteredFaq"
              :key="faq.id"
              class="bg-white rounded-2xl shadow-lg overflow-hidden"
            >
              <button
                @click="toggleFaq(faq.id)"
                class="w-full p-6 text-left flex items-center justify-between hover:bg-gray-50 transition-colors duration-300"
              >
                <h3 class="text-lg font-semibold text-gray-900">{{ faq.question }}</h3>
                <component
                  :is="faq.expanded ? ChevronUpIcon : ChevronDownIcon"
                  class="h-5 w-5 text-gray-500"
                />
              </button>
              <div
                v-if="faq.expanded"
                class="px-6 pb-6 border-t border-gray-100"
              >
                <p class="text-gray-600 leading-relaxed pt-4">{{ faq.answer }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 联系支持 -->
      <section class="text-center rounded-3xl p-10 bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50">
        <MessageSquareIcon class="h-12 w-12 mx-auto mb-4 text-pink-600" />
        <h2 class="text-2xl md:text-3xl font-bold mb-4 text-gray-900">仍需要帮助？</h2>
        <p class="text-base text-gray-600 mb-6 max-w-2xl mx-auto">
          如果您在帮助中心没有找到答案，我们的客服团队随时为您提供支持
        </p>
        <div class="space-x-4">
          <router-link
            to="/contact"
            class="inline-flex items-center bg-white text-pink-600 px-6 py-3 rounded-full font-semibold border border-gray-200 hover:bg-pink-50 transition-colors duration-300"
          >
            联系客服
          </router-link>
          <a
            href="mailto:support@godad.com"
            class="inline-flex items-center border-2 border-pink-600 text-pink-600 px-6 py-3 rounded-full font-semibold hover:bg-pink-600 hover:text-white transition-colors duration-300"
          >
            邮件支持
          </a>
        </div>
      </section>
    </PageContainer>
  </AppLayout>
</template>
