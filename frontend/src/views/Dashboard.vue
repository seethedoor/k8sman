<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useK8sStore } from '@/stores/k8s'
import { 
  Server, 
  Box, 
  Network, 
  HardDrive, 
  AlertCircle,
  Activity,
  CheckCircle,
  XCircle,
  Clock
} from 'lucide-vue-next'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'

const k8sStore = useK8sStore()
const healthChart = ref<HTMLElement>()
const podsChart = ref<HTMLElement>()

const statsCards = computed(() => [
  { icon: Server, label: '节点', value: k8sStore.clusterInfo?.nodeCount || 0, color: '#58a6ff' },
  { icon: Box, label: 'Pods', value: k8sStore.clusterInfo?.podCount || 0, color: '#3fb950' },
  { icon: Network, label: 'Services', value: k8sStore.clusterInfo?.serviceCount || 0, color: '#a371f7' },
  { icon: HardDrive, label: 'Deployments', value: k8sStore.clusterInfo?.deploymentCount || 0, color: '#f0883e' }
])

const healthStats = computed(() => [
  { label: '健康节点', value: k8sStore.clusterHealth?.nodesHealthy || 0, icon: CheckCircle, color: 'var(--color-success)' },
  { label: '异常节点', value: k8sStore.clusterHealth?.nodesNotReady || 0, icon: XCircle, color: 'var(--color-error)' },
  { label: '运行中 Pod', value: k8sStore.clusterHealth?.podsRunning || 0, icon: Activity, color: 'var(--color-success)' },
  { label: '等待中 Pod', value: k8sStore.clusterHealth?.podsPending || 0, icon: Clock, color: 'var(--color-warning)' },
  { label: '失败 Pod', value: k8sStore.clusterHealth?.podsFailed || 0, icon: AlertCircle, color: 'var(--color-error)' }
])

onMounted(async () => {
  await k8sStore.fetchClusterInfo()
  await k8sStore.fetchNamespaces()
  
  setTimeout(() => {
    initCharts()
  }, 100)
})

function initCharts() {
  if (healthChart.value && k8sStore.clusterHealth) {
    const chart = echarts.init(healthChart.value)
    const option: EChartsOption = {
      tooltip: { trigger: 'item' },
      legend: { 
        bottom: '5%',
        left: 'center',
        textStyle: { color: 'var(--color-text-secondary)' }
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: 'var(--color-bg-secondary)',
          borderWidth: 2
        },
        label: { show: false },
        emphasis: {
          label: { show: true, fontSize: 14, fontWeight: 'bold' }
        },
        data: [
          { value: k8sStore.clusterHealth.nodesHealthy, name: '健康节点', itemStyle: { color: '#3fb950' } },
          { value: k8sStore.clusterHealth.nodesNotReady, name: '异常节点', itemStyle: { color: '#f85149' } }
        ]
      }]
    }
    chart.setOption(option)
  }

  if (podsChart.value && k8sStore.clusterHealth) {
    const chart = echarts.init(podsChart.value)
    const option: EChartsOption = {
      tooltip: { trigger: 'item' },
      legend: { 
        bottom: '5%',
        left: 'center',
        textStyle: { color: 'var(--color-text-secondary)' }
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: 'var(--color-bg-secondary)',
          borderWidth: 2
        },
        label: { show: false },
        emphasis: {
          label: { show: true, fontSize: 14, fontWeight: 'bold' }
        },
        data: [
          { value: k8sStore.clusterHealth.podsRunning, name: 'Running', itemStyle: { color: '#3fb950' } },
          { value: k8sStore.clusterHealth.podsPending, name: 'Pending', itemStyle: { color: '#d29922' } },
          { value: k8sStore.clusterHealth.podsFailed, name: 'Failed', itemStyle: { color: '#f85149' } },
          { value: k8sStore.clusterHealth.podsSucceeded, name: 'Succeeded', itemStyle: { color: '#58a6ff' } }
        ]
      }]
    }
    chart.setOption(option)
  }
}
</script>

<template>
  <div class="dashboard">
    <div class="stats-grid">
      <div v-for="stat in statsCards" :key="stat.label" class="stat-card">
        <div class="stat-icon" :style="{ backgroundColor: stat.color + '20', color: stat.color }">
          <component :is="stat.icon" :size="24" />
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-card">
        <h3>节点状态</h3>
        <div ref="healthChart" class="chart"></div>
      </div>
      <div class="chart-card">
        <h3>Pod 状态分布</h3>
        <div ref="podsChart" class="chart"></div>
      </div>
    </div>

    <div class="health-grid">
      <div v-for="item in healthStats" :key="item.label" class="health-item">
        <component :is="item.icon" :size="18" :style="{ color: item.color }" />
        <span class="health-label">{{ item.label }}</span>
        <span class="health-value">{{ item.value }}</span>
      </div>
    </div>

    <div class="info-card">
      <h3>集群信息</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">版本</span>
          <span class="info-value mono">{{ k8sStore.clusterInfo?.version || '-' }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">命名空间</span>
          <span class="info-value">{{ k8sStore.clusterInfo?.namespaceCount || 0 }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.charts-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.chart-card {
  padding: 20px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.chart-card h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 16px;
}

.chart {
  height: 250px;
}

.health-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}

.health-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.health-label {
  flex: 1;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.health-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.info-card {
  padding: 20px;
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.info-card h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.info-value {
  font-size: 14px;
  color: var(--color-text-primary);
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .health-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
  
  .health-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
