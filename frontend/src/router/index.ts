import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'
import ResourceList from '@/views/ResourceList.vue'
import ResourceDetail from '@/views/ResourceDetail.vue'
import Logs from '@/views/Logs.vue'
import Events from '@/views/Events.vue'
import Terminal from '@/views/Terminal.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: { requiresAuth: true }
    },
    {
      path: '/resources/:kind',
      name: 'ResourceList',
      component: ResourceList,
      meta: { requiresAuth: true }
    },
    {
      path: '/resources/:kind/:namespace/:name',
      name: 'ResourceDetail',
      component: ResourceDetail,
      meta: { requiresAuth: true }
    },
    {
      path: '/logs',
      name: 'Logs',
      component: Logs,
      meta: { requiresAuth: true }
    },
    {
      path: '/events',
      name: 'Events',
      component: Events,
      meta: { requiresAuth: true }
    },
    {
      path: '/terminal',
      name: 'Terminal',
      component: Terminal,
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const sessionId = localStorage.getItem('sessionId')
  
  if (to.meta.requiresAuth && !sessionId) {
    next('/login')
  } else if (to.path === '/login' && sessionId) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
