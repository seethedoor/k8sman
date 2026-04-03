<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Menu, Sun, Moon, Bell } from 'lucide-vue-next'
import { useK8sStore } from '@/stores/k8s'
import { useThemeStore } from '@/stores/theme'

const emit = defineEmits<{
  'toggle-sidebar': []
}>()

const route = useRoute()
const k8sStore = useK8sStore()
const themeStore = useThemeStore()

const pageTitle = computed(() => {
  const path = route.path
  if (path === '/dashboard') return '仪表盘'
  if (path === '/logs') return '日志查看'
  if (path === '/events') return '事件列表'
  if (path === '/terminal') return '终端'
  if (path.startsWith('/resources/')) {
    const parts = path.split('/')
    if (parts.length >= 4) {
      return parts[3] || parts[2]
    }
    return parts[2]
  }
  return 'K8s Dashboard'
})

const clusterName = computed(() => {
  return k8sStore.clusterInfo?.name || '未连接'
})
</script>

<template>
  <header class="header">
    <div class="header-left">
      <button class="menu-btn" @click="emit('toggle-sidebar')">
        <Menu :size="20" />
      </button>
      <h1 class="page-title">{{ pageTitle }}</h1>
    </div>

    <div class="header-right">
      <div class="cluster-info">
        <span class="cluster-label">集群:</span>
        <span class="cluster-name">{{ clusterName }}</span>
      </div>

      <el-select
        v-model="k8sStore.currentNamespace"
        placeholder="选择命名空间"
        size="small"
        style="width: 160px"
        @change="k8sStore.setNamespace"
      >
        <el-option
          v-for="ns in k8sStore.namespaces"
          :key="ns.metadata.name"
          :label="ns.metadata.name"
          :value="ns.metadata.name"
        />
      </el-select>

      <button class="icon-btn" title="通知">
        <Bell :size="18" />
      </button>

      <button class="icon-btn" @click="themeStore.toggleTheme" :title="themeStore.theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'">
        <Sun v-if="themeStore.theme === 'dark'" :size="18" />
        <Moon v-else :size="18" />
      </button>
    </div>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 20px;
  background-color: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.menu-btn {
  display: none;
  width: 36px;
  height: 36px;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
}

.menu-btn:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.cluster-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: var(--color-bg-tertiary);
  border-radius: 6px;
}

.cluster-label {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.cluster-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.icon-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

@media (max-width: 768px) {
  .menu-btn {
    display: flex;
  }
  
  .cluster-info {
    display: none;
  }
}
</style>
