<template>
  <AppLayout footer-theme="tint">
    <!-- Hero -->
    <PageContainer background="tint" padding="xl">
      <div class="text-center max-w-3xl mx-auto">
        <h1 class="text-3xl sm:text-4xl font-bold text-gray-900 mb-3">关于我们</h1>
        <p class="text-gray-600 text-base sm:text-lg">专业的育儿指导平台，致力于为每个家庭提供科学、温暖的育儿支持</p>
      </div>
    </PageContainer>

    <!-- Content -->
    <PageContainer background="white" padding="xl">
      <div class="grid md:grid-cols-2 gap-6 mb-8">
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
          <h2 class="text-xl font-semibold text-gray-900 mb-3">我们的使命</h2>
          <p class="text-gray-600 leading-relaxed">GoDad 致力于成为每个家庭最信赖的育儿伙伴。我们相信，每个孩子都是独特的，每个家庭都有自己的育儿方式。我们的目标是提供科学、实用的育儿知识，帮助父母更自信、更从容。</p>
        </div>
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
          <h2 class="text-xl font-semibold text-gray-900 mb-3">我们的愿景</h2>
          <p class="text-gray-600 leading-relaxed">打造一个温暖、可信赖的育儿社区，用专业链接每个家庭的成长，让好的内容与经验被更好地看见与传递。</p>
        </div>
      </div>

      <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="rounded-xl border border-gray-100 bg-white p-5 shadow-sm">
          <div class="text-pink-600 font-bold text-2xl mb-1">专业</div>
          <p class="text-gray-600 text-sm">内容经专业团队审核，保证科学、可靠</p>
        </div>
        <div class="rounded-xl border border-gray-100 bg-white p-5 shadow-sm">
          <div class="text-orange-600 font-bold text-2xl mb-1">温暖</div>
          <p class="text-gray-600 text-sm">以用户为中心，关注情感支持与陪伴</p>
        </div>
        <div class="rounded-xl border border-gray-100 bg-white p-5 shadow-sm">
          <div class="text-purple-600 font-bold text-2xl mb-1">共创</div>
          <p class="text-gray-600 text-sm">邀请专家与用户共同建设优质内容</p>
        </div>
        <div class="rounded-xl border border-gray-100 bg-white p-5 shadow-sm">
          <div class="text-green-600 font-bold text-2xl mb-1">成长</div>
          <p class="text-gray-600 text-sm">陪伴孩子与家长在实践中持续成长</p>
        </div>
      </div>
    </PageContainer>

    <!-- Stats strip -->
    <PageContainer background="tint" padding="lg" :show-max-width="true" max-width="7xl">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="rounded-xl bg-white shadow-sm border border-gray-100 p-5 text-center">
          <div class="text-2xl font-bold text-pink-600">10K+</div>
          <div class="text-sm text-gray-600">注册用户</div>
        </div>
        <div class="rounded-xl bg-white shadow-sm border border-gray-100 p-5 text-center">
          <div class="text-2xl font-bold text-orange-600">500+</div>
          <div class="text-sm text-gray-600">优质文章</div>
        </div>
        <div class="rounded-xl bg-white shadow-sm border border-gray-100 p-5 text-center">
          <div class="text-2xl font-bold text-purple-600">50+</div>
          <div class="text-sm text-gray-600">专业作者</div>
        </div>
        <div class="rounded-xl bg-white shadow-sm border border-gray-100 p-5 text-center">
          <div class="text-2xl font-bold text-green-600">24/7</div>
          <div class="text-sm text-gray-600">在线支持</div>
        </div>
      </div>
    </PageContainer>

    <!-- CTA -->
    <PageContainer background="white" padding="xl">
      <div class="text-center rounded-2xl border border-gray-100 p-8 bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50">
        <h3 class="text-2xl font-bold text-gray-900 mb-2">与我们一起建设更好的育儿社区</h3>
        <p class="text-gray-600 mb-4">欢迎专业作者、机构合作伙伴与我们联系，共同为家庭提供更优质的内容与服务。</p>
        <router-link to="/contact" class="inline-flex items-center bg-white text-pink-600 px-6 py-3 rounded-full font-semibold border border-gray-200 hover:bg-pink-50 transition-colors">联系团队</router-link>
      </div>
    </PageContainer>
  </AppLayout>
</template>

<script setup lang="ts">
import { AppLayout, PageContainer } from '@/components/layout'
</script>
