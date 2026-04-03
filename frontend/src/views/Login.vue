<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Upload, Key } from 'lucide-vue-next'
import { useK8sStore } from '@/stores/k8s'
import { ElMessage } from 'element-plus'

const router = useRouter()
const k8sStore = useK8sStore()

const kubeconfig = ref('')
const contextName = ref('')
const loading = ref(false)
const activeTab = ref('file')

function handleFileUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      kubeconfig.value = e.target?.result as string
    }
    reader.readAsText(file)
  }
}

async function handleConnect() {
  if (!kubeconfig.value.trim()) {
    ElMessage.error('请输入或上传 kubeconfig 文件')
    return
  }

  loading.value = true
  try {
    const success = await k8sStore.connect(kubeconfig.value, contextName.value || undefined)
    if (success) {
      ElMessage.success('连接成功')
      router.push('/dashboard')
    } else {
      ElMessage.error(k8sStore.error || '连接失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <svg viewBox="0 0 24 24" fill="currentColor" width="48" height="48">
            <path d="M13.95 13.5h-.23c-.18.11-.26.32-.18.5l.95 2.35c.95-.53 1.72-1.32 2.22-2.27l-2.36-.94c-.09-.04-.18-.06-.27-.06-.19 0-.36.1-.44.27l-.23.15zm-3.67.27c-.09-.17-.26-.27-.44-.27-.09 0-.18.02-.27.06l-2.36.94c.5.95 1.27 1.74 2.22 2.27l.95-2.35c.08-.18 0-.39-.18-.5l.08-.15zm2.22-3.22c.09.18.26.27.44.27.09 0 .18-.02.27-.06l2.36-.94c-.5-.95-1.27-1.74-2.22-2.27l-.95 2.35c-.08.18 0 .39.18.5l-.08.15zm-3.89-.73c.09.04.18.06.27.06.19 0 .36-.1.44-.27l.23-.15c.18-.11.26-.32.18-.5l-.95-2.35c-.95.53-1.72 1.32-2.22 2.27l2.36.94zm3.39-1.32c.27 0 .5-.22.5-.5V5.5c-.33-.04-.66-.06-1-.06s-.67.02-1 .06v2.77c0 .28.22.5.5.5h1zm2.5 3.5c0 .28.22.5.5.5h2.77c.04-.33.06-.66.06-1s-.02-.67-.06-1H15.5c-.28 0-.5.22-.5.5v1zm-5.5 0c0-.28-.22-.5-.5-.5H6.73c-.04.33-.06.66-.06 1s.02.67.06 1H9.5c.28 0 .5-.22.5-.5v-1zm4.5 3.5c0-.28-.22-.5-.5-.5h-1c-.28 0-.5.22-.5.5v2.77c.33.04.66.06 1 .06s.67-.02 1-.06v-2.77zM12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/>
          </svg>
        </div>
        <h1>K8s Dashboard</h1>
        <p>连接到 Kubernetes 集群</p>
      </div>

      <el-tabs v-model="activeTab" class="login-tabs">
        <el-tab-pane label="上传文件" name="file">
          <div class="upload-area" @click="($refs.fileInput as HTMLInputElement).click()">
            <input
              ref="fileInput"
              type="file"
              accept=".yaml,.yml,.conf"
              @change="handleFileUpload"
              style="display: none"
            />
            <Upload :size="32" />
            <p>点击上传 kubeconfig 文件</p>
            <span>支持 .yaml, .yml, .conf 格式</span>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="粘贴内容" name="paste">
          <el-input
            v-model="kubeconfig"
            type="textarea"
            :rows="10"
            placeholder="粘贴 kubeconfig 内容..."
            class="kubeconfig-input"
          />
        </el-tab-pane>
      </el-tabs>

      <div class="context-select" v-if="kubeconfig">
        <el-input
          v-model="contextName"
          placeholder="Context 名称 (可选，默认使用当前 context)"
          clearable
        >
          <template #prepend>
            <Key :size="16" />
          </template>
        </el-input>
      </div>

      <el-button
        type="primary"
        size="large"
        :loading="loading"
        :disabled="!kubeconfig"
        @click="handleConnect"
        class="connect-btn"
      >
        连接集群
      </el-button>

      <div class="login-footer">
        <p>支持 Kubernetes 1.20+ 版本</p>
        <p>使用 K8s 原生认证，不存储凭证</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-bg-primary) 0%, var(--color-bg-secondary) 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 480px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
}

.login-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.login-tabs {
  margin-bottom: 20px;
}

.upload-area {
  border: 2px dashed var(--color-border);
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--color-text-secondary);
}

.upload-area:hover {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.upload-area p {
  margin: 12px 0 4px;
  font-size: 14px;
  color: var(--color-text-primary);
}

.upload-area span {
  font-size: 12px;
}

.kubeconfig-input :deep(textarea) {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
}

.context-select {
  margin-bottom: 20px;
}

.connect-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 500;
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

.login-footer p {
  font-size: 12px;
  color: var(--color-text-tertiary);
  margin: 4px 0;
}
</style>
