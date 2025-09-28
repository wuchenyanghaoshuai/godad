<template>
  <router-link :to="to" :class="linkClass">
    <span>{{ label }}</span>
    <span
      v-if="badge && badge > 0"
      class="ml-1 inline-flex items-center justify-center min-w-5 h-5 px-1 text-xs font-semibold bg-red-500 text-white rounded-full align-middle"
    >
      {{ displayBadge }}
    </span>
  </router-link>
  
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

interface Props {
  to: string
  label: string
  badge?: number
  exact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  badge: 0,
  exact: false,
})

const route = useRoute()

const isActive = computed(() => {
  const current = route.path || ''
  if (props.exact) return current === props.to
  // 高亮当前路径或其子路径
  return current === props.to || current.startsWith(props.to + '/')
})

const displayBadge = computed(() => Math.min(props.badge || 0, 99))

const base = 'relative px-3 py-2 text-sm font-medium rounded-md'
const activeClass = 'text-blue-600 bg-blue-50'
const inactiveClass = 'text-gray-600 hover:text-blue-600 hover:bg-gray-50'

const linkClass = computed(() => [base, isActive.value ? activeClass : inactiveClass])
</script>

<style scoped>
</style>

