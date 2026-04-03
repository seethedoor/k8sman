<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { resourceApi } from '@/api'
import { useK8sStore } from '@/stores/k8s'
import { RefreshCw, Search, Eye } from 'lucide-vue-next'
import type { KubernetesResource } from '@/types'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const k8sStore = useK8sStore()

const resources = ref<KubernetesResource[]>([])
const loading = ref(false)
const searchQuery = ref('')

const kind = computed(() => route.params.kind as string)

const isClusterResource = computed(() => {
  const clusterKinds = ['nodes', 'namespaces', 'clusterroles', 'clusterrolebindings', 'persistentvolumes', 'storageclasses', 'crds']
  return clusterKinds.includes(kind.value)
})

const filteredResources = computed(() => {
  if (!searchQuery.value) return resources.value
  const query = searchQuery.value.toLowerCase()
  return resources.value.filter(r => 
    r.metadata.name.toLowerCase().includes(query) ||
    r.metadata.namespace?.toLowerCase().includes(query)
  )
})

watch([kind, () => k8sStore.currentNamespace], () => {
  fetchResources()
})

onMounted(() => {
  fetchResources()
})

async function fetchResources() {
  loading.value = true
  try {
    const ns = isClusterResource.value ? undefined : k8sStore.currentNamespace
    resources.value = await resourceApi.list(kind.value, ns)
  } catch (e: unknown) {
    console.error('Failed to fetch resources:', e)
    ElMessage.error('获取资源列表失败')
  } finally {
    loading.value = false
  }
}

function viewDetail(resource: KubernetesResource) {
  const ns = resource.metadata.namespace || '_'
  router.push(`/resources/${kind.value}/${ns}/${resource.metadata.name}`)
}

function getAge(timestamp: string): string {
  const now = new Date()
  const created = new Date(timestamp)
  const diff = now.getTime() - created.getTime()
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  if (days > 0) return `${days}d`
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  if (hours > 0) return `${hours}h`
  
  const minutes = Math.floor(diff / (1000 * 60))
  return `${minutes}m`
}

function getStatus(resource: KubernetesResource): { status: string; color: string } {
  const status = resource.status as { phase?: string; conditions?: Array<{ type: string; status: string }> } | undefined
  
  if (kind.value === 'pods') {
    const phase = status?.phase || 'Unknown'
    const colorMap: Record<string, string> = {
      Running: 'var(--color-success)',
      Pending: 'var(--color-warning)',
      Failed: 'var(--color-error)',
      Succeeded: 'var(--color-info)'
    }
    return { status: phase, color: colorMap[phase] || 'var(--color-text-secondary)' }
  }
  
  if (kind.value === 'nodes') {
    const conditions = status?.conditions || []
    const ready = conditions.find(c => c.type === 'Ready')
    if (ready?.status === 'True') {
      return { status: 'Ready', color: 'var(--color-success)' }
    }
    return { status: 'NotReady', color: 'var(--color-error)' }
  }
  
  if (kind.value === 'deployments' || kind.value === 'statefulsets' || kind.value === 'daemonsets') {
    const replicas = (resource.status as { replicas?: number; readyReplicas?: number }) || {}
    if (replicas.readyReplicas === replicas.replicas) {
      return { status: 'Ready', color: 'var(--color-success)' }
    }
    return { status: 'Updating', color: 'var(--color-warning)' }
  }
  
  return { status: 'Active', color: 'var(--color-success)' }
}
</script>

<template>
  <div class="resource-list">
    <div class="list-header">
      <div class="search-box">
        <Search :size="16" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索资源..."
          class="search-input"
        />
      </div>
      <button class="refresh-btn" @click="fetchResources" :disabled="loading">
        <RefreshCw :size="16" :class="{ spinning: loading }" />
        刷新
      </button>
    </div>

    <div class="table-container">
      <el-table :data="filteredResources" v-loading="loading" stripe>
        <el-table-column label="名称" min-width="200">
          <template #default="{ row }">
            <div class="name-cell">
              <span class="resource-name">{{ row.metadata.name }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column v-if="!isClusterResource" label="命名空间" width="150">
          <template #default="{ row }">
            <span class="namespace-tag">{{ row.metadata.namespace }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <span class="status-badge" :style="{ backgroundColor: getStatus(row).color + '20', color: getStatus(row).color }">
              {{ getStatus(row).status }}
            </span>
          </template>
        </el-table-column>
        
        <el-table-column label="创建时间" width="120">
          <template #default="{ row }">
            <span class="age-text">{{ getAge(row.metadata.creationTimestamp) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <button class="action-btn" @click="viewDetail(row)" title="查看详情">
              <Eye :size="16" />
            </button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="list-footer">
      <span class="count-text">共 {{ filteredResources.length }} 个资源</span>
    </div>
  </div>
</template>

<style scoped>
.resource-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  width: 300px;
}

.search-box svg {
  color: var(--color-text-tertiary);
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
}

.search-input::placeholder {
  color: var(--color-text-tertiary);
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  color: var(--color-text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn:hover {
  background-color: var(--color-bg-tertiary);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.table-container {
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.resource-name {
  font-weight: 500;
  color: var(--color-text-primary);
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
}

.namespace-tag {
  padding: 2px 8px;
  background-color: var(--color-bg-tertiary);
  border-radius: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.status-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.age-text {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-family: 'JetBrains Mono', monospace;
}

.action-btn {
  width: 32px;
  height: 32px;
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

.action-btn:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-accent);
}

.list-footer {
  display: flex;
  justify-content: flex-end;
}

.count-text {
  font-size: 13px;
  color: var(--color-text-tertiary);
}
</style>
