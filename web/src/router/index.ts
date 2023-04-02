import { createRouter, createWebHistory } from 'vue-router'
import DashIndex from '@/views/dashboard/DashIndex.vue'

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
    }
  ]
})

export default router
