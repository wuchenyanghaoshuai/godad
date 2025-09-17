<script setup lang="ts">
import { ref } from 'vue'
import { HeartIcon, UsersIcon, TrendingUpIcon, StarIcon } from 'lucide-vue-next'

// 响应式数据
const teamMembers = ref([
  {
    name: '李医生',
    role: '儿科专家',
    description: '10年儿科临床经验，专注新生儿护理指导',
    avatar: 'https://images.unsplash.com/photo-1559839734-2b71ea197ec2?w=150&h=150&fit=crop&crop=face'
  },
  {
    name: '王营养师',
    role: '儿童营养师',
    description: '国家注册营养师，婴幼儿营养搭配专家',
    avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b417?w=150&h=150&fit=crop&crop=face'
  },
  {
    name: '张老师',
    role: '早教专家',
    description: '幼儿教育硕士，0-6岁早期教育指导师',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face'
  }
])

const milestones = ref([
  { year: '2023', event: 'GoDad平台正式上线', description: '为男性育儿提供专业指导' },
  { year: '2023', event: '用户突破1万', description: '越来越多的爸爸们加入我们' },
  { year: '2024', event: '专家团队扩充', description: '汇聚更多专业育儿专家' },
  { year: '2024', event: '功能持续优化', description: '不断完善用户体验' }
])
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 to-orange-50">
    <!-- Header -->
    <div class="bg-gradient-to-r from-pink-600 to-orange-600 text-white py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h1 class="text-4xl md:text-6xl font-bold mb-6">关于 GoDad</h1>
        <p class="text-xl md:text-2xl text-pink-100 max-w-3xl mx-auto">
          专为男性育儿知识分享打造的现代化社区平台
        </p>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">

      <!-- 我们的故事 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">我们的故事</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
        </div>

        <div class="grid md:grid-cols-2 gap-12 items-center">
          <div>
            <h3 class="text-2xl font-bold text-gray-900 mb-6">让我们一起出发，共同成长</h3>
            <p class="text-gray-600 text-lg leading-relaxed mb-6">
              GoDad（出发爸）诞生于一个简单而重要的理念：每一位父亲都应该积极参与孩子的成长过程。在传统的育儿观念中，
              母亲往往承担着更多的育儿责任，而父亲的角色常常被边缘化。
            </p>
            <p class="text-gray-600 text-lg leading-relaxed mb-6">
              我们相信，父亲在孩子的成长中扮演着不可替代的角色。通过专业的知识分享、经验交流和互助支持，
              我们希望帮助每一位父亲成为更好的育儿伙伴，让每一个家庭都能更好地"一起出发，共同成长"。
            </p>
            <p class="text-gray-600 text-lg leading-relaxed">
              目前平台还处于开发阶段，我们正在不断完善功能，汇聚更多专业资源，
              未来我们将考虑商业化运营，为更多家庭提供优质的育儿支持服务。
            </p>
          </div>
          <div class="relative">
            <img
              src="https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=600&h=400&fit=crop"
              alt="父子互动"
              class="rounded-2xl shadow-2xl"
            >
            <div class="absolute -bottom-6 -right-6 w-24 h-24 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center shadow-lg">
              <HeartIcon class="h-12 w-12 text-white" />
            </div>
          </div>
        </div>
      </section>

      <!-- 我们的使命 -->
      <section class="mb-20">
        <div class="bg-white rounded-3xl shadow-xl p-8 md:p-12">
          <div class="grid md:grid-cols-3 gap-8">
            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mx-auto mb-6">
                <HeartIcon class="h-8 w-8 text-white" />
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-4">专业可靠</h3>
              <p class="text-gray-600 leading-relaxed">
                汇聚儿科医生、营养师、教育专家的专业知识，确保每一条建议都科学可靠。
              </p>
            </div>

            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-6">
                <UsersIcon class="h-8 w-8 text-white" />
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-4">温暖社区</h3>
              <p class="text-gray-600 leading-relaxed">
                连接全球父母，分享育儿经验，在这里你永远不会感到孤单。
              </p>
            </div>

            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-r from-green-500 to-teal-500 rounded-full flex items-center justify-center mx-auto mb-6">
                <TrendingUpIcon class="h-8 w-8 text-white" />
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-4">持续成长</h3>
              <p class="text-gray-600 leading-relaxed">
                紧跟最新育儿研究成果，定期更新内容，让您的育儿知识始终保持前沿。
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- 专家团队 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">专家团队</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto">
            我们的专家团队由经验丰富的儿科医生、营养师和教育专家组成
          </p>
        </div>

        <div class="grid md:grid-cols-3 gap-8">
          <div
            v-for="member in teamMembers"
            :key="member.name"
            class="bg-white rounded-2xl shadow-lg p-8 text-center hover:shadow-xl transition-shadow duration-300"
          >
            <img
              :src="member.avatar"
              :alt="member.name"
              class="w-24 h-24 rounded-full mx-auto mb-6 object-cover"
            >
            <h3 class="text-xl font-bold text-gray-900 mb-2">{{ member.name }}</h3>
            <p class="text-pink-600 font-medium mb-4">{{ member.role }}</p>
            <p class="text-gray-600 leading-relaxed">{{ member.description }}</p>
          </div>
        </div>
      </section>

      <!-- 发展历程 -->
      <section class="mb-20">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">发展历程</h2>
          <div class="w-24 h-1 bg-gradient-to-r from-pink-600 to-orange-600 mx-auto mb-8"></div>
        </div>

        <div class="relative">
          <div class="absolute left-1/2 transform -translate-x-px h-full w-0.5 bg-gradient-to-b from-pink-600 to-orange-600"></div>

          <div class="space-y-12">
            <div
              v-for="(milestone, index) in milestones"
              :key="index"
              class="relative flex items-center"
              :class="index % 2 === 0 ? 'justify-start' : 'justify-end'"
            >
              <div
                class="bg-white rounded-2xl shadow-lg p-6 max-w-md"
                :class="index % 2 === 0 ? 'mr-8' : 'ml-8'"
              >
                <div class="flex items-center mb-4">
                  <div class="w-12 h-12 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mr-4">
                    <StarIcon class="h-6 w-6 text-white" />
                  </div>
                  <span class="text-2xl font-bold text-pink-600">{{ milestone.year }}</span>
                </div>
                <h3 class="text-lg font-bold text-gray-900 mb-2">{{ milestone.event }}</h3>
                <p class="text-gray-600">{{ milestone.description }}</p>
              </div>

              <div class="absolute left-1/2 transform -translate-x-1/2 w-4 h-4 bg-gradient-to-r from-pink-600 to-orange-600 rounded-full border-4 border-white shadow-lg"></div>
            </div>
          </div>
        </div>
      </section>

      <!-- Call to Action -->
      <section class="text-center bg-gradient-to-r from-pink-600 to-orange-600 rounded-3xl p-12 text-white">
        <h2 class="text-3xl md:text-4xl font-bold mb-6">加入我们的社区</h2>
        <p class="text-xl text-pink-100 mb-8 max-w-2xl mx-auto">
          成为GoDad社区的一员，与万千父母一起分享育儿经验，获得专业指导
        </p>
        <div class="space-x-4">
          <router-link
            to="/register"
            class="inline-flex items-center bg-white text-pink-600 px-8 py-4 rounded-full font-semibold hover:bg-pink-50 transition-colors duration-300"
          >
            立即注册
          </router-link>
          <router-link
            to="/articles"
            class="inline-flex items-center border-2 border-white text-white px-8 py-4 rounded-full font-semibold hover:bg-white hover:text-pink-600 transition-colors duration-300"
          >
            浏览文章
          </router-link>
        </div>
      </section>
    </div>
  </div>
</template>