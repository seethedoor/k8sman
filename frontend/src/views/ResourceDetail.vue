<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { resourceApi } from '@/api'
import { useK8sStore } from '@/stores/k8s'
import { ArrowLeft, Save, Trash2, RefreshCw } from 'lucide-vue-next'
import type { KubernetesResource } from '@/types'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as monaco from 'monaco-editor'

const route = useRoute()
const router = useRouter()
const k8sStore = useK8sStore()

const resource = ref<KubernetesResource | null>(null)
const yamlContent = ref('')
const originalYaml = ref('')
const loading = ref(false)
const saving = ref(false)
const editorContainer = ref<HTMLElement>()
let editor: monaco.editor.IStandaloneCodeEditor | null = null

const kind = computed(() => route.params.kind as string)
const namespace = computed(() => route.params.namespace as string)
const name = computed(() => route.params.name as string)

const isClusterResource = computed(() => {
  return namespace.value === '_'
})

const hasChanges = computed(() => {
  return yamlContent.value !== originalYaml.value
})

watch([kind, namespace, name], () => {
  fetchResource()
})

onMounted(async () => {
  await fetchResource()
  initEditor()
})

async function fetchResource() {
  loading.value = true
  try {
    const ns = isClusterResource.value ? '' : namespace.value
    resource.value = await resourceApi.get(kind.value, ns, name.value)
    const yaml = await resourceApi.getYaml(kind.value, ns, name.value)
    yamlContent.value = yaml
    originalYaml.value = yaml
    if (editor) {
      editor.setValue(yaml)
    }
  } catch (e: unknown) {
    console.error('Failed to fetch resource:', e)
    ElMessage.error('获取资源详情失败')
  } finally {
    loading.value = false
  }
}

function initEditor() {
  if (!editorContainer.value) return
  
  const darkTheme = document.documentElement.classList.contains('dark')
  
  editor = monaco.editor.create(editorContainer.value, {
    value: yamlContent.value,
    language: 'yaml',
    theme: darkTheme ? 'vs-dark' : 'vs',
    minimap: { enabled: false },
    fontSize: 13,
    fontFamily: 'JetBrains Mono, monospace',
    lineNumbers: 'on',
    scrollBeyondLastLine: false,
    automaticLayout: true,
    tabSize: 2
  })
  
  editor.onDidChangeModelContent(() => {
    yamlContent.value = editor?.getValue() || ''
  })
}

async function handleSave() {
  if (!hasChanges.value) return
  
  saving.value = true
  try {
    const yaml = await import('js-yaml')
    const data = yaml.load(yamlContent.value) as Record<string, unknown>
    const ns = isClusterResource.value ? '' : namespace.value
    await resourceApi.update(kind.value, ns, name.value, data)
    ElMessage.success('资源已更新')
    originalYaml.value = yamlContent.value
    await fetchResource()
  } catch (e: unknown) {
    console.error('Failed to update resource:', e)
    ElMessage.error('更新资源失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${name.value} 吗？此操作不可撤销。`,
      '删除确认',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ns = isClusterResource.value ? '' : namespace.value
    await resourceApi.delete(kind.value, ns, name.value)
    ElMessage.success('资源已删除')
    router.back()
  } catch (e: unknown) {
    if (e !== 'cancel') {
      console.error('Failed to delete resource:', e)
      ElMessage.error('删除资源失败')
    }
  }
}

function goBack() {
  router.back()
}
</script>

<template>
  <div class="resource-detail">
    <div class="detail-header">
      <div class="header-left">
        <button class="back-btn" @click="goBack">
          <ArrowLeft :size="18" />
        </button>
        <div class="title-info">
          <h2>{{ name }}</h2>
          <span v-if="!isClusterResource" class="namespace-badge">{{ namespace }}</span>
          <span class="kind-badge">{{ kind }}</span>
        </div>
      </div>
      <div class="header-actions">
        <button class="action-btn secondary" @click="fetchResource" :disabled="loading">
          <RefreshCw :size="16" :class="{ spinning: loading }" />
          刷新
        </button>
        <button class="action-btn danger" @click="handleDelete">
          <Trash2 :size="16" />
          删除
        </button>
        <button 
          class="action-btn primary" 
          @click="handleSave" 
          :disabled="!hasChanges || saving"
        >
          <Save :size="16" />
          保存
        </button>
      </div>
    </div>

    <div class="editor-container" ref="editorContainer" v-loading="loading"></div>

    <div class="change-indicator" v-if="hasChanges">
      <span class="dot"></span>
      有未保存的更改
    </div>
  </div>
</template>

<style scoped>
.resource-detail {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 100px);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--color-border);
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background-color: var(--color-bg-secondary);
  border-radius: 6px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.title-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-info h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.namespace-badge,
.kind-badge {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.namespace-badge {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
}

.kind-badge {
  background-color: var(--color-accent);
  color: white;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
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

.action-btn.primary {
  background-color: var(--color-accent);
  color: white;
}

.action-btn.primary:hover:not(:disabled) {
  background-color: var(--color-accent-hover);
}

.action-btn.primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.secondary {
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
}

.action-btn.secondary:hover:not(:disabled) {
  background-color: var(--color-bg-tertiary);
}

.action-btn.danger {
  background-color: transparent;
  border: 1px solid var(--color-error);
  color: var(--color-error);
}

.action-btn.danger:hover {
  background-color: var(--color-error);
  color: white;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.editor-container {
  flex: 1;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
  background-color: var(--color-bg-secondary);
}

.change-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background-color: var(--color-warning);
  color: #000;
  border-radius: 6px;
  margin-top: 16px;
  font-size: 14px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #000;
}
</style>
