<script setup lang="ts">
import { ref } from 'vue'
import { AppLayout, PageContainer } from '@/components/layout'
import { ShieldCheckIcon, UserCheckIcon, AlertTriangleIcon, ScaleIcon } from 'lucide-vue-next'

// 服务条款内容
const termssections = ref([
  {
    title: '接受条款',
    icon: UserCheckIcon,
    content: [
      '通过访问和使用GoDad（出发爸）平台，您同意遵守并受本服务条款的约束。',
      '如果您不同意这些条款，请立即停止使用我们的服务。',
      '我们保留随时修改这些条款的权利，修改后的条款将在发布后立即生效。',
      '继续使用我们的服务将被视为您接受修改后的条款。'
    ]
  },
  {
    title: '用户责任',
    icon: ShieldCheckIcon,
    content: [
      '您有责任确保您提供的所有信息都是准确、完整和最新的。',
      '您同意不使用虚假身份或冒充他人使用我们的服务。',
      '您有责任保护您的账户安全，包括密码的保密。',
      '您同意立即通知我们任何未经授权使用您账户的情况。'
    ]
  },
  {
    title: '禁止行为',
    icon: AlertTriangleIcon,
    content: [
      '发布违法、诽谤、骚扰、威胁或侵犯他人权利的内容。',
      '传播病毒、恶意软件或其他有害代码。',
      '试图未经授权访问我们的系统或其他用户的账户。',
      '干扰或破坏服务的正常运行。',
      '将平台用于商业目的而未经我们明确同意。',
      '发布虚假、误导性或不准确的信息。'
    ]
  },
  {
    title: '知识产权',
    icon: ScaleIcon,
    content: [
      '平台上的所有内容，包括但不限于文本、图像、标志、软件，均受版权法保护。',
      '您保留对自己发布内容的所有权，但授予我们使用、展示和分发这些内容的许可。',
      '未经许可，您不得复制、修改、分发或以其他方式使用我们的知识产权。',
      '如果您认为您的知识产权受到侵犯，请通过我们的投诉流程联系我们。'
    ]
  }
])

const liabilityTerms = ref([
  {
    title: '服务可用性',
    content: '我们努力保持服务的持续可用性，但不保证服务永远不会中断。'
  },
  {
    title: '内容准确性',
    content: '我们不对用户生成的内容的准确性、完整性或可靠性承担责任。'
  },
  {
    title: '第三方链接',
    content: '我们的服务可能包含指向第三方网站的链接，我们不对这些网站的内容负责。'
  },
  {
    title: '损害赔偿',
    content: '在法律允许的最大范围内，我们不对任何间接、偶然或后果性损害承担责任。'
  }
])

// 最后更新时间
const lastUpdated = ref('2024年1月15日')
</script>

<template>
  <AppLayout footer-theme="tint">
    <PageContainer background="tint" padding="xl" max-width="4xl">
      <div class="text-center">
        <ScaleIcon class="h-12 w-12 sm:h-16 sm:w-16 mx-auto mb-4 text-pink-600" />
        <h1 class="text-3xl sm:text-4xl font-bold text-gray-900 mb-2">服务条款</h1>
        <p class="text-gray-600">使用GoDad平台前请仔细阅读以下条款 · 最后更新：{{ lastUpdated }}</p>
      </div>
    </PageContainer>

    <PageContainer background="white" padding="xl" max-width="4xl">

      <!-- 引言 -->
      <section class="mb-12">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-6">重要说明</h2>
          <div class="prose prose-lg text-gray-600">
            <p class="mb-4">
              本服务条款（以下简称"条款"）是您与GoDad（出发爸）平台之间的法律协议。
              使用我们的服务表示您已阅读、理解并同意受这些条款的约束。
            </p>
            <p class="mb-4">
              我们建议您定期查看这些条款，因为我们可能会根据法律要求或业务需要进行更新。
              重要变更将通过平台通知或邮件形式告知您。
            </p>
            <p>
              如果您对这些条款有任何疑问，请通过
              <router-link to="/contact" class="text-pink-600 hover:text-orange-600 font-medium">
                联系我们页面
              </router-link>
              与我们取得联系。
            </p>
          </div>
        </div>
      </section>

      <!-- 详细条款 -->
      <section class="mb-12">
        <div class="space-y-8">
          <div
            v-for="section in termssections"
            :key="section.title"
            class="bg-white rounded-2xl shadow-lg p-8"
          >
            <div class="flex items-center mb-6">
              <div class="w-12 h-12 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mr-4">
                <component :is="section.icon" class="h-6 w-6 text-white" />
              </div>
              <h3 class="text-2xl font-bold text-gray-900">{{ section.title }}</h3>
            </div>
            <ul class="space-y-3">
              <li
                v-for="item in section.content"
                :key="item"
                class="flex items-start"
              >
                <div class="w-2 h-2 bg-pink-500 rounded-full mt-2 mr-3 flex-shrink-0"></div>
                <span class="text-gray-600 leading-relaxed">{{ item }}</span>
              </li>
            </ul>
          </div>
        </div>
      </section>

      <!-- 免责声明 -->
      <section class="mb-12">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">免责声明与责任限制</h3>
          <div class="grid md:grid-cols-2 gap-6">
            <div
              v-for="term in liabilityTerms"
              :key="term.title"
              class="border border-gray-200 rounded-lg p-4 hover:border-pink-300 transition-colors duration-300"
            >
              <h4 class="font-semibold text-gray-900 mb-2">{{ term.title }}</h4>
              <p class="text-sm text-gray-600">{{ term.content }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 账户终止 -->
      <section class="mb-12">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">账户终止</h3>
          <div class="prose text-gray-600">
            <p class="mb-4">
              您可以随时通过联系我们来终止您的账户。账户终止后，您将无法访问与该账户相关的内容和数据。
            </p>
            <p class="mb-4">
              我们保留在以下情况下暂停或终止您的账户的权利：
            </p>
            <ul class="list-disc list-inside space-y-2 mb-4">
              <li>违反本服务条款</li>
              <li>从事非法或有害活动</li>
              <li>损害平台或其他用户的利益</li>
              <li>长期不活跃（超过2年未登录）</li>
            </ul>
            <p>
              账户终止不会影响在终止之前产生的任何义务或责任。
            </p>
          </div>
        </div>
      </section>

      <!-- 争议解决 -->
      <section class="mb-12">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">争议解决</h3>
          <div class="prose text-gray-600">
            <p class="mb-4">
              如果您与我们之间产生争议，我们建议首先通过友好协商的方式解决。
              您可以通过以下方式联系我们：
            </p>
            <ul class="list-disc list-inside space-y-2 mb-4">
              <li>邮件：legal@godad.com</li>
              <li>电话：400-123-4567</li>
              <li>在线客服：通过平台联系我们</li>
            </ul>
            <p class="mb-4">
              如果无法通过协商解决争议，则应提交至北京市朝阳区人民法院解决。
            </p>
            <p>
              本条款受中华人民共和国法律管辖。
            </p>
          </div>
        </div>
      </section>

      <!-- 其他条款 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">其他重要条款</h3>
          <div class="space-y-6">
            <div>
              <h4 class="font-semibold text-gray-900 mb-2">完整协议</h4>
              <p class="text-gray-600">
                本服务条款构成您与我们之间关于使用服务的完整协议，
                取代所有先前的口头或书面协议。
              </p>
            </div>
            <div>
              <h4 class="font-semibold text-gray-900 mb-2">可分割性</h4>
              <p class="text-gray-600">
                如果本条款的任何部分被认定为无效或不可执行，
                其余部分将继续有效。
              </p>
            </div>
            <div>
              <h4 class="font-semibold text-gray-900 mb-2">语言版本</h4>
              <p class="text-gray-600">
                本条款的中文版本为正式版本。如有其他语言版本，
                以中文版本为准。
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- 联系信息 -->
      <section class="text-center rounded-3xl p-10 bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50">
        <ScaleIcon class="h-12 w-12 mx-auto mb-4 text-pink-600" />
        <h2 class="text-2xl md:text-3xl font-bold mb-4 text-gray-900">法律问题咨询</h2>
        <p class="text-base text-gray-600 mb-6 max-w-2xl mx-auto">
          如果您对我们的服务条款有任何法律相关问题，请随时与我们联系
        </p>
        <div class="space-y-3 text-gray-800">
          <div>
            <strong>法务部邮箱：</strong>
            <a href="mailto:legal@godad.com" class="text-pink-600 hover:text-pink-700">legal@godad.com</a>
          </div>
          <div>
            <strong>客服电话：</strong>
            <a href="tel:400-123-4567" class="text-pink-600 hover:text-pink-700">400-123-4567</a>
          </div>
          <div>
            <strong>办公地址：</strong>
            北京市朝阳区科技园区创新大厦15楼
          </div>
        </div>
      </section>
    </PageContainer>
  </AppLayout>
</template>

 
