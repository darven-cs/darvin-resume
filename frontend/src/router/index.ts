import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/editor/:id',
      name: 'editor',
      component: () => import('../views/EditorView.vue')
    }
  ]
})

export default router
