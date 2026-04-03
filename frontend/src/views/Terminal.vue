<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useK8sStore } from '@/stores/k8s'
import { resourceApi, execApi } from '@/api'
import { Terminal as TerminalIcon, X, Maximize2, Minimize2 } from 'lucide-vue-next'
import type { KubernetesResource } from '@/types'
import { ElMessage } from 'element-plus'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'

const k8sStore = useK8sStore()

const pods = ref<KubernetesResource[]>([])
const selectedPod = ref('')
const selectedContainer = ref('')
const containers = ref<string[]>([])
const loading = ref(false)
const terminalRef = ref<HTMLElement>()
const isFullscreen = ref(false)

let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null

watch(() => k8sStore.currentNamespace, () => {
  fetchPods()
})

watch(selectedPod, () => {
  fetchContainers()
})

onMounted(() => {
  fetchPods()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  disconnectTerminal()
  window.removeEventListener('resize', handleResize)
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

function initTerminal() {
  if (!terminalRef.value) return
  
  terminal = new Terminal({
    fontFamily: 'JetBrains Mono, monospace',
    fontSize: 14,
    theme: {
      background: '#0d1117',
      foreground: '#e6edf3',
      cursor: '#58a6ff',
      cursorAccent: '#0d1117',
      selectionBackground: 'rgba(88, 166, 255, 0.3)',
      black: '#0d1117',
      red: '#f85149',
      green: '#3fb950',
      yellow: '#d29922',
      blue: '#58a6ff',
      magenta: '#a371f7',
      cyan: '#39c5cf',
      white: '#e6edf3',
      brightBlack: '#6e7681',
      brightRed: '#f85149',
      brightGreen: '#3fb950',
      brightYellow: '#d29922',
      brightBlue: '#58a6ff',
      brightMagenta: '#a371f7',
      brightCyan: '#39c5cf',
      brightWhite: '#ffffff'
    },
    cursorBlink: true,
    cursorStyle: 'bar',
    scrollback: 10000
  })
  
  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.open(terminalRef.value)
  
  setTimeout(() => {
    fitAddon?.fit()
  }, 100)
  
  terminal.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  })
  
  terminal.onResize(({ cols, rows }) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ resize: { cols, rows } }))
    }
  })
}

function connectTerminal() {
  if (!selectedPod.value || !terminal) return
  
  disconnectTerminal()
  
  const sessionId = localStorage.getItem('sessionId')
  const url = `${execApi.getWebSocketUrl(k8sStore.currentNamespace, selectedPod.value, selectedContainer.value)}&sessionId=${sessionId}`
  
  ws = new WebSocket(url)
  
  ws.onopen = () => {
    terminal?.writeln('\x1b[32mConnected to terminal\x1b[0m')
    terminal?.writeln('')
    fitAddon?.fit()
  }
  
  ws.onmessage = (event) => {
    terminal?.write(event.data)
  }
  
  ws.onerror = () => {
    terminal?.writeln('\x1b[31mConnection error\x1b[0m')
    ElMessage.error('终端连接错误')
  }
  
  ws.onclose = () => {
    terminal?.writeln('\x1b[33mConnection closed\x1b[0m')
  }
}

function disconnectTerminal() {
  if (ws) {
    ws.close()
    ws = null
  }
}

function handleResize() {
  fitAddon?.fit()
}

function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value
  setTimeout(() => {
    fitAddon?.fit()
  }, 100)
}

function clearTerminal() {
  terminal?.clear()
}

function startSession() {
  initTerminal()
  connectTerminal()
}

function endSession() {
  disconnectTerminal()
  terminal?.dispose()
  terminal = null
}
</script>

<template>
  <div class="terminal-page" :class="{ fullscreen: isFullscreen }">
    <div class="controls" v-if="!isFullscreen">
      <div class="select-group">
        <el-select
          v-model="selectedPod"
          placeholder="选择 Pod"
          filterable
          style="width: 250px"
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
        />
      </div>
      
      <div class="action-buttons">
        <button 
          class="btn primary" 
          @click="startSession"
          :disabled="!selectedPod || !!terminal"
        >
          <TerminalIcon :size="16" />
          连接
        </button>
        
        <button 
          class="btn secondary" 
          @click="endSession"
          :disabled="!terminal"
        >
          <X :size="16" />
          断开
        </button>
      </div>
    </div>

    <div class="terminal-wrapper">
      <div class="terminal-header">
        <div class="header-left">
          <TerminalIcon :size="16" />
          <span v-if="selectedPod">{{ selectedPod }}</span>
          <span v-if="selectedContainer">/{{ selectedContainer }}</span>
        </div>
        <div class="header-actions">
          <button class="icon-btn" @click="clearTerminal" title="清空">
            清空
          </button>
          <button class="icon-btn" @click="toggleFullscreen">
            <Minimize2 v-if="isFullscreen" :size="16" />
            <Maximize2 v-else :size="16" />
          </button>
        </div>
      </div>
      
      <div class="terminal-container" ref="terminalRef">
        <div v-if="!terminal" class="empty-state">
          <TerminalIcon :size="48" />
          <p>选择 Pod 并点击"连接"开始终端会话</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.terminal-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 100px);
  gap: 16px;
}

.terminal-page.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  height: 100vh;
  z-index: 1000;
  background-color: var(--color-bg-primary);
  padding: 0;
}

.controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.select-group {
  display: flex;
  gap: 12px;
}

.action-buttons {
  display: flex;
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

.terminal-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #0d1117;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #161b22;
  border-bottom: 1px solid var(--color-border);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-secondary);
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  color: var(--color-text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.terminal-container {
  flex: 1;
  padding: 8px;
  overflow: hidden;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-tertiary);
  gap: 16px;
}

.empty-state p {
  font-size: 14px;
}
</style>
