<template>
  <div v-if="showSearch" class="hidden lg:flex flex-1 max-w-lg mx-8">
    <div class="relative w-full">
      <SearchIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
      <input
        v-model="searchQuery"
        @keyup.enter="performSearch"
        @focus="showSearchSuggestions = true"
        @blur="hideSearchSuggestions"
        type="text"
        placeholder="搜索文章、用户..."
        class="w-full pl-10 pr-4 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent bg-gray-50 focus:bg-white transition-all duration-200"
      />

      <!-- 搜索建议下拉框 -->
      <div
        v-if="showSearchSuggestions && searchSuggestions.length > 0"
        class="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-60 overflow-y-auto"
      >
        <div
          v-for="suggestion in searchSuggestions"
          :key="suggestion"
          @mousedown="searchWithSuggestion(suggestion)"
          class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-b-0"
        >
          <SearchIcon class="inline h-3 w-3 mr-2 text-gray-400" />
          {{ suggestion }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { SearchIcon } from 'lucide-vue-next'

interface Props {
  showSearch?: boolean
}

withDefaults(defineProps<Props>(), {
  showSearch: true
})

const router = useRouter()

// 搜索相关
const searchQuery = ref('')
const showSearchSuggestions = ref(false)
const searchSuggestions = ref([
  '育儿知识', '健康饮食', '早教方法', '亲子关系', '学习指导'
])

// 搜索相关方法
const performSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      path: '/search',
      query: { q: searchQuery.value.trim() }
    })
    searchQuery.value = ''
    showSearchSuggestions.value = false
  }
}

const searchWithSuggestion = (suggestion: string) => {
  searchQuery.value = suggestion
  performSearch()
}

const hideSearchSuggestions = () => {
  setTimeout(() => {
    showSearchSuggestions.value = false
  }, 200)
}
</script>