<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import Header from '@/components/Header.vue'
import { useK8sStore } from '@/stores/k8s'

const k8sStore = useK8sStore()
const sidebarCollapsed = ref(false)

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}
</script>

<template>
  <div class="app-container">
    <template v-if="k8sStore.isConnected">
      <Sidebar :collapsed="sidebarCollapsed" @toggle="toggleSidebar" />
      <div class="main-container" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
        <Header @toggle-sidebar="toggleSidebar" />
        <main class="content">
          <RouterView v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </RouterView>
        </main>
      </div>
    </template>
    <template v-else>
      <RouterView />
    </template>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 240px;
  transition: margin-left 0.3s ease;
  background-color: var(--color-bg-primary);
}

.main-container.sidebar-collapsed {
  margin-left: 64px;
}

.content {
  flex: 1;
  padding: 20px;
  overflow: auto;
}
</style>
