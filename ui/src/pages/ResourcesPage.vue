<template>
  <div class="min-h-screen bg-gray-50">
    <BaseHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
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
      <div v-if="filteredResources.length > 0" class="flex justify-center mt-12">
        <button
          class="px-8 py-3 bg-white border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors font-medium shadow-sm"
        >
          加载更多资源
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { BookOpenIcon } from 'lucide-vue-next'
import BaseHeader from '@/components/BaseHeader.vue'
import ResourceCard from '@/components/ResourceCard.vue'

// 响应式数据
const selectedCategory = ref('All')

// 分类数据
const categories = ['All', 'E-books', 'Videos', 'Tools']

// 分类标签映射
const categoryLabels: Record<string, string> = {
  'All': '全部',
  'E-books': '电子书',
  'Videos': '视频',
  'Tools': '工具'
}

// 资源数据
const resources = [
  {
    id: 1,
    title: '新手爸爸第一年指南',
    description: '这本综合性电子书涵盖了从新生儿护理到发育里程碑的所有内容，为新手父亲提供实用建议和支持。',
    image: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?auto=format&fit=crop&q=80&w=600',
    type: 'e-book' as const,
    buttonText: '立即下载',
    category: 'E-books'
  },
  {
    id: 2,
    title: '宝宝第一步：视频系列',
    description: '跟随这个视频系列，探索宝宝第一步的激动人心的旅程，包含专家技巧和演示。',
    image: 'https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?auto=format&fit=crop&q=80&w=600',
    type: 'video' as const,
    buttonText: '立即观看',
    category: 'Videos'
  },
  {
    id: 3,
    title: '睡眠训练工具',
    description: '我们的互动工具帮助您为宝宝制定个性化的睡眠训练计划，确保全家人都能安眠。',
    image: 'https://images.unsplash.com/photo-1559182671-0e6c3cb2f6b1?auto=format&fit=crop&q=80&w=600',
    type: 'tool' as const,
    buttonText: '开始使用',
    category: 'Tools'
  },
  {
    id: 4,
    title: '婴儿喂养指南',
    description: '这本电子书提供了婴儿喂养的详细指南，涵盖母乳喂养、配方奶喂养和辅食添加。',
    image: 'https://images.unsplash.com/photo-1476703993599-0035a21b17a9?auto=format&fit=crop&q=80&w=600',
    type: 'e-book' as const,
    buttonText: '立即下载',
    category: 'E-books'
  },
  {
    id: 5,
    title: '发育里程碑跟踪器',
    description: '使用我们易于使用的工具跟踪宝宝的发育里程碑，确保他们走在正确的道路上。',
    image: 'https://images.unsplash.com/photo-1555252333-9f8e92e65df9?auto=format&fit=crop&q=80&w=600',
    type: 'tool' as const,
    buttonText: '跟踪里程碑',
    category: 'Tools'
  },
  {
    id: 6,
    title: '双胞胎育儿技巧',
    description: '这个视频系列为双胞胎父母提供实用建议和支持，涵盖从喂养到睡眠时间表的所有内容。',
    image: 'https://images.unsplash.com/photo-1588392382834-a891154bca4d?auto=format&fit=crop&q=80&w=600',
    type: 'video' as const,
    buttonText: '立即观看',
    category: 'Videos'
  },
  {
    id: 7,
    title: '0-3岁营养搭配手册',
    description: '专业营养师编写的权威指南，帮助您为不同年龄段的宝宝提供最佳营养搭配。',
    image: 'https://images.unsplash.com/photo-1609220136736-443140cffec6?auto=format&fit=crop&q=80&w=600',
    type: 'e-book' as const,
    buttonText: '立即下载',
    category: 'E-books'
  },
  {
    id: 8,
    title: '宝宝安全评估工具',
    description: '全面的家庭安全检查工具，帮助您识别和解决家中的潜在安全隐患。',
    image: 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?auto=format&fit=crop&q=80&w=600',
    type: 'tool' as const,
    buttonText: '开始评估',
    category: 'Tools'
  }
]

// 计算属性：过滤后的资源
const filteredResources = computed(() => {
  return resources.filter(resource => {
    return selectedCategory.value === 'All' || resource.category === selectedCategory.value
  })
})

// 方法
const setSelectedCategory = (category: string) => {
  selectedCategory.value = category
}

const getCategoryLabel = (category: string): string => {
  return categoryLabels[category] || category
}

const handleResourceClick = (resource: any) => {
  console.log('打开资源:', resource.title)
  // 这里可以实现具体的资源打开逻辑
}
</script>