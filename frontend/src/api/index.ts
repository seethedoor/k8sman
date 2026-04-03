import axios from 'axios'
import type { ApiResponse, ClusterInfo, ClusterHealth, KubernetesResource, Event, ConnectRequest, ConnectResponse } from '@/types'

const API_BASE = '/api/v1'

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use((config) => {
  const sessionId = localStorage.getItem('sessionId')
  if (sessionId) {
    config.headers['X-Session-ID'] = sessionId
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('sessionId')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  connect: async (data: ConnectRequest): Promise<ConnectResponse> => {
    const response = await api.post<ApiResponse<ConnectResponse>>('/auth/connect', data)
    return response.data.data
  },
  
  disconnect: async (): Promise<void> => {
    await api.post('/auth/disconnect')
  }
}

export const clusterApi = {
  getInfo: async (): Promise<ClusterInfo> => {
    const response = await api.get<ApiResponse<ClusterInfo>>('/cluster/info')
    return response.data.data
  },
  
  getHealth: async (): Promise<ClusterHealth> => {
    const response = await api.get<ApiResponse<ClusterHealth>>('/cluster/health')
    return response.data.data
  }
}

export const resourceApi = {
  listNamespaces: async (): Promise<KubernetesResource[]> => {
    const response = await api.get<ApiResponse<KubernetesResource[]>>('/namespaces')
    return response.data.data
  },
  
  list: async (kind: string, namespace?: string): Promise<KubernetesResource[]> => {
    const params = namespace ? { namespace } : {}
    const response = await api.get<ApiResponse<KubernetesResource[]>>(`/resources/${kind}`, { params })
    return response.data.data
  },
  
  get: async (kind: string, namespace: string, name: string): Promise<KubernetesResource> => {
    const response = await api.get<ApiResponse<KubernetesResource>>(`/resources/${kind}/${namespace}/${name}`)
    return response.data.data
  },
  
  getYaml: async (kind: string, namespace: string, name: string): Promise<string> => {
    const response = await api.get<ApiResponse<{ yaml: string }>>(`/resources/${kind}/${namespace}/${name}/yaml`)
    return response.data.data.yaml
  },
  
  update: async (kind: string, namespace: string, name: string, data: Record<string, unknown>): Promise<KubernetesResource> => {
    const response = await api.put<ApiResponse<KubernetesResource>>(`/resources/${kind}/${namespace}/${name}`, data)
    return response.data.data
  },
  
  delete: async (kind: string, namespace: string, name: string): Promise<void> => {
    await api.delete(`/resources/${kind}/${namespace}/${name}`)
  }
}

export const eventApi = {
  list: async (namespace?: string, type?: string): Promise<Event[]> => {
    const params: Record<string, string> = {}
    if (namespace) params.namespace = namespace
    if (type) params.type = type
    const response = await api.get<ApiResponse<Event[]>>('/events', { params })
    return response.data.data
  }
}

export const logApi = {
  getWebSocketUrl: (namespace: string, podName: string, container?: string): string => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const baseUrl = `${protocol}//${window.location.host}${API_BASE}/pods/${namespace}/${podName}/logs`
    const params = new URLSearchParams()
    if (container) params.append('container', container)
    params.append('follow', 'true')
    params.append('tailLines', '100')
    return `${baseUrl}?${params.toString()}`
  }
}

export const execApi = {
  getWebSocketUrl: (namespace: string, podName: string, container?: string): string => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const baseUrl = `${protocol}//${window.location.host}${API_BASE}/pods/${namespace}/${podName}/exec`
    const params = new URLSearchParams()
    if (container) params.append('container', container)
    return `${baseUrl}?${params.toString()}`
  }
}

export default api
