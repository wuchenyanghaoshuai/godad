<template>
  <div :class="containerClasses">
    <div v-if="showMaxWidth" :class="maxWidthClasses">
      <slot />
    </div>
    <slot v-else />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Props
const props = withDefaults(defineProps<{
  // 容器样式
  background?: 'gray' | 'white' | 'gradient' | 'transparent' | 'tint' | 'accent'
  padding?: 'none' | 'sm' | 'md' | 'lg' | 'xl'

  // 最大宽度容器
  showMaxWidth?: boolean
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | '4xl' | '5xl' | '6xl' | '7xl'

  // 自定义类名
  customClass?: string
}>(), {
  background: 'gray',
  padding: 'md',
  showMaxWidth: true,
  maxWidth: '7xl',
  customClass: ''
})

// 背景样式映射
const backgroundClasses = {
  gray: 'bg-gray-50',
  white: 'bg-white',
  gradient: 'bg-gradient-to-br from-pink-50 to-orange-50',
  tint: 'bg-gradient-to-r from-pink-50 via-purple-50 to-orange-50',
  accent: 'bg-gradient-to-br from-pink-200 via-purple-200 to-orange-200',
  transparent: ''
}

// 内边距样式映射
const paddingClasses = {
  none: '',
  sm: 'py-4',
  md: 'py-8',
  lg: 'py-12',
  xl: 'py-16'
}

// 最大宽度样式映射
const maxWidthMapping = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  '2xl': 'max-w-2xl',
  '3xl': 'max-w-3xl',
  '4xl': 'max-w-4xl',
  '5xl': 'max-w-5xl',
  '6xl': 'max-w-6xl',
  '7xl': 'max-w-7xl'
}

// 容器样式
const containerClasses = computed(() => {
  const classes = []

  // 背景
  if (backgroundClasses[props.background]) {
    classes.push(backgroundClasses[props.background])
  }

  // 内边距
  if (paddingClasses[props.padding]) {
    classes.push(paddingClasses[props.padding])
  }

  // 自定义类名
  if (props.customClass) {
    classes.push(props.customClass)
  }

  return classes.join(' ')
})

// 最大宽度容器样式
const maxWidthClasses = computed(() => {
  const classes = []

  // 最大宽度
  if (maxWidthMapping[props.maxWidth]) {
    classes.push(maxWidthMapping[props.maxWidth])
  }

  // 居中和水平内边距
  classes.push('mx-auto px-4 sm:px-6 lg:px-8')

  return classes.join(' ')
})
</script>
