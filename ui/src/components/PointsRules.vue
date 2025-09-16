<template>
  <div class="space-y-4">
    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-4 text-gray-500">加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="text-center py-8 text-red-500">
      <p>{{ error }}</p>
      <button
        @click="loadRules"
        class="mt-4 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
      >
        重试
      </button>
    </div>

    <!-- 积分规则列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="rule in rules"
        :key="rule.id"
        class="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-4 border border-gray-200 hover:shadow-md transition-shadow"
      >
        <div class="flex items-center justify-between mb-2">
          <h4 class="font-medium text-gray-900">{{ rule.name }}</h4>
          <span class="flex items-center text-lg font-bold" :class="rule.points > 0 ? 'text-green-600' : 'text-red-600'">
            {{ rule.points > 0 ? '+' : '' }}{{ rule.points }}
            <span class="text-xs ml-1 text-gray-500">积分</span>
          </span>
        </div>
        <p class="text-sm text-gray-600 mb-3">{{ rule.description }}</p>
        <div class="flex items-center justify-between text-xs text-gray-500">
          <span>操作类型: {{ rule.action }}</span>
          <span v-if="rule.daily_limit > 0">
            每日限制: {{ rule.daily_limit }}次
          </span>
          <span v-else class="text-green-600">
            无限制
          </span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && !error && rules.length === 0" class="text-center py-8 text-gray-500">
      <p>暂无积分规则</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { PointsAPI, type PointsRule } from '@/api'

// 响应式数据
const rules = ref<PointsRule[]>([])
const loading = ref(false)
const error = ref<string>('')

// 方法
const loadRules = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await PointsAPI.getPointsRules()
    rules.value = response.data
  } catch (err: any) {
    error.value = err.message || '加载积分规则失败'
    console.error('加载积分规则失败:', err)
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  loadRules()
})
</script>

<style scoped>
/* 可以添加一些自定义样式 */
</style>