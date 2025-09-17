<script setup lang="ts">
import { ref } from 'vue'
import { ShieldCheckIcon, LockIcon, EyeIcon, DatabaseIcon } from 'lucide-vue-next'

// 隐私条款内容
const sections = ref([
  {
    title: '信息收集',
    icon: DatabaseIcon,
    content: [
      '我们收集您主动提供的信息，包括注册时的用户名、邮箱地址、头像等基本信息。',
      '我们会收集您在平台上的活动信息，如发布的文章、评论、点赞等行为数据。',
      '为了改善服务质量，我们可能收集设备信息、IP地址、浏览器类型等技术信息。',
      '我们不会收集敏感的个人信息，如身份证号、银行卡号等。'
    ]
  },
  {
    title: '信息使用',
    icon: EyeIcon,
    content: [
      '提供和维护GoDad平台服务，包括用户认证、内容展示、社区交流等功能。',
      '改善用户体验，包括个性化推荐、搜索优化、界面改进等。',
      '与您沟通，包括发送重要通知、回复咨询、处理反馈等。',
      '保护平台安全，包括防止欺诈、滥用、恶意行为等。',
      '遵守法律法规要求，在必要时配合相关部门调查。'
    ]
  },
  {
    title: '信息共享',
    icon: ShieldCheckIcon,
    content: [
      '我们不会向第三方出售、出租或以其他方式披露您的个人信息。',
      '在您明确同意的情况下，我们可能与合作伙伴分享必要信息。',
      '为了提供服务，我们可能与技术服务提供商共享部分信息，但会确保其遵守严格的保密协议。',
      '在法律要求或为保护我们及用户权益的情况下，我们可能披露必要信息。'
    ]
  },
  {
    title: '信息安全',
    icon: LockIcon,
    content: [
      '我们采用行业标准的安全措施保护您的个人信息，包括加密传输、访问控制等。',
      '定期进行安全审计和漏洞检测，及时修复潜在的安全风险。',
      '严格限制员工对个人信息的访问权限，仅在工作需要时才能访问。',
      '一旦发现安全事件，我们会立即采取措施并及时通知受影响的用户。'
    ]
  }
])

const userRights = ref([
  {
    title: '访问权',
    description: '您有权了解我们收集了您的哪些个人信息'
  },
  {
    title: '更正权',
    description: '您有权要求我们更正不准确或不完整的个人信息'
  },
  {
    title: '删除权',
    description: '在特定情况下，您有权要求我们删除您的个人信息'
  },
  {
    title: '限制处理权',
    description: '在特定情况下，您有权要求我们限制对您个人信息的处理'
  },
  {
    title: '数据可携权',
    description: '您有权以结构化、常用和机器可读的格式获取您的个人信息'
  },
  {
    title: '撤回同意权',
    description: '如果我们基于您的同意处理个人信息，您有权随时撤回同意'
  }
])

// 最后更新时间
const lastUpdated = ref('2024年1月15日')
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 to-orange-50">
    <!-- Header -->
    <div class="bg-gradient-to-r from-pink-600 to-orange-600 text-white py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <ShieldCheckIcon class="h-20 w-20 mx-auto mb-6" />
        <h1 class="text-4xl md:text-6xl font-bold mb-6">隐私政策</h1>
        <p class="text-xl md:text-2xl text-pink-100 max-w-3xl mx-auto">
          保护您的隐私是我们的首要责任
        </p>
        <p class="text-lg text-pink-200 mt-4">
          最后更新：{{ lastUpdated }}
        </p>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">

      <!-- 引言 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-6">引言</h2>
          <div class="prose prose-lg text-gray-600">
            <p class="mb-4">
              欢迎使用GoDad（出发爸）平台。我们深知个人信息对您的重要性，并会严格保护您的隐私安全。
              本隐私政策将向您详细说明我们如何收集、使用、储存和保护您的个人信息。
            </p>
            <p class="mb-4">
              使用我们的服务，即表示您同意本隐私政策的条款。如果您不同意本政策的任何内容，
              请不要使用我们的服务。我们可能会不时更新本政策，建议您定期查看。
            </p>
            <p>
              如果您对本隐私政策有任何疑问，请通过
              <router-link to="/contact" class="text-pink-600 hover:text-orange-600 font-medium">
                联系我们页面
              </router-link>
              与我们取得联系。
            </p>
          </div>
        </div>
      </section>

      <!-- 详细条款 -->
      <section class="mb-16">
        <div class="space-y-8">
          <div
            v-for="section in sections"
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

      <!-- Cookie政策 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">Cookie和类似技术</h3>
          <div class="prose text-gray-600">
            <p class="mb-4">
              我们使用Cookie和类似技术来改善您的用户体验。Cookie是存储在您设备上的小文件，
              帮助我们记住您的偏好设置、保持登录状态、分析网站流量等。
            </p>
            <p class="mb-4">
              我们使用的Cookie类型包括：
            </p>
            <ul class="list-disc list-inside space-y-2 mb-4">
              <li><strong>必要Cookie：</strong>确保网站正常运行的基本功能</li>
              <li><strong>功能Cookie：</strong>记住您的偏好设置和选择</li>
              <li><strong>分析Cookie：</strong>帮助我们了解网站使用情况，改善服务</li>
              <li><strong>广告Cookie：</strong>用于展示相关的个性化内容（如适用）</li>
            </ul>
            <p>
              您可以通过浏览器设置管理Cookie偏好，但请注意，禁用某些Cookie可能影响网站功能。
            </p>
          </div>
        </div>
      </section>

      <!-- 您的权利 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">您的权利</h3>
          <p class="text-gray-600 mb-6">
            根据相关法律法规，您对自己的个人信息享有以下权利：
          </p>
          <div class="grid md:grid-cols-2 gap-6">
            <div
              v-for="right in userRights"
              :key="right.title"
              class="border border-gray-200 rounded-lg p-4 hover:border-pink-300 transition-colors duration-300"
            >
              <h4 class="font-semibold text-gray-900 mb-2">{{ right.title }}</h4>
              <p class="text-sm text-gray-600">{{ right.description }}</p>
            </div>
          </div>
          <div class="mt-6 p-4 bg-pink-50 rounded-lg">
            <p class="text-sm text-gray-600">
              <strong>如何行使权利：</strong>
              如需行使上述权利，请通过邮件
              <a href="mailto:privacy@godad.com" class="text-pink-600 hover:text-orange-600">
                privacy@godad.com
              </a>
              联系我们，我们会在30天内回复您的请求。
            </p>
          </div>
        </div>
      </section>

      <!-- 儿童隐私 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">儿童隐私保护</h3>
          <div class="prose text-gray-600">
            <p class="mb-4">
              我们非常重视儿童隐私保护。GoDad平台主要面向成年用户，我们不会有意收集13岁以下儿童的个人信息。
            </p>
            <p class="mb-4">
              如果您是家长或监护人，发现您的孩子在未经您同意的情况下向我们提供了个人信息，
              请立即通过
              <a href="mailto:privacy@godad.com" class="text-pink-600 hover:text-orange-600">
                privacy@godad.com
              </a>
              联系我们，我们会立即删除相关信息。
            </p>
            <p>
              对于13-18岁的未成年用户，我们建议在家长或监护人的指导下使用我们的服务。
            </p>
          </div>
        </div>
      </section>

      <!-- 国际传输 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">跨境数据传输</h3>
          <div class="prose text-gray-600">
            <p class="mb-4">
              您的个人信息主要存储在中国境内的服务器上。为了提供更好的服务，
              我们可能需要将您的部分信息传输至其他国家或地区。
            </p>
            <p class="mb-4">
              在进行跨境传输时，我们会：
            </p>
            <ul class="list-disc list-inside space-y-2 mb-4">
              <li>确保接收方具备适当的数据保护水平</li>
              <li>签署标准合同条款或其他法律保障措施</li>
              <li>采取技术和管理措施保护数据安全</li>
              <li>在法律要求的情况下获得您的明确同意</li>
            </ul>
          </div>
        </div>
      </section>

      <!-- 政策更新 -->
      <section class="mb-16">
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h3 class="text-2xl font-bold text-gray-900 mb-6">政策更新</h3>
          <div class="prose text-gray-600">
            <p class="mb-4">
              我们可能会不时更新本隐私政策，以反映我们服务的变化或法律要求的变更。
              重要更新将通过以下方式通知您：
            </p>
            <ul class="list-disc list-inside space-y-2 mb-4">
              <li>在平台上发布显著通知</li>
              <li>通过邮件向注册用户发送通知</li>
              <li>在应用内推送重要更新通知</li>
            </ul>
            <p>
              建议您定期查看本政策。继续使用我们的服务即表示您接受更新后的政策条款。
            </p>
          </div>
        </div>
      </section>

      <!-- 联系信息 -->
      <section class="text-center bg-gradient-to-r from-pink-600 to-orange-600 rounded-3xl p-12 text-white">
        <LockIcon class="h-16 w-16 mx-auto mb-6" />
        <h2 class="text-3xl md:text-4xl font-bold mb-6">隐私问题咨询</h2>
        <p class="text-xl text-pink-100 mb-8 max-w-2xl mx-auto">
          如果您对我们的隐私政策有任何疑问或建议，请随时与我们联系
        </p>
        <div class="space-y-4">
          <div class="text-lg">
            <strong>隐私专员邮箱：</strong>
            <a href="mailto:privacy@godad.com" class="underline hover:no-underline">
              privacy@godad.com
            </a>
          </div>
          <div class="text-lg">
            <strong>客服电话：</strong>
            <a href="tel:400-123-4567" class="underline hover:no-underline">
              400-123-4567
            </a>
          </div>
          <div class="text-lg">
            <strong>办公地址：</strong>
            北京市朝阳区科技园区创新大厦15楼
          </div>
        </div>
      </section>
    </div>
  </div>
</template>