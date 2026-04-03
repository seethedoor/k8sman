<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useK8sStore } from '@/stores/k8s'
import { eventApi } from '@/api'
import { RefreshCw, AlertTriangle, Info } from 'lucide-vue-next'
import type { Event } from '@/types'
import { ElMessage } from 'element-plus'

const k8sStore = useK8sStore()

const events = ref<Event[]>([])
const loading = ref(false)
const filterType = ref('')
const filterNamespace = ref('')

watch([filterType, filterNamespace, () => k8sStore.currentNamespace], () => {
  fetchEvents()
})

onMounted(() => {
  fetchEvents()
})

async function fetchEvents() {
  loading.value = true
  try {
    const ns = filterNamespace.value || undefined
    const type = filterType.value || undefined
    events.value = await eventApi.list(ns, type)
  } catch (e: unknown) {
    console.error('Failed to fetch events:', e)
    ElMessage.error('获取事件列表失败')
  } finally {
    loading.value = false
  }
}

function formatTime(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

function getTypeClass(type: string): string {
  return type === 'Warning' ? 'warning' : 'normal'
}
</script>

<template>
  <div class="events-page">
    <div class="controls">
      <div class="filter-group">
        <el-select
          v-model="filterType"
          placeholder="事件类型"
          clearable
          style="width: 150px"
        >
          <el-option label="Normal" value="Normal" />
          <el-option label="Warning" value="Warning" />
        </el-select>
        
        <el-select
          v-model="filterNamespace"
          placeholder="命名空间"
          clearable
          filterable
          style="width: 200px"
        >
          <el-option
            v-for="ns in k8sStore.namespaces"
            :key="ns.metadata.name"
            :label="ns.metadata.name"
            :value="ns.metadata.name"
          />
        </el-select>
      </div>
      
      <button class="refresh-btn" @click="fetchEvents" :disabled="loading">
        <RefreshCw :size="16" :class="{ spinning: loading }" />
        刷新
      </button>
    </div>

    <div class="table-container">
      <el-table :data="events" v-loading="loading" stripe>
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <div class="type-cell" :class="getTypeClass(row.type)">
              <AlertTriangle v-if="row.type === 'Warning'" :size="16" />
              <Info v-else :size="16" />
              <span>{{ row.type }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="对象" min-width="200">
          <template #default="{ row }">
            <div class="object-cell">
              <span class="object-kind">{{ row.objectKind }}</span>
              <span class="object-name">{{ row.objectName }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="命名空间" width="150">
          <template #default="{ row }">
            <span class="namespace-tag">{{ row.namespace || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="原因" width="150">
          <template #default="{ row }">
            <span class="reason-text">{{ row.reason }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="消息" min-width="300">
          <template #default="{ row }">
            <span class="message-text">{{ row.message }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="次数" width="80">
          <template #default="{ row }">
            <span class="count-text">{{ row.count }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="时间" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ formatTime(row.lastTimestamp) }}</span>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="footer">
      <span class="count-text">共 {{ events.length }} 个事件</span>
    </div>
  </div>
</template>

<style scoped>
.events-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-group {
  display: flex;
  gap: 12px;
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

.refresh-btn:hover:not(:disabled) {
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

.type-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
}

.type-cell.warning {
  color: var(--color-warning);
}

.type-cell.normal {
  color: var(--color-info);
}

.object-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.object-kind {
  padding: 2px 6px;
  background-color: var(--color-bg-tertiary);
  border-radius: 4px;
  font-size: 11px;
  color: var(--color-text-secondary);
}

.object-name {
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: var(--color-text-primary);
}

.namespace-tag {
  padding: 2px 8px;
  background-color: var(--color-bg-tertiary);
  border-radius: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.reason-text {
  font-weight: 500;
  color: var(--color-text-primary);
}

.message-text {
  font-size: 13px;
  color: var(--color-text-secondary);
  word-break: break-word;
}

.count-text {
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: var(--color-text-primary);
}

.time-text {
  font-size: 12px;
  color: var(--color-text-tertiary);
  font-family: 'JetBrains Mono', monospace;
}

.footer {
  display: flex;
  justify-content: flex-end;
}

.footer .count-text {
  font-size: 13px;
  color: var(--color-text-tertiary);
}
</style>
