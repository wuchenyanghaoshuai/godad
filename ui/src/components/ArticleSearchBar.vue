<template>
  <div>
    <div class="relative">
      <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
      <input
        :value="internal"
        @input="onInput"
        @keydown.enter.prevent="onSubmit"
        :placeholder="placeholder || '输入关键词搜索文章...'"
        class="w-full pl-10 pr-20 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all text-sm sm:text-base placeholder-gray-500"
      />
      <button
        v-if="internal"
        @click="onClear"
        class="absolute right-16 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
        title="清空"
        aria-label="清空"
      >
        ×
      </button>
      <button
        @click="onSubmit"
        class="absolute right-2 top-1/2 -translate-y-1/2 px-3 py-1.5 bg-pink-600 text-white rounded-lg text-sm hover:bg-pink-700 transition-colors"
      >
        搜索
      </button>
    </div>

    <div v-if="suggestions && suggestions.length" class="mt-2 flex flex-wrap gap-2">
      <button
        v-for="(s, i) in suggestions.slice(0, 6)"
        :key="s + i"
        @click="applySuggestion(s)"
        class="px-2.5 py-1 text-xs bg-gray-100 text-gray-700 rounded-full hover:bg-gray-200 transition-colors"
      >
        #{{ s }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { SearchIcon } from 'lucide-vue-next'

interface Props {
  modelValue: string
  suggestions?: string[]
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  suggestions: () => [],
})

const emit = defineEmits<{
  'update:modelValue': [val: string]
  'submit': []
  'clear': []
}>()

const internal = ref(props.modelValue)
watch(() => props.modelValue, v => { internal.value = v })

const onInput = (e: Event) => {
  const v = (e.target as HTMLInputElement).value
  internal.value = v
  emit('update:modelValue', v)
}
const onSubmit = () => emit('submit')
const onClear = () => {
  internal.value = ''
  emit('update:modelValue', '')
  emit('clear')
}
const applySuggestion = (s: string) => {
  internal.value = s
  emit('update:modelValue', s)
  emit('submit')
}
</script>
