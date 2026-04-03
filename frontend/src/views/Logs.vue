<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useK8sStore } from '@/stores/k8s'
import { resourceApi, logApi } from '@/api'
import { Play, Pause, Download, Search, RefreshCw } from 'lucide-vue-next'
import type { KubernetesResource } from '@/types'
import { ElMessage } from 'element-plus'

const k8sStore = useK8sStore()

const pods = ref<KubernetesResource[]>([])
const selectedPod = ref('')
const selectedContainer = ref('')
const containers = ref<string[]>([])
const logs = ref<string[]>([])
const searchQuery = ref('')
const autoScroll = ref(true)
const follow = ref(true)
const loading = ref(false)
const ws = ref<WebSocket | null>(null)

const filteredLogs = computed(() => {
  if (!searchQuery.value) return logs.value
  return logs.value.filter(log => 
    log.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

watch(() => k8sStore.currentNamespace, () => {
  fetchPods()
})

watch(selectedPod, () => {
  fetchContainers()
})

onMounted(() => {
  fetchPods()
})

onUnmounted(() => {
  disconnectWebSocket()
})

async function fetchPods() {
  loading.value = true
  try {
    pods.value = await resourceApi.list('pods', k8sStore.currentNamespace)
    if (pods.value.length > 0 && !selectedPod.value) {
      selectedPod.value = pods.value[0].metadata.name
    }
  } catch (e: unknown) {
    console.error('Failed to fetch pods:', e)
    ElMessage.error('获取 Pod 列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchContainers() {
  if (!selectedPod.value) {
    containers.value = []
    return
  }
  
  try {
    const pod = await resourceApi.get('pods', k8sStore.currentNamespace, selectedPod.value)
    const spec = pod.spec as { containers: Array<{ name: string }> }
    containers.value = spec.containers.map(c => c.name)
    if (containers.value.length > 0 && !selectedContainer.value) {
      selectedContainer.value = containers.value[0]
    }
  } catch (e: unknown) {
    console.error('Failed to fetch containers:', e)
  }
}

function connectWebSocket() {
  if (!selectedPod.value) return
  
  disconnectWebSocket()
  logs.value = []
  
  const sessionId = localStorage.getItem('sessionId')
  const url = `${logApi.getWebSocketUrl(k8sStore.currentNamespace, selectedPod.value, selectedContainer.value)}&sessionId=${sessionId}`
  
  ws.value = new WebSocket(url)
  
  ws.value.onopen = () => {
    ElMessage.success('日志流已连接')
  }
  
  ws.value.onmessage = (event) => {
    logs.value.push(event.data)
    if (autoScroll.value) {
      scrollToBottom()
    }
  }
  
  ws.value.onerror = () => {
    ElMessage.error('日志流连接错误')
  }
  
  ws.value.onclose = () => {
    console.log('WebSocket closed')
  }
}

function disconnectWebSocket() {
  if (ws.value) {
    ws.value.close()
    ws.value = null
  }
}

function scrollToBottom() {
  const container = document.querySelector('.logs-content')
  if (container) {
    container.scrollTop = container.scrollHeight
  }
}

function toggleFollow() {
  follow.value = !follow.value
  if (follow.value) {
    connectWebSocket()
  } else {
    disconnectWebSocket()
  }
}

function downloadLogs() {
  const content = logs.value.join('\n')
  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${selectedPod.value}-${selectedContainer.value || 'all'}.log`
  a.click()
  URL.revokeObjectURL(url)
}

function clearLogs() {
  logs.value = []
}
</script>

<template>
  <div class="logs-page">
    <div class="controls">
      <div class="select-group">
        <el-select
          v-model="selectedPod"
          placeholder="选择 Pod"
          filterable
          style="width: 250px"
          @change="disconnectWebSocket"
        >
          <el-option
            v-for="pod in pods"
            :key="pod.metadata.name"
            :label="pod.metadata.name"
            :value="pod.metadata.name"
          />
        </el-select>
        
        <el-select
          v-model="selectedContainer"
          placeholder="选择容器"
          filterable
          style="width: 200px"
          @change="disconnectWebSocket"
        >
          <el-option
            v-for="container in containers"
            :key="container"
            :label="container"
            :value="container"
          />
        </el-select>
      </div>
      
      <div class="search-box">
        <Search :size="16" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索日志..."
          class="search-input"
        />
      </div>
      
      <div class="action-buttons">
        <button 
          class="btn" 
          :class="follow ? 'primary' : 'secondary'"
          @click="toggleFollow"
          :disabled="!selectedPod"
        >
          <Pause v-if="follow" :size="16" />
          <Play v-else :size="16" />
          {{ follow ? '暂停' : '开始' }}
        </button>
        
        <button class="btn secondary" @click="clearLogs">
          清空
        </button>
        
        <button class="btn secondary" @click="downloadLogs" :disabled="logs.length === 0">
          <Download :size="16" />
          下载
        </button>
        
        <el-checkbox v-model="autoScroll">自动滚动</el-checkbox>
      </div>
    </div>

    <div class="logs-container">
      <div class="logs-content">
        <div v-if="filteredLogs.length === 0" class="empty-state">
          <p v-if="!selectedPod">请选择一个 Pod 查看日志</p>
          <p v-else-if="!follow">点击"开始"按钮开始查看日志</p>
          <p v-else>等待日志...</p>
        </div>
        <div v-else class="logs-lines">
          <div
            v-for="(log, index) in filteredLogs"
            :key="index"
            class="log-line"
          >
            <span class="line-number">{{ index + 1 }}</span>
            <span class="log-text">{{ log }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.logs-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 100px);
}

.controls {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.select-group {
  display: flex;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background-color: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  flex: 1;
  max-width: 300px;
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

.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn.primary {
  background-color: var(--color-accent);
  color: white;
}

.btn.primary:hover:not(:disabled) {
  background-color: var(--color-accent-hover);
}

.btn.secondary {
  background-color: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
}

.btn.secondary:hover:not(:disabled) {
  background-color: var(--color-bg-secondary);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.logs-container {
  flex: 1;
  background-color: #0d1117;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
}

.logs-content {
  height: 100%;
  overflow: auto;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-tertiary);
}

.logs-lines {
  padding: 8px 0;
}

.log-line {
  display: flex;
  padding: 0 16px;
  background-color: transparent;
}

.log-line:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.line-number {
  width: 50px;
  color: #6e7681;
  text-align: right;
  padding-right: 16px;
  user-select: none;
}

.log-text {
  flex: 1;
  color: #e6edf3;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
