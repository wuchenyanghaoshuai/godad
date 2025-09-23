<template>
  <div
    :class="['rounded-full overflow-hidden relative select-none', wrapperClass]"
    :style="{ width: sizePx, height: sizePx }"
    :title="name || '用户'"
  >
    <img
      v-if="imgSrc"
      :src="imgSrc"
      class="w-full h-full object-cover"
      @error="handleError"
      :alt="name || 'avatar'"
    />
    <div v-else class="w-full h-full flex items-center justify-center bg-gray-200 text-gray-600">
      <span :style="{ fontSize: fontPx }" class="font-semibold tracking-wide">{{ initial }}</span>
    </div>
    <div
      v-if="showLetterOverlay"
      class="absolute inset-0 flex items-center justify-center pointer-events-none"
      :style="{ color: 'rgba(255,255,255,0.92)', textShadow: '0 1px 2px rgba(0,0,0,0.25)', fontSize: fontPx }"
    >
      <span class="font-semibold tracking-wide">{{ initial }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

// 使用你提供的10张图片（位于 ui/assets/default-avatars）
// 通过 new URL 相对路径导入，确保打包可用
const DEFAULTS: string[] = [
  new URL('../../assets/default-avatars/父与女共读.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/父子搭积木.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/父子乐.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/父子跑步.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/父子踢足球.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/妈妈给女儿梳头.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/母女读书.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/母女同乐.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/母女做饼干.jpg', import.meta.url).href,
  new URL('../../assets/default-avatars/喂宝宝吃饭.jpg', import.meta.url).href,
]

interface Props {
  avatar?: string
  name?: string
  size?: number // px
  wrapperClass?: string // 额外的类名（如阴影、边框）
}

const props = withDefaults(defineProps<Props>(), {
  avatar: '',
  name: 'U',
  size: 32,
  wrapperClass: ''
})

const sizePx = computed(() => `${props.size}px`)
const fontPx = computed(() => `${Math.max(12, Math.round(props.size * 0.45))}px`)

// 简单hash（djb2）
function hash(s: string): number {
  let h = 5381
  for (let i = 0; i < s.length; i++) {
    h = ((h << 5) + h) + s.charCodeAt(i)
    h |= 0
  }
  return Math.abs(h)
}

const initial = computed(() => (props.name || 'U').trim().charAt(0).toUpperCase())
const fallbackIdx = computed(() => hash((props.name || 'U').toLowerCase()) % DEFAULTS.length)
const fallbackSrc = computed(() => DEFAULTS[fallbackIdx.value])

const imgSrc = ref<string>('')
const isUsingFallback = ref<boolean>(false)

const computeSrc = () => {
  if (props.avatar && props.avatar.trim()) {
    imgSrc.value = props.avatar
    isUsingFallback.value = false
  } else {
    imgSrc.value = fallbackSrc.value
    isUsingFallback.value = true
  }
}

watch(() => props.avatar, computeSrc, { immediate: true })
watch(() => props.name, () => {
  if (!props.avatar) {
    imgSrc.value = fallbackSrc.value
    isUsingFallback.value = true
  }
})

const handleError = () => {
  // 头像加载失败，切换为默认
  imgSrc.value = fallbackSrc.value
  isUsingFallback.value = true
}

const showLetterOverlay = computed(() => isUsingFallback.value)
</script>

<style scoped>
</style>
