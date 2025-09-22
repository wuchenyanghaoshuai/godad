<script setup lang="ts">
import { ref } from 'vue'
import { AppLayout, PageContainer } from '@/components/layout'
import { MailIcon, PhoneIcon, MapPinIcon, ClockIcon, MessageCircleIcon, HeadphonesIcon } from 'lucide-vue-next'

// 响应式数据
const contactForm = ref({
  name: '',
  email: '',
  subject: '',
  message: ''
})

const contactMethods = ref([
  {
    icon: MailIcon,
    title: '邮箱联系',
    description: '我们会在24小时内回复您的邮件',
    contact: 'support@godad.com',
    link: 'mailto:support@godad.com'
  },
  {
    icon: PhoneIcon,
    title: '电话咨询',
    description: '工作日 9:00-18:00 提供电话支持',
    contact: '400-123-4567',
    link: 'tel:400-123-4567'
  },
  {
    icon: MessageCircleIcon,
    title: '在线客服',
    description: '7x24小时在线客服支持',
    contact: '点击开始对话',
    link: '#'
  }
])

const offices = ref([
  {
    city: '北京总部',
    address: '北京市朝阳区科技园区创新大厦15楼',
    phone: '010-12345678',
    email: 'beijing@godad.com'
  },
  {
    city: '上海分部',
    address: '上海市浦东新区陆家嘴金融中心20楼',
    phone: '021-87654321',
    email: 'shanghai@godad.com'
  },
  {
    city: '深圳分部',
    address: '深圳市南山区科技园软件大厦12楼',
    phone: '0755-98765432',
    email: 'shenzhen@godad.com'
  }
])

const faq = ref([
  {
    question: '如何注册GoDad账户？',
    answer: '点击页面右上角"注册"按钮，填写基本信息即可完成注册。注册后可以发布文章、参与讨论、关注其他用户。'
  },
  {
    question: '忘记密码怎么办？',
    answer: '在登录页面点击"忘记密码"，输入注册邮箱，我们会发送重置密码链接到您的邮箱。'
  },
  {
    question: '如何发布文章？',
    answer: '登录后点击"发布文章"按钮，选择分类，编写内容，支持Markdown格式和图片上传。'
  },
  {
    question: '平台收费吗？',
    answer: '目前GoDad平台完全免费使用，包括文章发布、社区交流等所有功能。'
  }
])

// 提交表单
const submitForm = () => {
  // 这里可以添加实际的表单提交逻辑
  alert('感谢您的反馈！我们会尽快与您联系。')
  // 重置表单
  contactForm.value = {
    name: '',
    email: '',
    subject: '',
    message: ''
  }
}
</script>

<template>
  <AppLayout footer-theme="tint">
    <!-- Hero -->
    <PageContainer background="tint" padding="xl">
      <div class="text-center max-w-3xl mx-auto">
        <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">联系我们</h1>
        <p class="text-gray-600">我们随时为您提供帮助和支持</p>
      </div>
    </PageContainer>

    <!-- Main Content -->
    <PageContainer background="white" padding="xl" max-width="7xl">

      <!-- 联系方式 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">多种联系方式</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto">
            选择最适合您的联系方式，我们的团队随时准备为您提供帮助
          </p>
        </div>

        <div class="grid md:grid-cols-3 gap-8">
          <div
            v-for="method in contactMethods"
            :key="method.title"
            class="bg-white rounded-2xl shadow-lg p-8 text-center hover:shadow-xl transition-shadow duration-300"
          >
            <div class="w-16 h-16 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mx-auto mb-6">
              <component :is="method.icon" class="h-8 w-8 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-4">{{ method.title }}</h3>
            <p class="text-gray-600 mb-4">{{ method.description }}</p>
            <a
              :href="method.link"
              class="text-pink-600 font-semibold hover:text-orange-600 transition-colors duration-300"
            >
              {{ method.contact }}
            </a>
          </div>
        </div>
      </section>

      <!-- 联系表单和办公地址 -->
      <div class="grid lg:grid-cols-2 gap-12 mb-20">
        <!-- 联系表单 -->
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">发送消息</h3>
          <form @submit.prevent="submitForm">
            <div class="grid md:grid-cols-2 gap-6 mb-6">
              <div>
                <label for="name" class="block text-sm font-medium text-gray-700 mb-2">姓名 *</label>
                <input
                  v-model="contactForm.name"
                  type="text"
                  id="name"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors duration-300"
                  placeholder="请输入您的姓名"
                >
              </div>
              <div>
                <label for="email" class="block text-sm font-medium text-gray-700 mb-2">邮箱 *</label>
                <input
                  v-model="contactForm.email"
                  type="email"
                  id="email"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors duration-300"
                  placeholder="请输入您的邮箱"
                >
              </div>
            </div>
            <div class="mb-6">
              <label for="subject" class="block text-sm font-medium text-gray-700 mb-2">主题 *</label>
              <input
                v-model="contactForm.subject"
                type="text"
                id="subject"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors duration-300"
                placeholder="请输入消息主题"
              >
            </div>
            <div class="mb-6">
              <label for="message" class="block text-sm font-medium text-gray-700 mb-2">消息内容 *</label>
              <textarea
                v-model="contactForm.message"
                id="message"
                rows="6"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors duration-300"
                placeholder="请详细描述您的问题或建议"
              ></textarea>
            </div>
            <button
              type="submit"
              class="w-full bg-gradient-to-r from-pink-600 to-orange-600 text-white py-3 px-6 rounded-lg font-semibold hover:from-pink-700 hover:to-orange-700 transition-all duration-300 transform hover:scale-105"
            >
              发送消息
            </button>
          </form>
        </div>

        <!-- 办公地址 -->
        <div>
          <h3 class="text-2xl font-bold text-gray-900 mb-6">办公地址</h3>
          <div class="space-y-6">
            <div
              v-for="office in offices"
              :key="office.city"
              class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300"
            >
              <h4 class="text-lg font-bold text-gray-900 mb-4">{{ office.city }}</h4>
              <div class="space-y-3">
                <div class="flex items-start">
                  <MapPinIcon class="h-5 w-5 text-pink-600 mr-3 mt-0.5" />
                  <span class="text-gray-600">{{ office.address }}</span>
                </div>
                <div class="flex items-center">
                  <PhoneIcon class="h-5 w-5 text-pink-600 mr-3" />
                  <span class="text-gray-600">{{ office.phone }}</span>
                </div>
                <div class="flex items-center">
                  <MailIcon class="h-5 w-5 text-pink-600 mr-3" />
                  <span class="text-gray-600">{{ office.email }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 工作时间 -->
          <div class="bg-gradient-to-r from-pink-600 to-orange-600 rounded-2xl p-6 text-white mt-6">
            <div class="flex items-center mb-4">
              <ClockIcon class="h-6 w-6 mr-3" />
              <h4 class="text-lg font-bold">工作时间</h4>
            </div>
            <div class="space-y-2">
              <div class="flex justify-between">
                <span>周一至周五</span>
                <span>9:00 - 18:00</span>
              </div>
              <div class="flex justify-between">
                <span>周六</span>
                <span>10:00 - 16:00</span>
              </div>
              <div class="flex justify-between">
                <span>周日</span>
                <span>休息</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 常见问题 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">常见问题</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto">
            快速找到您问题的答案
          </p>
        </div>

        <div class="max-w-4xl mx-auto">
          <div class="space-y-6">
            <div
              v-for="item in faq"
              :key="item.question"
              class="bg-white rounded-2xl shadow-lg overflow-hidden"
            >
              <div class="p-6">
                <h3 class="text-lg font-bold text-gray-900 mb-3">{{ item.question }}</h3>
                <p class="text-gray-600 leading-relaxed">{{ item.answer }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Call to Action -->
      <section class="text-center rounded-3xl p-10 bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50">
        <HeadphonesIcon class="h-12 w-12 mx-auto mb-4 text-pink-600" />
        <h2 class="text-2xl md:text-3xl font-bold mb-4 text-gray-900">还有其他问题？</h2>
        <p class="text-base text-gray-600 mb-6 max-w-2xl mx-auto">
          我们的客服团队随时准备为您提供帮助，无论是技术问题还是使用建议
        </p>
        <div class="space-x-3">
          <a
            href="mailto:support@godad.com"
            class="inline-flex items-center bg-white text-pink-600 px-6 py-3 rounded-full font-semibold border border-gray-200 hover:bg-pink-50 transition-colors duration-300"
          >
            立即邮件咨询
          </a>
          <a
            href="tel:400-123-4567"
            class="inline-flex items-center border-2 border-pink-600 text-pink-600 px-6 py-3 rounded-full font-semibold hover:bg-pink-600 hover:text-white transition-colors duration-300"
          >
            电话咨询
          </a>
        </div>
      </section>
    </PageContainer>
  </AppLayout>
</template>
