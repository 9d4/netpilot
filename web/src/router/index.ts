import { createRouter, createWebHistory } from 'vue-router'
import DashIndex from '@/views/dashboard/DashIndex.vue'
import { defineAsyncComponent } from 'vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/dashboard',
      component: DashIndex,
    },
    {
      path: '/dashboard/boards/:uuid?',
      component: defineAsyncComponent(() => import('@/views/dashboard/board/BoardIndex.vue'))
    }
  ]
})

export default router
