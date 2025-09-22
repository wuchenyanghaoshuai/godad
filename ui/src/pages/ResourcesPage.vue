<template>
  <AppLayout background="gray">
    <PageContainer background="transparent" padding="lg" max-width="7xl">
      <!-- Page Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-4">资源库</h1>
        <p class="text-gray-600 text-lg leading-relaxed max-w-4xl">
          探索我们精心策划的资源合集，旨在支持您的育儿之旅。从电子书和视频到互动工具，我们拥有您在抚养小宝贝的喜悦和挑战中所需的一切。
        </p>
      </div>

      <!-- Category Tabs -->
      <div class="mb-8">
        <div class="border-b border-gray-200">
          <nav class="flex space-x-8">
            <button
              v-for="category in categories"
              :key="category"
              @click="setSelectedCategory(category)"
              :class="[
                'py-4 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
                selectedCategory === category
                  ? 'border-pink-500 text-pink-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              {{ getCategoryLabel(category) }}
            </button>
          </nav>
        </div>
      </div>

      <!-- Resources Grid -->
      <div class="space-y-6">
        <div v-if="filteredResources.length === 0" class="text-center py-16">
          <div class="text-gray-400 mb-4">
            <BookOpenIcon class="h-16 w-16 mx-auto" />
          </div>
          <h3 class="text-lg font-medium text-gray-900 mb-2">该分类下暂无资源</h3>
          <p class="text-gray-500">请选择其他分类或稍后再试</p>
        </div>

        <ResourceCard
          v-for="resource in filteredResources"
          :key="resource.id"
          :resource="resource"
          @click="handleResourceClick"
        />
      </div>

      <!-- Load More -->
      <div v-if="filteredResources.length > 0 && hasMore" class="flex justify-center mt-12">
        <button
          @click="loadMore"
          :disabled="loading"
          :class="[
            'px-8 py-3 bg-white border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors font-medium shadow-sm',
            loading ? 'opacity-50 cursor-not-allowed' : ''
          ]"
        >
          {{ loading ? '加载中...' : '加载更多资源' }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading && filteredResources.length === 0" class="text-center py-16">
        <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-pink-500 mx-auto mb-4"></div>
        <p class="text-gray-500">正在加载资源...</p>
      </div>
    </PageContainer>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { BookOpenIcon } from 'lucide-vue-next'
import { AppLayout, PageContainer } from '@/components/layout'
import ResourceCard from '@/components/ResourceCard.vue'
import { ResourceApi, type Resource } from '@/api'

// 响应式数据
const selectedCategory = ref('All')
const resources = ref<Resource[]>([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const hasMore = ref(false)

// 分类数据
const categories = ['All', 'E-books', 'Videos', 'Tools']

// 分类标签映射
const categoryLabels: Record<string, string> = {
  'All': '全部',
  'E-books': '电子书',
  'Videos': '视频',
  'Tools': '工具'
}

// 加载资源数据
const loadResources = async (page = 1, append = false) => {
  try {
    loading.value = true
    const category = selectedCategory.value === 'All' ? undefined : selectedCategory.value

    const response = await ResourceApi.getPublishedResources({
      page,
      size: pageSize.value,
      category
    })

    if (response.code === 200) {
      const data = response.data
      if (append) {
        resources.value.push(...data.items)
      } else {
        resources.value = data.items
      }
      total.value = data.total
      currentPage.value = page
      hasMore.value = page < data.total_pages
    }
  } catch (error) {
    console.error('加载资源失败:', error)
  } finally {
    loading.value = false
  }
}

// 切换分类
const setSelectedCategory = async (category: string) => {
  selectedCategory.value = category
  currentPage.value = 1
  await loadResources(1, false)
}

// 加载更多
const loadMore = async () => {
  if (hasMore.value && !loading.value) {
    await loadResources(currentPage.value + 1, true)
  }
}

// 计算属性：过滤后的资源（现在直接使用API返回的数据）
const filteredResources = computed(() => resources.value)

// 方法
const getCategoryLabel = (category: string): string => {
  return categoryLabels[category] || category
}

const handleResourceClick = async (resource: Resource) => {
  try {
    // 调用下载接口（会增加下载次数）
    const response = await ResourceApi.downloadResource(resource.id)
    if (response.code === 200) {
      // 在新窗口打开资源文件
      window.open(response.data.file_url, '_blank')
    }
  } catch (error) {
    console.error('下载资源失败:', error)
    // 如果下载接口失败，直接使用原始URL
    if (resource.file_url) {
      window.open(resource.file_url, '_blank')
    }
  }
}

// 页面加载时获取资源数据
onMounted(() => {
  loadResources()
})
</script>