import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, clusterApi, resourceApi } from '@/api'
import type { ClusterInfo, ClusterHealth, KubernetesResource } from '@/types'

export const useK8sStore = defineStore('k8s', () => {
  const sessionId = ref<string | null>(localStorage.getItem('sessionId'))
  const clusterInfo = ref<ClusterInfo | null>(null)
  const clusterHealth = ref<ClusterHealth | null>(null)
  const namespaces = ref<KubernetesResource[]>([])
  const currentNamespace = ref<string>('')
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isConnected = computed(() => !!sessionId.value)

  async function connect(kubeconfig: string, contextName?: string) {
    loading.value = true
    error.value = null
    try {
      const response = await authApi.connect({ kubeconfig, contextName })
      sessionId.value = response.sessionId
      localStorage.setItem('sessionId', response.sessionId)
      await fetchClusterInfo()
      await fetchNamespaces()
      return true
    } catch (e: unknown) {
      const err = e as { response?: { data?: { message?: string } } }
      error.value = err.response?.data?.message || '连接失败'
      return false
    } finally {
      loading.value = false
    }
  }

  async function disconnect() {
    try {
      await authApi.disconnect()
    } finally {
      sessionId.value = null
      clusterInfo.value = null
      clusterHealth.value = null
      namespaces.value = []
      localStorage.removeItem('sessionId')
    }
  }

  async function fetchClusterInfo() {
    try {
      clusterInfo.value = await clusterApi.getInfo()
      clusterHealth.value = await clusterApi.getHealth()
    } catch (e: unknown) {
      console.error('Failed to fetch cluster info:', e)
    }
  }

  async function fetchNamespaces() {
    try {
      namespaces.value = await resourceApi.listNamespaces()
      if (namespaces.value.length > 0 && !currentNamespace.value) {
        currentNamespace.value = namespaces.value[0].metadata.name
      }
    } catch (e: unknown) {
      console.error('Failed to fetch namespaces:', e)
    }
  }

  function setNamespace(ns: string) {
    currentNamespace.value = ns
  }

  return {
    sessionId,
    clusterInfo,
    clusterHealth,
    namespaces,
    currentNamespace,
    loading,
    error,
    isConnected,
    connect,
    disconnect,
    fetchClusterInfo,
    fetchNamespaces,
    setNamespace
  }
}, {
  persist: {
    paths: ['sessionId']
  }
})
