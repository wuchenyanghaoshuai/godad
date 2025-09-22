<!--
  临时兼容组件 - 将逐步迁移到 AppLayout
  这个组件确保现有页面继续工作，同时我们逐步迁移到新的布局系统
-->
<template>
  <AppHeader
    :show-search="showSearch"
    :show-create-button="showCreateButton"
    :show-notifications="showNotifications"
    :show-user-points="showUserPoints"
    :show-navigation="showNavigation"
    :show-user-menu="showUserMenu"
  />
</template>

<script setup lang="ts">
import AppHeader from './layout/AppHeader.vue'

// Props - 与原 BaseHeader 兼容
withDefaults(defineProps<{
  showSearch?: boolean
  showCreateButton?: boolean
  showNotifications?: boolean
  showUserPoints?: boolean
  showNavigation?: boolean
  showUserMenu?: boolean
}>(), {
  showSearch: false,
  showCreateButton: false,
  showNotifications: false,
  showUserPoints: false,
  showNavigation: true,
  showUserMenu: true
})
</script>
