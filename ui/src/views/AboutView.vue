<template>
  <div class="min-h-screen bg-gradient-warm font-chinese">
    <!-- 导航栏 -->
    <Navbar />
    <div class="py-8">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- 页面标题 -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-soft-800 mb-4">关于我们</h1>
        <p class="text-xl text-soft-600 max-w-2xl mx-auto">
          专业的育儿指导平台，致力于为每个家庭提供科学、温暖的育儿支持
        </p>
      </div>

      <!-- 主要内容 -->
      <div class="bg-card-gradient rounded-3xl shadow-soft backdrop-blur-sm border border-white/20 p-8 mb-8">
        <div class="prose prose-lg max-w-none">
          <h2 class="text-2xl font-semibold text-soft-800 mb-6">我们的使命</h2>
          <p class="text-soft-700 mb-6 leading-relaxed">
            GoDad 致力于成为每个家庭最信赖的育儿伙伴。我们相信，每个孩子都是独特的，每个家庭都有自己的育儿方式。
            我们的目标是提供科学、实用的育儿知识，帮助父母们在育儿路上更加自信和从容。
          </p>

          <h2 class="text-2xl font-semibold text-soft-800 mb-6">我们的价值观</h2>
          <div class="grid md:grid-cols-2 gap-6 mb-8">
            <div class="bg-primary-50 p-6 rounded-2xl border border-primary-100/50 shadow-gentle hover:shadow-warm transition-all duration-300">
              <h3 class="text-lg font-semibold text-primary-800 mb-3">专业可靠</h3>
              <p class="text-primary-700">
                所有内容都经过专业团队审核，确保信息的准确性和科学性
              </p>
            </div>
            <div class="bg-secondary-50 p-6 rounded-2xl border border-secondary-100/50 shadow-gentle hover:shadow-warm transition-all duration-300">
              <h3 class="text-lg font-semibold text-secondary-800 mb-3">温暖陪伴</h3>
              <p class="text-secondary-700">
                我们不仅提供知识，更提供情感支持，陪伴每个家庭度过育儿的每个阶段
              </p>
            </div>
            <div class="bg-warm-50 p-6 rounded-2xl border border-warm-100/50 shadow-gentle hover:shadow-warm transition-all duration-300">
              <h3 class="text-lg font-semibold text-warm-800 mb-3">个性化服务</h3>
              <p class="text-warm-700">
                根据不同年龄段和个体差异，提供个性化的育儿建议和解决方案
              </p>
            </div>
            <div class="bg-gradient-to-br from-primary-50 to-secondary-50 p-6 rounded-2xl border border-primary-100/50 shadow-gentle hover:shadow-warm transition-all duration-300">
              <h3 class="text-lg font-semibold text-primary-800 mb-3">持续创新</h3>
              <p class="text-primary-700">
                不断学习最新的育儿理念和方法，为用户提供与时俱进的服务
              </p>
            </div>
          </div>

          <h2 class="text-2xl font-semibold text-soft-800 mb-6">联系我们</h2>
          <div class="bg-gradient-to-r from-warm-50 to-primary-50 p-6 rounded-2xl border border-warm-100/50 shadow-gentle">
            <div class="grid md:grid-cols-2 gap-6">
              <div>
                <h4 class="font-semibold text-soft-800 mb-2">客服邮箱</h4>
                <p class="text-soft-700">support@godad.com</p>
              </div>
              <div>
                <h4 class="font-semibold text-soft-800 mb-2">合作咨询</h4>
                <p class="text-soft-700">business@godad.com</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 返回首页按钮 -->
      <div class="text-center">
        <router-link
          to="/"
          class="inline-flex items-center bg-primary-500 hover:bg-primary-600 text-white px-8 py-3 rounded-2xl font-medium transition-all duration-300 shadow-warm hover:shadow-lg transform hover:-translate-y-0.5"
        >
          返回首页
        </router-link>
      </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Navbar from '@/components/Navbar.vue'
</script>