<template>
  <span 
    :class="[
      'tag-component',
      `tag-${size}`,
      { 
        'tag-clickable': clickable,
        'tag-selected': selected,
        'tag-removable': removable
      }
    ]"
    :style="{ 
      backgroundColor: backgroundColor,
      color: textColor,
      borderColor: borderColor
    }"
    @click="handleClick"
  >
    <span class="tag-text">{{ tag.name }}</span>
    <span v-if="showCount && tag.usage_count > 0" class="tag-count">
      {{ formatCount(tag.usage_count) }}
    </span>
    <button 
      v-if="removable" 
      class="tag-remove"
      @click.stop="handleRemove"
      :style="{ color: textColor }"
    >
      ×
    </button>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Tag } from '@/api/tags'

interface Props {
  tag: Tag
  size?: 'small' | 'medium' | 'large'
  clickable?: boolean
  selected?: boolean
  removable?: boolean
  showCount?: boolean
  variant?: 'solid' | 'outline' | 'ghost'
}

interface Emits {
  (e: 'click', tag: Tag): void
  (e: 'remove', tag: Tag): void
}

const props = withDefaults(defineProps<Props>(), {
  size: 'medium',
  clickable: false,
  selected: false,
  removable: false,
  showCount: false,
  variant: 'solid'
})

const emit = defineEmits<Emits>()

const backgroundColor = computed(() => {
  if (props.variant === 'outline') {
    return 'transparent'
  }
  if (props.variant === 'ghost') {
    return `${props.tag.color}20`
  }
  return props.selected ? '#3B82F6' : props.tag.color || '#3B82F6'
})

const textColor = computed(() => {
  if (props.variant === 'outline' || props.variant === 'ghost') {
    return props.tag.color || '#3B82F6'
  }
  return '#ffffff'
})

const borderColor = computed(() => {
  if (props.variant === 'outline') {
    return props.tag.color || '#3B82F6'
  }
  return 'transparent'
})

const handleClick = () => {
  if (props.clickable) {
    emit('click', props.tag)
  }
}

const handleRemove = () => {
  emit('remove', props.tag)
}

const formatCount = (count: number): string => {
  if (count < 1000) return count.toString()
  if (count < 10000) return `${Math.floor(count / 1000)}k`
  return `${Math.floor(count / 10000)}w`
}
</script>

<style scoped>
.tag-component {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border-radius: 12px;
  font-weight: 500;
  border: 1px solid;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.tag-small {
  padding: 2px 8px;
  font-size: 12px;
  line-height: 16px;
}

.tag-medium {
  padding: 4px 12px;
  font-size: 14px;
  line-height: 20px;
}

.tag-large {
  padding: 6px 16px;
  font-size: 16px;
  line-height: 24px;
}

.tag-clickable {
  cursor: pointer;
}

.tag-clickable:hover {
  opacity: 0.8;
  transform: translateY(-1px);
}

.tag-selected {
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3);
}

.tag-text {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tag-count {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 1px 6px;
  font-size: 12px;
  line-height: 16px;
}

.tag-remove {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 0;
  margin: 0;
  transition: opacity 0.2s ease;
}

.tag-remove:hover {
  opacity: 0.7;
}
</style>