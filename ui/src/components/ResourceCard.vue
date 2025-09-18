<template>
  <div class="p-4">
    <div class="flex items-stretch justify-between gap-4 rounded-xl bg-white border border-gray-200 hover:shadow-lg transition-shadow duration-200">
      <div class="flex flex-[2_2_0px] flex-col gap-4 p-6">
        <div class="flex flex-col gap-2">
          <h3 class="text-gray-900 text-lg font-bold leading-tight">
            {{ resource.title }}
          </h3>
          <p class="text-gray-600 text-sm leading-normal">
            {{ resource.description }}
          </p>
        </div>
        <button
          @click="$emit('click', resource)"
          class="flex min-w-[84px] max-w-[480px] cursor-pointer items-center justify-center overflow-hidden rounded-lg h-10 px-6 bg-gradient-to-r from-pink-500 to-rose-400 text-white text-sm font-medium leading-normal w-fit hover:from-pink-600 hover:to-rose-500 transition-all duration-200 shadow-md hover:shadow-lg transform hover:scale-105"
        >
          <span class="truncate">{{ resource.buttonText }}</span>
        </button>
      </div>
      <div class="flex-1 min-w-0">
        <div
          class="w-full h-full bg-center bg-no-repeat bg-cover rounded-xl min-h-[200px]"
          :style="{ backgroundImage: `url('${resource.image}')` }"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Props
interface Resource {
  id: number
  title: string
  description: string
  image: string
  type: 'e-book' | 'video' | 'tool'
  buttonText: string
  category: string
}

defineProps<{
  resource: Resource
}>()

// Emits
defineEmits<{
  click: [resource: Resource]
}>()
</script>