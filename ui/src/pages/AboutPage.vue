<script setup lang="ts">
import { ref } from 'vue'
import { AppLayout, PageContainer } from '@/components/layout'
import heroImg from '@/assets/about/hero-dad-baby.svg'
import avatarLi from '@/assets/about/team-li.svg'
import avatarWang from '@/assets/about/team-wang.svg'
import avatarZhang from '@/assets/about/team-zhang.svg'
import avatarPlaceholder from '@/assets/about/avatar-placeholder.svg'
import { HeartIcon, UsersIcon, TrendingUpIcon, StarIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { computed } from 'vue'

const teamMembers = ref([
  {
    name: '李医生',
    role: '儿科专家',
    description: '10年儿科临床经验，专注新生儿护理指导',
    avatar: avatarLi
  },
  {
    name: '王营养师',
    role: '儿童营养师',
    description: '国家注册营养师，婴幼儿营养搭配专家',
    avatar: avatarWang
  },
  {
    name: '张老师',
    role: '早教专家',
    description: '幼儿教育硕士，0-6岁早期教育指导师',
    avatar: avatarZhang
  }
])

const milestones = ref([
  { year: '2023', event: 'GoDad平台正式上线', description: '为男性育儿提供专业指导' },
  { year: '2023', event: '用户突破1万', description: '越来越多的爸爸们加入我们' },
  { year: '2024', event: '专家团队扩充', description: '汇聚更多专业育儿专家' },
  { year: '2024', event: '功能持续优化', description: '不断完善用户体验' }
])

// 成员头像兜底占位
const onMemberAvatarError = (e: Event, _name: string) => {
  const img = e.target as HTMLImageElement
  if (img && img.src !== avatarPlaceholder) {
    img.src = avatarPlaceholder
  }
}

// CTA 根据登录状态决定链接与文案
const authStore = useAuthStore()
const ctaTo = computed(() => authStore.isAuthenticated ? '/user-center' : '/register')
const ctaLabel = computed(() => authStore.isAuthenticated ? '前往个人中心' : '立即注册')

 
</script>

<template>
  <AppLayout footer-theme="tint">
    <!-- Hero -->
    <PageContainer background="tint" padding="xl">
      <div class="text-center max-w-3xl mx-auto">
        <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">关于 GoDad</h1>
        <p class="text-gray-600">专为男性育儿知识分享打造的现代化社区平台</p>
      </div>
    </PageContainer>

    <!-- Main Content -->
    <PageContainer background="white" padding="xl" max-width="7xl">
      <!-- 我们的故事 -->
      <section class="mb-24">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900">我们的故事</h2>
        </div>

        <div class="rounded-3xl bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50 p-6 md:p-8 border border-gray-100 mb-12 md:mb-16 lg:mb-20">
        <div class="grid md:grid-cols-2 gap-8 items-stretch">
          <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm h-full min-h-[280px] flex flex-col justify-center">
            <h3 class="text-xl font-semibold text-gray-900 mb-4">理念与初心</h3>
            <p class="text-gray-600 leading-relaxed mb-4">
              GoDad（出发爸）诞生于一个简单而重要的理念：每一位父亲都应该积极参与孩子的成长过程。在传统的育儿观念中，
              母亲往往承担着更多的育儿责任，而父亲的角色常常被边缘化。
            </p>
            <p class="text-gray-600 leading-relaxed mb-4">
              我们相信，父亲在孩子的成长中扮演着不可替代的角色。通过专业的知识分享、经验交流和互助支持，
              我们希望帮助每一位父亲成为更好的育儿伙伴，让每一个家庭都能更好地"一起出发，共同成长"。
            </p>
            <p class="text-gray-600 leading-relaxed">
              目前平台还处于开发阶段，我们正在不断完善功能，汇聚更多专业资源，未来我们将考虑商业化运营，
              为更多家庭提供优质的育儿支持服务。
            </p>
          </div>
          <div class="rounded-2xl border border-gray-100 shadow-sm overflow-hidden h-full min-h-[280px]">
            <img
              :src="heroImg"
              alt="父子互动"
              class="w-full h-full object-cover"
            >
          </div>
        </div>
        </div>
      </section>

      <!-- 我们的使命 -->
      <section class="mt-6 mb-16">
        <div class="rounded-2xl border border-gray-100 bg-white p-8 shadow-sm">
          <div class="grid md:grid-cols-3 gap-6">
            <div class="text-center">
              <div class="w-14 h-14 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <HeartIcon class="h-8 w-8 text-white" />
              </div>
              <h3 class="text-lg font-semibold text-gray-900 mb-2">专业可靠</h3>
              <p class="text-gray-600 text-sm leading-relaxed">汇聚儿科医生、营养师、教育专家的专业知识，确保每一条建议都科学可靠。</p>
            </div>
            <div class="text-center">
              <div class="w-14 h-14 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <UsersIcon class="h-8 w-8 text-white" />
              </div>
              <h3 class="text-lg font-semibold text-gray-900 mb-2">温暖社区</h3>
              <p class="text-gray-600 text-sm leading-relaxed">连接全球父母，分享育儿经验，在这里你永远不会感到孤单。</p>
            </div>
            <div class="text-center">
              <div class="w-14 h-14 bg-gradient-to-r from-green-500 to-teal-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <TrendingUpIcon class="h-8 w-8 text-white" />
              </div>
              <h3 class="text-lg font-semibold text-gray-900 mb-2">持续成长</h3>
              <p class="text-gray-600 text-sm leading-relaxed">紧跟最新育儿研究成果，定期更新内容，让您的育儿知识始终保持前沿。</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 专家团队 -->
      <section class="mb-16">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-3">专家团队</h2>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto">我们的专家团队由经验丰富的儿科医生、营养师和教育专家组成</p>
        </div>
        <div class="grid md:grid-cols-3 gap-8">
          <div
            v-for="member in teamMembers"
            :key="member.name"
            class="bg-white rounded-2xl border border-gray-100 shadow-sm p-8 text-center hover:shadow-md transition-shadow duration-300"
          >
            <img :src="member.avatar" :alt="member.name" class="w-24 h-24 rounded-full mx-auto mb-6 object-cover" @error="(e:any)=>onMemberAvatarError(e, member.name)" />
            <h3 class="text-xl font-bold text-gray-900 mb-2">{{ member.name }}</h3>
            <p class="text-pink-600 font-medium mb-4">{{ member.role }}</p>
            <p class="text-gray-600 leading-relaxed">{{ member.description }}</p>
          </div>
        </div>
      </section>

      <!-- 发展历程（纵向时间线卡片） -->
      <section class="mb-16">
        <div class="text-center mb-8">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900">发展历程</h2>
        </div>
        <div class="relative max-w-3xl mx-auto">
          <!-- vertical line -->
          <div class="absolute left-5 top-0 bottom-0 w-px bg-gradient-to-b from-pink-200 to-orange-200"></div>
          <div
            v-for="(milestone, index) in milestones"
            :key="index"
            class="relative pl-12 pb-8"
          >
            <!-- node dot -->
            <div class="absolute left-4.5 top-4 w-2.5 h-2.5 rounded-full bg-gradient-to-r from-pink-500 to-orange-500 border-2 border-white shadow"></div>
            <!-- year badge -->
            <span class="inline-block text-xs font-semibold px-2 py-0.5 rounded-full bg-pink-50 text-pink-700 border border-pink-100">{{ milestone.year }}</span>
            <!-- card -->
            <div class="mt-2 bg-white border border-gray-100 rounded-xl p-4 shadow-sm">
              <div class="text-gray-900 font-medium mb-1">{{ milestone.event }}</div>
              <div class="text-gray-600 text-sm leading-relaxed">{{ milestone.description }}</div>
            </div>
          </div>
        </div>
      </section>

      <!-- CTA -->
      <section class="text-center rounded-3xl p-10 bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50">
        <h2 class="text-2xl md:text-3xl font-bold mb-4 text-gray-900">加入我们的社区</h2>
        <p class="text-base text-gray-600 mb-6 max-w-2xl mx-auto">成为 GoDad 社区的一员，与万千父母一起分享育儿经验，获得专业指导</p>
        <div class="space-x-3">
          <router-link
            :to="ctaTo"
            class="inline-flex items-center bg-white text-pink-600 px-6 py-3 rounded-full font-semibold border border-gray-200 hover:bg-pink-50 transition-colors duration-300"
          >
            {{ ctaLabel }}
          </router-link>
          <router-link
            to="/articles"
            class="inline-flex items-center border-2 border-pink-600 text-pink-600 px-6 py-3 rounded-full font-semibold hover:bg-pink-600 hover:text-white transition-colors duration-300"
          >
            浏览文章
          </router-link>
        </div>
      </section>
    </PageContainer>
  </AppLayout>
</template>
