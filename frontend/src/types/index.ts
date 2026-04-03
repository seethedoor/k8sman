export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface ClusterInfo {
  name: string
  version: string
  nodeCount: number
  namespaceCount: number
  podCount: number
  deploymentCount: number
  serviceCount: number
}

export interface ClusterHealth {
  nodesHealthy: number
  nodesNotReady: number
  podsRunning: number
  podsPending: number
  podsFailed: number
  podsSucceeded: number
  podsUnknown: number
}

export interface ResourceMeta {
  name: string
  namespace: string
  uid: string
  creationTimestamp: string
  labels: Record<string, string>
  annotations: Record<string, string>
}

export interface KubernetesResource {
  apiVersion: string
  kind: string
  metadata: ResourceMeta
  spec?: Record<string, unknown>
  status?: Record<string, unknown>
}

export interface Event {
  type: 'Normal' | 'Warning'
  reason: string
  message: string
  objectKind: string
  objectName: string
  namespace: string
  count: number
  firstTimestamp: string
  lastTimestamp: string
}

export interface ConnectRequest {
  kubeconfig: string
  contextName?: string
}

export interface ConnectResponse {
  sessionId: string
}

export type ResourceKind = 
  | 'pods' | 'deployments' | 'statefulsets' | 'daemonsets' | 'replicasets' 
  | 'jobs' | 'cronjobs' | 'services' | 'ingresses' | 'endpoints' 
  | 'networkpolicies' | 'persistentvolumes' | 'persistentvolumeclaims' 
  | 'storageclasses' | 'configmaps' | 'secrets' | 'resourcequotas' 
  | 'limitranges' | 'serviceaccounts' | 'roles' | 'rolebindings' 
  | 'clusterroles' | 'clusterrolebindings' | 'nodes' | 'namespaces' | 'crds'

export interface ResourceCategory {
  name: string
  icon: string
  kinds: ResourceKind[]
}

export const resourceCategories: ResourceCategory[] = [
  {
    name: '工作负载',
    icon: 'Box',
    kinds: ['pods', 'deployments', 'statefulsets', 'daemonsets', 'replicasets', 'jobs', 'cronjobs']
  },
  {
    name: '网络',
    icon: 'Network',
    kinds: ['services', 'ingresses', 'endpoints', 'networkpolicies']
  },
  {
    name: '存储',
    icon: 'HardDrive',
    kinds: ['persistentvolumes', 'persistentvolumeclaims', 'storageclasses']
  },
  {
    name: '配置',
    icon: 'Settings',
    kinds: ['configmaps', 'secrets', 'resourcequotas', 'limitranges']
  },
  {
    name: '安全',
    icon: 'Shield',
    kinds: ['serviceaccounts', 'roles', 'rolebindings', 'clusterroles', 'clusterrolebindings']
  },
  {
    name: '集群',
    icon: 'Server',
    kinds: ['nodes', 'namespaces', 'crds']
  }
]
