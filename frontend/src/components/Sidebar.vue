<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  LayoutDashboard, 
  Box, 
  Network, 
  HardDrive, 
  Settings, 
  Shield, 
  Server,
  FileText,
  Terminal,
  AlertCircle,
  ChevronLeft,
  ChevronRight,
  LogOut
} from 'lucide-vue-next'
import { useK8sStore } from '@/stores/k8s'
import { resourceCategories } from '@/types'

defineProps<{
  collapsed: boolean
}>()

const emit = defineEmits<{
  toggle: []
}>()

const router = useRouter()
const route = useRoute()
const k8sStore = useK8sStore()

const iconMap: Record<string, unknown> = {
  Box,
  Network,
  HardDrive,
  Settings,
  Shield,
  Server
}

const menuItems = computed(() => [
  {
    title: '仪表盘',
    icon: LayoutDashboard,
    path: '/dashboard'
  },
  ...resourceCategories.map(cat => ({
    title: cat.name,
    icon: iconMap[cat.icon] || Box,
    children: cat.kinds.map(kind => ({
      title: getResourceTitle(kind),
      path: `/resources/${kind}`
    }))
  })),
  {
    title: '日志',
    icon: FileText,
    path: '/logs'
  },
  {
    title: '事件',
    icon: AlertCircle,
    path: '/events'
  },
  {
    title: '终端',
    icon: Terminal,
    path: '/terminal'
  }
])

function getResourceTitle(kind: string): string {
  const titles: Record<string, string> = {
    pods: 'Pods',
    deployments: 'Deployments',
    statefulsets: 'StatefulSets',
    daemonsets: 'DaemonSets',
    replicasets: 'ReplicaSets',
    jobs: 'Jobs',
    cronjobs: 'CronJobs',
    services: 'Services',
    ingresses: 'Ingresses',
    endpoints: 'Endpoints',
    networkpolicies: 'NetworkPolicies',
    persistentvolumes: 'PersistentVolumes',
    persistentvolumeclaims: 'PVCs',
    storageclasses: 'StorageClasses',
    configmaps: 'ConfigMaps',
    secrets: 'Secrets',
    resourcequotas: 'ResourceQuotas',
    limitranges: 'LimitRanges',
    serviceaccounts: 'ServiceAccounts',
    roles: 'Roles',
    rolebindings: 'RoleBindings',
    clusterroles: 'ClusterRoles',
    clusterrolebindings: 'ClusterRoleBindings',
    nodes: 'Nodes',
    namespaces: 'Namespaces',
    crds: 'CRDs'
  }
  return titles[kind] || kind
}

function isActive(path: string): boolean {
  return route.path === path || route.path.startsWith(path + '/')
}

async function handleLogout() {
  await k8sStore.disconnect()
  router.push('/login')
}
</script>

<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-header">
      <div class="logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="currentColor" width="24" height="24">
            <path d="M13.95 13.5h-.23c-.18.11-.26.32-.18.5l.95 2.35c.95-.53 1.72-1.32 2.22-2.27l-2.36-.94c-.09-.04-.18-.06-.27-.06-.19 0-.36.1-.44.27l-.23.15zm-3.67.27c-.09-.17-.26-.27-.44-.27-.09 0-.18.02-.27.06l-2.36.94c.5.95 1.27 1.74 2.22 2.27l.95-2.35c.08-.18 0-.39-.18-.5l.08-.15zm2.22-3.22c.09.18.26.27.44.27.09 0 .18-.02.27-.06l2.36-.94c-.5-.95-1.27-1.74-2.22-2.27l-.95 2.35c-.08.18 0 .39.18.5l-.08.15zm-3.89-.73c.09.04.18.06.27.06.19 0 .36-.1.44-.27l.23-.15c.18-.11.26-.32.18-.5l-.95-2.35c-.95.53-1.72 1.32-2.22 2.27l2.36.94zm3.39-1.32c.27 0 .5-.22.5-.5V5.5c-.33-.04-.66-.06-1-.06s-.67.02-1 .06v2.77c0 .28.22.5.5.5h1zm2.5 3.5c0 .28.22.5.5.5h2.77c.04-.33.06-.66.06-1s-.02-.67-.06-1H15.5c-.28 0-.5.22-.5.5v1zm-5.5 0c0-.28-.22-.5-.5-.5H6.73c-.04.33-.06.66-.06 1s.02.67.06 1H9.5c.28 0 .5-.22.5-.5v-1zm4.5 3.5c0-.28-.22-.5-.5-.5h-1c-.28 0-.5.22-.5.5v2.77c.33.04.66.06 1 .06s.67-.02 1-.06v-2.77zM12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/>
          </svg>
        </div>
        <span v-if="!collapsed" class="logo-text">K8s Dashboard</span>
      </div>
      <button class="collapse-btn" @click="emit('toggle')">
        <ChevronLeft v-if="!collapsed" :size="16" />
        <ChevronRight v-else :size="16" />
      </button>
    </div>

    <nav class="sidebar-nav">
      <template v-for="item in menuItems" :key="item.title">
        <router-link
          v-if="'path' in item"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path as string) }"
        >
          <component :is="item.icon" :size="18" />
          <span v-if="!collapsed">{{ item.title }}</span>
        </router-link>
        <div v-else class="nav-group">
          <div class="nav-group-title" v-if="!collapsed">
            <component :is="item.icon" :size="16" />
            <span>{{ item.title }}</span>
          </div>
          <router-link
            v-for="child in item.children"
            :key="child.path"
            :to="child.path"
            class="nav-item sub-item"
            :class="{ active: isActive(child.path) }"
          >
            <span v-if="!collapsed">{{ child.title }}</span>
          </router-link>
        </div>
      </template>
    </nav>

    <div class="sidebar-footer">
      <button class="logout-btn" @click="handleLogout">
        <LogOut :size="18" />
        <span v-if="!collapsed">断开连接</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 240px;
  background-color: var(--color-bg-secondary);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  z-index: 100;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--color-border);
  height: 60px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.collapse-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
}

.collapse-btn:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
}

.nav-item:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.nav-item.active {
  background-color: var(--color-bg-tertiary);
  color: var(--color-accent);
  border-right: 2px solid var(--color-accent);
}

.nav-item.sub-item {
  padding-left: 32px;
  font-size: 13px;
}

.nav-group {
  margin: 4px 0;
}

.nav-group-title {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--color-border);
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 16px;
  border: none;
  background: transparent;
  color: var(--color-error);
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
}

.logout-btn:hover {
  background-color: rgba(248, 81, 73, 0.1);
}
</style>
