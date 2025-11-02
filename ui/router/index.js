import { createRouter, createWebHistory } from 'vue-router'
import { useSettingStore } from '@/stores/setting'
import { useAuthStore } from '@/stores/auth'

import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/install',
      name: 'install',
      component: () => import('@/views/InstallView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
    },
    {
      path: '',
      component: HomeView,
      meta: { requiresLogin: true },
      children: [
        {
          path: '',
          name: 'root',
          component: () => import('@/views/RootView.vue'),
        },
        {
          path: 'issues',
          name: 'issue-list',
          component: () => import('@/views/IssueListView.vue'),
        },
        {
          path: 'issues/:issueId',
          name: 'issue-details',
          component: () => import('@/views/IssueDetailsView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, from, next) => {
  const settingStore = useSettingStore()
  if (!settingStore.isInstalled && to.name !== 'install') {
    return next({ name: 'install' })
  }

  const authStore = useAuthStore()
  if (to.meta.requiresLogin && !authStore.isLoggedIn) {
    return next({ name: 'login' })
  }

  return next()
})

export default router
